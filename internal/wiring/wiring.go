package wiring

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DataDog/datadog-go/statsd"
)

func LoadStatsd(addr, namespace string) (*statsd.Client, error) {
	stats, err := statsd.New(
		addr,
		statsd.WithNamespace(namespace),
		statsd.WithoutTelemetry(),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating StatsD client: %w", err)
	}
	return stats, nil
}

var ErrTerminated = errors.New("terminated")

func HandleTermination(ctx context.Context) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		return ErrTerminated
	case <-ctx.Done():
		return nil
	}
}
