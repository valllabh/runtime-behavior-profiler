package eventprocessortetragontype

import (
	eventtype "runtime-behavior-profiler/pkg/event/type"

	"github.com/cilium/tetragon/api/v1/tetragon"
)

type ProcessExit struct {
	*tetragon.ProcessExit
}

// IsHostEvent checks if the event is a host event.
func (e *ProcessExit) IsHostEvent() bool {
	return IsHostEvent(e.Process)
}

// GetContainer implements eventtype.IEvent.
func (e *ProcessExit) GetContainer() (*eventtype.Container, error) {
	return GetContainer(e.Process)
}

// GetNamespace implements eventtype.IEvent.
func (e *ProcessExit) GetNamespace() (*eventtype.Namespace, error) {
	return GetNamespace(e.Process)
}

// GetParentProcess implements eventtype.IEvent.
func (e *ProcessExit) GetParentProcess() (*eventtype.Process, error) {
	return GetProcess(e.Process)
}

// GetPod implements eventtype.IEvent.
func (e *ProcessExit) GetPod() (*eventtype.Pod, error) {
	return GetPod(e.Process)
}

// GetProcess implements eventtype.IEvent.
func (e *ProcessExit) GetProcess() (*eventtype.Process, error) {
	return GetProcess(e.Process)
}
