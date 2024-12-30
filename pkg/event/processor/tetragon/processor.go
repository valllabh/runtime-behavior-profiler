package eventprocessortetragon

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/signal"
	eventprocessortetragontype "runtime-behavior-profiler/pkg/event/processor/tetragon/type"
	"syscall"
	"time"

	"github.com/cilium/tetragon/api/v1/tetragon"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var (
	Debug         bool
	ServerAddress string
	Timeout       time.Duration = 10 * time.Second
	Retries       int           = 4
	Options       eventprocessortetragontype.Opts
)

// gRGC A6 - gRPC Retry Design (a.k.a. built in backoff retry)
// https://github.com/grpc/proposal/blob/master/A6-client-retries.md
// was implemented by https://github.com/grpc/grpc-go/pull/2111 but unusable
// for a long time since maxAttempts was limited to hardcoded 5
// (https://github.com/grpc/grpc-go/issues/4615), recent PR fixed that
// https://github.com/grpc/grpc-go/pull/7229.
//
// It's transparent to the user, to see it in action, make sure the gRPC server
// is unreachable (do not start tetragon for example), run tetra with:
// GRPC_GO_LOG_SEVERITY_LEVEL=warning <tetra cmd>
//
// Note that logs don't always have the time to be pushed before exit so output
// might be a bit off but the number of retries is respected (you can debug or
// synchronously print in the grpc/stream.c:shouldRetry or :withRetry to
// verify).
//
// Also note that the final backoff duration is completely random and chosen
// between 0 and the final duration that was computed via to the params:
// https://github.com/grpc/grpc-go/blob/v1.65.0/stream.go#L702
func retryPolicy(retries int) string {
	if retries < 0 {
		// gRPC should ignore the invalid retry policy but will issue a warning,
		return "{}"
	}
	// maxAttempt includes the first call
	maxAttempt := retries + 1
	// let's not limit backoff by hardcoding 1h in MaxBackoff
	// since we need to provide a value >0
	return fmt.Sprintf(`{
	"methodConfig": [{
	  "name": [{"service": "tetragon.FineGuidanceSensors"}],
	  "retryPolicy": {
		  "MaxAttempts": %d,
		  "InitialBackoff": "1s",
		  "MaxBackoff": "3600s",
		  "BackoffMultiplier": 2,
		  "RetryableStatusCodes": [ "UNAVAILABLE" ]
	  }
	}]}`, maxAttempt)
}

func CliRunErr(fn func(ctx context.Context, cli tetragon.FineGuidanceSensorsClient), fnErr func(err error)) {
	c, err := NewClientWithDefaultContextAndAddress()
	if err != nil {
		fnErr(err)
		return
	}
	defer c.Close()
	fn(c.Ctx, c.Client)
}

func CliRun(fn func(ctx context.Context, cli tetragon.FineGuidanceSensorsClient)) {
	CliRunErr(fn, func(_ error) {})
}

type ClientWithContext struct {
	Client tetragon.FineGuidanceSensorsClient
	// Ctx is a combination of the signal context and the timeout context
	Ctx context.Context
	// SignalCtx is only the signal context, you might want to use that context
	// when the command should never timeout (like a stream command)
	SignalCtx context.Context
	conn      *grpc.ClientConn
	// The signal context is the parent of the timeout context, so cancelling
	// signal will cancel its child, timeout
	signalCancel  context.CancelFunc
	timeoutCancel context.CancelFunc
}

// Close cleanup resources, it closes the connection and cancel the context
func (c ClientWithContext) Close() {
	c.conn.Close()
	c.signalCancel()
	c.timeoutCancel() // this should be a nop
}

// NewClientWithDefaultContextAndAddress returns a client to a tetragon
// server after resolving the server address using helpers, accompanied with an
// initialized context that can be used for the RPC call, caller must call
// Close() on the client.
func NewClientWithDefaultContextAndAddress() (*ClientWithContext, error) {
	return NewClient(context.Background(), ResolveServerAddress(), Timeout)
}

