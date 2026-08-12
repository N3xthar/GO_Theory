# Go Concurrency — Cancellation & Goroutine Leaks (Simple Notes)

## 1. Can one goroutine kill another goroutine?

**Answer: No.**

**Why not?**
Agar a goroutine beech mein force-kill ho jaye (jaise DB write ke beech, ya lock hold karte waqt), to data half-written reh sakta hai aur locks kabhi release nahi honge. System **inconsistent** ho jayega.

**So what does Go do instead?**

Go follows **cooperative cancellation** — main goroutine sirf ek "please stop" signal bhejta hai. Worker khud us signal ko check karke gracefully exit karti hai.

    Main → "please stop" (signal only)
    Worker → sees signal → cleans up → returns

# Go Concurrency — Cancellation (du4) & Chat Server (Simple Notes)

## PART A: The `du4` Cancellation Pattern (Poll-based)

### 1. The Core Problem

Tera program directory sizes calculate kar raha hai (`du` command jaisa), aur bahut saari goroutines chal rahi hain parallel mein. User beech mein **Enter** dabata hai — "bas, ruk jao."

**Why is this hard?**
- Ek goroutine doosri goroutine ko force-kill nahi kar sakti (shared state corrupt ho sakta hai)
- Agar tu `abort` channel pe values bhejta rahe, to kitni baar bhejni hain? Kitni goroutines chal rahi hain — pata nahi!
- Agar ek goroutine ne value le li, doosri ko wo value nahi dikhegi (channel se ek receive = ek consume)

**So the real need:** Ek signal jo **saath mein, baar-baar, sabko** dikhe — chahe abhi 5 goroutines chal rahi ho ya 500.

---

### 2. The Trick: `close(done)` Instead of Sending

```go
var done = make(chan struct{})

func cancelled() bool {
    select {
    case <-done:
        return true
    default:
        return false
    }
}
```

**Why this works (the actual mechanism):**
Jab channel **close** ho jata hai, uske baad koi bhi receive turant "ready" ho jata hai — **har baar**, **har goroutine ke liye**, forever.

    close(done)
    |
    Ab jitni baar bhi, jo bhi <-done kare
    |
    turant "true / ready" mil jayega

**`cancelled()` function kya karta hai?**
Ye ek **non-blocking check** hai — `default` ke saath `select`. Matlab:
- Agar `done` close ho chuka hai → `true` return
- Agar nahi hua → turant `false` return, wait nahi karta

Ye function goroutine kahin bhi call karke pooch sakti hai: **"Kya humein rukna hai?"** — bina block hue.

---

### 3. Two Ways to Listen for Cancellation

**Way 1: Poll at the start (cheap check before starting work)**
```go
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
    defer n.Done()
    if cancelled() {
        return   // cancel ho chuka? to shuru hi mat kar
    }
    // ... normal work ...
}
```
**Why:** Naye kaam **start hi na ho** agar cancellation already ho chuki hai. Isse naye goroutines banna band ho jate hain.

**Way 2: Listen inside a `select` (react immediately while waiting)**
```go
for {
    select {
    case <-done:
        // drain fileSizes so other goroutines don't get stuck sending
        for range fileSizes {
        }
        return
    case size, ok := <-fileSizes:
        // normal processing
    }
}
```

**Why "drain fileSizes" is important:**
Kuch goroutines abhi bhi `fileSizes` channel pe **bhejne** ki koshish kar rahi ho sakti hain. Agar tu turant `return` kar de bina channel khali kiye, to wo goroutines **hamesha ke liye atak jayengi** (send block ho jayega, koi receive nahi karega) → **goroutine leak**.

Isliye pehle drain karo (sab values discard karo), phir return karo — taaki koi bhi pending goroutine safely finish ho sake.

---

### 4. Cancelling a Blocking Operation (Semaphore)

```go
func dirents(dir string) []os.FileInfo {
    select {
    case sema <- struct{}{}: // token lene ki koshish
    case <-done:
        return nil // cancel ho gaya, wait mat karo
    }
    defer func() { <-sema }()
    // ...
}
```

