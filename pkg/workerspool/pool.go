package workerspool

import (
	"context"

	"git.proksy.io/golang/e2e/common"
	"git.proksy.io/golang/e2e/pkg/models"
)

// Worker interface for processing stage
type Worker interface {
	Register(metrics common.Meter)
	Start(
		ctx context.Context,
		client string,
		tester models.Tester,
		params []models.StateSelector,
		opts *models.Options,
	) (newParams []models.StateSelector, err error)
}

type pool struct {
	metrics common.Meter
	size    int
}

type result struct {
	id       int
	newState models.StateSelector
	err      error
}

type task struct {
	id    int
	state models.StateSelector
	opt   *models.Options
}

func (p *pool) Register(metrics common.Meter) {
	p.metrics = metrics
}

func (p *pool) Start(
	ctx context.Context,
	client string,
	tester models.Tester,
	params []models.StateSelector,
	opts *models.Options,
) (newParams []models.StateSelector, err error) {
	shootCount := len(params)
	resultChan := make(chan result, shootCount)
	taskCh := make(chan task, shootCount)

	// Set worker pool not greater then request count
	workerPool := p.size
	if workerPool > shootCount {
		workerPool = shootCount
	}
	for i := 0; i < workerPool; i++ {
		go runWorker(ctx, taskCh, resultChan, client, tester)
	}

	for i := 0; i < shootCount; i++ {
		task := task{
			id:    i,
			state: params[i],
			opt:   opts,
		}
		taskCh <- task
	}

	newParams = make([]models.StateSelector, shootCount)
	errs := make([]error, 0, shootCount)
	for i := 0; i < shootCount; i++ {
		res := <-resultChan
		newParams[res.id] = res.newState
		if res.err != nil {
			errs = append(errs, res.err)
		}
	}
	// FIXME aggregate all errors
	if len(errs) > 0 {
		err = errs[0]
	}
	close(taskCh)
	return
}

// NewPool create new Worker pool with size
func NewPool(size int) Worker {
	return &pool{size: size}
}

func runWorker(ctx context.Context, taskCh <-chan task, resultCh chan result, client string, tester models.Tester) {
	for task := range taskCh {
		newState, err := tester.Run(ctx, client, task.state, task.opt)
		resultCh <- result{
			id:       task.id,
			newState: newState,
			err:      err,
		}
	}
}
