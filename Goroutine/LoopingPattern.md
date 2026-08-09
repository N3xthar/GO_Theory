# Looping in parallel 

Looping in parallel means executing the independent iterations of a loop concurrently, usually by assigning different iterations to different goroutines.

    for _, f := range filenames {
        thumbnail.ImageFile(f)
    }

instead of 
    file1 → finish
            ↓
    file2 → finish
            ↓
    file3 → finish

we can execute like this 

    file1 ──→ goroutine 1
    file2 ──→ goroutine 2
    file3 ──→ goroutine 3

# 2. What is an Embarrassingly Parallel Problem?
    An embarrassingly parallel problem is a problem where the individual subproblems are completely independent of each other, so they can be executed concurrently without requiring communication or synchronization between the individual computations.

    all of them can be processed in parallel.

# 3. How Do We Know When Goroutines Finish?
    A goroutine doesn't automatically provide a direct "finished" notification to its caller.

One solution demonstrated by the material is to use a channel as a completion signal.   

create 
    ch := make(chan struct{})
Then each worker sends a value when it finishes:
ch <- struct{}{}
The main goroutine receives one completion event for every worker.

# 4 Returning error through a channel 

errors := make(chan error)

    Worker 1 ── error ──┐
    Worker 2 ── error ──┤
    Worker 3 ── error ──┼──→ errors channel → main
    Worker 4 ── error ──┘

# Goroutine Leak bug 

imagine

    Worker 1 → success
    Worker 2 → ERROR  main go routine returns the error dude 
    Worker 3 → success
    Worker 4 → success

main go routine return error 

But the channel is unbuffered and nobody is receiving.

Therefore those goroutines block forever.

This is a goroutine leak.

definitation 
    A goroutine leak occurs when a goroutine becomes permanently blocked or otherwise fails to terminate because the communication or synchronization it is waiting for can no longer happen.

1. Solution 
use buffer channel 
Each worker can send its result without needing the main goroutine to be actively receiving at that exact moment.

Therefore, even if the main goroutine returns early because of an error, workers can still place their results into the buffer rather than becoming blocked on the send.

# Definition: sync.WaitGroup

A sync.WaitGroup is a synchronization primitive used to wait for a collection of goroutines or tasks to finish.

It maintains a counter representing the number of active operations.
concept behind it 

    Add(1)
    ↓
    counter = 1

    Add(1)
    ↓
    counter = 2

    Done()
    ↓
    counter = 1

    Done()
    ↓
    counter = 0

    Wait()
    ↓
    continue

# Three Important WaitGroup Methods

Add()
    wg.Add(1)

    Increases the counter.

    Meaning:

    "One more goroutine/task is now active."

Done()
    wg.Done()

    Decreases the counter by one.

    It is equivalent to:

    wg.Add(-1)

    The source explicitly makes this connection.

Wait()
    wg.Wait()

    Blocks until the counter becomes zero.

    Meaning:

    "Wait until all registered workers have finished."

# The WaitGroup Pattern
    var wg sync.WaitGroup

    for ... {
        wg.Add(1)

        go func() {
            defer wg.Done()

            // work
        }()
    }

    wg.Wait()

conceptiually 


                 Add(1)
                   ↓
                Worker 1
                   ↓
                Done()

                 Add(1)
                   ↓
                Worker 2
                   ↓
                Done()

                 Add(1)
                   ↓
                Worker 3
                   ↓
                Done()

                   ↓
                counter = 0
                   ↓
                Wait returns

# Important theory questions 

1. Interview Definition — WaitGroup

If the interviewer asks:

"What is a WaitGroup?"

Say:

sync.WaitGroup is a synchronization primitive in Go used to wait for a collection of goroutines to complete. It maintains a counter representing active goroutines or tasks. Add() increments the counter, Done() decrements it, and Wait() blocks until the counter reaches zero. A common pattern is to call Add(1) before starting each worker, use defer wg.Done() inside the worker, and call wg.Wait() from the goroutine that needs to wait for all workers to finish.

That's a solid interview definition.

32. Interview Definition — Goroutine Leak

A goroutine leak occurs when a goroutine remains blocked indefinitely and cannot terminate because the communication or synchronization it is waiting for is no longer possible. For example, if worker goroutines send results to an unbuffered channel and the receiving goroutine returns early after encountering an error, the remaining workers can block forever trying to send their results.

33. Interview Definition — Embarrassingly Parallel

An embarrassingly parallel problem is a problem in which individual tasks are completely independent of one another, so they can be executed concurrently with minimal synchronization. Image processing is an example because generating a thumbnail for one image does not depend on generating a thumbnail for another image.

