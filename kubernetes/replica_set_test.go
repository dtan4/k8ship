package kubernetes

import (
	"reflect"
	"testing"
	"time"

	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func TestCreatedAt(t *testing.T) {
	raw := &v1beta1.ReplicaSet{
		ObjectMeta: v1.ObjectMeta{
			Name:      "deployment-1234567890",
			Namespace: "default",
			CreationTimestamp: unversioned.Time{
				Time: time.Date(2017, 12, 14, 16, 36, 17, 0, time.UTC),
			},
		},
	}
	r := &ReplicaSet{
		raw: raw,
	}

	got := r.CreatedAt()
	want := time.Date(2017, 12, 14, 16, 36, 17, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestImages(t *testing.T) {
	raw := &v1beta1.ReplicaSet{
		ObjectMeta: v1.ObjectMeta{
			Name:      "deployment-1234567890",
			Namespace: "default",
		},
		Spec: v1beta1.ReplicaSetSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:  "web",
							Image: "web:v1",
						},
						v1.Container{
							Name:  "nginx",
							Image: "nginx:1.13.7",
						},
					},
				},
			},
		},
	}
	r := &ReplicaSet{
		raw: raw,
	}

	got := r.Images()
	want := map[string]string{
		"web":   "web:v1",
		"nginx": "nginx:1.13.7",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %#v, got: %#v", want, got)
	}
}

func TestRevision(t *testing.T) {
	raw := &v1beta1.ReplicaSet{
		ObjectMeta: v1.ObjectMeta{
			Annotations: map[string]string{
				"deployment.kubernetes.io/revision": "1",
			},
			Name:      "deployment-1234567890",
			Namespace: "default",
		},
	}
	r := &ReplicaSet{
		raw: raw,
	}

	got := r.Revision()
	want := "1"
	if got != want {
		t.Errorf("want: %q, got: %q", want, got)
	}
}
