package eventtype

import (
	"fmt"

	"github.com/valllabh/ocsf-schema-golang/ocsf/v1_0_0/events/network"
	"github.com/valllabh/ocsf-schema-golang/ocsf/v1_0_0/events/system"
	"github.com/valllabh/ocsf-schema-golang/ocsf/v1_0_0/objects"
)

type OCSFEvent struct {
	OCSF_1_0_0 *OCSF_1_0_0 `json:"ocsf_1_0_0,omitempty"`
}

type OCSF_1_0_0 struct {
	FileActivity    *FileActivity    `json:"file_activity,omitempty"`
	NetworkActivity *NetworkActivity `json:"network_activity,omitempty"`
	ProcessActivity *ProcessActivity `json:"process_activity,omitempty"`
}

type FileActivity struct {
	system.FileActivity
	Resources []*Resource `json:"resources"`
}

type NetworkActivity struct {
	network.NetworkActivity
	Resources []*Resource `json:"resources"`
}

type ProcessActivity struct {
	system.ProcessActivity
	Resources []*Resource `json:"resources"`
}

type Resource struct {
	Type string `json:"type"`
	Name string `json:"name"`
	UID  string `json:"uid"`
}

// Functions

// func to get namespace from Event
func (e *OCSFEvent) GetNamespace() (*Namespace, error) {

	resource, err := e.getResource("kubernetes.namespace")

	if err != nil {
		return nil, err
	}

	return &Namespace{
		Name: resource.Name,
		Pods: map[string]*Pod{},
	}, nil
}

// func to get pod from Event
func (e *OCSFEvent) GetPod() (*Pod, error) {
	resource, err := e.getResource("kubernetes.pod")

	if err != nil {
		return nil, err
	}

	return &Pod{
		Name:       resource.Name,
		Containers: map[string]*Container{},
	}, nil
}

// func to get container from Event
func (e *OCSFEvent) GetContainer() (*Container, error) {
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
		return nil, fmt.Errorf("container not found")
	}

	return &Container{
		Name: container.Name,
		Image: &Image{
			Repo: container.Image.Name,
		},
		Processes: map[string]*Process{},
	}, nil
}

// func to get Process from Event
func (e *OCSFEvent) GetProcess() (*Process, error) {
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
		return nil, fmt.Errorf("process not found")
	}

	return &Process{
		Binary:         process.Name,
		Arguments:      process.CmdLine,
		ChildProcesses: map[string]*Process{},
	}, nil
}

// func to get Parent from Event
func (e *OCSFEvent) GetParentProcess() (*Process, error) {
	var process *objects.Process
	var tyepee string = e.GetType()

	if tyepee == "" {
		return nil, fmt.Errorf("event type not found")
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
		return &Process{
			Binary: "root",
		}, nil
	}

	return &Process{
		Binary:         process.Name,
		Arguments:      process.CmdLine,
		ChildProcesses: map[string]*Process{},
	}, nil
}

func (e *OCSFEvent) GetType() string {
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

func (e *OCSFEvent) getResource(t string) (*Resource, error) {
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

func (e *OCSFEvent) FileActivityGetResource(t string) (*Resource, error) {
	for _, resource := range e.OCSF_1_0_0.FileActivity.Resources {
		if resource.Type == t {
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource of type %s not found", t)
}

func (e *OCSFEvent) NetworkActivityGetResource(t string) (*Resource, error) {
	for _, resource := range e.OCSF_1_0_0.NetworkActivity.Resources {
		if resource.Type == t {
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource of type %s not found", t)
}

func (e *OCSFEvent) ProcessActivityGetResource(t string) (*Resource, error) {
	for _, resource := range e.OCSF_1_0_0.ProcessActivity.Resources {
		if resource.Type == t {
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource of type %s not found", t)
}
