package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/jsocol/nursery"
)

var urls = []string{
	"https://github.com/",
	"https://github.com/jsocol",
	"https://github.com/jsocol/nursery",
	"https://vorpus.org/",
	"https://vorpus.org/blog/notes-on-structured-concurrency-or-go-statement-considered-harmful/",
}

const batchSize = 2

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	st := time.Now()

	// limit concurrency by batching
	for start := 0; start <= len(urls)-1; start += batchSize {
		end := min(start+batchSize, len(urls))
		batch := urls[start:end]

		err := nursery.Open(ctx, func(n nursery.Nursery) error {
			for _, url := range batch {
				func(u string) {
					n.Start(func(ctx context.Context) error {
						s := time.Now()
						resp, err := http.Get(u)
						if err != nil {
							return err
						}
						dt := time.Since(s)
						fmt.Printf("fetched: %s [%d] in %0.3f seconds\n", resp.Request.URL, resp.StatusCode, dt.Seconds())
						return nil
					})
				}(url)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("error with batch: %#v\n", err)
			continue
		}
	}
	dt := time.Since(st)
	fmt.Printf("all batches done in %0.3f seconds\n", dt.Seconds())
}

func min(is ...int) int {
	m := math.Inf(1)
	for _, i := range is {
		j := float64(i)
		if j < m {
			m = j
		}
	}
	return int(m)
}
