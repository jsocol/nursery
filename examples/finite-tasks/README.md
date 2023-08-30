# finite-tasks

This example starts a number of tasks (defaults to 5) in a loop. Each task has
a finite duration and completes. When the final task completes, the open
nursery closes and the program exits.

Each task sleeps for a given number of seconds, so they complete in reverse
order. Tasks are also context-aware so if any are canceled early, they will all
exit.
