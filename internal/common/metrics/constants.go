package metrics

// Metrics keys
const (
	TemporalMetricsPrefix = "temporal_"

	WorkflowCompletedCounter     = TemporalMetricsPrefix + "workflow_completed"
	WorkflowCanceledCounter      = TemporalMetricsPrefix + "workflow_canceled"
	WorkflowFailedCounter        = TemporalMetricsPrefix + "workflow_failed"
	WorkflowContinueAsNewCounter = TemporalMetricsPrefix + "workflow_continue_as_new"
	WorkflowEndToEndLatency      = TemporalMetricsPrefix + "workflow_endtoend_latency" // measure workflow execution from start to close

	WorkflowTaskReplayLatency           = TemporalMetricsPrefix + "workflow_task_replay_latency"
	WorkflowTaskQueuePollEmptyCounter   = TemporalMetricsPrefix + "workflow_task_queue_poll_empty"
	WorkflowTaskQueuePollSucceedCounter = TemporalMetricsPrefix + "workflow_task_queue_poll_succeed"
	WorkflowTaskScheduleToStartLatency  = TemporalMetricsPrefix + "workflow_task_schedule_to_start_latency"
	WorkflowTaskExecutionLatency        = TemporalMetricsPrefix + "workflow_task_execution_latency"
	WorkflowTaskExecutionFailureCounter = TemporalMetricsPrefix + "workflow_task_execution_failed"
	WorkflowTaskNoCompletionCounter     = TemporalMetricsPrefix + "workflow_task_no_completion"

	ActivityPollNoTaskCounter             = TemporalMetricsPrefix + "activity_poll_no_task"
	ActivityScheduleToStartLatency        = TemporalMetricsPrefix + "activity_schedule_to_start_latency"
	ActivityExecutionFailedCounter        = TemporalMetricsPrefix + "activity_execution_failed"
	UnregisteredActivityInvocationCounter = TemporalMetricsPrefix + "unregistered_activity_invocation"
	ActivityExecutionLatency              = TemporalMetricsPrefix + "activity_execution_latency"
	ActivitySucceedEndToEndLatency        = TemporalMetricsPrefix + "activity_succeed_endtoend_latency"
	ActivityTaskErrorCounter              = TemporalMetricsPrefix + "activity_task_error"

	LocalActivityTotalCounter             = TemporalMetricsPrefix + "local_activity_total"
	LocalActivityCanceledCounter          = TemporalMetricsPrefix + "local_activity_canceled" // Deprecated: Use LocalActivityExecutionCanceledCounter instead.
	LocalActivityExecutionCanceledCounter = TemporalMetricsPrefix + "local_activity_execution_cancelled"
	LocalActivityFailedCounter            = TemporalMetricsPrefix + "local_activity_failed" // Deprecated: Use LocalActivityExecutionFailedCounter instead.
	LocalActivityExecutionFailedCounter   = TemporalMetricsPrefix + "local_activity_execution_failed"
	LocalActivityErrorCounter             = TemporalMetricsPrefix + "local_activity_error"
	LocalActivityExecutionLatency         = TemporalMetricsPrefix + "local_activity_execution_latency"
	LocalActivitySucceedEndToEndLatency   = TemporalMetricsPrefix + "local_activity_succeed_endtoend_latency"

	CorruptedSignalsCounter = TemporalMetricsPrefix + "corrupted_signals"

	WorkerStartCounter       = TemporalMetricsPrefix + "worker_start"
	WorkerTaskSlotsAvailable = TemporalMetricsPrefix + "worker_task_slots_available"
	WorkerTaskSlotsUsed      = TemporalMetricsPrefix + "worker_task_slots_used"
	PollerStartCounter       = TemporalMetricsPrefix + "poller_start"
	NumPoller                = TemporalMetricsPrefix + "num_pollers"

	TemporalRequest                      = TemporalMetricsPrefix + "request"
	TemporalRequestFailure               = TemporalRequest + "_failure"
	TemporalRequestLatency               = TemporalRequest + "_latency"
	TemporalLongRequest                  = TemporalMetricsPrefix + "long_request"
	TemporalLongRequestFailure           = TemporalLongRequest + "_failure"
	TemporalLongRequestLatency           = TemporalLongRequest + "_latency"
	TemporalRequestResourceExhausted     = TemporalRequest + "_resource_exhausted"
	TemporalLongRequestResourceExhausted = TemporalLongRequest + "_resource_exhausted"

	StickyCacheHit                 = TemporalMetricsPrefix + "sticky_cache_hit"
	StickyCacheMiss                = TemporalMetricsPrefix + "sticky_cache_miss"
	StickyCacheTotalForcedEviction = TemporalMetricsPrefix + "sticky_cache_total_forced_eviction"
	StickyCacheSize                = TemporalMetricsPrefix + "sticky_cache_size"

	WorkflowActiveThreadCount = TemporalMetricsPrefix + "workflow_active_thread_count"

	NexusPollNoTaskCounter          = TemporalMetricsPrefix + "nexus_poll_no_task"
	NexusTaskScheduleToStartLatency = TemporalMetricsPrefix + "nexus_task_schedule_to_start_latency"
	NexusTaskExecutionFailedCounter = TemporalMetricsPrefix + "nexus_task_execution_failed"
	NexusTaskExecutionLatency       = TemporalMetricsPrefix + "nexus_task_execution_latency"
	NexusTaskEndToEndLatency        = TemporalMetricsPrefix + "nexus_task_endtoend_latency"
)

// Metric tag keys
const (
	NamespaceTagName        = "namespace"
	ClientTagName           = "client_name"
	PollerTypeTagName       = "poller_type"
	WorkerTypeTagName       = "worker_type"
	WorkflowTypeNameTagName = "workflow_type"
	ActivityTypeNameTagName = "activity_type"
	NexusServiceTagName     = "nexus_service"
	NexusOperationTagName   = "nexus_operation"
	FailureReasonTagName    = "failure_reason"
	TaskQueueTagName        = "task_queue"
	OperationTagName        = "operation"
	CauseTagName            = "cause"
	RequestFailureCode      = "status_code"
)

// Metric tag values
const (
	NoneTagValue                 = "none"
	ClientTagValue               = "temporal_go"
	PollerTypeWorkflowTask       = "workflow_task"
	PollerTypeWorkflowStickyTask = "workflow_sticky_task"
	PollerTypeActivityTask       = "activity_task"
	PollerTypeNexusTask          = "nexus_task"
)
