# Go (Golang) Deep Study Guide — Backend + Interview Prep

> Simple English. Deep technical understanding. Real backend examples. Interview-ready.
> Built chapter by chapter, following the structure you uploaded.

---

# 4. Composite Types

Composite types are types built from other types. Instead of holding just one value (like an `int`), they hold many values together — a list, a table, a group of fields.

Covered here: Arrays, Slices, Maps, Structs, JSON, Text and HTML Templates.

---

## 4.1 Arrays

### 1. What is it?
```text
An array is a fixed-size list of values, all of the same type.

"Fixed-size" means: once you decide the size, it can never change.
```

### 2. Why do we need it?
Sometimes you know exactly how many items you will ever need — for example, a chess board always has 8x8 squares. An array gives you a block of memory of exactly that size, laid out next to each other, so access is fast and predictable.

### 3. What problem does it solve?
```text
Without a fixed-size structure:

You would use separate variables: x1, x2, x3, x4...
This does not scale and cannot be looped over.

With an array:
var scores [4]int
You can loop, index, and pass it as one unit.
```

### 4. How does it work?
```text
[10]int  →  a block of 10 ints, side by side in memory
             index:  0   1   2   3   4   5   6   7   8   9
```
The size is part of the type. `[4]int` and `[5]int` are **different types** — you cannot assign one to the other, and you cannot pass a `[4]int` where a `[5]int` is expected.

### 5. Simple Mental Model
```text
Array = a fixed box with a fixed number of slots.
The number of slots is written on the box, and it never changes.
```

### 6. Basic Go Example
```go
package main

import "fmt"

func main() {
	var nums [5]int      // array of 5 ints, all start at 0
	nums[0] = 10
	nums[1] = 20
	fmt.Println(nums)     // [10 20 0 0 0]
	fmt.Println(len(nums)) // 5
}
```

### 7. Explain the Code
```text
1. var nums [5]int creates an array of 5 ints. Go fills it with zero values (0) automatically.
2. nums[0] = 10 sets the first slot.
3. len(nums) always returns 5 — it is fixed.
```

### 8. Real-Life Problem
```text
Arrays are rare in everyday Go backend code because their size is fixed
and known at compile time. A real backend rarely knows in advance
"I will have exactly 5 users" — user counts change.

That is exactly why Go gives us Slices (next topic), which are
built on top of arrays but can grow.
```

### 9. When should I use it?
- When the size is truly fixed and known ahead of time (e.g. a SHA-256 hash is always `[32]byte`, an RGB color is `[3]byte`).
- When you want value semantics (copying the whole array by value) on purpose.

### 10. When should I NOT use it?
Do not use arrays for general-purpose collections in backend code — almost always you want a **slice** instead, because request counts, database rows, and list sizes are not known in advance.

### 11. Common Mistakes
- Trying to use an array like a dynamic list (append does not work on arrays the way it does on slices).
- Forgetting that arrays are **copied by value** when assigned or passed to functions — this can silently waste memory and cause confusing bugs (you edit a copy, not the original).
- Comparing `[3]int{}` with `[4]int{}` and expecting it to work — it's a compile error, they are different types.

### 12. Important Gotchas
- **Value semantics**: `b := a` copies the entire array. Changing `b` does NOT change `a`.
- Passing an array to a function copies it. If the array is large, this is expensive.
- Arrays ARE comparable with `==` if their element type is comparable — unlike slices, which are not comparable.

### 13. Internals
```text
Array value

+----+----+----+----+----+
| 10 | 20 |  0 |  0 |  0 |
+----+----+----+----+----+
```
### Go Language Guarantee
The array's length is fixed and is part of its type. Elements sit contiguously in memory.

### Implementation Detail
The Go compiler may keep small arrays on the stack instead of the heap for performance — this is an optimization, not a guarantee you should rely on in your reasoning about correctness.

### 14. Standard Library Connection
```text
[32]byte  → used by crypto packages like sha256.Sum256
[4]byte   → sometimes used for IPv4 addresses
```

### 15. Production Example
```go
func hashPassword(data []byte) [32]byte {
	return sha256.Sum256(data) // fixed-size 32-byte array result
}
```
This is used because a SHA-256 hash is *always* exactly 32 bytes — a perfect fit for a fixed-size array.

### 16. Performance
Copying a large array is expensive (it copies every byte). For big collections, always prefer a slice, which copies only a small header, not the whole data.

### 17. Related Concepts
| Concept | Meaning |
|---|---|
| Array | Fixed-size, value semantics |
| Slice | Dynamic-size, reference semantics (points to underlying array) |

### 18. Interview Questions

**Basic**
- Q: What is an array in Go? A: A fixed-size sequence of elements of the same type, where the size is part of the type itself.
- Q: Can you resize an array? A: No. Size is fixed at declaration/compile time.

**Intermediate**
- Q: What happens when you assign one array to another? A: The entire array is copied (value semantics), not referenced.
- Q: Are `[3]int` and `[5]int` the same type? A: No, different array lengths mean different types.

**Advanced**
- Q: Are arrays comparable? A: Yes, with `==`, as long as the element type is comparable. Slices are not comparable this way.

