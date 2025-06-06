package temporal

import (
	"errors"

	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/api/serviceerror"

	"go.temporal.io/sdk/internal"
)

/*
If activity fails then *ActivityError is returned to the workflow code. The error has important information about activity
and actual error which caused activity failure. This internal error can be unwrapped using errors.Unwrap() or checked using errors.As().
Below are the possible types of internal error:
1) *ApplicationError: (this should be the most common one)
	*ApplicationError can be returned in two cases:
		- If activity implementation returns *ApplicationError by using NewApplicationError()/NewNonRetryableApplicationError() API.
		  The error would contain a message and optional details. Workflow code could extract details to string typed variable, determine
		  what kind of error it was, and take actions based on it. The details are encoded payload therefore, workflow code needs to know what
          the types of the encoded details are before extracting them.
		- If activity implementation returns errors other than from NewApplicationError() API. In this case GetOriginalType()
		  will return original type of error represented as string. Workflow code could check this type to determine what kind of error it was
		  and take actions based on the type. These errors are retryable by default, unless error type is specified in retry policy.
2) *CanceledError:
	If activity was canceled, internal error will be an instance of *CanceledError. When activity cancels itself by
	returning NewCancelError() it would supply optional details which could be extracted by workflow code.
3) *TimeoutError:
	If activity was timed out (several timeout types), internal error will be an instance of *TimeoutError. The err contains
	details about what type of timeout it was.
4) *PanicError:
	If activity code panic while executing, temporal activity worker will report it as activity failure to temporal server.
	The SDK will present that failure as *PanicError. The error contains a string	representation of the panic message and
	the call stack when panic was happen.
Workflow code could handle errors based on different types of error. Below is sample code of how error handling looks like.

err := workflow.ExecuteActivity(ctx, MyActivity, ...).Get(ctx, nil)
if err != nil {
	var applicationErr *ApplicationError
	if errors.As(err, &applicationErr) {
		// retrieve error message
		fmt.Println(applicationErr.Error())

		// handle activity errors (created via NewApplicationError() API)
		var detailMsg string // assuming activity return error by NewApplicationError("message", true, "string details")
		applicationErr.Details(&detailMsg) // extract strong typed details

		// handle activity errors (errors created other than using NewApplicationError() API)
		switch applicationErr.Type() {
		case "CustomErrTypeA":
			// handle CustomErrTypeA
		case CustomErrTypeB:
			// handle CustomErrTypeB
		default:
			// newer version of activity could return new errors that workflow was not aware of.
		}
	}

	var canceledErr *CanceledError
	if errors.As(err, &canceledErr) {
		// handle cancellation
	}

	var timeoutErr *TimeoutError
	if errors.As(err, &timeoutErr) {
		// handle timeout, could check timeout type by timeoutErr.TimeoutType()
        switch timeoutErr.TimeoutType() {
        case enumspb.TIMEOUT_TYPE_SCHEDULE_TO_START:
			// Handle ScheduleToStart timeout.
        case enumspb.TIMEOUT_TYPE_SCHEDULE_TO_CLOSE:
			// Handle ScheduleToClose timeout.
        case enumspb.TIMEOUT_TYPE_START_TO_CLOSE:
            // Handle StartToClose timeout.
        case enumspb.TIMEOUT_TYPE_HEARTBEAT:
            // Handle heartbeat timeout.
        default:
        }
	}

	var panicErr *PanicError
	if errors.As(err, &panicErr) {
		// handle panic, message and stack trace are available by panicErr.Error() and panicErr.StackTrace()
	}
}
Errors from child workflow should be handled in a similar way, except that instance of *ChildWorkflowExecutionError is returned to
workflow code. It might contain *ActivityError in case if error comes from activity (which in turn will contain on of the errors above),
or *ApplicationError in case if error comes from child workflow itself.

When panic happen in workflow implementation code, SDK catches that panic and causing the workflow task timeout.
That workflow task will be retried at a later time (with exponential backoff retry intervals).
Workflow consumers will get an instance of *WorkflowExecutionError. This error will contain one of errors above.
*/

type (
	// ApplicationError returned from activity implementations with message and optional details.
	ApplicationError = internal.ApplicationError

	// CanceledError returned when operation was canceled.
	CanceledError = internal.CanceledError

	// ActivityError returned from workflow when activity returned an error.
	ActivityError = internal.ActivityError

	// ServerError can be returned from server.
	ServerError = internal.ServerError

	// ChildWorkflowExecutionError returned from workflow when child workflow returned an error.
	ChildWorkflowExecutionError = internal.ChildWorkflowExecutionError

	// NexusOperationError is an error returned when a Nexus Operation has failed.
	//
	// NOTE: Experimental
	NexusOperationError = internal.NexusOperationError

	// ChildWorkflowExecutionAlreadyStartedError is set as the cause of
	// ChildWorkflowExecutionError when failure is due the child workflow having
	// already started.
	ChildWorkflowExecutionAlreadyStartedError = internal.ChildWorkflowExecutionAlreadyStartedError

	// NamespaceNotFoundError is set as the cause when failure is due namespace not found.
	NamespaceNotFoundError = internal.NamespaceNotFoundError

	// WorkflowExecutionError returned from workflow.
	WorkflowExecutionError = internal.WorkflowExecutionError

	// TimeoutError returned when activity or child workflow timed out.
	TimeoutError = internal.TimeoutError

	// TerminatedError returned when workflow was terminated.
	TerminatedError = internal.TerminatedError

	// PanicError contains information about panicked workflow/activity.
	PanicError = internal.PanicError

	// UnknownExternalWorkflowExecutionError can be returned when external workflow doesn't exist
	UnknownExternalWorkflowExecutionError = internal.UnknownExternalWorkflowExecutionError

	// QueryRejectedError is a possible error that can be returned by
	// ClientOutboundInterceptor.QueryWorkflow to indicate that the query was rejected by the server.
	QueryRejectedError = internal.QueryRejectedError
)

