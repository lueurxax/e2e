package models

import "git.proksy.io/golang/e2e/common"

// Stage of test scenario
type Stage struct {
	Tester        string
	Client        string
	WantError     bool
	Error         string
	Params        []map[string]interface{}
	RequestsCount int
}

type State interface {
	Reset(count int) []StateSelector
	Prepare(state common.InitState)
	MergeToState(states []StateSelector)
	MergeToStateRepeat(selector StateSelector)
	AddToState(selectors []StateSelector, params ...map[string]interface{}) []StateSelector
}

type StateSelector interface {
	Index() int
}
