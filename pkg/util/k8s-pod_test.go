package util

import (
	"testing"
)

func TestExtractPodName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"my-nginx-554b9c67f9-c5cv4", "my-nginx"},
		{"nginx-554b9c67f9-c5cv4", "nginx"},
		{"nginx-c5cv4", "nginx"},
		{"nginx", "nginx"},
		{"app-server-abc1234567-xy123", "app-server"},
	}

	for _, test := range tests {
		result := ExtractPodName(test.input)
		if result != test.expected {
			t.Errorf("ExtractPodName(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}
