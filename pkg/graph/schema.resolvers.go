package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"git.proksy.io/golang/e2e/pkg/graph/generated"
	"git.proksy.io/golang/e2e/pkg/models"
)

func (r *completedTestResolver) Name(ctx context.Context, obj *models.CompletedTest) (string, error) {
	return obj.ScenarioName, nil
}

func (r *mutationResolver) RunTest(ctx context.Context, scenarios []string) (bool, error) {
	var err error
	if scenarios == nil {
		err = r.manager.RunAllTests()
	} else {
		err = r.manager.RunTests(scenarios)
	}
	return err == nil, err
}

func (r *queryResolver) AvailableScenarios(ctx context.Context) ([]string, error) {
	scenarios := r.manager.AllScenarios()
	tests := make([]string, len(scenarios))
	for i, scenario := range scenarios {
		tests[i] = scenario.Name
	}
	return tests, nil
}

func (r *queryResolver) CompletedScenarios(ctx context.Context) ([]*models.CompletedTest, error) {
	data := r.manager.CompletedScenarios()
	completed := make([]*models.CompletedTest, len(data))
	for i := range data {
		completed[i] = &data[i]
	}
	return completed, nil
}

func (r *queryResolver) LastReport(ctx context.Context) (*string, error) {
	info, err := r.manager.CurrentLaunch()
	if err != nil {
		return nil, err
	}
	link := makeGrafanaLink(r.grafana, info.ID)
	return &link, nil
}

func (r *subscriptionResolver) CurrentLaunchInfo(ctx context.Context) (<-chan *models.CompletedTest, error) {
	ch := make(chan *models.CompletedTest)
	go r.manager.SubscribeOnCompletedTests(ch)
	return ch, nil
}

// CompletedTest returns generated.CompletedTestResolver implementation.
func (r *Resolver) CompletedTest() generated.CompletedTestResolver { return &completedTestResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type completedTestResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
