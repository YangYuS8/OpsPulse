package collect

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type ServiceStatus struct {
	Name   string `json:"name"`
	Active string `json:"active"`
	Sub    string `json:"sub"`
	Status string `json:"status"`
}

type ContainerStatus struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	State  string `json:"state"`
	Status string `json:"status"`
}

type CheckConfig struct {
	Name           string
	Type           string
	Target         string
	Timeout        time.Duration
	ExpectedStatus int
}

type CheckStatus struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Target    string `json:"target"`
	Status    string `json:"status"`
	LatencyMS int64  `json:"latencyMs"`
	Error     string `json:"error"`
}

type UsageMetric struct {
	Used  uint64  `json:"used"`
	Total uint64  `json:"total"`
	Usage float64 `json:"usage"`
}

type LoadAverage struct {
	One     float64 `json:"one"`
	Five    float64 `json:"five"`
	Fifteen float64 `json:"fifteen"`
}

type DockerMetric struct {
	Running    int               `json:"running"`
	Exited     int               `json:"exited"`
	Containers []ContainerStatus `json:"containers"`
}

type Snapshot struct {
	Hostname string          `json:"hostname"`
	Uptime   int64           `json:"uptime"`
	CPUUsage float64         `json:"cpuUsage"`
	Memory   UsageMetric     `json:"memory"`
	Disk     UsageMetric     `json:"disk"`
	Load     LoadAverage     `json:"load"`
	Docker   DockerMetric    `json:"docker"`
	Services []ServiceStatus `json:"services"`
	Checks   []CheckStatus   `json:"checks"`
}

type Collector struct {
	services      []string
	checks        []CheckConfig
	dockerEnabled bool
	prevIdle      uint64
	prevTotal     uint64
	hasPrevCPU    bool
}

func New(serviceWhitelist []string, checks []CheckConfig, dockerEnabled bool) *Collector {
	return &Collector{services: serviceWhitelist, checks: checks, dockerEnabled: dockerEnabled}
}

func (c *Collector) Snapshot() (Snapshot, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return Snapshot{}, err
	}
	uptime, _ := readUptime()
	cpuUsage, _ := c.readCPUUsage()
	memory, _ := readMemoryUsage()
	disk, _ := readDiskUsage("/")
	load, _ := readLoadAverage()
	docker, _ := c.readDockerMetric()
	services := c.readServices()
	checks := c.runChecks()
	return Snapshot{
		Hostname: hostname,
		Uptime:   uptime,
		CPUUsage: cpuUsage,
		Memory:   memory,
		Disk:     disk,
		Load:     load,
		Docker:   docker,
		Services: services,
		Checks:   checks,
	}, nil
}

func readUptime() (int64, error) {
	content, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	parts := strings.Fields(string(content))
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid uptime")
	}
	value, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, err
	}
	return int64(value), nil
}

func (c *Collector) readCPUUsage() (float64, error) {
	content, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, err
	}
	line := strings.SplitN(string(content), "\n", 2)[0]
	parts := strings.Fields(line)
	if len(parts) < 5 {
		return 0, fmt.Errorf("invalid cpu stat")
	}
	values := make([]uint64, 0, len(parts)-1)
	for _, part := range parts[1:] {
		value, err := strconv.ParseUint(part, 10, 64)
		if err != nil {
			return 0, err
		}
		values = append(values, value)
	}
	var idle, total uint64
	idle = values[3]
	for _, value := range values {
		total += value
	}
	if !c.hasPrevCPU {
		c.prevIdle, c.prevTotal, c.hasPrevCPU = idle, total, true
		return 0, nil
	}
	idleDelta := idle - c.prevIdle
	totalDelta := total - c.prevTotal
	c.prevIdle, c.prevTotal = idle, total
	if totalDelta == 0 {
		return 0, nil
	}
	return 100 * (1 - float64(idleDelta)/float64(totalDelta)), nil
}

func readMemoryUsage() (UsageMetric, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return UsageMetric{}, err
	}
	defer file.Close()
	var total, available uint64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 2 {
			continue
		}
		value, _ := strconv.ParseUint(parts[1], 10, 64)
		switch strings.TrimSuffix(parts[0], ":") {
		case "MemTotal":
			total = value * 1024
		case "MemAvailable":
			available = value * 1024
		}
	}
	used := total - available
	return UsageMetric{Used: used, Total: total, Usage: percent(used, total)}, scanner.Err()
}

