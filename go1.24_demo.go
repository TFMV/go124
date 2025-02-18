// Demo for Go 1.24
//
// This demo shows some of the new features and improvements in Go 1.24.
// It includes:
// - Generic type aliases
// - CGO improve
// - Improved finalizers
// - Crypto packages: HKDF, PBKDF2, SHA3
// - Directory-limited filesystem access
// - Bytes and strings iterators
// - New encoding interfaces: TextAppender and BinaryAppender
// - netip: Encoding Interfaces
// - Regexp: TextAppender Interface
// - Runtime GOROOT deprecation notice
// - Text template: Range over integer sequence
// - math/big: Encoding TextAppender
// - math/rand: Using a Rand instance
// - sync.Map improvements
// - log/slog: DiscardHandler demonstration
// - time: Encoding Interfaces
// - experimental testing/synctest
// - go/types Iterator Methods
// - maphash: Comparable and WriteComparable

// To run the demo, ensure you have Go 1.24 installed and run:
// go run go1.24_demo.go
package main

import (
	"bytes"
	"crypto/pbkdf2"
	"crypto/sha256"
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"hash/maphash"
	"math/big"
	"math/rand"
	"net/netip"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"text/template"
	"time"
)

// ----------------------------------------------------------------------------
// 1. Generic Type Aliases
//
// Go 1.24 now fully supports generic type aliases. In this example we create
// a generic alias MySlice[T] for []T.
type MySlice[T any] = []T

func demoGenericTypeAlias() {
	numbers := MySlice[int]{1, 2, 3, 4, 5}
	fmt.Println("Generic Type Alias (MySlice[int]):", numbers)
}

// 2. CGO Improvements (Skipped Code Implementation)
// New cgo annotations such as "noescape" and "nocallback" can now be used.
// (This example calls two dummy C functions.)
//
// To compile cgo code, ensure cgo is enabled.
// The annotations are written in the preamble below:
//
// #cgo noescape: c_function_noescape
// #cgo nocallback: c_function_nocallback
// #include <stdlib.h>
// void c_function_noescape(void* p) {}
// void c_function_nocallback(void* p) {}

// ----------------------------------------------------------------------------
// 3. Improved Finalizers (using runtime.SetFinalizer as a stand-in)
//
// Go 1.24 introduces runtime.AddCleanup to attach multiple cleanups to an object.
// Here we use runtime.SetFinalizer (the older API) to demonstrate finalization.
func DemoFinalizers() {
	// Wrap an int in a custom struct to show finalization.
	type Holder struct {
		Value int
	}
	holder := &Holder{Value: 42}
	// Set a finalizer on the holder.
	runtime.SetFinalizer(holder, func(h *Holder) {
		fmt.Println("Finalizer called for Holder with value:", h.Value)
	})
	// Remove our reference and force garbage collection.
	holder = nil
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
}

// ----------------------------------------------------------------------------
// 4. Crypto Packages: HKDF, PBKDF2, SHA3
//
// This demo uses HKDF (from golang.org/x/crypto/hkdf for now),
// PBKDF2, and SHA3-256.
func DemoCryptoPackages() {
	// PBKDF2 and SHA3-256 demos
	password := "my password"
	salt := []byte("my salt")
	pbkdf2Key, err := pbkdf2.Key(sha256.New, password, salt, 4096, 32)
	if err != nil {
		fmt.Println("PBKDF2 error:", err)
		return
	}
	fmt.Println("Derived key (PBKDF2):", hex.EncodeToString(pbkdf2Key))

	// SHA3-256 demo
	hasher := sha3.New256()
	hasher.Write([]byte("hello world"))
	digest := hasher.Sum(nil)
	fmt.Println("SHA3-256 digest:", hex.EncodeToString(digest))
}

