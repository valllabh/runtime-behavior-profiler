package main

import (
	eventprocessortetragon "runtime-behavior-profiler/pkg/event/processor/tetragon"
)

func main() {
	err := eventprocessortetragon.ListenToEvents()
	if err != nil {
		println(err.Error())
	}
}
