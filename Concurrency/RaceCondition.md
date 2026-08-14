# 9. Concurrency with Shared Variables — Top 1% Interview Master Sheet

> Har topic mein 4 layers: **Definition → Real Example → Internals (kaise andar kaam karta hai) → Interview Q&A (junior se staff-level tak).** Ye depth zyada candidates nahi jaante — yahi tujhe alag banayega.

---

## 9.1 Race Conditions

### Definition
A **race condition** occurs when two or more goroutines access the same variable concurrently, and at least one of the accesses is a write, without proper synchronization. The final result depends on the unpredictable timing/interleaving of goroutine execution.

### Real-world example
```go
var balance int

func deposit(amount int) {
    balance += amount // NOT atomic: read, add, write = 3 steps
}

func main() {
    for i := 0; i < 1000; i++ {
        go deposit(1)
    }
    time.Sleep(time.Second)
    fmt.Println(balance) // rarely prints 1000 — classic race
}
```

### 🔬 Internals — what actually happens at the machine level
`balance += amount` on x86 compiles roughly to:
```asm
MOV  EAX, [balance]   ; LOAD
ADD  EAX, amount      ; MODIFY (in register, not memory)
MOV  [balance], EAX   ; STORE
```
Two goroutines running this on different CPU cores each have their **own core-local cache line copy** of `balance` (MESI cache coherence protocol). Goroutine A can LOAD balance=5 into its core's cache, goroutine B LOADs the same balance=5 into a different core's cache, both ADD 1, both STORE 6 — one increment is silently lost even though nothing "crashed." This is why race conditions are non-deterministic: they depend on cache coherence traffic timing, not just goroutine scheduling.

The Go **memory model formally defines "race"**: a program is race-free if, for every pair of conflicting accesses (same variable, at least one write), one happens-before the other. The compiler and CPU are legally allowed to reorder any instructions that don't violate a single goroutine's own sequential semantics — this is why "it works on my machine" for racy code is common; different CPU architectures (ARM vs x86) have different reordering aggressiveness, ARM being weaker-ordered and exposing races x86 might hide.

### Interview Q&A
- **Q: Difference between race condition and data race?** Data race = the technical condition (unsynchronized concurrent access, ≥1 write) that the race detector flags via shadow memory + vector clocks. Race condition = the broader outcome-depends-on-timing bug; can exist without a data race (e.g. two synchronized operations still racing at the logic level — TOCTOU bugs).
- **Q: Is `balance++` atomic in Go?** Never — regardless of int size, on any platform, by spec.
- **Staff-level Q: Why can a race bug pass 10,000 test runs and then fail in production?** Because the interleaving window is a handful of nanoseconds; test machines often have fewer cores or different load patterns than production, changing scheduling probability. This is why `-race` (dynamic instrumentation) matters more than "it passed CI" — it flags the *possibility*, not just observed failures.

---

## 9.2 Mutual Exclusion: sync.Mutex

### Definition
`sync.Mutex` ensures only one goroutine at a time enters a critical section. `Lock()` acquires, `Unlock()` releases; blocked goroutines wait until it's free.

### Real-world example
```go
var (
    balance int
    mu      sync.Mutex
)

func deposit(amount int) {
    mu.Lock()
    defer mu.Unlock()
    balance += amount
}
```

### 🔬 Internals — how Mutex is actually implemented (Go source: sync/mutex.go)
A `Mutex` is just two integers under the hood:
```go
type Mutex struct {
    state int32  // bit-packed: locked | woken | starving | waiter count
    sema  uint32 // OS-level semaphore for blocking/waking goroutines
}
```
Two operating modes, and this is genuinely rare knowledge:

