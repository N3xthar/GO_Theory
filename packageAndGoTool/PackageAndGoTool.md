# 10. Packages and the Go Tool — Complete Study Guide
### (Simple English • Deep Understanding • Interview Ready • Top 1%)

> Brother, don't worry. I will teach this like you are seeing Go for the very first time. No hard words without explaining them. Every big topic will follow this chain:
>
> **WHAT → WHY → PROBLEM WITHOUT IT → HOW → EXAMPLE → MISTAKE → PRODUCTION USE → INTERVIEW**

---

## Table of Contents

1. [Mental Model First](#1-mental-model-first)
2. [What is a Package?](#2-what-is-a-package)
3. [Package Declaration](#3-package-declaration)
4. [Scope (Package / File / Block)](#4-scope)
5. [Exported vs Unexported](#5-exported-vs-unexported)
6. [Import Paths](#6-import-paths)
7. [Import Declarations](#7-import-declarations)
8. [Blank Imports `_`](#8-blank-imports)
9. [init() and Package Initialization](#9-init-and-package-initialization)
10. [Package vs Module](#10-package-vs-module)
11. [go.mod](#11-gomod)
12. [go.sum](#12-gosum)
13. [Standard Library vs Third-Party](#13-standard-library-vs-third-party)
14. [Internal Packages](#14-internal-packages)
15. [Import Cycles](#15-import-cycles)
16. [The Go Tool (Commands)](#16-the-go-tool)
17. [go run vs go build](#17-go-run-vs-go-build)
18. [go install vs go get](#18-go-install-vs-go-get)
19. [go test](#19-go-test)
20. [go list, go doc, go vet](#20-go-list-go-doc-go-vet)
21. [Real Backend Project Structure](#21-real-backend-project-structure)
22. [Debugging Scenarios](#22-debugging-scenarios)
23. [Compiler-Level Mental Model](#23-compiler-level-mental-model)
24. [Anti-Patterns](#24-anti-patterns)
25. [Production Engineering](#25-production-engineering)
26. [One-Page Revision Sheet](#26-one-page-revision-sheet)
27. [Flashcards (40+)](#27-flashcards)
28. [Mock Interview (3 Rounds)](#28-mock-interview)
29. [Coding Exercises](#29-coding-exercises)
30. ["If You Remember Only This"](#30-if-you-remember-only-this)
31. [Beyond the Syllabus — Top 1%](#31-beyond-the-syllabus--top-1)
32. [Final Mental Map](#32-final-mental-map)

---

## 1. Mental Model First

Think of building a house.

```
Bricks          →  Go source file (.go)
One room        →  Package (many files, same room/purpose)
Whole house      →  Module (many rooms/packages, one go.mod)
Address of house →  Module path (github.com/you/myapp)
Address of room  →  Import path (github.com/you/myapp/handler)
Builder          →  go build (compiler + linker)
Finished house   →  Binary (the .exe / executable file)
```

**Simple chain:**

```
file → package → import → import path → module → go.mod → dependency graph → compiler → linker → binary
```

### Key words, explained like to a child

| Word | Simple meaning |
|---|---|
| **File** | One `.go` text file with code in it. |
| **Package** | A folder of `.go` files that work together as one unit. |
| **Import path** | The "address" you type to use someone else's package. |
| **Module** | A big folder that has a `go.mod` file. It can contain many packages. |
| **Module path** | The name of the module, written inside `go.mod`. |
| **Repository** | The place on the internet (like GitHub) where the code lives. |
| **Dependency** | Someone else's package that your code uses. |
| **Standard library** | Packages that come free with Go itself, like `fmt`. |
| **Third-party package** | Packages made by other people, downloaded from the internet. |
| **Executable** | A program you can run (has `package main`). |
| **Library** | A package meant to be *used* by other code, not run by itself. |

**Golden rule to remember:**
> A **module** is like a **book**. A **package** is like a **chapter**. A **file** is like a **page**.
> A book (module) has one `go.mod`. It can have many chapters (packages). Each chapter can have many pages (files).

When you open any Go repo, ask yourself these 5 questions in order:

1. Where is `go.mod`? → that tells me the **module**.
2. Which folders have `.go` files? → those are **packages**.
3. Which folder has `package main`? → that is the **entry point** (the program that runs).
4. What does `main` import? → that is the **dependency graph**.
5. Are those imports from stdlib, this module, or the internet? → that tells me **where code comes from**.

---

## 2. What is a Package?

### WHAT
A package is simply **a folder of `.go` files that all start with the same `package` line**, and Go treats them as *one unit* of code.

### WHY
Go needs a way to group related code together, and to control what other code is allowed to see. Packages do both jobs at once: **organizing code** and **hiding code**.

### PROBLEM WITHOUT IT
Imagine every function, every variable, from every file in your whole project, all mixed together with no boundaries. One giant file with 50,000 lines. You could not reuse code safely, you could not hide internal details, and name clashes (two functions with the same name) would break everything.

### HOW
- Every `.go` file must start with a `package` line — this is compulsory, not optional.
- All files that declare the **same package name** and sit in the **same folder** become one package.
- A package is compiled as **one single unit**. When you use one function from a package, the whole package gets compiled.
- **Rule:** All `.go` files in ONE folder must belong to the SAME package (with one small exception: test files can be `xxx_test`, explained later, and this is still tied to the same package).
- A folder = one package. You cannot have two different packages mixed in the same folder (except the special `_test` suffix case for external test packages).

### EXAMPLE

```
myapp/
├── main.go       →  package main
├── user.go       →  package main
├── database.go   →  package main
└── config.go     →  package main
```

All four files say `package main` at the top — so Go treats these 4 files as **one package** called `main`. You could write all this code in one file, or split across many files — Go does not care, it will glue them together at compile time.

```go
// file: user.go
package main

type User struct {
    Name string
}
```

```go
// file: main.go
package main

import "fmt"

func main() {
    u := User{Name: "Aman"}   // no import needed! same package = direct access
    fmt.Println(u.Name)
}
```

Notice: `main.go` uses `User` from `user.go` **without importing anything** — because they are the SAME package. Import is only needed when crossing package boundaries.

### UNDER THE HOOD
When you run `go build`, Go's tool first scans the folder, collects every `.go` file with a matching `package` name, and compiles them together as one translation unit. The compiler produces one "package object" (like `.a` file internally) for each package.

### MISTAKES
- ❌ Putting two different `package` names in the same folder (Go will refuse to build: "found packages main and user").
- ❌ Thinking each file is separately compiled/imported. It's not — the whole package is the unit.
- ❌ Forgetting the `package` line completely (compile error).

### PRODUCTION
Real backend repos split one package across many files by **responsibility**: `handler.go`, `service.go`, `repository.go` might each be their **own package** in their **own folder** (not the same package!) — but within `internal/handler/`, you might see `user_handler.go`, `order_handler.go`, `middleware.go` all as `package handler`.

### INTERVIEW
**Q: What is a package in Go?**
**A:** A package is a folder of Go source files that share the same `package` declaration. Go compiles them as one unit, and the package name controls what identifiers are visible to other packages.

---

## 3. Package Declaration

### WHAT
The very first non-comment line in every `.go` file: `package <name>`.

### WHY
Go needs to know, immediately, which "bucket" this file belongs to, before it even looks at imports or code.

### RULES (naming)
- Lowercase, short, no underscores, no camelCase. Examples: `fmt`, `http`, `json`, `user`, `auth`.
- The package name is **just an identifier** — it does NOT have to match the folder name (though by strong convention, it usually does).
- `package main` is special: it tells Go "this package produces an **executable**, not a library." It must contain a `func main()`.

### `package main` vs `package database`

```go
package main   // this folder builds into a RUNNABLE PROGRAM
```
```go
package database // this folder builds into a LIBRARY (cannot run directly)
```

If a folder says `package database`, running `go run .` inside it will give an error: *"go run: no main packages"* — because there's no entry point.

### Can directory name differ from package name?
**Yes**, technically. Example: folder `myapp/uuid-generator/` could contain `package uuidgen`. Go allows this, BUT it's confusing and against convention — most style guides say: **match folder name to package name** unless you have a strong reason (like `package main` files sitting inside a folder named after the binary, e.g. `cmd/server/` containing `package main`).

### Good vs Bad package names

| Good | Why good | Bad | Why bad |
|---|---|---|---|
| `user` | short, clear | `userpackage` | redundant word "package" |
| `http` | one clear idea | `utils` | too generic, becomes a dumping ground |
| `auth` | one clear idea | `common` | tells you nothing about content |
| `json` | matches purpose | `helper` | vague, invites everything to be dumped in |

**Golden rule:** a good package name describes *what it provides*, not *what it contains*. `strings` provides string functions. `helper` doesn't tell you anything.

### INTERVIEW
**Q: Can two files in the same folder have different package names?**
**A:** No — Go will refuse to compile. Every `.go` file in a directory (except special `_test` external test files) must declare the exact same package name.

---

## 4. Scope

Three levels, from small to big:

```
Block scope   → inside { } — smallest, e.g. inside an if/for
      ↓
File scope    → only imports are file-scoped (each file has its own import list)
      ↓
Package scope → anything declared outside a function is visible to EVERY file in the same package
```

### Example

```go
// file1.go
package main

var counter = 0     // PACKAGE scope — visible in file2.go too, no import needed

func increment() {
    counter++        // block/function scope changes happen here
}
```

```go
// file2.go
package main

import "fmt"

func show() {
    fmt.Println(counter)  // works! same package, package-level scope
}
```

But **imports are per-file**. If `file1.go` imports `"fmt"`, `file2.go` still needs its OWN `import "fmt"` line if it wants to use `fmt` — imports don't "leak" across files, only declarations (variables, functions, types) do.

---

## 5. Exported vs Unexported

### WHAT
Go decides "public" or "private" using **the first letter's capitalization**. No `public`/`private` keywords like Java.

```go
package user

type User struct {
    Name string   // Exported — starts with capital N — visible outside package
    age  int      // unexported — starts with lowercase a — hidden outside package
}

func NewUser() User { ... }   // Exported — callable from other packages
func validate() bool { ... }  // unexported — only usable inside package "user"
```

### WHY
This is Go's whole answer to "encapsulation" (hiding internal details). Instead of a keyword, capitalization IS the access rule. It's baked into the language grammar itself — you cannot forget to write `public`, and you cannot misuse it, because it's structural.

### PROBLEM WITHOUT IT
Without any visibility control, every internal helper function, every internal struct field, would be usable and breakable by any other package. Bugs would spread. You couldn't safely change internal code without breaking someone else's code that depended on your "private" details.

### From another package

```go
package main

import "myapp/user"

func main() {
    u := user.NewUser()     // OK — NewUser is exported
    fmt.Println(u.Name)     // OK — Name is exported
    fmt.Println(u.age)      // ❌ COMPILE ERROR — age is unexported, not visible here
}
```

### MISTAKES
- ❌ Exporting everything "just in case" — this destroys encapsulation and makes future changes dangerous (breaking changes for everyone using your package).
- ❌ Thinking unexported means "secure" — it's only a compile-time visibility rule, not real security.

### INTERVIEW
**Q: Why does Go use capitalization instead of `public`/`private` keywords?**
**A:** It keeps the language simpler — one rule (capital letter) replaces a whole keyword system, and it's visually obvious just by reading the name whether something is exported, without needing to check a separate keyword.

---

## 6. Import Paths

### WHAT
An import path is the "address string" you write inside `import "..."` to tell Go **which package you want to use**.

```go
import "fmt"                       // standard library
import "net/http"                  // standard library, nested
import "github.com/google/uuid"    // third-party, from GitHub
import "myapp/internal/handler"    // your own module's package
```

### The MOST IMPORTANT interview distinction

```
package name   ≠   import path   ≠   module path
```

Let's break this with a real example:

```go
module github.com/example/myapp     // <-- this is the MODULE PATH (written in go.mod)

// folder: myapp/internal/handler/user.go
package handler                     // <-- this is the PACKAGE NAME (used inside code as handler.Something)
```

To import that package from elsewhere:
```go
import "github.com/example/myapp/internal/handler"   // <-- this is the IMPORT PATH (the full address)
```

But when you USE it in code, you type the **package name**, not the whole path:
```go
handler.NewUserHandler()   // "handler" here is the package NAME, not the import path
```

**Simple analogy:**
- **Module path** = your house's full postal address.
- **Import path** = the address of one specific room inside that house (module path + folder).
- **Package name** = the nickname you call that room by, once you're inside.

### WHY three separate things?
Because they solve three DIFFERENT problems:
- **Module path** identifies WHERE the whole project lives (for downloading/versioning).
- **Import path** identifies WHICH exact package/folder you want.
- **Package name** is what you type in code — kept SHORT on purpose, so code isn't cluttered with long URLs everywhere.

### How Go resolves an import
1. Go reads the import path string, e.g. `"github.com/google/uuid"`.
2. It checks: is this the standard library? (no `.` in the first path segment → usually stdlib, e.g. `fmt`, `net/http`).
3. If not stdlib, Go checks `go.mod`: does the import path start with MY module path? → it's a **local package**, found on disk.
4. If it's an outside module path, Go checks `go.sum`/module cache (`$GOPATH/pkg/mod` or module cache) for a downloaded copy, matching the version pinned in `go.mod`.
5. If missing, `go mod tidy`/`go get`/`go build` will try to download it from the internet (via the configured proxy, usually `proxy.golang.org`).

### Local packages & relative imports
Modern Go **does not use relative imports** like `import "./utils"` — this was removed/discouraged because it broke the guarantee that an import path uniquely identifies a package regardless of who imports it. Instead:

```go
// ❌ Not used in modern Go modules
import "./utils"

// ✅ Correct — always full path from module root
import "github.com/example/myapp/internal/utils"
```

**Why relative imports are bad:** They make the meaning of an import depend on the *caller's location on disk*, which breaks reproducibility, breaks tooling, and makes packages non-portable.

### INTERVIEW
**Q: What's the difference between package name and import path?**
**A:** The import path is the full string used to locate and download the package (e.g. `"github.com/google/uuid"`). The package name is the short identifier declared inside the source code (`package uuid`) and is what you actually type in your code (`uuid.New()`).

---

## 7. Import Declarations

### All the forms

```go
import "fmt"                     // single import

import (                         // grouped import (idiomatic / preferred style)
    "fmt"
    "net/http"
)

import f "fmt"                   // aliased import — now you must call f.Println(), not fmt.Println()

import _ "some/package"          // blank import — explained deeply in section 8

import . "some/package"          // dot import — DANGEROUS, rarely used
```

### Why aliases exist
Sometimes two packages have the **same package name** but different import paths (e.g. two different "json" libraries), or the default name is too generic/clashes with a local variable. Alias fixes the name clash:

```go
import (
    stdjson "encoding/json"
    jsoniter "github.com/json-iterator/go"
)
```

### Why dot imports are discouraged
```go
import . "fmt"

func main() {
    Println("hello")   // no "fmt." prefix needed — everything gets dumped into current file's namespace
}
```
**Problem:** you lose clarity about *where* `Println` came from just by reading the code. It also risks silent name collisions. It's mainly used only in some test files (dot-importing a test-helper/matcher library), and even then, sparingly.

### MISTAKES
- ❌ Importing a package and never using it — Go **refuses to compile**: "imported and not used". This is a compile ERROR, not a warning (Go is strict on purpose, to keep code clean).
- ❌ Overusing aliases when not needed, making code harder to read.

### INTERVIEW
**Q: Why does Go treat an unused import as a compile error instead of a warning?**
**A:** Go's designers wanted to force clean code and prevent "just in case" imports from silently accumulating, which happens a lot in other languages. Being strict here keeps codebases tidy by default.

---

## 8. Blank Imports

This is a topic that confuses many beginners, so let's go slow.

### WHAT
```go
import _ "some/package"
```
The underscore `_` means: **"import this package ONLY for its side effects — I will never directly call anything from it by name."**

### WHY would you import something you never use?
Because a package can do useful work **just by being loaded** — through its `init()` function (explained next). You don't need to call any function from it; its mere PRESENCE in your program is enough.

### PROBLEM WITHOUT IT
Without blank imports, Go's "unused import = compile error" rule would make it *impossible* to import a package purely for its side effects — you'd be forced to write a fake reference just to satisfy the compiler, which is ugly and error-prone. Blank import is Go's clean, honest way of saying "yes, I really mean to import this and not use its names directly."

### The classic real example: database drivers

```go
import (
    "database/sql"
    _ "github.com/lib/pq"   // blank import — postgres driver
)

func main() {
    db, _ := sql.Open("postgres", "connection-string")
    // we never call anything like pq.Something() directly!
    // but the driver MUST be imported for sql.Open("postgres", ...) to work
}
```

### The complete flow (memorize this!)

```
main package
    ↓
blank import "github.com/lib/pq"
    ↓
Go loads that package (compiles + links it in)
    ↓
package-level init() inside lib/pq runs automatically
    ↓
inside init(), the driver calls sql.Register("postgres", &Driver{})
    ↓
this "registers" the driver into database/sql's internal map
    ↓
now sql.Open("postgres", ...) can find and use it
```

### Small custom example of "registration" pattern

```go
// file: plugin/plugin.go
package plugin

var registry = map[string]func() string{}

func Register(name string, fn func() string) {
    registry[name] = fn
}

func Run(name string) string {
    return registry[name]()
}
```

```go
// file: plugins/hello/hello.go
package hello

import "myapp/plugin"

func init() {
    plugin.Register("hello", func() string { return "Hello from plugin!" })
}
```

```go
// file: main.go
package main

import (
    "fmt"
    "myapp/plugin"
    _ "myapp/plugins/hello"   // blank import — triggers hello's init(), which registers itself
)

func main() {
    fmt.Println(plugin.Run("hello"))  // prints: Hello from plugin!
}
```

Notice: `main.go` never writes `hello.Something()`. It just imports it blankly, letting `init()` do the registration work silently.

### MISTAKES
- ❌ Forgetting the blank import entirely, then wondering why `sql.Open("postgres", ...)` fails with "unknown driver" — this is a VERY common real bug.
- ❌ Using blank imports for things that don't actually need side-effect registration — misuse hides real dependencies.

### PRODUCTION
Very common for: database drivers (`pq`, `mysql`), image format decoders (`image/png`, `image/jpeg` registering with `image.Decode`), and pprof profiling (`net/http/pprof` registers HTTP debug routes just by being imported).

### INTERVIEW
**Q: Why would someone use a blank import for a database driver?**
**A:** Because the driver needs to run its `init()` function to register itself with `database/sql`, but the calling code never needs to reference the driver package's exported names directly — `database/sql` finds it by name string ("postgres") at runtime, not by direct function calls.

---

## 9. init() and Package Initialization

### WHAT
`init()` is a special function that runs **automatically**, with no arguments, no return value, before `main()` starts.

```go
func init() {
    fmt.Println("I run before main!")
}
```

### Rules
- A package can have **multiple** `init()` functions (even in the same file!) — they run in the order they appear.
- `init()` cannot be called manually.
- Every file in a package can have its own `init()`.

### Full initialization order (memorize this chain)

```
1. Imported packages are initialized FIRST (their own vars + init(), recursively, deepest dependency first)
2. Package-level variables in THIS package are initialized (in dependency order between vars)
3. init() functions in THIS package run (in file order, top to bottom within a file)
4. Finally, main() runs — but ONLY in package main
```

### Example

```go
package main

import "fmt"

var x = compute()  // package var initialized BEFORE init()

func compute() int {
    fmt.Println("computing x")
    return 42
}

func init() {
    fmt.Println("init running, x =", x)
}

func main() {
    fmt.Println("main running")
}
```

**Output:**
```
computing x
init running, x = 42
main running
```

### When init() is good vs bad

| Good use | Bad use |
|---|---|
| Registering a driver/plugin (side-effect registration) | Doing heavy work like network calls (fails silently, hard to test) |
| Validating that required env vars exist at startup, and panicking early | Hidden business logic that surprises other developers |
| Setting up small internal lookup tables/maps | Anything that should be explicit and testable |

### MISTAKES
- ❌ Using `init()` for complex startup logic — it's invisible (nobody calls it directly), it makes testing hard (you can't skip it), and execution ORDER across packages can be surprising if you have many dependencies.
- ❌ Relying on `init()` order between DIFFERENT packages without understanding: Go initializes packages in dependency order (a package that is imported is initialized before the package that imports it) — but among *sibling* packages with no relation, order can be less obvious. Best practice: never depend on cross-package init timing beyond "imports finish first."

### INTERVIEW
**Q: When does init() run?**
**A:** After all package-level variables in that package are initialized, and after all imported packages have been fully initialized, but before `main()` runs.

---

## 10. Package vs Module

### The real-world confusion, cleared up

```
Module  = the whole project (one go.mod file at the root)
Package = ONE folder of related .go files inside that project
```

**One module → many packages.**

```
myproject/                          ← MODULE root (has go.mod)
├── go.mod                          ← module github.com/you/myproject
├── main.go                         ← package main
├── internal/
│   └── auth/                       ← package auth  (a separate package)
├── handler/                        ← package handler
├── service/                        ← package service
├── repository/                     ← package repository
└── model/                          ← package model
```

All of these packages live inside **ONE module**. They share the same `go.mod`, the same dependency versions, the same `go.sum`.

### Are repository and module the same thing?
**Not necessarily.** Usually one Git repository = one Go module (simplest, most common setup). But Go DOES support **multiple modules in one repository** (called a "multi-module repo" — less common, used by big projects like Kubernetes for special reasons). And a module does NOT have to be hosted anywhere online at all — you can have a purely local module with no remote repository.

### INTERVIEW
**Q: What's the difference between a package and a module?**
**A:** A module is a versioned collection of packages defined by a `go.mod` file — it's the unit of dependency management. A package is a folder of source files compiled together as one unit — it's the unit of code organization. One module usually contains many packages.

---

## 11. go.mod

### WHAT

```go
module github.com/example/myapp

go 1.23

require (
    github.com/google/uuid v1.6.0
    github.com/lib/pq v1.10.9
)
```

### Line by line

| Line | Meaning |
|---|---|
| `module github.com/example/myapp` | Declares the module's path/name. ALL internal import paths start with this. |
| `go 1.23` | The minimum Go language version this module needs (affects language features + toolchain behavior). |
| `require (...)` | Lists direct dependencies (and sometimes indirect ones, marked `// indirect`). |

### Why go.mod exists — the history (WHY matters here)

**Before modules (old `GOPATH` era):**
- All Go code had to live inside one global folder (`$GOPATH/src`).
- There was **no per-project versioning** — if two projects needed different versions of the same dependency, you were stuck; the whole `GOPATH` shared one copy.
- Reproducible builds were hard — "works on my machine" was common.

**After modules (`go.mod`, since Go 1.11+, default since Go 1.16):**
- Each project declares its OWN dependencies + exact versions, independent of any global folder.
- You can put your code **anywhere** on disk.
- Builds became reproducible — anyone who clones your repo gets the exact same dependency versions.

### Direct vs indirect dependencies
```go
require (
    github.com/google/uuid v1.6.0          // DIRECT — your code imports this directly
    github.com/some/lib v0.3.1 // indirect  // INDIRECT — a dependency OF your dependency
)
```

### `replace` and `exclude` and `retract`

```go
replace github.com/old/pkg => github.com/new/pkg v1.0.0   // swap one module for another (e.g. local fork, testing)
exclude github.com/broken/pkg v1.2.0                        // never select this specific broken version
retract v1.0.1                                              // module AUTHOR marks their own old version as bad
```

- `replace`: very useful for local development — point a dependency to your local folder while testing changes: `replace github.com/you/lib => ../lib`.
- `exclude`: tells Go's version selection to skip a specific bad version.
- `retract`: used by package AUTHORS (not consumers) inside THEIR OWN go.mod, to warn "please don't use this version I published, it has a bug."

### Commands

| Command | What it actually does |
|---|---|
| `go mod init <module-path>` | Creates a new `go.mod` file, sets the module path. |
| `go mod tidy` | Scans your code, adds missing `require` entries, REMOVES unused ones, updates `go.sum`. |
| `go mod download` | Downloads dependencies into the local module cache (without modifying go.mod). |
| `go mod graph` | Prints the full dependency graph (who requires what). |
| `go mod why <pkg>` | Explains WHY a certain package is in your dependency tree (which import chain pulled it in) — great for debugging. |

### MISTAKES
- ❌ Manually editing `require` versions without running `go mod tidy`/`go build` afterward — can create an inconsistent state.
- ❌ Forgetting to run `go mod tidy` before committing — leaves stale or missing entries.

### INTERVIEW
**Q: What problem did go.mod solve that GOPATH could not?**
**A:** GOPATH forced one single global version per dependency, shared across ALL projects on your machine, with no per-project isolation. `go.mod` lets each project pin its own exact dependency versions, making builds reproducible and projects self-contained, without needing a special folder structure.

---

## 12. go.sum

### WHAT
`go.sum` is a file full of **cryptographic checksums (hashes)** for every dependency (and its specific version) your module uses.

```
github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
```

### WHY it exists
`go.mod` says WHICH VERSION you want. `go.sum` proves that the code you actually downloaded is EXACTLY, byte-for-byte, the code that everyone else downloaded — nothing was tampered with.

### Difference from go.mod

| `go.mod` | `go.sum` |
|---|---|
| Says WHICH dependencies and versions you want | Says WHAT those exact downloaded files must hash to |
| You edit it (or `go mod tidy` edits it) | You basically never edit it by hand |
| Human-readable "intent" | Machine-verified "proof" |

### Should it be committed?
**YES — always commit go.sum to version control.** It's essential for reproducible, secure builds. Everyone who builds your project (including CI/CD) checks downloaded code against these hashes; if a dependency was compromised or changed, the build FAILS instead of silently using bad code.

### What happens when dependencies change?
When you add/update/remove a dependency (via `go get`, `go mod tidy`, editing `go.mod`), Go recomputes and updates the relevant lines in `go.sum` automatically.

### INTERVIEW
**Q: Why do we need go.sum if we already have go.mod?**
**A:** go.mod declares desired versions; go.sum cryptographically verifies the actual downloaded bytes match a known-good hash, protecting against supply-chain tampering and guaranteeing byte-for-byte reproducible builds across machines.

---

## 13. Standard Library vs Third-Party

```go
import "fmt"                       // Standard library — ships WITH Go, no download needed, no version, no go.sum entry
import "github.com/google/uuid"    // Third-party — external, versioned, needs go.mod + go.sum
```

| | Standard Library | Third-Party |
|---|---|---|
| Downloaded? | No — built into Go itself | Yes — fetched over network |
| Versioned per-project? | No — tied to your Go version | Yes — pinned in go.mod |
| Security review | Maintained by Go team | Maintained by random authors — YOU must vet it |
| Appears in go.sum? | No | Yes |

### Why minimizing dependencies matters for backend engineering
- Every dependency is a **supply-chain risk** — if it's compromised, your production server is compromised.
- Every dependency adds to **build time, binary size, and Docker image size**.
- Fewer dependencies = fewer versions to keep updated = fewer security patches to track.
- Standard library alone (`net/http`, `database/sql`, `encoding/json`) can build a surprisingly complete backend without heavy frameworks.

---

## 14. Internal Packages

### WHAT
Any package whose import path contains a folder literally named `internal` can **only** be imported by code that lives inside the same parent tree as that `internal` folder.

```
myapp/
├── internal/
│   ├── auth/         ← importable ONLY from within myapp/...
│   └── database/      ← importable ONLY from within myapp/...
└── cmd/
    └── server/
        └── main.go     ← CAN import myapp/internal/auth ✔
```

If another module (say `github.com/someone-else/theirapp`) tries to `import "github.com/you/myapp/internal/auth"` — **compile error**. Go's tool enforces this rule at build time; it's not just a convention, it's actually checked.

### WHY
You want to expose a clean **public API** to the outside world, while keeping true implementation details private — even from other people who legitimately import your module for its public packages.

### The rule precisely
> A package under a path with `.../internal/...` can be imported only by code rooted at the directory that is the **parent of `internal`**.

### PRODUCTION
Almost every serious Go backend uses `internal/` to protect `handler`, `service`, `repository`, `config` — things that are implementation detail, not meant to be imported by other projects.

### INTERVIEW
**Q: What is an internal package and why use it?**
**A:** It's a Go-enforced way to make a package private to a specific part of your module tree — anything under `internal/` can only be imported by code inside the parent directory of that `internal` folder. It's used to hide implementation details from external consumers while keeping code technically inside the same module.

---

## 15. Import Cycles

### WHAT
```
Package A imports Package B
Package B imports Package A   ← ❌ Go REFUSES to compile this
```

### WHY Go prohibits them
Go compiles packages in **dependency order** — B must be fully compiled before A can use it (A depends on B). If B also depends on A, there's no valid order to compile them in — it's a chicken-and-egg problem the compiler cannot solve.

### What it usually MEANS architecturally
An import cycle is almost always a sign of **bad separation of concerns** — two things that shouldn't know about each other directly are tangled together.

### Bad vs Good example

```
❌ BAD:
repository → service → repository   (repository importing service that imports repository back = cycle)

✅ GOOD (one-directional flow):
handler → service → repository
```

Data/control flows ONE way: handler calls service, service calls repository. Repository never needs to know about service or handler.

### How to FIX a cycle
1. **Extract the shared piece** into a third, lower-level package that both can depend on (e.g. move shared types into a `model` package).
2. **Use interfaces** — define an interface in the higher-level package, and let the lower-level package implement it, so the dependency only flows one direction (this is the classic "dependency inversion" fix).
3. **Merge the two packages** if they're really so tightly coupled they should be one thing.

### INTERVIEW
**Q: Why does Go reject import cycles, and what does an import cycle usually indicate?**
**A:** Go compiles packages bottom-up in dependency order, and a cycle makes that ordering impossible. It usually signals a design flaw — two packages are too tightly coupled and should either be merged, split differently, or decoupled using an interface owned by the higher-level package.

---

## 16. The Go Tool

The `go` command is a **complete toolchain** — one binary that does compiling, testing, formatting, downloading, and more.

```bash
go help          # see all commands
```

| Command | What / Why | Example |
|---|---|---|
| `go run` | Compiles + runs immediately, throwaway binary | `go run main.go` |
| `go build` | Compiles, produces a binary file on disk | `go build -o server .` |
| `go test` | Runs tests | `go test ./...` |
| `go vet` | Static analysis — finds suspicious code (not style, actual bugs) | `go vet ./...` |
| `go fmt` | Auto-formats code to Go's standard style | `go fmt ./...` |
| `go doc` | Shows documentation from the terminal | `go doc fmt.Println` |
| `go list` | Lists packages/modules info | `go list ./...` |
| `go env` | Shows Go's environment configuration | `go env GOPATH` |
| `go mod` | Module management (init, tidy, etc.) | `go mod tidy` |
| `go get` | Adds/updates/removes a dependency in go.mod | `go get github.com/google/uuid@latest` |
| `go install` | Compiles and installs a binary into `$GOBIN` | `go install github.com/x/tool@latest` |
| `go clean` | Removes build/cache artifacts | `go clean -cache` |

---

## 17. go run vs go build

| | `go run main.go` | `go build` |
|---|---|---|
| Output | Compiles to a TEMPORARY binary (in a temp dir), runs it, then deletes it | Produces a REAL, PERMANENT binary file on disk |
| Where the file goes | Nowhere you can keep — temp folder, cleaned up after | Current directory (or `-o` path you specify) |
| When to use | Quick local testing, "just run my code" | Producing something to DEPLOY or SHIP (e.g. into Docker) |
| Speed for repeated runs | Slower per run (recompiles from cache but still launches every time) | Build once, run the binary many times — fastest for repeated execution |
| Build caching | Yes, Go caches compiled packages either way (`$GOCACHE`) | Yes, same cache |

### Interview-ready answer
> "`go run` compiles your code into a temporary executable, runs it immediately, then throws it away — good for quick local iteration. `go build` compiles a permanent binary and saves it to disk, which is what you'd actually deploy to a server or bake into a Docker image. In production, you always use `go build` (or `go install`), never `go run`."

---

## 18. go install vs go get

**Modern behavior (Go 1.17+ / current):**

| | `go get` | `go install` |
|---|---|---|
| Purpose | Adds, updates, or removes a dependency **in your current module's go.mod** | Compiles and installs a binary (into `$GOBIN`), does NOT touch your project's go.mod |
| Typical use | `go get github.com/google/uuid@v1.6.0` inside your project | `go install github.com/x/some-cli-tool@latest` to install a command-line tool globally |
| Builds a binary? | No (as of modern Go) — it only manages dependencies | Yes — produces an executable |

**⚠️ Older Go behavior (before Go 1.17):**
`go get` used to ALSO build and install binaries (it did both jobs). This caused confusion, so Go split the responsibilities: `go get` = manage dependencies only, `go install` = build/install binaries only. **Do not follow old tutorials that use `go get` to install CLI tools — use `go install` for that today.**

### INTERVIEW
**Q: What's the difference between go get and go install today?**
**A:** `go get` manages dependency versions in go.mod for your current module. `go install` compiles a package into an executable binary placed in `$GOBIN`, independent of any particular project — this is the modern way to install Go-based CLI tools.

---

## 19. go test

```bash
go test               # run tests in current package
go test ./...         # run tests in EVERY package in the module (recursively)
go test -v            # verbose output — show each test name and result
go test -run TestName # run only tests matching this name/pattern
go test -race         # enable the race detector (finds concurrency bugs)
go test -cover        # show test coverage percentage
```

### How `go test ./...` discovers packages
The `...` is a wildcard meaning "this directory and all subdirectories." Go walks the whole module tree, finds every folder that is a package, and for each one that has `_test.go` files, it compiles and runs those tests — one package at a time, in parallel by default across packages.

### PRODUCTION
Real backend CI pipelines almost always run:
```bash
go vet ./...
go test ./... -race -cover
go build ./...
```
before allowing a merge or deploy — catching bugs, race conditions, and build failures early.

---

## 20. go list, go doc, go vet

### go list — for debugging large repos
```bash
go list ./...            # list all package import paths in the module
go list -m all           # list the FULL dependency tree (module versions)
go list -deps ./...      # list every package a given package depends on
```
**Real debugging use:** "Why is package X even in my build?" → `go list -deps ./cmd/server | grep X` shows if/how it's pulled in.

### go doc — read docs without leaving the terminal
```bash
go doc fmt            # shows package-level doc for fmt
go doc fmt.Println    # shows doc for just this one function
```
This reads the SAME comments that show up on `pkg.go.dev` — because Go documentation comments (written directly above declarations) are the single source of truth for both.

### go vet — finds real bugs, not style issues
```bash
go vet ./...
```
Catches things like: format string mismatches (`fmt.Printf("%d", someString)`), unreachable code, suspicious struct tags, and other correctness bugs — this is DIFFERENT from `go fmt`, which only fixes STYLE/formatting, not bugs.

---

## 21. Real Backend Project Structure

```
myapp/
├── go.mod
├── go.sum
├── cmd/
│   └── server/
│       └── main.go            ← package main, the ONE entry point
├── internal/
│   ├── handler/                ← HTTP layer: parses requests, calls service
│   ├── service/                ← business logic
│   ├── repository/             ← talks to the database
│   ├── model/                  ← shared structs (User, Order, etc.)
│   ├── middleware/              ← auth checks, logging, etc.
│   └── config/                  ← reads env vars / config files
├── pkg/                        ← (optional) code meant to be reused by OTHER projects
└── migrations/                  ← SQL migration files
```

### Why `cmd/`?
Lets you have MULTIPLE entry points in one module — e.g. `cmd/server/main.go` (the API) and `cmd/worker/main.go` (a background job runner) — each its own `package main`, sharing the same `internal/` code.

### Why `internal/`?
As explained above — protects implementation packages from being imported by outside projects, even if someone depends on your module.

### Is `pkg/` actually required?
**No.** This is a common misconception. `pkg/` has NO special meaning to the Go compiler (unlike `internal/`, which IS enforced). It's purely a **convention** some teams use for "code that's safe for other projects to import." Many respected Go projects skip `pkg/` entirely and just use `internal/` + top-level packages. **Don't blindly copy this structure** — for a small project or a single microservice, a flat structure with just a few top-level packages is often simpler and better. Use `cmd/` + `internal/` when you genuinely have multiple binaries or want to explicitly protect internals; skip the ceremony for small projects.

### Dependency direction (very important)
```
handler → service → repository → model
```
Each layer only calls the layer "below" it. `repository` never calls `service`. This one-directional flow is WHY there's no import cycle, and why the code is easy to reason about.

---

## 22. Debugging Scenarios

### Scenario 1: `package X is not in std`
**Meaning:** You typed an import path that Go thinks should be in the standard library, but isn't — usually a typo, or you forgot the full path for a third-party package.
**Fix:** Check spelling. If it's third-party, use the FULL path like `github.com/user/repo`, not just `repo`.

### Scenario 2: `cannot find module providing package ...`
**Meaning:** Go looked in go.mod, in the module cache, and possibly the internet, and could not find this package anywhere.
**Fix:** Run `go get <import-path>` to add it, or `go mod tidy` if it's used in code but missing from go.mod. Check for typos in the import path.

### Scenario 3: `import cycle not allowed`
**Meaning:** Explained in Section 15 — two (or more) packages depend on each other, directly or through a chain.
**Fix:** Use `go list -deps` or read the error's printed cycle chain (Go shows you the exact loop: `A imports B imports C imports A`). Break the cycle using the fixes from Section 15.

### Scenario 4: `undefined: SomeFunction`
**Meaning:** Either (a) you forgot to write/import it, (b) it's misspelled, or (c) it's UNEXPORTED (lowercase) and you're calling it from a different package.
**Fix:** Check spelling, check the import, check capitalization if calling across packages.

### Scenario 5: A function exists but another package cannot access it
**Meaning:** It starts with a lowercase letter — it's unexported, only visible inside its own package.
**Fix:** If it genuinely needs to be public, capitalize it (and think carefully — exporting is a promise to keep that API stable).

### Scenario 6: Two files, same directory, different `package` names
**Meaning:** Go REQUIRES every file in one folder to declare the exact same package. This is a hard compile error: "found packages X and Y."
**Fix:** Rename the package line in the mismatched file to match the rest of the folder, or physically move that file into its own separate folder if it's genuinely meant to be a different package.

### Scenario 7: Dependency works locally but fails elsewhere
**Possible causes:**
- `go.sum` wasn't committed, so checksums can't be verified on the other machine.
- Different Go version installed (check the `go` directive in go.mod).
- A `replace` directive pointing to a local path that doesn't exist on the other machine.
- Module cache differences / private module access (auth/proxy config differences, e.g. `GOPRIVATE` not set the same way).
**Fix:** Always commit `go.mod` + `go.sum`; use `go env` to compare environments; check for any local-only `replace` lines.

---

## 23. Compiler-Level Mental Model

```
Source code (.go files)
      ↓
Package discovery      (go tool scans folders, groups files by package declaration)
      ↓
Import/dependency resolution   (resolve every import path → stdlib, module cache, or local)
      ↓
Type checking            (does every expression/type/function call actually make sense?)
      ↓
Compilation               (each package compiled into intermediate object code)
      ↓
Package/object generation  (compiled packages cached, ready to link)
      ↓
Linking                    (combine all compiled packages + runtime into ONE binary)
      ↓
Executable binary
```

**Correction to the simplified version:** Go's compiler doesn't compile file-by-file independently — it compiles **package-by-package**, and it does so in **dependency order** (packages with no dependencies first, then packages that depend on them, and so on up the tree, ending with `main`). This is exactly why import cycles are impossible to compile — there's no valid "order" for a cycle. Also, Go statically links almost everything by default (unlike C, where linking external `.so` files is common) — this is WHY Go binaries are typically large, self-contained files with no external runtime dependencies needed to run them.

---

## 24. Anti-Patterns

| Anti-pattern | Why it's bad | When it might be OK | Better alternative |
|---|---|---|---|
| **Giant packages** (one package, everything) | Hard to navigate, everything coupled, unclear boundaries | Truly tiny projects/scripts | Split by responsibility |
| **Generic names** (`utils`, `common`, `helpers`) | Becomes a dumping ground for unrelated code, tells reader nothing | Never really "OK" — always rename by real purpose | Name by what it PROVIDES: `stringutil`, `timeutil` |
| **Excessive `pkg/`** | Adds ceremony with no compiler benefit; often copied blindly | Genuinely public library used by external projects | Just use `internal/` + flat top-level packages |
| **Circular dependencies** | Compiler literally rejects it | Never acceptable | Extract shared code, use interfaces |
| **Overusing `internal/`** | Makes even trivial shared code hard to reuse across sibling modules unnecessarily | Genuinely private implementation details | Only wrap things that are TRUE implementation detail |
| **Overusing `init()`** | Hidden, hard-to-test, order-dependent side effects | Simple registration only (drivers, plugins) | Explicit initialization functions called from `main()` |
| **Dot imports** | Hides WHERE a function came from, risks name collisions | Rare test-assertion-library cases | Normal (or aliased) imports |
| **Unnecessary dependencies** | Bigger attack surface, slower builds, bigger binaries | Well-known, well-maintained, saves real time | Prefer standard library when it's "good enough" |
| **Packages by technical layer, not domain** (e.g. all "controllers" mixed regardless of feature) | Related domain logic gets scattered across many folders | Small apps where this is genuinely simpler | Group by domain/feature for larger apps |
| **Exporting everything** | No real encapsulation, every internal change becomes a breaking change | Never really — always be deliberate | Export only what's a genuine, intended API |
| **Splitting files just to reduce file size** | Splits with no logical grouping create noise, not clarity | Genuinely too-long files that ARE logically separable | Split along logical/responsibility boundaries, not line-count |

---

## 25. Production Engineering

```
Local development
      ↓
go vet ./...           ← catch real bugs before tests even run
      ↓
go test ./... -race -cover   ← correctness + concurrency safety + coverage
      ↓
go build ./...          ← confirm everything compiles cleanly
      ↓
Docker image build       ← usually a MULTI-STAGE build: compile in one stage,
      ↓                      copy ONLY the final binary into a tiny final image
Production binary         (small image size, no Go toolchain needed at runtime)
```

**Why this matters for backend engineering specifically:**
- **Build reproducibility:** `go.sum` + pinned Go version = the exact same binary, byte-identical logic, every time, on every machine (dev laptop, CI server, production).
- **Docker builds:** Go binaries are statically compiled (usually), so the final Docker image can be TINY (e.g. based on `scratch` or `alpine`) — no need for a Go runtime installed at all, unlike Python/Node.
- **CI/CD:** `go vet`, `go test -race`, and `go build` are the standard 3-step gate before merging.
- **Dependency security:** in supply-chain-sensitive environments, `go list -m all` + tools like `govulncheck` scan for known vulnerable dependency versions.
- **Monorepos/microservices:** `internal/` boundaries + `cmd/` multiple entry points let ONE module cleanly produce several independently deployable binaries (e.g. `cmd/api`, `cmd/worker`, `cmd/migrate`) while sharing common `internal/` code.

---

## 26. One-Page Revision Sheet

| Concept | Definition | Why it exists | Rule | Interview Keyword |
|---|---|---|---|---|
| Package | Folder of .go files, same `package` name | Organize + hide code | All files in a folder = same package name | "unit of compilation" |
| Module | Versioned project with go.mod | Dependency management + reproducibility | One go.mod per module | "unit of dependency management" |
| Import path | Address string to locate a package | Uniquely identify/download packages | Full path, no relative imports | "locates the package" |
| Package name | Short identifier used in code | Keep code readable | Declared via `package x` | "used in code" |
| Exported/unexported | Capital = public, lowercase = private | Encapsulation without keywords | First letter capitalization | "visibility by case" |
| Blank import `_` | Import for side-effects only | Trigger `init()` without direct use | `import _ "pkg"` | "registration pattern" |
| `init()` | Auto-run setup function | Package-level setup before main | Runs after vars, before main | "runs before main" |
| go.mod | Module manifest | Replace GOPATH, per-project versions | `module`, `go`, `require` | "dependency manifest" |
| go.sum | Checksum ledger | Verify integrity/reproducibility | Always commit it | "cryptographic verification" |
| internal/ | Restricted-visibility folder | Hide implementation from outside importers | Enforced by compiler | "compiler-enforced privacy" |
| Import cycle | A→B→A dependency loop | Impossible to compile | Break with interfaces/extraction | "no valid compile order" |
| go run | Compile + run temp binary | Fast iteration | Not for production | "temporary binary" |
| go build | Compile to permanent binary | Deployment artifact | Used in Docker/production | "permanent binary" |
| go get | Manage go.mod dependencies | Add/update/remove deps | Doesn't build binaries anymore | "dependency management" |
| go install | Build + install a binary | Install CLI tools | Doesn't touch project go.mod | "installs executables" |

---

## 27. Flashcards

```
Q1: What is the difference between a package and a module?
A1: A package is a folder of .go files compiled as one unit. A module is a
    versioned collection of packages defined by one go.mod file.

Q2: What determines if an identifier is exported?
A2: Whether its first letter is uppercase (exported) or lowercase (unexported).

Q3: What does the blank identifier `_` do in an import?
A3: Imports the package purely for its side effects (init()), without using
    any of its exported names directly.

Q4: When does init() run relative to main()?
A4: After all package-level variables are initialized and all imported
    packages are fully initialized, but always before main().

Q5: Can two files in the same folder declare different packages?
A5: No — Go requires every file in a directory to declare the same package
    (this is a compile-time error otherwise).

Q6: What is an import path?
A6: The string used in an import statement to locate a package — could be a
    standard library name, or a full module-based path like a GitHub URL.

Q7: What's the difference between import path and package name?
A7: Import path is the full address used to locate/download the package;
    package name is the short identifier used inside code.

Q8: Why does Go reject unused imports?
A8: To force clean code and prevent silent accumulation of unnecessary
    imports.

Q9: What is a blank import commonly used for in real backends?
A9: Registering database drivers (e.g. postgres, mysql) via their init().

Q10: What problem did Go modules solve that GOPATH could not?
A10: Per-project dependency versioning; GOPATH forced one global version of
     each dependency across all projects on a machine.

Q11: What does go.sum guarantee?
A11: That downloaded dependency code exactly matches a known cryptographic
     hash, preventing tampering and guaranteeing reproducible builds.

Q12: Should go.sum be committed to version control?
A12: Yes, always — it is essential for reproducible, secure builds.

Q13: What is an internal package?
A13: A package under a folder literally named `internal/`, importable only
     by code rooted at internal's parent directory; enforced by the compiler.

Q14: Why does Go prohibit import cycles?
A14: Packages are compiled in dependency order; a cycle makes no valid
     compile order possible.

Q15: How do you fix an import cycle?
A15: Extract shared code into a third package, or invert the dependency
     using an interface owned by the higher-level package.

Q16: Difference between go run and go build?
A16: go run compiles to a temporary binary and immediately runs+deletes it;
     go build produces a permanent binary saved to disk.

Q17: Difference between go get and go install today?
A17: go get manages dependencies in the current module's go.mod; go install
     compiles and installs a standalone executable, independent of any
     project's go.mod.

Q18: What does `go mod tidy` do?
A18: Adds missing dependencies to go.mod/go.sum and removes unused ones,
     based on scanning your actual imports.

Q19: What does `go test ./...` do?
A19: Recursively runs tests in every package under the current directory
     tree.

Q20: What is `go vet` for?
A20: Static analysis that finds likely real bugs (e.g. format string
     mismatches), not just style issues.

Q21: What's the difference between go vet and go fmt?
A21: go vet finds correctness bugs; go fmt only reformats code style —
     they serve different purposes entirely.

Q22: Why are relative imports discouraged in modern Go modules?
A22: They make an import's meaning depend on the caller's location on disk,
     breaking reproducibility and portability.

Q23: Why is `package main` special?
A23: It marks a package as producing an executable program, requiring a
     func main() as the entry point — unlike library packages.

Q24: Can the directory name differ from the package name?
A24: Yes, technically allowed, but discouraged by convention unless there's
     a good reason (e.g. cmd/server/ containing package main).

Q25: What's a "good" Go package name?
A25: Short, lowercase, describes what it PROVIDES (e.g. "auth"), not vague
     or generic like "utils" or "common".

Q26: What happens if you import a package you never use?
A26: Compile error: "imported and not used" — Go enforces this strictly.

Q27: What is dependency inversion in the context of import cycles?
A27: Defining an interface in the higher-level package and letting the
     lower-level package implement it, so imports flow only one direction.

Q28: Why is pkg/ not compiler-enforced like internal/ is?
A28: pkg/ has zero special meaning to the Go tool — it's purely a naming
     convention, unlike internal/, which the compiler actually checks.

Q29: What's the recommended CI gate sequence for a Go backend?
A29: go vet ./... → go test ./... -race -cover → go build ./...

Q30: Why are Go binaries often used in tiny Docker images (scratch/alpine)?
A30: Go typically produces statically linked binaries, needing no separate
     runtime installed in the final image.

Q31: What is a direct vs indirect dependency in go.mod?
A31: Direct: your code imports it directly. Indirect: it's a dependency of
     one of your dependencies, pulled in transitively.

Q32: What does `replace` in go.mod do?
A32: Substitutes one module source for another — commonly used to point a
     dependency to a local folder during development.

Q33: What does `go mod why <pkg>` help with?
A33: Shows the exact import chain that causes a given package to be part
     of your dependency tree — useful for debugging bloat.

Q34: What's the risk of exporting everything in a package?
A34: You lose encapsulation — every internal change risks becoming a
     breaking change for anyone depending on your package.

Q35: What is package scope vs file scope?
A35: Package-scope declarations (vars, funcs, types outside functions) are
     visible to ALL files in the package; imports are file-scoped only.

Q36: How many init() functions can a package have?
A36: As many as needed, even multiple per file; they run in the order they
     appear (file-by-file, top-to-bottom within a file).

Q37: Why shouldn't init() do heavy work like network calls?
A37: It runs invisibly and automatically, can't be skipped or mocked
     easily, making testing and startup behavior unpredictable.

Q38: What determines the order packages are initialized in?
A38: Dependency order — an imported package is always fully initialized
     before the package that imports it.

Q39: What's the difference between go.mod's `go` directive and a
     `require` entry?
A39: `go X.Y` sets the minimum Go language/toolchain version for the
     module; `require` lists actual dependency packages and versions.

Q40: Why is `database/sql` combined with blank-imported drivers a good
     design example?
A40: It cleanly separates the generic SQL interface (database/sql) from
     specific driver implementations, which register themselves via side
     effects, keeping the core package driver-agnostic.
```

---

## 28. Mock Interview

> Rules: I'll ask, you answer FIRST (in your head or out loud), THEN read the expected answer below it. Don't peek early, brother — that's how you actually learn.

### 🟢 Round 1 — Basic (10 Questions)

1. What is a package in Go?
2. What is the special meaning of `package main`?
3. What determines if a name is exported?
4. What is an import path?
5. What happens if you import a package and never use it?
6. What is go.mod for?
7. What's the difference between `go run` and `go build`?
8. Can two files in one folder belong to different packages?
9. What is `go fmt` used for?
10. What is the standard library?

<details>
<summary>👉 Click for expected answers</summary>

1. A folder of `.go` files sharing the same `package` declaration, compiled as one unit.
2. It marks the package as an executable program with a required `func main()`.
3. Capitalization of the first letter of the identifier.
4. The string used in an `import` statement to locate a package (stdlib name or full module path).
5. Compile error: "imported and not used."
6. It declares the module path, Go version, and dependencies — Go's dependency manifest.
7. `go run` compiles + runs a temporary binary then deletes it; `go build` produces a permanent binary file.
8. No — Go requires every file in a directory to declare the same package name.
9. Auto-formats code into Go's standard style (not a correctness check).
10. Packages that ship WITH Go itself (like `fmt`, `net/http`), requiring no download.

**What a strong candidate mentions extra:** the WHY behind each rule, not just the WHAT — e.g. for Q5, mentioning that this is a deliberate design choice to keep code clean.
**Common mistakes:** confusing package with module (Q1/Q6), thinking go fmt catches bugs (Q9).
</details>

---

### 🟡 Round 2 — Intermediate (15 Questions)

1. Explain the difference between import path, package name, and module path.
2. Why does Go use a blank import for database drivers?
3. Walk through the complete package initialization order.
4. What is an internal package and how is it enforced?
5. Why are relative imports discouraged in Go modules?
6. What's the difference between `go get` and `go install` in modern Go?
7. What does `go mod tidy` actually do?
8. Explain go.sum's purpose in one sentence.
9. What causes an import cycle, and why does Go reject it?
10. What's the difference between direct and indirect dependencies?
11. Why is `pkg/` not enforced by the compiler, unlike `internal/`?
12. When is using `init()` appropriate vs risky?
13. What does `go vet` catch that `go build` won't?
14. Why should go.sum be committed to version control?
15. Explain the dependency direction in a typical `handler → service → repository` backend.

<details>
<summary>👉 Click for expected answers</summary>

1. Import path = full address to locate the package (e.g. github URL); package name = short identifier used in code (`fmt`, `handler`); module path = the root identity of the whole project declared in go.mod.
2. So the driver's `init()` runs, calling `sql.Register(...)`, making it discoverable by name string at runtime — without needing to call its exported functions directly.
3. Imported packages initialize first (deepest dependency first) → package-level vars initialize → init() functions run in file order → then main() runs.
4. A package under an `internal/` folder, importable only by code rooted at internal's parent directory — enforced by the `go` tool at build/compile time, not just convention.
5. They tie import meaning to the caller's disk location, breaking reproducibility and portability — modern Go always uses full module-based paths.
6. `go get` manages dependencies in go.mod; `go install` compiles and installs a standalone binary, unrelated to any project's go.mod.
7. Scans actual code imports, adds missing entries to go.mod/go.sum, and removes unused ones.
8. It cryptographically verifies downloaded dependency code hasn't been tampered with, ensuring reproducible builds.
9. Two or more packages depend on each other, directly or via a chain — Go compiles in dependency order, and a cycle has no valid order.
10. Direct: imported by your own code. Indirect: pulled in transitively by one of your direct dependencies.
11. `pkg/` has zero special meaning to the Go compiler — pure convention; `internal/` is actually checked and enforced by the toolchain.
12. Appropriate: side-effect registration (drivers/plugins), simple lookup table setup. Risky: heavy logic, network calls, anything that should be explicit/testable.
13. Real correctness bugs — e.g. format string/argument mismatches, suspicious struct tags — not just compile errors.
14. To guarantee anyone building the project (including CI) verifies dependencies against the same known-good hashes, catching tampering or inconsistency.
15. Data/control flows one direction only — handler calls service, service calls repository — repository never calls back up, preventing cycles and keeping responsibilities clear.

**What a strong candidate mentions extra:** real production examples (drivers, CI pipelines), and the WHY (design reasoning) behind each rule.
**Common mistakes:** mixing up go get/go install old vs new behavior; thinking go vet is the same as go fmt; vague answers on import cycles without mentioning WHY compile order matters.
</details>

---

### 🔴 Round 3 — Advanced (15 Questions)

1. You see `import cycle not allowed` involving 4 packages. How do you debug and fix it?
2. Explain exactly how `sql.Open("postgres", ...)` finds the postgres driver, end-to-end.
3. Why does Go statically link binaries by default, and how does this affect Docker builds?
4. A dependency builds fine locally but fails in CI. List every possible cause you'd check.
5. Design a package structure for a backend with 3 binaries (API, worker, migration tool) sharing common logic. Justify your folder choices.
6. Why might exporting an unnecessary field/function become a long-term liability?
7. Explain how `go build` differs conceptually from a C compiler+linker pipeline, focused on package-level compilation.
8. When would you deliberately use a `replace` directive in go.mod, and what are the risks of leaving it in for production?
9. Why does Go's toolchain treat unused imports and unused local variables as compile errors rather than warnings — what's the philosophy?
10. Explain how you'd use `go list -deps` and `go mod why` together to investigate unexpected binary bloat.
11. What's the architectural problem an import cycle almost always signals, and how does dependency inversion fix it?
12. Why is putting business logic inside `init()` considered a testability anti-pattern?
13. Explain the tradeoffs of using `internal/` aggressively vs sparingly in a growing codebase.
14. How does Go's "package name vs import path" distinction help avoid a common naming collision problem when using two different libraries with the same package name?
15. Walk through what happens, step by step, from typing `go build` to getting a runnable binary.

<details>
<summary>👉 Click for expected answers</summary>

1. Use the exact cycle Go prints in the error (it lists the full chain). Use `go list -deps` on the involved packages to confirm the chain. Then decide: can shared types/logic be extracted into a lower-level package both sides depend on? Or should one side depend on an interface instead of a concrete type, inverting the dependency?
2. `_ "github.com/lib/pq"` blank import triggers pq's `init()`, which calls `sql.Register("postgres", &Driver{})`, storing the driver in `database/sql`'s internal registry map. `sql.Open("postgres", dsn)` looks up "postgres" in that map and uses the registered driver to create connections.
3. Static linking bundles the Go runtime + all dependencies into one self-contained binary with no external shared libraries needed at runtime — this lets Docker images be built from `scratch` or minimal `alpine`, since no separate language runtime needs installing in the final image, drastically shrinking image size and attack surface.
4. go.sum not committed; different Go version (check `go` directive); a `replace` pointing to a local-only path; private module access/auth differences (`GOPRIVATE`, proxy config); OS/architecture-specific dependency behavior; module cache differences.
5. `cmd/api/main.go`, `cmd/worker/main.go`, `cmd/migrate/main.go` — each its own `package main`. Shared logic in `internal/service`, `internal/repository`, `internal/model`, `internal/config`. This lets three binaries share business/data logic without duplicating code, while `internal/` prevents external modules from depending on implementation details.
6. Every exported identifier becomes part of your PUBLIC API contract — changing or removing it later can break every consumer depending on it, forcing you into slower, more careful versioning/deprecation cycles just because something was exported "just in case."
7. Go's compiler treats each PACKAGE as the compilation unit (not each file, unlike typical C translation units per .c file), compiling in dependency order and producing package objects that the linker then combines — this package-level granularity is also why import cycles are structurally impossible to compile.
8. Use `replace` for local development against a not-yet-published fork/patch, or pointing to a local folder while iterating on a shared library. Risk of leaving it in for production: your build silently depends on a path/version that may not exist or match in other environments (like CI or a teammate's machine), breaking reproducibility — should be removed before merging/releasing.
9. Go's design philosophy strongly favors deterministic, always-clean compiled output and discourages "silence" around potentially dead/forgotten code — treating these as hard errors (not warnings developers can ignore) forces immediate cleanup rather than accumulating cruft.
10. `go list -deps ./cmd/server` lists every package pulled into that binary; grep for suspicious/heavy packages. Then `go mod why <suspected-package>` shows the exact import chain responsible, letting you decide whether it's necessary or can be removed/replaced.
11. It signals two packages/layers that are too tightly coupled — often a "lower" layer reaching back up to depend on a "higher" layer. Fix: the higher-level package defines an interface describing what it needs; the lower-level package implements that interface without importing the higher-level package at all — dependency now flows only one direction (interface owner doesn't need to import the implementer).
12. `init()` runs automatically and invisibly with no arguments — you cannot inject test doubles/mocks into it, cannot control WHEN it runs relative to test setup, and cannot skip it selectively in different test scenarios, making behavior hard to isolate and verify.
13. Aggressive `internal/`: strong protection of implementation details, but can make legitimately reusable code across separate modules/projects harder to share, and forces broader `internal/` trees as the codebase grows. Sparingly: easier sharing but weaker protection of true internals, risking third parties depending on things you meant to change freely later. Balance: use `internal/` for anything NOT intended as a stable, external contract.
14. Since a package's NAME (used in code, e.g. `json.Marshal`) is separate from its import PATH (the download address), Go lets you alias imports (`stdjson "encoding/json"`) when two dependencies happen to share the same default package name, avoiding a naming collision without renaming any actual code inside those libraries.
15. `go build` → scans/discovers packages in dependency order → resolves every import (stdlib/module cache/local) → type-checks each package → compiles each package into intermediate object code (using cached results where unchanged) → links all compiled packages + Go runtime together into one binary → writes the final executable to disk (current dir or `-o` path).

**What separates a top candidate here:** connecting mechanism (HOW) to reasoning (WHY) to real production consequence, every single time — not stopping at definitions.
</details>

---

## 29. Coding Exercises

### 🟩 Beginner — Package/Import Syntax
Create two files in the SAME package: `main.go` has `func main()`, and `greet.go` has a function `func Greet(name string) string`. Call `Greet` from `main()` without any import between them. Confirm it compiles and runs.

### 🟨 Intermediate — Multiple Packages
Create a module with two packages: `mathutil` (exports `Add(a, b int) int`) and `main` (imports `mathutil` and prints `Add(2,3)`). Practice writing the correct import path based on your `go.mod` module name.

### 🟧 Advanced — Design a Backend Package Structure
Design (just the folder tree + package names, no need for full working code) a small backend with: HTTP handling, a "todo" business service, an in-memory repository, and a shared `Todo` model struct. Decide where `internal/` fits and why.

**Sample target structure (compare against yours):**
```
todoapp/
├── go.mod
├── cmd/server/main.go
├── internal/
│   ├── handler/todo_handler.go
│   ├── service/todo_service.go
│   ├── repository/todo_repository.go
│   └── model/todo.go
```

### 🟥 Expert — Find the Bugs
A teammate gives you this broken mini-repo description. Identify every category of error (package, import, module, dependency, visibility, cycle, initialization):

```
go.mod: module github.com/x/app

// file: internal/service/service.go
package service
import "github.com/x/app/internal/repository"
func GetUser() string { return repository.getUser() }   // (A)

// file: internal/repository/repository.go
package repository
import "github.com/x/app/internal/service"               // (B)
func getUser() string { return "bob" }
func notify() { service.Alert() }                          // (C)

// file: internal/repository/helper.go
package repo                                                // (D)
func helperFunc() {}

// file: cmd/server/main.go
package main
import "./internal/service"                                  // (E)
func main() { service.GetUser() }
```

<details>
<summary>👉 Click for the answer key</summary>

- **(A)** `repository.getUser()` is UNEXPORTED (lowercase `g`) — `service` package cannot call it from outside. **Visibility error.** Fix: rename to `GetUser` (exported).
- **(B)** `repository` importing `service`, while `service` also imports `repository` (via A) → **Import cycle.** Fix: break the cycle — repository shouldn't depend on service at all in a clean layered design.
- **(C)** Depends on the cycle in (B); also assumes `service.Alert` exists and is exported — not shown, likely another bug once the cycle is fixed.
- **(D)** File `helper.go` declares `package repo`, but the rest of the `repository/` folder uses `package repository` → **Package name mismatch error** (Go: "found packages repository and repo").
- **(E)** `import "./internal/service"` — **relative import**, not allowed/used in modern Go modules → **Import path error.** Fix: use the full path `github.com/x/app/internal/service`.

</details>

---

## 30. "If You Remember Only This"

1. A package is the unit of compilation and organization in Go — a folder of files sharing one `package` name.
2. A module is a versioned collection of packages, defined by one `go.mod` file.
3. Import path locates a package; package name is what you type in code; module path identifies the whole project — these three are NOT the same thing.
4. Exported = starts with a capital letter; unexported = starts lowercase. That's Go's entire visibility system.
5. Blank import (`_`) runs a package's `init()` for side effects, without using its exported names directly — classic use: database drivers.
6. `init()` runs after package-level variables are set and after all imported packages are fully initialized, but always before `main()`.
7. Go strictly forbids unused imports and unused local variables — this is by design, to force clean code.
8. `go.mod` replaced the old GOPATH system, giving each project its OWN pinned dependency versions.
9. `go.sum` cryptographically verifies dependency integrity — always commit it.
10. `internal/` is compiler-enforced privacy — only code rooted at its parent directory can import it. `pkg/` is NOT enforced; it's just convention.
11. Import cycles are impossible to compile because Go compiles packages in dependency order — cycles have no valid order.
12. `go run` = temporary binary for quick iteration. `go build` = permanent binary for deployment.
13. `go get` manages dependencies in go.mod today; `go install` builds/installs standalone binaries — they used to be conflated in old Go, but not anymore.
14. Relative imports are not used in modern Go modules — always use the full import path.
15. `go vet` finds real bugs; `go fmt` only fixes style — never confuse the two.
16. Dependency direction in a clean backend flows ONE way: `handler → service → repository` — never backward.
17. Go binaries are typically statically linked, making them ideal for tiny, self-contained Docker images.
18. Good package names describe what they PROVIDE, not vague catch-alls like `utils` or `common`.
19. A CI pipeline for Go typically runs: `go vet ./...` → `go test ./... -race -cover` → `go build ./...`.
20. Every important Go design choice trades off simplicity/explicitness against flexibility — understanding WHY, not just WHAT, is what separates strong candidates from memorizers.

---

## 31. Beyond the Syllabus — Top 1%

```
Beginner        → knows package/import syntax
Intermediate     → understands exported vs unexported, go.mod basics
Job-ready        → can structure a small backend, debug common import errors
Strong developer  → understands WHY each rule exists, uses internal/ and
                    interfaces deliberately to prevent cycles, writes clean
                    CI pipelines
Top 1% candidate  → all of the above PLUS the points below
```

### What separates the top 1% specifically for this topic:

- **They think in terms of compile-order and dependency graphs**, not just "imports." They can predict, before running `go build`, whether a design will create a cycle — because they understand Go compiles packages bottom-up.
- **They understand `internal/` is a DESIGN tool, not just a rule to obey.** They deliberately draw the internal/public boundary around what's a genuine stable contract vs. implementation detail — thinking about it the way a library author thinks about semantic versioning.
- **They know go.sum and reproducibility are a supply-chain security concern**, not just a "housekeeping file" — they can explain WHY an uncommitted go.sum is a real risk in a CI/CD pipeline.
- **They avoid `init()` for anything beyond registration**, because they've been burned by (or have seen) hidden startup bugs that are hard to test — and they can explain the testability argument clearly, not just repeat "avoid init()" as a rule.
- **They understand the philosophical difference** between Go's compile-time-enforced rules (unused imports, `internal/`, import cycles) and Go's convention-only rules (`pkg/`, matching folder-to-package names) — and never confuse "the compiler enforces this" with "this is just a style guide."
- **They can read an unfamiliar Go repository in minutes** by immediately locating `go.mod`, scanning `cmd/` for entry points, and tracing dependency direction through `internal/` — instead of getting lost reading every file top to bottom.
- **They know the go get/go install split is RECENT history**, and actively correct outdated tutorials/Stack Overflow answers that teach the old, conflated behavior.
- **They can explain blank imports as a "registration pattern"** and generalize it beyond just database drivers — recognizing the same pattern in image format decoders, `net/http/pprof`, and plugin systems.
- **They treat package naming as an API design decision**, not an afterthought — because a bad package name (`utils`) invites bad architecture over time (everything gets dumped there).
- **They default to standard library first**, adding a third-party dependency only when it's clearly justified — thinking about supply-chain risk and long-term maintenance cost, not just "is there a package for this."

---

## 32. Final Mental Map

```
MODULE  (go.mod defines this — the whole project)
  │
  ├── go.mod
  │      ├── module path
  │      ├── go version directive
  │      ├── require (direct + indirect deps)
  │      ├── replace / exclude / retract
  │      └── verified against go.sum (checksums)
  │
  ├── PACKAGE  (folder of .go files, same package name)
  │      │
  │      ├── source files (.go)
  │      ├── package declaration (package X)
  │      ├── package-level scope (shared across files, no import needed)
  │      ├── exported identifiers   (Capitalized — public API)
  │      ├── unexported identifiers (lowercase — private)
  │      ├── init() functions       (auto-run setup, registration pattern)
  │      └── special case: package main → requires func main(), builds an executable
  │
  ├── IMPORT
  │      │
  │      ├── import path      (address: stdlib name OR full module-based path)
  │      ├── package name      (short identifier used in code)
  │      ├── forms: normal / aliased / blank (_) / dot (.)
  │      └── NO relative imports in modern Go modules
  │
  ├── BOUNDARIES & SAFETY
  │      ├── internal/       → compiler-enforced private packages
  │      ├── import cycles    → compiler-rejected (no valid compile order)
  │      └── pkg/             → convention only, NOT compiler-enforced
  │
  └── GO TOOLCHAIN  (the `go` command)
         │
         ├── go build     → permanent binary, for deployment
         ├── go run       → temporary binary, for quick iteration
         ├── go test       → runs tests, ./... = recursive discovery
         ├── go vet        → finds real bugs (not style)
         ├── go fmt        → fixes style only
         ├── go mod        → init / tidy / download / graph / why
         ├── go get        → manage dependencies in go.mod (modern)
         ├── go install     → build + install standalone binaries (modern)
         ├── go list        → inspect packages/modules, debug bloat
         └── go doc         → read documentation from the terminal
```

---

### 🎯 You're ready when you can, without looking:
- Draw this whole map from memory.
- Explain the difference between package/module/import path/package name out loud, in one breath.
- Diagnose all 7 debugging scenarios from Section 22 just from the error message.
- Answer every Round 3 advanced question with the WHY, not just the WHAT.

All the best, brother. Take this slow, one section a day if needed — this is a strong, complete foundation. 💪