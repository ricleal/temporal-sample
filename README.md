# temporal-sample

This project demonstrates running workflows in Temporal with parallel activities.

## Overview

This workflow is composed of 3 steps:
1. Run the `Activity` step
2. When the result of the `Activity` step is ready, 10 parallel activities (`ActivityParallel`) are run with input from the output of the previous step
3. When all activities are done, the workflow completes and returns the result

The project includes three executables:
- **worker** - Temporal worker that executes workflow and activity tasks
- **cmd/cli** - CLI tool to start a workflow and wait for its completion
- **cmd/server** - HTTP REST API to start workflows and query results

## Prerequisites

- Go 1.25+
- Docker and Docker Compose
- Temporal server (via docker-compose)

## Running the workflow

### Start the Temporal server

```sh
docker-compose up
```

Check the Temporal UI is running at: http://localhost:8233/

### Start the worker

The worker must be running to execute workflow tasks:

```sh
cd worker
go run .
```

### Start a workflow (CLI)

```sh
cd cmd/cli
go run .
```

View workflows in the UI: http://localhost:8233/namespaces/default/workflows

## Output

You should see:


### Worker

```sh
❯ go run .
9:58AM INF Started Worker Namespace=default TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@
9:59AM INF workflow started Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow name=Ricardo
(...)
9:59AM INF workflow completed Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow result="+[<Ricardo - 279> - 256]+[<Ricardo - 279> - 437]+[<Ricardo - 279> - 740]+[<Ricardo - 279> - 795]+[<Ricardo - 279> - 738]+[<Ricardo - 279> - 637]+[<Ricardo - 279> - 948]+[<Ricardo - 279> - 979]+[<Ricardo - 279> - 255]+[<Ricardo - 279> - 949]"
```

### Starter

```sh
❯ go run .
9:59AM INF Started workflow RunID=954a4c37-2f96-4978-b2fd-2b65d347400f WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7
9:59AM INF Workflow completed result={"Value":"+[\u003cRicardo - 279\u003e - 256]+[\u003cRicardo - 279\u003e - 437]+[\u003cRicardo - 279\u003e - 740]+[\u003cRicardo - 279\u003e - 795]+[\u003cRicardo - 279\u003e - 738]+[\u003cRicardo - 279\u003e - 637]+[\u003cRicardo - 279\u003e - 948]+[\u003cRicardo - 279\u003e - 979]+[\u003cRicardo - 279\u003e - 255]+[\u003cRicardo - 279\u003e - 949]"}
```


## Web service

The HTTP service provides REST endpoints for workflow management.

### Start the service

First, ensure the worker is running, then start the service:

```sh
cd cmd/server
go run .
```

The service listens on port 8888 by default (configurable via `MY_SERVICE_PORT` environment variable).

### Start a workflow (POST)

```sh
http -v POST localhost:8888 "name=ricardo"
```

Response:

```sh
HTTP/1.1 202 Accepted
Content-Type: application/json
Content-Length: 154
Date: Sun, 07 Jan 2024 21:11:42 GMT

{
    "message": "Workflow started",
    "run_id": "3dedcd0f-50bd-40c9-abd5-9214cd4e8e7c",
    "status": 202,
    "workflow_id": "workflow-b4f3b8b1-534c-45e5-b8fc-0f15c8a55114"
}
```

### Get workflow result (GET)

```sh
http 'localhost:8888/?run_id=3dedcd0f-50bd-40c9-abd5-9214cd4e8e7c&workflow_id=workflow-b4f3b8b1-534c-45e5-b8fc-0f15c8a55114'
```

Response:

```sh
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 506
Date: Sun, 07 Jan 2024 21:16:08 GMT

{
    "message": "Workflow completed",
    "result": "+[<ricardo - 227> - 969]+[<ricardo - 227> - 643]+[<ricardo - 227> - 964]+[<ricardo - 227> - 695]+[<ricardo - 227> - 30]+[<ricardo - 227> - 403]+[<ricardo - 227> - 29]+[<ricardo - 227> - 157]+[<ricardo - 227> - 928]+[<ricardo - 227> - 499]",
    "run_id": "3dedcd0f-50bd-40c9-abd5-9214cd4e8e7c",
    "status": 200,
    "workflow_id": "workflow-b4f3b8b1-534c-45e5-b8fc-0f15c8a55114"
}
```

## Load testing with Vegeta

You can perform load testing on the HTTP service using [Vegeta](https://github.com/tsenart/vegeta).

First, start the service:

```sh
cd cmd/server
go run .
```

Run the load test:

```sh
# run the load test
echo "POST http://localhost:8888/ Content-Type: application/json" | vegeta attack -body ./cmd/server/body.json -rate 100 -duration 1s | tee results.bin | vegeta report
# save the results in a json file (metrics.json)
vegeta report -type=json results.bin > metrics.json
# plot the results. View the plot.html file in a browser
cat results.bin | vegeta plot > plot.html
# plot the results as a histogram (in the terminal)
cat results.bin | vegeta report -type="hist[0,20ms,40ms,60ms,80ms,100ms,150ms,200ms]"
```