1. **Normal mode** (default): if `Lock()` fails, the goroutine doesn't sleep immediately — it **spins** briefly (active-spin, up to 4 iterations, only if on a multi-core machine and the lock is expected to free up soon) hoping to grab the lock without the expensive cost of parking on the OS semaphore. This is a huge throughput win under light contention. If spinning fails, the goroutine queues (FIFO) and parks via `runtime_SemacquireMutex`.
2. **Starvation mode**: if a goroutine waits **more than 1ms** for the lock, the mutex flips into starvation mode — new arrivals can no longer barge in front of the queue (no spinning, no "lucky" acquire); the lock is handed directly, FIFO, to the longest-waiting goroutine. This exists specifically to prevent goroutine starvation under heavy contention (a bug that plagued early Go's naive mutex). It flips back to normal mode once the queue is empty or the wait time drops below 1ms.

This means: under **low contention**, `sync.Mutex` behaves almost like a lock-free spinlock (very fast). Under **high contention**, it behaves like a strict FIFO ticket lock (fair but slower). Knowing this trade-off — and that Go engineered it deliberately — is a strong senior-level signal.

### Key rules
1. Never copy a mutex (copies reset the `state`/`sema` — pass structs by pointer).
2. Zero value is ready to use.
3. `Unlock()` in `defer` for panic/early-return safety.
4. Minimize critical section — no I/O, no heavy compute, no calling into user callbacks while holding the lock (risk of reentrant deadlock or long stalls).
5. Not reentrant.

### Interview Q&A
- **Q: What happens if you Unlock() an unlocked mutex?** Panic: "sync: unlock of unlocked mutex."
- **Q: Can a goroutine holding a lock lock it again?** No — self-deadlock; Go's Mutex is not reentrant (no owner tracking at all, unlike Java's monitor).
- **Q: Deadlock avoidance with multiple mutexes?** Global lock ordering across the whole codebase.
- **Staff-level Q: Why doesn't Go's Mutex track "owner goroutine"?** Deliberate design choice — tracking ownership adds overhead to the hot path (every Lock/Unlock) for a feature (reentrancy, deadlock detection) Go's designers consider an anti-pattern to lean on. Philosophy: restructure code to avoid needing reentrant locks rather than pay the cost for everyone.
- **Staff-level Q: What's the cost difference between contended vs uncontended Lock()?** Uncontended = single CAS (compare-and-swap) instruction, a few nanoseconds. Contended = CAS failure, possible spin loop, then a syscall-level futex wait — orders of magnitude slower (microseconds). This is why lock granularity/sharding matters at scale (e.g. sharded maps instead of one global mutex).

---

## 9.3 Read/Write Mutexes: sync.RWMutex

### Definition
`sync.RWMutex` allows multiple simultaneous readers OR one exclusive writer. `RLock()/RUnlock()` for readers, `Lock()/Unlock()` for writers.

### Real-world example
```go
type ConfigCache struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *ConfigCache) Get(key string) string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.data[key]
}

func (c *ConfigCache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}
```

### 🔬 Internals — the writer-starvation-prevention trick
```go
type RWMutex struct {
    w           Mutex  // held by writers, and by the first blocked writer against new readers
    writerSem   uint32
    readerSem   uint32
    readerCount atomic.Int32 // can go NEGATIVE — this is the clever part
    readerWait  atomic.Int32
}
```
- `RLock()` just does `atomic.AddInt32(&readerCount, 1)` — no full mutex needed for the common case, extremely cheap.
- When a **writer** calls `Lock()`: it acquires the internal `w` mutex (blocking other writers), then does `atomic.AddInt32(&readerCount, -rwmutexMaxReaders)` — this makes `readerCount` **go negative**. Any *new* `RLock()` call checks `if atomic.AddInt32(&readerCount,1) < 0 { block }` — so once a writer is waiting, every new reader immediately sees a negative count and blocks instead of barging ahead. This is the exact mechanism that prevents writer starvation without needing a separate "waiting writers" flag.
- The writer then waits (`readerWait`) only for **already-in-flight** readers to finish (`RUnlock()` decrements and signals when it hits zero), not new ones.

This negative-counter trick is a favorite "explain RWMutex internals" staff-interview question — most candidates only know the surface behavior, not *how* starvation is actually prevented.

