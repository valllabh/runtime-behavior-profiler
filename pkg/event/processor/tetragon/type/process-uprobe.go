package eventprocessortetragontype

import (
	eventtype "runtime-behavior-profiler/pkg/event/type"

	"github.com/cilium/tetragon/api/v1/tetragon"
)

type ProcessUprobe struct {
	*tetragon.ProcessUprobe
}

// IsHostEvent checks if the event is a host event.
func (e *ProcessUprobe) IsHostEvent() bool {
	return IsHostEvent(e.Process)
}

// GetContainer implements eventtype.IEvent.
func (e *ProcessUprobe) GetContainer() (*eventtype.Container, error) {
	return GetContainer(e.Process)
}

// GetNamespace implements eventtype.IEvent.
func (e *ProcessUprobe) GetNamespace() (*eventtype.Namespace, error) {
	return GetNamespace(e.Process)
}

// GetParentProcess implements eventtype.IEvent.
func (e *ProcessUprobe) GetParentProcess() (*eventtype.Process, error) {
	return GetProcess(e.Process)
}

// GetPod implements eventtype.IEvent.
func (e *ProcessUprobe) GetPod() (*eventtype.Pod, error) {
	return GetPod(e.Process)
}

// GetProcess implements eventtype.IEvent.
func (e *ProcessUprobe) GetProcess() (*eventtype.Process, error) {
	return GetProcess(e.Process)
}
