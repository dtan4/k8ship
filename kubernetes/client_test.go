package kubernetes

import (
	"strings"
	"testing"

	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/pkg/types"
)

func TestDetectTargetContainer_with_name(t *testing.T) {
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
								Name: "rails",
							},
						},
					},
				},
			},
		},
	}

	client := &Client{}

	testcases := []struct {
		name      string
		expectErr bool
		errMsg    string
	}{
		{
			name:      "rails",
			expectErr: false,
		},
		{
			name:      "foobar",
			expectErr: true,
			errMsg:    `container "foobar" does not exist`,
		},
	}

	for _, tc := range testcases {
		got, err := client.DetectTargetContainer(deployment, tc.name)

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
				continue
			}

			if got.Name() != tc.name {
				t.Errorf("expected deployment: %q, got: %q", tc.name, got.Name())
			}
		}
	}
}

func TestDetectTargetContainer_without_name(t *testing.T) {
	testcases := []struct {
		deployment *Deployment
		expectErr  bool
		errMsg     string
	}{
		{
			deployment: &Deployment{
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
										Name: "rails",
									},
								},
							},
						},
					},
				},
			},
			expectErr: false,
		},
		{
			deployment: &Deployment{
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
										Name: "rails",
									},
									v1.Container{
										Name: "foobar",
									},
								},
							},
						},
					},
				},
			},
			expectErr: true,
			errMsg:    `multiple containers ["rails" "foobar"] found in deployment "deployment"`,
		},
	}

	name := ""

	for _, tc := range testcases {
		client := &Client{}

		_, err := client.DetectTargetContainer(tc.deployment, name)

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
				continue
			}

			if !strings.Contains(err.Error(), tc.errMsg) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.errMsg)
			}
		} else {
			if err != nil {
				t.Errorf("got error: %s", err)
				continue
			}

			if got.Name() != tc.name {
				t.Errorf("expected deployment: %q, got: %q", tc.name, got.Name())
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

func TestGetDeployment(t *testing.T) {
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
	}{
		{
			name:      "deployment",
			expectErr: false,
		},
		{
			name:      "foobar",
			expectErr: true,
		},
	}

	namespace := "default"

	for _, tc := range testcases {
		got, err := client.GetDeployment(namespace, tc.name)

		if tc.expectErr {
			if err == nil {
				t.Error("got no error")
			}
		} else {
			if err != nil {
				t.Errorf("got error: %s", err)
				continue
			}

			if got.Name() != tc.name {
				t.Errorf("expected deployment: %q, got: %q", tc.name, got.Name())
			}
		}
	}
}

func TestListDeployments(t *testing.T) {
	deployments := []v1beta1.Deployment{
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
	}

	clientset := fake.NewSimpleClientset(&v1beta1.DeploymentList{
		Items: deployments,
	})
	client := &Client{
		clientset: clientset,
	}

	namespace := "default"

	got, err := client.ListDeployments(namespace)
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	expectedLength := 2
	if len(got) != expectedLength {
		t.Errorf("expected length: %d, got: %d", expectedLength, len(got))
	}
}

func TestListReplicaSets(t *testing.T) {
	replicasets := []v1beta1.ReplicaSet{
		v1beta1.ReplicaSet{
			ObjectMeta: v1.ObjectMeta{
				Name:      "deployment-1234567890",
				Namespace: "default",
				OwnerReferences: []v1.OwnerReference{
					v1.OwnerReference{
						UID: (types.UID)("0001"),
					},
				},
			},
		},
		v1beta1.ReplicaSet{
			ObjectMeta: v1.ObjectMeta{
				Name:      "foobar-9876543210",
				Namespace: "default",
				OwnerReferences: []v1.OwnerReference{
					v1.OwnerReference{
						UID: (types.UID)("0002"),
					},
				},
			},
		},
	}

	clientset := fake.NewSimpleClientset(&v1beta1.ReplicaSetList{
		Items: replicasets,
	})
	client := &Client{
		clientset: clientset,
	}

	deployment := &Deployment{
		raw: &v1beta1.Deployment{
			ObjectMeta: v1.ObjectMeta{
				Name:      "deployment",
				Namespace: "default",
				UID:       (types.UID)("0001"),
			},
			Spec: v1beta1.DeploymentSpec{
				Template: v1.PodTemplateSpec{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							v1.Container{
								Name: "rails",
							},
						},
					},
				},
			},
		},
	}

	got, err := client.ListReplicaSets(deployment)
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	expectedLength := 1
	if len(got) != expectedLength {
		t.Errorf("expected length: %d, got: %d", expectedLength, len(got))
	}

	expectedName := "deployment-1234567890"
	if got[0] != expectedName {
		t.Errorf("expected: %q, got: %q", expectedName, got[0])
	}
}

func TestReloadPods(t *testing.T) {
	raw := &v1beta1.Deployment{
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
							Image: "my-rails:v2",
						},
					},
				},
			},
		},
	}
	deployment := &Deployment{
		raw: raw,
	}

	clientset := fake.NewSimpleClientset(raw)
	client := &Client{
		clientset: clientset,
	}

	signature := "2017-12-05 12:18:31.789275051 +0900 JST"

	_, err := client.ReloadPods(deployment, signature)
	if err != nil {
		t.Errorf("got error: %s", err)
		return
	}

	// Unfortunally, there is no way to check the updated Deployment image...
}

func TestSetImage(t *testing.T) {
	raw := &v1beta1.Deployment{
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
							Image: "my-rails:v2",
						},
					},
				},
			},
		},
	}
	deployment := &Deployment{
		raw: raw,
	}

	clientset := fake.NewSimpleClientset(raw)
	client := &Client{
		clientset: clientset,
	}

	container := "rails"
	image := "my-rails:v3"
	cause := "k8ship test"

	_, err := client.SetImage(deployment, container, image, cause)
	if err != nil {
		t.Errorf("got error: %s", err)
		return
	}

	// Unfortunally, there is no way to check the updated Deployment image...
}
