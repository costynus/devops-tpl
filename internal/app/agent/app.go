package agent

import (
	"context"
	agent_config "devops-tpl/config/agent"
	"devops-tpl/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
)

func Run(cfg *agent_config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level)

	// Metrics -.
	metrics := NewMetrics()

	// Client -.
	client := resty.New().SetBaseURL(cfg.Agent.ServerSchema + cfg.Agent.ServerURL)

	// WebAPI -.
	webAPI := NewWebAPI(client, cfg.KEY)

	// Worker -.
	worker := NewWorker(
		webAPI,
		metrics,
		cfg.Agent.MetricFieldNames,
		l,
	)

	updateTicker := time.NewTicker(cfg.Agent.PollInterval)
	go worker.UpdateMetrics(context.TODO(), updateTicker)

	sendTicker := time.NewTicker(cfg.Agent.ReportInterval)
	go worker.SendMetrics(context.TODO(), sendTicker)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	s := <-interrupt
	l.Info("agent - stoped: " + s.String())
}
