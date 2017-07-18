package kubernetes

import (
	"reflect"
	"strings"
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

func TestContainerImageFromDeployment(t *testing.T) {
	deployment := &Deployment{
		raw: &v1beta1.Deployment{
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
		if got := deployment.ContainerImage(tc.container); got != tc.expected {
			t.Errorf("expected: %q, got: %q", tc.expected, got)
		}
	}
}

func TestDeployTargetContainer(t *testing.T) {
	testcases := []struct {
		deployment   *Deployment
		expectErr    bool
		expectedName string
		errMsg       string
	}{
		{
			deployment: &Deployment{
				raw: &v1beta1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name:      "deployment",
						Namespace: "default",
						Annotations: map[string]string{
							"deploy/target":           "1",
							"deploy/target-container": "rails",
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
				},
			},
			expectErr:    false,
			expectedName: "rails",
		},
		{
			deployment: &Deployment{
				raw: &v1beta1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name:      "deployment",
						Namespace: "default",
						Annotations: map[string]string{
							"deploy/target":           "1",
							"deploy/target-container": "nginx",
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
				},
			},
			expectErr: true,
			errMsg:    `container "nginx" does not exist in Deployment "deployment"`,
		},
		{
			deployment: &Deployment{
				raw: &v1beta1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name:        "deployment",
						Namespace:   "default",
						Annotations: map[string]string{},
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
				},
			},
			expectErr: true,
			errMsg:    `annotation "deploy/target-container" does not exist in Deployment "deployment"`,
		},
	}

	for _, tc := range testcases {
		got, err := tc.deployment.DeployTargetContainer()

		if tc.expectErr {
			if err == nil {
				t.Error("got no error")
				continue
			}

			if !strings.Contains(err.Error(), tc.errMsg) {
				t.Errorf("error %q does not container %q", err.Error(), tc.errMsg)
			}
		} else {
			if err != nil {
				t.Errorf("got error: %s", err)
			}

			if got.Name() != tc.expectedName {
				t.Errorf("expected name: %q, got: %q", tc.expectedName, got.Name())
			}
		}
	}
}

func TestIsDeployTarget(t *testing.T) {
	testcases := []struct {
		deployment *Deployment
		expected   bool
	}{
		{
			deployment: &Deployment{
				raw: &v1beta1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name:      "deployment",
						Namespace: "default",
						Annotations: map[string]string{
							"deploy/target": "1",
						},
						Labels: map[string]string{
							"app":   "rails-app",
							"color": "blue",
						},
					},
				},
			},
			expected: true,
		},
		{
			deployment: &Deployment{
				raw: &v1beta1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name:      "deployment",
						Namespace: "default",
						Annotations: map[string]string{
							"deploy/target": "true",
						},
						Labels: map[string]string{
							"app":   "rails-app",
							"color": "blue",
						},
					},
				},
			},
			expected: true,
		},
		{
			deployment: &Deployment{
				raw: &v1beta1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name:      "deployment",
						Namespace: "default",
						Annotations: map[string]string{
							"deploy/target": "false",
						},
						Labels: map[string]string{
							"app":   "rails-app",
							"color": "blue",
						},
					},
				},
			},
			expected: false,
		},
		{
			deployment: &Deployment{
				raw: &v1beta1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name:        "deployment",
						Namespace:   "default",
						Annotations: map[string]string{},
						Labels: map[string]string{
							"app":   "rails-app",
							"color": "blue",
						},
					},
				},
			},
			expected: false,
		},
	}

	for _, tc := range testcases {
		if got := tc.deployment.IsDeployTarget(); got != tc.expected {
			t.Errorf("expected: %t, got: %t", tc.expected, got)
		}
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

func TestRepositories(t *testing.T) {
	testcases := []struct {
		deployment *Deployment
		expectErr  bool
		expected   map[string]string
		errMsg     string
	}{
		{
			deployment: &Deployment{
				raw: &v1beta1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name:      "deployment",
						Namespace: "default",
						Annotations: map[string]string{
							"github": "rails=rails/rails",
						},
					},
				},
			},
			expectErr: false,
			expected: map[string]string{
				"rails": "rails/rails",
			},
		},
		{
			deployment: &Deployment{
				raw: &v1beta1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name:      "deployment",
						Namespace: "default",
						Annotations: map[string]string{
							"github": "rails=rails/rails,foobar=dtan4/foobar",
						},
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
			deployment: &Deployment{
				raw: &v1beta1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name:        "deployment",
						Namespace:   "default",
						Annotations: map[string]string{},
					},
				},
			},
			expectErr: true,
			errMsg:    `annotation "github" not found in Deployment "deployment"`,
		},
		{
			deployment: &Deployment{
				raw: &v1beta1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name:      "deployment",
						Namespace: "default",
						Annotations: map[string]string{
							"github": "rails=rails/rails,",
						},
					},
				},
			},
			expectErr: true,
			errMsg:    `invalid annotation "github" value "", must be "container=owner/repo"`,
		},
		{
			deployment: &Deployment{
				raw: &v1beta1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name:      "deployment",
						Namespace: "default",
						Annotations: map[string]string{
							"github": "foobarbaz",
						},
					},
				},
			},
			expectErr: true,
			errMsg:    `invalid annotation "github" value "foobarbaz", must be "container=owner/repo"`,
		},
	}

	for _, tc := range testcases {
		got, err := tc.deployment.Repositories()

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
