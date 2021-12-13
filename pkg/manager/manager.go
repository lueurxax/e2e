package manager

import (
	"context"
	"sync/atomic"

	"git.proksy.io/golang/e2e/common"
	"git.proksy.io/golang/e2e/pkg/log"
	"git.proksy.io/golang/e2e/pkg/models"
)

const (
	taskPool = 10
)

type scenariosGetter interface {
	GetScenarios() ([]models.Scenario, error)
}

type processor interface {
	Run(ctx context.Context, scenario models.Scenario, launchID string) (err error)
}

// Manager manage test runs and history
type Manager interface {
	AllScenarios() (scenarios []models.Scenario)
	ScenariosByNames(names []string) (scenarios []models.Scenario, err error)
	CurrentLaunch() (info *models.LaunchInfo, err error)
	RunAllTests() (err error)
	RunTests(names []string) (err error)
	CompletedScenarios() (completed []models.CompletedTest)
	SubscribeOnCompletedTests(ch chan<- *models.CompletedTest)
	Start()
	Stop()
}

type state struct {
	history                 []models.LaunchInfo
	running                 int32
	scenarios               []models.Scenario
	scenariosIndex          map[string]int
	processor               processor
	log                     log.Logger
	taskQueue               chan task
	stopped                 chan struct{}
	completedTasks          chan completedTask
	currentLaunch           *models.LaunchInfo
	listenersCompletedTests []chan<- *models.CompletedTest
}

func (s *state) CompletedScenarios() (completed []models.CompletedTest) {
	return s.currentLaunch.CompletedTests
}

func (s *state) SubscribeOnCompletedTests(ch chan<- *models.CompletedTest) {
	s.listenersCompletedTests = append(s.listenersCompletedTests, ch)
}

func (s *state) ScenariosByNames(names []string) (scenarios []models.Scenario, err error) {
	panic("implement me")
}

func (s *state) CurrentLaunch() (info *models.LaunchInfo, err error) {
	if s.currentLaunch != nil {
		return s.currentLaunch, nil
	}
	return nil, common.ErrTestsDidNotRun()
}

func (s *state) History() []models.LaunchInfo {
	return s.history
}

func (s *state) RunAllTests() error {
	if s.IsRunning() {
		return common.ErrTestsAlreadyRunning()
	}
	s.currentLaunch = models.NewLaunchInfo()
	s.taskQueue <- task{
		launchID:  s.currentLaunch.ID,
		scenarios: s.scenarios,
	}
	return nil
}

func (s *state) RunTests(names []string) (err error) {
	if s.IsRunning() {
		return common.ErrTestsAlreadyRunning()
	}
	s.currentLaunch = models.NewLaunchInfo()
	scenarios := make([]models.Scenario, len(names))
	for i, scenarioName := range names {
		scenarioIndex, ok := s.scenariosIndex[scenarioName]
		if !ok {
			return common.ErrUnknownScenario(scenarioName)
		}
		scenarios[i] = s.scenarios[scenarioIndex]
		s.taskQueue <- task{
			launchID:  s.currentLaunch.ID,
			scenarios: scenarios,
		}
	}
	return nil
}

// Start manager
func (s *state) Start() {
	go s.loop()
	go s.broadcast()
}

func (s *state) Stop() {
	close(s.taskQueue)
	<-s.stopped
}

func (s *state) AllScenarios() (scenarios []models.Scenario) {
	return s.scenarios
}

func (s *state) IsRunning() bool {
	return atomic.LoadInt32(&s.running) == 1
}

func (s *state) broadcast() {
	for result := range s.completedTasks {
		for _, listener := range s.listenersCompletedTests {
			listener <- &result.CompletedTest
		}
		s.currentLaunch.CompletedTests = append(s.currentLaunch.CompletedTests, result.CompletedTest)
	}
}

func (s *state) loop() {
	var err error
	for task := range s.taskQueue {
		if !atomic.CompareAndSwapInt32(&s.running, 0, 1) {
			continue
		}
		for _, scenario := range task.scenarios {
			err = s.processor.Run(context.Background(), scenario, task.launchID)
			status := models.StatusCompleted
			var e string
			if err != nil {
				status = models.StatusAborted
				e = err.Error()
			}
			result := completedTask{
				launchID: task.launchID,
				CompletedTest: models.CompletedTest{
					ScenarioName: scenario.Name,
					Status:       status,
					Error:        e,
				},
			}
			s.completedTasks <- result
		}
		atomic.StoreInt32(&s.running, 0)
	}
	close(s.completedTasks)
	s.stopped <- struct{}{}
}

// New construct new manager
func New(conf scenariosGetter, proc processor, logger log.Logger) (man Manager, err error) {
	var scenarios []models.Scenario
	scenarios, err = conf.GetScenarios()
	scenariosIndex := make(map[string]int, len(scenarios))
	for i, scenario := range scenarios {
		scenariosIndex[scenario.Name] = i
	}
	if err != nil {
		return
	}
	return &state{
		processor:      proc,
		scenarios:      scenarios,
		scenariosIndex: scenariosIndex,
		log:            logger,
		taskQueue:      make(chan task, 2),
		stopped:        make(chan struct{}),
		completedTasks: make(chan completedTask, taskPool),
	}, nil
}
