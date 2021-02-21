package saga

import "fmt"

// Logger interface is used throughout the saga code in this repo.
type Logger interface {
	Info(value ...interface{})
	Error(value ...interface{})
}

// SimpleLogger is as good as using fmt to print on console.
type SimpleLogger struct{}

// NewSimpleLogger creates a new simple logger object.
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{}
}

// Info formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It ignores any error while writing string to io.
func (sl *SimpleLogger) Info(args ...interface{}) {
	str := fmt.Sprint(args...)
	fmt.Print(str)
}

// Error formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It ignores any error while writing string to io.
func (sl *SimpleLogger) Error(args ...interface{}) {
	str := fmt.Sprint(args...)
	fmt.Print(str)
}

// DummyLogger just satisfies the saga.Logger interface and does nothin.
// This is the default Logger used in saga.
type DummyLogger struct{}

// NewDummyLogger returns new DummyLogger object.
func NewDummyLogger() *DummyLogger {
	return &DummyLogger{}
}

// Info does nothing.
func (dl *DummyLogger) Info(args ...interface{}) {}

// Error does nothing.
func (dl *DummyLogger) Error(args ...interface{}) {}
