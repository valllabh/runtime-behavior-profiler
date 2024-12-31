package eventprocessortetragon

import (
	eventprocessortetragontype "runtime-behavior-profiler/pkg/event/processor/tetragon/type"
	eventtype "runtime-behavior-profiler/pkg/event/type"

	"github.com/cilium/tetragon/api/v1/tetragon"
)

func ProcessProcessExec(e *tetragon.ProcessExec) eventtype.IEvent {
	if eventprocessortetragontype.IsHostEvent(e.Process) {
		return nil
	}

	return &eventprocessortetragontype.ProcessExec{
		ProcessExec: e,
	}
}

func ProcessProcessExit(e *tetragon.ProcessExit) eventtype.IEvent {
	if eventprocessortetragontype.IsHostEvent(e.Process) {
		return nil
	}

	return &eventprocessortetragontype.ProcessExit{
		ProcessExit: e,
	}
}

func ProcessProcessLoader(e *tetragon.ProcessLoader) eventtype.IEvent {
	if eventprocessortetragontype.IsHostEvent(e.Process) {
		return nil
	}

	return &eventprocessortetragontype.ProcessLoader{
		ProcessLoader: e,
	}
}

func ProcessProcessKprobe(e *tetragon.ProcessKprobe) eventtype.IEvent {
	if eventprocessortetragontype.IsHostEvent(e.Process) {
		return nil
	}

	return &eventprocessortetragontype.ProcessKprobe{
		ProcessKprobe: e,
	}
}

func ProcessProcessTracepoint(e *tetragon.ProcessTracepoint) eventtype.IEvent {
	if eventprocessortetragontype.IsHostEvent(e.Process) {
		return nil
	}

	return &eventprocessortetragontype.ProcessTracepoint{
		ProcessTracepoint: e,
	}
}

func ProcessProcessUprobe(e *tetragon.ProcessUprobe) eventtype.IEvent {
	if eventprocessortetragontype.IsHostEvent(e.Process) {
		return nil
	}

	return &eventprocessortetragontype.ProcessUprobe{
		ProcessUprobe: e,
	}
}

func ProcessProcessLsm(e *tetragon.ProcessLsm) eventtype.IEvent {
	if eventprocessortetragontype.IsHostEvent(e.Process) {
		return nil
	}

	return &eventprocessortetragontype.ProcessLsm{
		ProcessLsm: e,
	}
}
