# Go (Golang) Deep Study Guide — Chapter 6: Methods

> Simple English. Deep technical understanding. Real backend examples. Interview-ready.
> Companion file to Chapter 4 (Composite Types) and Chapter 5 (Functions) — same format, same depth.

---

# 6. Methods

Go doesn't have classes, but it has **methods** — functions attached to a specific type. This chapter covers how methods are declared, the crucial choice between value and pointer receivers, how Go achieves code reuse through struct embedding instead of inheritance, treating methods as values, a full worked example, and how Go does encapsulation without `private`/`public` keywords.

Covered here: Method Declarations, Methods with a Pointer Receiver, Composing Types by Struct Embedding, Method Values and Expressions, Example: Bit Vector Type, Encapsulation.

---

## 6.1 Method Declarations

### 1. What is it?
```text
A method is a function attached to a specific type, using a
special "receiver" parameter written between `func` and the method name.
```

### 2. Why do we need it?
Go doesn't have classes, but real programs still need to group behavior with the data it operates on — a `User` needs a way to validate itself, a `Point` needs a way to compute distance to another point. Methods let you attach that behavior directly to the type, called with the familiar `value.Method()` syntax.

### 3. What problem does it solve?
```text
Without methods:
distance := calculateDistance(p1, p2) // behavior is disconnected from the type

With methods:
distance := p1.DistanceTo(p2) // behavior lives on the type itself,
                                // reads naturally, and is discoverable
```

### 4. How does it work?
```text
func (receiverName ReceiverType) MethodName(params) returnType {
    // body — can access receiverName like a normal parameter
}
```
The **receiver** (`receiverName ReceiverType`) is what makes this a method instead of a plain function — it says "this function belongs to `ReceiverType`."

### 5. Simple Mental Model
```text
Method = a regular function, except it's "attached" to a type,
so you call it as value.Method() instead of Method(value).
```

### 6. Basic Go Example
```go
package main

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

func (p Point) DistanceTo(q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

func main() {
	p1 := Point{0, 0}
	p2 := Point{3, 4}
	fmt.Println(p1.DistanceTo(p2)) // 5
}
```

### 7. Explain the Code
```text
1. func (p Point) DistanceTo(q Point) float64 declares a method
   on the Point type. `p` is the receiver — it's `p1` when we call p1.DistanceTo(p2).
2. Inside the method, `p` and `q` are just normal values, usable like
   any function parameters.
3. p1.DistanceTo(p2) reads naturally: "p1, give me your distance to p2."
```

### 8. Real-Life Problem
```text
Backend example: attaching validation and formatting behavior
directly to a domain type.

type Money struct {
    Cents int64
}

func (m Money) Dollars() float64 {
    return float64(m.Cents) / 100
}

func (m Money) String() string {
    return fmt.Sprintf("$%.2f", m.Dollars())
}
```
Now every place in the codebase that has a `Money` value can call `.String()` or `.Dollars()` directly, instead of scattering conversion logic everywhere.

### 9. When should I use it?
Whenever behavior conceptually "belongs to" a type — formatting, validation, calculations tied to that type's data, or satisfying an interface (Chapter 7).

### 10. When should I NOT use it?
If a function operates equally on multiple unrelated types, or doesn't naturally "belong" to any one type, a plain function is clearer than forcing it into a method on an arbitrary receiver.

### 11. Common Mistakes
- Confusing the receiver with a regular parameter conceptually — it's specifically what makes `value.Method()` syntax work.
- Declaring a method on a type defined in another package — Go **only allows methods on types defined in the same package** as the method.
- Naming the receiver inconsistently across methods on the same type (e.g., `p` in one method, `pt` in another) — Go convention is to keep the receiver name short and consistent.

### 12. Important Gotchas
- You can only declare a method on a **named type** defined in the current package — you cannot add a method to `int` or `string` directly, or to a type from another package, without first defining your own named type wrapping it.
- The receiver name is entirely up to you (not a keyword like `this` or `self`), but idiomatic Go uses a short abbreviation of the type name.
- Method sets differ between value receivers and pointer receivers — a crucial distinction fully explored in 6.2.

### 13. Internals
```text
Go Language Guarantee:
- A method is syntactic sugar: p.DistanceTo(q) is equivalent to
  calling a function with p passed as the first (receiver) argument.

Implementation Detail:
- The compiler generates the underlying function; the receiver
  syntax is purely a readability/organization feature at the
  language level, not a different execution mechanism from a
  plain function call.
```

### 14. Standard Library Connection
```text
time.Time has methods like t.Year(), t.Add(duration)
strings.Builder has methods like b.WriteString(s)
Nearly every standard library type exposes its behavior via methods.
```

### 15. Production Example
```go
type Order struct {
	Items []LineItem
}

func (o Order) Total() float64 {
	var total float64
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}
```
Business logic (`Total()`) lives right next to the data it operates on (`Order`), instead of being a loose function that has to be discovered separately.

### 16. Performance
Method calls in Go have essentially the same cost as regular function calls — no special runtime overhead is introduced just by using method syntax (this is different from, say, virtual method dispatch through an interface, covered in Chapter 7).

### 17. Related Concepts
| Concept | Meaning |
|---|---|
| Receiver | The `(p Point)` part that attaches a function to a type |
| Method set | The set of methods callable on a given value (depends on value vs pointer receiver, see 6.2) |
| Interface satisfaction (Ch. 7) | Types satisfy interfaces by having the right methods |

### 18. Interview Questions

**Basic**
- Q: What is a method in Go? A: A function with a special receiver parameter that attaches it to a specific type.
- Q: How do you call a method? A: `value.MethodName(args)`.

**Intermediate**
- Q: Can you add a method to a type from another package (e.g., `int`)? A: No — methods can only be declared on named types defined in your own package.

**Advanced**
- Q: Is `p.DistanceTo(q)` fundamentally different from a plain function call at the language level? A: No — it's syntactic sugar; the compiler treats it essentially as calling a function with `p` as an implicit first argument.

**Tricky**
- Q: If you want to add a method to `int`, how would you do it? A: Define a new named type based on `int` (e.g. `type Age int`) in your own package, and declare the method on that new type — you cannot add methods directly to `int` itself since it isn't defined in your package.

### 19. Interview Follow-Up Questions
```text
Q: What is a method, and how is it different from a plain function?
Q: What is a receiver?
Q: Can you declare a method on any type?
Q: How would you add "method-like" behavior to a built-in type like int?
Q: How does this connect to interface satisfaction? (bridges into Ch. 7)
```

