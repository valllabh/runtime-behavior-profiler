package main

import (
	"fmt"
	eventprocessortetragon "runtime-behavior-profiler/pkg/event/processor/tetragon"
)

var version string

func main() {
	fmt.Printf("Version: %s\n", version)
	err := eventprocessortetragon.ListenToEvents()
	if err != nil {
		println(err.Error())
	}
}