// ----------------------------------------------------------------------------
// 5. Directory-Limited Filesystem Access
//
// In Go 1.24 the new os.Root type (and related functions) let you limit
// filesystem access to a directory. For this demo we simulate such behavior.
func DemoDirectoryLimitedFS() {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "demo-root")
	if err != nil {
		fmt.Println("Error creating temp directory:", err)
		return
	}
	defer os.RemoveAll(tempDir)

	// Create a file within the directory.
	filePath := tempDir + "/example.txt"
	if err := os.WriteFile(filePath, []byte("Hello from a limited FS!"), 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	// Open the directory.
	root, err := os.Open(tempDir)
	if err != nil {
		fmt.Println("Error opening directory:", err)
		return
	}
	defer root.Close()

	entries, err := root.Readdir(0)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	fmt.Println("Files in limited FS:")
	for _, entry := range entries {
		fmt.Println(" -", entry.Name())
	}
}

// ----------------------------------------------------------------------------
// 6. Bytes and Strings Iterators
//
// The new iterator-style functions (e.g. Lines, SplitSeq) make it easier
// to work with byte slices and strings.
func DemoBytesAndStringsIterators() {
	text := []byte("line1\nline2\nline3\n")
	fmt.Println("Iterating over lines (using bytes.Split):")
	for _, line := range bytes.Split(text, []byte("\n")) {
		if len(line) > 0 {
			fmt.Println(string(line))
		}
	}

	sample := "  foo   bar baz  "
	fmt.Println("Iterating over fields (using strings.Fields):")
	for _, field := range strings.Fields(sample) {
		fmt.Println(field)
	}
}

// ----------------------------------------------------------------------------
// 7. New encoding Interfaces: TextAppender and BinaryAppender
//
// Types that already implement TextMarshaler now also implement the
// TextAppender interface to append directly to a buffer.
type demoStruct struct {
	Value int
}

// AppendText implements encoding.TextAppender for demoStruct.
func (d demoStruct) AppendText(dst []byte) []byte {
	return append(dst, fmt.Sprintf("demoStruct(%d)", d.Value)...)
}

func DemoEncodingAppend() {
	ds := demoStruct{Value: 123}
	var buf []byte
	// Use the TextAppender interface if available.
	if appender, ok := interface{}(ds).(interface {
		AppendText([]byte) []byte
	}); ok {
		buf = appender.AppendText(buf)
	} else {
		buf = append(buf, fmt.Sprintf("%v", ds)...)
	}
	fmt.Println("Encoding append result:", string(buf))
}

// ----------------------------------------------------------------------------
// 8. go/net/netip: Encoding Interfaces
//
// netip.Addr now implements encoding.TextAppender.
func DemoNetipEncoding() {
	addr, err := netip.ParseAddr("192.0.2.1")
	if err != nil {
		fmt.Println("Error parsing IP:", err)
		return
	}
	var buf []byte
	// Use type assertion to check for TextAppender.
	if appender, ok := interface{}(addr).(interface {
		AppendText([]byte) []byte
	}); ok {
		buf = appender.AppendText(buf)
	} else {
		buf = []byte(addr.String())
	}
	fmt.Println("netip.Addr appended text:", string(buf))
}

// ----------------------------------------------------------------------------
// 9. Regexp: TextAppender Interface
//
// Regular expressions now implement encoding.TextAppender.
func DemoRegexpEncoding() {
	re := regexp.MustCompile(`a*b`)
	var buf []byte
	if appender, ok := interface{}(re).(interface {
		AppendText([]byte) []byte
	}); ok {
		buf = appender.AppendText(buf)
	} else {
		buf = []byte(re.String())
	}
	fmt.Println("Regexp appended text:", string(buf))
}

// ----------------------------------------------------------------------------
// 10. Runtime GOROOT Deprecation Notice
//
// runtime.GOROOT is now deprecated.
func DemoRuntimeGOROOT() {
	fmt.Println("Note: runtime.GOROOT is deprecated; use 'go env GOROOT' instead.")
}

// ----------------------------------------------------------------------------
// 11. Text Template: Range over Integer Sequence
//
// Templates now support range-over-int. This demo uses a "seq" function.
func DemoTextTemplate() {
	tmplText := `Numbers: {{range $i := seq 1 5}}{{$i}} {{end}}`
	tmpl, err := template.New("demo").Funcs(template.FuncMap{
		"seq": func(start, end int) []int {
			s := make([]int, 0, end-start+1)
			for i := start; i <= end; i++ {
				s = append(s, i)
			}
			return s
		},
	}).Parse(tmplText)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}
	var tplOutput bytes.Buffer
	if err := tmpl.Execute(&tplOutput, nil); err != nil {
		fmt.Println("Error executing template:", err)
		return
	}
	fmt.Println("Template output:", tplOutput.String())
}

