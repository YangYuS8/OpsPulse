package collect

import (
	"bufio"
	"bytes"
	"fmt"
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

type Snapshot struct {
	Hostname string          `json:"hostname"`
	Uptime   int64           `json:"uptime"`
	CPUUsage float64         `json:"cpuUsage"`
	Memory   UsageMetric     `json:"memory"`
	Disk     UsageMetric     `json:"disk"`
	Load     LoadAverage     `json:"load"`
	Docker   DockerMetric    `json:"docker"`
	Services []ServiceStatus `json:"services"`
}

type Collector struct {
	services      []string
	dockerEnabled bool
	prevIdle      uint64
	prevTotal     uint64
	hasPrevCPU    bool
}

func New(serviceWhitelist []string, dockerEnabled bool) *Collector {
	return &Collector{services: serviceWhitelist, dockerEnabled: dockerEnabled}
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
	return Snapshot{
		Hostname: hostname,
		Uptime:   uptime,
		CPUUsage: cpuUsage,
		Memory:   memory,
		Disk:     disk,
		Load:     load,
		Docker:   docker,
		Services: services,
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
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.State}}")
	output, err := cmd.Output()
	if err != nil {
		return DockerMetric{}, nil
	}
	var metric DockerMetric
	for _, line := range bytes.Split(bytes.TrimSpace(output), []byte("\n")) {
		state := strings.TrimSpace(string(line))
		switch state {
		case "running":
			metric.Running++
		case "exited":
			metric.Exited++
		}
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
