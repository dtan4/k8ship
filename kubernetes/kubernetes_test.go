package kubernetes

import (
	"strings"
	"testing"

	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func TestGetTargetImage(t *testing.T) {
	testcases := []struct {
		containers map[string]*Container
		expectErr  bool
		expected   string
		errMsg     string
	}{
		{
			containers: map[string]*Container{
				"web": &Container{
					raw: &v1.Container{
						Name:  "web",
						Image: "my-rails:v3",
					},
				},
			},
			expectErr: false,
			expected:  "my-rails",
		},
		{
			containers: map[string]*Container{
				"web": &Container{
					raw: &v1.Container{
						Name:  "web",
						Image: "my-rails:v3",
					},
				},
				"worker": &Container{
					raw: &v1.Container{
						Name:  "worker",
						Image: "my-rails:v3",
					},
				},
			},
			expectErr: false,
			expected:  "my-rails",
		},
		{
			containers: map[string]*Container{
				"web": &Container{
					raw: &v1.Container{
						Name:  "web",
						Image: "my-rails:v3",
					},
				},
				"worker": &Container{
					raw: &v1.Container{
						Name:  "worker",
						Image: "my-rails:abc123",
					},
				},
			},
			expectErr: false,
			expected:  "my-rails",
		},
		{
			containers: map[string]*Container{
				"web": &Container{
					raw: &v1.Container{
						Name:  "web",
						Image: "my-rails:v3",
					},
				},
				"nginx": &Container{
					raw: &v1.Container{
						Name:  "nginx",
						Image: "nginx:latest",
					},
				},
			},
			expectErr: true,
			// order of repository list is random
			// because this list is extracted from map[string]bool
			errMsg: `all target containers must use the same image`,
		},
		{
			containers: map[string]*Container{},
			expectErr:  true,
			errMsg:     `no image found`,
		},
	}

	for _, tc := range testcases {
		got, err := GetTargetImage(tc.containers)

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

			if got != tc.expected {
				t.Errorf("expected: %q, got: %q", tc.expected, got)
			}
		}
	}
}

func TestGetTargetRepository(t *testing.T) {
	testcases := []struct {
		deployments []*Deployment
		containers  map[string]*Container
		expectErr   bool
		expected    string
		errMsg      string
	}{
		{
			deployments: []*Deployment{
				&Deployment{
					raw: &v1beta1.Deployment{
						ObjectMeta: v1.ObjectMeta{
							Name:      "web",
							Namespace: "default",
							Annotations: map[string]string{
								"github": "web=dtan4/my-rails",
							},
						},
					},
				},
			},
			containers: map[string]*Container{
				"web": &Container{
					raw: &v1.Container{
						Name:  "web",
						Image: "my-rails:v3",
					},
				},
			},
			expectErr: false,
			expected:  "dtan4/my-rails",
		},
		{
			deployments: []*Deployment{
				&Deployment{
					raw: &v1beta1.Deployment{
						ObjectMeta: v1.ObjectMeta{
							Name:      "web",
							Namespace: "default",
							Annotations: map[string]string{
								"github": "web=dtan4/my-rails",
							},
						},
					},
				},
				&Deployment{
					raw: &v1beta1.Deployment{
						ObjectMeta: v1.ObjectMeta{
							Name:      "worker",
							Namespace: "default",
							Annotations: map[string]string{
								"github": "worker=dtan4/my-rails",
							},
						},
					},
				},
			},
			containers: map[string]*Container{
				"web": &Container{
					raw: &v1.Container{
						Name:  "web",
						Image: "my-rails:v3",
					},
				},
				"worker": &Container{
					raw: &v1.Container{
						Name:  "worker",
						Image: "my-rails:v3",
					},
				},
			},
			expectErr: false,
			expected:  "dtan4/my-rails",
		},
		{
			deployments: []*Deployment{
				&Deployment{
					raw: &v1beta1.Deployment{
						ObjectMeta: v1.ObjectMeta{
							Name:        "web",
							Namespace:   "default",
							Annotations: map[string]string{},
						},
					},
				},
			},
			containers: map[string]*Container{
				"web": &Container{
					raw: &v1.Container{
						Name:  "web",
						Image: "my-rails:v3",
					},
				},
			},
			expectErr: true,
			errMsg:    `failed to retrieve repositories of Deployment "web"`,
		},
		{
			deployments: []*Deployment{
				&Deployment{
					raw: &v1beta1.Deployment{
						ObjectMeta: v1.ObjectMeta{
							Name:      "web",
							Namespace: "default",
							Annotations: map[string]string{
								"github": "foo=dtan4/foo",
							},
						},
					},
				},
			},
			containers: map[string]*Container{
				"web": &Container{
					raw: &v1.Container{
						Name:  "web",
						Image: "my-rails:v3",
					},
				},
			},
			expectErr: true,
			errMsg:    `GitHub repository for container "web" in Deployment "web" is not set`,
		},
		{
			deployments: []*Deployment{
				&Deployment{
					raw: &v1beta1.Deployment{
						ObjectMeta: v1.ObjectMeta{
							Name:      "web",
							Namespace: "default",
							Annotations: map[string]string{
								"github": "web=dtan4/my-rails",
							},
						},
					},
				},
				&Deployment{
					raw: &v1beta1.Deployment{
						ObjectMeta: v1.ObjectMeta{
							Name:      "nginx",
							Namespace: "default",
							Annotations: map[string]string{
								"github": "nginx=dtan4/my-nginx",
							},
						},
					},
				},
			},
			containers: map[string]*Container{
				"web": &Container{
					raw: &v1.Container{
						Name:  "web",
						Image: "my-rails:v3",
					},
				},
				"nginx": &Container{
					raw: &v1.Container{
						Name:  "nginx",
						Image: "nginx:latest",
					},
				},
			},
			expectErr: true,
			// order of repository list is random
			// because this list is extracted from map[string]bool
			errMsg: `all target containers must use the same repository`,
		},
		{
			deployments: []*Deployment{
				&Deployment{
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
			},
			containers: map[string]*Container{},
			expectErr:  true,
			errMsg:     `no container found in Deployment "deployment"`,
		},
	}

	for _, tc := range testcases {
		got, err := GetTargetRepository(tc.deployments, tc.containers)

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

			if got != tc.expected {
				t.Errorf("expected: %q, got: %q", tc.expected, got)
			}
		}
	}
}
