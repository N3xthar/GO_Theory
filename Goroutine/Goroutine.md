# 1. Goroutine 

1. What is a Goroutine?

    A goroutine is a lightweight concurrent function in Go.

    It allows multiple functions to execute concurrently.
    
    Managed by the Go runtime, not directly by the operating system.

Interview Definition:

    A goroutine is a lightweight thread of execution managed by the Go runtime.

# 2. Main Goroutine 

Every Go program starts with one goroutine .

    func main() {
    // This runs in the main goroutine.
    }

It is called the main goroutine because it executes the main() function.

# 3. Creating a Goroutine

Normal function call:

    f()
    Executes synchronously.
    Caller waits until f() finishes.

Goroutine:

    go f()
    
    Creates a new goroutine.
    Executes concurrently.
    Caller does not wait.

# 4 . Syntax 

    go functionname(parameter)

# 5. Concurrent Execution

Suppose we have:

    go spinner()
    fib(45)

Now,

    spinner() runs in one goroutine.
    fib(45) runs in the main goroutine.

    Both execute at the same time (concurrently).

# 6 What happens when main() finishes?
    
    func main(){

    }

returns 
    the program executes 
    all the goroutine are terminated immediately , even if they havent finished the task or not 

# 7 can one goroutine stop the another ? 

No. It should communicate a stop request (typically via channels or context) so the other goroutine can exit gracefully.


Instead, goroutines communicate using mechanisms like:

    channels
    context cancellation
    shared signals

The goroutine should stop itself after receiving the signal.

# 8 Goroutine vs Thread 
| Thread          | Goroutine                                         |
| --------------- | ------------------------------------------------- |
| Managed by OS   | Managed by Go runtime                             |
| Heavyweight     | Lightweight                                       |
| More memory     | Very small initial stack (~2 KB, grows as needed) |
| Slower creation | Faster creation                                   |
| Fewer threads   | Thousands to millions of goroutines possible      |


Synchronization means coordinating multiple goroutines so they execute certain operations in the correct order.


# 9 Concurrency  + Networking 
Why is concurrency important in networking?
    A server usually serves multiple clients simultaneously.
    Each client connection is independent.
    Concurrency prevents one slow client from blocking others.

Interview Answer:

    Concurrency allows a server to handle multiple client requests simultaneously, improving responsiveness and scalability.

# 10.3 what is the net packages ? 

Standard Go package for networking.
    
    Supports:
    TCP
    UDP
    Unix Domain Sockets

net/http is built on top of net.


3. What does net.Listen() do?
    Starts a server.
    Opens a network port.
    Returns a net.Listener.
    listener, err := net.Listen("tcp", "localhost:8000")

4. What is net.Listener?
    Represents a listening server.
    Waits for incoming client connections.

Important method:

    listener.Accept()

5. What does Accept() do?
    Waits for a client to connect.
    Returns a net.Conn.

Important: Accept() is a blocking call.

6. What is net.Conn?

Represents a TCP connection between client and server.

Implements:

    io.Reader
    io.Writer

    So you can:

    Read()
    Write()
    io.Copy()
    io.WriteString()

    directly on it.

7. Sequential Server

    Without goroutines:

    handleConn(conn)

    Only one client is handled at a time.

    Other clients must wait.

8. Concurrent Server

    go handleConn(conn)

    Each client gets its own goroutine.

    Clients are handled simultaneously.

    This is the standard Go server pattern.

9. One Goroutine Per Connection

Most Go servers follow:

    Client 1
        ↓
    Goroutine 1

    Client 2
        ↓
    Goroutine 2

    Client 3
        ↓
    Goroutine 3

    Benefits:

    High concurrency
    Better scalability
    Better performance

10. Why use defer conn.Close()?

    Ensures the connection is closed automatically when the function exits.

    Prevents resource leaks.

11. What is net.Dial()?

    Used by the client.

    conn, err := net.Dial("tcp", "localhost:8000")

    Connects to a remote server.

12. What is io.Copy()?

        Copies data from one stream to another until:

        EOF
        Error

        Example:

        io.Copy(os.Stdout, conn)

        Server → Terminal

13. What is io.WriteString()?

    Writes a string to any io.Writer.

    Since net.Conn implements io.Writer, we can write directly to the client.

14. What is bufio.Scanner?

    Used to read input line by line.

    scanner := bufio.NewScanner(conn)

    Commonly used for:

    Chat servers
    Echo servers
    Command-line protocols
15. What is an Echo Server?

    An echo server sends back whatever it receives.

    Client:

    Hello

    Server:

    Hello

    Used to demonstrate network communication.

16. Sequential Echo Server
    echo(message)

    Processes one message at a time.

17. Concurrent Echo Server
    go echo(message)

    Each message runs in a separate goroutine.

    Messages overlap in execution.

18. Goroutines can be used at multiple levels

    One goroutine per client connection.
    One goroutine per request/message inside a connection.

    This demonstrates Go's flexible concurrency model.


19. Blocking Operations

Blocking means the current goroutine stops executing until the operation finishes.

    Accept()
    Scanner.Scan()
    io.Copy()
    time.Sleep()

20. Function arguments in goroutines
go echo(c, input.Text(), time.Second)

    Before starting the goroutine, Go evaluates the function arguments and passes their values to the goroutine. Because the goroutine receives a copy of those values, any changes made to the original variables later do not affect the values that were already passed to the goroutine.

21. Concurrency Safety

Not every type in Go is safe for concurrent access.

net.Conn is designed to allow concurrent reads and writes.

For shared data structures (maps, slices, variables), use synchronization such as mutexes or channels when needed.
| API                  | Purpose                        |
| -------------------- | ------------------------------ |
| `net.Listen()`       | Start TCP server               |
| `net.Listener`       | Listen for clients             |
| `Accept()`           | Accept incoming client         |
| `net.Conn`           | TCP connection                 |
| `net.Dial()`         | Connect to server              |
| `io.Copy()`          | Copy stream data               |
| `io.WriteString()`   | Write string                   |
| `bufio.NewScanner()` | Read line by line              |
| `defer Close()`      | Close connection automatically |
| `go`                 | Create a goroutine             |

Q2. What is net.Conn?

    A TCP connection that implements both io.Reader and io.Writer.

Q3. Is Accept() blocking?

    Yes. It waits until a client connects.

Q4. Difference between net.Listen() and net.Dial()?
    Listen() → Server
    Dial() → Client
Q5. Why do we write go handleConn(conn)?

    To process each client in a separate goroutine.

Q6. Why use defer conn.Close()?

    To ensure the connection is closed even if the function returns because of an error.

Q7. What is the difference between the clock server and the echo server?

Clock server: Sends the current time periodically to connected clients.
Echo server: Reads data from a client and sends the same data back.
One-line Summary

    Go networking follows a simple pattern: net.Listen() creates a server,
     Accept() waits for clients, each client is represented by a net.Conn, and using go handleConn(conn) allows the server to handle many clients concurrently with lightweight goroutines.