**Tricky**
- Q: If you pass an array to a function and modify it inside, does the caller see the change? A: No — the function received a copy, so the original is untouched (unless you pass a pointer to the array).

### 19. Interview Follow-Up Questions
```text
Q: What is an array?
Q: Why is size part of the type?
Q: What happens on assignment?
Q: How is this different from a slice?
Q: When would you actually choose an array over a slice in production?
```

### 20. Interview Answer
> "In Go, an array is a fixed-size collection where the length is baked into the type itself. Arrays use value semantics, so copying or passing one copies all its data. In practice, I rarely use raw arrays in backend code — I reach for slices, which are dynamic. I do use arrays when the size is truly fixed by the problem, like a 32-byte hash output."

### 21. Quick Revision
```text
WHAT?      → Fixed-size list, same type, size is part of the type
WHY?       → Predictable memory layout for known, fixed sizes
PROBLEM?   → Avoids separate variables for a fixed group of values
HOW?       → Contiguous memory block, indexed 0..n-1
REAL USE?  → sha256.Sum256 returns [32]byte
GOTCHA?    → Copied by value on assignment/function calls
INTERVIEW? → Arrays are rare in real backend code — slices are the default
```

### 22. Code Challenge
> Declare an array of 7 `int` values representing daily step counts for a week. Write a function that takes the array and returns the total steps. Try it once passing by value, and once passing a pointer to the array — observe that modifying inside the pointer version affects the caller.

---

## 4.2 Slices

### 1. What is it?
```text
A slice is a flexible, resizable view into an array.

Think of it as: "a window that looks at part (or all) of an array,
and can grow."
```

### 2. Why do we need it?
Real backend programs almost never know sizes in advance — number of users, number of rows returned from a database, number of items in a cart. We need a list that can grow and shrink. Arrays can't do that; slices can.

### 3. What problem does it solve?
```text
Without slices:
You'd need to know array size ahead of time — impossible for dynamic data
(e.g., "how many rows will this SQL query return?").

With slices:
users := []User{}
users = append(users, newUser) // grows as needed
```

### 4. How does it work?
A slice is a small struct with 3 fields: a pointer to an underlying array, a length, and a capacity.
```text
Slice header
+---------+--------+----------+
| pointer | length | capacity |
+---------+--------+----------+
      |
      v
 underlying array: [ 10, 20, 30, 40, _, _ ]
                     ^length=3   ^capacity=6
```

### 5. Simple Mental Model
```text
Slice = a "view" (pointer + length + capacity) into an array.
Multiple slices can share the SAME underlying array.
```

### 6. Basic Go Example
```go
package main

import "fmt"

func main() {
	nums := []int{1, 2, 3}
	nums = append(nums, 4)
	fmt.Println(nums, len(nums), cap(nums))
}
```

### 7. Explain the Code
```text
1. []int{1,2,3} creates a slice with length 3.
2. append adds 4. If capacity allows, Go reuses the same array.
   If not, Go allocates a new, bigger array and copies everything over.
3. len() tells current count, cap() tells how much room exists before
   the next grow-and-copy happens.
```

### 8. Real-Life Problem
```text
Backend example: reading rows from a database.

rows, _ := db.Query("SELECT * FROM users")
var users []User
for rows.Next() {
    var u User
    rows.Scan(&u.ID, &u.Name)
    users = append(users, u)   // grows dynamically, row by row
}
```
You never know row count ahead of time — slices make this natural.

### 9. When should I use it?
Basically always, for lists in Go. It is the default collection type for ordered data.

### 10. When should I NOT use it?
- When you need fast key lookup — use a **map** instead.
- When the exact fixed size matters for the type system itself (rare) — use an array.

### 11. Common Mistakes
- Believing slices are always independent copies — they often **share** the same underlying array, so mutating one can silently mutate another.
- Appending to a slice and assuming the original variable is always updated everywhere — `append` may or may not reallocate, so always reassign the result: `s = append(s, x)`.
- Forgetting `nil` slices are usable (`len` is 0, `append` works fine on `nil`).

### 12. Important Gotchas
```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3]     // b shares memory with a
b[0] = 999      // this changes a[1] too!
```
- **Slicing shares memory** — a classic Go interview trap.
- `append` beyond capacity allocates a **new** array; below capacity, it mutates the existing one in place. This makes bugs appear only sometimes, depending on capacity — hard to catch in testing.
- A `nil` slice and an empty slice (`[]int{}`) are different in identity but behave the same for `len`, `range`, and `append`.

### 13. Internals
```text
Slice header (this part IS copied when you pass a slice around)
+---------+--------+----------+
| pointer |   len  |   cap    |
+---------+--------+----------+

Underlying array (this part is SHARED, not copied)
```
### Go Language Guarantee
A slice has a length and capacity, and indexing beyond length panics.

### Implementation Detail
Go's growth strategy (e.g., roughly doubling capacity when small) is an implementation detail of the runtime and can change between Go versions — don't hardcode assumptions about exact growth factors.

### 14. Standard Library Connection
```text
[]byte   → used everywhere: io.Reader, io.Writer, bytes.Buffer
[]string → used in strings.Split, flag parsing, etc.
sort.Slice → sorts any slice using a custom comparison function
```

