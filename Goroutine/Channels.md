# Channels 
1. What is a Channel?

Definition

    A channel is a typed communication mechanism that allows goroutines to safely send and receive data.

Remember

    Used for communication.
    Used for synchronization.
    Connects two or more goroutines.

2. Why do we use Channels?

Channels solve two problems:

    Communication between goroutines.
    Synchronization between goroutines.

    Instead of sharing variables directly, goroutines exchange data through channels.

3. Channel Characteristics

    A channel is a reference type(A reference type is a data type that stores a memory address pointing to an object rather than holding the actual data directly).

    Created using the make() function.
    Every channel has an element type.
    Zero value of a channel is nil.
    Two channels can be compared using ==.

4. Creating a Channel

ch := make(chan int)

    Creates an unbuffered channel of integers.

5. Channel Operations

    A channel supports only three operations:

Send
    ch <- value

    Sends a value into the channel.

Receive
    
    value := <-ch

Receives a value from the channel.

Close
    close(ch)

Indicates that no more values will be sent.

6. What happens after closing a channel?
    Sending on a closed channel causes a panic.
    Receiving from a closed channel is allowed.
    Once all values are consumed, further receives return the zero value of the element type.


7. Types of Channels

Unbuffered Channel
    make(chan int)

    Capacity = 0

    Communication happens immediately between sender and receiver.

Buffered Channel
    make(chan int, 5)

    Capacity = 5

    Can store values before a receiver reads them.

    Unbuffered Channels
8. What is an Unbuffered Channel?

    An unbuffered channel has no internal storage.

    A send operation and a receive operation must occur together

10. Receive Operation

    A receive operation blocks until another goroutine performs a send.

11. Why are Unbuffered Channels called Synchronous Channels?

    Because:

    Sender waits for Receiver.
    Receiver waits for Sender.
    Both synchronize before continuing.
12. Synchronization

A channel not only transfers data but also synchronizes goroutines.

Communication guarantees that both goroutines meet before proceeding.

13. Happens-Before Relationship

Channel communication establishes a happens-before relationship.

Meaning:

All operations performed before sending on a channel are guaranteed to be visible after the corresponding receive.

This is one of the most important concepts in Go's memory model.

14. Concurrent Events

If there is no synchronization between two goroutines,

their execution order is not guaranteed.

Such operations are said to be concurrent.

15. Channel as a Signal

Sometimes the transmitted value is not important.

The communication itself is important.

Channels are commonly used to signal:

Completion
Cancellation
Notification
Readiness
16. Why chan struct{} is commonly used?

struct{} occupies zero bytes.

It is used when only synchronization is required and no actual data needs to be transferred.

17. Anonymous Goroutine

An anonymous function can be started as a goroutine.

go func() {

}()

Used for short background tasks.

18. Waiting for Goroutines

    A channel can be used to make one goroutine wait until another finishes.

This is a common synchronization pattern.

Interview Definitions

What is a channel?

    A channel is a typed communication mechanism that enables safe communication and synchronization between goroutines.

What is an unbuffered channel?

    An unbuffered channel has no capacity. A send blocks until a receive occurs, and a receive blocks until a send occurs.

What is synchronization?

    Synchronization is the coordination of goroutines so that they execute certain operations in a guaranteed order.

What is the happens-before relationship?

    It guarantees that operations performed before sending on a channel are completed and visible before the corresponding receive continues.

Why use channels instead of shared variables?

    Channels simplify concurrent programming by enabling safe communication and synchronization, reducing the need for explicit locking.

What happens if you send on a closed channel?

    The program panics.

Can you receive from a closed channel?

    Yes. Remaining values are received first; after that, receives immediately return the zero value of the channel's element type.

19. Directional Channels

    A channel can be restricted to send-only or receive-only.

    Send-only Channel
    func send(ch chan<- int) {
        ch <- 10
    }
Can only send data.
Cannot receive.

Receive-only Channel
    func receive(ch <-chan int) {
        value := <-ch
        fmt.Println(value)
    }
    Can only receive data.
    Cannot send.

Interview Question

Q: Why use directional channels?

Answer:

    They improve type safety by preventing unintended send or receive operations.

20. Blocking Behavior

    A channel operation blocks until it can proceed.

    Send blocks
    ch <- 10

    Blocks until another goroutine receives.

    Receive blocks
    x := <-ch

    Blocks until another goroutine sends.

