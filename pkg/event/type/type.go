package eventtype

import (
	"github.com/valllabh/ocsf-schema-golang/ocsf/v1_0_0/events/network"
	"github.com/valllabh/ocsf-schema-golang/ocsf/v1_0_0/events/system"
)

// These structs are shared across all event types.

type Event struct {
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

type Process struct {
	Binary    string `json:"binary"`
	Arguments string `json:"arguments"`
}

type Namespace struct {
	Name string `json:"name"`
}

type Pod struct {
	Name string `json:"name"`
}

type Container struct {
	Name  string `json:"name"`
	Image Image  `json:"image"`
}

type Image struct {
	Name string `json:"name"`
}
