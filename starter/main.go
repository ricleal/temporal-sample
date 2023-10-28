package main

import (
	"context"
	"temporal-sample/common"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	logrusadapter "logur.dev/adapter/zerolog"
	"logur.dev/logur"
)

func main() {

	log.Logger = common.Logger()

	logger := logur.LoggerToKV(logrusadapter.New(common.Logger()))
	c, err := client.Dial(client.Options{Logger: logger})
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create client")
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "workflow-" + uuid.New().String(),
		TaskQueue: common.TASK_QUEUE,
	}

	wInput := &common.WorkflowInput{
		Name: "Ricardo",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, common.Workflow, wInput)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to execute workflow")
	}
	log.Info().
		Str("WorkflowID", we.GetID()).
		Str("RunID", we.GetRunID()).
		Msg("Started workflow")

	// Synchronously wait for the workflow completion.
	var result common.WorkflowOutput
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable get workflow result")
	}
	log.Info().Interface("result", result).Msg("Workflow completed")
}
