package kubernetes

import (
	"reflect"
	"testing"

	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func TestDeploymentAnnotations(t *testing.T) {
	raw := &v1beta1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "deployment",
			Namespace: "default",
			Annotations: map[string]string{
				"deploy/target":           "1",
				"deploy/target-container": "rails",
				"github":                  "dtan4/rails-app",
			},
			Labels: map[string]string{
				"app":   "rails-app",
				"color": "blue",
			},
		},
		Spec: v1beta1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:  "rails",
							Image: "my-rails:v3",
						},
					},
				},
			},
		},
	}
	deployment := &Deployment{
		raw: raw,
	}

	expected := map[string]string{
		"deploy/target":           "1",
		"deploy/target-container": "rails",
		"github":                  "dtan4/rails-app",
	}
	if got := deployment.Annotations(); !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %q, got: %q", expected, got)
	}
}

func TestDeploymentContainers(t *testing.T) {
	raw := &v1beta1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "deployment",
			Namespace: "default",
			Annotations: map[string]string{
				"deploy/target":           "1",
				"deploy/target-container": "rails",
				"github":                  "dtan4/rails-app",
			},
			Labels: map[string]string{
				"app":   "rails-app",
				"color": "blue",
			},
		},
		Spec: v1beta1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:  "rails",
							Image: "my-rails:v3",
						},
					},
				},
			},
		},
	}
	deployment := &Deployment{
		raw: raw,
	}

	expected := []*Container{
		&Container{
			raw: &v1.Container{
				Name:  "rails",
				Image: "my-rails:v3",
			},
		},
	}
	if got := deployment.Containers(); !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %q, got: %q", expected, got)
	}
}

func TestDeploymentLabels(t *testing.T) {
	raw := &v1beta1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "deployment",
			Namespace: "default",
			Annotations: map[string]string{
				"deploy/target":           "1",
				"deploy/target-container": "rails",
				"github":                  "dtan4/rails-app",
			},
			Labels: map[string]string{
				"app":   "rails-app",
				"color": "blue",
			},
		},
		Spec: v1beta1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:  "rails",
							Image: "my-rails:v3",
						},
					},
				},
			},
		},
	}
	deployment := &Deployment{
		raw: raw,
	}

	expected := map[string]string{
		"app":   "rails-app",
		"color": "blue",
	}
	if got := deployment.Labels(); !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %q, got: %q", expected, got)
	}
}

func TestDeploymentName(t *testing.T) {
	raw := &v1beta1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "deployment",
			Namespace: "default",
			Annotations: map[string]string{
				"deploy/target":           "1",
				"deploy/target-container": "rails",
				"github":                  "dtan4/rails-app",
			},
			Labels: map[string]string{
				"app":   "rails-app",
				"color": "blue",
			},
		},
		Spec: v1beta1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:  "rails",
							Image: "my-rails:v3",
						},
					},
				},
			},
		},
	}
	deployment := &Deployment{
		raw: raw,
	}

	expected := "deployment"
	if got := deployment.Name(); got != expected {
		t.Errorf("expected: %q, got: %q", expected, got)
	}
}

func TestDeploymentNamespace(t *testing.T) {
	raw := &v1beta1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "deployment",
			Namespace: "default",
			Annotations: map[string]string{
				"deploy/target":           "1",
				"deploy/target-container": "rails",
				"github":                  "dtan4/rails-app",
			},
			Labels: map[string]string{
				"app":   "rails-app",
				"color": "blue",
			},
		},
		Spec: v1beta1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:  "rails",
							Image: "my-rails:v3",
						},
					},
				},
			},
		},
	}
	deployment := &Deployment{
		raw: raw,
	}

	expected := "default"
	if got := deployment.Namespace(); got != expected {
		t.Errorf("expected: %q, got: %q", expected, got)
	}
}
