# package samplepkg
Package samplepkg is a fixture for the docs-generator parser tests.
It deliberately mixes exported and unexported declarations so tests
can verify that only exported symbols are reported.

## Index
- [func Add](#func-add)
- [type Counter](#type-counter)
- [type Greeter](#type-greeter)
- [type StringList](#type-stringlist)
- [const DefaultLimit](#const-defaultlimit)
- [var Version](#var-version)

## Functions
### func Add
```go
func Add(a, b int) int
```
Add returns the sum of a and b.

## Types
### type Counter
```go
type Counter struct {
	// Count is the current value.
	Count int
	// contains filtered or unexported fields
}
```
Counter counts up from zero.

#### Fields
| Name | Type | Doc |
| --- | --- | --- |
| Count | `int` | Count is the current value. |

- method `func (*Counter) Inc()` — Inc increases the counter by one.
- method `func (Counter) Value() int` — Value returns the current count.

### type Greeter
```go
type Greeter interface {
	// Greet returns a greeting for name.
	Greet(name string) string
}
```
Greeter can greet by name.

### type StringList
```go
type StringList []string
```
StringList is a list of strings.

## Constants
### const DefaultLimit
```go
const DefaultLimit = 10
```
DefaultLimit is the default cap used by the examples.

## Variables
### var Version
```go
var Version = "0.1.0"
```
Version is the fixture package version string.

