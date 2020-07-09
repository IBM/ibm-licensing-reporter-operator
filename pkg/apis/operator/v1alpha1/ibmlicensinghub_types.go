//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IBMLicensingHubSpec defines the desired state of IBMLicensingHub
type IBMLicensingHubSpec struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Secret name used to store application token, either one that exists, or one that will be created
	APISecretToken string `json:"apiSecretToken,omitempty"`

	// Storage class used by database to provide persistency
	StorageClass string `json:"storageClass,omitempty"`
}

// IBMLicensingHubStatus defines the observed state of IBMLicensingHub
type IBMLicensingHubStatus struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	LicensingHubPods []corev1.PodStatus `json:"LicensingHubPods"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IBMLicensingHub is the Schema for the ibmlicensinghubs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=ibmlicensinghubs,scope=Namespaced
type IBMLicensingHub struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IBMLicensingHubSpec   `json:"spec,omitempty"`
	Status IBMLicensingHubStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IBMLicensingHubList contains a list of IBMLicensingHub
type IBMLicensingHubList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IBMLicensingHub `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IBMLicensingHub{}, &IBMLicensingHubList{})
}
