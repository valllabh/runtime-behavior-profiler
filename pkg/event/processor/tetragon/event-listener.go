package eventprocessortetragon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/signal"
	eventprocessortetragontype "runtime-behavior-profiler/pkg/event/processor/tetragon/type"
	eventtype "runtime-behavior-profiler/pkg/event/type"
	"syscall"

	"github.com/cilium/tetragon/api/v1/tetragon"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type tetragonEventListener struct {
	Options               eventprocessortetragontype.TetragonEventListerOptions
	Cluster               *eventtype.Cluster
	GRPCClientWithContext *eventprocessortetragontype.ClientWithContext
}

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
func (tel *tetragonEventListener) retryPolicy() string {
	if tel.Options.Retries < 0 {
		// gRPC should ignore the invalid retry policy but will issue a warning,
		return "{}"
	}
	// maxAttempt includes the first call
	maxAttempt := tel.Options.Retries + 1
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

// initGRPCClientWithContext returns a client to a tetragon
// server after resolving the server address using helpers, accompanied with an
// initialized context that can be used for the RPC call, caller must call
// Close() on the client.
func (tel *tetragonEventListener) initGRPCClientWithContext(ctx context.Context) error {
	tel.GRPCClientWithContext = &eventprocessortetragontype.ClientWithContext{}

	tel.GRPCClientWithContext.Ctx, tel.GRPCClientWithContext.Cancel = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	var err error
	tel.GRPCClientWithContext.Conn, err = grpc.NewClient(tel.Options.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // no TLS
		grpc.WithDefaultServiceConfig(tel.retryPolicy()),         // gRPC A6 - gRPC Retry Design
		grpc.WithMaxCallAttempts(tel.Options.Retries+1),          // maxAttempt includes the first call
	)

	if err != nil {
		return fmt.Errorf("failed to create gRPC client with address %s: %w", tel.Options.ServerAddress, err)
	}

	tel.GRPCClientWithContext.Client = tetragon.NewFineGuidanceSensorsClient(tel.GRPCClientWithContext.Conn)

	return nil
}

func (tel *tetragonEventListener) getRequest() *tetragon.GetEventsRequest {
	var fieldFilters []*tetragon.FieldFilter
	if len(tel.Options.IncludeFields) > 0 {
		fieldFilters = append(fieldFilters, &tetragon.FieldFilter{
			EventSet: []tetragon.EventType{},
			Fields: &fieldmaskpb.FieldMask{
				Paths: tel.Options.IncludeFields,
			},
			Action: tetragon.FieldFilterAction_INCLUDE,
		})
	}
	if len(tel.Options.ExcludeFields) > 0 {
		fieldFilters = append(fieldFilters, &tetragon.FieldFilter{
			EventSet: []tetragon.EventType{},
			Fields: &fieldmaskpb.FieldMask{
				Paths: tel.Options.ExcludeFields,
			},
			Action: tetragon.FieldFilterAction_EXCLUDE,
		})
	}

	return &tetragon.GetEventsRequest{
		FieldFilters: fieldFilters,
		AllowList:    []*tetragon.Filter{tel.getFilter()},
	}
}

func (tel *tetragonEventListener) getFilter() *tetragon.Filter {
	if tel.Options.Host {
		// Host events can be matched by an empty namespace string.
		tel.Options.Namespaces = append(tel.Options.Namespaces, "")
	}
	// Only set these filters if they are not empty. We currently rely on Protobuf to
	// marshal empty lists as nil for filters to function properly. It doesn't work with
	// stdin mode since it doesn't go over the wire, causing all events to get filtered
	// out because empty allowlist does not match anything.
	filter := tetragon.Filter{}
	if len(tel.Options.Processes) > 0 {
		filter.BinaryRegex = tel.Options.Processes
	}
	if len(tel.Options.Namespaces) > 0 {
		filter.Namespace = tel.Options.Namespaces
	}
	if len(tel.Options.Pods) > 0 {
		filter.PodRegex = tel.Options.Pods
	}
	// Is used to filter on the event types i.e. PROCESS_EXEC, PROCESS_EXIT etc.
	if len(tel.Options.EventTypes) > 0 {
		var eventType tetragon.EventType

		for _, v := range tel.Options.EventTypes {
			eventType = tetragon.EventType(tetragon.EventType_value[v])
			filter.EventSet = append(filter.EventSet, eventType)
		}
	}
	if len(tel.Options.PolicyNames) > 0 {
		filter.PolicyNames = tel.Options.PolicyNames
	}

	return &filter
}
func (tel *tetragonEventListener) OnEndListeningEvent() {
	println("Tetragon Event Listener is done listening to events")

	json, _ := json.MarshalIndent(tel.Cluster, "", "  ")
	println("\n" + string(json))

	tel.GRPCClientWithContext.Conn.Close()
	tel.GRPCClientWithContext.Cancel()
}

func (tel *tetragonEventListener) ListenToEvents() error {

	fmt.Printf("Listening to tetragon events on %s\n", tel.Options.ServerAddress)

	err := tel.initGRPCClientWithContext(context.Background())

	if err != nil {
		fmt.Printf("failed to create gRPC client: %v\n", err)
		return err
	}
	defer tel.OnEndListeningEvent()

	request := tel.getRequest()

	stream, err := tel.GRPCClientWithContext.Client.GetEvents(tel.GRPCClientWithContext.Ctx, request)

	if err != nil {
		fmt.Errorf("failed to call GetEvents: %v", err)
		return err
	}

	for {
		response, err := stream.Recv()
		if err != nil {
			if !errors.Is(err, context.Canceled) && status.Code(err) != codes.Canceled && !errors.Is(err, io.EOF) {
				fmt.Errorf("failed to receive events: %w", err)
			}
			return err // if not returned will go in infinite loop
		}
		eventType := response.EventType()

		var iEvent eventtype.IEvent

		switch eventType {
		case tetragon.EventType_PROCESS_EXEC:
			iEvent = ProcessProcessExec(response.GetProcessExec())
		case tetragon.EventType_PROCESS_EXIT:
			iEvent = ProcessProcessExit(response.GetProcessExit())
		case tetragon.EventType_PROCESS_LOADER:
			iEvent = ProcessProcessLoader(response.GetProcessLoader())
		case tetragon.EventType_PROCESS_KPROBE:
			iEvent = ProcessProcessKprobe(response.GetProcessKprobe())
		case tetragon.EventType_PROCESS_TRACEPOINT:
			iEvent = ProcessProcessTracepoint(response.GetProcessTracepoint())
		case tetragon.EventType_PROCESS_UPROBE:
			iEvent = ProcessProcessUprobe(response.GetProcessUprobe())
		case tetragon.EventType_PROCESS_LSM:
			iEvent = ProcessProcessLsm(response.GetProcessLsm())
		}

		tel.Cluster.SinkEvent(iEvent)

	}

}

func GetDefaultOptions() eventprocessortetragontype.TetragonEventListerOptions {
	return eventprocessortetragontype.TetragonEventListerOptions{
		Retries:       5,
		ServerAddress: "localhost:54321",
		Namespaces:    []string{},
	}
}

func NewEventListener(cluster *eventtype.Cluster) *tetragonEventListener {
	return NewEventListenerWithOptions(cluster, GetDefaultOptions())
}

func NewEventListenerWithOptions(cluster *eventtype.Cluster, options eventprocessortetragontype.TetragonEventListerOptions) *tetragonEventListener {
	return &tetragonEventListener{
		Options: options,
		Cluster: cluster,
	}
}
