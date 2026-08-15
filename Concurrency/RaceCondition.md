# 9. Concurrency with Shared Variables — Top 1% Interview Master Sheet

> Every topic has 4 parts: **What it means → A real example → How it works inside (the deep stuff) → Interview questions and answers (easy to hard).** Most candidates don't know this deep stuff. Knowing it is what makes you stand out.

---

## 9.1 Race Conditions

### What it means
A **race condition** happens when two or more goroutines (Go's lightweight workers) touch the same variable at the same time, and at least one of them is writing (changing) the value, without any rule to keep them in order. The final answer depends on pure luck — on which goroutine happens to run first, second, and so on. That's a big problem, because your program should always give the same, correct answer.

### Real-world example
```go
var balance int

func deposit(amount int) {
    balance += amount // NOT one single step: it's really read, add, write = 3 steps
}

func main() {
    for i := 0; i < 1000; i++ {
        go deposit(1)
    }
    time.Sleep(time.Second)
    fmt.Println(balance) // rarely prints 1000 — this is the classic race bug
}
```

### 🔬 How it works inside — what really happens in the machine
`balance += amount` looks like one step in your code. But the computer actually breaks it into three separate steps:
```asm
MOV  EAX, [balance]   ; STEP 1: LOAD — read balance into a fast little storage spot called a register
ADD  EAX, amount      ; STEP 2: ADD — add the amount, but only inside that register, not in memory yet
MOV  [balance], EAX   ; STEP 3: STORE — now save the new number back into memory
```
Each CPU core keeps its own small, fast copy of memory nearby, called a **cache**. Imagine goroutine A reads balance = 5 into its own core's cache. At almost the same moment, goroutine B reads the SAME balance = 5 into a different core's cache. Both of them add 1, so both now have 6. Both save 6 back. One of the two additions is completely lost — and nothing "crashed" or gave an error. This is why race bugs are so sneaky: they depend on tiny timing differences between cores, not on anything that looks broken.

Go has an official rulebook called the **Go Memory Model**. It says clearly: a program is only "race-free" if, for any two operations that touch the same variable (and at least one of them writes to it), we can always say for sure which one happens first. If we can't say that for sure, it's a race. Also, the compiler and the CPU are BOTH allowed to rearrange the order of your instructions, as long as it doesn't change what a single goroutine sees on its own. This is exactly why "it works fine on my computer" is common with buggy race code — different types of CPUs (like ARM vs the common x86) rearrange instructions differently, so a bug that's hidden on one type of computer shows up clearly on another.

### Interview Q&A
- **Q: What is the difference between a race condition and a data race?** A "data race" is the exact technical problem: two goroutines touching the same variable at the same time with no ordering rule, and at least one is a write. Special tools can catch this using tricks like shadow memory and vector clocks (explained more in 9.6). A "race condition" is the bigger, more general problem: your program's result depends on timing. You can have a race condition even without a data race — for example, two well-protected steps that are still done in the wrong order from a logic point of view.
- **Q: Is `balance++` a single safe step in Go?** No — never. No matter what type of number it is, or what computer you're using, this rule never changes.
- **Staff-level Q: Why can a race bug pass 10,000 test runs and then still fail once it reaches production?** Because the tiny window of time where things can go wrong is just a few nanoseconds. Test computers often have fewer processor cores, or different amounts of work happening, compared to the production computers — so the "wrong timing" almost never happens during tests, but shows up later. This is why the `-race` tool (which watches your code as it runs) matters more than just "the tests passed" — it looks for the *possibility* of a race, not only races that actually caused visible failures.

---

## 9.2 Mutual Exclusion: sync.Mutex

### What it means
`sync.Mutex` is like a single key to a room. Only one goroutine can be inside that "room" (called a critical section) at a time. `Lock()` means "I'm taking the key and going in." `Unlock()` means "I'm done, here's the key back." Any goroutine that wants to go in while someone else is inside has to wait outside until the key is free.

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

### 🔬 How Mutex actually works inside (from Go's own source code)
A `Mutex` is really just two numbers hiding under the hood:
```go
type Mutex struct {
    state int32  // packed with info: is it locked? did someone just wake up? is it in "starving" mode? how many are waiting?
    sema  uint32 // a signal the operating system uses to put goroutines to sleep and wake them up
}
```
It actually works in two different modes — and this is something most people don't know:

1. **Normal mode** (the usual case): If `Lock()` can't get in right away, the goroutine doesn't fall asleep immediately. First, it **spins** — it keeps checking, very fast, for a few tries (up to 4), hoping the lock frees up soon. This only happens on computers with more than one core, and only when it seems likely the lock will be free very soon. This little bit of "checking instead of sleeping" makes things much faster when there isn't too much competition for the lock. If spinning doesn't work, the goroutine gets added to a waiting line (first come, first served) and truly goes to sleep, waiting for the operating system to wake it up.
2. **Starving mode**: If a goroutine has been waiting for the lock for **more than 1 millisecond**, the mutex switches into "starving mode." In this mode, new goroutines can no longer cut in line or get lucky — the lock is handed directly, in strict first-come-first-served order, to whoever has been waiting the longest. This exists specifically to stop some unlucky goroutines from waiting forever while others keep jumping ahead — a real bug that happened in Go's very early, simpler mutex design. Once the waiting line is empty, or the wait time drops low again, it switches back to normal mode.

In simple words: when **not many** goroutines want the lock, `sync.Mutex` behaves almost like a super-fast check-and-go lock. When **a lot** of goroutines want the lock, it behaves like a strict, fair, first-come-first-served line. Knowing this — and knowing Go's designers built it this way on purpose — is a strong sign of deep knowledge in an interview.

### Key rules
1. Never copy a mutex (copying resets its inside numbers — pass structs containing a mutex using a pointer, not a copy).
2. A brand-new mutex is already ready to use — you don't need to set it up.
3. Put `Unlock()` inside a `defer` so it still runs even if there's a panic or an early return.
4. Keep the locked section small — don't do slow file/network work, heavy number crunching, or call other functions you don't fully control while holding the lock (it can cause long waits or even a deadlock, where things freeze forever).
5. A mutex is not reentrant — a goroutine can't lock it twice in a row, even if it's the same goroutine.

### Interview Q&A
- **Q: What happens if you call Unlock() on a mutex that isn't locked?** The program panics with the message: "sync: unlock of unlocked mutex."
- **Q: Can a goroutine that already holds the lock, lock it again?** No — it would freeze itself forever (a self-deadlock). Go's Mutex doesn't keep track of "who owns it," unlike some other languages' locks.
- **Q: How do you avoid a deadlock when using several mutexes at once?** Always lock them in the same fixed order, everywhere in your code.
- **Staff-level Q: Why doesn't Go's Mutex keep track of which goroutine owns it?** This was a deliberate choice. Tracking "who owns the lock" would add a small cost to every single Lock/Unlock call — which happens extremely often — just to support a feature (like allowing the same goroutine to lock twice) that Go's designers think you shouldn't rely on anyway. Their philosophy: fix your code's structure instead of leaning on that feature.
- **Staff-level Q: What's the cost difference between locking when nobody else wants the lock, vs. when many goroutines want it at once?** With no competition: it's just one small "compare-and-swap" instruction — a few nanoseconds, extremely fast. With heavy competition: that first attempt fails, there may be some spinning, and then it falls back to a slower operating-system-level wait — this can be a thousand times slower (microseconds instead of nanoseconds). This is why, at large scale, people split one big lock into several smaller ones (for example, splitting one big map into several smaller maps, each with its own lock) instead of using a single global lock.

---

## 9.3 Read/Write Mutexes: sync.RWMutex

### What it means
`sync.RWMutex` is a smarter lock. It allows many goroutines to **read** at the same time, OR exactly one goroutine to **write** at a time (and while it writes, nobody else — not even other readers — can touch the data). Readers use `RLock()`/`RUnlock()`. Writers use `Lock()`/`Unlock()`, just like a normal mutex.

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

### 🔬 How it works inside — the clever trick that stops writers from waiting forever
```go
type RWMutex struct {
    w           Mutex  // held by writers, and also used to block new readers once a writer is waiting
    writerSem   uint32
    readerSem   uint32
    readerCount atomic.Int32 // this number can go NEGATIVE — that's the clever trick
    readerWait  atomic.Int32
}
```
- `RLock()` simply adds 1 to `readerCount` in one quick, safe step. No full lock needed for the normal case — this is extremely cheap and fast.
- When a **writer** calls `Lock()`, it first grabs the internal `w` lock (which blocks other writers from getting in). Then it subtracts a huge number from `readerCount`, which pushes it **into negative numbers**. Now, whenever a *new* reader calls `RLock()`, it checks: if the count after adding 1 is still negative, it knows a writer is waiting, and it blocks itself instead of going ahead. This one simple trick is exactly how Go stops writers from waiting forever while readers keep cutting in line — no extra "someone is waiting" flag is needed.
- The writer then waits only for the readers that were **already in progress** to finish up (tracked by `readerWait` — each finishing reader lowers this count, and the writer proceeds once it hits zero). New readers, as we said, are already blocked.

This negative-number trick is a favorite question in senior-level interviews about RWMutex. Most people only know the basic behavior (many readers, one writer), not *how* Go actually stops writers from starving.

### Interview Q&A
- **Q: When is RWMutex actually worse than a normal Mutex?** When there are lots of writes, or when the locked section is very tiny. RWMutex has extra bookkeeping (special safe number operations and two separate signals) that costs more than it saves in those cases. Rule of thumb: RWMutex is worth it only when reads happen far more often than writes, AND the locked section isn't super tiny.
- **Q: Can a writer be stuck waiting forever (starve)?** No — thanks to the negative-counter trick, new readers block behind a waiting writer instead of cutting in front.
- **Q: Are Go's built-in maps safe to use from multiple goroutines?** No, never by default, if there's any mix of reading and writing. You must add an RWMutex, or use `sync.Map`, yourself.
- **Staff-level Q: Why can `RLock()` still be somewhat costly even when there's no writer around, under heavy load?** Because every `RLock`/`RUnlock` still touches a shared number using a special safe operation. With lots of CPU cores doing this at once, they end up fighting over the same small piece of memory (this is called "false sharing," explained more in section 9.9) — and that fighting can slow things down even without any actual blocking.

---

## 9.4 Memory Synchronization

### What it means
Memory synchronization means making sure that when one goroutine changes some data, other goroutines can actually SEE that change, and see it in the correct order.

Here's the surprising part: even if you don't have a race condition in the strict sense, a goroutine might still not "see" another goroutine's changes properly, because of things like CPU caches and code reordering by the compiler. There needs to be a proper "this happened before that" relationship (called **happens-before**) for changes to be guaranteed visible. Locks and other tools don't just stop two goroutines from bumping into each other — they also guarantee this visibility.

### Real-world example
```go
var ready bool
var data string

func producer() {
    data = "hello"
    ready = true // no synchronization tool used here — the order or visibility of this write is NOT guaranteed
}

func consumer() {
    for !ready { }
    fmt.Println(data) // this might print "" (empty) even though ready looks true
}
```

### 🔬 How it works inside — the official Go Memory Model rules (worth knowing by name)
The [Go Memory Model](https://go.dev/ref/mem) is Go's official rulebook. It lists exactly which situations guarantee that one goroutine's changes are visible to another. The most commonly asked ones are:
1. **Program order** — inside a single goroutine, its own statements happen in the order they're written, but only from that same goroutine's point of view. The compiler is still allowed to secretly reorder things, as long as that one goroutine can't tell the difference.
2. **The `go` statement** — starting a goroutine with `go f()` is guaranteed to happen before `f` actually starts running.
3. **A goroutine finishing** — this does NOT automatically mean other goroutines can see everything it did. You must use a proper tool (like WaitGroup) to be sure. This surprises a lot of people, who wrongly assume "the goroutine is done" is automatically visible everywhere.
4. **A channel send happens before the matching receive finishes** (for unbuffered channels) — this is the whole foundation behind Go's famous idea "share memory by communicating."
5. **For a channel with room for C items, the kth item received happens before the (k+C)th item is fully sent** — this describes the ordering for channels that have some storage space (buffered channels).
6. **Unlocking a `Mutex` happens before the next `Lock()`** on that same mutex finishes — this is exactly WHY mutexes give you visibility, not just "one at a time" protection.
7. **When `sync.Once.Do(f)` finishes running `f` for the first caller, that is guaranteed to happen before ANY other call to `Do` returns** — so every caller is guaranteed to see whatever `f` did.

Why does the `ready bool` example fail? Because there is NO happens-before connection at all between the producer's write and the consumer's read — no channel, no mutex, no atomic tool is used. Go's official rules give **zero guarantees** here — this isn't just "it might be a little slow," it's genuinely broken. In fact, the compiler is technically allowed to read `ready` just once, keep it in a fast register, and loop forever, never even checking memory again — since Go 1.19 made this danger extra clear in the documentation.

### Interview Q&A
- **Q: What does "happens-before" mean?** It's a rule about the order of memory operations, guaranteeing visibility. If A happens-before B, then B is guaranteed to see everything A wrote.
- **Q: Why doesn't a plain true/false flag work for signaling "I'm done"?** Because there's no happens-before connection — nothing stops the compiler or the CPU from reordering or caching things in a way that hides the change.
- **Q: What's the safest way to signal "I'm done" between goroutines?** Closing a channel, using `sync.WaitGroup`, or using `context.Done()`.
- **Staff-level Q: Does `sync/atomic` alone guarantee proper visibility (happens-before), or does it only guarantee the operation itself won't be interrupted halfway (atomicity)?** Since Go 1.19, operations from `sync/atomic` (including newer types like `atomic.Int32`) DO guarantee proper visibility, according to the official memory model — this was less clearly promised in earlier Go versions. Knowing about this 1.19 update shows you've actually read Go's official rules, not just blog posts.

---

## 9.5 Lazy Initialization: sync.Once

### What it means
`sync.Once` guarantees that a piece of setup code runs exactly one time, even if many goroutines try to trigger it at the same time. You call `Do(f)`. Whichever goroutine gets there first actually runs `f`. Everyone else who calls `Do` waits until that first run is finished, and then they all simply return without running `f` again.

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

### 🔬 How it works inside — the "double-checked locking" pattern, done correctly
```go
type Once struct {
    done atomic.Uint32
    m    Mutex
}

func (o *Once) Do(f func()) {
    if o.done.Load() == 0 { // fast path: just one quick, safe read — no full lock at all
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
This is the textbook-correct way to do something called **double-checked locking** — you check first without a lock (fast), and only if needed, check again while holding a lock (safe). The everyday path (`Do`) is just a single quick, safe read, with basically no locking cost, for the huge majority of calls after the very first one. Only the *first* call (or a few goroutines racing to be first) pay the cost of the actual lock. Interviewers love this example because it tests whether you understand why doing this same trick badly (in languages without Go's clear memory rules) is broken in general — but is correct here, specifically because `atomic.Uint32` gives the proper visibility guarantee.

### Interview Q&A
- **Q: What happens if f() panics halfway through?** `done` still gets marked as finished (because of the `defer`) — so `Once` is now permanently "used up," and `f` will never run again, even though it never fully completed. (Note: this `defer o.done.Store(1)` running even during a panic is intentional — some people write their own simpler version and accidentally skip this detail.)
- **Q: Is Once itself safe to use from multiple goroutines?** Yes — that's the entire reason it exists.
- **Q: Why not just write `if !initialized { ... }`?** Because "check, then act" is not one single safe step — this is a classic race bug pattern (called TOCTOU: "time of check to time of use").
- **Staff-level Q: Why is the everyday fast path just a plain safe read, instead of a full Lock()?** For speed. `Once.Do` often gets called on every single request (like fetching a lazily-created singleton object). Making the normal case just one quick, uncontested safe read (a matter of nanoseconds), instead of a full lock (which involves at least one more expensive safe operation), is a deliberate choice to make the common case as cheap as possible.

---

## 9.6 The Race Detector

### What it means

The Race Detector is a Go tool that checks, while your program is actually running, whether multiple goroutines are touching the same piece of data in an unsafe way.

In other words: Go's Race Detector tells you whether your code has a data race problem or not.

Turned on with Go's `-race` flag, it watches memory access as your program runs and can tell you the exact goroutines and exact lines of code that raced on a particular variable.

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

### 🔬 How it actually finds races inside (not commonly known)
The race detector is built on top of Google's **ThreadSanitizer (TSan)** tool, connected to Go through cgo. It uses two main ideas:
1. **Shadow memory** — for every single byte of your program's real memory, there's a hidden "shadow" area tracking extra info: which goroutine touched it last, and at what logical time.
2. **Vector clocks** — every goroutine carries its own logical clock (like a personal timestamp counter). Every time goroutines synchronize with each other (locking/unlocking a mutex, sending/receiving on a channel, using a WaitGroup, and so on), their clocks get merged together, which builds up the proper happens-before order. Then, on every single memory access, the tool checks: "has some other goroutine touched this same spot, in a way that isn't properly ordered compared to me?" If the answer is yes, it reports a race.

This is why the tool only finds races **dynamically** — meaning only on the exact paths your program actually ran during that specific execution. Any concurrent code path that never got triggered during your test stays completely invisible to it. This is also why it's expensive to run: every single memory access now needs an extra shadow-memory check and a clock comparison — which is why programs typically run 2 to 10 times slower, and use 5 to 10 times more memory, with `-race` turned on.

### Interview Q&A
- **Q: Does it catch every race?** No — only the ones that actually happen while that specific run is executing. You need good test coverage of your concurrent code paths (or fuzz testing) for it to be really effective.
- **Q: What's the performance cost?** About 2-10 times slower, and 5-10 times more memory — this is meant for testing and CI only, never for a real production server.
- **Q: Should you use it in CI (automated testing)?** Yes — it's standard practice to run it on every pull request or every nightly build.
- **Staff-level Q: Why can't the race detector just stay turned on all the time in production, as a safety net?** Because the extra checking (the shadow memory and clock comparisons on every single memory access) is simply too expensive to run at real production scale. It's meant as a debugging tool, not a live safety guard. Real production safety comes from good *design* — proper locking tools, careful code review, and running `-race` during testing — not from checking for races live while serving real traffic.

---

## 9.7 Example: Concurrent Non-Blocking Cache

### What it means / the pattern

This pattern stores answers to expensive questions so you don't have to redo the expensive work every time. And when many requests ask for the exact same missing answer at the exact same time, only ONE of them actually does the expensive work — the rest just wait and share that one result.

Cache = "This work is already done, and I've kept the answer."

Singleflight = "This work is already happening right now — don't start the same work again, just wait for it."

#     Singleflight — Simple Meaning

Singleflight is a pattern where, if many goroutines ask for the exact same piece of data at the exact same time, only one of them actually goes and does the real work, while the rest simply wait and then share that same result.

This is a cache where `Get` and `Set` for different pieces of data (keys) never block each other, and where several goroutines asking for the same missing piece of data at the same time won't cause the expensive work to be repeated. This combined idea — remembering results plus avoiding repeated work — is often called the "singleflight" pattern in the Go world. There's even an official, production-ready version of it: `golang.org/x/sync/singleflight`.

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

### Why it's designed this way
1. The lock only protects the small step of looking things up in the map — not the slow work itself. This means different keys never have to wait for each other.
2. The very first goroutine that asks for a given key creates a placeholder entry, lets go of the lock, and THEN does the slow work — without holding up anyone else.
3. If other goroutines ask for the SAME key while that work is happening, they just wait on a channel — no duplicate slow work gets done.
4. `close(e.ready)` wakes up every single one of those waiting goroutines all at the same time.

### 🔬 Why `close()` is the right tool for waking everyone up at once
Closing a channel is Go's special way to let **many goroutines all wake up from one single event at the same moment**, without needing a loop or repeated checking. Once a channel is closed, trying to receive from it (`<-ch`) always instantly returns a zero value — every single goroutine that was waiting wakes up together, in one operation. Compare this to an older, trickier tool called `sync.Cond.Broadcast()`, which needs more careful, manual condition-checking while holding a lock. Closing a `chan struct{}` is the modern, simpler way to do most of what `Cond` used to be needed for.

The official production version (`x/sync/singleflight`) adds a few more things: it recovers safely from a panic (so one broken piece of work doesn't leave everyone else waiting forever), it has a `Forget(key)` method to cancel an in-progress entry, and it tells you whether your result was *shared* with others (meaning duplicate work was avoided). Mentioning this official package by name is a good way to show you know more than just the textbook example.

### Interview Q&A
- **Q: Why not just hold the lock during the whole fetch()?** Because that would force every single request for every single key to wait in line one at a time, destroying all the concurrency benefits.
- **Q: Why use a channel instead of sync.Cond?** It's simpler and much less error-prone for waking up many waiting goroutines at once.
- **Q: What happens if fetch() panics?** In the simple version shown here, all the waiting goroutines would hang forever, because the channel never gets closed. The fix: use `defer close(e.ready)` together with `recover()`, and store the panic as an error result instead.
- **Staff-level Q: How would you add an expiry time (TTL) or cache invalidation to this, without breaking the "no duplicate work" guarantee?** Store an expiry time inside `entry`. When `Get` is called, if the entry has expired, remove it from the map while holding the lock, and treat it like a fresh cache miss (fetch again). But be careful: any goroutine that had already grabbed a reference to the *old* entry, and is still waiting on its channel, should still get that old (now outdated) result rather than getting an error — this keeps the design simple and predictable. Mentioning this small but real detail unprompted is a good sign of production experience.

---

## 9.8 Goroutines and Threads

### What it means
A goroutine is a lightweight worker that Go's own runtime manages — not something the operating system manages directly, like a normal thread. Go can run huge numbers of goroutines using only a small number of real operating system threads, thanks to something called the **M:N scheduler**. This works using the **GMP model**: G (a goroutine — a unit of work), M (an actual OS thread, or "machine"), and P (a logical processor slot, which holds a small local queue of work — an M needs to hold a P before it can run any Go code).

### Comparison table

| Thing being compared | Goroutine | OS Thread |
|---|---|---|
| Starting stack size | About 2KB, and it grows or shrinks automatically as needed | 1-8MB, fixed |
| Who creates it | Go's own runtime (cheap, no need to ask the operating system) | The operating system's kernel (needs a special request, which is expensive) |
| How it's scheduled | Many goroutines share few threads, mixing cooperative pauses with some automatic interruption | One-to-one with real threads, fully controlled by the operating system |
| Cost to switch between them | About tens of nanoseconds, no trip into the operating system needed | About 1-2 microseconds, and it does need to go through the operating system |
| How many you can typically have | Millions | Thousands |
| How they usually talk to each other | Channels (Go's preferred way) | Shared memory plus operating-system tools |

### 🔬 GMP scheduler deep dive (serious, senior-level knowledge)
- **G (Goroutine)**: a small structure holding its own stack, its current position in the code, and its status (things like "ready to run," "currently running," "waiting," "finished," and so on).
- **M (Machine)**: an actual real operating system thread. It can only run Go code while it's holding on to a P.
- **P (Processor)**: a logical slot — the total number of these equals `GOMAXPROCS`. Each P holds a **local queue** of up to 256 pieces of ready-to-run work, and can also reach into a shared **global queue**.
- **Work-stealing**: whenever a P runs out of work in its own local queue, it tries, in this order: (1) "steal" half the work from another random P's queue, (2) check the shared global queue, (3) check if any goroutines are ready because their network request finished, (4) if there's truly nothing to do anywhere, the M goes idle and rests. This work-stealing design is exactly why Go scales up so smoothly with more CPU cores for goroutine-heavy workloads — there's no single bottleneck everyone has to wait on.
- **Blocking system calls**: when a goroutine makes a slow, blocking request directly to the operating system (like reading a file), the runtime detaches its M from its P, and either wakes up (or creates) a different M to take over that P and keep other goroutines running, or lets that P sit idle briefly. Network requests are handled very differently: Go's **netpoller** (built using tools like epoll, kqueue, or IOCP depending on the operating system) means a goroutine waiting on the network does NOT block a whole operating system thread at all — it's simply parked, and the runtime gets notified through an efficient event system once the data is actually ready, and then reschedules that goroutine. This is the real trick behind how Go can cheaply handle millions of open network connections at once.
- **Preemption (interrupting a goroutine that's hogging the CPU)**: before Go version 1.14, the scheduler was purely **cooperative** — meaning a goroutine stuck in a tight loop, with no function calls inside it, could actually block the ENTIRE program from making progress. This was a real, embarrassing type of bug. Since **Go 1.14**, the runtime can use **asynchronous preemption**, using operating system signals (called SIGURG on Unix-based systems). A background helper thread called `sysmon` notices if a goroutine has been running for too long (more than 10 milliseconds) and sends a signal that forces it to pause, even in the middle of a tight loop. This fixed a long-standing gap, and knowing about it is a nice way to show you follow Go's recent history.
- **sysmon**: a special background thread, separate from the normal GMP scheduling, whose job is to: force long-running goroutines to pause, trigger garbage collection at regular intervals, and take back Ps from Ms that have been stuck too long in a blocking system call.

### Interview Q&A
- **Q: Realistically, how many goroutines can you run?** Millions, on modern hardware.
- **Q: What does GOMAXPROCS control?** The number of Ps — meaning the number of operating system threads allowed to run Go code at the exact same time. It defaults to `runtime.NumCPU()` (the number of CPU cores your machine has).
- **Q: Can you explain GMP in one line?** G is a unit of work, M is the real OS thread that runs it, and P is the scheduling slot that makes work-stealing possible.
- **Q: Does one blocking system call block every goroutine?** No — the M simply detaches from its P, and other goroutines keep running on other Ms.
- **Staff-level Q: Why doesn't network I/O need this same M-detaching trick?** Because network operations go through the netpoller (built on something like epoll) instead of being a real blocking request from the goroutine's point of view. The goroutine gets completely parked (removed from any M at all), and a single dedicated netpoller thread watches all the pending network activity, waking up only the specific goroutine whose data is finally ready. This is far cheaper than dedicating one whole OS thread to every single blocked network connection.
- **Staff-level Q: What changed in Go 1.14 about preemption, and why did it matter?** Before 1.14, preemption was only "cooperative," meaning a tight loop with no function calls, channel operations, or memory allocations inside it could block garbage collection and other goroutines on that same P forever — a genuine risk in real programs, like tight number-crunching loops. The new signal-based preemption fixed this at the runtime level, without requiring anyone to change their own code.

---

## 9.9 Bonus — False Sharing, Atomics & Benchmarking (extra material, but interviewers love it)

### False sharing — a bug most Go developers have never even heard of
CPU caches don't work byte by byte — they work in chunks called **cache lines**, usually 64 bytes each. If two totally unrelated variables, used by *different* goroutines running on *different* CPU cores, happen to land inside the same 64-byte cache line, then writing to one of them will force the other CPU core to throw away and reload its cached copy — even though there's no real logical race happening. This can make code that uses fast, "lock-free" atomic counters run *slower* than you'd expect, once you have a lot of concurrent activity.

```go
type Counters struct {
    a atomic.Int64 // used by goroutine A
    b atomic.Int64 // used by goroutine B — probably sitting in the same cache line as 'a'!
}

// Fix: add padding to force them onto separate cache lines
type PaddedCounters struct {
    a   atomic.Int64
    _   [56]byte // unused space, just to push 'b' onto its own 64-byte cache line
    b   atomic.Int64
}
```
Bringing up false sharing on your own, when discussing high-traffic counters or metrics code, is a strong signal that you've worked on real, large-scale systems before.

### sync/atomic — when to reach for it instead of a Mutex
Use `atomic.Int64`, `atomic.Bool`, `atomic.Value`, or `atomic.Pointer[T]` (the newer, typed versions available since Go 1.19) for simple things: a single counter, a single true/false flag, or swapping a single pointer. These are cheaper than a mutex, because there's no waiting-line machinery involved — just one quick, CPU-level safe operation. But don't reach for atomics when you need to keep several related fields in sync together (for example, "these two numbers must always update at the same time") — that genuinely needs a mutex, because atomic tools only guarantee safety for one single operation at a time, not across several related ones.

### Benchmarking discipline (mention this in interviews to sound experienced)
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
Always run your benchmarks with `-race` turned OFF (to get real, honest timing numbers), and try `-cpu=1,2,4,8` to see how a tool performs as you add more cores. A single mutex protecting a shared counter often actually gets *worse* past a certain number of cores, because of competition for the lock — while splitting the work up, or switching to atomics, keeps scaling nicely. Being able to explain this is what separates "knows the syntax" from "has actually measured real Go code in production."

---

## Mastery Checklist — How to Actually Own This Section

1. Be able to explain race condition vs. data race smoothly, without pausing — and mention shadow memory and vector clocks if someone asks how the detector actually works.
2. Be able to draw the GMP model from memory, including work-stealing and the netpoller — this comes up in almost every mid-to-senior level Go interview.
3. Be able to write the Memo/singleflight pattern (from section 9.7) from scratch, in under 10 minutes — this is one of the most repeated Go coding interview questions.
4. Know Go's famous saying — "Don't communicate by sharing memory; share memory by communicating" — and then be able to explain *when you'd still reach for a mutex anyway* (simple shared state, not a full pipeline of steps). This shows real understanding, not just memorized quotes.
5. Mention `-race` on your own, any time you're writing or discussing concurrent code live.
6. Practice spotting the classic bug where `wg.Add()` is called *inside* a goroutine instead of before it (it should be called right before `go func()`, not inside it — otherwise `Wait()` might finish before all the goroutines have even started).
7. Know when to use a Mutex versus a Channel: use a mutex to protect simple shared state (like counters or caches) with small locked sections; use channels to manage the lifecycle of goroutines, build pipelines, or hand off ownership of data.
8. Know Go 1.14's automatic preemption update and Go 1.19's official atomic visibility rules by their version numbers — this shows you actually follow changes to the Go runtime, not just the basic syntax.
9. Bring up false sharing and cache-line padding on your own when talking about high-traffic atomic counters — almost nobody mentions this without being asked.
10. Be ready to explain *why* `sync.Mutex` has a starvation mode, and *how* RWMutex stops writers from starving using the negative-counter trick — both of these are real "have you actually read the standard library's source code" questions used at big tech companies.
