/*
Copyright 2021 hollicube.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DemoHollicubeSpec defines the desired state of DemoHollicube
type DemoHollicubeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of DemoHollicube. Edit DemoHollicube_types.go to remove/update
	//Foo string `json:"foo,omitempty"`
	Image         string `json:"image"`
	Replicas      *int32 `json:"replicas"`
	ContainerPort int32  `json:"containerPort"`
	Protocol      string `json:"protocol"`
	CPURequest    int32  `json:"cpuRequest"`
	MEMRequest    int32  `json:"memRequest"`
	CPULimit       int32 `json:"cpuLimit"`
	MEMLimit      int32  `json:"memLimit"`
}

// DemoHollicubeStatus defines the observed state of DemoHollicube
type DemoHollicubeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// DemoHollicube is the Schema for the demohollicubes API
type DemoHollicube struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DemoHollicubeSpec   `json:"spec,omitempty"`
	Status DemoHollicubeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DemoHollicubeList contains a list of DemoHollicube
type DemoHollicubeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DemoHollicube `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DemoHollicube{}, &DemoHollicubeList{})
}
