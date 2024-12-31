package eventprocessortetragontype

import (
	"context"

	"github.com/cilium/tetragon/api/v1/tetragon"
	"google.golang.org/grpc"
)

type TetragonEvent struct {
}

type TetragonEventListerOptions struct {
	Output        string
	IncludeFields []string
	EventTypes    []string
	ExcludeFields []string
	Namespaces    []string
	Processes     []string
	Pods          []string
	Host          bool
	Timestamps    bool
	TTYEncode     string
	StackTraces   bool
	ImaHash       bool
	PolicyNames   []string

	Debug         bool
	ServerAddress string
	Retries       int
}

type ClientWithContext struct {
	Client tetragon.FineGuidanceSensorsClient
	// Ctx is a combination of the signal context and the timeout context
	Ctx context.Context
	// SignalCtx is only the signal context, you might want to use that context
	// when the command should never timeout (like a stream command)
	// SignalCtx context.Context
	Conn *grpc.ClientConn
	// The signal context is the parent of the timeout context, so cancelling
	// signal will cancel its child, timeout
	Cancel context.CancelFunc
	// TimeoutCancel context.CancelFunc
}
