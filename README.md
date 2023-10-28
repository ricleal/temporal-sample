# temporal-sample

This shows running a workflow in temporal.
This workflow is composed by 2 steps:
1. Run the `Activity` step
2. When the result of the `Activity` step is ready, `N` activities (`ActivityParallel`) are run with input from the output of the previous step.
3. When all activities are done, the workflow is done and the result is printed.

## Running the workflow

Run the temporal server:
```sh
docker-compose up
```

Register the worker (this thread does all the work):
```sh
temporal-sample/worker
❯ go run .
```

Start the workflow:
```sh
temporal-sample/starter
❯ go run .
```

See the UI: http://localhost:8080/namespaces/default/workflows

## Output

You should see:


### Worker

```sh
❯ go run .
9:58AM INF Started Worker Namespace=default TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@
9:59AM INF workflow started Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow name=Ricardo
9:59AM INF Activity started Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow input={"Name":"Ricardo"}
9:59AM DBG ExecuteActivity ActivityID=5 ActivityType=Activity Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow
9:59AM INF Activity ActivityID=5 ActivityType=Activity Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow input={"Name":"Ricardo"}
Working......
9:59AM DBG ExecuteActivity ActivityID=11 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow
9:59AM DBG ExecuteActivity ActivityID=12 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow
9:59AM DBG ExecuteActivity ActivityID=13 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow
9:59AM DBG ExecuteActivity ActivityID=14 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow
9:59AM DBG ExecuteActivity ActivityID=15 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow
9:59AM DBG ExecuteActivity ActivityID=16 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow
9:59AM DBG ExecuteActivity ActivityID=17 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow
9:59AM DBG ExecuteActivity ActivityID=18 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow
9:59AM DBG ExecuteActivity ActivityID=19 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow
9:59AM DBG ExecuteActivity ActivityID=20 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow
9:59AM INF Parallel Activities started Attempt=1 N=10 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow
9:59AM INF Activity Parallel ActivityID=20 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow input={"Name":"\u003cRicardo - 279\u003e"}
Working in parallel......
9:59AM INF Activity Parallel ActivityID=14 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow input={"Name":"\u003cRicardo - 279\u003e"}
Working in parallel......
9:59AM INF Activity Parallel ActivityID=18 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow input={"Name":"\u003cRicardo - 279\u003e"}
Working in parallel......
9:59AM INF Activity Parallel ActivityID=19 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow input={"Name":"\u003cRicardo - 279\u003e"}
Working in parallel......
9:59AM INF Activity Parallel ActivityID=12 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow input={"Name":"\u003cRicardo - 279\u003e"}
Working in parallel......
9:59AM INF Activity Parallel ActivityID=15 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow input={"Name":"\u003cRicardo - 279\u003e"}
Working in parallel......
9:59AM INF Activity Parallel ActivityID=17 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow input={"Name":"\u003cRicardo - 279\u003e"}
Working in parallel......
9:59AM INF Activity Parallel ActivityID=13 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow input={"Name":"\u003cRicardo - 279\u003e"}
Working in parallel......
9:59AM INF Activity Parallel ActivityID=16 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow input={"Name":"\u003cRicardo - 279\u003e"}
Working in parallel......
9:59AM INF Activity Parallel ActivityID=11 ActivityType=ActivityParallel Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow input={"Name":"\u003cRicardo - 279\u003e"}
Working in parallel......
9:59AM INF Activity returned with result Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow result={"Value":"[\u003cRicardo - 279\u003e - 256]"}
9:59AM INF Activity returned with result Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow result={"Value":"[\u003cRicardo - 279\u003e - 437]"}
9:59AM INF Activity returned with result Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow result={"Value":"[\u003cRicardo - 279\u003e - 740]"}
9:59AM INF Activity returned with result Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow result={"Value":"[\u003cRicardo - 279\u003e - 795]"}
9:59AM INF Activity returned with result Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow result={"Value":"[\u003cRicardo - 279\u003e - 738]"}
9:59AM INF Activity returned with result Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow result={"Value":"[\u003cRicardo - 279\u003e - 637]"}
9:59AM INF Activity returned with result Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow result={"Value":"[\u003cRicardo - 279\u003e - 948]"}
9:59AM INF Activity returned with result Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow result={"Value":"[\u003cRicardo - 279\u003e - 979]"}
9:59AM INF Activity returned with result Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow result={"Value":"[\u003cRicardo - 279\u003e - 255]"}
9:59AM INF Activity returned with result Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow result={"Value":"[\u003cRicardo - 279\u003e - 949]"}
9:59AM INF workflow completed Attempt=1 Namespace=default RunID=954a4c37-2f96-4978-b2fd-2b65d347400f TaskQueue=sample_task_queue WorkerID=80267@MacBook-Pro@ WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7 WorkflowType=Workflow result="+[<Ricardo - 279> - 256]+[<Ricardo - 279> - 437]+[<Ricardo - 279> - 740]+[<Ricardo - 279> - 795]+[<Ricardo - 279> - 738]+[<Ricardo - 279> - 637]+[<Ricardo - 279> - 948]+[<Ricardo - 279> - 979]+[<Ricardo - 279> - 255]+[<Ricardo - 279> - 949]"
```

### Starter

```sh
❯ go run .
9:59AM INF Started workflow RunID=954a4c37-2f96-4978-b2fd-2b65d347400f WorkflowID=workflow-6a7885e9-f5cc-45aa-8b2d-ea3bff3605f7
9:59AM INF Workflow completed result={"Value":"+[\u003cRicardo - 279\u003e - 256]+[\u003cRicardo - 279\u003e - 437]+[\u003cRicardo - 279\u003e - 740]+[\u003cRicardo - 279\u003e - 795]+[\u003cRicardo - 279\u003e - 738]+[\u003cRicardo - 279\u003e - 637]+[\u003cRicardo - 279\u003e - 948]+[\u003cRicardo - 279\u003e - 979]+[\u003cRicardo - 279\u003e - 255]+[\u003cRicardo - 279\u003e - 949]"}
```