34. Interview Definition — Parallel Loop

A parallel loop is a loop whose independent iterations are executed concurrently, typically by launching goroutines for the individual iterations. Because the iterations execute independently, the program can utilize multiple CPUs and overlap operations such as I/O. However, the program must still provide synchronization to ensure that all required work completes and to safely collect results.

35. Most Important Interview Questions
# Go Concurrency — Parallel Thumbnail Generation

## Basic

### 1. What is looping in parallel?

Looping in parallel means processing independent loop iterations concurrently instead of processing them one by one.

```go
for _, file := range files {
    go process(file)
}
```

Each iteration can run in a separate goroutine.

---

### 2. What is an embarrassingly parallel problem?

An embarrassingly parallel problem is a problem where the work can be divided into independent tasks that require little or no communication between them.

**Example:** Generating thumbnails for multiple images.

```text
Image 1 → Thumbnail 1
Image 2 → Thumbnail 2
Image 3 → Thumbnail 3
```

**Key idea:**

> Independent tasks → easy to execute in parallel.

---

### 3. Why is thumbnail generation suitable for parallelism?

Each image can be processed independently.

Image 1 does not need the result of Image 2.

```text
Image 1 ──→ Thumbnail 1
Image 2 ──→ Thumbnail 2
Image 3 ──→ Thumbnail 3
```

Therefore, multiple images can be processed at the same time.

---

### 4. Why is simply adding `go` not enough?

If we only write:

```go
for _, file := range files {
    go process(file)
}
```

the main goroutine may finish before the worker goroutines complete.

When `main()` exits, the program terminates and unfinished goroutines are stopped.

Therefore, we need synchronization to wait for the goroutines.

---

### 5. How can we wait for goroutines using a channel?

We can use a channel as a completion signal.

```go
done := make(chan struct{})

go func() {
    process()
    done <- struct{}{}
}()

<-done
```

The main goroutine blocks at `<-done` until the worker sends the completion signal.

---

### 6. Why use `chan struct{}` for completion signaling?

We only need to communicate:

> "The work is finished."

We don't need to send any actual data.

`struct{}` is commonly used for signaling because it carries no meaningful data.

```go
done := make(chan struct{})
```

So:

```go
done <- struct{}{}
```

means:

> "I am done."

---

### 7. What is loop-variable capture?

Loop-variable capture happens when a goroutine closure refers to a loop variable instead of receiving its value directly.

The goroutine may observe a value different from the one we intended.

This is especially important when working with older Go loop-variable semantics.

---

### 8. How do you safely pass the loop variable to a goroutine?

Pass the loop variable as a function argument.

```go
for _, f := range files {
    go func(f string) {
        process(f)
    }(f)
}
```

Now each goroutine receives its own value of `f`.

---

# Intermediate

## 9. How can a goroutine return a value to the main goroutine?

A goroutine cannot directly return a value to the goroutine that started it.

Instead, use a channel.

```go
result := make(chan int)

go func() {
    result <- 42
}()

value := <-result

fmt.Println(value)
```

The worker sends the value through the channel, and the main goroutine receives it.

---

## 10. How can errors be communicated through channels?

Use an error channel.

```go
errors := make(chan error)

go func() {
    _, err := process()

    errors <- err
}()

err := <-errors
```

Workers send errors through the channel.

```text
Worker 1 ── error ──┐
Worker 2 ── error ──┤
Worker 3 ── error ──┼──→ errors channel → main
Worker 4 ── error ──┘
```

---

## 11. What is a goroutine leak?

A goroutine leak occurs when a goroutine remains blocked forever and cannot finish.

Example:

```go
ch := make(chan int)

go func() {
    ch <- 10
}()
```

If nobody receives from `ch`, the goroutine blocks forever.

```text
Goroutine
    ↓
ch <- 10
    ↓
No receiver
    ↓
Blocked forever
```

---

## 12. How does `makeThumbnails4` leak goroutines?

If workers send results through an unbuffered channel but the receiving goroutine does not receive enough results, workers can remain blocked.

```text
Worker
   ↓
result channel
   ↓
No receiver
   ↓
Blocked goroutine
   ↓
Goroutine leak
```

The problem is that sending to an unbuffered channel requires a receiver to be ready.

---

## 13. How does buffering solve that particular problem?

A buffered channel can store results temporarily.

```go
results := make(chan result, len(filenames))
```

Workers can send results into the buffer without immediately waiting for a receiver.

```text
Worker
   ↓
Buffered channel
   ↓
[ result ][ result ][ result ]
   ↓
Main receives later
```

