package common

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

func Logger() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: false,
	}).With().
		Timestamp().
		Logger()
}

// Workflow input
type WorkflowInput struct {
	Name string
}

type WorkflowOutput struct {
	Value string
}

// Activity input
type ActivityInput struct {
	Name string
}

type ActivityOutput struct {
	Value string
}

// Workflow implementation using activity and parallel activities
func Workflow(ctx workflow.Context, wInput *WorkflowInput) (*WorkflowOutput, error) {

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("workflow started", "name", wInput.Name)

	aInput := &ActivityInput{
		Name: wInput.Name,
	}

	logger.Info("Activity started", "input", aInput)
	we := workflow.ExecuteActivity(ctx, Activity, aInput)

	var aResult *ActivityOutput
	// block until activity completes - the result will be available
	err := we.Get(ctx, &aResult)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return &WorkflowOutput{}, err
	}

	newActivityInput := &ActivityInput{
		Name: aResult.Value,
	}
	// launch several activities in parallel
	var futures []workflow.Future
	for i := 0; i < N; i++ {
		future := workflow.ExecuteActivity(ctx, ActivityParallel, newActivityInput)
		futures = append(futures, future)
	}
	logger.Info("Parallel Activities started", "N", N)

	// accumulate results
	var results []ActivityOutput
	for _, future := range futures {
		var aOut ActivityOutput
		err = future.Get(ctx, &aOut)
		logger.Info("Activity returned with result", "result", aOut)
		if err != nil {
			logger.Error("Activity Parallel failed", "Error", err)
			return &WorkflowOutput{}, err
		}
		results = append(results, aOut)
	}

	v := ""
	for _, r := range results {
		v += "+" + r.Value
	}

	logger.Info("workflow completed", "result", v)

	return &WorkflowOutput{Value: v}, nil
}

func Activity(ctx context.Context, input *ActivityInput) (*ActivityOutput, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "input", *input)

	// random sleep
	fmt.Println("Working......")
	r := rand.Intn(1000)
	time.Sleep(time.Duration(r) * time.Millisecond)

	return &ActivityOutput{Value: fmt.Sprintf("<%s - %d>", input.Name, r)}, nil
}

func ActivityParallel(ctx context.Context, input *ActivityInput) (*ActivityOutput, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity Parallel", "input", *input)

	// random sleep
	fmt.Println("Working in parallel......")
	r := rand.Intn(1000)
	time.Sleep(time.Duration(r) * time.Millisecond)

	return &ActivityOutput{Value: fmt.Sprintf("[%s - %d]", input.Name, r)}, nil
}
