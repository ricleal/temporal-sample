package common

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func Logger() *slog.Logger {
	// Get log level from environment variable, default to INFO
	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		levelStr = "INFO"
	}

	var level slog.Level
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN", "WARNING":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	return slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      level,
		TimeFormat: time.Kitchen,
	}))
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
	// Configure activity options with retry policy
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("workflow started", "name", wInput.Name)

	aInput := &ActivityInput{
		Name: wInput.Name,
	}

	logger.Info("Activity started", "input", aInput)
	var aResult ActivityOutput
	err := workflow.ExecuteActivity(ctx, Activity, aInput).Get(ctx, &aResult)
	if err != nil {
		logger.Error("Activity failed", "error", err)
		return nil, fmt.Errorf("activity failed: %w", err)
	}

	newActivityInput := &ActivityInput{
		Name: aResult.Value,
	}

	// Launch several activities in parallel using selector for efficient execution
	logger.Info("Parallel Activities started", "count", ParallelActivities)
	selector := workflow.NewSelector(ctx)
	var results []ActivityOutput
	var pending int = ParallelActivities

	for i := 0; i < ParallelActivities; i++ {
		future := workflow.ExecuteActivity(ctx, ActivityParallel, newActivityInput)
		selector.AddFuture(future, func(f workflow.Future) {
			var aOut ActivityOutput
			if err := f.Get(ctx, &aOut); err != nil {
				logger.Error("Activity Parallel failed", "error", err)
				return
			}
			logger.Info("Activity returned with result", "result", aOut)
			results = append(results, aOut)
			pending--
		})
	}

	// Wait for all activities to complete
	for pending > 0 {
		selector.Select(ctx)
	}

	// Check if we got all results
	if len(results) != ParallelActivities {
		return nil, fmt.Errorf("expected %d results, got %d", ParallelActivities, len(results))
	}

	// Accumulate results
	var v string
	for _, r := range results {
		v += "+" + r.Value
	}

	logger.Info("workflow completed", "result", v)
	return &WorkflowOutput{Value: v}, nil
}

func Activity(ctx context.Context, input *ActivityInput) (*ActivityOutput, error) {
	logger := activity.GetLogger(ctx)
	info := activity.GetInfo(ctx)
	logger.Info("Activity", "input", *input, "attempt", info.Attempt)

	// Simulate work with random duration
	sleepDuration := rand.Intn(1000)
	logger.Debug("Activity working", "duration_ms", sleepDuration)

	// Use activity.RecordHeartbeat for long-running activities
	if sleepDuration > 500 {
		activity.RecordHeartbeat(ctx, "processing")
	}

	time.Sleep(time.Duration(sleepDuration) * time.Millisecond)

	return &ActivityOutput{Value: fmt.Sprintf("<%s - %d>", input.Name, sleepDuration)}, nil
}

func ActivityParallel(ctx context.Context, input *ActivityInput) (*ActivityOutput, error) {
	logger := activity.GetLogger(ctx)
	info := activity.GetInfo(ctx)
	logger.Info("Activity Parallel", "input", *input, "attempt", info.Attempt)

	// Simulate work with random duration
	sleepDuration := rand.Intn(1000)
	logger.Debug("Activity Parallel working", "duration_ms", sleepDuration)

	// Use activity.RecordHeartbeat for long-running activities
	if sleepDuration > 500 {
		activity.RecordHeartbeat(ctx, "processing")
	}

	time.Sleep(time.Duration(sleepDuration) * time.Millisecond)

	return &ActivityOutput{Value: fmt.Sprintf("[%s - %d]", input.Name, sleepDuration)}, nil
}
