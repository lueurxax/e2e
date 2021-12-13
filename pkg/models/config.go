package models

import "git.proksy.io/golang/e2e/common"

// Config of config structure
type Config struct {
	Clients []common.Client `yaml:"clients"`
	Tests   []Test          `yaml:"tests"`
}
