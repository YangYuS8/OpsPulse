package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"opspulse/agent/internal/collect"
)

type Client struct {
	endpoint string
	token    string
	http     *http.Client
}

func New(serverURL, token string) *Client {
	return &Client{
		endpoint: strings.TrimRight(serverURL, "/") + "/api/v1/agents/heartbeat",
		token:    token,
		http:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) Send(nodeID string, snapshot collect.Snapshot) error {
	payload := map[string]any{
		"nodeId": nodeID,
		"sentAt": time.Now().UTC(),
		"metrics": snapshot,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("heartbeat rejected: %s", resp.Status)
	}
	return nil
}
