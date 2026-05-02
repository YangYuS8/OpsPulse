package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"opspulse/server/internal/core"
)

type Store struct {
	db             *sql.DB
	offlineTimeout time.Duration
}

func Open(path string, offlineTimeout time.Duration) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	store := &Store{db: db, offlineTimeout: offlineTimeout}
	if err := store.migrate(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS nodes (
		node_id TEXT PRIMARY KEY,
		hostname TEXT NOT NULL,
		last_seen TEXT NOT NULL,
		status TEXT NOT NULL,
		uptime INTEGER NOT NULL,
		cpu_usage REAL NOT NULL,
		memory_json TEXT NOT NULL,
		disk_json TEXT NOT NULL,
		load_json TEXT NOT NULL,
		docker_json TEXT NOT NULL,
		services_json TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		node_id TEXT NOT NULL,
		level TEXT NOT NULL,
		message TEXT NOT NULL,
		created_at TEXT NOT NULL
	);
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *Store) UpsertHeartbeat(ctx context.Context, report core.NodeReport) (core.NodeRecord, core.Event, error) {
	status, summary, level := deriveNodeStatus(report.Metrics)
	lastSeen := report.SentAt.UTC()
	memoryJSON, _ := json.Marshal(report.Metrics.Memory)
	diskJSON, _ := json.Marshal(report.Metrics.Disk)
	loadJSON, _ := json.Marshal(report.Metrics.Load)
	dockerJSON, _ := json.Marshal(report.Metrics.Docker)
	servicesJSON, _ := json.Marshal(report.Metrics.Services)

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO nodes (node_id, hostname, last_seen, status, uptime, cpu_usage, memory_json, disk_json, load_json, docker_json, services_json)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(node_id) DO UPDATE SET
			hostname = excluded.hostname,
			last_seen = excluded.last_seen,
			status = excluded.status,
			uptime = excluded.uptime,
			cpu_usage = excluded.cpu_usage,
			memory_json = excluded.memory_json,
			disk_json = excluded.disk_json,
			load_json = excluded.load_json,
			docker_json = excluded.docker_json,
			services_json = excluded.services_json
	`, report.NodeID, report.Metrics.Hostname, lastSeen.Format(time.RFC3339), status, report.Metrics.Uptime, report.Metrics.CPUUsage, string(memoryJSON), string(diskJSON), string(loadJSON), string(dockerJSON), string(servicesJSON))
	if err != nil {
		return core.NodeRecord{}, core.Event{}, err
	}

	event := core.Event{
		NodeID:    report.NodeID,
		Level:     level,
		Message:   summary,
		CreatedAt: lastSeen,
	}
	res, err := s.db.ExecContext(ctx, `INSERT INTO events (node_id, level, message, created_at) VALUES (?, ?, ?, ?)`, event.NodeID, event.Level, event.Message, event.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return core.NodeRecord{}, core.Event{}, err
	}
	event.ID, _ = res.LastInsertId()

	node, err := s.GetNode(ctx, report.NodeID)
	return node, event, err
}

func (s *Store) GetOverview(ctx context.Context) (core.Overview, error) {
	nodes, err := s.ListNodes(ctx)
	if err != nil {
		return core.Overview{}, err
	}
	var overview core.Overview
	for _, node := range nodes {
		overview.NodesTotal++
		switch node.Status {
		case core.NodeStatusOnline:
			overview.NodesOnline++
		case core.NodeStatusWarn:
			overview.NodesWarn++
		default:
			overview.NodesOffline++
		}
		overview.ContainersUp += node.Docker.Running
		for _, svc := range node.Services {
			if svc.Status != "healthy" {
				overview.ServicesDown++
			}
		}
	}
	return overview, nil
}

func (s *Store) HasNodes(ctx context.Context) (bool, error) {
	row := s.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM nodes`)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Store) ListNodes(ctx context.Context) ([]core.NodeRecord, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT node_id, hostname, last_seen, status, uptime, cpu_usage, memory_json, disk_json, load_json, docker_json, services_json FROM nodes ORDER BY hostname ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	nodes := make([]core.NodeRecord, 0)
	for rows.Next() {
		node, err := scanNode(rows, s.offlineTimeout)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, rows.Err()
}

func (s *Store) GetNode(ctx context.Context, nodeID string) (core.NodeRecord, error) {
	row := s.db.QueryRowContext(ctx, `SELECT node_id, hostname, last_seen, status, uptime, cpu_usage, memory_json, disk_json, load_json, docker_json, services_json FROM nodes WHERE node_id = ?`, nodeID)
	return scanNode(row, s.offlineTimeout)
}

func (s *Store) ListEvents(ctx context.Context, limit int) ([]core.Event, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, node_id, level, message, created_at FROM events ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := make([]core.Event, 0)
	for rows.Next() {
		var event core.Event
		var createdAt string
		if err := rows.Scan(&event.ID, &event.NodeID, &event.Level, &event.Message, &createdAt); err != nil {
			return nil, err
		}
		event.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

type scanner interface{ Scan(dest ...any) error }

func scanNode(row scanner, offlineTimeout time.Duration) (core.NodeRecord, error) {
	var node core.NodeRecord
	var lastSeenRaw, statusRaw string
	var memoryJSON, diskJSON, loadJSON, dockerJSON, servicesJSON string
	err := row.Scan(&node.NodeID, &node.Hostname, &lastSeenRaw, &statusRaw, &node.Uptime, &node.CPUUsage, &memoryJSON, &diskJSON, &loadJSON, &dockerJSON, &servicesJSON)
	if err != nil {
		return core.NodeRecord{}, err
	}
	node.LastSeen, err = time.Parse(time.RFC3339, lastSeenRaw)
	if err != nil {
		return core.NodeRecord{}, err
	}
	node.Status = core.NodeStatus(statusRaw)
	if err := json.Unmarshal([]byte(memoryJSON), &node.Memory); err != nil {
		return core.NodeRecord{}, err
	}
	if err := json.Unmarshal([]byte(diskJSON), &node.Disk); err != nil {
		return core.NodeRecord{}, err
	}
	if err := json.Unmarshal([]byte(loadJSON), &node.Load); err != nil {
		return core.NodeRecord{}, err
	}
	if err := json.Unmarshal([]byte(dockerJSON), &node.Docker); err != nil {
		return core.NodeRecord{}, err
	}
	if err := json.Unmarshal([]byte(servicesJSON), &node.Services); err != nil {
		return core.NodeRecord{}, err
	}
	node.HeartbeatAge = int64(time.Since(node.LastSeen).Seconds())
	node.Offline = time.Since(node.LastSeen) > offlineTimeout
	if node.Offline {
		node.Status = core.NodeStatusOffline
		node.StatusSummary = fmt.Sprintf("Node offline for %ds", node.HeartbeatAge)
		return node, nil
	}
	_, summary, _ := deriveNodeStatus(core.NodeMetrics{CPUUsage: node.CPUUsage, Memory: node.Memory, Disk: node.Disk, Services: node.Services})
	node.StatusSummary = summary
	return node, nil
}

func deriveNodeStatus(metrics core.NodeMetrics) (core.NodeStatus, string, string) {
	problems := 0
	for _, svc := range metrics.Services {
		if svc.Status != "healthy" {
			problems++
		}
	}
	if metrics.CPUUsage >= 90 || metrics.Memory.Usage >= 92 || metrics.Disk.Usage >= 92 || problems > 0 {
		return core.NodeStatusWarn, fmt.Sprintf("%s reporting stress: cpu %.1f%%, mem %.1f%%, disk %.1f%%, services %d unhealthy", metrics.Hostname, metrics.CPUUsage, metrics.Memory.Usage, metrics.Disk.Usage, problems), "warn"
	}
	return core.NodeStatusOnline, fmt.Sprintf("%s heartbeat ok", metrics.Hostname), "info"
}
