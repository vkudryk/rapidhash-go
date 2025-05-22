package rapidhash

import (
	"fmt"
	"hash"
	"testing"
)

func TestBasicHashing(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"single_byte", "a"},
		{"short", "hello"},
		{"medium", "hello world"},
		{"long", "The quick brown fox jumps over the lazy dog"},
		{"very_long", "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := []byte(tt.input)

			// Test Hash function
			h1 := Hash(data)

			// Test HashWithSeed function
			h2 := HashWithSeed(data, RapidSeed)

			// Should be the same when using default seed
			if h1 != h2 {
				t.Errorf("Hash() and HashWithSeed() with default seed should be equal: %x != %x", h1, h2)
			}

			// Test String function
			h3 := String(tt.input)
			if h1 != h3 {
				t.Errorf("Hash() and String() should be equal: %x != %x", h1, h3)
			}

			// Test with different seed should produce different result
			h4 := HashWithSeed(data, 12345)
			if h1 == h4 && len(data) > 0 {
				t.Errorf("Different seeds should produce different hashes for non-empty input")
			}

			t.Logf("Input: %q -> Hash: %016x", tt.input, h1)
		})
	}
}

func TestConsistency(t *testing.T) {
	data := []byte("test data for consistency check")

	// Hash should be consistent across multiple calls
	h1 := Hash(data)
	h2 := Hash(data)

	if h1 != h2 {
		t.Errorf("Hash should be consistent: %x != %x", h1, h2)
	}
}

func TestEdgeCases(t *testing.T) {
	// Test various input lengths around boundaries
	for i := 0; i <= 64; i++ {
		data := make([]byte, i)
		for j := range data {
			data[j] = byte(j)
		}

		h := Hash(data)
		t.Logf("Length %d: %016x", i, h)

		// Should not panic
		if h == 0 && i > 0 {
			// Zero hash is unlikely but not impossible, just log it
			t.Logf("Got zero hash for length %d", i)
		}
	}
}

func TestHasherInterface(t *testing.T) {
	// Test that Hasher implements hash.Hash64
	var _ hash.Hash64 = (*Hasher)(nil)

	hasher := NewHasher()

	// Test writing data
	data := []byte("hello world")
	n, err := hasher.Write(data)
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Write returned wrong count: %d != %d", n, len(data))
	}

	// Test Sum64
	h1 := hasher.Sum64()

	// Test direct hash
	h2 := Hash(data)

	if h1 != h2 {
		t.Errorf("Hasher.Sum64() and Hash() should be equal: %x != %x", h1, h2)
	}

	// Test Sum
	sum := hasher.Sum(nil)
	if len(sum) != 8 {
		t.Errorf("Sum should return 8 bytes: %d", len(sum))
	}

	// Test Reset
	hasher.Reset()
	n, _ = hasher.Write([]byte("different"))
	h3 := hasher.Sum64()

	if h3 == h1 {
		t.Errorf("Hash should be different after reset and different input")
	}

	// Test Size and BlockSize
	if hasher.Size() != 8 {
		t.Errorf("Size should be 8: %d", hasher.Size())
	}

	if hasher.BlockSize() != 48 {
		t.Errorf("BlockSize should be 48: %d", hasher.BlockSize())
	}
}

func TestSpecificTypes(t *testing.T) {
	// Test Uint64 hashing
	val64 := uint64(0x123456789abcdef0)
	h1 := Uint64(val64)
	h2 := Uint64WithSeed(val64, RapidSeed)

	if h1 != h2 {
		t.Errorf("Uint64() and Uint64WithSeed() with default seed should be equal")
	}

	// Test Uint32 hashing
	val32 := uint32(0x12345678)
	h3 := Uint32(val32)
	h4 := Uint32WithSeed(val32, RapidSeed)

	if h3 != h4 {
		t.Errorf("Uint32() and Uint32WithSeed() with default seed should be equal")
	}

	t.Logf("Uint64(%016x) -> %016x", val64, h1)
	t.Logf("Uint32(%08x) -> %016x", val32, h3)
}

func TestLargeInputs(t *testing.T) {
	// Test with large inputs to exercise the chunked processing
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i)
		}

		h := Hash(data)
		t.Logf("Size %d: %016x", size, h)

		// Verify it's the same when split across multiple writes
		hasher := NewHasher()
		chunkSize := 1000
		for i := 0; i < len(data); i += chunkSize {
			end := i + chunkSize
			if end > len(data) {
				end = len(data)
			}
			hasher.Write(data[i:end])
		}

		h2 := hasher.Sum64()
		if h != h2 {
			t.Errorf("Chunked hash should equal single hash for size %d: %x != %x", size, h, h2)
		}
	}
}

func TestSeedVariation(t *testing.T) {
	data := []byte("test data")
	seeds := []uint64{0, 1, 42, 0xdeadbeef, ^uint64(0)}

	hashes := make(map[uint64]bool)

	for _, seed := range seeds {
		h := HashWithSeed(data, seed)
		if hashes[h] {
			t.Errorf("Duplicate hash %x for seed %x", h, seed)
		}
		hashes[h] = true
		t.Logf("Seed %016x -> %016x", seed, h)
	}
}

// Benchmark tests
func BenchmarkHash(b *testing.B) {
	data := []byte("The quick brown fox jumps over the lazy dog")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Hash(data)
	}
}

func BenchmarkHashShort(b *testing.B) {
	data := []byte("hello")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Hash(data)
	}
}

func BenchmarkHashLong(b *testing.B) {
	data := make([]byte, 10000)
	for i := range data {
		data[i] = byte(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Hash(data)
	}
}

func BenchmarkString(b *testing.B) {
	s := "The quick brown fox jumps over the lazy dog"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = String(s)
	}
}

func BenchmarkUint64(b *testing.B) {
	val := uint64(0x123456789abcdef0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Uint64(val)
	}
}

func BenchmarkHasher(b *testing.B) {
	data := []byte("The quick brown fox jumps over the lazy dog")
	hasher := NewHasher()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hasher.Reset()
		hasher.Write(data)
		_ = hasher.Sum64()
	}
}

// Example usage
func ExampleHash() {
	data := []byte("hello world")
	h := Hash(data)
	fmt.Printf("%016x\n", h)
}

func ExampleString() {
	h := String("hello world")
	fmt.Printf("%016x\n", h)
}

func ExampleHasher() {
	hasher := NewHasher()
	hasher.Write([]byte("hello"))
	hasher.Write([]byte(" "))
	hasher.Write([]byte("world"))
	h := hasher.Sum64()
	fmt.Printf("%016x\n", h)
}

// Test vector verification (these would need to be compared against the C implementation)
func TestKnownVectors(t *testing.T) {
	// Note: These test vectors would need to be generated from the original C implementation
	// For now, we just test that the hashes are deterministic

	vectors := []struct {
		input string
		seed  uint64
	}{
		{"", RapidSeed},
		{"a", RapidSeed},
		{"abc", RapidSeed},
		{"message digest", RapidSeed},
		{"abcdefghijklmnopqrstuvwxyz", RapidSeed},
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", RapidSeed},
		{"hello", 42},
	}

	for i, v := range vectors {
		h := HashWithSeed([]byte(v.input), v.seed)
		t.Logf("Vector %d: %q (seed=%d) -> %016x", i, v.input, v.seed, h)

		// Verify consistency
		hash2 := HashWithSeed([]byte(v.input), v.seed)
		if h != hash2 {
			t.Errorf("Vector %d is not consistent", i)
		}
	}
}
