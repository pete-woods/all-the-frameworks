package main

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"all-the-frameworks/cmd/httprouter/api"
	"all-the-frameworks/internal/wiring"
)

func main() {
	err := run()
	if err != nil && !errors.Is(err, wiring.ErrTerminated) {
		log.WithError(err).Fatal()
	}
}

func run() error {
	stats, err := wiring.LoadStatsd("127.0.0.1:8125", "httprouter.")
	if err != nil {
		return err
	}
	defer stats.Close()

	g, ctx := errgroup.WithContext(context.Background())

	a := api.New("localhost:8083", stats)
	g.Go(func() error {
		return a.Run(ctx)
	})

	g.Go(func() error {
		return wiring.HandleTermination(ctx)
	})

	return g.Wait()
}
