package models

import (
	"context"
)

// Tester method for making request
type Tester interface {
	MethodName() (name string)
	RequiredFields() (fields []string) // RequiredFields fields of state required for run method
	ReturnedFields() (fields []string)
	Run(ctx context.Context, client string, selector StateSelector, opts *Options) (StateSelector, error)
}
