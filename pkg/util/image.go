package util

import "strings"

// ExtractImageParts splits an image string into its registry, repository, and tag parts.
// The image string should be in the format: [registry/][repository/]image[:tag]
// If the registry is not included, it will default to an empty string.
// If the tag is not included, it will default to "latest".
func ExtractImageParts(image string) (string, string, string) {
	var registry, repository, tag string

	// Split into main parts: registry/repo and tag
	parts := strings.Split(image, ":")
	if len(parts) == 2 {
		tag = parts[1]
	} else {
		tag = "latest" // Default tag
	}

	imageParts := strings.Split(parts[0], "/")

	// Determine if registry is included
	if len(imageParts) > 2 || strings.Contains(imageParts[0], ".") || strings.Contains(imageParts[0], ":") {
		registry = imageParts[0]
		repository = strings.Join(imageParts[1:], "/")
	} else {
		registry = "" // Default registry
		repository = parts[0]
	}

	return registry, repository, tag
}
