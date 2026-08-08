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