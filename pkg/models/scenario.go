package models

// Scenario content parameters for run test scenario
type Scenario struct {
	Name       string
	BeforeTest []Stage
	Action     *Stage
	Check      *Stage
	AfterTest  *Stage
	Config     *Test
}

// ScenarioMeta scenario meta information
type ScenarioMeta struct {
	ID   string
	Name string
}
