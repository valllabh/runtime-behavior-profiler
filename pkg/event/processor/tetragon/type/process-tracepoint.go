package eventprocessortetragontype

import (
	eventtype "runtime-behavior-profiler/pkg/event/type"

	"github.com/cilium/tetragon/api/v1/tetragon"
)

type ProcessTracepoint struct {
	*tetragon.ProcessTracepoint
}

// IsHostEvent checks if the event is a host event.
func (e *ProcessTracepoint) IsHostEvent() bool {
	return IsHostEvent(e.Process)
}

// GetContainer implements eventtype.IEvent.
func (e *ProcessTracepoint) GetContainer() (*eventtype.Container, error) {
	return GetContainer(e.Process)
}

// GetNamespace implements eventtype.IEvent.
func (e *ProcessTracepoint) GetNamespace() (*eventtype.Namespace, error) {
	return GetNamespace(e.Process)
}

// GetParentProcess implements eventtype.IEvent.
func (e *ProcessTracepoint) GetParentProcess() (*eventtype.Process, error) {
	return GetProcess(e.Process)
}

// GetPod implements eventtype.IEvent.
func (e *ProcessTracepoint) GetPod() (*eventtype.Pod, error) {
	return GetPod(e.Process)
}

// GetProcess implements eventtype.IEvent.
func (e *ProcessTracepoint) GetProcess() (*eventtype.Process, error) {
	return GetProcess(e.Process)
}
