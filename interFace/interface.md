# 7. Interfaces — Complete Go Study Guide (Interview + Production Ready)

> Goal: After reading this, you should be able to say —
> "I know WHY interfaces exist, WHAT problem they solve, HOW Go implements them internally, HOW the standard library uses them, and HOW to explain all this in an interview."

---

## How to use this file

Each topic below follows the same pattern:

`WHAT → WHY → PROBLEM SOLVED → HOW IT WORKS → CODE → REAL-LIFE USE → WHEN TO USE / NOT USE → MISTAKES → GOTCHAS → STDLIB → INTERVIEW Qs → QUICK REVISION`

Read one topic at a time. Do the practice exercise before moving to the next.

---

# 7.1 Interfaces as Contracts

## 1. What is it?

```text
An interface is a CONTRACT.

The contract says:
"If your type has these methods, you are allowed to be used here."
```

An interface does **not** care what the type *is* (a struct, a map, an int...). It only cares what the type **can do** (its methods).

## 2. Why do we need it?

Without interfaces, every function would have to say "give me exactly this struct." That means your code becomes glued to one specific type forever.

With interfaces, functions say "give me anything that can do X." This is called **behavior-based programming** instead of **type-based programming**.

## 3. What problem does it solve?

```text
Without interface:

PrintReport(db *PostgresDB)

Now PrintReport can NEVER work with MySQL, or a mock database for testing.
It is stuck forever with PostgresDB.
```

```text
With interface:

PrintReport(db DataSource)

Now PrintReport works with PostgreSQL, MySQL, a test mock,
or anything else that satisfies DataSource.
```

This solves the problem of **tight coupling** — code being permanently stuck to one specific implementation.

## 4. How does it work?

```text
Step 1: Define an interface — a list of method names + signatures.
Step 2: Any type that has ALL those methods automatically satisfies the interface.
Step 3: Go does this check at COMPILE TIME, not runtime.
Step 4: No "implements" keyword needed — it's automatic ("structural typing").
```

## 5. Simple Mental Model

```text
Interface = a checklist of abilities

If a type has every ability on the checklist,
Go allows that type to be used wherever the interface is expected.
```

## 6. Basic Go Example

```go
package main

import "fmt"

type Speaker interface {
    Speak() string
}

type Dog struct{}

func (Dog) Speak() string { return "Woof" }

type Robot struct{}

func (Robot) Speak() string { return "Beep boop" }

func announce(s Speaker) {
    fmt.Println(s.Speak())
}

func main() {
    announce(Dog{})
    announce(Robot{})
}
```

## 7. Explain the Code

```text
1. Speaker interface requires one method: Speak() string.
2. Dog has Speak() -> Dog satisfies Speaker.
3. Robot has Speak() -> Robot satisfies Speaker.
4. announce() does not know or care about Dog or Robot.
5. It only knows "I got something that can Speak()".
6. This is polymorphism: different types, same interface, different behavior.
```

**Polymorphism** simply means: many different types can be used in the same place, as long as they share the same behavior (methods).

## 8. Real-Life Problem

```text
A backend service needs to send notifications.
Today it's email. Tomorrow it might be SMS or Slack.

Business Logic
      |
      v
Notifier interface  { Send(msg string) error }
      |
      +-- EmailNotifier
      +-- SMSNotifier
      +-- SlackNotifier
```

Business logic never changes when you add a new notification channel — you just add a new type that satisfies `Notifier`.

## 9. When should I use it?

- When a piece of code needs to work with **multiple** possible implementations (real DB vs test mock, multiple payment gateways, etc).
- When you want to swap implementations without touching the calling code.
- When you're designing a boundary between packages (e.g., between your business logic and external systems).

## 10. When should I NOT use it?

- **Do not** create an interface if there is only ever going to be **one** implementation. That's just extra indirection with no benefit.
- Do not create interfaces "just in case you need them later." Go's advice: **create interfaces where they are used (consumer side), not where the type is defined.**

> Rule of thumb: "Accept interfaces, return structs." Functions should accept interfaces (flexible input), but usually return concrete types (clear output).

## 11. Common Mistakes