### Interview Q&A
- **Q: When is RWMutex worse than Mutex?** Frequent writes, or tiny critical sections — RWMutex's extra bookkeeping (atomic ops + two semaphores) costs more than it saves. Rule of thumb: RWMutex wins only when reads are the vast majority AND critical sections aren't trivially short.
- **Q: Can a writer starve?** No — new readers block behind a waiting writer (negative counter trick above).
- **Q: Map safety?** Go maps are never concurrency-safe by default for mixed read/write; RWMutex or `sync.Map` must be added explicitly.
- **Staff-level Q: Why can RLock() still be "expensive" under heavy write contention even without an active writer?** Because every RLock/RUnlock is an atomic CAS/add — under high core counts, atomic ops on a shared cache line cause cache-line ping-ponging (see 9.9 False Sharing below), which can dominate cost even without lock blocking.

---

## 9.4 Memory Synchronization

### Definition
Memory synchronization is the process of making sure that when one goroutine changes some data, other goroutines can see those changes correctly and in the proper order.

Even with correct locking, CPU caches, compiler reordering, and memory visibility can cause a goroutine to not "see" another's writes without a proper **happens-before** relationship. Synchronization primitives establish this ordering guarantee — not just mutual exclusion, but *visibility*.

### Real-world example
```go
var ready bool
var data string

func producer() {
    data = "hello"
    ready = true // no sync — compiler/CPU may reorder, or the write may never become visible
}

func consumer() {
    for !ready { }
    fmt.Println(data) // may print "" even after ready==true is observed
}
```