func readDiskUsage(path string) (UsageMetric, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(filepath.Clean(path), &stat); err != nil {
		return UsageMetric{}, err
	}
	total := stat.Blocks * uint64(stat.Bsize)
	available := stat.Bavail * uint64(stat.Bsize)
	used := total - available
	return UsageMetric{Used: used, Total: total, Usage: percent(used, total)}, nil
}

func readLoadAverage() (LoadAverage, error) {
	content, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return LoadAverage{}, err
	}
	parts := strings.Fields(string(content))
	if len(parts) < 3 {
		return LoadAverage{}, fmt.Errorf("invalid loadavg")
	}
	one, _ := strconv.ParseFloat(parts[0], 64)
	five, _ := strconv.ParseFloat(parts[1], 64)
	fifteen, _ := strconv.ParseFloat(parts[2], 64)
	return LoadAverage{One: one, Five: five, Fifteen: fifteen}, nil
}

func (c *Collector) readDockerMetric() (DockerMetric, error) {
	if !c.dockerEnabled {
		return DockerMetric{}, nil
	}
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}\t{{.Image}}\t{{.State}}\t{{.Status}}")
	output, err := cmd.Output()
	if err != nil {
		return DockerMetric{}, nil
	}
	var metric DockerMetric
	for _, line := range bytes.Split(bytes.TrimSpace(output), []byte("\n")) {
		parts := strings.Split(string(line), "\t")
		if len(parts) < 4 {
			continue
		}
		container := ContainerStatus{
			Name:   strings.TrimSpace(parts[0]),
			Image:  strings.TrimSpace(parts[1]),
			State:  strings.TrimSpace(parts[2]),
			Status: strings.TrimSpace(parts[3]),
		}
		state := container.State
		switch state {
		case "running":
			metric.Running++
		case "exited":
			metric.Exited++
		}
		metric.Containers = append(metric.Containers, container)
	}
	return metric, nil
}

func (c *Collector) readServices() []ServiceStatus {
	services := make([]ServiceStatus, 0, len(c.services))
	for _, name := range c.services {
		cmd := exec.Command("systemctl", "show", name, "--property=ActiveState,SubState,LoadState", "--no-page")
		output, err := cmd.Output()
		if err != nil {
			services = append(services, ServiceStatus{Name: name, Active: "unknown", Sub: "unknown", Status: "unavailable"})
			continue
		}
		status := ServiceStatus{Name: name, Status: "healthy"}
		for _, line := range strings.Split(string(output), "\n") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			switch parts[0] {
			case "ActiveState":
				status.Active = parts[1]
			case "SubState":
				status.Sub = parts[1]
			case "LoadState":
				if parts[1] != "loaded" {
					status.Status = "missing"
				}
			}
		}
		if status.Active != "active" {
			status.Status = "degraded"
		}
		services = append(services, status)
	}
	return services
}

func percent(used, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return float64(used) * 100 / float64(total)
}

func (c *Collector) runChecks() []CheckStatus {
	results := make([]CheckStatus, 0, len(c.checks))
	for _, check := range c.checks {
		results = append(results, executeCheck(check))
	}
	return results
}

func executeCheck(check CheckConfig) CheckStatus {
	if check.Timeout == 0 {
		check.Timeout = 5 * time.Second
	}
	result := CheckStatus{Name: check.Name, Type: check.Type, Target: check.Target, Status: "healthy"}
	started := time.Now()
	switch check.Type {
	case "http":
		client := &http.Client{Timeout: check.Timeout}
		resp, err := client.Get(check.Target)
		if err != nil {
			result.Status = "failed"
			result.Error = err.Error()
			break
		}
		defer resp.Body.Close()
		expected := check.ExpectedStatus
		if expected == 0 {
			expected = http.StatusOK
		}
		if resp.StatusCode != expected {
			result.Status = "failed"
			result.Error = fmt.Sprintf("expected %d got %d", expected, resp.StatusCode)
		}
	case "tcp":
		conn, err := net.DialTimeout("tcp", check.Target, check.Timeout)
		if err != nil {
			result.Status = "failed"
			result.Error = err.Error()
			break
		}
		_ = conn.Close()
	default:
		result.Status = "failed"
		result.Error = "unsupported check type"
	}
	result.LatencyMS = time.Since(started).Milliseconds()
	return result
}

func SleepUntilNextTick(interval time.Duration, fn func() error) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	if err := fn(); err != nil {
		return err
	}
	for range ticker.C {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}
