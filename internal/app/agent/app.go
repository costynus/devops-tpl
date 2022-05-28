package agent

import (
	"context"
	"devops-tpl/config"
	"devops-tpl/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level)

	// Metrics -.
	metrics := NewMetrics()

	// Client -.
	client := resty.New().SetBaseURL(cfg.Agent.ServerURL)

	// WebAPI -.
	webAPI := NewWebAPI(client)

	// Worker -.
	worker := NewWorker(
		webAPI,
		metrics,
		cfg.Agent.MetricFieldNames,
		l,
	)

	updateTicker := time.NewTicker(time.Duration(cfg.Agent.PollInterval) * time.Second)
	go worker.UpdateMetrics(context.Background(), updateTicker)

	sendTicker := time.NewTicker(time.Duration(cfg.Agent.ReportInterval) * time.Second)
	go worker.SendMetrics(context.Background(), sendTicker)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	s := <-interrupt
	l.Info("agent - stoped: " + s.String())
}