For this particular problem, giving the channel enough capacity prevents workers from getting stuck while sending their results.

---

## 14. Why use a result struct containing both result and error?

Instead of using separate channels:

```go
results chan string
errors  chan error
```

we can combine the result and error into one struct.

```go
type result struct {
    thumbnail string
    err       error
}
```

Then:

```go
results <- result{
    thumbnail: thumbnail,
    err:       err,
}
```

This keeps the result and its associated error together.

---

## 15. Why are results returned in arbitrary order?

Because goroutines execute concurrently.

For example:

```text
Worker 1 → slow
Worker 2 → fast
Worker 3 → medium
```

The results may arrive as:

```text
Worker 2
Worker 3
Worker 1
```

The order depends on goroutine scheduling and execution time.

Therefore, concurrent results are not necessarily returned in the same order as the input.

---

# Advanced

## 16. Why can't we use `len(filenames)` when filenames arrive through a channel?

When filenames arrive through a channel, the number of values may not be known in advance.

```go
for f := range filenames {
    process(f)
}
```

The producer can send values dynamically.

Therefore, we cannot depend on `len(filenames)` to determine how many values will arrive.

The receiver simply continues until the channel is closed.

---

## 17. What is `sync.WaitGroup`?

`sync.WaitGroup` is a synchronization mechanism used to wait for a collection of goroutines to finish.

```go
var wg sync.WaitGroup
```

It maintains a counter representing the number of active goroutines.

---

## 18. Explain `Add`, `Done`, and `Wait`.

### `Add(n)`

Increases the WaitGroup counter.

```go
wg.Add(1)
```

Means:

> One more goroutine needs to finish.

### `Done()`

Decreases the counter by one.

```go
wg.Done()
```

Means:

> One goroutine has finished.

### `Wait()`

Blocks until the counter becomes zero.

```go
wg.Wait()
```

Means:

> Wait until all registered goroutines finish.

---

## 19. Why must `wg.Add(1)` happen before launching the goroutine?

We should register the goroutine with the WaitGroup before starting it.

Correct:

```go
wg.Add(1)

go func() {
    defer wg.Done()
    process()
}()
```

If `Wait()` sees the counter as zero before `Add(1)` happens, it may return before the goroutine has been properly registered.

Therefore:

```text
Add(1)
   ↓
Start goroutine
   ↓
Done()
   ↓
Wait()
```

---

## 20. Why use `defer wg.Done()`?

Because `Done()` must always be called when the goroutine finishes.

```go
go func() {
    defer wg.Done()

    process()
}()
```

`defer` ensures `wg.Done()` executes when the goroutine exits, even if the function returns early.

It also reduces the chance of forgetting to call `Done()`.

---

## 21. What is the purpose of the closer goroutine?

The closer goroutine waits for all workers to finish and then closes the result channel.

```go
go func() {
    wg.Wait()
    close(results)
}()
```

Its purpose is to tell the receiver:

> "All workers are finished. No more results will arrive."

---

## 22. Why can't `wg.Wait()` happen before the `range sizes` loop?

Workers need to receive values from the `sizes` channel to do their work.

If the main goroutine waits for workers before providing the work properly, workers may be waiting for input while the main goroutine is waiting for workers.

This can cause a deadlock.

```text
Main
 ↓
wg.Wait()
 ↓
Waiting for workers

Workers
 ↓
Waiting for input
```

Neither side can make progress.

---

## 23. Why can't `wg.Wait()` simply happen after the `range sizes` loop?

Because workers may still be trying to send results while the main goroutine is waiting.

With an unbuffered results channel:

```text
Main
 ↓
wg.Wait()
 ↓
Waiting for workers

Worker
 ↓
results <- result
 ↓
No receiver
 ↓
Blocked
```

The worker cannot finish because it is blocked sending the result.

The main goroutine cannot finish waiting because the worker is blocked.

This can cause a deadlock.

---

## 24. Why does the closer goroutine call `close(results)`?

Once all workers finish:

```go
wg.Wait()
```

there will be no more results.

Therefore, the closer goroutine closes the channel:

```go
close(results)
```

This allows the receiver to finish:

```go
for r := range results {
    // handle result
}
```

The range loop stops when `results` is closed and all buffered values have been received.

---

# Complete `makeThumbnails6` Architecture

The architecture has three major components:

1. A channel for distributing work.
2. Multiple worker goroutines.
3. A result channel for collecting results.

