/*
Copyright 2019 IBM.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CredSource specifies a credential source either a Secret or a ConfigMap
type CredSource struct {
	// Selects a key of a ConfigMap in the local namespace.
	// +optional
	ConfigMapKeyRef *v1.ConfigMapKeySelector `json:"configMapKeyRef,omitempty" protobuf:"bytes,3,opt,name=configMapKeyRef"`

	// Selects a key of a secret in the local namespace
	// +optional
	SecretKeyRef *v1.SecretKeySelector `json:"secretKeyRef,omitempty" protobuf:"bytes,4,opt,name=secretKeyRef"`
}

// BindingSource specifies a Binding source
type BindingSource struct {
	// The Secret to select from.
	v1.LocalObjectReference `json:",inline" protobuf:"bytes,1,opt,name=localObjectReference"`
}

// EsIndexSpec defines the desired state of EsIndex
type EsIndexSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Name of Index to be created on elastic search
	IndexName string `json:"indexName"`

	// Binding resource name that holds the secret for elastic search credentials
	// +optional
	BindingFrom BindingSource `json:"bindingFrom,omitempty"`

	// EsURIComposed is the URI of elasticesearch resource in the format https://<user>:<passwd>@hostname:port. Cannot be used if BindingFrom is not empty.
	// +optional
	EsURIComposed CredSource `json:"esURIComposed,omitempty"`

	// Bind to an existing index if true, default value false
	// +optional
	BindOnly bool `json:"bindOnly,omitempty"`

	// Number of shards, default value 1
	// +optional
	NumberOfShards int64 `json:"numberOfShards,omitempty"`

	// Number of replicas, default value 1
	// +optional
	NumberOfReplicas int64 `json:"numberOfReplicas,omitempty"`
}

// EsIndexStatus defines the observed state of EsIndex
type EsIndexStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// State of the resource instance
	State string `json:"state,omitempty"`

	// Message
	Message string `json:"message,omitempty"`

	// Generation of the implemented spec
	Generation int64 `json:"generation"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EsIndex is the Schema for the esindices API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type EsIndex struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EsIndexSpec   `json:"spec,omitempty"`
	Status EsIndexStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EsIndexList contains a list of EsIndex
type EsIndexList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EsIndex `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EsIndex{}, &EsIndexList{})
}

// GetStatus returns the function status
func (r *EsIndex) GetStatus() *EsIndexStatus {
	return &r.Status
}
