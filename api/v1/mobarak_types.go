/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ContainerSpec defines the desired state of Container
type ContainerSpec struct {
	Image string `json:"image"`
	Port  int32  `json:"port"`
}

// ServiceSpec defines the desired state of Service
type ServiceSpec struct {
	// +optional
	ServiceName string `json:"serviceName,omitempty"`
	ServiceType string `json:"serviceType"`
	ServicePort int32  `json:"servicePort"`
	// +optional
	ServiceNodePort int32 `json:"serviceNodePort,omitempty"`
}

// BookServerSpec defines the desired state of BookServer
type BookServerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// DeploymentName represents the name of the deployment we will create using CustomCrd
	// +optional
	DeploymentName string `json:"deploymentName,omitempty"`
	// Replicas defines number of pods will be running in the deployment
	Replicas *int32 `json:"replicas"`
	// Container contains Image and Port
	Container ContainerSpec `json:"container"`
	// Service contains ServiceName, ServiceType, ServiceNodePort
	// +optional
	Service ServiceSpec `json:"service,omitempty"`
}

// BookServerStatus defines the observed state of BookServer
type BookServerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +optional
	AvailableReplicas *int32 `json:"availableReplicas,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BookServer is the Schema for the mobaraks API
type BookServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BookServerSpec   `json:"spec,omitempty"`
	Status BookServerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BookServerList contains a list of BookServer
type BookServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BookServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BookServer{}, &BookServerList{})
}
