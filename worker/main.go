package main

import (
	"github.com/rs/zerolog/log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	logrusadapter "logur.dev/adapter/zerolog"
	"logur.dev/logur"

	"temporal-sample/common"
)

func main() {

	logger := logur.LoggerToKV(logrusadapter.New(common.Logger()))
	c, err := client.Dial(client.Options{Logger: logger})
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create client")
	}
	defer c.Close()

	w := worker.New(c, common.TASK_QUEUE, worker.Options{})

	w.RegisterWorkflow(common.Workflow)
	w.RegisterActivity(common.Activity)
	w.RegisterActivity(common.ActivityParallel)

	// Run the worker in a blocking fashion. Stop the worker when interruptCh receives signal.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to start worker")
	}

	log.Info().Msg("Shutting down worker")
}
