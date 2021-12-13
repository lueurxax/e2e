package models

import (
	uuid "github.com/satori/go.uuid"
)

// ErrorTest contain errors of completed tests
type ErrorTest struct {
	ScenarioName string
	Error        string
}

// LaunchInfo info about launch
type LaunchInfo struct {
	ID string
	Status
	CompletedTests []CompletedTest
	Errors         []ErrorTest
}

// NewLaunchInfo construct LaunchInfo struct
func NewLaunchInfo() *LaunchInfo {
	launchID := uuid.NewV4()
	return &LaunchInfo{
		ID:             launchID.String(),
		Status:         StatusRunning,
		CompletedTests: make([]CompletedTest, 0),
	}
}
