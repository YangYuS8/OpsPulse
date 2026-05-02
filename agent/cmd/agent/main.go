package main

import (
	"flag"
	"log"

	"opspulse/agent/internal/collect"
	"opspulse/agent/internal/config"
	"opspulse/agent/internal/report"
)

func main() {
	configPath := flag.String("config", "./agent.yaml", "path to agent config")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	collector := collect.New(cfg.ServiceWhitelist, cfg.DockerEnabled)
	client := report.New(cfg.ServerURL, cfg.Token)

	err = collect.SleepUntilNextTick(cfg.Interval, func() error {
		snapshot, err := collector.Snapshot()
		if err != nil {
			return err
		}
		if err := client.Send(cfg.NodeID, snapshot); err != nil {
			log.Printf("report failed: %v", err)
			return nil
		}
		log.Printf("heartbeat sent for %s", cfg.NodeID)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
