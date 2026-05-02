package core

import "time"

type NodeStatus string

const (
	NodeStatusOnline  NodeStatus = "online"
	NodeStatusOffline NodeStatus = "offline"
	NodeStatusWarn    NodeStatus = "warn"
)

type ServiceStatus struct {
	Name   string `json:"name"`
	Active string `json:"active"`
	Sub    string `json:"sub"`
	Status string `json:"status"`
}

type NodeMetrics struct {
	Hostname string          `json:"hostname"`
	Uptime   int64           `json:"uptime"`
	CPUUsage float64         `json:"cpuUsage"`
	Memory   UsageMetric     `json:"memory"`
	Disk     UsageMetric     `json:"disk"`
	Load     LoadAverage     `json:"load"`
	Docker   DockerMetric    `json:"docker"`
	Services []ServiceStatus `json:"services"`
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
	Running int `json:"running"`
	Exited  int `json:"exited"`
}

type NodeReport struct {
	NodeID     string      `json:"nodeId"`
	AgentToken string      `json:"-"`
	SentAt     time.Time   `json:"sentAt"`
	Metrics    NodeMetrics `json:"metrics"`
}

type NodeRecord struct {
	NodeID        string          `json:"nodeId"`
	Hostname      string          `json:"hostname"`
	Status        NodeStatus      `json:"status"`
	LastSeen      time.Time       `json:"lastSeen"`
	Offline       bool            `json:"offline"`
	CPUUsage      float64         `json:"cpuUsage"`
	Memory        UsageMetric     `json:"memory"`
	Disk          UsageMetric     `json:"disk"`
	Load          LoadAverage     `json:"load"`
	Docker        DockerMetric    `json:"docker"`
	Services      []ServiceStatus `json:"services"`
	Uptime        int64           `json:"uptime"`
	HeartbeatAge  int64           `json:"heartbeatAgeSeconds"`
	StatusSummary string          `json:"statusSummary"`
}

type Event struct {
	ID        int64     `json:"id"`
	NodeID    string    `json:"nodeId"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

type Overview struct {
	NodesTotal   int `json:"nodesTotal"`
	NodesOnline  int `json:"nodesOnline"`
	NodesOffline int `json:"nodesOffline"`
	NodesWarn    int `json:"nodesWarn"`
	ServicesDown int `json:"servicesDown"`
	ContainersUp int `json:"containersUp"`
}
