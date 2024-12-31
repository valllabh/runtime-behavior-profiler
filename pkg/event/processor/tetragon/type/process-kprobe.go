package eventprocessortetragontype

import (
	eventtype "runtime-behavior-profiler/pkg/event/type"

	"github.com/cilium/tetragon/api/v1/tetragon"
)

type ProcessKprobe struct {
	*tetragon.ProcessKprobe
}

// IsHostEvent checks if the event is a host event.
func (e *ProcessKprobe) IsHostEvent() bool {
	return IsHostEvent(e.Process)
}

// GetContainer implements eventtype.IEvent.
func (e *ProcessKprobe) GetContainer() (*eventtype.Container, error) {
	return GetContainer(e.Process)
}

// GetNamespace implements eventtype.IEvent.
func (e *ProcessKprobe) GetNamespace() (*eventtype.Namespace, error) {
	return GetNamespace(e.Process)
}

// GetParentProcess implements eventtype.IEvent.
func (e *ProcessKprobe) GetParentProcess() (*eventtype.Process, error) {
	return GetProcess(e.Process)
}

// GetPod implements eventtype.IEvent.
func (e *ProcessKprobe) GetPod() (*eventtype.Pod, error) {
	return GetPod(e.Process)
}

// GetProcess implements eventtype.IEvent.
func (e *ProcessKprobe) GetProcess() (*eventtype.Process, error) {
	return GetProcess(e.Process)
}