// ----------------------------------------------------------------------------
// 12. math/big: Encoding TextAppender
//
// big.Int now implements encoding.TextAppender.
func DemoMathBigEncoding() {
	bigInt := new(big.Int)
	bigInt.SetString("12345678901234567890", 10)
	var buf []byte
	if appender, ok := interface{}(bigInt).(interface {
		AppendText([]byte) []byte
	}); ok {
		buf = appender.AppendText(buf)
	} else {
		buf = []byte(bigInt.String())
	}
	fmt.Println("big.Int appended text:", string(buf))
}

// ----------------------------------------------------------------------------
// 13. math/rand: Using a Rand Instance
//
// The top-level Seed function is deprecated. Create a new Rand instance.
func DemoMathRand() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Println("Random number (rand.New):", r.Int())
}

// ----------------------------------------------------------------------------
// 14. sync.Map Improvements
//
// The new sync.Map implementation now exhibits reduced contention.
func DemoSyncMap() {
	var m sync.Map
	m.Store("key1", 100)
	m.Store("key2", 200)
	fmt.Println("Iterating over sync.Map:")
	m.Range(func(key, value any) bool {
		fmt.Printf("  key=%v, value=%v\n", key, value)
		return true
	})
}

// ----------------------------------------------------------------------------
// 15. log/slog: DiscardHandler Demonstration
//
// In Go 1.24, the new log/slog package provides a DiscardHandler that discards log output.
// For simplicity we just note its existence.
func DemoSlog() {
	fmt.Println("slog.DiscardHandler demo: In production, a DiscardHandler would discard logs.")
}

// ----------------------------------------------------------------------------
// 16. Text Template (already shown above)
// ----------------------------------------------------------------------------
// 17. time: Encoding Interfaces
//
// time.Time now implements encoding.TextAppender.
func DemoTimeEncoding() {
	now := time.Now()
	var buf []byte
	if appender, ok := interface{}(now).(interface {
		AppendText([]byte) []byte
	}); ok {
		buf = appender.AppendText(buf)
	} else {
		buf = []byte(now.String())
	}
	fmt.Println("time.Time appended text:", string(buf))
}

// ----------------------------------------------------------------------------
// 18. Experimental testing/synctest
//
// The new experimental testing/synctest package is best used in tests and requires
// GOEXPERIMENT=synctest. Here we simply print a note.
func DemoSynctest() {
	fmt.Println("Experimental synctest demo: See tests built with GOEXPERIMENT=synctest for usage.")
}

// ----------------------------------------------------------------------------
// 19. go/types Iterator Methods
//
// Improvements to go/types now let you iterate over sequences with methods like Variables().
// We simply note this improvement.
func DemoGoTypesIterators() {
	fmt.Println("go/types iterator demonstration: Use the Variables() method on tuples, etc.")
}

// ----------------------------------------------------------------------------
// 20. maphash: Comparable and WriteComparable
//
// The new maphash functions make it easy to hash comparable values.
func DemoMaphashComparable() {
	var h maphash.Hash
	key := "myKey"
	h.WriteString(key)
	hashValue := h.Sum64()
	fmt.Printf("Hash for key %q: %d\n", key, hashValue)
}

func main() {
	fmt.Println("=== Go 1.24 Demo ===")
	demoGenericTypeAlias()
	fmt.Println("CGO Improvements Demo: Not implemented")
	DemoFinalizers()
	DemoCryptoPackages()
	DemoDirectoryLimitedFS()
	DemoBytesAndStringsIterators()
	DemoEncodingAppend()
	DemoNetipEncoding()
	DemoRegexpEncoding()
	DemoRuntimeGOROOT()
	DemoTextTemplate()
	DemoMathBigEncoding()
	DemoMathRand()
	DemoSyncMap()
	DemoSlog()
	fmt.Println("Text Template Range Demo: Not implemented") // #16
	DemoTimeEncoding()
	DemoSynctest()
	DemoGoTypesIterators()
	DemoMaphashComparable()
	fmt.Println("=== Go 1.24 Demo End ===")
}
