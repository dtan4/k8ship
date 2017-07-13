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
		got, err := client.GetDeployment(namespace, tc.name)

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
			}

			if got.Name != tc.name {
				t.Errorf("expected deployment: %q, got: %q", tc.name, got.Name)
			}
		}
	}
}
