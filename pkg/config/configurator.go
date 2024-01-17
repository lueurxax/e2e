package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/lueurxax/e2e/common"
	"github.com/lueurxax/e2e/pkg/models"
)

// Configurator read and init config for work with e2e framework
type Configurator interface {
	Read() (err error)
	Validate() (err error)
	Init()
	GetScenarios() (scenarios []models.Scenario, err error)
	Clients() (clients []common.Client)
}

type config struct {
	configPath string
	data       *models.Config
}

func (c *config) Clients() []common.Client {
	return c.data.Clients
}

// Read config for yaml file
func (c *config) Read() (err error) {
	var absPath string
	absPath, err = filepath.Abs(c.configPath)
	if err != nil {
		return
	}
	var data []byte
	data, err = ioutil.ReadFile(absPath)
	if err != nil {
		return
	}
	c.data = &models.Config{}
	err = yaml.Unmarshal(data, c.data)
	return
}

// Init config
func (c *config) Init() {
	for i := range c.data.Tests {
		c.data.Tests[i].Prepare()
	}
}

// Validate config
func (c *config) Validate() (err error) {
	for i := range c.data.Tests {
		if err = c.data.Tests[i].Validate(); err != nil {
			return
		}
	}
	return
}

// GetScenarios scenarios names list
func (c *config) GetScenarios() (scenarios []models.Scenario, err error) {
	scenarios = make([]models.Scenario, len(c.data.Tests))
	for i, scenario := range c.data.Tests {
		data := models.Scenario{
			Name:       scenario.Name,
			Config:     &c.data.Tests[i],
			BeforeTest: make([]models.Stage, len(scenario.BeforeTest)),
		}
		var shootCount int
		if scenario.StressLoad != nil {
			shootCount = scenario.StressLoad.GetShootCount()
		} else {
			shootCount = scenario.Repeat
		}
		for i, stageConf := range scenario.BeforeTest {
			stageData := getStage(scenario.InitState.GlobalParams, stageConf, shootCount)
			data.BeforeTest[i] = *stageData
		}

		data.Action = getStage(scenario.InitState.GlobalParams, scenario.Action, shootCount)
		data.Check = getStage(scenario.InitState.GlobalParams, scenario.Check, shootCount)
		data.AfterTest = getStage(scenario.InitState.GlobalParams, scenario.AfterTest, shootCount)
		scenarios[i] = data
	}
	return
}

func getStage(globalParams map[string]interface{}, conf models.TestStage, stressLoadCount int) (st *models.Stage) {
	requestCount := stressLoadCount
	if conf.Once {
		requestCount = 1
	}
	st = &models.Stage{
		Tester:        conf.Name,
		Client:        conf.Client,
		Params:        []map[string]interface{}{globalParams, conf.Params},
		RequestsCount: requestCount,
	}
	if conf.Error != nil {
		st.Error = *conf.Error
		st.WantError = true
	}
	return
}

// NewConfig construct Configurator
func NewConfig(configPath string) Configurator {
	return &config{
		configPath: configPath,
	}
}
