package processor

import (
	"context"

	"git.proksy.io/golang/e2e/common"
	"git.proksy.io/golang/e2e/pkg/log"
	"git.proksy.io/golang/e2e/pkg/models"
	"git.proksy.io/golang/e2e/pkg/testerspool"
)

// StageProcessor interface process test stage
type StageProcessor interface {
	Validate(globalFields []string, stage *models.TestStage, name string) (resultFields []string, err error)
	Run(
		ctx context.Context,
		test *models.Stage, state []models.StateSelector, config *models.Test, stressLoad bool,
	) (newParams []models.StateSelector, err error)
}

type stageProcessor struct {
	state      models.State
	testers    testerspool.TestersPool
	pandora    worker
	workerPool worker
	metrics    common.Meter
	logger     log.Logger
}

func (s *stageProcessor) Validate(
	globalFields []string,
	stage *models.TestStage,
	scenarioName string,
) (resultFields []string, err error) {
	fields := make([]string, len(globalFields))
	copy(fields, globalFields)
	for key := range stage.Params {
		fields = append(fields, key)
	}
	if err = s.validateFields(stage, fields, scenarioName); err != nil {
		return
	}
	resultFields, err = s.getReturnedFields(stage)
	return
}

func (s *stageProcessor) Run(
	ctx context.Context,
	stage *models.Stage,
	scenarioState []models.StateSelector,
	config *models.Test,
	stressLoad bool,
) (newParams []models.StateSelector, err error) {
	state := s.state.AddToState(scenarioState, stage.Params...)
	opts := &models.Options{
		Conf: config,
	}

	var tester models.Tester
	tester, err = s.testers.Get(stage.Tester)
	if err != nil {
		return
	}
	if stressLoad {
		if newParams, err = s.pandora.Start(
			ctx,
			stage.Client,
			tester,
			state,
			opts); err != nil {
			return
		}
	} else {
		newParams, err = s.workerPool.Start(
			ctx,
			stage.Client,
			tester,
			state[:stage.RequestsCount],
			opts,
		)
		if err != nil {
			return
		}
	}
	return
}

func (s *stageProcessor) getReturnedFields(stage *models.TestStage) (fields []string, err error) {
	tester, err := s.testers.Get(stage.Name)
	if err != nil {
		return
	}
	return tester.ReturnedFields(), nil
}

func (s *stageProcessor) validateFields(stage *models.TestStage, params []string, scenarioName string) (err error) {
	var tester models.Tester
	tester, err = s.testers.Get(stage.Name)
	if err != nil {
		return
	}
	paramsMap := map[string]struct{}{}
	for _, param := range params {
		paramsMap[param] = struct{}{}
	}
	for _, requiredKey := range tester.RequiredFields() {
		if _, ok := paramsMap[requiredKey]; !ok {
			return common.ErrRequiredFieldDidntSet(stage.Name, scenarioName, requiredKey)
		}
	}
	return
}

func newStageProcessor(
	state models.State,
	testers testerspool.TestersPool,
	pandora worker, workerPool worker,
	metrics common.Meter,
	logger log.Logger,
) StageProcessor {
	return &stageProcessor{
		state:      state,
		testers:    testers,
		pandora:    pandora,
		workerPool: workerPool,
		metrics:    metrics,
		logger:     logger,
	}
}
