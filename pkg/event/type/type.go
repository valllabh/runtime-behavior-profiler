package eventtype

type SinkOperation string

const (
	SinkOperationIgnored  SinkOperation = "IGNORED"
	SinkOperationInserted SinkOperation = "INSERTED"
	SinkOperationUpdated  SinkOperation = "UPDATED"
)

type IEvent interface {
	GetNamespace() (*Namespace, error)
	GetPod() (*Pod, error)
	GetContainer() (*Container, error)
	GetParentProcess() (*Process, error)
	GetProcess() (*Process, error)
}

type SinkResult struct {
	Operation SinkOperation `json:"operation"`
	Path      []string      `json:"path"`
}

type Cluster struct {
	Name       string                `json:"name"`
	Namespaces map[string]*Namespace `json:"namespaces"`
}

type Namespace struct {
	Name string          `json:"name"`
	Pods map[string]*Pod `json:"pods"`
}

type Pod struct {
	Name       string                `json:"name"`
	Containers map[string]*Container `json:"containers"`
}

type Container struct {
	Name      string              `json:"name"`
	Image     *Image              `json:"image"`
	Processes map[string]*Process `json:"processes"`
}

type Image struct {
	Repo     string    `json:"repo"`
	Tag      string    `json:"tag"`
	Registry *Registry `json:"registry"`
}

type Registry struct {
	Name string `json:"name"`
}

type Process struct {
	Binary         string              `json:"binary"`
	Arguments      string              `json:"arguments"`
	ChildProcesses map[string]*Process `json:"child_processes"`
}
