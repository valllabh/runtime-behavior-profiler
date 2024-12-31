package eventprocessortetragontype

import (
	eventtype "runtime-behavior-profiler/pkg/event/type"

	"github.com/cilium/tetragon/api/v1/tetragon"
)

type ProcessExec struct {
	*tetragon.ProcessExec
}

// IsHostEvent checks if the event is a host event.
func (e *ProcessExec) IsHostEvent() bool {
	return IsHostEvent(e.Process)
}

// GetContainer implements eventtype.IEvent.
func (e *ProcessExec) GetContainer() (*eventtype.Container, error) {
	return GetContainer(e.Process)
}

// GetNamespace implements eventtype.IEvent.
func (e *ProcessExec) GetNamespace() (*eventtype.Namespace, error) {
	return GetNamespace(e.Process)
}

// GetParentProcess implements eventtype.IEvent.
func (e *ProcessExec) GetParentProcess() (*eventtype.Process, error) {
	return GetProcess(e.Process)
}

// GetPod implements eventtype.IEvent.
func (e *ProcessExec) GetPod() (*eventtype.Pod, error) {
	return GetPod(e.Process)
}

// GetProcess implements eventtype.IEvent.
func (e *ProcessExec) GetProcess() (*eventtype.Process, error) {
	return GetProcess(e.Process)
}
