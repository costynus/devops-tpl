package agent

import (
	"context"
	"devops-tpl/config"
	"devops-tpl/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level)

	// Metrics -.
	metrics := NewMetrics()

	// WebAPI -.
	webAPI := NewWebAPI(l, cfg.Agent.ServerURL)

	// Worker -.
	worker := NewWorker(l, metrics, webAPI)

	updateTicker := time.NewTicker(time.Duration(cfg.Agent.PollInterval) * time.Second)
	go worker.UpdateMetrics(context.Background(), updateTicker)

	sendTicker := time.NewTicker(time.Duration(cfg.Agent.ReportInterval) * time.Second)
	go worker.SendMetrics(context.Background(), sendTicker)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case s := <-interrupt:
		l.Info("agent - Run signal: " + s.String())
	}
}