### 20. Interview Answer
> "A method in Go is just a function with a receiver — a special parameter before the method name that attaches it to a specific type, letting me call it with `value.Method()` syntax. Under the hood it's essentially sugar for a regular function call with the receiver as the first argument. I use methods to keep behavior next to the data it operates on, like validation or formatting logic on a domain type, and methods are also how Go decides whether a type satisfies an interface."

### 21. Quick Revision
```text
WHAT?      → A function attached to a type via a receiver parameter
WHY?       → Keeps behavior next to the data it naturally belongs to
PROBLEM?   → Avoids disconnected free functions for type-specific logic
HOW?       → func (recv Type) Name(params) returnType { ... }
REAL USE?  → Money.Dollars(), Order.Total(), time.Time methods
GOTCHA?    → Can only declare methods on types in your own package
INTERVIEW? → Know that method syntax is sugar over a function call
```

### 22. Code Challenge
> Define a `Rectangle` struct with `Width` and `Height`, and add `Area()` and `Perimeter()` methods.

---

## 6.2 Methods with a Pointer Receiver

### 1. What is it?
```text
A pointer receiver method is declared with a pointer to the type
(e.g., `func (p *Point) Method()`) instead of the type itself —
meaning the method can modify the actual value the caller has.
```

### 2. Why do we need it?
Without a pointer receiver, a method receives a **copy** of the value (since Go is pass-by-value everywhere) — any changes inside the method are lost once it returns. A pointer receiver lets a method genuinely mutate the original value, and also avoids copying large structs on every call.

### 3. What problem does it solve?
```text
Without a pointer receiver:
func (p Point) MoveBy(dx, dy float64) {
    p.X += dx // modifies the COPY, caller's original Point is unchanged!
}

With a pointer receiver:
func (p *Point) MoveBy(dx, dy float64) {
    p.X += dx // modifies the ACTUAL Point the caller has
}
```

### 4. How does it work?
```text
p := Point{1, 1}
p.MoveBy(2, 2) // Go automatically takes &p here, since MoveBy needs a *Point

fmt.Println(p) // {3 3} — the original was actually modified
```
Go automatically converts `p.MoveBy(...)` to `(&p).MoveBy(...)` when needed, **as long as `p` is addressable** (a plain local variable is; a value returned directly from a function call, for example, is not).

### 5. Simple Mental Model
```text
Value receiver  = the method gets a photocopy — edits don't stick.
Pointer receiver = the method gets the original's address — edits DO stick.
```

### 6. Basic Go Example
```go
package main

import "fmt"

type Counter struct {
	count int
}

func (c *Counter) Increment() {
	c.count++
}

func main() {
	c := Counter{}
	c.Increment()
	c.Increment()
	fmt.Println(c.count) // 2
}
```

### 7. Explain the Code
```text
1. Increment has a pointer receiver (*Counter), so it can modify
   the real Counter, not a throwaway copy.
2. c.Increment() automatically becomes (&c).Increment() by Go.
3. Because the SAME underlying Counter is modified each call,
   count correctly reaches 2.
```

### 8. Real-Life Problem
```text
Backend example: a struct representing shared, mutable state,
like an in-memory cache entry's hit counter.

type CacheEntry struct {
    Value string
    Hits  int
}

func (e *CacheEntry) RecordHit() {
    e.Hits++
}
```
Every access to this cache entry needs to genuinely update the SAME hit counter — a value receiver here would silently do nothing useful.

### 9. When should I use it?
- When the method needs to **mutate** the receiver.
- When the receiver is a **large struct**, to avoid the cost of copying it on every method call.
- For consistency: if ANY method on a type needs a pointer receiver, it's Go convention to make **all** methods on that type use pointer receivers, for a predictable, consistent API.

### 10. When should I NOT use it?
For small, simple, immutable-in-spirit types (like a `Point{X, Y float64}` used read-only) where copying is cheap and mutation is never needed, value receivers are simpler and avoid any risk of unexpected aliasing.

