// Command worker starts a Temporal worker that listens for and executes workflow tasks.
// It registers the workflow and activity implementations and processes tasks from
// the configured task queue until interrupted.
package main

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"

	"temporal-sample/common"
)

func main() {
	logger := common.Logger()
	c, err := client.Dial(client.Options{Logger: log.NewStructuredLogger(logger)})
	if err != nil {
		logger.Error("Unable to create client", "error", err)
		panic(err)
	}
	defer c.Close()

	w := worker.New(c, common.TASK_QUEUE, worker.Options{})

	w.RegisterWorkflow(common.Workflow)
	w.RegisterActivity(common.Activity)
	w.RegisterActivity(common.ActivityParallel)

	// Run the worker in a blocking fashion. Stop the worker when interruptCh receives signal.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		logger.Error("Unable to start worker", "error", err)
		panic(err)
	}

	logger.Info("Shutting down worker")
}
