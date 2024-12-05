package eventtype

import (
	"fmt"

	"github.com/valllabh/ocsf-schema-golang/ocsf/v1_0_0/objects"
)

// func to get namespace from Event
func (e *Event) GetNamespace() (Namespace, error) {

	resource, err := e.getResource("kubernetes.namespace")

	if err != nil {
		return Namespace{}, err
	}

	return Namespace{
		Name: resource.Name,
	}, nil
}

// func to get pod from Event
func (e *Event) GetPod() (Pod, error) {
	resource, err := e.getResource("kubernetes.pod")

	if err != nil {
		return Pod{}, err
	}

	return Pod{
		Name: resource.Name,
	}, nil
}

// func to get container from Event
func (e *Event) GetContainer() (Container, error) {
	var container *objects.Container
	switch e.GetType() {
	case "FILE_EVENT":
		container = e.OCSF_1_0_0.FileActivity.Actor.Process.Container
	case "NETWORK_EVENT":
		container = e.OCSF_1_0_0.NetworkActivity.Actor.Process.Container
	case "PROCESS_EVENT":
		container = e.OCSF_1_0_0.ProcessActivity.Actor.Process.Container
	}

	if container == nil {
		return Container{}, fmt.Errorf("container not found")
	}

	return Container{
		Name:  container.Name,
		Image: Image{Name: container.Image.Name},
	}, nil
}

// func to get Process from Event
func (e *Event) GetProcess() (Process, error) {
	var process *objects.Process
	switch e.GetType() {
	case "FILE_EVENT":
		process = e.OCSF_1_0_0.FileActivity.Actor.Process
	case "NETWORK_EVENT":
		process = e.OCSF_1_0_0.NetworkActivity.Actor.Process
	case "PROCESS_EVENT":
		process = e.OCSF_1_0_0.ProcessActivity.Actor.Process

	}

	if process == nil {
		return Process{}, fmt.Errorf("process not found")
	}

	return Process{
		Binary:    process.Name,
		Arguments: process.CmdLine,
	}, nil
}

// func to get Parent from Event
func (e *Event) GetParent() (Process, error) {
	var process *objects.Process
	var tyepee string = e.GetType()

	if tyepee == "" {
		return Process{}, fmt.Errorf("event type not found")
	}

	switch tyepee {
	case "FILE_EVENT":
		process = e.OCSF_1_0_0.FileActivity.Actor.Process.ParentProcess
	case "NETWORK_EVENT":
		process = e.OCSF_1_0_0.NetworkActivity.Actor.Process.ParentProcess
	case "PROCESS_EVENT":
		process = e.OCSF_1_0_0.ProcessActivity.Actor.Process.ParentProcess
	}

	if process == nil {
		return Process{
			Binary: "root",
		}, nil
	}

	return Process{
		Binary:    process.Name,
		Arguments: process.CmdLine,
	}, nil
}

func (e *Event) GetType() string {
	if e.OCSF_1_0_0 == nil {
		return ""
	}

	if e.OCSF_1_0_0.FileActivity != nil {
		return "FILE_EVENT"
	}

	if e.OCSF_1_0_0.NetworkActivity != nil {
		return "NETWORK_EVENT"
	}

	if e.OCSF_1_0_0.ProcessActivity != nil {
		return "PROCESS_EVENT"
	}

	return ""
}

func (e *Event) getResource(t string) (*Resource, error) {
	switch e.GetType() {
	case "FILE_EVENT":
		return e.FileActivityGetResource(t)

	case "NETWORK_EVENT":
		return e.NetworkActivityGetResource(t)

	case "PROCESS_EVENT":
		return e.ProcessActivityGetResource(t)
	}

	return nil, fmt.Errorf("resource of type %s not found", t)
}

func (e *Event) FileActivityGetResource(t string) (*Resource, error) {
	for _, resource := range e.OCSF_1_0_0.FileActivity.Resources {
		if resource.Type == t {
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource of type %s not found", t)
}

func (e *Event) NetworkActivityGetResource(t string) (*Resource, error) {
	for _, resource := range e.OCSF_1_0_0.NetworkActivity.Resources {
		if resource.Type == t {
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource of type %s not found", t)
}

func (e *Event) ProcessActivityGetResource(t string) (*Resource, error) {
	for _, resource := range e.OCSF_1_0_0.ProcessActivity.Resources {
		if resource.Type == t {
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource of type %s not found", t)
}
