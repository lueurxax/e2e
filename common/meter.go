package common

// Meter metering test scenario
type Meter interface {
	NewLaunch(launchID string)
	NewStress(scenarioName string, conf StressLoad)
	AddRequest(data *RequestData)
	Reset()
}
