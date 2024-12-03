package clusterbehaviourprofiletype

type ClusterBehaviourProfile struct {
	Cluster    string               `json:"cluster"`
	Namespaces map[string]Namespace `json:"namespaces"`
}

type Namespace struct {
	Name string         `json:"name"`
	Pods map[string]Pod `json:"pods"`
}

type Pod struct {
	Name       string               `json:"name"`
	Containers map[string]Container `json:"containers"`
}

type Container struct {
	Name      string             `json:"name"`
	Image     Image              `json:"image"`
	Processes map[string]Process `json:"processes"`
}

type Image struct {
	Repo     string   `json:"repo"`
	Tag      string   `json:"tag"`
	Registry Registry `json:"registry"`
}

type Registry struct {
	Name string `json:"name"`
}

type Process struct {
	Binary         string             `json:"binary"`
	Arguments      string             `json:"arguments"`
	ChildProcesses map[string]Process `json:"child_processes"`
}