**Why needed?**
Agar saare semaphore tokens already use ho rahe hain, to `sema <- struct{}{}` **block** ho jayega — goroutine wahi phas jayegi, cancellation signal aane ke baad bhi.

`select` ismein ek **exit door** add kar deta hai: "token milne ka wait karo, YA cancel signal ka — jo pehle aaye."

**Real impact (from the book):** Isse cancellation ka response time **sekdo milliseconds se ghatke sirf 10s of milliseconds** ho gaya. Matlab — jahan bhi tera code kisi cheez ka wait kar raha hai, wahan `select` + `done` daal do, taaki wo jagah bhi turant cancel ho sake.

---

### 5. Why "Poll in a Few Key Places" is Enough

Har single line ke baad cancellation check karna practically possible nahi hai. Book ka point:

> Cancellation is a trade-off — faster response = more code changes. Lekin usually **kuch important jagah** check karna (loop start, blocking calls) hi zyada tar fayda de deta hai.

**Interview-friendly line:**
> "Full cancellation coverage requires checking everywhere, which is expensive. In practice, checking at the start of expensive operations and at blocking points (like semaphore acquisition) captures most of the benefit."

---

## PART B: Chat Server — `select` Handling 3 Event Types

### 6. The Real-World Problem

Ek chat server hai jahan **multiple clients** connect hote hain. Server ko ek saath 3 cheezein handle karni hain:
1. Koi naya client aaya (**entering**)
2. Koi client chala gaya (**leaving**)
3. Koi message aaya jo sabko bhejna hai (**messages**)

Ye teeno **kabhi bhi, kisi bhi order mein** ho sakte hain. Single-threaded logic se ye handle karna mushkil hai.

### 7. The `broadcaster` Goroutine — The Heart of the Server

```go
func broadcaster() {
    clients := make(map[client]bool)
    for {
        select {
        case msg := <-messages:
            for cli := range clients {
                cli <- msg
            }
        case cli := <-entering:
            clients[cli] = true
        case cli := <-leaving:
            delete(clients, cli)
            close(cli)
        }
    }
}
```

**Why `select` here?**
Ek hi goroutine (`broadcaster`) teen alag channels ko **ek saath** sun raha hai. Jo bhi event pehle aaye, `select` usi case ko turant chala deta hai.

**Why does this need NO locks (`sync.Mutex`)?**
Kyunki `clients` map ko sirf **ek hi goroutine** (broadcaster) touch karta hai — koi doosra goroutine directly isse access nahi karta. Sab kuch channels ke through hota hai. Isko **"confinement"** kehte hain — data ek hi goroutine tak confined (band) hai, isliye race condition ka risk hi nahi.

**Interview gold line:**
> "The clients map is only ever accessed by the broadcaster goroutine, so there's no need for locks — this is called confinement."

---

### 8. Per-Client Goroutines: `handleConn` + `clientWriter`

Har client ke liye **2 goroutines** banti hain:

**`cancelled()` function kya karta hai?**
Ye ek **non-blocking check** hai — `default` ke saath `select`. Matlab:
- Agar `done` close ho chuka hai → `true` return
- Agar nahi hua → turant `false` return, wait nahi karta

Ye function goroutine kahin bhi call karke pooch sakti hai: **"Kya humein rukna hai?"** — bina block hue.

---

### 3. Two Ways to Listen for Cancellation

**Way 1: Poll at the start (cheap check before starting work)**
```go
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
    defer n.Done()
    if cancelled() {
        return   // cancel ho chuka? to shuru hi mat kar
    }
    // ... normal work ...
}
```
**Why:** Naye kaam **start hi na ho** agar cancellation already ho chuki hai. Isse naye goroutines banna band ho jate hain.

**Way 2: Listen inside a `select` (react immediately while waiting)**
```go
for {
    select {
    case <-done:
        // drain fileSizes so other goroutines don't get stuck sending
        for range fileSizes {
        }
        return
    case size, ok := <-fileSizes:
        // normal processing
    }
}
```

