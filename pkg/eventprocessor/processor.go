package eventprocessor

import (
	clusterbehaviourprofiletype "runtime-behavior-profiler/pkg/clusterbehaviourprofile/type"
	eventtype "runtime-behavior-profiler/pkg/event/type"
	eventprocessortype "runtime-behavior-profiler/pkg/eventprocessor/type"
)

// ProcessEvent processes a raw event and updates the ClusterBehaviourProfile accordingly.
// It extracts and organizes the event data into a hierarchical structure of namespaces,
// pods, containers, and processes. If any of these entities do not exist in the profile,
// they are created and added to the appropriate parent entity.
func ProcessEvent(event eventtype.Event, cbp clusterbehaviourprofiletype.ClusterBehaviourProfile) (eventprocessortype.ProcessEventResult, error) {

	// Add the event to the ClusterBehaviourProfile
	err := cbp.Add(event)
	if err != nil {
		return eventprocessortype.ProcessEventResult{}, err
	}

	// Return the result
	return eventprocessortype.ProcessEventResult{}, nil
}
