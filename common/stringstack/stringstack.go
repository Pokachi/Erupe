package stringstack

import (
	"errors"
)

// StringStack is a basic LIFO "stack" for storing strings.
type StringStack struct {
	stack []string
}

// New creates a new instance of StringStack
func New() *StringStack {
	return &StringStack{}
}

// Set sets up a new StringStack
func (s *StringStack) Set(v string) {
	s.stack = []string{v}
}

// Push pushes a string onto the stack.
func (s *StringStack) Push(v string) {
	s.stack = append(s.stack, v)
}

// Pop pops a string from the stack.
func (s *StringStack) Pop() (string, error) {
	if len(s.stack) == 0 {
		return "", errors.New("no items on stack")
	}

	x := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]

	return x, nil
}