### 11. Common Mistakes
- Mixing value and pointer receivers inconsistently across methods of the same type — this can cause confusing behavior about the type's interface satisfaction (see gotcha below).
- Expecting a value-receiver method to mutate the original — it never can; it only ever sees a copy.
- Trying to call a pointer-receiver method on a value that is **not addressable** (e.g., a map value directly: `myMap[key].Increment()` often fails to compile because map values aren't addressable).

### 12. Important Gotchas
```text
Method Set Rule:
- The method set of type T includes only methods with a VALUE receiver.
- The method set of type *T includes BOTH value and pointer receiver methods.
```
- This means: if a type has ANY pointer-receiver methods, only `*T` (not `T`) satisfies an interface requiring those methods — a very common interview trap (fully explored again in Chapter 7).
- Map values are **not addressable** in Go — you cannot call a pointer-receiver method directly on `myMap[key]`; you'd need to retrieve, modify, and reassign, or store pointers in the map instead (`map[string]*Counter`).
- Calling a pointer-receiver method on an addressable value (like a local variable) is automatically handled by Go (`c.Increment()` → `(&c).Increment()`), but this automatic conversion does NOT happen for values inside interfaces or non-addressable expressions.

### 13. Internals
```text
Go Language Guarantee:
- A pointer receiver method operates on the actual value via its address;
  a value receiver method operates on a copy.
- Automatic &value conversion happens only for addressable values.

Implementation Detail:
- The compiler decides, based on receiver type, whether to pass
  the receiver by value or by address at the call site — this
  compilation detail is invisible to you as long as you understand
  the addressability rules above.
```

### 14. Standard Library Connection
```text
bytes.Buffer's methods (like Write) use pointer receivers, since
writing to a buffer must mutate its internal state.
sync.Mutex's Lock/Unlock use pointer receivers for the same reason.
```

### 15. Production Example
```go
type Account struct {
	Balance float64
}

func (a *Account) Withdraw(amount float64) error {
	if amount > a.Balance {
		return errors.New("insufficient funds")
	}
	a.Balance -= amount
	return nil
}
```
A `Withdraw` method that didn't use a pointer receiver would be a serious bug — the balance would appear to update inside the method, but the caller's actual account balance would never change.

### 16. Performance
For large structs, a pointer receiver avoids copying the entire struct on every method call — this can matter for structs with many fields or embedded large data, though for small structs (a few small fields), the difference is negligible.

### 17. Related Concepts
| Concept | Meaning |
|---|---|
| Value receiver | Method gets a copy; cannot mutate the original |
| Pointer receiver | Method gets the address; can mutate the original |
| Method set | Which methods a value can call, depends on receiver types used |
| Addressability | Whether `&value` is legal — determines automatic conversion |

### 18. Interview Questions

**Basic**
- Q: What's the difference between a value receiver and a pointer receiver? A: A value receiver operates on a copy; a pointer receiver operates on the original value via its address, allowing mutation.
- Q: How do you declare a pointer receiver? A: `func (p *Type) Method() {...}`.

**Intermediate**
- Q: Why might you choose a pointer receiver even if the method doesn't mutate anything? A: To avoid copying a large struct on every call.

**Advanced**
- Q: What is the "method set" rule for value vs pointer types? A: `T`'s method set includes only value-receiver methods; `*T`'s method set includes both value- and pointer-receiver methods — meaning a type with any pointer-receiver methods requires `*T` to satisfy an interface needing those methods.

**Tricky**
- Q: Why can't you call a pointer-receiver method directly on a map value like `myMap[key].Increment()`? A: Map values are not addressable in Go — there's no stable memory address to take, since the map implementation may move values around internally, so the compiler refuses to let you take `&myMap[key]`.

### 19. Interview Follow-Up Questions
```text
Q: When do you use a pointer receiver vs a value receiver?
Q: What is the method set rule?
Q: Why does this matter for interface satisfaction?
Q: What does "addressable" mean, and where does it bite you?
Q: Why are map values not addressable?
```

### 20. Interview Answer
> "I use a pointer receiver whenever a method needs to mutate the receiver, or when the struct is large enough that copying it on every call would be wasteful. The rule I always keep in mind is the method set rule: a value type's method set only includes value-receiver methods, while a pointer type's method set includes both — so if any method on a type needs a pointer receiver, only a pointer to that type will satisfy an interface requiring it. I also watch out for addressability — you can't call a pointer-receiver method directly on a map value, since map values aren't addressable in Go."

### 21. Quick Revision
```text
WHAT?      → Receiver as *Type instead of Type — operates on the original
WHY?       → Enables mutation, and avoids copying large structs
PROBLEM?   → Value receivers can never actually change the caller's data
HOW?       → Go auto-converts p.M() to (&p).M() when p is addressable
REAL USE?  → Account.Withdraw, Counter.Increment, sync.Mutex.Lock
GOTCHA?    → Method set rule affects interface satisfaction; map values aren't addressable
INTERVIEW? → THE most common Go methods interview trap — know it cold
```

### 22. Code Challenge
> Convert the `Rectangle` from 6.1 into a `Scale(factor float64)` method using a pointer receiver that actually resizes the rectangle in place. Then try calling it on a value stored inside a `map[string]Rectangle` and observe/explain the compile error.

---

## 6.3 Composing Types by Struct Embedding

### 1. What is it?
```text
Struct embedding lets you include one struct type INSIDE another,
without giving it a field name — the outer struct then automatically
gets access to the inner struct's fields and methods, as if they
were its own.
```

### 2. Why do we need it?
Go deliberately has **no inheritance** (no "class extends class"). But real programs still benefit from reusing behavior across related types. Embedding gives Go a form of code reuse through **composition** — "has-a" relationships that behave conveniently like "is-a" for field/method access — without the complexity and pitfalls of classical inheritance.

### 3. What problem does it solve?
```text
Without embedding:
type Base struct { ID int; CreatedAt time.Time }
type User struct {
    Base   Base   // named field — must write user.Base.ID every time
    Name   string
}

With embedding:
type User struct {
    Base            // embedded — anonymous field
    Name   string
}
u.ID          // works directly! promoted from Base
u.CreatedAt   // also promoted
```

### 4. How does it work?
```text
type Base struct {
    ID int
}

func (b Base) Describe() string {
    return fmt.Sprintf("ID: %d", b.ID)
}

type User struct {
    Base       // embedded, no field name — just the type
    Name string
}

u := User{Base: Base{ID: 1}, Name: "Aman"}
u.ID          // promoted field access
u.Describe()  // promoted method access
```
The embedded type's fields and methods are **promoted** to the outer type, so they can be accessed as if they belonged to it directly.

### 5. Simple Mental Model
```text
Embedding = "baking one struct into another" so the outer struct
automatically inherits access to the inner one's fields and methods —
composition that FEELS like inheritance for field/method access,
but is NOT true inheritance.
```

### 6. Basic Go Example
```go
package main

import "fmt"

type Animal struct {
	Name string
}

func (a Animal) Describe() string {
	return "I am " + a.Name
}

type Dog struct {
	Animal   // embedded
	Breed string
}

func main() {
	d := Dog{Animal: Animal{Name: "Rex"}, Breed: "Labrador"}
	fmt.Println(d.Describe()) // "I am Rex" — promoted method
	fmt.Println(d.Name)        // "Rex" — promoted field
}
```

### 7. Explain the Code
```text
1. Dog embeds Animal (no field name given — just the type Animal).
2. d.Name works directly, even though Name actually lives on Animal —
   Go "promotes" it automatically.
3. d.Describe() calls Animal's method, again via promotion.
4. Importantly: Dog is NOT an Animal — this is composition, not
   subtyping. You cannot pass a Dog where an Animal is explicitly required,
   the way you could with true inheritance/subtyping in other languages.
```

### 8. Real-Life Problem
```text
Backend example: sharing common fields/behavior across multiple
domain models, like audit fields on every database entity.

type AuditFields struct {
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (a AuditFields) Age() time.Duration {
    return time.Since(a.CreatedAt)
}

type Order struct {
    AuditFields
    ID     int64
    Total  float64
}

type User struct {
    AuditFields
    ID    int64
    Email string
}
```
Both `Order` and `User` get `.Age()` and the audit timestamp fields "for free," without duplicating the fields or the method.

### 9. When should I use it?
When multiple types genuinely share common fields/behavior that makes sense as a reusable building block — audit fields, common metadata, or wrapping/extending a type from another package (e.g., embedding `sync.Mutex` to add locking to your own type).

### 10. When should I NOT use it?
Don't reach for embedding just to save typing a field name, or to simulate class-style inheritance hierarchies — Go's philosophy favors small, explicit composition over deep, implicit hierarchies. If the relationship isn't a genuine "built from this reusable part," a plain named field is often clearer.

### 11. Common Mistakes
- Believing embedding is inheritance — a `Dog` cannot be passed where an `Animal` is expected as a distinct type; only the fields/methods are promoted, not the type identity.
- Ambiguous promotion: if two embedded types both have a field/method with the same name, accessing it directly at the outer level is a compile error — you must disambiguate explicitly (`d.Animal.Name` vs `d.OtherType.Name`).
- Overusing embedding to flatten deep hierarchies of shared structs, making it hard to trace where a field or method actually comes from.

### 12. Important Gotchas
- **Promotion is shallow by name, not deep by "is-a" semantics** — this is a real composition mechanism, not polymorphism. A function expecting an `Animal` parameter will NOT accept a `Dog`, even though `Dog` embeds `Animal`.
- If the embedded type has a pointer-receiver method, embedding an addressable value still gets the method promoted appropriately (following the same method-set rules from 6.2) — but embedding by value vs by pointer (`Animal` vs `*Animal`) changes some of these details subtly.
- Ambiguous field/method names from multiple embedded types must be resolved explicitly, or Go will refuse to compile the ambiguous access.

### 13. Internals
```text
Go Language Guarantee:
- An embedded type's exported fields/methods are promoted to the
  outer type for direct access, following the standard field/method
  lookup and ambiguity rules.

Implementation Detail:
- Internally, the embedded value is just a regular (anonymous) field;
  promotion is a compile-time name-resolution feature, not a
  different runtime representation.
```

### 14. Standard Library Connection
```text
bufio.ReadWriter embeds a *Reader and a *Writer, gaining both
sets of promoted methods.
Many HTTP framework "context" or "request" types embed *http.Request
to add extra convenience methods on top.
```

### 15. Production Example
```go
type BaseRepository struct {
	db *sql.DB
}

func (r *BaseRepository) Ping() error {
	return r.db.Ping()
}

type UserRepository struct {
	BaseRepository // embedded — gets Ping() for free
}

func (r *UserRepository) FindByID(id int64) (*User, error) {
	// uses r.db, inherited via the embedded BaseRepository
	...
}
```

### 16. Performance
Embedding has no runtime performance cost beyond what a regular field of that type would have — promotion is resolved at compile time, not through any dynamic dispatch mechanism.

### 17. Related Concepts
| Concept | Meaning |
|---|---|
| Embedding | Including a type anonymously inside another struct |
| Promotion | Automatic access to an embedded type's fields/methods |
| Interface embedding (Ch. 7) | The same idea, applied to combining interfaces |
| Composition vs Inheritance | Embedding is "has-a" reuse, not "is-a" subtyping |

### 18. Interview Questions

**Basic**
- Q: What is struct embedding? A: Including a type anonymously (no field name) inside another struct, so its fields/methods are promoted to the outer type.
- Q: Does Go have inheritance? A: No — Go uses composition via embedding instead.

**Intermediate**
- Q: Can a `Dog` that embeds `Animal` be passed where an `Animal` parameter is expected? A: No — embedding is composition, not subtyping; the types remain distinct.

**Advanced**
- Q: What happens if two embedded types have a field with the same name? A: Accessing it directly at the outer level becomes ambiguous and is a compile error, requiring explicit disambiguation like `d.TypeA.Field`.

**Tricky**
- Q: Why does Go prefer embedding over classical inheritance? A: Go's designers favored explicit, flat composition over deep implicit hierarchies, which tend to create fragile coupling and the "diamond problem" seen in some inheritance-based languages — embedding gives code reuse benefits while keeping type relationships simple and explicit.

### 19. Interview Follow-Up Questions
```text
Q: What is struct embedding?
Q: How is it different from inheritance?
Q: What happens with promoted field/method name collisions?
Q: When would you choose embedding over a named field?
Q: How does this concept extend to interfaces? (bridges into Ch. 7)
```

### 20. Interview Answer
> "Go doesn't have classical inheritance — instead, struct embedding lets me include one type anonymously inside another, and its fields and methods get promoted to the outer type automatically. It's composition, not subtyping — a type that embeds Animal isn't itself an Animal as far as the type system is concerned; it just gets convenient access to Animal's fields and methods. I use it for genuinely shared building blocks, like audit fields or a base repository with a shared DB connection, rather than trying to build deep inheritance-like hierarchies, which isn't idiomatic Go."

### 21. Quick Revision
```text
WHAT?      → Anonymous field embedding a type inside another struct
WHY?       → Code reuse via composition, since Go has no inheritance
PROBLEM?   → Avoids duplicating shared fields/methods across types
HOW?       → Embedded type's fields/methods are "promoted" to the outer type
REAL USE?  → Shared audit fields, base repository with common DB logic
GOTCHA?    → NOT inheritance — Dog isn't substitutable for Animal
INTERVIEW? → Emphasize composition vs inheritance explicitly
```

### 22. Code Challenge
> Create a `Logger` struct with a `Log(msg string)` method. Embed it into a `Service` struct, and call `service.Log(...)` directly via promotion. Then add a second embedded type also with a `Log` method and observe the resulting ambiguity error.

---

## 6.4 Method Values and Expressions

### 1. What is it?
```text
A method VALUE binds a specific receiver to a method, producing
a function value you can call without repeating the receiver.

A method EXPRESSION produces a function value where the receiver
becomes an explicit first parameter, instead of being bound.
```

### 2. Why do we need it?
Sometimes you want to pass "this specific object's method" around as a callback, without repeating the receiver every time — that's a method value. Other times you want a generic function that takes ANY receiver of that type as its first argument — that's a method expression. Both let you treat methods as flexible function values (building on 5.5).

### 3. What problem does it solve?
```text
Without method values:
callback := func() { p.Distance(origin) }  // manually wrap the receiver

With a method value:
callback := p.Distance   // p is captured automatically; call as callback(origin)
```

### 4. How does it work?
```text
Method value:
f := p.Distance     // p is bound; f has type func(Point) float64
f(q)                 // equivalent to p.Distance(q)

Method expression:
f := Point.Distance  // receiver becomes an explicit parameter
f(p, q)               // equivalent to p.Distance(q)
```

### 5. Simple Mental Model
```text
Method value       = "remember WHO (which specific value), just tell me WHAT to do to it later."
Method expression  = "don't remember who yet — I'll tell you both WHO and WHAT each time."
```

### 6. Basic Go Example
```go
package main

import "fmt"

type Point struct{ X, Y float64 }

func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

func main() {
	p := Point{1, 2}
	q := Point{3, 4}

	// Method value: receiver p is already bound
	addToP := p.Add
	fmt.Println(addToP(q)) // {4 6}

	// Method expression: receiver becomes an explicit first argument
	add := Point.Add
	fmt.Println(add(p, q)) // {4 6}
}
```

### 7. Explain the Code
```text
1. addToP := p.Add creates a method VALUE — p is captured right now,
   so addToP only needs the remaining argument (q) when called later.
2. add := Point.Add creates a method EXPRESSION — nothing is bound yet,
   so BOTH the receiver (p) and the argument (q) must be supplied at call time.
3. Both ultimately produce the same result, but they're useful in
   different situations, as shown next.
```

### 8. Real-Life Problem
```text
Backend example: registering a specific object's method as an
event handler/callback, without wrapping it manually.

type OrderService struct{ ... }

func (s *OrderService) HandleOrderPlaced(order Order) {
    ...
}

eventBus.Subscribe("order.placed", service.HandleOrderPlaced) // method value
```
`service.HandleOrderPlaced` is a method value: it already "remembers" which `service` instance to call it on, so the event bus can treat it as a plain callback function.

### 9. When should I use it?
- **Method values**: when passing "this specific object's behavior" as a callback (event handlers, goroutine bodies, functional-style pipelines).
- **Method expressions**: rarer in everyday code — useful for generic helper functions that need to apply a method across different receiver values, or in some reflection/generic-programming scenarios.

### 10. When should I NOT use it?
Don't reach for method expressions just to look clever if a simple direct method call would be clearer — they're a fairly advanced/rare tool in typical backend code.

### 11. Common Mistakes
- Confusing a method value with just calling the method immediately — `p.Add` (no parens) is a function value; `p.Add(q)` (with parens) is an actual call.
- Forgetting a method value captures the receiver **at the time the method value is created**, not at call time — if the receiver is a value type, later changes to the original variable won't affect the already-captured method value.
- Using a method expression when a method value would be simpler and clearer for the situation.

### 12. Important Gotchas
- If the receiver is a **value type**, `p.Add` copies `p`'s value at the moment the method value is created — mutating `p` afterward does not affect the already-created method value's behavior.
- If the receiver is a **pointer type**, `p.Method` (where `p` is addressable) captures the pointer, so later mutations to the underlying value ARE visible through the method value, since it's still pointing at the same data.
- Method expressions require you to pass the receiver explicitly as the first argument every time — this is intentional flexibility, not a shortcut.

### 13. Internals
```text
Go Language Guarantee:
- p.Method (method value) has type func(remaining params) returnType,
  with the receiver already bound.
- Type.Method (method expression) has type func(ReceiverType, remaining params) returnType,
  with the receiver as an explicit leading parameter.

Implementation Detail:
- The compiler generates the appropriate closure (for method values)
  or plain function (for method expressions) under the hood —
  this is a compile-time transformation, not a different runtime mechanism.
```

### 14. Standard Library Connection
```text
Method values are commonly used when passing a bound method as a
callback to functions expecting a func(...) type, such as event
handlers, HTTP handler registration in some frameworks, or goroutine bodies.
```

### 15. Production Example
```go
type MetricsCollector struct {
	name string
}

func (m *MetricsCollector) Record(value float64) {
	fmt.Printf("[%s] recorded: %.2f\n", m.name, value)
}

func processValues(values []float64, record func(float64)) {
	for _, v := range values {
		record(v)
	}
}

// usage:
collector := &MetricsCollector{name: "latency"}
processValues([]float64{1.2, 3.4, 5.6}, collector.Record) // method value
```

### 16. Performance
Method values involve a small amount of closure-like allocation to bind the receiver, similar to anonymous function closures (5.6) — negligible for typical backend usage, but worth knowing if used in extremely hot loops.

### 17. Related Concepts
| Concept | Meaning |
|---|---|
| Method value | `p.Method` — receiver bound, produces a callable function value |
| Method expression | `Type.Method` — receiver becomes an explicit parameter |
| Function value (5.5) | The broader concept that method values/expressions build on |

### 18. Interview Questions

**Basic**
- Q: What is a method value? A: A function value created from `receiver.Method`, with the receiver already bound.
- Q: What is a method expression? A: A function value created from `Type.Method`, where the receiver becomes an explicit first parameter.

**Intermediate**
- Q: If you create a method value from a value-type receiver, does mutating the original variable afterward affect the method value? A: No — the receiver was copied at the time the method value was created.

**Advanced**
- Q: What's the practical use case for a method expression over a method value? A: When you want a generic function that can apply the same method logic across MULTIPLE different receiver values, supplying the receiver explicitly each call, rather than being bound to one specific instance.

**Tricky**
- Q: If the receiver is a pointer type, does a method value see later mutations to the underlying data? A: Yes — because the pointer itself (not a copy of the data) is captured, so the method value still points at the same, potentially-mutated data.

### 19. Interview Follow-Up Questions
```text
Q: What's the difference between a method value and a method expression?
Q: Does a method value capture the receiver by value or reference?
Q: When would a value-receiver method value NOT reflect later changes?
Q: When would a pointer-receiver method value reflect later changes?
Q: Where are method values commonly used in real code?
```

### 20. Interview Answer
> "A method value is what you get when you take `instance.Method` without calling it — it's a function value with the receiver already bound, which is really convenient for passing an object's specific behavior as a callback, like an event handler. A method expression, `Type.Method`, is less common — it leaves the receiver as an explicit first parameter instead of binding it, which is useful for generic helper functions that need to apply the same method across different instances. One subtlety I keep in mind: if the receiver is a value type, the method value captures a copy at creation time, but if it's a pointer, later mutations are still visible through it."

### 21. Quick Revision
```text
WHAT?      → Method value: bound receiver + method as a function value
           → Method expression: receiver as an explicit first parameter
WHY?       → Pass an object's behavior around like any other function value
PROBLEM?   → Avoids manually wrapping receiver.Method() in a closure
HOW?       → p.Method (value) vs Type.Method (expression)
REAL USE?  → Registering collector.Record as a callback function
GOTCHA?    → Value receiver captured by copy; pointer receiver captured live
INTERVIEW? → Be ready to explain the value-vs-pointer capture difference
```

### 22. Code Challenge
> Create a `Counter` struct with an `Increment()` pointer-receiver method. Store `counter.Increment` as a method value, call it three times through the stored function value, and confirm the original counter's count actually increased.

---

## 6.5 Example: Bit Vector Type

### 1. What is it?
```text
A worked example: building a simple "IntSet" (a set of non-negative
integers) using a bit vector — a []uint64 slice where each individual
BIT represents whether one specific number is present in the set.
```

### 2. Why do we need it?
This is a classic Go teaching example (from "The Go Programming Language" book) that ties together everything from this chapter: struct definition, pointer-receiver methods for mutation, and method design — while showing a genuinely useful, memory-efficient data structure.

### 3. What problem does it solve?
```text
Without a bit vector:
A map[int]bool or []bool uses far more memory per element than
necessary, since each entry takes at least a byte (often more,
due to overhead), even though a "present/absent" flag only needs 1 bit.

With a bit vector:
Each uint64 word holds 64 individual "present/absent" flags —
64x more memory-efficient per element than a naive []bool, in principle.
```

### 4. How does it work?
```text
IntSet stores numbers using their bit POSITION.

To find if number `x` is in the set:
    word  = x / 64   (which uint64 "word" holds this number)
    bit   = x % 64    (which bit position within that word)

    present if: (words[word] >> bit) & 1 == 1
```
```text
Example: storing {0, 2, 64}

words[0] = 0b...00000101   (bits 0 and 2 are set)
words[1] = 0b...00000001   (bit 0 of word 1 = number 64)
```

### 5. Simple Mental Model
```text
IntSet = a row of light switches (bits), grouped into panels
(uint64 words) of 64 switches each. Flipping switch N to "on"
means "number N is in the set."
```

### 6. Basic Go Example
```go
package main

import "fmt"

type IntSet struct {
	words []uint64
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func main() {
	var set IntSet
	set.Add(1)
	set.Add(9)
	set.Add(144)
	fmt.Println(set.Has(9))   // true
	fmt.Println(set.Has(123)) // false
}
```

### 7. Explain the Code
```text
1. IntSet holds a []uint64 — a growable slice of "panels of 64 bits."
2. Has(x): finds which word and which bit within that word represent x,
   and checks if that specific bit is 1. `word < len(s.words)` guards
   against reading past the end of the slice for large x not yet stored.
3. Add(x): grows the words slice (with append) until it's long enough
   to hold x's word, then sets the correct bit using `|=` (bitwise OR-assign).
4. Both methods use a POINTER receiver — Has() doesn't strictly need
   one (it doesn't mutate), but Add() absolutely does, since it must
   grow the slice, and by convention, once one method needs a pointer
   receiver, all methods on the type typically use pointer receivers too (6.2).
```

### 8. Real-Life Problem
```text
Backend example: tracking "which user IDs have already been processed
today" in a batch job, when IDs are dense and numerous — a bit vector
is far more memory-efficient than a map[int64]bool for large, dense ID ranges.

Real systems (e.g., Redis's bitmap/bitset commands, Bloom filters,
feature flag systems tracking millions of user IDs) use exactly this
bit-vector idea for compact membership tracking.
```

### 9. When should I use it?
When you need to track membership (present/absent) for a large, dense range of non-negative integers, and memory efficiency matters — feature flags per user ID, deduplication of numeric IDs, visited-node tracking in graph algorithms.

### 10. When should I NOT use it?
If the numbers are sparse (a few large, widely-scattered values) or not naturally non-negative small-ish integers, a `map[int]bool` or `map[int]struct{}` is simpler and won't waste memory on huge unused ranges of bits.

### 11. Common Mistakes
- Forgetting to grow the underlying slice before setting a bit in `Add` — this would panic with an index-out-of-range error.
- Using `Has` on an index that would require growing the slice — a naive implementation might panic instead of correctly returning `false` for "not present" (this example correctly guards with `word < len(s.words)`).
- Forgetting that `1 << bit` needs `bit` to be an appropriately-sized unsigned type to avoid subtle type issues in bit-shifting.

### 12. Important Gotchas
- This is a great real interview example for demonstrating **why pointer receivers matter**: `Add` MUST use a pointer receiver, because it needs to grow (reassign) the slice — a value receiver would only grow a local copy, silently doing nothing useful to the caller's set.
- Bitwise operators (`&`, `|`, `<<`, `>>`) used here are a common area of quiet unfamiliarity — worth being comfortable reading and writing them for interviews.
- This structure only supports **non-negative** integers as written — negative numbers would require a different indexing scheme.

### 13. Internals
```text
Go Language Guarantee:
- Bitwise operators (&, |, <<, >>) behave per Go's spec on unsigned integers.
- append grows the slice as shown in Chapter 4.2.

Implementation Detail:
- The exact growth strategy of the underlying []uint64 slice (how much
  extra capacity append reserves) is a runtime implementation detail,
  not something this design depends on for correctness.
```

### 14. Standard Library Connection
```text
math/bits provides bit-manipulation helper functions (like
bits.OnesCount64) that pair naturally with a hand-rolled bit-vector
type like this one, for operations like counting set bits.
```

### 15. Production Example
```go
type ProcessedIDs struct {
	words []uint64
}

func (p *ProcessedIDs) MarkProcessed(id int) {
	word, bit := id/64, uint(id%64)
	for word >= len(p.words) {
		p.words = append(p.words, 0)
	}
	p.words[word] |= 1 << bit
}

func (p *ProcessedIDs) IsProcessed(id int) bool {
	word, bit := id/64, uint(id%64)
	return word < len(p.words) && p.words[word]&(1<<bit) != 0
}
```
A batch job tracking which of millions of sequential record IDs have already been handled, using a fraction of the memory a `map[int]bool` would need.

### 16. Performance
```text
Memory: 1 bit per number vs. much more per entry in a map or []bool
         — dramatically more memory-efficient for dense integer ranges.
Speed:  Has/Add are both O(1) (ignoring the occasional slice growth in Add),
         very fast compared to map hashing overhead.
```

### 17. Related Concepts
| Concept | Meaning |
|---|---|
| Bit vector / bitset | Compact structure using individual bits for membership |
| Pointer receiver (6.2) | Required here since Add must mutate/grow the slice |
| Bitwise operators | `&` AND, `\|` OR, `<<`/`>>` shift — the mechanism behind bit manipulation |

### 18. Interview Questions

**Basic**
- Q: What is a bit vector used for here? A: Efficiently tracking which non-negative integers are present in a set, using individual bits instead of bytes or map entries.
- Q: Why is `[]uint64` used instead of `[]bool`? A: Each `uint64` packs 64 individual flags into 8 bytes, far more memory-efficient than a bool per entry.

**Intermediate**
- Q: Why must `Add` use a pointer receiver? A: Because it needs to grow the underlying slice, which requires modifying the actual struct field, not a local copy.

**Advanced**
- Q: How do you compute which word and bit represent a given integer x? A: `word = x / 64`, `bit = x % 64` — dividing the integer's position by the word size (64 bits) and taking the remainder.

**Tricky**
- Q: If `Has` were implemented without the `word < len(s.words)` bounds check, what would happen when checking a large, never-added number? A: It would panic with an index-out-of-range error instead of correctly returning `false`, since the word slot for that number was never allocated.

### 19. Interview Follow-Up Questions
```text
Q: How does this bit-vector set work internally?
Q: Why is a pointer receiver required for Add but not strictly for Has?
Q: How would you extend this to support Remove or Union operations?
Q: What are the memory trade-offs vs a map[int]bool?
Q: When would a bit vector be a bad fit for a real dataset?
```

### 20. Interview Answer
> "This example implements a set of non-negative integers using a bit vector — a slice of uint64 'words,' where each individual bit represents whether one number is present. To check membership, I compute which word and bit position a number maps to with division and modulo by 64, then test that bit. Adding a number may need to grow the underlying slice, which is exactly why the Add method must use a pointer receiver — a value receiver would only grow a throwaway copy. It's a great real-world illustration of when pointer receivers are not optional, and it's dramatically more memory-efficient than a map or bool slice for dense integer ranges."

### 21. Quick Revision
```text
WHAT?      → Set of ints stored as individual bits inside []uint64
WHY?       → Massive memory savings vs map[int]bool for dense ranges
PROBLEM?   → 1 bit per number instead of a byte+ per map/bool entry
HOW?       → word = x/64, bit = x%64; test/set via &, |, <<
REAL USE?  → Tracking processed IDs, feature flags, dedup at scale
GOTCHA?    → Add MUST use a pointer receiver to grow the slice
INTERVIEW? → Walk through word/bit math and the pointer-receiver reasoning
```

### 22. Code Challenge
> Add a `Remove(x int)` method (clear the bit) and a `Len() int` method (count how many bits are set, using `math/bits.OnesCount64` on each word) to the `IntSet` type above.

---

## 6.6 Encapsulation

### 1. What is it?
```text
Encapsulation means hiding a type's internal details and only
exposing a controlled, intentional public interface — in Go,
this is done using CAPITALIZATION, not private/public keywords.
```

### 2. Why do we need it?
If every field of every type is freely accessible and modifiable from anywhere, it's very easy for other code to put a type into an invalid or inconsistent state, and any internal change to the type risks breaking unrelated code that depended on those details. Encapsulation lets a type control exactly how it can be used and changed.

### 3. What problem does it solve?
```text
Without encapsulation:
type Account struct {
    Balance float64 // anyone, anywhere, can set this to any invalid value
}
acc.Balance = -500 // nothing stops this nonsensical state

With encapsulation:
type Account struct {
    balance float64 // unexported — only this package can touch it directly
}
func (a *Account) Withdraw(amount float64) error {
    if amount > a.balance {
        return errors.New("insufficient funds")
    }
    a.balance -= amount
    return nil
}
```

### 4. How does it work?
```text
Go's rule is simple and purely based on the FIRST LETTER of an
identifier's name:

Capitalized (e.g., Balance, Withdraw) → exported, visible outside the package
lowercase   (e.g., balance, validate)  → unexported, visible only within the same package
```
There's no `private`/`public`/`protected` keyword — visibility is determined entirely by capitalization, applied consistently to types, fields, functions, and methods.

### 5. Simple Mental Model
```text
Capital letter = "the outside world can see and use this."
Lowercase letter = "this is an internal detail, package-only."
```

### 6. Basic Go Example
```go
package bank

type Account struct {
	Owner   string // exported: anyone can read/set this
	balance float64 // unexported: only code inside package "bank" can touch it
}

func NewAccount(owner string) *Account {
	return &Account{Owner: owner}
}

func (a *Account) Deposit(amount float64) {
	a.balance += amount
}

func (a *Account) Balance() float64 {
	return a.balance // controlled, read-only access from outside the package
}
```

### 7. Explain the Code
```text
1. Owner is exported (capital O) — freely accessible from other packages.
2. balance is unexported (lowercase b) — code outside package "bank"
   cannot access account.balance directly at all; it's a compile error.
3. Deposit() and Balance() are exported METHODS that provide
   controlled access to the unexported field — this is the
   standard Go "getter" pattern (note: Go convention omits "Get" —
   it's Balance(), not GetBalance()).
```

### 8. Real-Life Problem
```text
Backend example: preventing invalid state in a domain type used
across many parts of a large codebase.

type Order struct {
    items []LineItem // unexported — must go through AddItem
}

func (o *Order) AddItem(item LineItem) error {
    if item.Quantity <= 0 {
        return errors.New("quantity must be positive")
    }
    o.items = append(o.items, item)
    return nil
}
```
By keeping `items` unexported, `Order` guarantees every item that ever gets added has passed validation — no other code, anywhere in the program, can bypass `AddItem` and push an invalid item directly into the slice.

### 9. When should I use it?
For any field where you want to enforce invariants (valid states), control how/when it changes, or simply hide implementation details that might change later without wanting to break every caller across the codebase.

### 10. When should I NOT use it?
For simple, passive data-holder types (DTOs, configuration structs) where there's no invariant to protect and direct field access is genuinely simpler and clearer, exporting fields directly is fine and idiomatic — don't over-encapsulate trivial structs.

### 11. Common Mistakes
- Exporting a field "just in case," when it should really be controlled through a method to protect an invariant.
- Writing a Go "getter" as `GetBalance()` — idiomatic Go omits the `Get` prefix; it's simply `Balance()`.
- Forgetting that encapsulation in Go is at the **package** level, not the **type** level — any code within the SAME package can access unexported fields/methods of any type in that package, not just the type's own methods.

### 12. Important Gotchas
- Go's encapsulation boundary is the **package**, not the type/class. Unlike languages with `private` (visible only inside the class), Go's lowercase visibility means "visible anywhere in this package" — a different, coarser granularity worth understanding clearly.
- JSON marshaling (Chapter 4.5) can ONLY see **exported** fields — an unexported field will simply be silently skipped during `json.Marshal`/`Unmarshal`, which is sometimes a source of confusing bugs for beginners.
- There is no way to make something visible to "friend" packages selectively, the way some languages support — Go's visibility model is intentionally simple: exported (all packages) or unexported (this package only).

### 13. Internals
```text
Go Language Guarantee:
- Identifiers starting with an uppercase letter are exported
  (accessible from other packages); lowercase are unexported
  (accessible only within the declaring package).

Implementation Detail:
- This is enforced entirely at compile time by the Go compiler
  examining identifier names — there's no separate runtime
  access-control mechanism involved.
```

### 14. Standard Library Connection
```text
Many standard library types encapsulate internal state this way —
e.g., time.Time has unexported internal fields, exposing only
methods like .Year(), .Month(), .Add() for controlled access.
```

### 15. Production Example
```go
type RateLimiter struct {
	mu       sync.Mutex // unexported — internal implementation detail
	requests map[string]int // unexported
	limit    int             // unexported
}

func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{requests: make(map[string]int), limit: limit}
}

func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests[key]++
	return r.requests[key] <= r.limit
}
```
Callers of `RateLimiter` only ever see `NewRateLimiter` and `Allow` — the mutex, map, and limit are entirely hidden implementation details that can be changed later without breaking any caller.

### 16. Performance
Encapsulation itself has zero runtime performance cost in Go — visibility rules are purely a compile-time, source-level concept, not something enforced or checked at runtime.

### 17. Related Concepts
| Concept | Meaning |
|---|---|
| Exported identifier | Capitalized — accessible from other packages |
| Unexported identifier | Lowercase — accessible only within the declaring package |
| Getter convention | `Balance()`, not `GetBalance()`, is idiomatic Go |
| Package-level visibility | Go's granularity is the package, not the individual type |

### 18. Interview Questions

**Basic**
- Q: How does Go control visibility/encapsulation? A: By capitalization — capitalized identifiers are exported (public), lowercase are unexported (package-private).
- Q: Does Go have `private`/`public` keywords? A: No — visibility is purely based on identifier capitalization.

**Intermediate**
- Q: What is the idiomatic Go getter method name for a field called `balance`? A: `Balance()`, without a "Get" prefix.

**Advanced**
- Q: Is Go's encapsulation boundary the type or the package? A: The package — any code within the same package can access unexported fields/methods of any type declared in that package, not just the type's own methods.

**Tricky**
- Q: Why might a struct with unexported fields silently produce unexpected JSON output? A: `encoding/json` can only see exported fields via reflection, so unexported fields are simply omitted from `Marshal` output (and ignored during `Unmarshal`) without any error, which can confuse developers expecting them to appear.

### 19. Interview Follow-Up Questions
```text
Q: How does Go implement encapsulation without private/public keywords?
Q: What's the idiomatic naming convention for a getter?
Q: Is the encapsulation boundary the type or the package?
Q: How does encapsulation interact with JSON marshaling?
Q: Why might you keep a field unexported even without strict validation needs?
```

### 20. Interview Answer
> "Go doesn't have private or public keywords — visibility is determined purely by capitalization. A capitalized field, function, or type is exported and visible from other packages; lowercase means it's unexported, visible only within the same package. I use this to protect invariants — for example, keeping a balance field unexported and only allowing changes through a Withdraw method that validates the amount first. One thing to remember is that Go's encapsulation boundary is the package, not the individual type, so any code in the same package can still reach unexported fields directly."

### 21. Quick Revision
```text
WHAT?      → Capitalization controls exported vs unexported visibility
WHY?       → Protects invariants and hides internal implementation details
PROBLEM?   → Prevents external code from putting a type into an invalid state
HOW?       → Capital letter = exported; lowercase = package-private
REAL USE?  → Account.balance unexported, only changed via Withdraw/Deposit
GOTCHA?    → Boundary is the PACKAGE, not the type; unexported fields skip JSON
INTERVIEW? → Know the Balance() vs GetBalance() convention
```

### 22. Code Challenge
> Refactor the `IntSet` type from 6.5 so its `words` field is unexported (it already is), then add an exported `Len()` method and an exported `String() string` method that formats the set's contents nicely, without ever exposing the raw `words` slice directly.

---

# End of Chapter 6 — Methods

## Quick Chapter Summary
```text
Method Declarations   → func (recv Type) Name(...) — attaches behavior to a type
Pointer Receivers     → needed to mutate the original, or avoid copying large structs
Struct Embedding      → composition-based code reuse; NOT inheritance
Method Values/Exprs   → treating methods as function values, bound or unbound
Bit Vector Example    → a full worked example tying receivers + methods together
Encapsulation         → capitalization controls exported vs unexported visibility
```

## How These Connect
```text
Method Declaration (attach behavior to a type)
   ↓
Pointer Receiver (decide: copy, or operate on the real thing?)
   ↓
Struct Embedding (reuse fields/methods across related types)
   ↓
Method Values (treat a bound method as a portable function value)
   ↓
Encapsulation (control what's actually exposed to the outside world)
```

## Final Revision

### Most Important Concepts
- The **method set rule**: `T`'s method set has only value-receiver methods; `*T`'s method set has both — this decides interface satisfaction (bridges directly into Chapter 7).
- Embedding is **composition**, never inheritance — no implicit subtyping.
- Go's encapsulation granularity is the **package**, not the type.

### Must Remember
```text
Method = function + receiver
Value receiver  = works on a copy
Pointer receiver = works on the real thing
Embedding = "baked-in" reuse, fields/methods promoted, but NOT is-a
Capital letter = exported; lowercase = package-private
```

### Common Traps
- Writing a mutating method with a value receiver — it silently does nothing to the caller's data.
- Believing an embedded type creates a true "is-a" relationship.
- Trying to call a pointer-receiver method on a non-addressable value (like a map value).
- Writing `GetBalance()` instead of the idiomatic `Balance()`.
- Assuming unexported fields are hidden from other types in the SAME package too (they aren't).

### Top Interview Questions
- What's the method set rule, and why does it matter for interfaces?
- Why must a slice-growing method use a pointer receiver?
- How is Go's embedding different from classical inheritance?
- What's the difference between a method value and a method expression?
- How does Go implement encapsulation without access-modifier keywords?

### Advanced Questions
- Walk through exactly why `myMap[key].Increment()` fails to compile when `Increment` has a pointer receiver.
- Design a small library of composable structs using embedding, and explain the ambiguous-promotion edge cases that could arise.
- Explain how the `IntSet` bit-vector example would need to change to support negative integers.
- Why does Go's package-level (not type-level) encapsulation change how you structure large codebases compared to class-based languages?

### One-Minute Chapter Explanation
> "Go attaches behavior to types through methods, which are just functions with a receiver parameter. The biggest decision with any method is whether to use a value or pointer receiver — pointer receivers are required to mutate the original data or to avoid copying large structs, and this choice also determines which method set a type has, which directly affects interface satisfaction later. Since Go has no inheritance, code reuse between related types happens through struct embedding, which promotes an embedded type's fields and methods to the outer type — but this is composition, not subtyping, so an embedding type is never substitutable for the type it embeds. Methods can also be treated as portable function values, either bound to a specific receiver or left generic as method expressions. And for hiding implementation details, Go skips private/public keywords entirely in favor of a simple rule: capitalized names are exported, lowercase names are package-private — letting a type protect its invariants while still keeping the whole language refreshingly simple."

---

*Series so far: Chapter 4 (Composite Types) in `README.md`, Chapter 5 (Functions) in `chapter5-functions.md`, and this file for Chapter 6 (Methods). Say "continue" if you want to move on to Chapter 7 (Interfaces) as the next installment.*