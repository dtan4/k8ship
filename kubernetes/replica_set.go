package kubernetes

import (
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
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

// Name returns the name of ReplicaSet
func (r *ReplicaSet) Name() string {
	return r.raw.Name
}

// Namespace returns the namespace of ReplicaSet
func (r *ReplicaSet) Namespace() string {
	return r.raw.Namespace
}
