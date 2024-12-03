package eventtype

import (
	"errors"
)

// func to get namespace from Event
func (e *Event) GetNamespace() (string, error) {
	process, err := e.GetProcess()
	if err != nil {
		return "", err
	}
	return process.Pod.Namespace, nil
}

// func to get pod from Event
func (e *Event) GetPod() (Pod, error) {
	process, err := e.GetProcess()
	if err != nil {
		return Pod{}, err
	}
	return process.Pod, nil
}

// func to get container from Event
func (e *Event) GetContainer() (Container, error) {
	process, err := e.GetProcess()
	if err != nil {
		return Container{}, err
	}
	return process.Pod.Container, nil
}

// func to get Process from Event
func (e *Event) GetProcess() (Process, error) {

	// for process kprobe event
	if e.ProcessKprobe.Process.ExecID != "" {
		return e.ProcessKprobe.Process, nil
	}

	// for process exit event
	if e.ProcessExit.Process.ExecID != "" {
		return e.ProcessExit.Process, nil
	}

	// for process exec event
	if e.ProcessExec.Process.ExecID != "" {
		return e.ProcessExec.Process, nil
	}

	// return error if no process found
	return Process{}, errors.New("no process found in event")
}

// func to get Parent from Event
func (e *Event) GetParent() (Process, error) {
	// for process kprobe event
	if e.ProcessKprobe.Process.ExecID != "" {
		return e.ProcessKprobe.Parent, nil
	}

	// for process exit event
	if e.ProcessExit.Process.ExecID != "" {
		return e.ProcessExit.Parent, nil
	}

	// for process exec event
	if e.ProcessExec.Process.ExecID != "" {
		return e.ProcessExec.Parent, nil
	}

	return Process{}, errors.New("no parent process found in event")
}
