package eventprocessortetragontype

type TetragonEvent struct {
}

type Opts struct {
	Output        string
	Color         string
	IncludeFields []string
	EventTypes    []string
	ExcludeFields []string
	Namespaces    []string
	Processes     []string
	Pods          []string
	Host          bool
	Timestamps    bool
	TTYEncode     string
	StackTraces   bool
	ImaHash       bool
	PolicyNames   []string
}
