package main

import (
	"context"
	"errors"

	"all-the-frameworks/cmd/gin/api"
	"all-the-frameworks/internal/wiring"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	err := run()
	if err != nil && !errors.Is(err, wiring.ErrTerminated) {
		log.WithError(err).Fatal()
	}
}

func run() error {
	stats, err := wiring.LoadStatsd("127.0.0.1:8125", "gin.")
	if err != nil {
		return err
	}
	defer stats.Close()

	g, ctx := errgroup.WithContext(context.Background())

	a := api.New("localhost:8082", stats)
	g.Go(func() error {
		return a.Run(ctx)
	})

	g.Go(func() error {
		return wiring.HandleTermination(ctx)
	})

	return g.Wait()
}