### 15. Production Example
```go
type UserRepository struct{}

func (r *UserRepository) FindActive() ([]User, error) {
	var users []User
	rows, err := r.db.Query("SELECT id, name FROM users WHERE active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
```

### 16. Performance
- `append` in a loop without a known final size can cause multiple reallocations. If you know the approximate size ahead of time, pre-allocate: `make([]User, 0, expectedCount)` avoids repeated copying.
- Passing a slice to a function is cheap (only the 3-field header is copied), even if the underlying data is huge.

### 17. Related Concepts
| Concept | Meaning |
|---|---|
| `len` | Number of elements currently in use |
| `cap` | Room available before a new array must be allocated |
| `make([]T, len, cap)` | Creates a slice with explicit length and capacity |
| Slicing `a[i:j]` | Creates a new slice header sharing the same array |

### 18. Interview Questions

**Basic**
- Q: What is a slice? A: A resizable view over an underlying array, made of a pointer, length, and capacity.
- Q: How do you add an element to a slice? A: `append(slice, element)`, and reassign the result.

**Intermediate**
- Q: What happens when `append` exceeds capacity? A: Go allocates a new, larger array, copies the old data, and returns a slice pointing to the new array.
- Q: Do two slices from the same array share memory? A: Yes, until one of them triggers a reallocation via `append`.

**Advanced**
- Q: Why might mutating a sub-slice unexpectedly affect the original slice? A: Because slicing doesn't copy the underlying array — both slices point to the same memory until a reallocation happens.
- Q: How do you safely copy a slice? A: Use `copy(dst, src)` into a slice with enough length, or `append([]T{}, src...)`.

**Tricky**
- Q: If you `append` to a slice inside a function, does the caller's slice variable change? A: Only if you return the new slice and the caller reassigns it — Go passes the slice header by value, so the caller's variable itself is untouched unless reassigned. The underlying array *might* be mutated in place, but the caller's `len`/`cap` view won't update automatically.

### 19. Interview Follow-Up Questions
```text
Q: What is a slice made of internally?
Q: What's the difference between length and capacity?
Q: What happens on append when capacity runs out?
Q: Why do two slices sometimes share memory unexpectedly?
Q: How do you avoid this shared-memory bug in production code?
```

### 20. Interview Answer
> "A slice in Go is a lightweight view over an array — it holds a pointer, a length, and a capacity. This makes slices cheap to pass around, but it also means slices can share underlying memory, which is a common source of bugs if you're not careful. `append` grows the slice, reallocating a new array only when the current capacity is exceeded. In production, I preallocate capacity when I know the approximate size, to avoid repeated copying."

### 21. Quick Revision
```text
WHAT?      → Resizable view (pointer+len+cap) over an array
WHY?       → Dynamic-size lists, since real data sizes aren't known ahead
PROBLEM?   → Arrays can't grow; slices can
HOW?       → append grows in place if capacity allows, else reallocates
REAL USE?  → Reading DB rows into a growing list
GOTCHA?    → Slices can share memory — mutating one can mutate another
INTERVIEW? → Always mention pointer+len+cap and the reallocation rule
```

### 22. Code Challenge
> Write a function `removeAt(s []int, index int) []int` that removes an element at a given index without leaving gaps. Then write a test proving that slicing shares memory by modifying a sub-slice and checking the parent slice.

---

## 4.3 Maps

### 1. What is it?
```text
A map is a collection of key-value pairs.

Like a dictionary: you look up a word (key) to get its meaning (value).
```

### 2. Why do we need it?
We often need fast lookup by some identifier — find a user by ID, find a config value by name, count word frequency. Scanning a slice for a match is slow (checks every item). A map gives near-instant lookup.

### 3. What problem does it solve?
```text
Without a map:
for _, u := range users {
    if u.ID == targetID { ... }   // O(n) — slow for large lists
}

With a map:
usersByID[targetID]               // O(1) average — fast
```

### 4. How does it work?
Internally, Go maps are hash tables: the key is hashed to decide where to store the value, so lookup, insert, and delete are all fast on average.

### 5. Simple Mental Model
```text
Map = a set of labeled boxes.
The label (key) is hashed to find the right box instantly,
instead of checking every box one by one.
```

### 6. Basic Go Example
```go
package main

import "fmt"

func main() {
	ages := map[string]int{"Aman": 22, "Riya": 25}
	ages["Karan"] = 30            // add
	fmt.Println(ages["Aman"])      // 22
	age, ok := ages["Missing"]     // 0, false — safe lookup
	fmt.Println(age, ok)
}
```

### 7. Explain the Code
```text
1. map[string]int means: keys are strings, values are ints.
2. ages["Karan"] = 30 inserts a new key-value pair.
3. The "comma ok" pattern (age, ok := ages["Missing"]) tells you
   whether the key actually existed, avoiding confusion with a
   real zero value.
```

### 8. Real-Life Problem
```text
Backend example: caching user data by ID to avoid repeated DB hits.

cache := map[int64]*User{}

func getUser(id int64) *User {
    if u, ok := cache[id]; ok {
        return u   // fast path, no DB call
    }
    u := fetchFromDB(id)
    cache[id] = u
    return u
}
```

