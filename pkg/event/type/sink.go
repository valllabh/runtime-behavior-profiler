package eventtype

import (
	"log"
	"runtime-behavior-profiler/pkg/util"
	"time"

	"github.com/gofrs/uuid"
)

const (
	namespaceType = "namespace"
	podType       = "pod"
	containerType = "container"
	processType   = "process"
)

// SinkOperation represents the operation performed by the SinkEvent function.
// It can be either SinkOperationInserted or SinkOperationUpdated.
// SinkOperationInserted indicates that a new entity was inserted into the Cluster.
// SinkOperationUpdated indicates that an existing entity was updated in the Cluster.
// The SinkResult struct contains the operation and the path to the entity that was inserted or updated.
func (cluster *Cluster) SinkEvent(rawEvent IEvent) (*SinkResult, error) {

	// If the raw event is nil, return an SinkOperationIgnored.
	if rawEvent == nil {
		return &SinkResult{
			Operation: SinkOperationIgnored,
			Path:      []string{},
		}, nil
	}

	startTime := time.Now()

	sinkResult := SinkResult{
		Operation: SinkOperationUpdated,
		Path:      []string{},
	}

	// Namespace
	namespaceRaw, err := rawEvent.GetNamespace()
	if err != nil {
		return nil, err
	}
	namespaceKey := namespaceRaw.GetKey()

	namespace, ok := cluster.Namespaces[namespaceKey]

	if !ok {
		namespace = &Namespace{
			Name: namespaceRaw.Name,
			Pods: map[string]*Pod{},
		}
		sinkResult.Inserted(namespace.GetKey())
		cluster.Namespaces[namespaceKey] = namespace
	}

	// Pod
	podRaw, err := rawEvent.GetPod()
	if err != nil {
		return nil, err
	}
	podKey := podRaw.GetKey()

	pod, ok := namespace.Pods[podKey]
	if !ok {
		pod = &Pod{
			Name:       podRaw.GetName(),
			Containers: map[string]*Container{},
		}
		sinkResult.Inserted(pod.GetKey())
		namespace.Pods[podKey] = pod
	}

	// Container
	containerRaw, err := rawEvent.GetContainer()
	if err != nil {
		return nil, err
	}
	containerKey := containerRaw.GetKey()

	container, ok := pod.Containers[containerKey]
	if !ok {
		container = &Container{
			Name:      containerRaw.Name,
			Image:     newImage(util.ExtractImageParts(containerRaw.Image.Repo)),
			Processes: map[string]*Process{},
		}
		sinkResult.Inserted(container.GetKey())
		pod.Containers[containerKey] = container
	}

	// Parent
	parentRaw, err := rawEvent.GetParentProcess()
	if err != nil {
		return nil, err
	}
	parentRawKey := parentRaw.GetKey()

	parent, ok := container.Processes[parentRawKey]
	if !ok {
		parent = &Process{
			Binary:         parentRaw.Binary,
			Arguments:      parentRaw.Arguments,
			ChildProcesses: map[string]*Process{},
		}
		sinkResult.Inserted(parent.GetKey())
		container.Processes[parentRawKey] = parent
	}

	// Process
	processRaw, err := rawEvent.GetProcess()
	if err != nil {
		return nil, err
	}

	processRawKey := processRaw.GetKey()

	process, ok := parent.ChildProcesses[processRawKey]
	if !ok {
		process = &Process{
			Binary:         processRaw.Binary,
			Arguments:      processRaw.Arguments,
			ChildProcesses: map[string]*Process{},
		}
		sinkResult.Inserted(process.GetKey())
		parent.ChildProcesses[processRawKey] = process
	}

	util.Noop(process)

	elapsedTime := time.Since(startTime)
	log.Printf("Add function took %s", elapsedTime)

	return nil, nil
}

func (pod *Pod) GetName() string {
	return util.ExtractPodName(pod.Name)
}

func (cluster *Cluster) GetKey() string {
	return key("cluster", cluster.Name)
}

func (namespace *Namespace) GetKey() string {
	return key(namespaceType, namespace.Name)
}

func (pod *Pod) GetKey() string {
	return key(podType, pod.GetName())
}

func (container *Container) GetKey() string {
	return key(containerType, container.Name)
}

func (process *Process) GetKey() string {
	hashValue := process.Binary + ":" + process.Arguments
	name := uuid.NewV5(uuid.NamespaceURL, hashValue).String()
	return key(processType, name)
}

// key generates a unique key for a given type and value.
// It concatenates the type and value with a colon separator.
func key(t string, value string) string {
	return t + ":" + value
}

// newImage creates a new Image object with the given registry, repository, and tag.
// It initializes the Registry object with the given registry name.
func newImage(registry string, repo string, tag string) *Image {
	return &Image{
		Repo: repo,
		Tag:  tag,
		Registry: &Registry{
			Name: registry,
		},
	}
}

// Inserted appends the given path to the SinkResult path and sets the operation to SinkOperationInserted.
func (sinkResult *SinkResult) Inserted(path string) {
	sinkResult.Path = append(sinkResult.Path, path)
	sinkResult.Operation = SinkOperationInserted
}
