// Command service provides an HTTP API for Temporal workflow management.
// It exposes REST endpoints to start new workflow executions (POST) and
// query workflow results (GET), running on port 8888 by default.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"

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
	once   sync.Once
	c      client.Client
	logger *slog.Logger
)

func main() {
	// Create a context that will be canceled when termination signals are received
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger = common.Logger()

	myServicePort := os.Getenv("MY_SERVICE_PORT")
	if myServicePort == "" {
		myServicePort = "8888"
	}
	srv := &http.Server{
		Addr: "127.0.0.1:" + myServicePort,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
	http.HandleFunc("/", handleWorkflow)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Unable to start server", "error", err)
			panic(err)
		}
	}()

	logger.Info("Server started", "port", myServicePort)

	// Wait for either the server to close gracefully or the context to be canceled
	<-ctx.Done()
	logger.Info("Shutting down server")

	srv.Shutdown(ctx)
	logger.Info("Server gracefully stopped")
	// close the connection to Temporal server
	c.Close()
}

// handleWorkflow handles the workflow execution: POST to start a new workflow execution, GET to get the result of a workflow execution
func handleWorkflow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "GET" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}

	ctx := r.Context()

	once.Do(func() {
		var err error
		c, err = client.Dial(client.Options{
			Logger: log.NewStructuredLogger(logger),
		})
		if err != nil {
			httpError(w, http.StatusInternalServerError, "Unable to create client")
			logger.Error("Unable to create client", "error", err)
			return
		}
	})

	// Create a new workflow execution
	if r.Method == "POST" {
		if err := handlePOST(ctx, w, r, c); err != nil {
			logger.Error("Unable to handle POST", "error", err)
		}
		return
	}

	// Get workflow execution result
	if r.Method == "GET" {
		if err := handleGET(ctx, w, r, c); err != nil {
			logger.Error("Unable to handle GET", "error", err)
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
	logger.Info("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	w.Header().Set("Content-Type", "application/json")
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: http.StatusOK, Message: "Workflow completed", WorkflowID: we.GetID(), RunID: we.GetRunID(), Result: result.Value})
	return nil
}

func httpError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Status: status, Message: message})
}