### 9. When should I use it?
When you need fast lookup, existence checks, counting, grouping, or deduplication by a key.

### 10. When should I NOT use it?
- When order matters — maps in Go have **no guaranteed iteration order**.
- When the collection is small and you'll only ever scan it once — a slice may be simpler and cache-friendlier.
- For concurrent read/write from multiple goroutines without protection — plain maps are **not** safe for concurrent writes.

### 11. Common Mistakes
- Assuming map iteration order is consistent — it is intentionally randomized by Go.
- Writing to a map from multiple goroutines without a `sync.Mutex` or `sync.Map` — this crashes with "concurrent map writes" at runtime.
- Reading a missing key without the "comma ok" check and confusing the zero value with "key exists but is zero."

### 12. Important Gotchas
- Maps are **reference types** — passing a map to a function does not copy the data; both refer to the same underlying map.
- A `nil` map can be read from (returns zero value) but writing to a `nil` map **panics**. Always `make(map[K]V)` before writing.
- Map keys must be comparable types (no slices or maps as keys — only types supporting `==`, like strings, ints, structs of comparable fields).

### 13. Internals
```text
Go Language Guarantee:
- Lookup, insert, delete by key.
- No guaranteed order during iteration (range).

Implementation Detail:
- Uses a hash table with buckets internally.
- Exact bucket/resizing strategy can change between Go versions.
```

### 14. Standard Library Connection
```text
map[string]interface{} → common pattern for dynamic JSON data
encoding/json uses maps when decoding into interface{}
```

### 15. Production Example
```go
type RateLimiter struct {
	mu       sync.Mutex
	requests map[string]int // key: client IP, value: request count
}

func (r *RateLimiter) Allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests[ip]++
	return r.requests[ip] <= 100
}
```
This shows the real production pattern: map + mutex together, because concurrent backend requests hit this map from many goroutines at once.