Example
    package main

    import "fmt"

    func main() {
        ch := make(chan int)

        go func() {
            ch <- 100
        }()

        value := <-ch

        fmt.Println(value)
    }

Output

100

21. Deadlock

    A deadlock occurs when all goroutines are waiting forever.

Example

    package main

    func main() {
        ch := make(chan int)

        ch <- 10
    }

Output

fatal error:
all goroutines are asleep - deadlock!

Why?

Because nobody is receiving from the channel.

Another Example

package main

func main() {
    ch := make(chan int)

    <-ch
}

Nobody sends.

Receiver waits forever.

Deadlock.

Interview Definition

Deadlock occurs when goroutines are blocked indefinitely because no goroutine can make progress.

22. Buffered Channel

A buffered channel stores values temporarily.

ch := make(chan int, 3)

Capacity = 3

Example

package main

import "fmt"

func main() {

    ch := make(chan int, 3)

    ch <- 10
    ch <- 20
    ch <- 30

    fmt.Println(<-ch)
    fmt.Println(<-ch)
    fmt.Println(<-ch)
}

Output

10
20
30

Notice

No receiver was needed while sending because the buffer had space.

23. Capacity and Length

Capacity

cap(ch)

Maximum values the channel can hold.

Example

ch := make(chan int, 5)

fmt.Println(cap(ch))

Output

5

Length

len(ch)

Current number of values stored.

Example

ch := make(chan int, 5)

ch <- 10
ch <- 20

fmt.Println(len(ch))

Output

2
24. Iterating Over Channels

    Channels can be used with range.

Example

    package main

    import "fmt"

    func main() {

        ch := make(chan int)

        go func() {
            ch <- 1
            ch <- 2
            ch <- 3
            close(ch)
        }()

        for value := range ch {
            fmt.Println(value)
        }
    }

    Output

    1
    2
    3

Range stops automatically after the channel is closed.

25. Why Close a Channel?

    Closing tells receivers:

    "No more values will be sent."

    Without closing,

    range never stops.

26. Detecting a Closed Channel

Go provides a second return value.

value, ok := <-ch

If

ok == true

Channel is still open.

If

ok == false

Channel is closed.

Example

package main

import "fmt"

func main() {

    ch := make(chan int)

    close(ch)

    value, ok := <-ch

    fmt.Println(value)
    fmt.Println(ok)
}

Output

0
false
27. Unbuffered vs Buffered

| Unbuffered               | Buffered                                 |
| ------------------------ | ---------------------------------------- |
| Capacity = 0             | Capacity > 0                             |
| Sender waits             | Sender waits only when buffer is full    |
| Receiver waits           | Receiver waits only when buffer is empty |
| Used for synchronization | Used for communication and buffering     |

28. Channel Ownership

    Best practice:

    The goroutine that creates the channel should usually close it.

    Receivers generally should not close channels.

29. Communication over Shared Memory

Go Philosophy

Don't communicate by sharing memory; share memory by communicating.

Meaning

Instead of

Multiple goroutines

↓

Same variable

Use

Goroutine

↓

Channel

↓

Goroutine
30. Real-Life Uses of Channels

    Channels are commonly used for

    Worker pools
    Producer-Consumer pattern
    Pipelines
    Background jobs
    Event notifications
    Task scheduling
    Graceful shutdown
    Synchronization
⭐ Must-Know Interview Code

1. Basic Send & Receive

    package main

    import "fmt"

    func main() {

        ch := make(chan int)

        go func() {
            ch <- 100
        }()

        value := <-ch

        fmt.Println(value)
    }

2. Deadlock Example

    package main

    func main() {

        ch := make(chan int)

        ch <- 10
    }
3. Buffered Channel

    package main

    import "fmt"

    func main() {

        ch := make(chan int, 2)

        ch <- 10
        ch <- 20

        fmt.Println(<-ch)
        fmt.Println(<-ch)
    }

4. Range Over Channel

    package main

    import "fmt"

    func main() {

        ch := make(chan int)

        go func() {
            for i := 1; i <= 5; i++ {
                ch <- i
            }
            close(ch)
        }()

        for value := range ch {
            fmt.Println(value)
        }
    }

5. Detect Closed Channel

    package main

    import "fmt"

    func main() {

        ch := make(chan int)

        close(ch)

        value, ok := <-ch

        fmt.Println(value)
        fmt.Println(ok)
    }

Go's philosophy: "Don't communicate by sharing memory; share memory by communicating."