**Why "drain fileSizes" is important:**
Kuch goroutines abhi bhi `fileSizes` channel pe **bhejne** ki koshish kar rahi ho sakti hain. Agar tu turant `return` kar de bina channel khali kiye, to wo goroutines **hamesha ke liye atak jayengi** (send block ho jayega, koi receive nahi karega) → **goroutine leak**.

Isliye pehle drain karo (sab values discard karo), phir return karo — taaki koi bhi pending goroutine safely finish ho sake.

---

### 4. Cancelling a Blocking Operation (Semaphore)

```go
func dirents(dir string) []os.FileInfo {
    select {
    case sema <- struct{}{}: // token lene ki koshish
    case <-done:
        return nil // cancel ho gaya, wait mat karo
    }
    defer func() { <-sema }()
    // ...
}
```

**Why needed?**
Agar saare semaphore tokens already use ho rahe hain, to `sema <- struct{}{}` **block** ho jayega — goroutine wahi phas jayegi, cancellation signal aane ke baad bhi.

`select` ismein ek **exit door** add kar deta hai: "token milne ka wait karo, YA cancel signal ka — jo pehle aaye."

**Real impact (from the book):** Isse cancellation ka response time **sekdo milliseconds se ghatke sirf 10s of milliseconds** ho gaya. Matlab — jahan bhi tera code kisi cheez ka wait kar raha hai, wahan `select` + `done` daal do, taaki wo jagah bhi turant cancel ho sake.

---

### 5. Why "Poll in a Few Key Places" is Enough

Har single line ke baad cancellation check karna practically possible nahi hai. Book ka point:

> Cancellation is a trade-off — faster response = more code changes. Lekin usually **kuch important jagah** check karna (loop start, blocking calls) hi zyada tar fayda de deta hai.

**Interview-friendly line:**
> "Full cancellation coverage requires checking everywhere, which is expensive. In practice, checking at the start of expensive operations and at blocking points (like semaphore acquisition) captures most of the benefit."

---

## PART B: Chat Server — `select` Handling 3 Event Types

### 6. The Real-World Problem

Ek chat server hai jahan **multiple clients** connect hote hain. Server ko ek saath 3 cheezein handle karni hain:
1. Koi naya client aaya (**entering**)
2. Koi client chala gaya (**leaving**)
3. Koi message aaya jo sabko bhejna hai (**messages**)

Ye teeno **kabhi bhi, kisi bhi order mein** ho sakte hain. Single-threaded logic se ye handle karna mushkil hai.

### 7. The `broadcaster` Goroutine — The Heart of the Server

```go
func broadcaster() {
    clients := make(map[client]bool)
    for {
        select {
        case msg := <-messages:
            for cli := range clients {
                cli <- msg
            }
        case cli := <-entering:
            clients[cli] = true
        case cli := <-leaving:
            delete(clients, cli)
            close(cli)
        }
    }
}
```

**Why `select` here?**
Ek hi goroutine (`broadcaster`) teen alag channels ko **ek saath** sun raha hai. Jo bhi event pehle aaye, `select` usi case ko turant chala deta hai.

**Why does this need NO locks (`sync.Mutex`)?**
Kyunki `clients` map ko sirf **ek hi goroutine** (broadcaster) touch karta hai — koi doosra goroutine directly isse access nahi karta. Sab kuch channels ke through hota hai. Isko **"confinement"** kehte hain — data ek hi goroutine tak confined (band) hai, isliye race condition ka risk hi nahi.

**Interview gold line:**
> "The clients map is only ever accessed by the broadcaster goroutine, so there's no need for locks — this is called confinement."

---

### 8. Per-Client Goroutines: `handleConn` + `clientWriter`

    Har client ke liye **2 goroutines** banti hain:

    Client connects
    ↓
    handleConn goroutine → reads messages FROM this client
    ↓
    clientWriter goroutine → sends messages TO this client

