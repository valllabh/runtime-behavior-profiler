package main

import (
	eventprocessortetragon "runtime-behavior-profiler/pkg/event/processor/tetragon"
	eventtype "runtime-behavior-profiler/pkg/event/type"
)

func main() {

	cluster := eventtype.Cluster{
		Name:       "test-cluster",
		Namespaces: map[string]*eventtype.Namespace{},
	}

	eventLister := eventprocessortetragon.NewEventListener(&cluster)

	err := eventLister.ListenToEvents()
	if err != nil {
		println(err.Error())
	}
}
