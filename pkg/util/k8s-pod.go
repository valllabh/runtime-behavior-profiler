package util

import "regexp"

// ExtractPodName extracts the base pod name from a given Kubernetes pod name string.
// The base pod name is the first segment of the pod name before the hash/identifier.
// The base pod name is used to group pods by their type or purpose.
// For example, the base pod name of "my-nginx-554b9c67f9-c5cv4" is "my-nginx".
func ExtractPodName(input string) string {
	// Regular expression to match the pattern:
	// - Last segment after hyphen: 5 alphanumeric characters
	// - Second last segment (optional): 10 alphanumeric characters
	re := regexp.MustCompile(`^([a-zA-Z0-9\-]+)-[a-zA-Z0-9]{10}-[a-zA-Z0-9]{5}$|^([a-zA-Z0-9\-]+)-[a-zA-Z0-9]{5}$|^([a-zA-Z0-9\-]+)$`)
	match := re.FindStringSubmatch(input)

	// Return the first matched group with a valid base name
	if match[1] != "" {
		return match[1] // Matches with 10 and 5 alphanumeric segments
	}
	if match[2] != "" {
		return match[2] // Matches with only 5 alphanumeric segment
	}
	if match[3] != "" {
		return match[3] // Matches with no trailing hash/identifier
	}
	return ""
}
