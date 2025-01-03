package util

import (
	"testing"
)

func TestExtractImageParts(t *testing.T) {
	tests := []struct {
		image      string
		registry   string
		repository string
		tag        string
	}{
		{"repository/repository/image:tag", "", "repository/repository/image", "tag"},
		{"repository/image:tag", "", "repository/image", "tag"},
		{"image:tag", "", "image", "tag"},
		{"docker.io/my-repo/test-image:latest", "docker.io", "my-repo/test-image", "latest"},
		{"docker.io/my-repo/sample-service:1.0", "docker.io", "my-repo/sample-service", "1.0"},
		{"docker.io/my-repo/ci-cd-test:dev", "docker.io", "my-repo/ci-cd-test", "dev"},
		{"docker.io/my-repo/debug-container:alpha", "docker.io", "my-repo/debug-container", "alpha"},
		{"my-private-registry.com/test-repo/test-image:1.0.0", "my-private-registry.com", "test-repo/test-image", "1.0.0"},
		{"my-private-registry.com/test-repo/mock-service:qa", "my-private-registry.com", "test-repo/mock-service", "qa"},
		{"my-private-registry.com/test-repo/sample-app:staging", "my-private-registry.com", "test-repo/sample-app", "staging"},
		{"123456789012.dkr.ecr.us-east-1.amazonaws.com/test-repo:test-latest", "123456789012.dkr.ecr.us-east-1.amazonaws.com", "test-repo", "test-latest"},
		{"123456789012.dkr.ecr.us-east-1.amazonaws.com/sample-repo:integration-v1", "123456789012.dkr.ecr.us-east-1.amazonaws.com", "sample-repo", "integration-v1"},
		{"gcr.io/my-project/test-repo/test-image:v1", "gcr.io", "my-project/test-repo/test-image", "v1"},
		{"gcr.io/my-project/sample-service:test-dev", "gcr.io", "my-project/sample-service", "test-dev"},
		{"myregistry.azurecr.io/test-repo/sample-app:staging", "myregistry.azurecr.io", "test-repo/sample-app", "staging"},
		{"myregistry.azurecr.io/test-repo/test-image:latest", "myregistry.azurecr.io", "test-repo/test-image", "latest"},
		{"myregistry.azurecr.io/test-repo/test-image", "myregistry.azurecr.io", "test-repo/test-image", "latest"},
		{"qregistry:8080/ikhanqualys/performance", "qregistry:8080", "ikhanqualys/performance", "latest"},
		{"qregistry:8080/ikhanqualys/performance:celery", "qregistry:8080", "ikhanqualys/performance", "celery"},
		{"qregistry:8080/ikhanqualys/performance@sha256:4c1c50d0ffc614f90b93b07d778028dc765548e823f676fb027f61d281ac380d", "qregistry:8080", "ikhanqualys/performance", "sha256:4c1c50d0ffc614f90b93b07d778028dc765548e823f676fb027f61d281ac380d"},
		{"docker.io/ikhanqualys/performance:celery", "docker.io", "ikhanqualys/performance", "celery"},
		{"ikhanqualys/performance:celery", "", "ikhanqualys/performance", "celery"},
		{"qregistry:8080/ikhanqualys/a/b/c/d/e/f/performance:celery", "qregistry:8080", "ikhanqualys/a/b/c/d/e/f/performance", "celery"},
		{"art-hq.intranet.qualys.com:5006/secure/oraclelinux:8-slim", "art-hq.intranet.qualys.com:5006", "secure/oraclelinux", "8-slim"},
		{"art-hq.intranet.qualys.com:5001/cs/build/golang-cgo:oel8", "art-hq.intranet.qualys.com:5001", "cs/build/golang-cgo", "oel8"},
	}

	for _, test := range tests {
		registry, repository, tag := ExtractImageParts(test.image)
		if registry != test.registry || repository != test.repository || tag != test.tag {
			t.Errorf("ExtractImageParts(%q) = %q, %q, %q; want %q, %q, %q",
				test.image, registry, repository, tag, test.registry, test.repository, test.tag)
		}
	}
}
