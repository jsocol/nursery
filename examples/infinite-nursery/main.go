package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/metrics"
	"time"

	"github.com/jsocol/nursery"
)

const goMetric = "/sched/goroutines:goroutines"

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	err := nursery.Open(ctx, func(n nursery.Nursery) error {
		n.Start(func(ctx context.Context) error {
			shutdown := make(chan os.Signal, 1)
			signal.Notify(shutdown, os.Interrupt)

			<-shutdown
			cancel()
			fmt.Println("shutting down")
			return nil
		})

		n.Start(func(ctx context.Context) error {
			sample := make([]metrics.Sample, 1)
			sample[0].Name = goMetric

			for {
				select {
				case <-time.After(5 * time.Second):
					metrics.Read(sample)
					fmt.Printf("METRIC: live goroutines = %v\n", sample[0].Value.Uint64())
				case <-ctx.Done():
					return nil
				}
			}
		})

		for {
			select {
			case <-time.After(1000 * time.Millisecond):
				fmt.Println("tick started")
				err := n.Start(func(context.Context) error {
					time.Sleep(2500 * time.Millisecond)
					fmt.Println("tick complete")
					return nil
				})
				if err != nil {
					fmt.Printf("error starting task: %v\n", err)
				}
			case <-ctx.Done():
				return nil
			}
		}
	})

	fmt.Printf("exiting with error: %v\n", err)
}
