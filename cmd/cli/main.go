// Command starter executes a Temporal workflow synchronously.
// It connects to the Temporal server, starts a workflow execution,
// and waits for the workflow to complete before printing the result.
package main

import (
	"context"

	"temporal-sample/common"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
)

func main() {
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
		Name: "Ricardo",
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
	logger.Info("Workflow completed", "result", result)
}
