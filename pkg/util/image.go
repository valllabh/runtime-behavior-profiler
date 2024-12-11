package util

import "strings"

// ExtractImageParts splits an image string into its registry, repository, and tag parts.
// The image string should be in the format: [registry/][repository/]image[:tag]
// If the registry is not included, it will default to an empty string.
// If the tag is not included, it will default to "latest".
func ExtractImageParts(image string) (string, string, string) {
	var registry, repository, tag string

	// Split image by `@` for digests
	if strings.Contains(image, "@") {
		parts := strings.SplitN(image, "@", 2)
		image = parts[0]
		tag = parts[1]
	} else if strings.Contains(image, ":") {
		// Split for tags (considering port numbers and image tags)
		lastColonIndex := strings.LastIndex(image, ":")
		if lastColonIndex > strings.Index(image, "/") {
			// Colon belongs to the tag
			tag = image[lastColonIndex+1:]
			image = image[:lastColonIndex]
		} else {
			// No tag provided; use default
			tag = "latest"
		}
	} else {
		tag = "latest" // Default tag if none provided
	}

	// Split into registry and repository
	imageParts := strings.SplitN(image, "/", 2)
	if len(imageParts) > 1 && (strings.Contains(imageParts[0], ".") || strings.Contains(imageParts[0], ":")) {
		registry = imageParts[0]
		repository = imageParts[1]
	} else {
		registry = "" // Default registry is implied
		repository = image
	}

	return registry, repository, tag
}