- Creating a big interface with 10 methods when only 1 method is actually needed by the caller.
- Defining the interface next to the concrete type (in Go, it's usually better to define it next to the code that *uses* it).
- Thinking Go needs an `implements` keyword — it doesn't. Satisfaction is automatic.

## 12. Important Gotchas

- Interface satisfaction is checked at **compile time** — if a method is missing, your code won't build.
- A type can satisfy many interfaces at once, with no declaration needed anywhere.

## 13. Internals

### Go Language Guarantee
Any type (not just structs) that implements all methods of an interface satisfies that interface. This includes named basic types, maps, slices, funcs, etc.

### Implementation Detail
Internally, an interface value is stored as a pair: (type, value). We cover this fully in section 7.5 (Interface Values).

## 14. Standard Library Connection

```text
io.Reader     -> Read(p []byte) (n int, err error)
io.Writer     -> Write(p []byte) (n int, err error)
error         -> Error() string
fmt.Stringer  -> String() string
```

These tiny interfaces are used **everywhere** in Go's standard library.

## 15. Production Example

```go
type UserRepository interface {
    FindByID(id int64) (*User, error)
    Save(user *User) error
}

type service struct {
    repo UserRepository
}

func (s *service) GetUser(id int64) (*User, error) {
    return s.repo.FindByID(id)
}
```

`service` doesn't know if `repo` is backed by Postgres, MySQL, or an in-memory map used in unit tests. That's the power of the contract.

## 16. Performance

Interfaces have a small performance cost (an extra pointer indirection, and sometimes a heap allocation — see 7.5). For 99% of backend code this cost is irrelevant. Only worry about it in extremely hot loops, and measure before optimizing.

## 17. Related Concepts

| Concept | Meaning |
|---|---|
| Interface | The contract (list of required methods) |
| Concrete type | The real type (Dog, Robot) that implements the methods |
| Interface value | A "box" holding a concrete type + its value at runtime |

## 18. Interview Questions

**Basic**
- Q: What is an interface in Go?
- A: A set of method signatures. Any type that implements all those methods automatically satisfies the interface — no explicit declaration needed.

**Intermediate**
- Q: How is Go's interface satisfaction different from Java/C#?
- A: Go uses **structural typing** (implicit) — no `implements` keyword. Java/C# use **nominal typing** (explicit) — you must declare `implements InterfaceName`.

**Advanced**
- Q: Why does Go encourage small interfaces?
- A: Small interfaces (1-2 methods) are easier to satisfy, easier to mock in tests, and force you to depend only on the behavior you actually need (Interface Segregation Principle).

**Tricky**
- Q: Can a non-struct type satisfy an interface?
- A: Yes. Any type — int, string, map, func, slice — can have methods and thus satisfy an interface, as long as methods are defined on it.

## 19. Interview Follow-Up Questions

```text
Q: What is an interface?
Q: How does a type satisfy an interface?
Q: Does Go have an "implements" keyword?
Q: Where should interfaces be defined — producer or consumer side?
Q: Why does Go prefer small interfaces?
```

## 20. How to Explain This in an Interview

> "In Go, an interface is a contract made up of method signatures. Any type that implements all of those methods automatically satisfies the interface — there's no explicit 'implements' keyword like in Java. This lets us write functions that depend on behavior rather than a specific concrete type, which makes code more flexible and easier to test, because we can swap in mocks or different implementations without changing the calling code."

## 21. Quick Revision

```text
WHAT?        -> A contract: a list of required methods.
WHY?         -> Decouple code from specific implementations.
PROBLEM?     -> Tight coupling to one concrete type.
HOW?         -> Implicit satisfaction, checked at compile time.
REAL USE?    -> Swappable DB, notification, payment implementations.
GOTCHA?      -> No "implements" keyword; satisfaction is automatic.
INTERVIEW?   -> Structural typing vs nominal typing.
```

## 22. Practice

> Create an interface `PaymentProcessor` with method `Pay(amount float64) error`. Implement it with `UPI` and `CardPayment` types. Write a function `Checkout(p PaymentProcessor, amount float64)` that uses it.

---

# 7.2 Interface Types

## 1. What is it?

An **interface type** is a type whose values can hold **any concrete type** that satisfies its method set. It's declared with the `interface` keyword.

## 2. Why do we need it?

We need a way to say "this variable can hold different kinds of values, as long as they behave a certain way" — this is what an interface type gives us, similar to how `int` gives us a way to hold numbers.

## 3. What problem does it solve?

Without interface types, a variable/parameter can only ever hold **one specific concrete type**. Interface types let a single variable hold *different* concrete types over its lifetime, as long as they all satisfy the interface.

## 4. How does it work?

```text
interface { M1(); M2() }   <- an interface TYPE (a set of methods)

var x SomeInterface         <- x can hold ANY value whose type has M1 and M2
```

An interface type can also be built from other interfaces (**embedding**):

```go
type ReadWriter interface {
    Reader // embeds io.Reader's methods
    Writer // embeds io.Writer's methods
}
```

## 5. Simple Mental Model

```text
Interface type = an empty box shape

Any concrete value that "fits the shape" (has the right methods)
can be placed inside that box.
```

## 6. Basic Go Example

```go
package main

import "fmt"

type Animal interface {
    Sound() string
}

type Cat struct{}
func (Cat) Sound() string { return "Meow" }

type Cow struct{}
func (Cow) Sound() string { return "Moo" }

func main() {
    var a Animal   // interface type variable
    a = Cat{}
    fmt.Println(a.Sound())
    a = Cow{}      // same variable, different concrete type
    fmt.Println(a.Sound())
}
```

## 7. Explain the Code

```text
1. Animal is an interface TYPE.
2. Variable "a" is declared with type Animal.
3. First we store a Cat inside "a".
4. Later we store a Cow inside the SAME variable "a".
5. This works because both Cat and Cow satisfy Animal.
```

## 8. Real-Life Problem

```text
Logging system needs to support multiple log destinations:

type Logger interface {
    Log(msg string)
}

+-- ConsoleLogger
+-- FileLogger
+-- CloudLogger (Datadog, CloudWatch, etc.)

The rest of the app only depends on the Logger interface type.
```

## 9. When should I use it?

Use an interface type whenever a variable, parameter, or struct field needs to hold "some type that behaves this way," rather than one fixed type.

## 10. When should I NOT use it?

Don't declare a variable as an interface type if you always know the exact concrete type and never need to swap it — just use the concrete type directly for simplicity and better performance.

## 11. Common Mistakes

- Thinking an interface type is like a class — it is NOT. It has no fields, no state, only method signatures.
- Forgetting that an interface type itself is never instantiated directly (`Animal{}` is invalid — you must use a concrete type that satisfies it).

## 12. Important Gotchas

- An interface type can embed other interface types — this just merges their method sets, it's not inheritance in the OOP sense.
- Two interface types are considered identical if they have the same method set, regardless of name.

## 13. Internals

### Go Language Guarantee
An interface type is fully described by its method set (name + signature of each method).

### Implementation Detail
At runtime, an interface-typed variable is a 2-word structure (type pointer + data pointer) — covered fully in 7.5.

## 14. Standard Library Connection

```go
type ReadWriter interface {
    Reader
    Writer
}
```
`io.ReadWriter` is a real stdlib example of interface embedding.

## 15. Production Example

```go
type Cache interface {
    Get(key string) (string, bool)
    Set(key, value string)
}
```
Used by services to plug in Redis, in-memory map, or Memcached — same interface type, different implementation.

## 16. Performance

No extra cost just from *declaring* an interface type. Cost appears only when you *store a value in it* (see 7.5).

## 17. Related Concepts

| Concept | Meaning |
|---|---|
| Interface type | The shape/contract itself |
| Interface value | An actual instance holding a concrete type + value |
| Embedding | Combining multiple interfaces into one bigger interface |

## 18. Interview Questions

**Basic** — Q: Can an interface type have fields? A: No, only method signatures.

**Intermediate** — Q: What does embedding one interface inside another do? A: It merges the method sets into a single larger interface.

**Advanced** — Q: Are two interface types with the same methods but different names treated as equal by Go's type system? A: Yes, Go compares interfaces structurally — identical method sets means compatible, regardless of interface name.

**Tricky** — Q: Can you create a value of an interface type directly, like `Animal{}`? A: No — interfaces have no data of their own; you must assign a concrete type that satisfies it.

## 19. Interview Follow-Up Questions

```text
Q: What is an interface type?
Q: What does embedding do?
Q: Can two differently-named interfaces be structurally equal?
Q: Can an interface type hold state?
```

## 20. How to Explain This in an Interview

> "An interface type in Go is just a named set of method signatures. Any variable declared with that type can hold any concrete value that implements those methods. Interface types can also embed other interfaces to build bigger contracts, like io.ReadWriter combining Read and Write."

## 21. Quick Revision

```text
WHAT?      -> A named set of method signatures.
WHY?       -> Lets one variable hold different concrete types.
HOW?       -> Declared with `interface { ... }`, can embed others.
GOTCHA?    -> No fields, no state, no direct instantiation.
INTERVIEW? -> Structural equality, embedding merges method sets.
```

## 22. Practice

> Define `Shape` interface with `Area() float64`. Define `Circle` and `Rectangle`. Store both, one at a time, in a variable of type `Shape`.

---

# 7.3 Interface Satisfaction

## 1. What is it?

**Interface satisfaction** is the rule Go uses to decide: "Does this concrete type qualify to be used as this interface?"

```text
Rule: A type satisfies an interface if it has ALL the methods
that the interface requires — with matching names and signatures.
```

## 2. Why do we need it?

We need a clear rule so the compiler can check, automatically and safely, whether a value can be used where an interface is expected — without us writing that check manually.

## 3. What problem does it solve?

It removes the need for manual "does this object support this behavior?" checks. The compiler does it for you, at compile time, so bugs are caught early instead of causing runtime crashes.

## 4. How does it work?

```text
1. Look at interface's required methods.
2. Look at concrete type's method set (see 7.12/method sets note below).
3. If type has every required method with the exact matching signature -> satisfied.
4. If even one method is missing or has wrong signature -> compile error.
```

**Important detail: pointer vs value receivers.**

```text
If method has VALUE receiver  (func (t T) M())   -> both T and *T satisfy the interface.
If method has POINTER receiver (func (t *T) M()) -> only *T satisfies the interface (T does NOT).
```

This is because when a method has a pointer receiver, the compiler cannot always guarantee it can take the address of a value (e.g. a value stored in a map isn't addressable).

## 5. Simple Mental Model

```text
Satisfaction = "Do you have everything on my checklist?"

Value receiver methods  -> everyone with a copy or address can do it.
Pointer receiver methods -> only someone holding the address can do it.
```

## 6. Basic Go Example

```go
package main

import "fmt"

type Greeter interface {
    Greet() string
}

type Person struct{ Name string }

func (p *Person) Greet() string { // pointer receiver
    return "Hi, I'm " + p.Name
}

func main() {
    p := Person{Name: "Aman"}
    // var g Greeter = p     // COMPILE ERROR: Person does not satisfy Greeter
    var g Greeter = &p       // OK: *Person satisfies Greeter
    fmt.Println(g.Greet())
}
```

## 7. Explain the Code

```text
1. Greet() is defined with a POINTER receiver (*Person).
2. So only *Person satisfies Greeter, not plain Person.
3. Assigning "p" (a Person value) to Greeter fails to compile.
4. Assigning "&p" (a *Person) works fine.
```

## 8. Real-Life Problem

```text
Many Go structs use pointer receivers when methods modify internal state
(e.g. a Logger that increments a counter, or a Cache that mutates a map).

If you forget this rule, you'll get a confusing compile error like:
"Person does not implement Greeter (Greet method has pointer receiver)"
```
Understanding this rule saves hours of confusion for beginners.

## 9. When should I use it?

Use pointer receivers when the method needs to modify the struct, or when the struct is large (avoid copying). Use value receivers for small, immutable, read-only data.

## 10. When should I NOT use it?

Don't mix receiver types inconsistently on the same type without a reason — it creates confusing satisfaction rules. Pick one style (usually pointer receiver) for a given type and stick with it.

## 11. Common Mistakes

- Trying to assign a value type to an interface when the method uses a pointer receiver.
- Assuming "it compiled with one method" means "it satisfies the whole interface" — it needs ALL methods.

## 12. Important Gotchas

- A `map` value's elements are **not addressable**, so you can't call pointer-receiver methods directly on `m[key]` if `m[key]` is a struct.
- If a type embeds another type, it can "inherit" that type's methods into its own method set (see embedding, related in composition topics) — this can silently make satisfaction work.

## 13. Internals

### Go Language Guarantee
- Value receiver method → included in method set of both `T` and `*T`.
- Pointer receiver method → included in method set of `*T` only.

### Implementation Detail
The compiler generates the method table for each type at compile time; the exact addressability rules are part of the spec, not just an implementation quirk (so actually this is a language guarantee, not merely an internal detail — worth remembering precisely for interviews).

## 14. Standard Library Connection

```go
func (b *bytes.Buffer) Write(p []byte) (n int, err error)
```
`*bytes.Buffer` satisfies `io.Writer`, but a plain `bytes.Buffer` value does not (because `Write` uses a pointer receiver).

## 15. Production Example

```go
type Store interface {
    Save(v string)
}

type MemStore struct{ data []string }

func (m *MemStore) Save(v string) {
    m.data = append(m.data, v)
}

func NewStore() Store {
    return &MemStore{} // must return pointer to satisfy Store
}
```

## 16. Performance

Passing a pointer avoids copying large structs — generally good for performance. But pointers can cause heap allocation (escape analysis) — usually a non-issue unless profiling shows otherwise.

## 17. Related Concepts

| Concept | Meaning |
|---|---|
| Method set | List of methods usable on a given type (value or pointer) |
| Value receiver | Method works on a copy |
| Pointer receiver | Method works on the original, via address |

## 18. Interview Questions

**Basic** — Q: What determines if a type satisfies an interface? A: Having all required methods with matching signatures.

**Intermediate** — Q: Does `T` satisfy an interface if a method uses a pointer receiver `*T`? A: No, only `*T` satisfies it.

**Advanced** — Q: Why can't Go automatically take the address of every value to make it satisfy pointer-receiver interfaces? A: Because not every value is addressable (e.g., map values, function return values), so Go can't guarantee it always, hence the rule is strict.

**Tricky** — Q: If `T` has a pointer-receiver method, can you call that method on a value `t T` directly (not through an interface)? A: Yes, if `t` is an addressable variable, Go automatically takes `&t` for you — but this convenience does NOT extend to interface satisfaction.

## 19. Interview Follow-Up Questions

```text
Q: What is interface satisfaction?
Q: What's the difference in method sets between T and *T?
Q: Why does Person{} fail but &Person{} work for a pointer-receiver interface?
Q: What is "addressability" in Go?
```

## 20. How to Explain This in an Interview

> "A type satisfies an interface when it has every method the interface requires, with matching signatures. There's a subtlety with receivers: if a method has a pointer receiver, only the pointer type satisfies the interface, not the value type — because Go can't always guarantee it can take the address of a value. This is a very common gotcha for people switching from other languages."

## 21. Quick Revision

```text
WHAT?      -> The rule for "does this type qualify for this interface".
WHY?       -> Compile-time safety without manual checks.
HOW?       -> Match ALL methods; pointer receivers restrict to *T.
GOTCHA?    -> Value type doesn't satisfy interface if method has pointer receiver.
INTERVIEW? -> Explain addressability + method sets precisely.
```

## 22. Practice

> Create a `Counter` struct with a pointer-receiver method `Increment()`. Try assigning a `Counter` value (not pointer) to an interface requiring `Increment()`, and observe the compile error.

---

# 7.4 Parsing Flags with flag.Value

## 1. What is it?

`flag.Value` is a small interface defined by Go's standard `flag` package. It lets you create **custom command-line flag types** — not just the built-in `string`, `int`, `bool`.

```go
type Value interface {
    String() string
    Set(string) error
}
```

## 2. Why do we need it?

CLI tools often need flags that aren't simple strings/numbers — like a duration (`1h30m`), a comma-separated list, or a custom enum (e.g. `--log-level=debug`). `flag.Value` lets you plug in **your own parsing logic** for such flags.

## 3. What problem does it solve?

```text
Without flag.Value:
You'd parse raw strings manually after flag.Parse(),
scattered validation logic everywhere in main().

With flag.Value:
Parsing + validation is bundled into one type,
and flag.Parse() calls it automatically.
```

## 4. How does it work?

```text
1. You define a custom type (e.g. type celsiusFlag float64).
2. You implement String() string  -> how to display current value.
3. You implement Set(string) error -> how to parse a new value from CLI text.
4. You register it with flag.Var(&myFlag, "name", "usage").
5. flag.Parse() calls Set() automatically when it sees --name=value.
```

## 5. Simple Mental Model

```text
flag.Value = "teach the flag package how to read MY custom type"
```

## 6. Basic Go Example

```go
package main

import (
    "flag"
    "fmt"
)

type celsiusFlag float64

func (c *celsiusFlag) String() string {
    return fmt.Sprintf("%g°C", float64(*c))
}

func (c *celsiusFlag) Set(s string) error {
    var f float64
    _, err := fmt.Sscanf(s, "%g", &f)
    if err != nil {
        return err
    }
    *c = celsiusFlag(f)
    return nil
}

func main() {
    var temp celsiusFlag
    flag.Var(&temp, "temp", "temperature in Celsius")
    flag.Parse()
    fmt.Println("Temperature:", temp)
}
```

Run: `go run main.go -temp=30`

## 7. Explain the Code

```text
1. celsiusFlag is a custom float64-based type.
2. String() defines how it's printed (used in help text and defaults).
3. Set() defines how to parse the CLI string into our type.
4. flag.Var registers our custom type as a flag named "temp".
5. When the CLI has -temp=30, flag.Parse() calls Set("30") automatically.
```

## 8. Real-Life Problem

```text
A backend CLI tool needs a --log-level flag that only accepts:
"debug", "info", "warn", "error" (and rejects anything else).

type logLevel string

func (l *logLevel) Set(s string) error {
    switch s {
    case "debug", "info", "warn", "error":
        *l = logLevel(s)
        return nil
    default:
        return fmt.Errorf("invalid log level: %s", s)
    }
}
```
This gives you built-in validation, right at flag parsing time — before your program even starts doing real work.

## 9. When should I use it?

When your CLI tool needs flags with custom types, custom validation, or custom parsing (durations, enums, IP addresses, lists).

## 10. When should I NOT use it?

Don't use it for simple built-in types (`string`, `int`, `bool`, `float64`, `time.Duration`) — the `flag` package already provides `flag.String`, `flag.Int`, etc. Only build a custom `flag.Value` when the built-ins don't fit.

## 11. Common Mistakes

- Forgetting the receiver must be a **pointer** (`*celsiusFlag`), otherwise `Set()` can't modify the value.
- Forgetting to register with `flag.Var` (people sometimes just implement the interface but never call `flag.Var`).

## 12. Important Gotchas

- `flag.Value`'s `Set` method must have a **pointer receiver**, because it mutates the underlying value — this connects directly back to the pointer-receiver satisfaction rule from 7.3.
- If `Set()` returns an error, `flag.Parse()` will print usage and exit the program (calls `os.Exit(2)`).

## 13. Internals

### Go Language Guarantee
`flag.Value` is just a regular interface — no magic beyond normal interface satisfaction rules.

### Implementation Detail
Internally, the `flag` package stores your value along with its `String()`/`Set()` methods in a `flag.Flag` struct, and calls `Set()` whenever it encounters that flag on the command line.

## 14. Standard Library Connection

This is one of the best real examples of a **tiny, purpose-built interface** in the standard library — showing Go's philosophy of small interfaces over big frameworks.

## 15. Production Example

```go
type stringList []string

func (s *stringList) String() string { return strings.Join(*s, ",") }
func (s *stringList) Set(v string) error {
    *s = append(*s, v)
    return nil
}
```
Used for flags like `-tag=foo -tag=bar` to build up a list.

## 16. Performance

Negligible — this runs once at program startup during flag parsing, not in a hot path.

## 17. Related Concepts

| Concept | Meaning |
|---|---|
| flag.Value | Interface for custom CLI flag types |
| fmt.Stringer | Similar `String()` method, used for printing |
| encoding.TextUnmarshaler | Similar idea, used for text-based decoding elsewhere |

## 18. Interview Questions

**Basic** — Q: What two methods does `flag.Value` require? A: `String() string` and `Set(string) error`.

**Intermediate** — Q: Why must the receiver be a pointer? A: Because `Set()` needs to modify the underlying value.

**Advanced** — Q: How does this example demonstrate Go's interface philosophy? A: It shows how a tiny 2-method interface can plug custom behavior into a generic system (the flag package) without any inheritance or big framework.

**Tricky** — Q: What happens if `Set()` returns an error during flag parsing? A: `flag.Parse()` prints the error and usage message, then exits the program.

## 19. Interview Follow-Up Questions

```text
Q: What is flag.Value?
Q: Why pointer receiver here specifically?
Q: How would you build a custom "list of strings" flag?
Q: What's the difference between flag.Value and fmt.Stringer?
```

## 20. How to Explain This in an Interview

> "flag.Value is a small standard-library interface with just two methods — String and Set — that lets you plug custom types into Go's command-line flag parser. Instead of manually parsing raw strings after flag.Parse, you implement Set to convert and validate the string, and flag.Parse calls it automatically. It's a great real example of how tiny interfaces let you extend standard library behavior without modifying it."

## 21. Quick Revision

```text
WHAT?      -> Interface for custom CLI flag types: String() + Set().
WHY?       -> Built-in types don't cover custom parsing/validation needs.
HOW?       -> Implement both methods (pointer receiver), register with flag.Var.
GOTCHA?    -> Must use pointer receiver; Set() error exits the program.
INTERVIEW? -> Great example of small, purpose-built interfaces.
```

## 22. Practice

> Build a custom `flag.Value` type `durationRange` that parses `"1h-3h"` into two `time.Duration` values.

---

# 7.5 Interface Values

## 1. What is it?

An **interface value** is what actually gets stored when you assign a concrete value to an interface-typed variable. Internally, it's a pair:

```text
Interface Value = (Dynamic Type, Dynamic Value)
```

## 2. Why do we need it?

The interface needs to "remember" what concrete type is currently inside it, so Go knows which method implementation to call when you invoke a method on the interface.

## 3. What problem does it solve?

Without this dynamic pairing, Go wouldn't know **which** concrete implementation's method to run when you call a method through an interface variable. The (type, value) pair is what makes **dynamic dispatch** possible.

## 4. How does it work?

```text
+-------------------+
|  Dynamic Type      |   e.g. *os.File, Dog, MyError
+-------------------+
|  Dynamic Value      |   the actual data
+-------------------+

var w io.Writer = os.Stdout
   Dynamic Type  = *os.File
   Dynamic Value = pointer to Stdout's file descriptor data
```

When you call `w.Write(...)`, Go looks at the dynamic type, finds `*os.File`'s `Write` method, and calls it with the dynamic value.

## 5. Simple Mental Model

```text
Interface value = a labeled box

Label  = the real type inside
Content = the real data inside

Calling a method opens the box, reads the label,
and runs that type's version of the method.
```

## 6. Basic Go Example

```go
package main

import "fmt"

type Shape interface {
    Area() float64
}

type Circle struct{ R float64 }
func (c Circle) Area() float64 { return 3.14 * c.R * c.R }

func main() {
    var s Shape
    fmt.Println(s == nil) // true: both type and value are nil

    s = Circle{R: 2}
    fmt.Printf("%T %v\n", s, s) // dynamic type + value
    fmt.Println(s.Area())
}
```

## 7. Explain the Code

```text
1. "var s Shape" creates an interface value with (nil type, nil value).
2. s == nil is true ONLY when BOTH parts are nil.
3. Assigning Circle{R:2} sets:
      dynamic type  = Circle
      dynamic value = {R:2}
4. %T prints the dynamic type; %v prints the dynamic value.
5. s.Area() runs Circle's Area() method using the stored value.
```

## 8. Real-Life Problem

```text
When debugging production issues, engineers often print %T on an
interface variable (like an "err" of type error) to find out
EXACTLY which concrete error type caused a failure —
this only works because of the dynamic type stored inside the interface value.
```

## 9. When should I use it?

You don't "choose" to use interface values directly — they happen automatically any time you assign a concrete value to an interface variable. Understanding them helps you reason about nil checks, type assertions, and debugging.

## 10. When should I NOT use it?

N/A — this is a core mechanism, not an optional feature. But you should avoid **relying on interface value internals** (like assuming layout/size) since those are implementation details, not guaranteed by the spec.

## 11. Common Mistakes

- Comparing an interface value to `nil` and expecting `true`, when it actually holds a non-nil type with a nil value (the classic "nil interface" gotcha — see below).
- Assuming `%v` always shows enough to debug — sometimes you need `%T` too.

## 12. Important Gotchas — THE FAMOUS ONE

```go
type MyError struct{}
func (*MyError) Error() string { return "boom" }

func doSomething() error {
    var err *MyError = nil
    return err // returns a NON-nil interface!
}

func main() {
    err := doSomething()
    fmt.Println(err == nil) // false! surprising to beginners
}
```

```text
WHY?
err's dynamic type = *MyError  (NOT nil)
err's dynamic value = nil       (nil pointer)

Interface value == nil requires BOTH type AND value to be nil.
Here only the value is nil, so the interface itself is NOT nil.
```

**This is one of the most common Go interview questions.** Remember it well.

## 13. Internals

### Go Language Guarantee
An interface value is `nil` if and only if both its dynamic type and dynamic value are `nil`.

### Implementation Detail
The two-word (type pointer, data pointer) internal representation is a common implementation approach in most Go compilers, but the spec itself only guarantees the *behavior* (the nil rule above), not the exact memory layout.

## 14. Standard Library Connection

```text
This exact gotcha shows up constantly with the `error` interface —
functions returning a concrete *SomeError type as `error`
can accidentally return a "non-nil" error even when logically there's no error.
```

## 15. Production Example

```go
// BAD pattern that causes the bug:
func validate() error {
    var verr *ValidationError
    if somethingWrong {
        verr = &ValidationError{}
    }
    return verr // BUG: always non-nil interface, even if verr is nil!
}

// GOOD pattern:
func validate() error {
    if somethingWrong {
        return &ValidationError{}
    }
    return nil // explicit nil interface
}
```

## 16. Performance

Assigning a concrete value to an interface **may cause a heap allocation** if the value doesn't fit in a machine word or if it escapes (this is why very hot-path code sometimes avoids unnecessary interface conversions). Measure before optimizing — don't guess.

## 17. Related Concepts

| Concept | Meaning |
|---|---|
| Interface value | (dynamic type, dynamic value) pair |
| Nil interface | Both parts nil |
| Typed nil in interface | Type is non-nil, value is nil pointer — interface is NOT nil |

## 18. Interview Questions

**Basic** — Q: What does an interface value contain? A: A dynamic type and a dynamic value.

**Intermediate** — Q: When is an interface value equal to nil? A: Only when both the dynamic type and dynamic value are nil.

**Advanced** — Q: Why can a function returning `error` with a nil `*MyError` still produce a non-nil error? A: Because the returned interface has a non-nil dynamic type (`*MyError`), even though the underlying pointer value is nil.

**Tricky** — Q: How do you avoid the typed-nil bug? A: Return a literal `nil` explicitly when there's no error, instead of returning a typed nil pointer variable.

## 19. Interview Follow-Up Questions

```text
Q: What is an interface value made of?
Q: When is it truly nil?
Q: What's the typed-nil-pointer-in-interface bug?
Q: How do you avoid it in production code?
```

## 20. How to Explain This in an Interview

> "An interface value in Go is really a pair — a dynamic type and a dynamic value. It's only equal to nil when both parts are nil. This causes a classic gotcha: if you return a nil pointer of a concrete error type as an `error` interface, the interface itself is NOT nil, because it still carries a non-nil type. The fix is to explicitly return a literal nil when there's no error, instead of returning a possibly-nil typed variable."

## 21. Quick Revision

```text
WHAT?      -> (Dynamic Type, Dynamic Value) pair.
WHY?       -> Needed for dynamic dispatch (right method call at runtime).
GOTCHA?    -> Typed nil pointer inside interface != nil interface.
INTERVIEW? -> #1 most common interfaces gotcha question — know it cold.
```

## 22. Practice

> Reproduce the typed-nil bug yourself: write a function returning `error` from a nil `*MyError`, print `err == nil`, then fix it.

---

# 7.6 Sorting with sort.Interface

## 1. What is it?

`sort.Interface` is a 3-method interface from Go's `sort` package that lets you sort **any** collection, as long as you tell Go how to compare and swap its elements.

```go
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}
```

## 2. Why do we need it?

Go doesn't have generics-based universal sorting in older code (and even with generics, this pattern is still foundational). `sort.Interface` lets the `sort` package sort **any custom data structure** without knowing anything about its internal fields.

## 3. What problem does it solve?

```text
Without sort.Interface:
You'd write a custom sort algorithm for every single data type. Repetitive and error-prone.

With sort.Interface:
You write Len/Less/Swap once for your type,
and reuse Go's battle-tested sort algorithm (introsort) for free.
```

## 4. How does it work?

```text
1. Define a named slice type, e.g. type ByAge []Person.
2. Implement Len(), Less(i,j), Swap(i,j) on ByAge.
3. Call sort.Sort(ByAge(people)).
4. sort.Sort repeatedly calls Less() and Swap() to reorder the data.
```

## 5. Simple Mental Model

```text
sort.Interface = "teach me how to compare 2 items and swap them,
                   and I'll handle the sorting algorithm for you."
```

## 6. Basic Go Example

```go
package main

import (
    "fmt"
    "sort"
)

type Person struct {
    Name string
    Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func main() {
    people := []Person{{"Amit", 30}, {"Sara", 22}, {"Raj", 45}}
    sort.Sort(ByAge(people))
    fmt.Println(people)
}
```

## 7. Explain the Code

```text
1. ByAge is just []Person with extra methods attached.
2. Len() tells sort how many elements exist.
3. Less() tells sort the ordering rule (younger first here).
4. Swap() tells sort how to exchange two elements.
5. sort.Sort uses these 3 methods to fully sort the slice in place.
```

## 8. Real-Life Problem

```text
An admin dashboard needs to sort users by: name, signup date, or activity score
— chosen dynamically by the logged-in admin.

Instead of writing 3 separate sort functions,
you write 3 small types (ByName, ByDate, ByScore),
each implementing sort.Interface, and swap which one you use.
```

## 9. When should I use it?

When sorting custom structs/slices by custom rules, especially before generics, or when you need the flexibility of swapping the comparison logic at runtime.

## 10. When should I NOT use it?

For simple, one-off sorts of built-in types (`[]int`, `[]string`), just use `sort.Ints`, `sort.Strings`, or (in modern Go) `sort.Slice` / `slices.Sort` with generics — no need to define a whole named type with 3 methods.

## 11. Common Mistakes

- Forgetting `Swap` must actually swap the underlying data, not just a copy.
- Confusing `sort.Sort` (needs `sort.Interface`) with `sort.Slice` (takes a plain slice + `less` func — much less boilerplate for one-off sorts).

## 12. Important Gotchas

- `sort.Sort` is **not guaranteed to be stable** (equal elements might get reordered). Use `sort.Stable` if stability matters.
- `sort.Reverse` wraps your `sort.Interface` to flip the order, by cleverly swapping `Less(i,j)` to call `Less(j,i)`.

## 13. Internals

### Go Language Guarantee
`sort.Sort` will produce a fully sorted result based on your `Less` function, but doesn't guarantee stability.

### Implementation Detail
Go's `sort` package internally uses an introsort-like hybrid algorithm (quicksort + heapsort + insertion sort) for performance — this can change between Go versions since it's not part of the language spec, just the implementation.

## 14. Standard Library Connection

```text
sort.Sort, sort.Stable, sort.Reverse, sort.Search
-- all built around sort.Interface.
```

This is a textbook example of designing a whole package around one small interface.

## 15. Production Example

```go
type ByScore []User

func (b ByScore) Len() int           { return len(b) }
func (b ByScore) Less(i, j int) bool { return b[i].Score > b[j].Score } // descending
func (b ByScore) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

sort.Sort(ByScore(leaderboard))
```

## 16. Performance

Sorting is O(n log n) on average. The main performance factor is how expensive your `Less`/`Swap` functions are — keep them cheap (avoid heavy computation inside them).

## 17. Related Concepts

| Concept | Meaning |
|---|---|
| sort.Interface | 3-method contract: Len, Less, Swap |
| sort.Slice | Modern shortcut — no need to define a named type |
| sort.Stable | Guarantees equal elements keep original order |

## 18. Interview Questions

**Basic** — Q: What 3 methods does `sort.Interface` require? A: `Len()`, `Less(i, j)`, `Swap(i, j)`.

**Intermediate** — Q: What's the difference between `sort.Sort` and `sort.Slice`? A: `sort.Sort` needs a full `sort.Interface` implementation (named type + 3 methods); `sort.Slice` just needs a slice and a `less` function, using reflection internally — quicker for one-off sorts.

**Advanced** — Q: How does `sort.Reverse` work internally? A: It wraps your `sort.Interface` in a struct that swaps the arguments of `Less(i,j)` to `Less(j,i)`, flipping the sort order without duplicating logic.

**Tricky** — Q: Is `sort.Sort` stable? A: No — use `sort.Stable` if you need equal elements to preserve their original relative order.

## 19. Interview Follow-Up Questions

```text
Q: What is sort.Interface?
Q: How would you sort in descending order using it?
Q: Difference between sort.Sort and sort.Slice?
Q: What does sort.Stable guarantee that sort.Sort doesn't?
```

## 20. How to Explain This in an Interview

> "sort.Interface is a 3-method contract — Len, Less, and Swap — that lets Go's sort package sort any custom type without knowing its internal structure. You implement those 3 methods on a named slice type, then call sort.Sort. It's a great example of designing small interfaces that let a package be reused across completely unrelated data types."

## 21. Quick Revision

```text
WHAT?      -> Len/Less/Swap contract for custom sorting.
WHY?       -> Reuse Go's sort algorithm for any data type.
HOW?       -> Define named slice type + implement 3 methods.
GOTCHA?    -> sort.Sort isn't stable by default; use sort.Stable.
INTERVIEW? -> Know sort.Slice as the modern shortcut alternative.
```

## 22. Practice

> Given a slice of `Product{Name string, Price float64}`, implement `ByPrice` satisfying `sort.Interface`, then sort ascending and descending using `sort.Reverse`.

---

# 7.7 The http.Handler Interface

## 1. What is it?

`http.Handler` is the core interface behind Go's entire `net/http` package — anything that can "handle" an HTTP request satisfies it.

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

## 2. Why do we need it?

A web server needs a uniform way to say "when a request comes in, run this logic." `http.Handler` provides that single, uniform contract for every possible request handler — middleware, routers, static file servers, everything.

## 3. What problem does it solve?

```text
Without http.Handler:
Every web framework would need its own custom way to register request logic.

With http.Handler:
Anything satisfying ServeHTTP() can plug into net/http's server,
router, or middleware — total interoperability.
```

## 4. How does it work?

```text
1. You implement ServeHTTP(w http.ResponseWriter, r *http.Request).
2. You register your handler with a mux (router) or http.ListenAndServe.
3. When a request arrives, the server calls YourHandler.ServeHTTP(w, r).
4. You write the response using w, and read request data using r.
```

## 5. Simple Mental Model

```text
http.Handler = "Give me a request, I'll produce a response."
```

## 6. Basic Go Example

```go
package main

import (
    "fmt"
    "net/http"
)

type helloHandler struct{}

func (helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, World!")
}

func main() {
    var h http.Handler = helloHandler{}
    http.ListenAndServe(":8080", h)
}
```

## 7. Explain the Code

```text
1. helloHandler is a struct with no fields.
2. It has a ServeHTTP method -> satisfies http.Handler.
3. ListenAndServe uses "h" to handle EVERY incoming request.
4. w is used to write the response; r holds request info.
```

## 8. Real-Life Problem

```text
Middleware pattern (logging, auth, rate limiting) is built
entirely on top of http.Handler:

func withLogging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
```
This is how nearly every Go web framework (chi, gin-adjacent patterns, standard net/http) builds middleware chains — because everything speaks the same `http.Handler` language.

## 9. When should I use it?

Any time you're building HTTP servers, custom middleware, or need to plug custom logic into Go's HTTP stack.

## 10. When should I NOT use it?

For a single simple route with no reusable state, `http.HandlerFunc` (a function type that already satisfies `http.Handler`) is simpler than defining a whole struct.

## 11. Common Mistakes

- Forgetting `http.HandlerFunc` exists — it lets you use a plain function as a handler without writing a struct.
- Writing to `w` after already sending headers, causing subtle bugs.

## 12. Important Gotchas

- `http.HandlerFunc` is a **function type** that implements `ServeHTTP` by simply calling itself — a beautiful example of a function satisfying an interface:
```go
type HandlerFunc func(ResponseWriter, *Request)
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) { f(w, r) }
```
- Middleware chains are just **nested http.Handlers wrapping http.Handlers**.

## 13. Internals

### Go Language Guarantee
Any type with a matching `ServeHTTP(ResponseWriter, *Request)` method satisfies `http.Handler` — including function types, thanks to `HandlerFunc`.

### Implementation Detail
`net/http`'s server internally runs each request in its own goroutine, calling your handler's `ServeHTTP` — this concurrency model is implementation behavior, not something the `Handler` interface itself dictates.

## 14. Standard Library Connection

This entire section IS standard library — `http.Handler`, `http.HandlerFunc`, `http.ServeMux` are all core `net/http` types built around one interface.

## 15. Production Example

```go
mux := http.NewServeMux()
mux.Handle("/users", userHandler{})
mux.Handle("/health", http.HandlerFunc(healthCheck))

var handler http.Handler = withLogging(withAuth(mux))
http.ListenAndServe(":8080", handler)
```

## 16. Performance

`ServeHTTP` runs per-request in its own goroutine — cheap to create, but avoid heavy blocking work inside it without care (use context timeouts, avoid unbounded goroutine spawning).

## 17. Related Concepts

| Concept | Meaning |
|---|---|
| http.Handler | Interface: ServeHTTP(w, r) |
| http.HandlerFunc | Function type that also satisfies Handler |
| Middleware | A Handler that wraps another Handler |

## 18. Interview Questions

**Basic** — Q: What method must a type implement to be an `http.Handler`? A: `ServeHTTP(http.ResponseWriter, *http.Request)`.

**Intermediate** — Q: What is `http.HandlerFunc`? A: A function type that implements `http.Handler` by calling itself inside its own `ServeHTTP` method — letting plain functions act as handlers.

**Advanced** — Q: How does Go's middleware pattern rely on interfaces? A: Middleware functions take a `Handler` and return a new `Handler` that wraps it, forming a chain — all made possible because everything shares the same `Handler` interface.

**Tricky** — Q: Can a function satisfy an interface directly in Go? A: Not directly, but via a named function type with a method (like `HandlerFunc`), a function value can satisfy an interface — a classic and elegant Go pattern.

## 19. Interview Follow-Up Questions

```text
Q: What is http.Handler?
Q: What's http.HandlerFunc and why does it exist?
Q: How is middleware implemented using this interface?
Q: Can a bare func satisfy an interface without a wrapper type?
```

## 20. How to Explain This in an Interview

> "http.Handler is the single interface that Go's entire net/http package is built around — anything with a ServeHTTP(w, r) method can handle requests. What's elegant is http.HandlerFunc, a function type that also implements Handler by calling itself — so you can use a plain function as a handler. Middleware is just handlers wrapping handlers, all speaking the same interface, which is why Go's HTTP ecosystem composes so cleanly."

## 21. Quick Revision

```text
WHAT?      -> ServeHTTP(w, r) contract for handling HTTP requests.
WHY?       -> Uniform way to plug logic into the HTTP server/router.
HOW?       -> Implement ServeHTTP, or use HandlerFunc for plain functions.
GOTCHA?    -> HandlerFunc lets a function satisfy an interface — remember this trick.
INTERVIEW? -> Explain middleware as nested Handlers.
```

## 22. Practice

> Write a middleware `withRequestID` that adds a random request ID to the context before calling the next handler.

---

# 7.8 The error Interface

## 1. What is it?

`error` is Go's built-in interface for representing failure.

```go
type error interface {
    Error() string
}
```

## 2. Why do we need it?

Go has no exceptions (no try/catch). Errors are just **regular values** — the `error` interface is what lets any type represent "something went wrong," in a uniform way that functions can return and callers must explicitly check.

## 3. What problem does it solve?

```text
Without a shared error interface:
Every package would invent its own way to represent failures — chaos.

With error:
Every function that can fail returns (result, error),
and callers use one consistent pattern: `if err != nil { ... }`
```

## 4. How does it work?

```text
1. Any type with an Error() string method satisfies error.
2. Functions return error as their last return value by convention.
3. Caller checks: if err != nil { handle it }.
4. errors.New() and fmt.Errorf() create simple error values.
```

## 5. Simple Mental Model

```text
error = "just another interface, with one method: Error() string"

No magic. No exceptions. Just values you must check.
```

## 6. Basic Go Example

```go
package main

import (
    "errors"
    "fmt"
)

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)
}
```

## 7. Explain the Code

```text
1. divide returns (float64, error).
2. If b is 0, we return a zero value + a real error.
3. errors.New creates a simple error value satisfying the error interface.
4. The caller MUST check err != nil before trusting "result".
```

## 8. Real-Life Problem

```text
A backend API call to a database can fail for many reasons:
connection lost, record not found, permission denied.

Using the error interface + wrapping (fmt.Errorf with %w),
we build a clear chain of "what failed, and why", which
can be logged and traced through multiple layers of the app.
```

## 9. When should I use it?

For any function that can fail — file I/O, network calls, DB queries, parsing, validation. It's Go's default failure-reporting mechanism.

## 10. When should I NOT use it?

Don't use `error` for expected, non-exceptional control flow best expressed with normal return values (like "found = false" in a lookup) unless the "not found" case really represents a failure state that callers must handle explicitly.

## 11. Common Mistakes

- Ignoring errors (`_ = someFunc()`), leading to silent failures.
- Comparing errors with `==` when they've been wrapped (use `errors.Is`/`errors.As` instead — see 7.11).
- The typed-nil-in-interface bug (from 7.5) applies directly to `error` — a very common real-world bug.

## 12. Important Gotchas

- `err != nil` check must happen **immediately** after the call, before using the result.
- Custom error types let you attach extra data (error codes, HTTP status, etc.) beyond just a message.

## 13. Internals

### Go Language Guarantee
`error` is a built-in predeclared interface type: `type error interface { Error() string }`. Any type with that method satisfies it.

### Implementation Detail
`errors.New` returns a pointer to an internal struct holding just a string — a minimal implementation, not something the language spec mandates in detail.

## 14. Standard Library Connection

`error` is used **everywhere** — `os.Open`, `json.Unmarshal`, `http.Get`, virtually every stdlib function that can fail returns one.

## 15. Production Example

```go
type NotFoundError struct {
    ID int64
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("record with id %d not found", e.ID)
}

func FindUser(id int64) (*User, error) {
    u, ok := db[id]
    if !ok {
        return nil, &NotFoundError{ID: id}
    }
    return u, nil
}
```
Callers can later check `errors.As(err, &notFoundErr)` to react specifically to a "not found" case (covered in 7.11).

## 16. Performance

Creating errors is cheap. Avoid creating errors in extremely hot loops unnecessarily, but generally this is a non-issue for backend code.

## 17. Related Concepts

| Concept | Meaning |
|---|---|
| error | Built-in interface: Error() string |
| errors.New / fmt.Errorf | Ways to create error values |
| Custom error types | Structs implementing Error() with extra fields |

## 18. Interview Questions

**Basic** — Q: What method must a type have to satisfy `error`? A: `Error() string`.

**Intermediate** — Q: How does Go's error handling differ from exceptions in Java/Python? A: Errors are ordinary return values checked explicitly by the caller, not thrown/caught control flow.

**Advanced** — Q: Why does Go prefer explicit error returns over exceptions? A: It makes failure paths visible in the function signature and forces callers to consciously handle them, improving code clarity and avoiding hidden control-flow jumps.

**Tricky** — Q: Can returning a nil pointer of a custom error type cause a bug? A: Yes — see the typed-nil-in-interface issue from 7.5; a nil `*MyError` returned as `error` is a non-nil interface.

## 19. Interview Follow-Up Questions

```text
Q: What is the error interface?
Q: How do you create custom error types?
Q: Why does Go avoid exceptions?
Q: What's the typed-nil bug with error?
```

## 20. How to Explain This in an Interview

> "error is just a built-in interface with a single method, Error() string. Go doesn't have exceptions — functions that can fail return an error value as their last return, and callers explicitly check it with `if err != nil`. This makes failure handling visible and predictable. You can also define custom error types with extra fields, which is useful for attaching context like error codes."

## 21. Quick Revision

```text
WHAT?      -> Built-in interface: Error() string.
WHY?       -> Explicit, value-based failure handling instead of exceptions.
HOW?       -> Return error as last value; check err != nil.
GOTCHA?    -> Typed-nil-in-interface bug applies here too.
INTERVIEW? -> Compare vs exceptions; mention custom error types.
```

## 22. Practice

> Create a custom `ValidationError` type with an `Error()` method that includes a field name and message. Return it from a validation function.

---

# 7.9 Example: Expression Evaluator

## 1. What is it?

This is a **worked example** (from Go teaching material) showing how interfaces let you build an **expression tree evaluator** — e.g., parsing and evaluating math expressions like `x*x + y*y` — where every node type (literal, variable, binary operation, function call) satisfies the same `Expr` interface.

## 2. Why do we need it?

We need a way to represent very different kinds of expression nodes (numbers, variables, operators) uniformly, so a single `Eval()` function can work on ANY of them without a giant if/else chain checking types manually.

## 3. What problem does it solve?

```text
Without an Expr interface:
You'd need a huge switch statement checking types manually everywhere
you want to evaluate, print, or check an expression.

With Expr interface:
Each node type implements Eval() itself — the tree "evaluates itself"
recursively, node by node.
```

## 4. How does it work?

```text
type Expr interface {
    Eval(env Env) float64
}

Node types: literal, Var, unary, binary, call
Each implements Eval() differently:
  literal.Eval()  -> just returns its number
  Var.Eval()      -> looks up value in environment (a map)
  binary.Eval()   -> evaluates left + right, applies operator
  call.Eval()     -> evaluates arguments, applies a function (sin, sqrt...)
```

## 5. Simple Mental Model

```text
Expr = "I know how to compute my own value, given the variables around me."

A tree of Exprs evaluates itself recursively —
each node asks its children to Eval(), then combines the results.
```

## 6. Basic Go Example

```go
package main

import "fmt"

type Env map[string]float64

type Expr interface {
    Eval(env Env) float64
}

type literal float64
func (l literal) Eval(_ Env) float64 { return float64(l) }

type variable string
func (v variable) Eval(env Env) float64 { return env[string(v)] }

type binary struct {
    op          byte
    left, right Expr
}

func (b binary) Eval(env Env) float64 {
    switch b.op {
    case '+':
        return b.left.Eval(env) + b.right.Eval(env)
    case '*':
        return b.left.Eval(env) * b.right.Eval(env)
    }
    panic("unsupported operator")
}

func main() {
    // Represents: x*x + y*y
    expr := binary{'+',
        binary{'*', variable("x"), variable("x")},
        binary{'*', variable("y"), variable("y")},
    }
    env := Env{"x": 3, "y": 4}
    fmt.Println(expr.Eval(env)) // 25
}
```

## 7. Explain the Code

```text
1. Expr interface requires just one method: Eval(env) float64.
2. literal, variable, binary each implement Eval() differently.
3. binary holds two child Exprs (left, right) — this is RECURSION via interfaces.
4. Calling expr.Eval(env) triggers a chain of nested Eval() calls.
5. Each node doesn't need to know about the OTHER node types — total decoupling.
```

## 8. Real-Life Problem

```text
This exact pattern (interface-based expression trees) is used in:

- SQL query builders/parsers
- Configuration rule engines (e.g. feature-flag conditions)
- Template engines
- Compilers/interpreters (AST evaluation)

Business Rule Engine Example:
Rule: "user.age > 18 AND user.country == 'IN'"
Represented as a tree of Expr nodes (AndExpr, CompareExpr, FieldExpr),
each implementing Eval(user) bool.
```

## 9. When should I use it?

When you need to represent and evaluate a **recursive, tree-shaped** structure with multiple different node "kinds" — parsers, rule engines, calculators, template systems.

## 10. When should I NOT use it?

Don't build a whole expression-tree/interface system for something that's really just a simple, flat calculation — that's over-engineering. Use this pattern only when the structure is genuinely recursive/tree-shaped.

## 11. Common Mistakes

- Forgetting to handle all node types (leads to `panic` at runtime instead of compile-time safety, since the `switch` on operator/type is a runtime check).
- Making the tree mutable when it should be immutable (harder to reason about, especially with concurrency).

## 12. Important Gotchas

- This is a great example of a **type switch** in real use — when evaluating unknown expression types, code often falls back to a `switch e := expr.(type)` (covered fully in 7.13).
- Deep recursive trees can cause **stack overflow** on pathologically deep input — worth validating input depth in production parsers.

## 13. Internals

### Go Language Guarantee
Interface method dispatch always calls the correct concrete type's method, no matter how deeply nested the recursive structure is.

### Implementation Detail
Each `Eval()` call is a normal Go function call using the interface's dynamic dispatch — the Go runtime doesn't do anything special for recursive interface trees; it's just regular recursive function calls.

## 14. Standard Library Connection

```text
go/ast package -> the Go compiler itself represents Go source code
                   as a tree of interface-typed AST nodes (ast.Expr, ast.Stmt)!
```
This is literally how the Go compiler parses Go code — the same pattern you just learned.

## 15. Production Example

```go
type Rule interface {
    Evaluate(ctx Context) bool
}

type AndRule struct{ Rules []Rule }
func (a AndRule) Evaluate(ctx Context) bool {
    for _, r := range a.Rules {
        if !r.Evaluate(ctx) {
            return false
        }
    }
    return true
}
```
A feature-flag or access-control system built exactly like this.

## 16. Performance

Recursive interface calls have small overhead per call (dynamic dispatch), but this is rarely a bottleneck compared to actual business logic. For very hot evaluation paths (e.g. compiling rules once instead of re-walking trees), consider caching or compiling the tree.

## 17. Related Concepts

| Concept | Meaning |
|---|---|
| Expr interface | Contract: Eval(env) float64 |
| Recursive interface structure | Tree of interfaces, each holding child interfaces |
| Type switch | Used to add behavior (e.g. printing/checking) without changing the interface |

## 18. Interview Questions

**Basic** — Q: What single method does the `Expr` interface need? A: `Eval(env Env) float64`.

**Intermediate** — Q: How does the expression tree evaluate itself without one giant switch statement? A: Each node type implements its own `Eval()`, and evaluation happens through recursive interface calls — the "switch" is distributed across types instead of centralized.

**Advanced** — Q: What real Go tool uses this exact same pattern? A: The Go compiler's own `go/ast` package represents source code as a tree of interface-typed nodes.

**Tricky** — Q: What happens if you add a new expression type later? A: You just implement `Eval()` on the new type — no existing code needs to change (this is the Open/Closed Principle in action).

## 19. Interview Follow-Up Questions

```text
Q: What is the Expr interface here?
Q: How does recursion work through interfaces?
Q: What's a real-world use case for this pattern?
Q: How would you add a new operator/node type?
```

## 20. How to Explain This in an Interview

> "This pattern models an expression as a tree of interface values, where each node type — literal, variable, binary operation — implements the same Eval method differently. Evaluating the whole tree just means calling Eval on the root, which recursively calls Eval on its children. It avoids a giant type-checking switch statement and makes the system open for extension — you can add new node types without touching existing code. It's the same idea behind Go's own AST package."

## 21. Quick Revision

```text
WHAT?      -> Tree of interface nodes, each self-evaluating.
WHY?       -> Avoid giant switch statements; enable extension.
HOW?       -> Each node type implements Eval() recursively.
REAL USE?  -> Rule engines, parsers, Go's own go/ast package.
INTERVIEW? -> Mention Open/Closed Principle + go/ast connection.
```

## 22. Practice

> Extend the expression evaluator to support a `call` node type for functions like `sqrt(x)`.

---

# 7.10 Type Assertions

## 1. What is it?

A **type assertion** lets you extract the concrete value stored inside an interface — or check what concrete type it holds.

```go
v, ok := x.(T)
```

## 2. Why do we need it?

Sometimes you have a value stored inside a generic interface (like `error` or `interface{}`), but you need access to the **specific concrete type's** extra fields or methods that aren't part of the interface itself.

## 3. What problem does it solve?

```text
Without type assertion:
Once a value is "boxed" into an interface, you can only call
the interface's own methods — you lose access to extra fields/methods.

With type assertion:
You can safely "unbox" the interface back to its concrete type
when you need that type's specific capabilities.
```

## 4. How does it work?

```text
x.(T)        -> PANICS if x does not hold a T. (single-result form)
v, ok := x.(T) -> SAFE form: ok is false if it doesn't hold T, v is zero value.
```

## 5. Simple Mental Model

```text
Type assertion = "Let me check what's really inside this box,
                   and open it if it matches what I expect."
```

## 6. Basic Go Example

```go
package main

import "fmt"

func main() {
    var i interface{} = "hello"

    s, ok := i.(string)
    fmt.Println(s, ok) // hello true

    n, ok := i.(int)
    fmt.Println(n, ok) // 0 false (safe, no panic)

    // n2 := i.(int) // PANIC! this is the unsafe single-result form
}
```

## 7. Explain the Code

```text
1. i is an interface{} (holds any type) containing a string "hello".
2. i.(string) with two results safely checks AND extracts the value.
3. i.(int) fails safely because i doesn't hold an int -> ok is false.
4. Using the single-result form (i.(int) without ", ok") would PANIC here.
```

## 8. Real-Life Problem

```text
Handling errors from different libraries:

if pathErr, ok := err.(*os.PathError); ok {
    fmt.Println("failed path:", pathErr.Path)
} else {
    fmt.Println("generic error:", err)
}
```
This lets you react to SPECIFIC error types with extra detail, while still handling generic errors gracefully.

## 9. When should I use it?

When you need to check/access a specific concrete type from an interface value — especially for error handling, or working with `interface{}`/`any` values from JSON decoding, plugin systems, etc.

## 10. When should I NOT use it?

Avoid overusing type assertions as a substitute for good interface design — if you're constantly asserting types to get at behavior, you probably need a better interface or a type switch (7.13) instead.

## 11. Common Mistakes

- Using the single-result form `x.(T)` in code where the type isn't guaranteed — causes a runtime panic.
- Forgetting that a **failed assertion with the 2-result form** returns the **zero value**, not `nil` (important for numeric/struct types).

## 12. Important Gotchas

- Type assertions check **exact concrete type match** (or interface implementation, if `T` is itself an interface).
- `errors.As` (in `errors` package) is essentially a smarter, chain-aware version of a type assertion for errors (see 7.11).

## 13. Internals

### Go Language Guarantee
`x.(T)` succeeds if `x`'s dynamic type is exactly `T` (when `T` is concrete), or if `x`'s dynamic type implements `T` (when `T` is an interface).

### Implementation Detail
The runtime compares the interface value's stored type pointer against `T`'s type descriptor — a fast pointer comparison in most implementations.

## 14. Standard Library Connection

```go
if pathErr, ok := err.(*fs.PathError); ok { ... }
```
Extremely common pattern across the standard library for extracting specific error details.

## 15. Production Example

```go
func handleError(err error) {
    if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
        log.Println("network timeout, retrying...")
        return
    }
    log.Println("unhandled error:", err)
}
```

## 16. Performance

Type assertions are very fast (a pointer/type comparison) — no meaningful performance concern in typical backend code.

## 17. Related Concepts

| Concept | Meaning |
|---|---|
| Type assertion | Extract/check ONE specific type from an interface |
| Type switch | Handle MULTIPLE possible types in one construct |
| errors.As | Chain-aware type assertion for wrapped errors |

## 18. Interview Questions

**Basic** — Q: What does `v, ok := x.(T)` do? A: Safely checks if `x` holds a `T`; if yes, `v` is the value and `ok` is true, else `v` is the zero value and `ok` is false.

**Intermediate** — Q: What happens with the single-result form if the assertion fails? A: It panics.

**Advanced** — Q: How does a type assertion behave differently if `T` is an interface vs a concrete type? A: If `T` is concrete, it checks for an EXACT dynamic type match. If `T` is an interface, it checks whether the dynamic type implements that interface.

**Tricky** — Q: If `x` is `nil`, what does `x.(T)` return with the two-result form? A: `ok` is false and `v` is the zero value of `T` — no panic, since the safe form handles this gracefully.

## 19. Interview Follow-Up Questions

```text
Q: What is a type assertion?
Q: Safe form vs unsafe form — differences?
Q: What happens when T is an interface, not a concrete type?
Q: How does this relate to errors.As?
```

## 20. How to Explain This in an Interview

> "A type assertion lets you extract the concrete value stored inside an interface, or check whether it matches a specific type. The safe two-result form, v, ok := x.(T), never panics — ok tells you if the assertion succeeded. The single-result form panics if it fails, so it should only be used when you're certain of the type. This is heavily used for handling specific error types while still gracefully falling back for generic ones."

## 21. Quick Revision

```text
WHAT?      -> Extract/check a specific concrete type from an interface.
WHY?       -> Access type-specific fields/methods hidden by the interface.
HOW?       -> v, ok := x.(T) (safe) vs x.(T) (panics on failure).
GOTCHA?    -> Failed 2-result assertion gives zero value, not nil.
INTERVIEW? -> Concrete T = exact match; interface T = implements check.
```

## 22. Practice

> Write a function that accepts `interface{}` and prints a custom message for `int`, `string`, and everything else, using type assertions (then compare with type switch in 7.13).

---

# 7.11 Discriminating Errors with Type Assertions

## 1. What is it?

This is the **pattern** of using type assertions (or `errors.As`) specifically to figure out **which kind of error** occurred, so your code can react differently to different failure types.

## 2. Why do we need it?

Not all errors should be handled the same way. A "file not found" error might mean "create the file," while a "permission denied" error might mean "alert an admin." You need to **discriminate** (tell apart) between error types to react correctly.

## 3. What problem does it solve?

```text
Without discrimination:
if err != nil { log.Fatal(err) }  <- treats EVERY error the same way.

With discrimination:
Different errors get different, appropriate handling —
retry, ignore, alert, fallback, etc.
```

## 4. How does it work?

```text
1. A function returns different concrete error types for different failures.
2. Caller uses type assertion (or type switch) to check WHICH type it is.
3. Caller reacts differently based on the result.
```

Modern Go (1.13+) prefers `errors.Is` / `errors.As` because they understand **wrapped** errors (errors wrapped with `%w`), while a plain type assertion only checks the outermost error.

## 5. Simple Mental Model

```text
Discriminating errors = "Don't just ask IF something failed —
                          ask WHAT KIND of failure it was."
```

## 6. Basic Go Example

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

func main() {
    _, err := os.Open("missing.txt")

    var pathErr *os.PathError
    if errors.As(err, &pathErr) {
        fmt.Println("Path error on:", pathErr.Path)
    } else {
        fmt.Println("Unknown error:", err)
    }
}
```

## 7. Explain the Code

```text
1. os.Open fails because the file doesn't exist.
2. errors.As checks if err (possibly wrapped) contains a *os.PathError.
3. If yes, pathErr is populated and we can access its Path field.
4. errors.As works even through multiple layers of wrapped errors — unlike a plain type assertion.
```

## 8. Real-Life Problem

```text
A backend service calling a payment gateway needs to:

- Retry on TimeoutError
- Alert on FraudError
- Return "insufficient funds" message on InsufficientFundsError
- Log everything else generically

var timeoutErr *TimeoutError
var fraudErr *FraudError
switch {
case errors.As(err, &timeoutErr):
    retry()
case errors.As(err, &fraudErr):
    alertSecurityTeam()
default:
    logGeneric(err)
}
```

## 9. When should I use it?

Whenever different error causes need genuinely different handling logic — retries, user-facing messages, alerting, fallback behavior.

## 10. When should I NOT use it?

Don't build elaborate error-type hierarchies for a simple script where `if err != nil { return err }` is sufficient. Over-engineering error handling adds complexity without real benefit for small tools.

## 11. Common Mistakes

- Using `==` to compare wrapped errors (`err == someSentinelErr`) instead of `errors.Is`.
- Using a plain type assertion (`err.(*MyError)`) instead of `errors.As` when the error might be wrapped — the plain assertion will fail to see through wrapping.

## 12. Important Gotchas

- `errors.Is` checks for a specific **error value** (sentinel errors like `sql.ErrNoRows`).
- `errors.As` checks for a specific **error type** (like `*os.PathError`) and extracts it.
- Both walk the "wrapped error chain" created by `fmt.Errorf("...: %w", err)`.

## 13. Internals

### Go Language Guarantee
`errors.Is`/`errors.As` use the `Unwrap() error` method (if present) to walk through a chain of wrapped errors, checking each layer.

### Implementation Detail
`fmt.Errorf` with `%w` internally creates a wrapped error struct implementing `Unwrap()` — this mechanism was added in Go 1.13.

## 14. Standard Library Connection

```text
sql.ErrNoRows, io.EOF, os.ErrNotExist
-- classic sentinel errors checked with errors.Is throughout the stdlib and real codebases.
```

## 15. Production Example

```go
var ErrNotFound = errors.New("not found")

func FindUser(id int64) (*User, error) {
    if !exists(id) {
        return nil, fmt.Errorf("finding user %d: %w", id, ErrNotFound)
    }
    // ...
}

// caller:
if errors.Is(err, ErrNotFound) {
    return c.JSON(404, "user not found")
}
```

## 16. Performance

Negligible — error-chain walking is a short linked-list traversal, not a performance concern.

## 17. Related Concepts

| Concept | Meaning |
|---|---|
| errors.Is | Compare against a specific sentinel error VALUE |
| errors.As | Extract a specific error TYPE from the chain |
| %w in fmt.Errorf | Wraps an error, preserving the chain for Is/As |

## 18. Interview Questions

**Basic** — Q: What's the difference between `errors.Is` and `errors.As`? A: `Is` checks for a specific error value; `As` checks for and extracts a specific error type.

**Intermediate** — Q: Why is `errors.As` preferred over a plain type assertion for error handling? A: Because it can see through wrapped errors (via `Unwrap()`), while a plain assertion only checks the immediate outer error.

**Advanced** — Q: How does error wrapping work internally in Go? A: `fmt.Errorf("...: %w", err)` creates a wrapping error type with an `Unwrap() error` method; `errors.Is`/`As` repeatedly call `Unwrap()` to walk the chain.

**Tricky** — Q: Can you compare two wrapped errors with `==`? A: No — wrapping creates a new error value, so direct `==` comparison will fail even if the underlying cause is the same; use `errors.Is` instead.

## 19. Interview Follow-Up Questions

```text
Q: What does "discriminating errors" mean?
Q: Difference between errors.Is and errors.As?
Q: How does error wrapping/unwrapping work?
Q: Why not use a plain type assertion for error handling?
```

## 20. How to Explain This in an Interview

> "Discriminating errors means reacting differently based on WHAT kind of error occurred, not just whether one occurred. In modern Go, we do this with errors.Is for checking a specific sentinel error value, and errors.As for extracting a specific concrete error type — both understand wrapped errors created with %w, unlike a plain type assertion, which only sees the outermost error."

## 21. Quick Revision

```text
WHAT?      -> Reacting differently based on the kind of error.
WHY?       -> Different failures need different handling (retry, alert, etc).
HOW?       -> errors.Is (value) / errors.As (type), both wrap-aware.
GOTCHA?    -> Plain == or type assertion breaks with wrapped errors.
INTERVIEW? -> Explain Unwrap() chain mechanism clearly.
```

## 22. Practice

> Create 2 custom error types (`RateLimitError`, `AuthError`). Write a handler function using `errors.As` in a switch-like structure to react to each differently, including when wrapped with `fmt.Errorf("...: %w", err)`.

---

# 7.12 Querying Behaviors with Interface Type Assertions

## 1. What is it?

This is the pattern of using a type assertion to check if a value **optionally** supports **extra** behavior beyond its primary interface — "does this ALSO support X?"

## 2. Why do we need it?

Sometimes a function receives a general interface (like `io.Writer`), but wants to use a **more specific, optional capability** IF the concrete value happens to support it (like `io.StringWriter` for a more efficient string write).

## 3. What problem does it solve?

```text
Without optional behavior querying:
You'd need a giant interface with every possible method,
forcing every implementer to support features they may not need.

With it:
The main interface stays small and required.
Extra capabilities are optional and checked only when useful.
```

## 4. How does it work?

```text
1. Accept the general interface (e.g. io.Writer) as your parameter.
2. Internally, try a type assertion to a MORE SPECIFIC interface (e.g. io.StringWriter).
3. If the assertion succeeds, use the more efficient/specific method.
4. If it fails, fall back to the general method.
```

## 5. Simple Mental Model

```text
"I need you to be able to Write().
 But if you ALSO happen to know how to WriteString() efficiently,
 I'll use that instead."
```

## 6. Basic Go Example

```go
package main

import (
    "fmt"
    "io"
    "strings"
)

func writeGreeting(w io.Writer, name string) {
    msg := "Hello, " + name

    if sw, ok := w.(io.StringWriter); ok {
        sw.WriteString(msg) // more efficient path
        fmt.Println("(used WriteString)")
        return
    }
    w.Write([]byte(msg)) // generic fallback
    fmt.Println("(used generic Write)")
}

func main() {
    var sb strings.Builder
    writeGreeting(&sb, "Aman")
    fmt.Println(sb.String())
}
```

## 7. Explain the Code

```text
1. writeGreeting takes a general io.Writer.
2. It asserts whether w ALSO satisfies io.StringWriter.
3. strings.Builder DOES implement WriteString, so the efficient path runs.
4. If we passed something without WriteString, it would fall back to Write().
```

## 8. Real-Life Problem

```text
HTTP response writers sometimes optionally support http.Flusher
(for streaming responses like Server-Sent Events):

if flusher, ok := w.(http.Flusher); ok {
    flusher.Flush() // stream data immediately
}
```
This is exactly how Go's `net/http` package lets handlers optionally support flushing without forcing every `ResponseWriter` implementation to support it.

## 9. When should I use it?

When you want to offer an optimization or extra feature ONLY for implementations that support it, without bloating the main required interface.

## 10. When should I NOT use it?

Don't overuse this pattern to hide required behavior behind "optional" checks — if something is *always* needed, put it in the main interface instead of querying for it via assertion every time.

## 11. Common Mistakes

- Forgetting the fallback path — always handle the case where the optional interface ISN'T supported.
- Overusing this pattern instead of properly designing interface hierarchies.

## 12. Important Gotchas

- This pattern is sometimes called **"optional interfaces"** — very idiomatic in Go's standard library (`http.Flusher`, `http.Hijacker`, `io.StringWriter`, `io.ReaderFrom`).
- The assertion cost is small but non-zero — don't do it in extremely tight loops if avoidable.

## 13. Internals

### Go Language Guarantee
A type assertion to check for optional interface support behaves exactly like any other interface-to-interface type assertion — succeeds only if the concrete type actually implements it.

### Implementation Detail
No special runtime mechanism — it's a normal type assertion, just used in an "optional feature check" style.

## 14. Standard Library Connection

```text
io.StringWriter, io.ReaderFrom, io.WriterTo
http.Flusher, http.Hijacker, http.Pusher (HTTP/2)
```
All are "optional" interfaces queried via type assertion inside generic code paths.

## 15. Production Example

```go
func copyEfficient(dst io.Writer, src io.Reader) (int64, error) {
    if rf, ok := dst.(io.ReaderFrom); ok {
        return rf.ReadFrom(src) // let dst do the copy efficiently
    }
    return io.Copy(dst, src) // generic fallback
}
```
This is essentially how `io.Copy` itself is implemented internally in the standard library.

## 16. Performance

This pattern EXISTS specifically for performance — letting code opt into a faster path when available, while still working correctly (just slower) when it's not.

## 17. Related Concepts

| Concept | Meaning |
|---|---|
| Optional interface | An extra capability, checked via assertion, not required |
| io.Copy | Real stdlib function using exactly this technique |
| Type assertion | The underlying mechanism used to check |

## 18. Interview Questions

**Basic** — Q: What is an "optional interface" in Go? A: An extra capability a type MAY support, checked at runtime with a type assertion, beyond its main required interface.

**Intermediate** — Q: Give a real standard library example. A: `http.Flusher` — an HTTP handler checks if `ResponseWriter` also supports `Flush()` to stream data early.

**Advanced** — Q: Why does Go prefer this pattern over one giant interface with all possible methods? A: It keeps the primary interface small (easy to implement/mock) while still enabling performance optimizations or extra features only where actually supported — following Go's "small interfaces" philosophy.

**Tricky** — Q: What happens if you skip the fallback path for an optional interface check? A: The code will break/panic or silently do nothing for implementations that don't support the optional interface — always write a fallback.

## 19. Interview Follow-Up Questions

```text
Q: What is an "optional interface" pattern?
Q: Real examples from net/http or io packages?
Q: Why not just require the extra method in the main interface?
Q: What must you always include when using this pattern?
```

## 20. How to Explain This in an Interview

> "Sometimes a function accepts a general interface but wants to use a more specific capability if the concrete value happens to support it — like checking if an io.Writer also implements io.StringWriter for a faster write path. You do this with a type assertion to the more specific interface, with a fallback to the general behavior if it's not supported. It's how io.Copy and HTTP streaming with Flusher work internally, and it keeps interfaces small while still allowing optimizations."

## 21. Quick Revision

```text
WHAT?      -> Check if a value ALSO supports extra optional behavior.
WHY?       -> Keep main interfaces small; allow optional optimizations.
HOW?       -> Type assertion to a more specific interface + fallback.
REAL USE?  -> http.Flusher, io.ReaderFrom, io.StringWriter.
INTERVIEW? -> Mention io.Copy's internal use of this exact pattern.
```

## 22. Practice

> Write a `writeMessage(w io.Writer, msg string)` function that checks for `io.StringWriter` support first, falling back to `Write` otherwise. Test with both `strings.Builder` and `bytes.Buffer`.

---

# 7.13 Type Switches

## 1. What is it?

A **type switch** is a special `switch` statement that checks the **dynamic type** of an interface value against multiple possible types in one clean construct.

```go
switch v := x.(type) {
case int:
    // v is int here
case string:
    // v is string here
default:
    // unknown type
}
```

## 2. Why do we need it?

Without a type switch, checking multiple possible types with type assertions means a long chain of `if/else if` blocks — a type switch makes this much cleaner and more readable.

## 3. What problem does it solve?

```text
Without type switch:

if s, ok := x.(string); ok {
    ...
} else if n, ok := x.(int); ok {
    ...
} else if b, ok := x.(bool); ok {
    ...
}
-- repetitive and hard to read.

With type switch:
One clean construct handles all cases.
```

## 4. How does it work?

```text
1. switch v := x.(type) { ... } — special syntax, only valid in switch.
2. Each "case TYPE:" checks if x's dynamic type matches TYPE.
3. Inside each case, "v" automatically has THAT specific type.
4. "default:" catches anything that didn't match any case.
```

## 5. Simple Mental Model

```text
Type switch = "Look inside the box, and run different code
               depending on exactly WHAT is inside."
```

## 6. Basic Go Example

```go
package main

import "fmt"

func describe(x interface{}) {
    switch v := x.(type) {
    case int:
        fmt.Println("It's an int:", v*2)
    case string:
        fmt.Println("It's a string of length:", len(v))
    case bool:
        fmt.Println("It's a bool:", !v)
    case nil:
        fmt.Println("It's nil")
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}

func main() {
    describe(42)
    describe("hello")
    describe(true)
    describe(nil)
    describe(3.14)
}
```

## 7. Explain the Code

```text
1. x.(type) is special syntax ONLY valid inside a switch statement.
2. Go checks x's dynamic type against each "case" in order.
3. Inside case int, "v" is already typed as int (no extra assertion needed).
4. "case nil" specifically catches a truly nil interface value.
5. "default" handles any type we didn't explicitly list.
```

## 8. Real-Life Problem

```text
JSON decoding into interface{} often needs a type switch
to figure out what kind of JSON value you actually got:

switch val := data.(type) {
case string:
    // handle string
case float64:      // JSON numbers decode as float64 by default!
    // handle number
case map[string]interface{}:
    // handle nested object
case []interface{}:
    // handle array
case nil:
    // handle JSON null
}
```

## 9. When should I use it?

When you need to handle several DIFFERENT possible concrete types stored in an interface, especially 3+ cases — cleaner than chained type assertions.

## 10. When should I NOT use it?

If you only need to check ONE specific type, a plain type assertion (`v, ok := x.(T)`) is simpler and clearer than a type switch with just one case + default.

## 11. Common Mistakes

- Forgetting the `default` case, silently ignoring unexpected types.
- Confusing `case nil:` (checks for a nil interface) with checking a nil pointer stored inside a non-nil interface (see the gotcha from 7.5 — a type switch's `nil` case only matches a TRULY nil interface).

## 12. Important Gotchas

- You CAN combine multiple types in one case (`case int, int64:`), but then `v`'s type inside that case is the ORIGINAL interface type (`interface{}`), not narrowed — because Go can't know which of the two types it actually is.
- Order matters if types overlap through interface satisfaction (e.g., checking a specific interface type before a broader one).

## 13. Internals

### Go Language Guarantee
Each `case` in a type switch is checked against the dynamic type in the order written; the first match wins.

### Implementation Detail
Under the hood, the compiler generates a similar sequence of type comparisons as a chain of type assertions — a type switch is essentially syntactic sugar for that.

## 14. Standard Library Connection

```text
go/ast walking, encoding/json decoding into interface{},
fmt package's internal formatting logic —
all commonly use type switches to handle multiple concrete types.
```

## 15. Production Example

```go
func handleEvent(e Event) {
    switch ev := e.(type) {
    case *UserCreatedEvent:
        sendWelcomeEmail(ev.UserID)
    case *OrderPlacedEvent:
        notifyWarehouse(ev.OrderID)
    case *PaymentFailedEvent:
        alertBillingTeam(ev.OrderID, ev.Reason)
    default:
        log.Printf("unhandled event type: %T", ev)
    }
}
```
Common pattern in event-driven backend systems.

## 16. Performance

Similar cost to a sequence of type assertions — small and rarely a bottleneck. Order your most common cases first for a tiny optimization if you have MANY cases in a very hot path.

## 17. Related Concepts

| Concept | Meaning |
|---|---|
| Type switch | Handle MULTIPLE possible types cleanly |
| Type assertion | Check/extract ONE specific type |
| case nil | Matches a truly nil interface value |

## 18. Interview Questions

**Basic** — Q: What is a type switch used for? A: Checking an interface value's dynamic type against multiple possible types in one clean construct.

**Intermediate** — Q: What's the type of `v` inside `case int, int64:`? A: The original interface type, NOT narrowed to `int` or `int64`, because Go doesn't know which of the two it is.

**Advanced** — Q: How does `case nil:` behave differently from checking a nil pointer of a concrete type? A: `case nil:` matches only when the interface itself is truly nil (both dynamic type and value nil) — it will NOT match a non-nil interface holding a nil concrete pointer (the classic typed-nil gotcha from 7.5).

**Tricky** — Q: Is a type switch essentially syntactic sugar for something else? A: Yes — it behaves like a sequence of type assertions checked in order, just with cleaner syntax and automatic narrowing per case.

## 19. Interview Follow-Up Questions

```text
Q: What is a type switch?
Q: How does it differ from a normal switch on values?
Q: What happens with multiple types in one case?
Q: How does "case nil" interact with the typed-nil gotcha?
```

## 20. How to Explain This in an Interview

> "A type switch lets you cleanly handle multiple possible concrete types stored in an interface value, instead of chaining several type assertions. Inside each case, the variable is automatically narrowed to that specific type — except when a case lists multiple types together, in which case it stays as the original interface type. It's heavily used for things like JSON decoding into interface{}, or dispatching on event types in event-driven systems."

## 21. Quick Revision

```text
WHAT?      -> switch on the dynamic type of an interface value.
WHY?       -> Cleaner than chained type assertions for multiple types.
HOW?       -> switch v := x.(type) { case T1: ... case T2: ... }
GOTCHA?    -> Multi-type case doesn't narrow "v"; case nil is strict.
INTERVIEW? -> Explain it as sugar over sequential type assertions.
```

## 22. Practice

> Write a function `describeJSON(v interface{})` that handles `string`, `float64`, `bool`, `map[string]interface{}`, `[]interface{}`, and `nil` using a type switch.

---

# 7.14 Example: Token-Based XML Decoding

## 1. What is it?

This is a worked example of using Go's `encoding/xml` **token-based streaming decoder** (`xml.Decoder`), where each XML token (start element, end element, character data) is represented as a value satisfying the `xml.Token` interface.

## 2. Why do we need it?

Loading an ENTIRE XML document into memory (like `xml.Unmarshal` into a struct) doesn't work well for HUGE XML files, or when you need custom, low-level, streaming control over parsing. Token-based decoding processes the document piece by piece.

## 3. What problem does it solve?

```text
Without token-based decoding:
Large XML files must be fully loaded into memory before processing —
expensive for gigabyte-sized files, or impossible on memory-constrained systems.

With token-based decoding:
You read and process ONE token at a time,
using constant memory regardless of file size.
```

## 4. How does it work?

```text
type Token interface{}  // actually satisfied by: StartElement, EndElement, CharData, Comment, etc.

1. Create an xml.Decoder from an io.Reader.
2. Call decoder.Token() repeatedly in a loop.
3. Each call returns the NEXT token + an error (io.EOF when done).
4. Use a type switch to check what kind of token you got.
```

## 5. Simple Mental Model

```text
Token-based decoding = "Read the XML like a stream of small events,
                         one piece at a time, instead of loading it all at once."
```

## 6. Basic Go Example

```go
package main

import (
    "encoding/xml"
    "fmt"
    "strings"
)

func main() {
    data := `<book><title>Go in Action</title></book>`
    decoder := xml.NewDecoder(strings.NewReader(data))

    for {
        tok, err := decoder.Token()
        if err != nil {
            break // io.EOF or real error
        }
        switch t := tok.(type) {
        case xml.StartElement:
            fmt.Println("Start:", t.Name.Local)
        case xml.EndElement:
            fmt.Println("End:", t.Name.Local)
        case xml.CharData:
            text := strings.TrimSpace(string(t))
            if text != "" {
                fmt.Println("Text:", text)
            }
        }
    }
}
```

## 7. Explain the Code

```text
1. xml.NewDecoder wraps any io.Reader (file, network stream, string, etc).
2. decoder.Token() reads ONE token at a time — memory-efficient.
3. tok is of type xml.Token — an interface satisfied by several concrete types.
4. The type switch checks WHICH kind of XML token we just read.
5. Loop continues until an error (typically io.EOF) signals "no more tokens."
```

## 8. Real-Life Problem

```text
A backend system ingesting large third-party XML feeds
(e.g. product catalogs, bank statement exports, RSS feeds)
CANNOT load a multi-gigabyte file entirely into memory.

Token-based decoding processes the feed in a streaming fashion,
extracting only the fields it needs, using constant memory.
```

## 9. When should I use it?

For very large XML documents, streaming data sources (network feeds), or when you need fine-grained custom parsing control that `xml.Unmarshal` doesn't give you.

## 10. When should I NOT use it?

For small, simple XML documents that fit comfortably in memory, `xml.Unmarshal` into a Go struct is MUCH simpler and less error-prone. Only reach for token-based decoding when memory/streaming genuinely matters.

## 11. Common Mistakes

- Forgetting to check for `io.EOF` correctly (it's a normal, expected "end of stream" signal, not necessarily a real failure).
- Not handling nested elements correctly — you often need a stack-like structure to track "where you are" in the XML tree while streaming.

## 12. Important Gotchas

- `xml.CharData` often includes whitespace/newlines between tags — you usually need to `TrimSpace` it.
- The `xml.Token` "interface" is really just `interface{}` (an alias historically) with several concrete types satisfying it via convention — this is a great real example of a type switch being ESSENTIAL, not optional, because there's no shared method to call otherwise.

## 13. Internals

### Go Language Guarantee
`decoder.Token()` returns tokens strictly in the order they appear in the underlying XML document.

### Implementation Detail
Internally, the decoder reads and buffers just enough bytes from the `io.Reader` to produce the next token — the exact buffering strategy is an implementation detail that can change between Go versions.

## 14. Standard Library Connection

This entire topic IS the standard library (`encoding/xml`) — showing a real production-grade use of interfaces + type switches together for streaming data processing.

## 15. Production Example

```go
func extractPrices(r io.Reader) ([]float64, error) {
    var prices []float64
    decoder := xml.NewDecoder(r)
    inPrice := false

    for {
        tok, err := decoder.Token()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }
        switch t := tok.(type) {
        case xml.StartElement:
            inPrice = t.Name.Local == "price"
        case xml.CharData:
            if inPrice {
                var p float64
                fmt.Sscanf(string(t), "%f", &p)
                prices = append(prices, p)
            }
        }
    }
    return prices, nil
}
```

## 16. Performance

Token-based decoding uses **constant memory** regardless of document size — a major performance win over full unmarshaling for very large documents. CPU cost is similar or slightly higher due to manual token handling.

## 17. Related Concepts

| Concept | Meaning |
|---|---|
| xml.Token | Interface satisfied by StartElement, EndElement, CharData, etc. |
| xml.Unmarshal | Simpler, whole-document decoding into a struct |
| Streaming parsing | Processing data piece-by-piece instead of all at once |

## 18. Interview Questions

**Basic** — Q: What does `xml.NewDecoder(r).Token()` return? A: The next XML token (start element, end element, char data, etc.) and an error.

**Intermediate** — Q: Why use token-based decoding instead of `xml.Unmarshal`? A: For very large or streaming documents where loading everything into memory at once isn't practical.

**Advanced** — Q: How do you track your position in a nested XML tree while streaming tokens? A: Typically with a stack (slice) that pushes on `StartElement` and pops on `EndElement`, letting you know the current nesting context.

**Tricky** — Q: Why is a type switch essential here, rather than optional? A: Because `xml.Token` doesn't define any shared methods to call — the only way to tell tokens apart is by checking their concrete type via a type switch.

## 19. Interview Follow-Up Questions

```text
Q: What is xml.Token?
Q: Why choose token-based decoding over Unmarshal?
Q: How do you track nesting depth while streaming?
Q: Why is a type switch necessary here specifically?
```

## 20. How to Explain This in an Interview

> "Token-based XML decoding uses xml.Decoder to read an XML document as a stream of tokens — start elements, end elements, character data — instead of loading the whole document into memory. You call decoder.Token() in a loop and use a type switch to handle each kind of token. It's ideal for very large or streaming XML sources where full unmarshaling into a struct isn't practical, and it's a great real example of using a type switch as the PRIMARY way to distinguish behavior, since the token interface itself has no shared methods."

## 21. Quick Revision

```text
WHAT?      -> Stream XML as tokens instead of loading it all at once.
WHY?       -> Constant memory usage for huge/streaming documents.
HOW?       -> decoder.Token() in a loop + type switch on token kind.
GOTCHA?    -> CharData often needs trimming; track nesting with a stack.
INTERVIEW? -> Explain why a type switch is essential here.
```

## 22. Practice

> Write a token-based decoder that counts how many `<item>` elements appear in a large XML document, without ever holding the whole document in memory as a struct.

---

# 7.15 A Few Words of Advice

## 1. What is it?

This is Go's own design philosophy around interfaces — a set of practical guidelines for using interfaces well, instead of overusing them.

## 2. The Core Advice

```text
1. Only introduce a new interface when you actually need ONE (or more)
   concrete types to satisfy it. Don't create interfaces "just in case."

2. Define interfaces where they are USED (consumer side),
   not where the concrete type is defined (producer side).

3. Keep interfaces SMALL. 1-3 methods is common and healthy.
   Small interfaces are easy to implement, easy to mock, easy to understand.

4. "Accept interfaces, return structs" — function parameters should be
   flexible (interfaces); return values should usually be concrete types.

5. Don't force interfaces onto every type. If there's only ONE
   implementation and no plan for more, a concrete type is simpler.
```

## 3. Why this matters (in simple English)

```text
Overusing interfaces makes code HARDER to read, not easier.
You end up jumping between an interface definition and its one single
implementation, adding indirection with zero real benefit.

Interfaces are a tool for FLEXIBILITY, not a rule to follow everywhere.
```

## 4. Real-Life Example

```text
BAD (interface with only ONE implementation, defined near the struct,
     never actually swapped):

type UserServicer interface {
    GetUser(id int) (*User, error)
}

type UserService struct{}
func (UserService) GetUser(id int) (*User, error) { ... }

-- if UserService NEVER has a second implementation and is never mocked,
   this interface adds no value. Just use UserService directly.


GOOD (interface defined where it's CONSUMED, because
      multiple implementations genuinely exist — Postgres, Mock):

// in the package that USES the repository:
type UserRepository interface {
    GetUser(id int) (*User, error)
}
```

## 5. When Interfaces DO Help

- Testing: you can pass a mock implementation instead of a real database/API.
- Swappable implementations: multiple storage backends, multiple payment providers.
- Package boundaries: keeping one package from depending directly on another package's concrete types.

## 6. When Interfaces DON'T Help

- Single implementation, no testing need for mocking, no plan to swap it.
- Adding "an interface for everything" out of habit from other languages (Java/C# where interfaces are used far more heavily even for single implementations).

## 7. Go's Philosophy Summed Up

```text
Go prefers:
  Small interfaces over big ones
  Composition over inheritance
  Explicit over implicit (except interface satisfaction itself, which is implicit by design)
  Simplicity over "just in case" flexibility
```

## 8. Interview Questions

**Basic** — Q: Where should interfaces usually be defined in Go — near the implementation or near the usage? A: Near the usage (consumer side).

**Intermediate** — Q: What does "accept interfaces, return structs" mean? A: Function parameters should be interfaces for flexibility; return values should usually be concrete types for clarity.

**Advanced** — Q: Why does Go discourage interfaces with only one implementation and no testing need? A: Because it adds indirection and mental overhead without providing any real flexibility benefit — it's unnecessary abstraction.

**Tricky** — Q: Is defining an interface "for future flexibility" that isn't needed yet good practice in Go? A: Generally no — Go favors adding the interface later, exactly when a second implementation or mocking need actually appears (YAGNI: "You Aren't Gonna Need It").

## 9. How to Explain This in an Interview

> "Go's advice on interfaces is: keep them small, define them where they're consumed rather than where they're implemented, and only introduce one when you genuinely need multiple implementations or the ability to mock something in tests. Overusing interfaces — especially ones with a single implementation — adds indirection without real benefit. The idiom 'accept interfaces, return structs' captures this well: be flexible about what you take in, but clear and concrete about what you give back."

## 10. Quick Revision

```text
WHAT?      -> Practical guidelines for using interfaces well.
WHY?       -> Prevent unnecessary abstraction/indirection.
RULE 1?    -> Define interfaces at the point of use, not near the type.
RULE 2?    -> Keep interfaces small (1-3 methods).
RULE 3?    -> Accept interfaces, return structs.
RULE 4?    -> Don't create interfaces "just in case."
INTERVIEW? -> Explain WHY, not just WHAT — shows real understanding.
```

---

# Final Revision — Chapter 7: Interfaces

## Most Important Concepts

```text
1. Interface = a contract of required methods; satisfaction is IMPLICIT.
2. Pointer vs value receivers change WHICH type satisfies an interface.
3. Interface value = (dynamic type, dynamic value) — nil only if BOTH are nil.
4. Typed-nil-in-interface is Go's #1 interfaces gotcha (error/*T bug).
5. Type assertion = check/extract ONE type. Type switch = handle MANY types.
6. error is just an interface — no exceptions, explicit checking.
7. errors.Is / errors.As are wrap-aware; plain == / assertions are not.
8. Small interfaces (io.Reader, io.Writer, sort.Interface, http.Handler,
   flag.Value) are Go's signature design style.
9. "Optional interfaces" pattern lets code opt into extra behavior safely.
10. Define interfaces at the point of USE, keep them small, avoid over-abstraction.
```

## Must Remember (Mental Models)

```text
Interface           = Contract / checklist of abilities
Interface value      = Labeled box (type label + data inside)
Type assertion       = "Let me check/open the box for one specific type"
Type switch          = "Open the box and branch on whatever's really inside"
Typed nil gotcha      = Box with a label but empty content is NOT an empty box
Optional interfaces   = "Can you also do X? If yes, great. If not, fallback."
```

## Common Traps (Danger Zone for Interviews)

```text
1. Returning a nil pointer of a concrete error type as `error`
   -> interface is NOT nil (classic bug).
2. Assuming a value with a pointer-receiver method satisfies
   an interface even when passed by value -> compile error.
3. Comparing wrapped errors with == instead of errors.Is.
4. Using unsafe single-result type assertion where failure is possible -> panic.
5. Case with multiple types in a type switch does NOT narrow the variable's type.
```

## Top Interview Questions

```text
Q1: What is an interface, and how is satisfaction determined in Go?
Q2: Why does a nil pointer stored in an interface make the interface non-nil?
Q3: Difference between type assertion and type switch?
Q4: How do pointer vs value receivers affect interface satisfaction?
Q5: Difference between errors.Is and errors.As, and why do they exist?
Q6: What is the "accept interfaces, return structs" idiom?
Q7: How does http.HandlerFunc let a plain function satisfy an interface?
Q8: Why does Go favor small interfaces over big ones?
```

## Advanced Questions (Strong Backend Interviewer Level)

```text
Q1: Walk through exactly what happens in memory when you assign
    a concrete struct to an interface variable.
Q2: Design a plugin-style system in Go using a small interface +
    optional interfaces for extra capabilities (like io.ReaderFrom).
Q3: Explain how Go's own compiler uses interfaces in go/ast for
    representing source code as a tree.
Q4: How would you design error types for a backend service so that
    callers can both log generically AND react to specific failures?
Q5: Why might introducing an interface for a struct with a single
    implementation actually make code WORSE, not better?
```

## One-Minute Chapter Explanation

> "Interfaces in Go are contracts made of method signatures. Any type that implements all those methods automatically satisfies the interface — no 'implements' keyword needed. Under the hood, an interface value is a pair: a dynamic type and a dynamic value, and it's only nil when both are nil — which causes Go's most famous gotcha, where a nil pointer of a concrete error type becomes a non-nil error interface. Type assertions let you extract or check for one specific type from an interface; type switches let you cleanly handle several possible types at once. Go's standard library builds almost everything — sorting, HTTP handling, error handling, flag parsing, XML decoding — around small, focused interfaces like sort.Interface, http.Handler, error, and flag.Value. The overall philosophy is: keep interfaces small, define them where they're used, and only introduce one when you genuinely need multiple implementations or testability — not just because Go supports them."

---

*End of Chapter 7 — Interfaces. Practice each code challenge before moving to the next Go chapter.*