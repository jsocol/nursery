# batching

Nurseries have a natural way to limit concurrency: batch `Task`s into slices of
the maximum permitted concurrency. In this example, 5 URLs are fetched in
batches of 2, ensuring that no more than 2 requests are in flight at a time.

Because `nursery.Open()` doesn't return until all `Task`s started are done, it
can be used in a simple `for` loop. There is no need to maintain a pool of
workers or enforce concurrency limits externally.
