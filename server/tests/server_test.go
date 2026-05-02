package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"opspulse/server/internal/api"
	"opspulse/server/internal/core"
)

func TestHeartbeatAndNodeListing(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "ops.db")
	srv, err := api.New(api.Config{DatabasePath: dbPath, OfflineTimeout: time.Minute, Token: "secret"})
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Close()

	report := core.NodeReport{
		NodeID: "node-1",
		SentAt: time.Now().UTC(),
		Metrics: core.NodeMetrics{
			Hostname: "worker-a",
			CPUUsage: 22,
			Memory:   core.UsageMetric{Used: 2, Total: 8, Usage: 25},
			Disk:     core.UsageMetric{Used: 10, Total: 100, Usage: 10},
			Services: []core.ServiceStatus{{Name: "docker", Status: "healthy"}},
		},
	}
	body, _ := json.Marshal(report)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/agents/heartbeat", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer secret")
	resp := httptest.NewRecorder()
	srv.Handler().ServeHTTP(resp, req)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", resp.Code)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/nodes", nil)
	listResp := httptest.NewRecorder()
	srv.Handler().ServeHTTP(listResp, listReq)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", listResp.Code)
	}
	var nodes []core.NodeRecord
	if err := json.Unmarshal(listResp.Body.Bytes(), &nodes); err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 || nodes[0].Hostname != "worker-a" {
		t.Fatalf("unexpected nodes response: %+v", nodes)
	}
}
