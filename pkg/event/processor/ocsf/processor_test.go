package eventprocessorocsf

import (
	"encoding/json"
	"os"
	eventprocessorocsftype "runtime-behavior-profiler/pkg/event/processor/ocsf/type"
	eventtype "runtime-behavior-profiler/pkg/event/type"
	"testing"
)

func TestProcessEvent(t *testing.T) {

	path := "../../../../testdata/raw_events.json"

	// Read events from the JSON file
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	var events []eventprocessorocsftype.OCSFEvent
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if err := json.Unmarshal(data, &events); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(events) == 0 {
		t.Fatalf("no events found in the file")
	}

	cluster := eventtype.Cluster{
		Name:       "test-cluster",
		Namespaces: map[string]*eventtype.Namespace{},
	}

	// Process each event
	for _, event := range events {

		_, err := ProcessEvent(&event, &cluster)

		// fail test if error
		if err != nil {
			t.Fatalf("failed to process event: %v", err)
		}
	}

	// Print
	json, err := json.MarshalIndent(cluster, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal cluster behaviour profile: %v", err)
	}

	t.Log("\n" + string(json))
}
