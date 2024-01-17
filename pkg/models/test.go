package models

import "github.com/lueurxax/e2e/common"

// CompletedTest model
type CompletedTest struct {
	ScenarioName string `json:"scenarioName"`
	Status       Status `json:"status"`
	Error        string `json:"error"`
}

// Test config struct
type Test struct {
	Name                    string                 `yaml:"name"`
	Params                  map[string]interface{} `yaml:"params"`
	TestData                *TestData              `yaml:"testdata"`
	WaiterDelayMilliseconds int                    `yaml:"waiter_delay_milliseconds"`
	FailOnError             bool                   `yaml:"fail_on_error"`
	Repeat                  int                    `yaml:"repeat"`
	StressLoad              *common.StressLoad     `yaml:"stress_load"`
	BeforeTest              []TestStage            `yaml:"before_test"`
	Action                  TestStage              `yaml:"action"`
	Check                   TestStage              `yaml:"check"`
	AfterTest               TestStage              `yaml:"after_test"`
	InitState               common.InitState       `yaml:"-"`
}

// TestStage config for stage
type TestStage struct {
	Name   string
	Once   bool
	Client string
	Error  *string
	Params map[string]interface{}
}

// TestData contain info for find test data for tests
type TestData struct {
	BucketName string `yaml:"bucket_name"`
	Key        string `yaml:"key"`
}

// Prepare test config for use
func (t *Test) Prepare() {
	t.computeShootCount()
	t.InitState = common.InitState{
		Random:       make([]string, 0),
		Increment:    make([]string, 0),
		GlobalParams: map[string]interface{}{},
	}
	for key, param := range t.Params {
		field, ok := param.(string)
		if !ok {
			t.InitState.GlobalParams[key] = param
			continue
		}
		switch field {
		case "$random":
			t.InitState.Random = append(t.InitState.Random, key)
		case "$increment":
			t.InitState.Increment = append(t.InitState.Increment, key)
		default:
			// check unimplemented generator value
			if field[0] != '$' {
				t.InitState.GlobalParams[key] = field
			}
		}
	}
}

// compute shoot count by from, to and duration(only for line type)
//
//	    /|
//	   /*|
//	  /**|
//	 /***|
//	/****|
//	|****|to
//
// from|****|
//
//	 |****|
//	______
//	duration
func (t *Test) computeShootCount() {
	if t.StressLoad == nil {
		return
	}
	t.StressLoad.ShootCount = (t.StressLoad.From/2 + t.StressLoad.To/2) * t.StressLoad.Duration * 2
}

// Validate config
func (t *Test) Validate() error {
	if t.StressLoad != nil && t.Repeat > 1 {
		return common.ErrInvalidConfig("cannot use stress load with repeated requests")
	}
	return nil
}
