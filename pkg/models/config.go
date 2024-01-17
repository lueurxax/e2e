package models

import "github.com/lueurxax/e2e/common"

// Config of config structure
type Config struct {
	Clients []common.Client `yaml:"clients"`
	Tests   []Test          `yaml:"tests"`
}