func NewClient(ctx context.Context, address string, timeout time.Duration) (*ClientWithContext, error) {
	c := &ClientWithContext{}

	c.SignalCtx, c.signalCancel = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	c.Ctx, c.timeoutCancel = context.WithTimeout(c.SignalCtx, timeout)

	var err error
	c.conn, err = grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(retryPolicy(Retries)),
		grpc.WithMaxCallAttempts(Retries+1), // maxAttempt includes the first call
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client with address %s: %w", address, err)
	}
	c.Client = tetragon.NewFineGuidanceSensorsClient(c.conn)

	return c, nil
}

func getRequest(includeFields, excludeFields []string, filter *tetragon.Filter) *tetragon.GetEventsRequest {
	var fieldFilters []*tetragon.FieldFilter
	if len(includeFields) > 0 {
		fieldFilters = append(fieldFilters, &tetragon.FieldFilter{
			EventSet: []tetragon.EventType{},
			Fields: &fieldmaskpb.FieldMask{
				Paths: includeFields,
			},
			Action: tetragon.FieldFilterAction_INCLUDE,
		})
	}
	if len(excludeFields) > 0 {
		fieldFilters = append(fieldFilters, &tetragon.FieldFilter{
			EventSet: []tetragon.EventType{},
			Fields: &fieldmaskpb.FieldMask{
				Paths: excludeFields,
			},
			Action: tetragon.FieldFilterAction_EXCLUDE,
		})
	}

	return &tetragon.GetEventsRequest{
		FieldFilters: fieldFilters,
		AllowList:    []*tetragon.Filter{filter},
	}
}

var GetFilter = func() *tetragon.Filter {
	if Options.Host {
		// Host events can be matched by an empty namespace string.
		Options.Namespaces = append(Options.Namespaces, "")
	}
	// Only set these filters if they are not empty. We currently rely on Protobuf to
	// marshal empty lists as nil for filters to function properly. It doesn't work with
	// stdin mode since it doesn't go over the wire, causing all events to get filtered
	// out because empty allowlist does not match anything.
	filter := tetragon.Filter{}
	if len(Options.Processes) > 0 {
		filter.BinaryRegex = Options.Processes
	}
	if len(Options.Namespaces) > 0 {
		filter.Namespace = Options.Namespaces
	}
	if len(Options.Pods) > 0 {
		filter.PodRegex = Options.Pods
	}
	// Is used to filter on the event types i.e. PROCESS_EXEC, PROCESS_EXIT etc.
	if len(Options.EventTypes) > 0 {
		var eventType tetragon.EventType

		for _, v := range Options.EventTypes {
			eventType = tetragon.EventType(tetragon.EventType_value[v])
			filter.EventSet = append(filter.EventSet, eventType)
		}
	}
	if len(Options.PolicyNames) > 0 {
		filter.PolicyNames = Options.PolicyNames
	}

	return &filter
}

func ListenToEvents() error {
	c, err := NewClientWithDefaultContextAndAddress()
	if err != nil {
		fmt.Printf("failed to create gRPC client: %v\n", err)
		return err
	}
	defer c.Close()

	request := getRequest(Options.IncludeFields, Options.ExcludeFields, GetFilter())

	stream, err := c.Client.GetEvents(c.Ctx, request)

	if err != nil {
		fmt.Errorf("failed to call GetEvents: %v", err)
		return err
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			if !errors.Is(err, context.Canceled) && status.Code(err) != codes.Canceled && !errors.Is(err, io.EOF) {
				fmt.Errorf("failed to receive events: %w", err)
			}
			return err
		}

		println(res.GetEvent())

	}

}

func ResolveServerAddress() string {
	if ServerAddress == "" {
		return "localhost:54321"
	}

	return ServerAddress
}
