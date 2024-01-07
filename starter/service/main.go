package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	logrusadapter "logur.dev/adapter/zerolog"
	"logur.dev/logur"

	"temporal-sample/common"
)

type Response struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	WorkflowID string `json:"workflow_id,omitempty"`
	RunID      string `json:"run_id,omitempty"`
	Result     string `json:"result,omitempty"`
}

type Body struct {
	Name string `json:"name"`
}

var (
	once sync.Once
	c    client.Client
)

func main() {
	// Create a context that will be canceled when termination signals are received
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	srv := &http.Server{
		Addr: "127.0.0.1:8088",
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
	http.HandleFunc("/", handleWorkflow)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Unable to start server")
		}
	}()

	// Wait for either the server to close gracefully or the context to be canceled
	<-ctx.Done()
	log.Info().Msg("Shutting down server")

	srv.Shutdown(ctx)
	log.Info().Msg("Server gracefully stopped")
	// close the connection to Temporal server
	c.Close()
}

// handleWorkflow handles the workflow execution: POST to start a new workflow execution, GET to get the result of a workflow execution
func handleWorkflow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}

	ctx := r.Context()

	once.Do(func() {
		log.Logger = common.Logger()
		logger := logur.LoggerToKV(logrusadapter.New(common.Logger()))
		var err error
		c, err = client.Dial(client.Options{
			Logger: logger,
		})
		if err != nil {
			httpError(w, http.StatusInternalServerError, "Unable to create client")
			log.Error().Err(err).Msg("Unable to create client")
			return
		}
	})

	// Create a new workflow execution
	if r.Method == "POST" {
		if err := handlePOST(ctx, w, r, c); err != nil {
			log.Error().Err(err).Msg("Unable to handle POST")
		}
		return
	}

	// Get workflow execution result
	if r.Method == "GET" {
		if err := handleGET(ctx, w, r, c); err != nil {
			log.Error().Err(err).Msg("Unable to handle GET")
		}
	}
}

// handlePOST starts a new workflow execution
func handlePOST(ctx context.Context, w http.ResponseWriter, r *http.Request, c client.Client) error {
	workflowOptions := client.StartWorkflowOptions{
		ID:        "workflow-" + uuid.New().String(),
		TaskQueue: common.TASK_QUEUE,
	}

	var body Body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpError(w, http.StatusBadRequest, "Unable to parse request body")
		return fmt.Errorf("unable to parse request body: %v", err)
	}

	wInput := &common.WorkflowInput{
		Name: body.Name,
	}

	we, err := c.ExecuteWorkflow(ctx, workflowOptions, common.Workflow, wInput)
	if err != nil {
		httpError(w, http.StatusInternalServerError, "Unable to execute workflow")
		return fmt.Errorf("unable to execute workflow: %v", err)
	}
	log.Info().
		Str("WorkflowID", we.GetID()).
		Str("RunID", we.GetRunID()).
		Msg("Started workflow")

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(Response{Status: http.StatusAccepted, Message: "Workflow started", WorkflowID: we.GetID(), RunID: we.GetRunID()})
	return nil
}

// handleGET gets the result of a workflow execution
func handleGET(ctx context.Context, w http.ResponseWriter, r *http.Request, c client.Client) error {
	workflowID := r.URL.Query().Get("workflow_id")
	runID := r.URL.Query().Get("run_id")

	if workflowID == "" || runID == "" {
		httpError(w, http.StatusBadRequest, "Missing workflow_id or run_id")
		return fmt.Errorf("missing workflow_id or run_id")
	}

	we := c.GetWorkflow(ctx, workflowID, runID)

	var result common.WorkflowOutput
	if err := we.Get(ctx, &result); err != nil {
		httpError(w, http.StatusInternalServerError, "Unable get workflow result")
		return fmt.Errorf("unable get workflow result: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: http.StatusOK, Message: "Workflow completed", WorkflowID: we.GetID(), RunID: we.GetRunID(), Result: result.Value})
	return nil
}

func httpError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Status: status, Message: message})
}
