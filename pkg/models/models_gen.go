// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"fmt"
	"io"
	"strconv"
)

type Mutation struct {
}

type Query struct {
}

type Subscription struct {
}

type Status string

const (
	StatusCompleted Status = "COMPLETED"
	StatusAborted   Status = "ABORTED"
	StatusRunning   Status = "RUNNING"
)

var AllStatus = []Status{
	StatusCompleted,
	StatusAborted,
	StatusRunning,
}

func (e Status) IsValid() bool {
	switch e {
	case StatusCompleted, StatusAborted, StatusRunning:
		return true
	}
	return false
}

func (e Status) String() string {
	return string(e)
}

func (e *Status) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Status(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

func (e Status) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
