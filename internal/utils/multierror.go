package utils

import (
	"fmt"
	"strings"
)

type MultiError struct {
	Errors []error
}

func (m *MultiError) Join(err error) {
	m.Errors = append(m.Errors, err)
}

func (m MultiError) IsEmpty() bool {
	return len(m.Errors) == 0
}

func (m MultiError) Error() string {
	var sb strings.Builder
	for i, err := range m.Errors {
		sb.WriteString(fmt.Sprintf("(%d) %v\n", i+1, err))
	}
	return sb.String()
}
