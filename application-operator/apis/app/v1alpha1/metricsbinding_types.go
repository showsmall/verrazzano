// Copyright (c) 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&MetricsBinding{}, &MetricsBindingList{})
}

// MetricsBindingList contains a list of metrics binding resources
// +kubebuilder:object:root=true
type MetricsBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MetricsBinding `json:"items"`
}

// MetricsBinding specifies the metrics binding API
// +kubebuilder:object:root=true
// +genclient
type MetricsBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec MetricsBindingSpec `json:"spec"`
}

// MetricsBindingSpec specifies the desired state of a metrics binding
type MetricsBindingSpec struct {
	// Identifies a namespace and name for a metricsTemplate resource
	MetricsTemplate NamespaceName `json:"metricsTemplate"`

	// Identifies a namespace and name for a Prometheus configMap resource
	PrometheusConfigMap NamespaceName `json:"prometheusConfigMap,omitempty"`

	// Identifies a namespace, name and key for a secret containing the Prometheus config
	PrometheusConfigSecret SecretKey `json:"prometheusConfigSecret,omitempty"`

	// Identifies the name and type for a workload
	Workload Workload `json:"workload"`
}

// NamespaceName identifies a namespace and name pair for a resource
type NamespaceName struct {
	// Namespace of a resource
	Namespace string `json:"namespace"`

	// Name of a resource
	Name string `json:"name"`
}

// SecretKey identifies a value in a Kubernetes Secret by its namespace, name and key in the Secret.
type SecretKey struct {
	// Namespace of the Secret
	Namespace string `json:"namespace"`

	// Name of the Secret
	Name string `json:"name"`

	// Key in the Secret whose value this object represents
	Key string `json:"key"`
}

// Workload identifies the name and type of a workload
type Workload struct {
	// Name of a resource
	Name string `json:"name"`

	// TypeMeta of a resource
	TypeMeta metav1.TypeMeta `json:"typeMeta"`
}
