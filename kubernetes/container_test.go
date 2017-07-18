package kubernetes

import (
	"testing"

	"k8s.io/client-go/pkg/api/v1"
)

func TestContainerImage(t *testing.T) {
	raw := &v1.Container{
		Name:  "rails",
		Image: "my-rails:v3",
	}
	container := &Container{
		raw: raw,
	}

	expected := "my-rails:v3"
	if got := container.Image(); got != expected {
		t.Errorf("expected: %q, got: %q", expected, got)
	}
}

func TestContainerName(t *testing.T) {
	raw := &v1.Container{
		Name:  "rails",
		Image: "my-rails:v3",
	}
	container := &Container{
		raw: raw,
	}

	expected := "rails"
	if got := container.Name(); got != expected {
		t.Errorf("expected: %q, got: %q", expected, got)
	}
}
