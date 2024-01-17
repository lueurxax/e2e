package manager

import (
	"github.com/lueurxax/e2e/pkg/models"
)

type task struct {
	launchID  string
	scenarios []models.Scenario
}

type completedTask struct {
	launchID string
	models.CompletedTest
}