### 🔬 Internals — the formal Go Memory Model rules (know these by name)
The [Go Memory Model](https://go.dev/ref/mem) defines happens-before edges precisely. The ones asked most often:
1. **Program order** — within a single goroutine, statements happen-before each other in source order (from that goroutine's own view only — compiler CAN still reorder as long as this goroutine can't observe it).
2. **`go` statement** — the `go f()` statement happens-before `f` begins executing.
3. **Goroutine exit** — is not guaranteed to happen-before anything; you must synchronize explicitly (e.g. WaitGroup) — this trips people up, they assume "goroutine finished" is automatically visible.
4. **Channel send happens-before the corresponding receive completes** (unbuffered) — this is the backbone of "share memory by communicating."
5. **The kth receive on a channel of capacity C happens-before the (k+C)th send completes** — buffered channel back-pressure ordering.
6. **`Mutex.Unlock()` happens-before a subsequent `Lock()`** (on the same mutex) returns — this is *why* mutexes give visibility, not just exclusion.
7. **`sync.Once.Do(f)` return happens-before any call to `Do` returns** — every caller sees f's effects.

Why `ready bool` fails: there's no happens-before edge between `producer`'s write and `consumer`'s read at all — no channel, no mutex, no atomic. The Go spec makes **zero guarantees** in this case, not just "might be slow" — the compiler is legally allowed to cache `ready` in a register forever inside the `for` loop (infinite loop) since Go 1.19 clarified this hazard.

### Interview Q&A
- **Q: What is happens-before?** A partial order over memory operations guaranteeing visibility — if A happens-before B, A's writes are visible to B.
- **Q: Why doesn't a raw bool flag work for signaling?** No happens-before edge exists — compiler/CPU reordering and caching are unconstrained.
- **Q: Safest way to signal "done"?** Channel close, `sync.WaitGroup`, or `context.Done()`.
- **Staff-level Q: Does `sync/atomic` alone give happens-before guarantees, or just atomicity?** Since Go 1.19, `sync/atomic` operations (and the new `atomic.Int32` etc. types) **do** provide sequentially-consistent happens-before semantics per the memory model — earlier informal guidance was murkier. Knowing the 1.19 memory-model formalization is a strong signal you've read the spec, not just blog posts.

---

## 9.5 Lazy Initialization: sync.Once

### Definition
`sync.Once` guarantees a function runs exactly once, even under concurrent calls. `Do(f)` — first caller runs `f`; everyone else blocks until done, then all return without re-running.

### Real-world example
```go
var (
    once   sync.Once
    dbConn *sql.DB
)

func GetDB() *sql.DB {
    once.Do(func() {
        dbConn, _ = sql.Open("postgres", dsn)
    })
    return dbConn
}
```

### 🔬 Internals — the double-checked-locking pattern done right
```go
type Once struct {
    done atomic.Uint32
    m    Mutex
}

func (o *Once) Do(f func()) {
    if o.done.Load() == 0 { // fast path: single atomic read, no lock
        o.doSlow(f)
    }
}

func (o *Once) doSlow(f func()) {
    o.m.Lock()
    defer o.m.Unlock()
    if o.done.Load() == 0 {
        defer o.done.Store(1)
        f()
    }
}
```
This is the textbook correct **double-checked locking** pattern — the fast path (`Do`) is a single atomic load with no mutex overhead at all for the 99.999% of calls after the first. Only the *first* call (or callers racing to be first) pay the mutex cost. This is *the* canonical example interviewers use to test whether you understand why naive double-checked locking (`if !initialized { lock(); if !initialized {...} }` in languages without a proper memory model) is broken in general but correct here — because `atomic.Uint32` gives the happens-before guarantee C/C++ pre-C11 or naive Java lacked.

### Interview Q&A
- **Q: What if f() panics?** `done` is set (via `defer`) regardless — Once is permanently "used up," f will never run again. (Note: in older Go versions the `defer o.done.Store(1)` ordering is deliberate — it runs even on panic, unlike some naive reimplementations people write that skip it.)
- **Q: Is Once goroutine-safe by itself?** Yes — that's its entire purpose.
- **Q: Why not `if initialized {}`?** Check-then-act isn't atomic — classic TOCTOU race.
- **Staff-level Q: Why is the fast path a plain atomic load and not a full Lock()?** Performance — `Once.Do` is often called on every request (e.g. lazy singleton access). Making the common case a single uncontended atomic load (nanoseconds) instead of a mutex acquire (which still involves a CAS at minimum) is a deliberate hot-path optimization — a great example of "pay for what you use" design.

---

## 9.6 The Race Detector

### Definition

Race Detector Go ka ek tool hai jo program run karte time check karta hai ki multiple goroutines same data ko unsafe way mein access toh nahi kar rahi hain.

other way 

Go ka Race Detector hume batata hai ki hamare code mein data race ho rahi hai ya nahi."

Go's `-race` flag instruments memory accesses to detect data races dynamically at runtime, reporting exactly which goroutines/lines raced on which variable.

### Real-world usage
```bash
go run -race main.go
go test -race ./...
go build -race -o app
```

### Sample output
```
WARNING: DATA RACE
Write at 0x00c0000140a0 by goroutine 7:
  main.deposit()
      /main.go:12 +0x3c

Previous write at 0x00c0000140a0 by goroutine 6:
  main.deposit()
      /main.go:12 +0x3c
```

### 🔬 Internals — how it actually detects races (this is rare knowledge)
The race detector is built on Google's **ThreadSanitizer (TSan)**, linked in via cgo. Two core mechanisms:
1. **Shadow memory**: every byte of real memory gets a shadow region tracking metadata — which goroutine last accessed it, with what clock value.
2. **Vector clocks**: every goroutine carries a vector clock (logical timestamp per goroutine). On every synchronization event (mutex lock/unlock, channel send/receive, WaitGroup, etc.) clocks are merged, establishing happens-before order. On every memory access, TSan checks: "is there a prior access to this address from another goroutine that's NOT ordered by happens-before relative to me?" If yes → race reported.

This is why it's **dynamic, not static**: it only catches races on code paths *actually executed* during that run — untested concurrent paths remain invisible. It's also why it's expensive: every single memory access now involves a shadow-memory lookup and vector-clock comparison (hence 2-10x slowdown, 5-10x memory).

### Interview Q&A
- **Q: Does it catch all races?** No — only ones triggered during that execution; needs good concurrent test coverage/fuzzing to be effective.
- **Q: Performance cost?** 2-10x slower, 5-10x memory — CI/testing only, never production.
- **Q: Use in CI?** Yes, standard practice on every PR / nightly build.
- **Staff-level Q: Why can't the race detector be always-on in production as a safety net?** The instrumentation overhead (shadow memory + vector clock per access) makes it economically infeasible at scale — it's a debugging tool, not a runtime guard rail; production safety comes from *design* (proper sync primitives, code review, `-race` in CI) not runtime detection.

---

## 9.7 Example: Concurrent Non-Blocking Cache

### Definition / Pattern
    
    This pattern stores frequently requested data in a cache and ensures that when many requests ask for the same missing data at the same time, only one request performs the expensive work and the others share its result

    Cache = "already kaam ho chuka hai, result rakha hai."

    Singleflight = "kaam already chal raha hai, doosre same kaam ko mat chalao."

#     Singleflight — Simple Definition

    Singleflight is a pattern where, if multiple goroutines request the same data at the same time, only one goroutine performs the actual work, while the others wait and share the same result.

A cache where `Get`/`Set` for different keys don't block each other, and duplicate concurrent requests for the same missing key don't trigger redundant work — memoization + request deduplication (also called "singleflight" pattern in Go's ecosystem, `golang.org/x/sync/singleflight` is the production-grade version of this).

### Real-world example
```go
type result struct {
    value string
    err   error
}

type Memo struct {
    mu    sync.Mutex
    cache map[string]*entry
}

type entry struct {
    res   result
    ready chan struct{}
}

func NewMemo() *Memo {
    return &Memo{cache: make(map[string]*entry)}
}

func (m *Memo) Get(key string, fetch func(string) (string, error)) (string, error) {
    m.mu.Lock()
    e, ok := m.cache[key]
    if !ok {
        e = &entry{ready: make(chan struct{})}
        m.cache[key] = e
        m.mu.Unlock()

        e.res.value, e.res.err = fetch(key)
        close(e.ready)
    } else {
        m.mu.Unlock()
        <-e.ready
    }
    return e.res.value, e.res.err
}
```

### Why this design
1. Mutex protects only the map access, not the fetch itself — different keys never block each other.
2. First goroutine for a key creates the entry, releases lock, does slow work outside the lock.
3. Concurrent goroutines for the *same* key wait on the channel — no duplicate fetch.
4. `close(e.ready)` broadcasts to all waiters simultaneously.

### 🔬 Internals — why `close()` is the right broadcast primitive
Closing a channel is the only Go primitive that lets **N goroutines all unblock from one event simultaneously** without a loop or polling. A closed channel's receive (`<-ch`) always returns immediately with the zero value — every blocked receiver wakes at once (the runtime moves all waiters from the channel's wait queue to runnable in one operation). Compare to `sync.Cond.Broadcast()` (older, more error-prone, needs manual predicate-checking under a lock) — `close()` on a `chan struct{}` is the idiomatic modern replacement for most Cond use cases.

Production-grade version (`x/sync/singleflight`) adds: panic recovery (so one panicking fetch doesn't hang all waiters forever), a `Forget(key)` method to evict in-flight entries, and returns whether the result was *shared* (dedup'd) — worth mentioning by name in interviews to show you know the ecosystem beyond stdlib.

### Interview Q&A
- **Q: Why not hold the lock during fetch()?** Would serialize all requests for all keys, killing concurrency.
- **Q: Why channel instead of sync.Cond?** Simpler, less error-prone broadcast to N waiters.
- **Q: What if fetch() panics?** Naive version: waiters hang forever (channel never closes). Fix: `defer close(e.ready)` plus `recover()`, storing the panic as an error.
- **Staff-level Q: How would you add a TTL / cache invalidation to this without breaking the dedup guarantee?** Store an expiry timestamp in `entry`; on `Get`, if expired, delete from map under the lock and re-fetch as if it were a cache miss — but must be careful: readers that already got a reference to the *old* `entry` and are waiting on its `ready` channel should still get that (now-stale) result rather than erroring, to keep the contract simple. This is a real production nuance worth raising unprompted.

---

## 9.8 Goroutines and Threads

### Definition
A goroutine is a lightweight, user-space thread managed by the Go runtime, not the OS. Go multiplexes many goroutines onto few OS threads via the **M:N scheduler**, using the **GMP model**: G (goroutine), M (OS thread/"machine"), P (logical processor — holds a local run queue, required for an M to execute Go code).

### Comparison table

| Aspect | Goroutine | OS Thread |
|---|---|---|
| Initial stack | ~2KB, grows/shrinks dynamically (segmented → contiguous copying since Go 1.4) | 1-8MB fixed |
| Created by | Go runtime (cheap, no syscall) | OS kernel (syscall, expensive) |
| Scheduling | M:N, cooperative + async-preemptive | 1:1, OS preemptive |
| Context switch | ~tens of ns, no kernel trap | ~1-2 microseconds, kernel involved |
| Typical count | Millions | Thousands |
| Communication | Channels (idiomatic) | Shared memory + OS primitives |

### 🔬 Internals — GMP scheduler deep dive (staff-level territory)
- **G (Goroutine)**: struct with its own stack, program counter, status (`_Grunnable`, `_Grunning`, `_Gwaiting`, `_Gdead`, etc.)
- **M (Machine)**: an actual OS thread. Can only run Go code while holding a P.
- **P (Processor)**: a logical resource, count = `GOMAXPROCS`. Holds a **local run queue (LRQ)** of up to 256 runnable Gs, plus access to a **global run queue (GRQ)**.
- **Work-stealing**: when a P's local queue is empty, it tries, in order: (1) steal half of another random P's LRQ, (2) check the GRQ, (3) check the network poller for ready I/O goroutines, (4) if truly nothing, the M goes idle and parks. This work-stealing design is *why* Go scales near-linearly with cores for CPU-bound goroutine-heavy workloads — no central bottleneck.
- **Blocking syscalls**: when a G makes a blocking syscall (e.g. file I/O), the runtime detaches the M from its P (`handoff`) and either wakes/creates a new M to take that P and keep running other Gs, or lets the P sit idle briefly. This is different from network I/O — Go's **netpoller** (epoll/kqueue/IOCP under the hood) means network-blocked goroutines *don't* block an OS thread at all; they're parked and the runtime is notified via the event loop when the fd is ready, then rescheduled. This is the actual mechanism behind Go handling millions of concurrent network connections cheaply.
- **Preemption**: pre-Go-1.14, the scheduler was purely **cooperative** — a goroutine stuck in a tight loop with no function calls could starve the whole program (a real, embarrassing bug class). Since **Go 1.14**, the runtime uses **asynchronous preemption** via OS signals (SIGURG on Unix) — `sysmon` (a background monitor thread) detects a G running too long (>10ms) and sends a signal to force a preemption check, even mid-loop. This closed a long-standing gap and is a great "do you know recent runtime history" interview flex.
- **sysmon**: a special M running outside the normal GMP scheduling, dedicated to: forcing preemption of long-running Gs, triggering GC at intervals, retaking Ps from Ms blocked too long in syscalls.

### Interview Q&A
- **Q: How many goroutines realistically?** Millions on modern hardware.
- **Q: GOMAXPROCS?** Number of Ps — OS threads that can run Go code simultaneously; defaults to `runtime.NumCPU()`.
- **Q: GMP in one line?** G = work unit, M = OS thread executor, P = scheduling context enabling work-stealing.
- **Q: Does a blocking syscall block all goroutines?** No — M detaches from P; other Gs keep running on other Ms.
- **Staff-level Q: Why doesn't network I/O need this M-detach handoff?** Because network I/O goes through the netpoller (epoll-based) instead of a real blocking syscall from the goroutine's perspective — the G is parked (removed from any M entirely) and the single netpoller thread watches all pending fds, waking the relevant G only when data's ready. Far cheaper than one M per blocked connection.
- **Staff-level Q: What changed in Go 1.14 regarding preemption, and why did it matter?** Cooperative-only preemption (pre-1.14) meant a tight loop with no function calls / channel ops / allocations could indefinitely starve GC and other goroutines on that P — a real production hazard (e.g. tight numeric loops). Async signal-based preemption fixed this at the runtime level without requiring code changes.

---

## 9.9 Bonus — False Sharing, Atomics & Benchmarking (goes beyond the book, but interviewers love it)

### False sharing — the bug most Go devs have never heard of
CPU caches operate in **cache lines** (typically 64 bytes). If two unrelated variables used by *different* goroutines on *different* cores happen to sit in the same 64-byte cache line, writes to one invalidate the other core's cached copy — causing constant cache-coherence traffic even though there's no logical race. This can make "lock-free" atomic-counter code *slower* than expected under high concurrency.

```go
type Counters struct {
    a atomic.Int64 // used by goroutine A
    b atomic.Int64 // used by goroutine B — likely same cache line as 'a'!
}

// Fix: pad to force separate cache lines
type PaddedCounters struct {
    a   atomic.Int64
    _   [56]byte // padding to push b onto its own 64-byte line
    b   atomic.Int64
}
```
Mentioning false sharing unprompted when discussing high-throughput counters/metrics code is a strong signal of production-scale experience.

### sync/atomic — when to reach for it instead of Mutex
Use `atomic.Int64`, `atomic.Bool`, `atomic.Value`, `atomic.Pointer[T]` (typed generics API since Go 1.19) for simple single-variable counters/flags/pointer swaps — cheaper than a mutex (no queueing/parking machinery, just a CPU-level CAS/LOCK-prefixed instruction). Don't reach for atomics to protect multi-field invariants (e.g. "these two fields must update together") — that needs a mutex; atomics only guarantee atomicity per-operation, not across multiple related operations.

### Benchmarking discipline (say this in interviews to sound senior)
```go
func BenchmarkMutex(b *testing.B) {
    var mu sync.Mutex
    var n int
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            mu.Lock()
            n++
            mu.Unlock()
        }
    })
}
```
Always benchmark with `-race` off (for real numbers) and `-cpu=1,2,4,8` to see how a primitive scales with core count — a mutex-protected counter often *degrades* past a certain core count due to contention, while a sharded/atomic approach keeps scaling. This is the kind of answer that separates "knows the syntax" from "has actually profiled Go code in production."

---

## Mastery Checklist — How to Actually Own This Section

1. Explain race condition vs data race without pausing — mention shadow memory / vector clocks if asked how the detector works.
2. Draw the GMP model from memory, including work-stealing and the netpoller — this is asked at nearly every mid-to-senior Go interview.
3. Write the Memo/singleflight pattern (9.7) from scratch in under 10 minutes — one of the most repeated Go coding interview questions.
4. Say Go's proverb — "Don't communicate by sharing memory; share memory by communicating" — then explain *when you still need mutexes anyway* (simple shared state, not full pipelines). Shows depth over memorization.
5. Mention `-race` unprompted whenever discussing concurrent code you write live.
6. Practice spotting the classic `wg.Add()`-inside-goroutine bug (should be called before `go func()`, not inside it — otherwise `Wait()` can return before all goroutines even start).
7. Mutex vs Channel: mutex for protecting shared state (counters/caches, short critical sections); channels for orchestrating goroutine lifecycles, pipelines, ownership transfer.
8. Know Go 1.14's async preemption and Go 1.19's atomic memory-model formalization by version number — shows you track runtime changes, not just syntax.
9. Mention false sharing and cache-line padding when discussing high-throughput atomic counters — very few candidates bring this up unprompted.
10. Be ready to explain *why* `sync.Mutex` has starvation mode and *how* RWMutex prevents writer starvation with the negative-counter trick — both are real "have you read the stdlib source" filters used at Google-style interviews.

---

*Next up in the course: Section 10 — likely Channels & Select, whenever you're ready.*