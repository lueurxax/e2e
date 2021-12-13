package pandoraconnector

import (
	"git.proksy.io/golang/e2e/pkg/models"
)

// Ammo test params for shoot
type Ammo struct {
	Conf   *models.Test
	Params models.StateSelector
}
