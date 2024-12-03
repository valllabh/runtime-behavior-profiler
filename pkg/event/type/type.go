package eventtype

import "time"

// These structs are shared across all event types.

type Event struct {
	ProcessKprobe ProcessKprobe `json:"process_kprobe,omitempty"`
	ProcessExit   ProcessExit   `json:"process_exit,omitempty"`
	ProcessExec   ProcessExec   `json:"process_exec,omitempty"`
	NodeName      string        `json:"node_name"`
	Time          time.Time     `json:"time"`
}

type Process struct {
	ExecID       string    `json:"exec_id"`
	PID          int       `json:"pid"`
	UID          int       `json:"uid"`
	CWD          string    `json:"cwd"`
	Binary       string    `json:"binary"`
	Arguments    string    `json:"arguments"`
	Flags        string    `json:"flags"`
	StartTime    time.Time `json:"start_time"`
	AUID         int       `json:"auid"`
	Pod          Pod       `json:"pod"`
	Docker       string    `json:"docker"`
	ParentExecID string    `json:"parent_exec_id"`
	Refcnt       int       `json:"refcnt,omitempty"`
}

type Pod struct {
	Namespace string                 `json:"namespace"`
	Name      string                 `json:"name"`
	Container Container              `json:"container"`
	PodLabels map[string]interface{} `json:"pod_labels,omitempty"`
}

type Container struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Image     Image     `json:"image"`
	StartTime time.Time `json:"start_time"`
	PID       int       `json:"pid"`
}

type Image struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PROCESSKPROBE STRUCTS
// These structs are specific to ProcessKprobe events.

type ProcessKprobe struct {
	Process      Process     `json:"process"`
	Parent       Process     `json:"parent"`
	FunctionName string      `json:"function_name"`
	Args         []Argument  `json:"args"`
	Return       ReturnValue `json:"return"`
	Action       string      `json:"action"`
}

type Argument struct {
	StringArg string `json:"string_arg"`
}

type ReturnValue struct {
	IntArg int `json:"int_arg"`
}

// PROCESSEXIT STRUCTS
// These structs are specific to ProcessExit events.

type ProcessExit struct {
	Process Process `json:"process"`
	Parent  Process `json:"parent"`
}

// PROCESSEXEC STRUCTS
// These structs are specific to ProcessExec events.

type ProcessExec struct {
	Process Process `json:"process"`
	Parent  Process `json:"parent"`
}
