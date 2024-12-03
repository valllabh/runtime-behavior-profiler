package eventprocessor

import (
	"encoding/json"
	"os"
	"runtime-behavior-profiler/pkg/clusterbehaviourprofile"
	eventtype "runtime-behavior-profiler/pkg/event/type"
	"testing"
)

func TestProcessEvent(t *testing.T) {

	// Read events from the JSON file
	file, err := os.Open("../../testdata/raw_events.json")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	var events []eventtype.Event
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&events); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if len(events) == 0 {
		t.Fatalf("no events found in the file")
	}

	// Get the cluster behaviour profile
	cbp := clusterbehaviourprofile.GetClusterBehaviourProfile("cluster1")

	// Process each event
	for _, event := range events {
		_, err := ProcessEvent(event, cbp)

		// fail test if error
		if err != nil {
			t.Fatalf("failed to process event: %v", err)
		}
	}

	// Print
	json, err := json.MarshalIndent(cbp, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal cluster behaviour profile: %v", err)
	}

	t.Log("\n" + string(json))
}