### 16. Performance
- Average O(1) lookup/insert/delete, but worst-case can degrade with poor hash distribution (rare with Go's built-in hashing).
- Maps have more memory overhead per entry than slices — don't use a map when a slice would do.

### 17. Related Concepts
| Concept | Meaning |
|---|---|
| Map | Key → value, unordered, hash-based |
| Slice | Ordered list, index-based |
| `sync.Map` | Concurrency-safe map for specific high-read patterns |

### 18. Interview Questions

**Basic**
- Q: What is a map in Go? A: An unordered collection of key-value pairs backed by a hash table.
- Q: How do you check if a key exists? A: `value, ok := m[key]`; `ok` is `false` if missing.

**Intermediate**
- Q: Is Go map iteration order guaranteed? A: No, it's intentionally randomized each run.
- Q: What happens if you write to a `nil` map? A: It panics at runtime.

**Advanced**
- Q: Are maps safe for concurrent use? A: No — concurrent writes (or a write with any read) from multiple goroutines cause a runtime panic unless protected by a mutex or you use `sync.Map`.
- Q: Can a slice be a map key? A: No, slices aren't comparable, so they can't be used as map keys (arrays can, though).

**Tricky**
- Q: If you pass a map to a function and add a key inside, does the caller see it? A: Yes — maps are reference types, so both caller and function share the same underlying map.

### 19. Interview Follow-Up Questions
```text
Q: What is a map?
Q: How is it implemented internally?
Q: Is order guaranteed?
Q: What happens with concurrent access?
Q: How would you make map access safe in a multi-goroutine backend?
```

### 20. Interview Answer
> "A Go map is a hash table giving average O(1) lookup, insert, and delete by key. It's a reference type, so passing it around doesn't copy the data — everyone shares the same map. The two biggest gotchas I watch for are: iteration order is never guaranteed, and concurrent writes from multiple goroutines will panic unless I protect the map with a mutex or use `sync.Map`."

### 21. Quick Revision
```text
WHAT?      → Key-value collection, hash table based
WHY?       → Fast lookup by key instead of scanning a list
PROBLEM?   → Avoids O(n) search for "find by ID" type operations
HOW?       → Hashing the key decides where the value is stored
REAL USE?  → In-memory cache, rate limiter, counting/grouping
GOTCHA?    → No order guarantee; concurrent writes panic; nil map write panics
INTERVIEW? → Always mention concurrency safety — it's the #1 follow-up
```

### 22. Code Challenge
> Write a function that counts word frequency in a sentence using `map[string]int`. Then wrap it with a mutex and simulate two goroutines updating the same map safely.

---

## 4.4 Structs

### 1. What is it?
```text
A struct groups related fields together into one custom type.

Example: a "User" isn't just one value — it's a Name, an Email,
an Age, all bundled together.
```

### 2. Why do we need it?
Real-world data is naturally made of multiple related fields. Without structs, you'd manage separate variables or parallel slices for every entity, which is messy and error-prone.

### 3. What problem does it solve?
```text
Without struct:
name := "Aman"
age := 22
email := "a@example.com"
// no relationship between these three variables

With struct:
type User struct {
    Name  string
    Age   int
    Email string
}
u := User{Name: "Aman", Age: 22, Email: "a@example.com"}
// one value, clearly grouped
```

### 4. How does it work?
```text
Struct value

+--------+-----+------------------+
|  Name  | Age |      Email       |
+--------+-----+------------------+
| "Aman" | 22  | "a@example.com"  |
+--------+-----+------------------+
```
Fields are laid out together in memory. A struct is a value type — copying it copies all fields.

### 5. Simple Mental Model
```text
Struct = a labeled bundle of fields, treated as one single value.
```

### 6. Basic Go Example
```go
package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func main() {
	u := User{Name: "Aman", Age: 22}
	fmt.Println(u.Name, u.Age)
}
```

### 7. Explain the Code
```text
1. type User struct {...} defines a new type with two fields.
2. User{Name: "Aman", Age: 22} creates a value using named fields
   (recommended — order doesn't matter, and it's self-documenting).
3. u.Name accesses a field with dot notation.
```

### 8. Real-Life Problem
```text
Backend example: representing a database row as a Go value.

type Order struct {
    ID       int64
    UserID   int64
    Total    float64
    Status   string
}
```
Every layer of a backend — DB scanning, JSON responses, business logic — passes `Order` structs around as one coherent unit instead of loose variables.

### 9. When should I use it?
Whenever you're modeling a real-world entity with multiple related fields: users, orders, config, API requests/responses.

### 10. When should I NOT use it?
Don't create a struct for a single, standalone value that has no related fields — a plain variable or a named type (`type UserID int64`) is simpler.

### 11. Common Mistakes
- Forgetting struct assignment copies all fields (value semantics) — mutating a copy doesn't affect the original unless you use a pointer.
- Comparing structs with `==` when they contain uncomparable fields (like slices or maps) — this is a compile error.
- Exporting all fields (capitalized) when some should be private (lowercase) implementation details.

### 12. Important Gotchas
- Structs are compared field-by-field with `==`, but **only if every field is itself comparable**.
- Passing a large struct by value to a function copies the whole thing — for large structs, pass a pointer (`*User`) instead.
- Struct field order affects memory layout and padding (advanced/performance topic) — reordering fields can reduce memory usage.

### 13. Internals
```text
Go Language Guarantee:
- Fields are accessed via dot notation.
- A struct is a value type: assignment and function calls copy it.

Implementation Detail:
- The Go compiler may add "padding" bytes between fields
  for memory alignment — this affects struct size but not behavior.
```

### 14. Standard Library Connection
```text
http.Request  → a struct bundling Method, URL, Header, Body, etc.
time.Time     → internally a struct
```

### 15. Production Example
```go
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CreateUser(req CreateUserRequest) (*User, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	// ... save to DB
	return &User{Name: req.Name, Email: req.Email}, nil
}
```

### 16. Performance
- Large structs passed by value can be expensive to copy repeatedly — prefer pointers (`*User`) once a struct grows beyond a few fields, or is passed frequently in hot paths.
- Struct field ordering can reduce padding/memory footprint (place larger fields before smaller ones) — a micro-optimization, not usually necessary early on.

### 17. Related Concepts
| Concept | Meaning |
|---|---|
| Struct | Value type bundling named fields |
| Pointer to struct | Reference to the same struct, avoids copying |
| Embedding (6.3) | Reusing one struct's fields/methods inside another |

### 18. Interview Questions

**Basic**
- Q: What is a struct? A: A composite type that groups named fields together.
- Q: How do you access a struct field? A: Dot notation, e.g. `u.Name`.

**Intermediate**
- Q: Are structs value types or reference types? A: Value types — assigning or passing them copies all fields.
- Q: When should you use a pointer to a struct instead of the struct itself? A: When the struct is large, or when you need the function to mutate the original.

**Advanced**
- Q: When are two struct values comparable with `==`? A: When every field is itself a comparable type (no slices, maps, or functions as fields).

**Tricky**
- Q: If a struct field is a slice, and you copy the struct, does the slice also get deep-copied? A: No — the slice header is copied, but it still points to the same underlying array, so both struct copies share the slice's data.

### 19. Interview Follow-Up Questions
```text
Q: What is a struct?
Q: Is it a value type or reference type?
Q: What happens to slice/map fields when a struct is copied?
Q: When do you pass a struct by pointer instead of by value?
Q: How does struct embedding relate to structs? (bridges into 6.3)
```

### 20. Interview Answer
> "A struct in Go groups related fields into a single named type — it's how we model real entities like a User or an Order. Structs are value types, so copying one copies every field. In production, I pass small structs by value and larger ones (or ones I need to mutate) by pointer, to avoid unnecessary copying."

### 21. Quick Revision
```text
WHAT?      → Named bundle of fields forming one type
WHY?       → Models real entities cleanly instead of loose variables
PROBLEM?   → Avoids disconnected variables for related data
HOW?       → Fields laid out together; value type; copied on assignment
REAL USE?  → http.Request, DB row models, API request/response types
GOTCHA?    → Slice/map fields aren't deep-copied when struct is copied
INTERVIEW? → Know when to use pointer receivers/params for structs
```

### 22. Code Challenge
> Define a `Product` struct with `Name`, `Price`, and `Tags []string`. Copy a `Product` value, modify the `Tags` slice on the copy, and observe that the original is also affected. Explain why.

---

## 4.5 JSON

### 1. What is it?
```text
JSON (JavaScript Object Notation) is a text format for representing
structured data — used everywhere in backend APIs to send/receive data.

Go's encoding/json package converts between Go structs and JSON text.
```

### 2. Why do we need it?
Backends talk to frontends, mobile apps, and other services over HTTP, and JSON is the universal language for that conversation. Go needs a reliable way to turn its typed structs into JSON, and JSON back into structs.

### 3. What problem does it solve?
```text
Without JSON support:
You'd manually parse text — extremely error-prone and slow to write.

With encoding/json:
json.Marshal(myStruct)     // Go struct → JSON bytes
json.Unmarshal(data, &obj) // JSON bytes → Go struct
```

### 4. How does it work?
Go uses **reflection** to look at your struct's fields and their `json:"..."` tags at runtime, and maps them to/from JSON keys.
```text
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

Marshal:    User{"Aman", 22}  →  {"name":"Aman","age":22}
Unmarshal:  {"name":"Aman","age":22}  →  User{"Aman", 22}
```

### 5. Simple Mental Model
```text
Struct tags are labels telling encoding/json:
"When you convert this field, call it THIS name in JSON."
```

### 6. Basic Go Example
```go
package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	u := User{Name: "Aman", Age: 22}

	data, _ := json.Marshal(u)
	fmt.Println(string(data)) // {"name":"Aman","age":22}

	var u2 User
	json.Unmarshal(data, &u2)
	fmt.Println(u2) // {Aman 22}
}
```

### 7. Explain the Code
```text
1. json.Marshal(u) turns the struct into JSON bytes.
2. json:"name" tag controls the JSON key name (without it, Go would
   use the Go field name "Name" as-is).
3. json.Unmarshal(data, &u2) fills u2 by reading the JSON —
   note the & because Unmarshal must modify u2 directly.
```

### 8. Real-Life Problem
```text
Backend HTTP handler example:

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }
    // req.Name, req.Email now populated from client's JSON body
}
```
Every REST API in Go leans on this — reading a JSON request body into a struct, and writing a struct back as a JSON response.

### 9. When should I use it?
For any HTTP API request/response body, config files in JSON format, or storing structured data as JSON in a database column.

### 10. When should I NOT use it?
- For very high-performance/low-latency internal services, consider a faster binary format (protobuf) instead of JSON.
- Don't use `map[string]interface{}` when you know the shape of the data — define a proper struct instead, for type safety.

### 11. Common Mistakes
- Forgetting the `&` when calling `Unmarshal` (it needs a pointer to modify the target).
- Not handling the `error` returned by `Marshal`/`Unmarshal`.
- Using unexported (lowercase) struct fields and wondering why they never appear in JSON — `encoding/json` can only see exported (capitalized) fields.

### 12. Important Gotchas
- Only **exported fields** (starting with a capital letter) are marshaled/unmarshaled — private fields are silently skipped.
- `omitempty` in a tag (`json:"age,omitempty"`) skips the field in output if it's the zero value — but this means a real `0` and "not set" become indistinguishable in the output.
- JSON numbers unmarshal into `interface{}` as `float64` by default — a classic Go interview trap when working with `map[string]interface{}`.

### 13. Internals
```text
Go Language Guarantee:
- Marshal/Unmarshal use struct tags and exported fields as documented.

Implementation Detail:
- Uses reflection internally, which has some CPU cost compared to
  hand-written serialization code. This is generally acceptable for
  typical backend workloads.
```

### 14. Standard Library Connection
```text
encoding/json — Marshal, Unmarshal, NewEncoder, NewDecoder
Used directly in almost every Go HTTP backend.
```

### 15. Production Example
```go
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
```

### 16. Performance
- Reflection-based marshaling has overhead compared to generated/hand-written code — for very high-throughput services, some teams use code generation (e.g. `easyjson`) to skip reflection.
- Use `json.NewDecoder(r.Body).Decode(&v)` (streaming) instead of reading the whole body into memory first, for large payloads.

### 17. Related Concepts
| Concept | Meaning |
|---|---|
| `Marshal` | Go value → JSON bytes |
| `Unmarshal` | JSON bytes → Go value |
| Struct tags | Control field naming/behavior during (un)marshaling |
| `Encoder`/`Decoder` | Stream JSON to/from `io.Writer`/`io.Reader` |

### 18. Interview Questions

**Basic**
- Q: How do you convert a Go struct to JSON? A: `json.Marshal(value)`.
- Q: How do you convert JSON into a Go struct? A: `json.Unmarshal(data, &value)`.

**Intermediate**
- Q: Why must you pass a pointer to `Unmarshal`? A: Because it needs to modify the target value directly; passing by value would only modify a copy.
- Q: What does `omitempty` do? A: Omits the field from JSON output if it holds its zero value.

**Advanced**
- Q: Why do unexported struct fields not appear in JSON? A: `encoding/json` uses reflection, which can only access exported (capitalized) fields from outside the package.
- Q: What type does a JSON number become when unmarshaled into `interface{}`? A: `float64`, even if it looks like an integer in the JSON text.

**Tricky**
- Q: If two struct fields end up with the same JSON tag name, what happens? A: It's ambiguous and can lead to one field being silently ignored during unmarshaling — Go's json package has tie-breaking rules based on field depth, but it's a bug-prone situation to avoid.

### 19. Interview Follow-Up Questions
```text
Q: How does Marshal/Unmarshal work?
Q: Why do struct tags matter?
Q: Why is a pointer needed for Unmarshal?
Q: What happens with unexported fields?
Q: How would you handle a JSON field that could be a number or a string?
```

### 20. Interview Answer
> "Go's encoding/json package converts between structs and JSON using reflection and struct tags. Marshal serializes a Go value to JSON bytes, Unmarshal does the reverse into a pointer target. In my APIs, I define explicit request/response structs with json tags rather than using generic maps, so I get compile-time type safety, and I always check the returned error since malformed JSON is common with real clients."

### 21. Quick Revision
```text
WHAT?      → Convert Go structs <-> JSON text
WHY?       → JSON is the standard format for API request/response bodies
PROBLEM?   → Manual text parsing is error-prone; encoding/json automates it
HOW?       → Reflection reads struct tags to map fields to JSON keys
REAL USE?  → HTTP handlers decoding request bodies, encoding responses
GOTCHA?    → Only exported fields work; numbers become float64 in interface{}
INTERVIEW? → Know the pointer requirement for Unmarshal cold
```

### 22. Code Challenge
> Define a `Product` struct with `json` tags including `omitempty` on an optional `Discount float64` field. Marshal a `Product` with `Discount: 0` and observe it's omitted. Then unmarshal a JSON string with an extra unknown field and confirm Go ignores it silently by default.

---

## 4.6 Text and HTML Templates

### 1. What is it?
```text
Templates let you generate text (or HTML) by mixing a fixed template
string with dynamic data — like a "mail merge" for Go.

text/template → for any plain text output (emails, config files, CLI output)
html/template → same idea, but auto-escapes data to prevent HTML/JS injection
```

### 2. Why do we need it?
Backends often need to generate dynamic content: a welcome email with the user's name, an HTML page with data filled in, a config file generated from variables. Hand-building this with string concatenation is messy and, for HTML, dangerous (XSS risk).

### 3. What problem does it solve?
```text
Without templates:
"Hello " + user.Name + ", your balance is " + strconv.Itoa(balance)
— messy, and if used for HTML, a user's Name like <script> could
  inject malicious code into the page (XSS vulnerability).

With html/template:
{{.Name}} is automatically escaped, so <script> becomes harmless text.
```

### 4. How does it work?
```text
Template string:  "Hello {{.Name}}, you have {{.Count}} new messages."
Data:             struct{ Name string; Count int }{"Aman", 3}

Template engine walks the template, replaces {{.Field}} placeholders
with the matching field from the data, and (for html/template)
escapes dangerous characters based on WHERE in the HTML they appear
(inside a tag, inside a script, inside an attribute, etc).
```

### 5. Simple Mental Model
```text
Template = a fill-in-the-blanks form.
{{.Field}} = a blank that gets filled with real data.
html/template additionally makes sure nothing dangerous slips through.
```

### 6. Basic Go Example
```go
package main

import (
	"os"
	"text/template"
)

func main() {
	const tpl = "Hello {{.Name}}, you have {{.Count}} new messages.\n"
	t := template.Must(template.New("greeting").Parse(tpl))

	data := struct {
		Name  string
		Count int
	}{"Aman", 3}

	t.Execute(os.Stdout, data)
}
```

### 7. Explain the Code
```text
1. template.New("greeting") creates a named template.
2. .Parse(tpl) parses the template text, which contains {{.Name}}
   and {{.Count}} placeholders.
3. template.Must panics if parsing fails — used for templates you
   know are hardcoded and correct (fail fast during startup).
4. t.Execute(os.Stdout, data) fills the placeholders using `data`
   and writes the result to stdout (or any io.Writer — e.g. an
   http.ResponseWriter in a real server).
```

### 8. Real-Life Problem
```text
Backend example: rendering an HTML page for a web app, or
generating a "welcome" email body.

tmpl := template.Must(template.ParseFiles("welcome_email.html"))
tmpl.Execute(w, EmailData{Name: user.Name, ActivationLink: link})
```
Using html/template here (not text/template) is critical, because
`user.Name` came from user input — if it contained `<script>...</script>`,
html/template automatically neutralizes it. text/template would not.

### 9. When should I use it?
- `html/template` for any output that will be rendered as HTML in a browser — server-rendered pages, emails with HTML content.
- `text/template` for plain-text output: config files, CLI output, plain-text emails, code generation.

### 10. When should I NOT use it?
For simple, fixed strings with no dynamic data, plain string formatting (`fmt.Sprintf`) is simpler and doesn't need the template machinery.

### 11. Common Mistakes
- Using `text/template` to render HTML — this is a real **security bug**, because it does not escape user input, opening the door to XSS attacks.
- Forgetting to handle the error from `Parse`/`Execute` (using `Must` hides parse errors as panics, which is fine at startup but not for user-supplied templates).
- Mismatching field names in the template (`{{.name}}` lowercase) against exported struct fields (`Name` uppercase) — templates can only access exported fields.

### 12. Important Gotchas
- `html/template` escaping is **context-aware** — it escapes differently depending on whether the data lands inside an HTML tag, an attribute, a URL, or a `<script>` block. This is more sophisticated (and safer) than a single blanket escape function.
- Template parsing errors and execution errors are separate — a template can parse fine but fail at `Execute` time if data doesn't match what the template expects.
- Nested templates and `{{define}}`/`{{template}}` blocks let you compose larger pages from smaller partials — useful in real web apps but adds structure you must design deliberately.

### 13. Internals
```text
Go Language Guarantee:
- html/template guarantees contextual auto-escaping of untrusted data
  based on where it appears in the HTML/JS/CSS/URL context.

Implementation Detail:
- Internally, html/template wraps text/template's parsing/execution
  engine and inserts escaping logic during execution — this pipeline
  detail can evolve, but the escaping guarantee itself is the API contract.
```

### 14. Standard Library Connection
```text
text/template and html/template share the same {{ }} syntax and
core engine — html/template is literally built on top of text/template,
adding the escaping safety layer.
```

### 15. Production Example
```go
type PageData struct {
	Title string
	Items []string
}

func renderPage(w http.ResponseWriter, data PageData) {
	tmpl := template.Must(template.ParseFiles("templates/page.html"))
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}
```

### 16. Performance
- Parsing a template is relatively expensive; **parse once at startup** (e.g. in an `init` or main setup) and reuse the parsed `*Template` for every request, rather than re-parsing on every HTTP request.
- Executing an already-parsed template is fast and suitable for per-request use in a web server.

### 17. Related Concepts
| Concept | Meaning |
|---|---|
| `text/template` | Plain-text generation, no escaping |
| `html/template` | HTML generation, context-aware auto-escaping |
| `{{.Field}}` | Access a field of the data passed to `Execute` |
| `{{range}}` / `{{if}}` | Loops and conditionals inside templates |

### 18. Interview Questions

**Basic**
- Q: What's the difference between `text/template` and `html/template`? A: `html/template` automatically escapes data to prevent HTML/JS injection; `text/template` does not escape anything.
- Q: How do you insert a struct field into a template? A: `{{.FieldName}}`, and it must be exported.

**Intermediate**
- Q: Why would using `text/template` for a web page be dangerous? A: User-supplied data wouldn't be escaped, opening an XSS vulnerability.
- Q: When should you parse a template — once at startup, or per request? A: Once at startup; parsing is comparatively expensive, execution is cheap and safe to repeat per request.

**Advanced**
- Q: How does `html/template`'s escaping differ depending on context? A: It inspects where in the HTML/CSS/JS/URL structure the data lands and applies the correct escaping rules for that specific context, rather than one generic escape.

**Tricky**
- Q: If you accidentally use `text/template` to render an HTML page containing user input, what's the real-world risk? A: A malicious user could submit input like `<script>...</script>` that gets rendered unescaped — this is a classic stored/reflected XSS vulnerability, potentially letting attackers run JavaScript in other users' browsers.

### 19. Interview Follow-Up Questions
```text
Q: What is a Go template?
Q: What's the difference between text/template and html/template?
Q: Why does that difference matter for security?
Q: When should templates be parsed vs executed?
Q: How would you structure templates for a multi-page web app? (partials/layouts)
```

### 20. Interview Answer
> "Go's template packages let me generate dynamic text or HTML by filling placeholders with real data. For anything rendered in a browser, I always use html/template instead of text/template, because it automatically, contextually escapes untrusted data — protecting against XSS. In production, I parse templates once at startup and reuse the parsed template object across requests, since parsing is the expensive part and execution is cheap."

### 21. Quick Revision
```text
WHAT?      → Fill-in-the-blanks engine for generating text/HTML
WHY?       → Avoids messy, unsafe manual string building
PROBLEM?   → html/template also prevents XSS via auto-escaping
HOW?       → {{.Field}} placeholders filled from data at Execute time
REAL USE?  → Rendering server-side HTML pages, generating emails
GOTCHA?    → Never use text/template for HTML with user-supplied data
INTERVIEW? → Be ready to explain WHY html/template exists (security)
```

### 22. Code Challenge
> Write an HTML template that displays a list of usernames using `{{range .Users}}`. Pass a username containing `<b>` as one of the values and confirm `html/template` escapes it, while a `text/template` version does not.

---

# End of Chapter 4 — Composite Types

## Quick Chapter Summary
```text
Array   → fixed size, value semantics, rare in real backend code
Slice   → dynamic size, reference semantics, the default Go list type
Map     → key-value hash table, fast lookup, not concurrency-safe by default
Struct  → bundles related fields into one type, value semantics
JSON    → Marshal/Unmarshal structs <-> JSON, powers almost every Go API
Templates → text/template (plain) vs html/template (auto-escaped, safe for browsers)
```

## How These Connect
```text
Struct (models your data)
   ↓
Slice of structs (a list of records)
   ↓
Map (fast lookup of those records by key/ID)
   ↓
JSON (send that data over HTTP as an API response)
   ↓
html/template (render that same data as a web page, safely)
```

---

*Next up: Chapter 5 (Functions) and Chapter 6 (Methods), continuing in the same format. Say "continue" and I'll add them to this file.*