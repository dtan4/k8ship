package kubernetes

import (
	"reflect"
	"strings"
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

func TestRepositoriesFromDeployment(t *testing.T) {
	testcases := []struct {
		deployment *v1beta1.Deployment
		expectErr  bool
		expected   map[string]string
		errMsg     string
	}{
		{
			deployment: &v1beta1.Deployment{
				ObjectMeta: v1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
					Labels: map[string]string{
						"github": "rails:rails/rails",
					},
				},
			},
			expectErr: false,
			expected: map[string]string{
				"rails": "rails/rails",
			},
		},
		{
			deployment: &v1beta1.Deployment{
				ObjectMeta: v1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
					Labels: map[string]string{
						"github": "rails:rails/rails,foobar:dtan4/foobar",
					},
				},
			},
			expectErr: false,
			expected: map[string]string{
				"rails":  "rails/rails",
				"foobar": "dtan4/foobar",
			},
		},
		{
			deployment: &v1beta1.Deployment{
				ObjectMeta: v1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
					Labels:    map[string]string{},
				},
			},
			expectErr: true,
			errMsg:    `label "github" not found in deployment "deployment"`,
		},
		{
			deployment: &v1beta1.Deployment{
				ObjectMeta: v1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
					Labels: map[string]string{
						"github": "rails:rails/rails,",
					},
				},
			},
			expectErr: true,
			errMsg:    `invalid label "github" value "", must be "container=owner/repo"`,
		},
		{
			deployment: &v1beta1.Deployment{
				ObjectMeta: v1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
					Labels: map[string]string{
						"github": "foobarbaz",
					},
				},
			},
			expectErr: true,
			errMsg:    `invalid label "github" value "foobarbaz", must be "container=owner/repo"`,
		},
	}

	for _, tc := range testcases {
		got, err := RepositoriesFromDeployment(tc.deployment)

		if tc.expectErr {
			if err == nil {
				t.Errorf("got no error")
				continue
			}

			if !strings.Contains(err.Error(), tc.errMsg) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.errMsg)
			}
		} else {
			if err != nil {
				t.Errorf("got error: %s", err)
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("expected: %q, got: %q", tc.expected, got)
			}
		}
	}
}