**Why 2 separate goroutines per client, not 1?**
- `handleConn` **client se padhta hai** (blocking read)
- `clientWriter` **client ko likhta hai** (blocking write)

Agar dono kaam ek hi goroutine mein hote, to jab tak read ho raha hai, tab tak write nahi ho sakta — server client ko real-time messages nahi bhej payega jab wo khud kuch type kar raha ho. **Alag goroutines = simultaneous read aur write.**

```go
func handleConn(conn net.Conn) {
    ch := make(chan string)
    go clientWriter(conn, ch)          // separate writer goroutine

    who := conn.RemoteAddr().String()
    ch <- "You are " + who
    messages <- who + " has arrived"    // tell broadcaster
    entering <- ch                      // register this client

    input := bufio.NewScanner(conn)
    for input.Scan() {
        messages <- who + ": " + input.Text()   // forward every line
    }

    leaving <- ch                       // tell broadcaster client left
    messages <- who + " has left"
    conn.Close()
}
```

**How does the writer goroutine know when to stop?**
```go
func clientWriter(conn net.Conn, ch <-chan string) {
    for msg := range ch {
        fmt.Fprintln(conn, msg)
    }
}
```
`for msg := range ch` — jab tak `ch` **close** nahi hota, ye loop chalta rehta hai. Broadcaster jab `leaving` case handle karta hai, tab `close(cli)` karta hai — isi se writer goroutine ka loop **automatically** khatam ho jata hai. Koi manual signal nahi bhejna padta.

---

### 9. Why This Design Needs "No Explicit Locking"

**n clients ke liye = `2n + 2` goroutines** (har client ke 2, plus `main` aur `broadcaster`).

Phir bhi koi `Mutex` nahi use hua kyunki:
- `clients` map sirf broadcaster ke andar hai (confined)
- Baaki sab sirf **channels** ke through communicate karte hain (jo khud thread-safe hote hain)

**Interview line:**
> "The chat server handles concurrency safely without locks because shared state is confined to a single goroutine, and all communication happens over channels — which are inherently concurrency-safe."

---

## 10. Quick Comparison Table

| Pattern | Used For | Key Mechanism |
|---|---|---|
| `close(done)` + poll (`cancelled()`) | Cancel many goroutines with one signal | Closed channel = always-ready receive |
| `select` with `done` inside loop | React to cancellation while waiting for work | Extra exit-door in select |
| `select` on semaphore + `done` | Cancel even while blocked on acquiring a resource | Prevents stuck blocking calls |
| Drain before return | Avoid leaking goroutines still trying to send | Empty the channel before exiting |
| `broadcaster` with `select` on 3 channels | Coordinate multiple event types in one goroutine | One case per event type |
| Confinement (no locks) | Safe concurrent access to shared state | Only 1 goroutine touches the data directly |

---

## 11. Interview-Ready Answers 

**Q: How do you cancel an arbitrary number of goroutines?**
> "Close a shared 'done' channel instead of sending values. Closing broadcasts to every goroutine listening, and every future receive on a closed channel returns immediately — so it works regardless of how many goroutines are involved."

**Q: What's the danger if you return immediately on cancellation without draining a channel?**
> "Other goroutines still trying to send on that channel would block forever, causing a goroutine leak. You should drain the channel before returning."

**Q: How do you make a blocking operation (like acquiring a semaphore) cancellable?**
> "Wrap it in a select with a case for acquiring the resource and a case for the done channel — whichever is ready first wins, so cancellation can interrupt the wait."

**Q: How does the chat server avoid using locks for shared state?**
> "The shared clients map is confined to a single goroutine (the broadcaster). All other goroutines only interact with it indirectly through channels, which are concurrency-safe by design."

**Q: Why does each client get two goroutines instead of one?**
> "One goroutine (handleConn) reads incoming messages from the client, and another (clientWriter) writes outgoing messages to the client. Separating them allows simultaneous reading and writing on the same connection."