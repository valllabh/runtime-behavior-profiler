package eventprocessor

import (
	"encoding/json"
	"fmt"
	eventtype "runtime-behavior-profiler/pkg/event/type"
	eventprocessortype "runtime-behavior-profiler/pkg/eventprocessor/type"
)

// ProcessEvent processes a raw event and updates the ClusterBehaviourProfile accordingly.
// It extracts and organizes the event data into a hierarchical structure of namespaces,
// pods, containers, and processes. If any of these entities do not exist in the profile,
// they are created and added to the appropriate parent entity.
func ProcessEvent(event *eventtype.OCSFEvent, cluster *eventtype.Cluster) (eventprocessortype.ProcessEventResult, error) {

	// Add the event to the ClusterBehaviourProfile
	sinkResult, err := cluster.SinkEvent(event)

	// Pring the sinkResult in JSON
	sinkResultJSON, err := json.Marshal(sinkResult)

	// Print the sinkResultJSON
	_, err = fmt.Println(string(sinkResultJSON))

	if err != nil {
		return eventprocessortype.ProcessEventResult{}, err
	}

	// Return the result
	return eventprocessortype.ProcessEventResult{}, nil
}
