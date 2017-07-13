package kubernetes

import (
	"strings"
	"testing"

	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func TestDetectTargetDeployment_with_name(t *testing.T) {
	deployment := &v1beta1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "deployment",
			Namespace: "default",
		},
	}

	clientset := fake.NewSimpleClientset(deployment)
	client := &Client{
		clientset: clientset,
	}

	testcases := []struct {
		name      string
		expectErr bool
		errMsg    string
	}{
		{
			name:      "deployment",
			expectErr: false,
		},
		{
			name:      "foobar",
			expectErr: true,
			errMsg:    `failed to retrieve Deployment "foobar"`,
		},
	}

	namespace := "default"

	for _, tc := range testcases {
		got, err := client.DetectTargetDeployment(namespace, tc.name)

		if tc.expectErr {
			if err == nil {
				t.Error("got no error")
			}

			if !strings.Contains(err.Error(), tc.errMsg) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.errMsg)
			}
		} else {
			if err != nil {
				t.Errorf("got error: %s", err)
				continue
			}

			if got.Name != tc.name {
				t.Errorf("expected deployment: %q, got: %q", tc.name, got.Name)
			}
		}
	}
}

func TestDetectTargetDeployment_without_name(t *testing.T) {
	testcases := []struct {
		deployments []v1beta1.Deployment
		expectErr   bool
		errMsg      string
	}{
		{
			deployments: []v1beta1.Deployment{
				v1beta1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name:      "deployment",
						Namespace: "default",
					},
				},
			},
			expectErr: false,
		},
		{
			deployments: []v1beta1.Deployment{},
			expectErr:   true,
			errMsg:      `no Deployment found in namespace "default"`,
		},
		{
			deployments: []v1beta1.Deployment{
				v1beta1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name:      "deployment",
						Namespace: "default",
					},
				},
				v1beta1.Deployment{
					ObjectMeta: v1.ObjectMeta{
						Name:      "foobar",
						Namespace: "default",
					},
				},
			},
			expectErr: true,
			errMsg:    `multiple Deployments ["deployment" "foobar"] found in namespace "default"`,
		},
	}

	name := ""
	namespace := "default"

	for _, tc := range testcases {
		clientset := fake.NewSimpleClientset(&v1beta1.DeploymentList{
			Items: tc.deployments,
		})
		client := &Client{
			clientset: clientset,
		}

		_, err := client.DetectTargetDeployment(namespace, name)

		if tc.expectErr {
			if err == nil {
				t.Error("got no error")
				continue
			}

			if !strings.Contains(err.Error(), tc.errMsg) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.errMsg)
			}
		} else {
			if err != nil {
				t.Errorf("got error: %s", err)
			}
		}
	}
}
