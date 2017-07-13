package kubernetes

import (
	"testing"

	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func TestContainerImageFromDeployment(t *testing.T) {
	deployment := &v1beta1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "deployment",
			Namespace: "default",
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

	testcases := []struct {
		container string
		expected  string
	}{
		{
			container: "rails",
			expected:  "my-rails:v3",
		},
		{
			container: "foobar",
			expected:  "",
		},
	}

	for _, tc := range testcases {
		if got := ContainerImageFromDeployment(deployment, tc.container); got != tc.expected {
			t.Errorf("expected: %q, got: %q", tc.expected, got)
		}
	}
}
