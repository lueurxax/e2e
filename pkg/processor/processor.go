package processor

import (
	"context"

	"github.com/pkg/errors"

	"github.com/lueurxax/e2e/common"
	"github.com/lueurxax/e2e/pkg/log"
	"github.com/lueurxax/e2e/pkg/models"
	"github.com/lueurxax/e2e/pkg/pandoraconnector"
	"github.com/lueurxax/e2e/pkg/testerspool"
	"github.com/lueurxax/e2e/pkg/workerspool"
)

// Processor interface process tests
type Processor interface {
	ValidateScenario(scenario models.Scenario) error
	Run(ctx context.Context, scenario models.Scenario, launchID string) (err error)
}

type worker interface {
	Register(metrics common.Meter)
	Start(
		ctx context.Context,
		client string,
		tester models.Tester,
		params []models.StateSelector,
		opts *models.Options,
	) (newParams []models.StateSelector, err error)
}

type processor struct {
	state models.State

	stageProcessor StageProcessor
	logger         log.Logger
	metrics        common.Meter
}

// ValidateScenario map scenario params and validate current scenario params
func (p *processor) ValidateScenario(scenario models.Scenario) (err error) {
	globalFields := make([]string, 0)
	var fields []string
	for key := range scenario.Config.Params {
		globalFields = append(globalFields, key)
	}

	// validate before test params
	for i := range scenario.Config.BeforeTest {
		fields, err = p.stageProcessor.Validate(globalFields, &scenario.Config.BeforeTest[i], scenario.Name)
		if err != nil {
			return
		}
		globalFields = append(globalFields, fields...)
	}

	// validate action params
	fields, err = p.stageProcessor.Validate(globalFields, &scenario.Config.Action, scenario.Name)
	if err != nil {
		return
	}
	globalFields = append(globalFields, fields...)

	// validate check params
	fields, err = p.stageProcessor.Validate(globalFields, &scenario.Config.Check, scenario.Name)
	if err != nil {
		return
	}
	globalFields = append(globalFields, fields...)

	// validate after test params
	_, err = p.stageProcessor.Validate(globalFields, &scenario.Config.AfterTest, scenario.Name)
	return
}

// Run test scenario
func (p *processor) Run(ctx context.Context, scenario models.Scenario, launchID string) (err error) {
	p.logger.WithField("scenario", scenario.Name).WithField("id", launchID).Info("run")
	p.metrics.NewLaunch(launchID)
	stressLoad := false
	shootCount := 1
	if scenario.Config.StressLoad != nil {
		p.metrics.NewStress(scenario.Name, *scenario.Config.StressLoad)
		stressLoad = true
		shootCount = scenario.Config.StressLoad.ShootCount
	}
	if scenario.Config.Repeat > 1 {
		shootCount = scenario.Config.Repeat
	}

	// init scenario state
	selectors := p.state.Reset(shootCount)
	p.state.Prepare(scenario.Config.InitState)

	// clean instance after tests
	defer func(scenario models.Scenario, metr common.Meter) {
		_, err2 := p.stageProcessor.Run(ctx, scenario.AfterTest, selectors, scenario.Config, false)
		metr.Reset()
		if err2 != nil {
			p.logger.WithField("scenario", scenario.Name).WithError(err2).
				Warn("failed run after test for scenario")
			return
		}
	}(scenario, p.metrics)

	// run before test
	for i, stage := range scenario.BeforeTest {
		newStates, err := p.stageProcessor.Run(ctx, &scenario.BeforeTest[i], selectors, scenario.Config, false)
		if err != nil {
			p.logger.WithField("scenario", scenario.Name).WithError(err).
				Warn("failed run before test for scenario")
		}

		// merge state if it possible
		if stage.RequestsCount == shootCount {
			p.state.MergeToState(newStates)
		}
		if stage.RequestsCount == 1 {
			p.state.MergeToStateRepeat(newStates[0])
		}
	}

	// run test actions
	var newSelectors []models.StateSelector
	newSelectors, err = p.stageProcessor.Run(ctx, scenario.Action, selectors, scenario.Config, stressLoad)
	if err != nil {
		p.logger.WithField("scenario", scenario.Name).WithError(err).
			Warn("failed run action test for scenario")
		err = errors.Wrap(err, "on action")
		return
	}

	// merge state if it is possible
	if scenario.Action.RequestsCount == shootCount && shootCount == len(newSelectors) {
		p.state.MergeToState(newSelectors)
	}

	// run checks of this test
	newSelectors, err = p.stageProcessor.Run(ctx, scenario.Check, newSelectors, scenario.Config, false)
	if (err != nil && !scenario.Check.WantError && err.Error() != scenario.Check.Error) ||
		(err == nil && scenario.Check.WantError) {
		// FIXME write real error
		p.logger.WithField("scenario", scenario.Name).WithError(err).
			Warn("failed run check for scenario")
		return errors.Wrap(err, "on check")
	}
	// merge state if it possible
	if scenario.Check.RequestsCount == shootCount {
		p.state.MergeToState(newSelectors)
	}

	return nil
}

// NewProcessor construct new Processor
func NewProcessor(
	state models.State,
	testers testerspool.TestersPool,
	l log.Logger,
	metrics common.Meter,
	workerPoolSize int,
) (proc Processor, err error) {
	pandora := pandoraconnector.NewConnector(l.WithField("receiver", "pandora"))

	pandora.Register(metrics)

	proc = &processor{
		state: state,
		stageProcessor: newStageProcessor(
			state, testers, pandora, workerspool.NewPool(workerPoolSize), metrics, l,
		),
		logger:  l,
		metrics: metrics,
	}
	return
}
