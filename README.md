# Nursery

Nursery is an implementation of the concurrency primitive of the same name from
Nathaniel J. Smith's (a.k.a [@vorpalsmith][1] on Twitter, [njsmith][2] on
GitHub) 2018 [blog post on structured concurrency][3]. This is the same concept
as in Smith's Python library [Trio][4].

Nurseries are a control flow structure for concurrent execution that provides
(as much as possible) the same guarantee of other control flow mechanisms:
program execution enters the structure in one place, and exits the structure in
one place. What happens in between those places can happen concurrently (or in
parallel) and is guaranteed (as much as possible) to be complete by the time the
Nursery is closed again.

**Nursery should be considered ALPHA quality software.**

[1]: https://twitter.com/vorpalsmith
[2]: https://github.com/njsmith/
[3]: https://vorpus.org/blog/notes-on-structured-concurrency-or-go-statement-considered-harmful/
[4]: https://trio.readthedocs.io/

## Quick Start

To use a Nursery, there are two types of function to provide: an `Initializer`
which starts some number of `Tasks`.

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/jsocol/nursery"
)

func main() {
    ctx := context.Background()
    err := nursery.Open(ctx, func(n nursery.Nursery) error {
        n.Start(func(context.Context) error {
            time.Sleep(1 * time.Second)
            fmt.Println("task 1")
            return nil
        })

        n.Start(func(context.Context) error {
            fmt.Println("task 2, will error")
            return fmt.Errorf("always errors")
        })

        return nil
    })
    // both tasks have completed
    if err != nil {
        fmt.Println("Error! " + err)
    }
}
```

Unlike `sync.WaitGroup`, tasks can be added to the Nursery at any point during
its lifetime. Attempting to use a Nursery's `Start` method _after_ `Open` has
returned will return an `ErrStopped`.

## Context-aware Tasks

All `Task` functions are passed a `context.Context` that is a child of the
context passed to `Open`. `Task` functions should use the given context to
control any cancellable operations they start, like network activity. In the
event of an error from one of the `Tasks`, these contexts will be canceled, and
the initial error will be returned.

The initial context passed to `Open` can itself be cancellable or with a timeout
or deadline.
