package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jsocol/nursery"
)

func main() {
	ctx := context.Background()

	const tasks = 5

	nursery.Open(ctx, func(n nursery.Nursery) error {
		for i := 0; i < tasks; i++ {
			func(j int) {
				n.Start(func(ctx context.Context) error {
					select {
					case <-time.After(time.Duration(tasks-j) * time.Second):
						fmt.Printf("finished task %d\n", j)
					case <-ctx.Done():
						fmt.Printf("ending task %d early\n", j)
					}
					return nil
				})
			}(i + 1)
		}
		return nil
	})
}
