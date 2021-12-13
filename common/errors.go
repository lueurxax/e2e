package common

import "fmt"

type errParameterRequired struct {
	parameter string
}

// Error return string error
func (e *errParameterRequired) Error() string {
	return fmt.Sprintf("parameter %s is required for this method", e.parameter)
}

// ErrParameterRequired error
func ErrParameterRequired(parameter string) (err error) {
	return &errParameterRequired{parameter: parameter}
}

type errParameterHasIncorrectType struct {
	name         string
	expectedType string
}

// Error return string error
func (e *errParameterHasIncorrectType) Error() string {
	return fmt.Sprintf("parameter  %s has incorrect type, expected %s", e.name, e.expectedType)
}

// ErrParameterHasIncorrectType error
func ErrParameterHasIncorrectType(name, expectedType string) (err error) {
	return &errParameterHasIncorrectType{name: name, expectedType: expectedType}
}

// ErrObjectsDidntEqual error
type ErrObjectsDidntEqual struct {
}

// Error return string error
func (e *ErrObjectsDidntEqual) Error() string {
	return "test data and result objects are different"
}

type errIncorrectObjectCounts struct {
	expected int
	actual   int
}

// Error return string error
func (e *errIncorrectObjectCounts) Error() string {
	return fmt.Sprintf("incorrect object counts expected %d, actual %d", e.expected, e.actual)
}

// ErrIncorrectObjectCounts error
func ErrIncorrectObjectCounts(expected, actual int) error {
	return &errIncorrectObjectCounts{
		expected: expected,
		actual:   actual,
	}
}

type errIncorrectObjectKey struct {
	expected string
	actual   string
}

// Error return string error
func (e *errIncorrectObjectKey) Error() string {
	return fmt.Sprintf("incorrect object key expected %s, actual %s", e.expected, e.actual)
}

// ErrIncorrectObjectKey error
func ErrIncorrectObjectKey(expected, actual string) error {
	return &errIncorrectObjectKey{expected: expected, actual: actual}
}

// ErrTestsAlreadyRunning error
func ErrTestsAlreadyRunning() error {
	return fmt.Errorf("tests already running")
}

type errUnknownScenario struct {
	name string
}

// Error return string error
func (e errUnknownScenario) Error() string {
	return fmt.Sprintf("unknown scenario %s", e.name)
}

// ErrUnknownScenario error
func ErrUnknownScenario(name string) error {
	return &errUnknownScenario{name: name}
}

type errTestMethodNotImplemented struct {
	name string
}

// ErrTestMethodNotImplemented error
func ErrTestMethodNotImplemented(name string) error {
	return &errTestMethodNotImplemented{name: name}
}

// Error return error string
func (e *errTestMethodNotImplemented) Error() string {
	return fmt.Sprintf("method %s not implemented", e.name)
}

type errUnknownMethod struct {
	name string
}

// ErrUnknownMethod error
func ErrUnknownMethod(name string) error {
	return &errUnknownMethod{name: name}
}

// Error return error string
func (e *errUnknownMethod) Error() string {
	return fmt.Sprintf("method %s didn't found", e.name)
}

type errTestsDidNotRun struct {
}

// ErrTestsDidNotRun error
func ErrTestsDidNotRun() error {
	return &errTestsDidNotRun{}
}

// Error return error string
func (e *errTestsDidNotRun) Error() string {
	return "tests did not run"
}

type errInvalidConfig struct {
	reason string
}

// ErrInvalidConfig error
func ErrInvalidConfig(reason string) error {
	return &errInvalidConfig{reason: reason}
}

// Error return error string
func (e *errInvalidConfig) Error() string {
	return fmt.Sprintf("config is invalid, reason: %s", e.reason)
}

type errRequiredFieldDidntSet struct {
	method, scenario, requiredField string
}

// ErrRequiredFieldDidntSet error
func ErrRequiredFieldDidntSet(method, scenario, requiredField string) error {
	return &errRequiredFieldDidntSet{method: method, scenario: scenario, requiredField: requiredField}
}

// Error return error string
func (e *errRequiredFieldDidntSet) Error() string {
	return fmt.Sprintf(
		"config is invalid, required field %s, didn't found in scenario %s, method %s",
		e.requiredField, e.scenario, e.method)
}
