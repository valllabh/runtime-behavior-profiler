package eventprocessortetragontype

import (
	eventtype "runtime-behavior-profiler/pkg/event/type"

	"github.com/cilium/tetragon/api/v1/tetragon"
)

// IsHostEvent checks if the event is a host event.
func IsHostEvent(process *tetragon.Process) bool {
	return process.Pod == nil
}

// GetContainer implements eventtype.IEvent.
func GetContainer(process *tetragon.Process) (*eventtype.Container, error) {
	return &eventtype.Container{
		Name: process.Pod.Container.Name,
		Image: &eventtype.Image{
			Repo: process.Pod.Container.Image.Name,
		},
		Processes: map[string]*eventtype.Process{},
	}, nil
}

// GetNamespace implements eventtype.IEvent.
func GetNamespace(process *tetragon.Process) (*eventtype.Namespace, error) {
	return &eventtype.Namespace{
		Name: process.Pod.Namespace,
		Pods: map[string]*eventtype.Pod{},
	}, nil
}

// GetPod implements eventtype.IEvent.
func GetPod(process *tetragon.Process) (*eventtype.Pod, error) {
	return &eventtype.Pod{
		Name:       process.Pod.Name,
		Containers: map[string]*eventtype.Container{},
	}, nil
}

// GetProcess implements eventtype.IEvent.
func GetProcess(process *tetragon.Process) (*eventtype.Process, error) {
	return &eventtype.Process{
		Binary:         process.Binary,
		Arguments:      process.Arguments,
		ChildProcesses: map[string]*eventtype.Process{},
	}, nil
}
