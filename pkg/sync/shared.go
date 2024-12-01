package sync

import (
	"fmt"
	"strings"
)

// Used to capture future result in a channel for async processing
type resultWithIndex[T any, E error] struct {
	result T
	err    *E
	index  int
}

// Errors is used in place of an array of errors. A future can't return an array of errors, due to the type definitions.
// Instead we use this type
type Errors[E error] struct {
	Errors []*E
}

// Error converts an array of errors to a string. This therefore implements the error interface
func (e Errors[E]) Error() string {
	builder := strings.Builder{}
	for i, err := range e.Errors {
		if err == nil {
			builder.WriteString(fmt.Sprintf("%d: nil\n", i))
		} else {
			builder.WriteString(fmt.Sprintf("%d: %s\n", i, (*err).Error()))
		}
	}

	return builder.String()
}
