# infinite-nursery

This example mimics a typical server setup, which would have an infinite loop
accepting new connections and handling them.

In this example, we simulate new connections with a loop of periodic ticks,
each one kicking off a new `Task` that pauses (and yields) before completing.
Graceful shutdown is also implemented as a `Task`, as is a monitor that prints
the current number of live goroutines.