```text
                         filenames
                             │
                             ▼
                     ┌──────────────┐
                     │  sizes chan  │
                     └──────┬───────┘
                            │
              ┌─────────────┼─────────────┐
              ▼             ▼             ▼
          Worker 1      Worker 2      Worker 3
              │             │             │
              └─────────────┼─────────────┘
                            │
                            ▼
                     ┌──────────────┐
                     │ results chan │
                     └──────┬───────┘
                            │
                            ▼
                           Main
```

The `WaitGroup` tracks the workers:

```text
Worker 1 ──→ wg.Done()
Worker 2 ──→ wg.Done()
Worker 3 ──→ wg.Done()
                 │
                 ▼
              wg.Wait()
                 │
                 ▼
           close(results)
```

---

## Complete Flow

### Step 1 — Send work

Filenames are sent through the `sizes` channel.

```text
filenames → sizes
```

### Step 2 — Workers receive work

Multiple goroutines receive filenames.

```text
sizes → Worker 1
sizes → Worker 2
sizes → Worker 3
```

### Step 3 — Process independently

Each worker generates the thumbnail for its file.

```text
Worker → ImageFile()
```

### Step 4 — Send result

Workers send their result to the results channel.

```text
Worker → results
```

### Step 5 — Track completion

Each worker calls:

```go
defer wg.Done()
```

### Step 6 — Wait for all workers

A separate goroutine executes:

```go
wg.Wait()
```

### Step 7 — Close results

After all workers finish:

```go
close(results)
```

### Step 8 — Main receives results

The main goroutine ranges over the results:

```go
for r := range results {
    // handle result
}
```

When the channel is closed and all results are consumed, the loop ends.

---

# Interview Cheat Sheet

| Concept | Short Answer |
|---|---|
| Parallel loop | Execute independent loop iterations in parallel |
| Embarrassingly parallel | Independent tasks with little/no communication |
| Thumbnail generation | Each image can be processed independently |
| `go` alone | Doesn't wait for goroutines |
| Channel | Communication between goroutines |
| `chan struct{}` | Signal without carrying data |
| Loop capture | Closure can capture a loop variable unexpectedly |
| Safe loop variable | Pass it as a function argument |
| Goroutine result | Return it through a channel |
| Error communication | Send `error` through a channel |
| Goroutine leak | Goroutine blocked forever |
| Buffered channel | Temporarily stores values and can prevent blocking |
| Result struct | Keeps result and error together |
| Arbitrary order | Goroutines finish at different times |
| `WaitGroup` | Waits for multiple goroutines |
| `Add(1)` | Register a goroutine |
| `Done()` | Mark a goroutine finished |
| `Wait()` | Wait for counter to become zero |
| `defer Done()` | Ensure `Done()` is called on exit |
| Closer goroutine | Waits for workers and closes result channel |
| `close(results)` | Signals that no more results will arrive |

---

# One-Line Mental Model

```text
Independent Work
      ↓
Goroutines
      ↓
Channels = Communication
      ↓
WaitGroup = Completion Tracking
      ↓
Closer Goroutine
      ↓
close(results)
      ↓
Main receives all results
```

> **Embarrassingly parallel = the problem is easy to parallelize because its tasks are independent.**
36. The Core Theory You Should Remember

If you remember only the conceptual progression of this section, remember this:

Sequential loop
      ↓
Add goroutines
      ↓
Oops — function returns too early
      ↓
Need completion synchronization
      ↓
Use channel as completion signal
      ↓
Need to pass loop variable correctly
      ↓
Need to return results/errors
      ↓
Early return can cause goroutine leak
      ↓
Buffered channel can prevent workers
from blocking on result sends
      ↓
But what if number of workers is unknown?
      ↓
Use sync.WaitGroup
      ↓
Workers send results through channel
      ↓
Closer goroutine waits for WaitGroup
      ↓
Closer closes result channel
      ↓
Main ranges over result channel

That progression is the real theory of this section.

And for an interview, the three patterns I'd know extremely well are:

// Pattern 1: completion channel

    ch := make(chan struct{})

    go func() {
        // work
        ch <- struct{}{}
    }()

    <-ch
// Pattern 2: WaitGroup
    var wg sync.WaitGroup

    wg.Add(1)

    go func() {
        defer wg.Done()
        // work
    }()

    wg.Wait()
// Pattern 3: WaitGroup + result channel
    var wg sync.WaitGroup
    results := make(chan Result)

    wg.Add(1)
    go func() {
        defer wg.Done()
        results <- result
    }()

    go func() {
        wg.Wait()
        close(results)
    }()

    for result := range results {
        // consume result
    }

Pattern 3 is the big one from this section. It combines parallel workers + synchronization + result collection + channel closure, and understanding why the closer goroutine exists is exactly the kind of thing an interviewer can probe deeply.