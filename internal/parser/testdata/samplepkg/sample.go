// Package samplepkg is a fixture for the docs-generator parser tests.
// It deliberately mixes exported and unexported declarations so tests
// can verify that only exported symbols are reported.
package samplepkg

// DefaultLimit is the default cap used by the examples.
const DefaultLimit = 10

// Version is the fixture package version string.
var Version = "0.1.0"

// Add returns the sum of a and b.
func Add(a, b int) int {
	return a + b
}

func helper() int {
	return 0
}

// Counter counts up from zero.
type Counter struct {
	// Count is the current value.
	Count int
	limit int
}

// Inc increases the counter by one.
func (c *Counter) Inc() {
	c.Count++
}

// Value returns the current count.
func (c Counter) Value() int {
	return c.Count
}

// Greeter can greet by name.
type Greeter interface {
	// Greet returns a greeting for name.
	Greet(name string) string
}

// StringList is a list of strings.
type StringList []string
