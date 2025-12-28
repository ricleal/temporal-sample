package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"temporal-sample/common"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
)

func main() {
	// Parse command-line flags
	fibIndex := flag.Int("n", 20, "Fibonacci index to calculate (0-1000)")
	flag.Parse()

	ctx := context.Background()
	logger := common.Logger()

	c, err := client.Dial(client.Options{
		Logger: log.NewStructuredLogger(logger),
	})
	if err != nil {
		logger.Error("Unable to create client", "error", err)
		panic(err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "workflow-" + uuid.New().String(),
		TaskQueue: common.TaskQueue,
	}

	wInput := &common.WorkflowInput{
		FibonacciIndex: *fibIndex,
	}

	we, err := c.ExecuteWorkflow(ctx, workflowOptions, common.Workflow, wInput)
	if err != nil {
		logger.Error("Unable to execute workflow", "error", err)
		panic(err)
	}
	logger.Info("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion.
	var result common.WorkflowOutput
	err = we.Get(ctx, &result)
	if err != nil {
		logger.Error("Unable get workflow result", "error", err)
		panic(err)
	}

	// Marshal result to JSON for readable output
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logger.Error("Unable to marshal result to JSON", "error", err)
		panic(err)
	}

	logger.Info("Workflow completed")
	fmt.Println("\nResult:")
	fmt.Println(string(resultJSON))
}
