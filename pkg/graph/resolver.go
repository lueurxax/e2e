package graph

import (
	"fmt"

	"git.proksy.io/golang/e2e/pkg/log"
	"git.proksy.io/golang/e2e/pkg/models"
)

//go:generate go run ../scripts/gqlgen.go

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver struct, root of resolvers
type Resolver struct {
	logger  log.Logger
	manager manager
	grafana string
}

type manager interface {
	AllScenarios() (scenarios []models.Scenario)
	ScenariosByNames(names []string) (scenarios []models.Scenario, err error)
	CurrentLaunch() (info *models.LaunchInfo, err error)
	RunAllTests() (err error)
	RunTests(names []string) (err error)
	CompletedScenarios() []models.CompletedTest
	SubscribeOnCompletedTests(chan<- *models.CompletedTest)
}

// NewResolver construct new resolver
func NewResolver(
	man manager,
	logger log.Logger,
	grafana string,
) *Resolver {
	return &Resolver{logger: logger, manager: man, grafana: grafana}
}

func makeGrafanaLink(link, id string) string {
	return fmt.Sprintf("%s?var-launch_id=%s", link, id)
}
