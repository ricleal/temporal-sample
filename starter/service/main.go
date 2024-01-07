package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

func main() {
	http.HandleFunc("/", handlePost)
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to start server")
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{Status: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}

	ctx := r.Context()
	log.Logger = common.Logger()

	logger := logur.LoggerToKV(logrusadapter.New(common.Logger()))

	c, err := client.Dial(client.Options{
		Logger: logger,
	})
	if err != nil {
		httpError(w, http.StatusInternalServerError, "Unable to create client")
		log.Error().Err(err).Msg("Unable to create client")
		return
	}
	defer c.Close()

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