var (
	// ErrNoData is returned when trying to extract strong typed data while there is no data available.
	ErrNoData = internal.ErrNoData

	// ErrScheduleAlreadyRunning can be returned when a schedule ID is reused
	ErrScheduleAlreadyRunning = internal.ErrScheduleAlreadyRunning

	// ErrSkipScheduleUpdate is used by a user if they want to skip updating a schedule.
	ErrSkipScheduleUpdate = internal.ErrSkipScheduleUpdate
)

// ApplicationErrorOptions should be used to set all the desired attributes of a new ApplicationError
// To get a new instance use ErrorAttributes function
type ApplicationErrorOptions = internal.ApplicationErrorOptions

// NewApplicationErrorWithOptions creates new instance of *ApplicationError type, all the options of the
// newly created error could be controlled through instance of ApplicationErrorOptions.
// The options structure also receives some extra requests. See activity.ApplicationErrorOptions for details.
func NewApplicationErrorWithOptions(msg, errType string, options ApplicationErrorOptions) error {
	return internal.NewApplicationErrorWithOptions(msg, errType, options)
}

// NewApplicationError creates new instance of retryable *ApplicationError with message, type, and optional details.
// Use ApplicationError for any use case specific errors that cross activity and child workflow boundaries.
// errType can be used to control if error is retryable or not. Add the same type in to RetryPolicy.NonRetryableErrorTypes
// to avoid retrying of particular error types.
func NewApplicationError(message, errType string, details ...interface{}) error {
	return internal.NewApplicationErrorWithOptions(message, errType, ApplicationErrorOptions{Details: details})
}

// NewApplicationErrorWithCause creates new instance of retryable *ApplicationError with message, type, cause, and optional details.
// Use ApplicationError for any use case specific errors that cross activity and child workflow boundaries.
// errType can be used to control if error is retryable or not. Add the same type in to RetryPolicy.NonRetryableErrorTypes
// to avoid retrying of particular error types.
func NewApplicationErrorWithCause(message, errType string, cause error, details ...interface{}) error {
	return internal.NewApplicationErrorWithOptions(
		message, errType, ApplicationErrorOptions{NonRetryable: false, Cause: cause, Details: details},
	)
}

// NewNonRetryableApplicationError creates new instance of non-retryable *ApplicationError with message, type, and optional cause and details.
// Use ApplicationError for any use case specific errors that cross activity and child workflow boundaries.
func NewNonRetryableApplicationError(message, errType string, cause error, details ...interface{}) error {
	return internal.NewApplicationErrorWithOptions(
		message, errType, ApplicationErrorOptions{NonRetryable: true, Cause: cause, Details: details},
	)
}

// NewCanceledError creates CanceledError instance.
// Return this error from activity or child workflow to indicate that it was successfully canceled.
func NewCanceledError(details ...interface{}) error {
	return internal.NewCanceledError(details...)
}

// IsApplicationError return if the err is a ApplicationError
func IsApplicationError(err error) bool {
	var applicationError *ApplicationError
	return errors.As(err, &applicationError)
}

// IsWorkflowExecutionAlreadyStartedError return if the err is a
// WorkflowExecutionAlreadyStartedError or if an error in the chain is a
// ChildWorkflowExecutionAlreadyStartedError.
func IsWorkflowExecutionAlreadyStartedError(err error) bool {
	if _, ok := err.(*serviceerror.WorkflowExecutionAlreadyStarted); ok {
		return ok
	}
	var childError *ChildWorkflowExecutionAlreadyStartedError
	return errors.As(err, &childError)
}

// IsCanceledError return if the err is a CanceledError
func IsCanceledError(err error) bool {
	var cancelError *CanceledError
	return errors.As(err, &cancelError)
}

// IsTimeoutError return if the err is a TimeoutError
func IsTimeoutError(err error) bool {
	var timeoutError *TimeoutError
	return errors.As(err, &timeoutError)
}

// IsTerminatedError return if the err is a TerminatedError
func IsTerminatedError(err error) bool {
	var terminateError *TerminatedError
	return errors.As(err, &terminateError)
}

// IsPanicError return if the err is a PanicError
func IsPanicError(err error) bool {
	var panicError *PanicError
	return errors.As(err, &panicError)
}

// NewTimeoutError creates TimeoutError instance.
// Use NewHeartbeatTimeoutError to create heartbeat TimeoutError
// WARNING: This function is public only to support unit testing of workflows.
// It shouldn't be used by application level code.
func NewTimeoutError(timeoutType enumspb.TimeoutType, lastErr error, details ...interface{}) error {
	return internal.NewTimeoutError("Test timeout", timeoutType, lastErr, details...)
}

// NewHeartbeatTimeoutError creates TimeoutError instance
// WARNING: This function is public only to support unit testing of workflows.
// It shouldn't be used by application level code.
func NewHeartbeatTimeoutError(details ...interface{}) error {
	return internal.NewHeartbeatTimeoutError(details...)
}

// ApplicationErrorCategory sets the category of the error. The category of the error
// maps to logging/metrics SDK behaviours, does not impact server-side logging/metrics.
type ApplicationErrorCategory = internal.ApplicationErrorCategory

const (
	// ApplicationErrorCategoryUnspecified represents an error with an unspecified category.
	ApplicationErrorCategoryUnspecified = internal.ApplicationErrorCategoryUnspecified
	// ApplicationErrorCategoryBenign indicates an error that is expected under normal operation and should not trigger alerts.
	ApplicationErrorCategoryBenign = internal.ApplicationErrorCategoryBenign
)
