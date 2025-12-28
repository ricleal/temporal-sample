package common

import (
	"context"
	"fmt"
	"log/slog"
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
	FibonacciIndex int // The Fibonacci number to calculate
}

type WorkflowOutput struct {
	Name    string        `json:"name"`
	Results map[int]int64 `json:"results"` // Map of index to Fibonacci result
}

// Activity input
type ActivityInput struct {
	Index int // Fibonacci index to compute
}

// FibonacciInfo contains metadata about the Fibonacci computation
type FibonacciInfo struct {
	Index        int       `json:"index"`
	Result       int64     `json:"result"`
	ComputedAt   time.Time `json:"computed_at"`
	Iterations   int       `json:"iterations"`
	IsLarge      bool      `json:"is_large"`      // true if index > 30
	ResultDigits int       `json:"result_digits"` // number of digits in result
}

type ActivityOutput struct {
	Value  string        `json:"value"`
	Result int64         `json:"result"` // Fibonacci result
	Info   FibonacciInfo `json:"info"`   // Detailed computation info
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
	logger.Info("workflow started", "fibonacci_index", wInput.FibonacciIndex)

	// Validate input
	if wInput.FibonacciIndex < 0 {
		return nil, fmt.Errorf("fibonacci index must be non-negative")
	}
	if wInput.FibonacciIndex > 1000 {
		return nil, fmt.Errorf("fibonacci index too large (max 1000)")
	}

	aInput := &ActivityInput{
		Index: wInput.FibonacciIndex,
	}

	logger.Info("Activity started", "input", aInput)
	var aResult ActivityOutput
	err := workflow.ExecuteActivity(ctx, Activity, aInput).Get(ctx, &aResult)
	if err != nil {
		logger.Error("Activity failed", "error", err)
		return nil, fmt.Errorf("activity failed: %w", err)
	}

	// Calculate the number of parallel activities based on input
	// We'll compute Fibonacci numbers from (index - ParallelActivities + 1) to index
	startIndex := wInput.FibonacciIndex - ParallelActivities + 1
	if startIndex < 0 {
		startIndex = 0
	}

	// Launch parallel activities to compute multiple Fibonacci numbers
	logger.Info("Parallel Activities started", "count", ParallelActivities, "start_index", startIndex)
	selector := workflow.NewSelector(ctx)
	var results []ActivityOutput
	var pending int = 0

	for i := 0; i < ParallelActivities; i++ {
		index := startIndex + i
		if index > wInput.FibonacciIndex {
			break
		}
		pending++
		newInput := &ActivityInput{
			Index: index,
		}
		future := workflow.ExecuteActivity(ctx, ActivityParallel, newInput)
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

	// Collect all results into a map
	resultsMap := make(map[int]int64)
	for _, r := range results {
		resultsMap[r.Info.Index] = r.Info.Result
	}

	logger.Info("workflow completed", "results_count", len(resultsMap))
	return &WorkflowOutput{Name: "fibonacci", Results: resultsMap}, nil
}

func Activity(ctx context.Context, input *ActivityInput) (*ActivityOutput, error) {
	logger := activity.GetLogger(ctx)
	info := activity.GetInfo(ctx)
	logger.Info("Activity", "input", *input, "attempt", info.Attempt)

	// Compute Fibonacci number with heartbeat for long computations
	computedAt := time.Now()
	result := fibonacciWithHeartbeat(ctx, input.Index)
	logger.Debug("Activity computed Fibonacci", "index", input.Index, "result", result)

	// Calculate result digits
	resultDigits := len(fmt.Sprintf("%d", result))

	return &ActivityOutput{
		Value:  fmt.Sprintf("Fib(%d)=%d", input.Index, result),
		Result: result,
		Info: FibonacciInfo{
			Index:        input.Index,
			Result:       result,
			ComputedAt:   computedAt,
			Iterations:   input.Index,
			IsLarge:      input.Index > 30,
			ResultDigits: resultDigits,
		},
	}, nil
}

func ActivityParallel(ctx context.Context, input *ActivityInput) (*ActivityOutput, error) {
	logger := activity.GetLogger(ctx)
	info := activity.GetInfo(ctx)
	logger.Info("Activity Parallel", "input", *input, "attempt", info.Attempt)

	// Compute Fibonacci number with heartbeat for long computations
	computedAt := time.Now()
	result := fibonacciWithHeartbeat(ctx, input.Index)
	logger.Debug("Activity Parallel computed Fibonacci", "index", input.Index, "result", result)

	// Calculate result digits
	resultDigits := len(fmt.Sprintf("%d", result))

	return &ActivityOutput{
		Value:  fmt.Sprintf("Fib(%d)", input.Index),
		Result: result,
		Info: FibonacciInfo{
			Index:        input.Index,
			Result:       result,
			ComputedAt:   computedAt,
			Iterations:   input.Index,
			IsLarge:      input.Index > 30,
			ResultDigits: resultDigits,
		},
	}, nil
}

// fibonacci computes the nth Fibonacci number iteratively
func fibonacci(n int) int64 {
	if n <= 1 {
		return int64(n)
	}

	var a, b int64 = 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// fibonacciWithHeartbeat computes the nth Fibonacci number and records heartbeats for larger values
func fibonacciWithHeartbeat(ctx context.Context, n int) int64 {
	if n <= 1 {
		return int64(n)
	}

	var a, b int64 = 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
		// Record heartbeat every 1000 iterations for large numbers
		if i%1000 == 0 {
			activity.RecordHeartbeat(ctx, fmt.Sprintf("computing iteration %d/%d", i, n))
		}
	}
	return b
}
