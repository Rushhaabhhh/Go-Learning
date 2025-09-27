# Go Language: Comprehensive Notes

## Table of Contents
1. [Language Overview & Philosophy](#language-overview--philosophy)
2. [Execution Model](#execution-model)
3. [Memory Management & Storage](#memory-management--storage)
4. [Type System & Variables](#type-system--variables)
5. [Structs & Memory Alignment](#structs--memory-alignment)
6. [Pointers & Reference Semantics](#pointers--reference-semantics)

---

## 1. Language Overview & Philosophy

### Why Go? The Problem Statement

Go was designed by Robert Griesemer, Rob Pike, and Ken Thompson at Google (2007, released 2009) to address specific pain points:

- **Slow compilation** (C++, Java compilation times)
- **Complexity in concurrent systems** (multi-threading overhead, race conditions)
- **Memory safety without garbage collection overhead** (C's raw performance, but safer)
- **Cross-platform deployment** (single binary, no runtime dependencies)

### Key Advantages

| Aspect | Advantage |
|--------|-----------|
| **Compilation** | Compiled to native machine code; single binary deployment |
| **Concurrency** | Lightweight goroutines (100k+ simultaneously); channels for safe communication |
| **Memory Safety** | Garbage collection; no manual memory management |
| **Type Safety** | Strong static typing; catches errors at compile time |
| **Execution** | Direct hardware execution (no VM); uses native OS threads + goroutine scheduler |
| **Simplicity** | Minimal syntax; easy to read and maintain |
| **Performance** | Near C/C++ performance for compute-bound tasks |

### Drawbacks & Limitations

- **Garbage Collection pauses** (predictability issues in latency-critical systems)
- **No generics** (pre-1.18; now available with type parameters, Go 1.18+)
- **Less expressive type system** (compared to Rust, TypeScript)
- **Large binary size** (compared to C; smaller than JVM languages)
- **Limited standard library** (purposefully minimal; philosophy: "less is more")

### Go's Revolution in Modern Development

- **Cloud-native computing** (Kubernetes, Docker built in Go)
- **Microservices** (simplicity + performance ideal for services)
- **DevOps tooling** (Terraform, Prometheus, Grafana, Hugo)
- **Concurrency model** (goroutines changed how we think about parallelism)

---

## 2. Execution Model

### Binary Execution & Hardware Architecture

Go compiles to **native machine code** specific to the target architecture:

```bash
# Go detects architecture at compile time
go build -o myapp                    # Compiles for local machine
GOOS=linux GOARCH=amd64 go build    # Cross-compile for Linux x64
GOOS=windows GOARCH=386 go build    # Windows x86 (32-bit)
```

**Memory Block Size (Default Allocation):**
- **x64 (amd64)**: 8-byte blocks (64 bits)
- **x86 (386)**: 4-byte blocks (32 bits)
- **ARM**: Varies (32-bit: 4-byte, 64-bit: 8-byte)

The Go runtime respects these native alignments for CPU efficiency.

### Runtime Execution Flow

```
Source Code (.go)
      ↓
Go Compiler (gc)
      ↓
Machine Code (ELF/PE/Mach-O)
      ↓
OS Kernel + Runtime (goroutine scheduler, GC)
      ↓
Direct CPU Execution (no VM layer)
```

---

## 3. Memory Management & Storage

### Memory Model: Stack vs. Heap

Go uses two primary memory regions:

| Memory Region | Characteristics | Lifetime | Use Case |
|---------------|-----------------|----------|----------|
| **Stack** | LIFO allocation; blazing fast; per-goroutine | Function scope | Local variables, function calls, parameters |
| **Heap** | Dynamic allocation; GC-managed; shared across goroutines | Until unreferenced | Pointers, large structs, data shared across functions |

### Stack Memory Architecture

**Stack Frame Structure:**

```
Higher Addresses
┌─────────────────┐
│ Function Args   │ (Parameters passed to function)
├─────────────────┤
│ Return Address  │ (Where to return after function)
├─────────────────┤
│ Local Variables │ (Variables declared in function)
├─────────────────┤
│ (Unused Space)  │ (Expansion room for child functions)
└─────────────────┘
Lower Addresses
```

**Key Characteristics:**

- **Fixed size known at compile time**: Compiler calculates stack space for each function
- **Per-goroutine allocation**: Each goroutine gets its own stack (2-4 KB initially, grows dynamically)
- **Zero allocation cost**: Stack pointer increment is a single CPU instruction
- **Automatic cleanup**: Variables deallocated when function returns (stack pointer decrements)

**Stack Space Efficiency:**

```go
// Stack: ~24 bytes (3 x int64 on x64)
func example() {
    var a int64 = 10      // 8 bytes
    var b int64 = 20      // 8 bytes
    var c int64 = 30      // 8 bytes
}

// Goroutine stack growth:
// Initial: 2 KB | Can grow to multiple MB as needed
```

### Heap Allocation & Garbage Collection

**When variables escape to heap:**

```go
func heapAllocation() *int {
    x := 42
    return &x  // x escapes to heap (address returned to caller)
}

func stackAllocation() int {
    y := 100
    return y   // y stays on stack (value copied, no escape)
}
```

**Escape Analysis (Compiler optimization):**

Go's compiler uses "escape analysis" to determine if variables must live on heap:

```go
// Stack allocation (no escape)
func stackVar() {
    var buf [1024]byte
    _ = buf
}

// Heap allocation (escapes)
func heapVar() *[]byte {
    buf := make([]byte, 1024)
    return &buf  // Escapes
}
```

**Garbage Collection:**

- **Mark & Sweep algorithm** (with concurrent collection in recent versions)
- **Triggers**: Heap pressure, allocations, manual `runtime.GC()`
- **Tuning**: `GOGC` environment variable (default 100 = collect when heap doubles)

```go
import "runtime"

runtime.GC()  // Manual GC trigger
var m runtime.MemStats
runtime.ReadMemStats(&m)
fmt.Println("Alloc:", m.Alloc)  // Bytes allocated
```

---

## 4. Type System & Variables

### Basic Types & Sizing

```go
// Integer types (sign & size specify)
var a int      // Architecture-dependent (32 or 64 bit)
var b int8     // 1 byte  (-128 to 127)
var c int16    // 2 bytes (-32768 to 32767)
var d int32    // 4 bytes (rune alias)
var e int64    // 8 bytes

// Unsigned integers
var f uint     // Architecture-dependent
var g uint8    // 1 byte (byte alias)
var h uint16   // 2 bytes
var i uint32   // 4 bytes
var j uint64   // 8 bytes

// Floating-point
var k float32  // 4 bytes (IEEE 754)
var l float64  // 8 bytes (IEEE 754)

// Other
var m bool     // 1 byte
var n string   // String header (24 bytes: ptr + len + cap)
```

### Variable Declaration Patterns

**Standard Declaration (Explicit Type):**

```go
var age int = 25
var name string = "Alice"
var pi float64 = 3.14159
```

**Type Inference (Short-hand, no `var` keyword):**

```go
age := 25           // int (inferred)
name := "Alice"     // string (inferred)
pi := 3.14159       // float64 (inferred)
isActive := true    // bool (inferred)
```

**Rules for short-hand (`:=`):**
- Only inside functions (not at package level)
- Variables must be new (can't reassign with `:=`)
- Type inferred from right-hand expression

**Implicit Type Conversion (Not Allowed):**

```go
var a int32 = 10
var b int64 = a     // ❌ COMPILE ERROR: cannot assign int32 to int64
```

Go enforces explicit conversion—no implicit type coercion.

### Explicit Type Conversion & Casting

**Numeric Type Conversion:**

```go
// Syntax: targetType(value)
var a int32 = 100
var b int64 = int64(a)          // int32 → int64

var x float32 = 3.14
var y float64 = float64(x)      // float32 → float64

var n uint = 42
var m int = int(n)              // uint → int (risky if n > math.MaxInt)
```

**String to Number:**

```go
import "strconv"

str := "123"
num, err := strconv.Atoi(str)           // string → int
if err != nil {
    fmt.Println("Conversion failed")
}

str64 := "123"
num64, _ := strconv.ParseInt(str64, 10, 64)  // string → int64
result := int(num64)

str_float := "3.14"
numFloat, _ := strconv.ParseFloat(str_float, 64)  // string → float64
```

**Number to String:**

```go
import "strconv"

num := 123
str := strconv.Itoa(num)                // int → string

num64 := int64(456)
str64 := strconv.FormatInt(num64, 10)   // int64 → string (base 10)

numFloat := 3.14
strFloat := strconv.FormatFloat(numFloat, 'f', 2, 64)  // float64 → string
```

**Byte/Rune to String:**

```go
// rune = int32 (Unicode codepoint)
r := rune('A')
str := string(r)                // rune → string = "A"

// []byte to string
bytes := []byte{72, 101, 108, 108, 111}
str := string(bytes)            // []byte → string = "Hello"

// string to []byte
str := "Hello"
bytes := []byte(str)            // string → []byte
```

**Type Assertion (Interface to Concrete Type):**

```go
var x interface{} = "Hello"

// Type assertion syntax: value.(ConcreteType)
str, ok := x.(string)           // Safe assertion with comma-ok
if ok {
    fmt.Println("It's a string:", str)
}

num, ok := x.(int)              // Will be false
if !ok {
    fmt.Println("Not an int")
}

// Unsafe assertion (panics if wrong)
str := x.(string)               // Direct (panic if type mismatch)
```

---

## 5. Structs & Memory Alignment

### Memory Alignment & Padding

**Why Alignment Matters:**

CPUs read data most efficiently when aligned to natural boundaries:
- `int8`: 1-byte alignment (any address)
- `int16`: 2-byte alignment (addresses divisible by 2)
- `int32`: 4-byte alignment (addresses divisible by 4)
- `int64`: 8-byte alignment (addresses divisible by 8)
- `float64`: 8-byte alignment

Misaligned access is slower (or causes crashes on some architectures).

**Struct Padding Example:**

```go
// ❌ INEFFICIENT: 24 bytes (with padding)
type BadStruct struct {
    A bool        // 1 byte @ offset 0 | padding 1: [1 byte]
    B int64       // 8 bytes @ offset 8 | (aligned to 8-byte boundary)
    C int8        // 1 byte @ offset 16 | padding 7: [7 bytes]
}

// ✅ EFFICIENT: 16 bytes (no wasted padding)
type GoodStruct struct {
    B int64       // 8 bytes @ offset 0
    A bool        // 1 byte @ offset 8
    C int8        // 1 byte @ offset 9 | padding 6: [6 bytes]
}
```

**Alignment Rule (Compiler-enforced):**

```
Field offset = ceil(field_offset / field_alignment) * field_alignment
Struct size = ceil(last_field_offset / max_field_alignment) * max_field_alignment
```

### Memory Layout Calculation

For `GoodStruct` above:

```
Offset  Field     Size    Alignment   Total Size
─────────────────────────────────────────────────
0       B (int64)  8 bytes  8-byte    = 8 bytes
8       A (bool)   1 byte   1-byte    = 1 byte
9       C (int8)   1 byte   1-byte    = 1 byte
────────────────────────────────────────────────
Struct size: ceil(10 / 8) × 8 = 16 bytes
```

### Ordering Fields: Largest to Smallest

**Strategy: Place largest fields first to minimize padding**

```go
// ❌ POOR: 40 bytes
type Config struct {
    Name   string    // 24 bytes (slice header)
    Count  int8      // 1 byte
    ID     int64     // 8 bytes (padding before = 7)
    Active bool      // 1 byte
}

// ✅ OPTIMAL: 40 bytes (same size, but deliberate)
type Config struct {
    Name   string    // 24 bytes @ offset 0
    ID     int64     // 8 bytes @ offset 24
    Count  int8      // 1 byte @ offset 32
    Active bool      // 1 byte @ offset 33
    // padding 6 bytes to align struct size to 40
}
```

### Using `unsafe.Sizeof()` to Inspect

```go
import (
    "fmt"
    "unsafe"
)

type Person struct {
    Name    string    // 24 bytes
    Age     int32     // 4 bytes
    Active  bool      // 1 byte
}

fmt.Println("Size:", unsafe.Sizeof(Person{}))  // 32 bytes (24+4+1+3 padding)

// Inspect individual field offsets
p := Person{}
fmt.Println("Name offset:", unsafe.Offsetof(p.Name))    // 0
fmt.Println("Age offset:", unsafe.Offsetof(p.Age))      // 24
fmt.Println("Active offset:", unsafe.Offsetof(p.Active)) // 28
```

### Anonymous Structs

**Embedding Anonymous Structs (Composition):**

```go
// Anonymous struct type (no name, inline declaration)
var person struct {
    Name string
    Age  int
}

person.Name = "Alice"
person.Age = 30
```

**Nested Anonymous Structs:**

```go
var config struct {
    Database struct {
        Host string
        Port int
    }
    Server struct {
        Timeout int
    }
}

config.Database.Host = "localhost"
config.Database.Port = 5432
```

**Anonymous Struct in Functions:**

```go
func getUserInfo() struct {
    Name  string
    Email string
} {
    return struct {
        Name  string
        Email string
    }{
        Name:  "Alice",
        Email: "alice@example.com",
    }
}
```

**Embedding Named Structs (Composition Pattern):**

```go
type Address struct {
    Street string
    City   string
}

type Person struct {
    Name    string
    Address Address  // Named struct field
}

p := Person{
    Name: "Bob",
    Address: Address{
        Street: "Main St",
        City:   "NYC",
    },
}

// Access: p.Address.Street
```

**Method Receivers on Anonymous Structs:**

```go
// Can't define methods on anonymous types directly
// Solution: Define named type

type Config struct {
    Host string
    Port int
}

func (c Config) String() string {
    return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
```

---

## 6. Pointers & Reference Semantics

### Pass by Value (Default in Go)

**All function parameters are copied:**

```go
func increment(x int) {
    x++  // Modifies copy, not original
}

func main() {
    num := 5
    increment(num)
    fmt.Println(num)  // Still 5
}
```

**Memory Flow:**

```
Stack of main():
├─ num = 5 @ address 0x1000

Function call increment(x int):
├─ x = 5 (copy) @ address 0x2000  [NEW STACK FRAME]
└─ x++ increments copy only

Return to main():
├─ num still = 5 @ address 0x1000
```

### Pass by Address (Pointer Semantics)

**Syntax: `*Type` declares pointer, `&variable` gets address**

```go
func increment(x *int) {
    *x++  // *x dereferences pointer, modifies original
}

func main() {
    num := 5
    increment(&num)    // Pass address
    fmt.Println(num)   // Now 6
}
```

**Memory Flow:**

```
Stack of main():
├─ num = 5 @ address 0x1000

Function call increment(x *int):
├─ x = 0x1000 (address copy) @ address 0x2000
└─ *x++ dereferences to address 0x1000, increments value there

Return to main():
├─ num now = 6 @ address 0x1000 (modified via pointer)
```

**Pointer Declaration & Dereferencing:**

```go
var ptr *int            // Declare pointer to int
var x int = 42
ptr = &x                // ptr holds address of x
fmt.Println(*ptr)       // Dereference: prints 42

*ptr = 100              // Modify value at address
fmt.Println(x)          // x now 100
```

### Nil Pointer

```go
var ptr *int            // Uninitialized pointer is nil
fmt.Println(ptr)        // <nil>
fmt.Println(ptr == nil) // true

// Dereferencing nil causes panic
// value := *ptr        // ❌ PANIC: runtime error
```

**Safe nil check:**

```go
var ptr *int

if ptr != nil {
    fmt.Println(*ptr)
} else {
    fmt.Println("Pointer is nil")
}
```

### Pointer to Struct (Method Receivers)

```go
type Rectangle struct {
    Width  float64
    Height float64
}

// Value receiver (receives copy)
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Pointer receiver (can modify struct)
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    
    area := rect.Area()        // Works with value
    fmt.Println(area)          // 50
    
    rect.Scale(2)              // Must use pointer to modify
    fmt.Println(rect.Area())   // Now 200
}
```

### Pointer Comparison & Aliasing

```go
var x, y int = 10, 10
var px, py *int = &x, &y

px == py      // false (different addresses)
*px == *py    // true (same value: 10)

var pz = px
px == pz      // true (same address—aliasing)
```

### Common Pointer Pitfalls

**❌ Returning address of local variable (escapes to heap):**

```go
func getPointer() *int {
    x := 42
    return &x  // x escapes to heap; compiler handles this
}

func main() {
    ptr := getPointer()
    fmt.Println(*ptr)  // Safe (Go moves x to heap)
}
```

**✅ Go's escape analysis handles this automatically.**

**❌ Pointer to slice element (slice may reallocate):**

```go
s := []int{1, 2, 3}
ptr := &s[0]       // Risky: slice may grow, reallocate
s = append(s, 4)   // Reallocation → ptr now invalid
```

### Pointer-to-Pointer

```go
var x int = 42
var px *int = &x
var ppx **int = &px

fmt.Println(**ppx)  // 42 (dereference twice)

**ppx = 100         // Modify x through double pointer
fmt.Println(x)      // 100
```

**Use Case: Modifying pointers themselves**

```go
func updatePointer(pp **int, newValue int) {
    newVar := newValue
    *pp = &newVar      // Change what pp points to
}
```

---

## Summary Table: Quick Reference

| Concept | Syntax | Use Case |
|---------|--------|----------|
| Declare variable | `var x int` or `x := 10` | Store data on stack |
| Type conversion | `int64(x)` | Convert between types |
| Declare pointer | `var ptr *int` | Store memory address |
| Get address | `ptr := &x` | Create pointer to variable |
| Dereference | `*ptr` | Access value at address |
| Pass by value | `func(x int)` | Function receives copy |
| Pass by reference | `func(x *int)` | Function receives address |
| Method receiver (value) | `func (r Rect) Method()` | Immutable receiver |
| Method receiver (pointer) | `func (r *Rect) Method()` | Mutable receiver |
| Check nil | `if ptr != nil` | Safe pointer access |

---