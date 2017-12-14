package kubernetes

import (
	"time"

	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

const (
	revisionAnnotation = "deployment.kubernetes.io/revision"
)

// ReplicaSet represents the wrapper of Kubernetes ReplicaSet
type ReplicaSet struct {
	raw *v1beta1.ReplicaSet
}

// NewReplicaSet creates ne ReplicaSet object
func NewReplicaSet(raw *v1beta1.ReplicaSet) *ReplicaSet {
	return &ReplicaSet{
		raw: raw,
	}
}

// CreatedAt returns the creation timestamp
func (r *ReplicaSet) CreatedAt() time.Time {
	return r.raw.CreationTimestamp.Time
}

// Name returns the name of ReplicaSet
func (r *ReplicaSet) Name() string {
	return r.raw.Name
}

// Namespace returns the namespace of ReplicaSet
func (r *ReplicaSet) Namespace() string {
	return r.raw.Namespace
}

// Revision returns the revision signature
func (r *ReplicaSet) Revision() string {
	return r.raw.Annotations[revisionAnnotation]
}
