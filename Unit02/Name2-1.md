Name 
  a name must be begin with the letter or underscore 
  
there are 25 keywords in go lang 
  break default func interface select
case defer go map struct
chan else goto package switch
const fallthrough if range type
continue for import return var 


| Keyword       | Crisp Definition                                                 |
| ------------- | ---------------------------------------------------------------- |
| `break`       | Exits the nearest loop or `switch` immediately.                  |
| `default`     | Runs when no `case` matches in a `switch` or `select`.           |
| `func`        | Declares a function or method.                                   |
| `interface`   | Defines a set of method signatures that types can implement.     |
| `select`      | Waits on multiple channel operations and executes the ready one. |
| `case`        | Represents a branch inside `switch` or `select`.                 |
| `defer`       | Delays a function call until the surrounding function returns.   |
| `go`          | Starts a new goroutine (lightweight concurrent function).        |
| `map`         | Built-in key-value collection with unique keys.                  |
| `struct`      | User-defined type that groups related fields together.           |
| `chan`        | Declares a channel used for communication between goroutines.    |
| `else`        | Executes when the `if` condition is false.                       |
| `goto`        | Jumps to a labeled statement within the same function.           |
| `package`   | Groups related Go files and provides namespaces.                 |
| `switch`      | Selects one execution path among multiple conditions.            |
| `const`       | Declares an immutable compile-time value.                        |
| `fallthrough` | Forces execution of the next `case` in a `switch`.               |
| `if`          | Executes code only when a condition is true.                     |
| `range`       | Iterates over arrays, slices, maps, strings, and channels.       |
| `type`        | Declares a new type or creates a type alias.                     |
| `continue`    | Skips the current loop iteration and moves to the next.          |
| `for`         | The only looping construct in Go.                                |
| `import`      | Brings packages into the current file.                           |
| `return`      | Exits a function and optionally returns values.                  |
| `var`         | Declares a variable whose value can change.                      |


Flow Control: if, else, switch, case, default, fallthrough, break, continue, goto, return
Loops: for, range
Functions: func, defer
Concurrency: go, chan, select
Types: struct, interface, type
Variables: var, const
Collections: map
Organization: package, import 

# 2. Predeclared Constants

| Constant | Definition |
|----------|------------|
| true | Boolean true value. |
| false | Boolean false value. |
| iota | Auto-incrementing constant generator inside a const block. |
| nil | Zero value for pointers, slices, maps, channels, functions, and interfaces. |

---

# 3. Predeclared Types

## Signed Integers

| Type | Definition |
|------|------------|
| int | Signed integer (32 or 64 bits depending on architecture). |
| int8 | 8-bit signed integer. |
| int16 | 16-bit signed integer. |
| int32 | 32-bit signed integer. |
| int64 | 64-bit signed integer. |

---

## Unsigned Integers

| Type | Definition |
|------|------------|
| uint | Unsigned integer (32 or 64 bits depending on architecture). |
| uint8 | 8-bit unsigned integer. |
| uint16 | 16-bit unsigned integer. |
| uint32 | 32-bit unsigned integer. |
| uint64 | 64-bit unsigned integer. |
| uintptr | Integer large enough to hold a pointer value. |

---

## Floating-Point Types

| Type | Definition |
|------|------------|
| float32 | 32-bit floating-point number. |
| float64 | 64-bit floating-point number. |

---

## Complex Types

| Type | Definition |
|------|------------|
| complex64 | Complex number with float32 components. |
| complex128 | Complex number with float64 components. |

---

## Other Built-in Types

| Type | Definition |
|------|------------|
| bool | Represents true or false. |
| byte | Alias for uint8; stores raw bytes. |
| rune | Alias for int32; represents a Unicode code point. |
| string | Immutable sequence of UTF-8 bytes. |
| error | Built-in interface representing an error. |

---

# 4. Predeclared Functions

## Memory & Initialization

| Function | Definition |
|----------|------------|
| make() | Initializes slices, maps, and channels. |
| new() | Allocates memory and returns a pointer to the zero value. |

---

## Collections

| Function | Definition |
|----------|------------|
| append() | Appends elements to a slice. |
| copy() | Copies elements from one slice to another. |
| delete() | Removes a key from a map. |
| clear() | Removes all elements from a map or zeroes a slice. |

---

## Size & Capacity

| Function | Definition |
|----------|------------|
| len() | Returns the number of elements. |
| cap() | Returns the capacity of a slice, array, or channel. |

---

## Channels

| Function | Definition |
|----------|------------|
| close() | Closes a channel so no more values can be sent. |

---

## Complex Numbers

| Function | Definition |
|----------|------------|
| complex() | Creates a complex number. |
| real() | Returns the real part of a complex number. |
| imag() | Returns the imaginary part of a complex number. |

---


## Error Handling

| Function | Definition |
|----------|------------|
| panic() | Stops normal execution and starts stack unwinding. |
| recover() | Recovers from a panic inside a deferred function. |

---

## Debugging

| Function | Definition |
|----------|------------|
| print() | Prints values for debugging. |
| println() | Prints values followed by a newline for debugging. |

---

## Comparison (Go 1.21+)

| Function | Definition |
|----------|------------|
| min() | Returns the smallest value. |
| max() | Returns the largest value. |

---

