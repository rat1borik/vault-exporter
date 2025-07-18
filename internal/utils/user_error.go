package utils

import (
	"errors"
	"fmt"
)

type UserError struct {
	Message string
}

func (e *UserError) Error() string {
	return fmt.Sprintf("Ошибка импорта: %s", e.Message)
}

func UserErrorf(format string, a ...any) *UserError {
	return &UserError{
		Message: fmt.Sprintf(format, a...),
	}
}

type userErrorCollection struct {
	collection []error
}

func NewUserErrorCollection() *userErrorCollection {
	return &userErrorCollection{
		collection: make([]error, 0),
	}
}

func (col *userErrorCollection) Add(e error, debug bool) {
	var dummy *UserError
	if errors.As(e, &dummy) || debug {
		col.collection = append(col.collection, e)
	}
}

func (col *userErrorCollection) Collection() []error {
	if len(col.collection) > 0 {
		return col.collection
	}
	return nil
}
