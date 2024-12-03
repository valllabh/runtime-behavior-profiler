package clusterbehaviourprofiletype

import (
	eventtype "runtime-behavior-profiler/pkg/event/type"
	"runtime-behavior-profiler/pkg/util"

	"github.com/gofrs/uuid"
)

const (
	namespaceType = "namespace"
	podType       = "pod"
	containerType = "container"
	processType   = "process"
)

// Add processes a raw event and updates the ClusterBehaviourProfile accordingly.
// It extracts and organizes the event data into a hierarchical structure of namespaces,
// pods, containers, and processes. If any of these entities do not exist in the profile,
// they are created and added to the appropriate parent entity.
func (cbp *ClusterBehaviourProfile) Add(rawEvent eventtype.Event) error {

	// Namespace
	namespaceRaw, err := rawEvent.GetNamespace()
	if err != nil {
		return err
	}
	namespaceKey := key(namespaceType, namespaceRaw)

	namespace, ok := cbp.Namespaces[namespaceKey]

	if !ok {
		namespace = Namespace{
			Name: namespaceRaw,
			Pods: map[string]Pod{},
		}
		cbp.Namespaces[namespaceKey] = namespace
	}

	// Pod
	podRaw, err := rawEvent.GetPod()
	if err != nil {
		return err
	}

	podName := util.ExtractPodName(podRaw.Name)
	podKey := key(podType, podName)

	pod, ok := namespace.Pods[podKey]
	if !ok {
		pod = Pod{
			Name:       podName,
			Containers: map[string]Container{},
		}
		namespace.Pods[podKey] = pod
	}

	// Container
	containerRaw, err := rawEvent.GetContainer()
	if err != nil {
		return err
	}

	containerKey := key(containerType, containerRaw.Name)

	container, ok := pod.Containers[containerKey]
	if !ok {
		container = Container{
			Name:      containerRaw.Name,
			Image:     newImage(util.ExtractImageParts(containerRaw.Image.Name)),
			Processes: map[string]Process{},
		}
		pod.Containers[containerKey] = container
	}

	// Parent
	parentRaw, err := rawEvent.GetParent()
	if err != nil {
		return err
	}

	parentRawKey := key(processType, getProcessKey(parentRaw))

	parent, ok := container.Processes[parentRawKey]
	if !ok {
		parent = Process{
			Binary:         parentRaw.Binary,
			Arguments:      parentRaw.Arguments,
			ChildProcesses: map[string]Process{},
		}
		container.Processes[parentRawKey] = parent
	}

	// Process
	processRaw, err := rawEvent.GetProcess()
	if err != nil {
		return err
	}

	processRawKey := key(processType, getProcessKey(processRaw))

	process, ok := parent.ChildProcesses[processRawKey]
	if !ok {
		process = Process{
			Binary:         processRaw.Binary,
			Arguments:      processRaw.Arguments,
			ChildProcesses: map[string]Process{},
		}
		parent.ChildProcesses[processRawKey] = process
	}

	util.Noop(process)

	return nil
}

// key generates a unique key for a given type and value.
// It concatenates the type and value with a colon separator.
func key(t string, value string) string {
	return t + ":" + value
}

// newImage creates a new Image object with the given registry, repository, and tag.
// It initializes the Registry object with the given registry name.
func newImage(registry string, repo string, tag string) Image {
	return Image{
		Repo: repo,
		Tag:  tag,
		Registry: Registry{
			Name: registry,
		},
	}
}

// getProcessKey generates a unique key for a given process based on its binary and arguments.
// It concatenates the process binary and arguments with a colon separator and then creates
// a UUID version 5 (namespace-based) using the concatenated string and the URL namespace.
func getProcessKey(process eventtype.Process) string {
	name := process.Binary + ":" + process.Arguments
	return uuid.NewV5(uuid.NamespaceURL, name).String()
}
