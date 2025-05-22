# RapidHash Go Implementation

A Go port of the [rapidhash](https://github.com/Nicoshev/rapidhash) algorithm - a very fast, high quality, platform-independent hashing algorithm that is the official successor to wyhash.

## Features

- **Fast**: Extremely fast for both short and large inputs
- **High Quality**: Passes all tests in SMHasher and SMHasher3
- **Platform Independent**: Works on all platforms without machine-specific instructions
- **Go Native**: Implements standard Go interfaces like `hash.Hash64`
- **Memory Safe**: No unsafe operations in the core algorithm (except for string optimization)

## Installation

```bash
go get github.com/vkudryk/rapidhash-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/vkudryk/rapidhash-go"
)

func main() {
    // Hash bytes
    data := []byte("hello world")
    hash := rapidhash.Hash(data)
    fmt.Printf("Hash: %016x\n", hash)
    
    // Hash string directly (optimized)
    hash2 := rapidhash.String("hello world")
    fmt.Printf("String hash: %016x\n", hash2)
    
    // Hash with custom seed
    hash3 := rapidhash.HashWithSeed(data, 12345)
    fmt.Printf("Seeded hash: %016x\n", hash3)
}
```

## API Reference

### Basic Functions

- `Hash(data []byte) uint64` - Hash bytes with default seed
- `HashWithSeed(data []byte, seed uint64) uint64` - Hash bytes with custom seed
- `String(s string) uint64` - Hash string with default seed (optimized)
- `StringWithSeed(s string, seed uint64) uint64` - Hash string with custom seed

### Type-Specific Functions

- `Uint64(value uint64) uint64` - Hash a single uint64 value
- `Uint64WithSeed(value, seed uint64) uint64` - Hash uint64 with custom seed
- `Uint32(value uint32) uint64` - Hash a single uint32 value
- `Uint32WithSeed(value uint32, seed uint64) uint64` - Hash uint32 with custom seed

### Hasher Interface

The `Hasher` type implements Go's standard `hash.Hash64` interface:

```go
hasher := rapidhash.NewHasher()
hasher.Write([]byte("hello"))
hasher.Write([]byte(" world"))
hash := hasher.Sum64()

// Or with custom seed
hasher2 := rapidhash.NewHasherWithSeed(12345)
```

## Examples

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/vkudryk/rapidhash-go"
)

func main() {
    // Different ways to hash the same data
    data := []byte("example")
    
    fmt.Printf("Hash: %016x\n", rapidhash.Hash(data))
    fmt.Printf("String: %016x\n", rapidhash.String("example"))
    
    // These should produce the same result
    hasher := rapidhash.NewHasher()
    hasher.Write(data)
    fmt.Printf("Hasher: %016x\n", hasher.Sum64())
}
```

### Using as a HashMap Hash Function

```go
package main

import (
    "fmt"
    "github.com/vkudryk/rapidhash-go"
)

type MyKey struct {
    ID   uint64
    Name string
}

func (k MyKey) Hash() uint64 {
    hasher := rapidhash.NewHasher()
    
    // Hash the ID
    hasher.Write((*[8]byte)(unsafe.Pointer(&k.ID))[:])
    
    // Hash the name
    hasher.Write([]byte(k.Name))
    
    return hasher.Sum64()
}

func main() {
    key := MyKey{ID: 12345, Name: "example"}
    hash := key.Hash()
    fmt.Printf("Key hash: %016x\n", hash)
}
```

### Streaming Large Data

```go
package main

import (
    "fmt"
    "io"
    "os"
    "github.com/vkudryk/rapidhash-go"
)

func hashFile(filename string) (uint64, error) {
    file, err := os.Open(filename)
    if err != nil {
        return 0, err
    }
    defer file.Close()
    
    hasher := rapidhash.NewHasher()
    
    buffer := make([]byte, 32*1024) // 32KB buffer
    for {
        n, err := file.Read(buffer)
        if n > 0 {
            hasher.Write(buffer[:n])
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            return 0, err
        }
    }
    
    return hasher.Sum64(), nil
}

func main() {
    hash, err := hashFile("example.txt")
    if err != nil {
        panic(err)
    }
    fmt.Printf("File hash: %016x\n", hash)
}
```

## Performance

RapidHash is designed to be extremely fast. Here are some typical benchmark results (Apple M1 Pro):

```
BenchmarkHash-10         	171148759	         6.852 ns/op
BenchmarkHashShort-10    	256948310	         4.666 ns/op
BenchmarkHashLong-10     	 1000000	      1022 ns/op
BenchmarkString-10       	160144948	         7.468 ns/op
BenchmarkUint64-10       	593450163	         2.013 ns/op
BenchmarkHasher-10       	145452333	         8.258 ns/op
```

## Constants and Configuration

```go
const RapidSeed uint64 = 0xbdd89aa982704029  // Default seed

// Default secret parameters (same as C implementation)
var rapidSecret = [3]uint64{
    0x2d358dccaa6c78a5,
    0x8bb84b93962eacc9,
    0x4b33a62ed433d4a3,
}
```

## Testing

Run the test suite:

```bash
go test -v
go test -bench=.
```

## Implementation Notes

- This implementation maintains compatibility with the original C rapidhash
- Uses Go's `bits.Mul64` for efficient 64x64→128 bit multiplication
- Properly handles endianness for cross-platform compatibility
- The `String()` function uses a small unsafe optimization to avoid allocation
- No external dependencies

## License

BSD 2-Clause License - same as the original rapidhash implementation.

## Credits

Go implementation authored by [Vladyslav Kudryk](https://github.com/vkudryk)

Based on rapidhash by Nicolas De Carli, which is itself based on wyhash by Wang Yi.