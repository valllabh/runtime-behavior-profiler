package eventprocessorocsf

import (
	"encoding/json"
	"fmt"
	eventprocessorocsftype "runtime-behavior-profiler/pkg/event/processor/ocsf/type"
	eventtype "runtime-behavior-profiler/pkg/event/type"
)

// ProcessEvent processes a raw event and updates the ClusterBehaviourProfile accordingly.
// It extracts and organizes the event data into a hierarchical structure of namespaces,
// pods, containers, and processes. If any of these entities do not exist in the profile,
// they are created and added to the appropriate parent entity.
func ProcessEvent(event *eventprocessorocsftype.OCSFEvent, cluster *eventtype.Cluster) (*eventtype.SinkResult, error) {

	// Add the event to the ClusterBehaviourProfile
	sinkResult, err := cluster.SinkEvent(event)

	if err != nil {
		return nil, err
	}

	// Pring the sinkResult in JSON
	sinkResultJSON, _ := json.Marshal(sinkResult)

	// Print the sinkResultJSON
	fmt.Println(string(sinkResultJSON))

	return sinkResult, nil
}
