package nursery_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jsocol/nursery"
)

func TestNurseryBasic(t *testing.T) {
	ctx := context.Background()
	child1 := false
	child2 := false
	err := nursery.Open(ctx, func(n nursery.Nursery) error {
		n.Start(func(context.Context) error {
			t.Log("inside child1")
			child1 = true
			return nil
		})
		n.Start(func(context.Context) error {
			t.Log("inside child2")
			time.Sleep(1 * time.Millisecond)
			child2 = true
			return nil
		})
		return nil
	})

	if !child1 {
		t.Error("func child1 did not complete")
	}
	if !child2 {
		t.Error("func child2 did not complete")
	}
	if err != nil {
		t.Errorf("nursery returned unexpected error: %#v", err)
	}
}

func TestNurseryChildError(t *testing.T) {
	ctx := context.Background()

	err := nursery.Open(ctx, func(n nursery.Nursery) error {
		// never fails
		n.Start(func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		})

		// always fails
		n.Start(func(context.Context) error {
			time.Sleep(1 * time.Millisecond)
			return errors.New("uh oh")
		})

		return nil
	})

	if err == nil {
		t.Error("nursery did not return error")
	}
	if err.Error() != "uh oh" {
		t.Errorf("got unexpected error: %#v", err)
	}
}

func TestNurseryChildCompletesAfterError(t *testing.T) {
	ctx := context.Background()
	child1 := false
	err := nursery.Open(ctx, func(n nursery.Nursery) error {
		n.Start(func(context.Context) error {
			time.Sleep(2 * time.Millisecond)
			child1 = true
			return nil
		})
		n.Start(func(context.Context) error {
			return errors.New("too soon")
		})
		return nil
	})

	if !child1 {
		t.Error("func child1 did not complete")
	}
	if err == nil {
		t.Error("nursery did not return expected error")
	}
	if err.Error() != "too soon" {
		t.Errorf("got unexpected error: %#v", err)
	}
}

func TestNurseryEmpty(t *testing.T) {
	ctx := context.Background()

	err := nursery.Open(ctx, func(n nursery.Nursery) error {
		return nil
	})
	if err != nil {
		t.Errorf("got unexpected error: %#v", err)
	}
}

func TestNurseryLongRunning(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := nursery.Open(ctx, func(n nursery.Nursery) error {
	loop:
		for {
			select {
			case <-time.After(1 * time.Millisecond):
				n.Start(func(context.Context) error {
					time.Sleep(2 * time.Millisecond)
					return nil
				})
			case <-ctx.Done():
				break loop
			}
		}
		return nil
	})

	if err == nil {
		t.Error("nursery did not return error")
	}
	if err != context.DeadlineExceeded {
		t.Errorf("got unexpected error: %#v", err)
	}
}

func ExampleOpen() {
	ctx := context.Background()

	err := nursery.Open(ctx, func(n nursery.Nursery) error {
		n.Start(func(ctx context.Context) error {
			return errors.New("I always fail")
		})

		n.Start(func(ctx context.Context) error {
			fmt.Println("I always happen")
			return nil
		})

		return nil
	})
	if err != nil {
		fmt.Println("Err: " + err.Error())
	}
	// Output:
	// I always happen
	// Err: I always fail
}
