/*
Copyright 2020 The Flux CD contributors.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ImageRepositoryKind = "ImageRepository"

// ImageRepositorySpec defines the parameters for scanning an image
// repository, e.g., `fluxcd/flux`.
type ImageRepositorySpec struct {
	// Image is the name of the image repository
	// +required
	Image string `json:"image,omitempty"`
	// ScanInterval is the (minimum) length of time to wait between
	// scans of the image repository.
	// +optional
	ScanInterval *metav1.Duration `json:"scanInterval,omitempty"`

	// SecretRef can be given the name of a secret containing
	// credentials to use for the image registry. The secret should be
	// created with `kubectl create secret docker-registry`, or the
	// equivalent.
	SecretRef *corev1.LocalObjectReference `json:"secretRef,omitempty"`

	// This flag tells the controller to suspend subsequent image scans.
	// It does not apply to already started scans. Defaults to false.
	// +optional
	Suspend bool `json:"suspend,omitempty"`
}

type ScanResult struct {
	TagCount int `json:"tagCount"`
}

// ImageRepositoryStatus defines the observed state of ImageRepository
type ImageRepositoryStatus struct {
	// +optional
	Conditions []Condition `json:"conditions,omitempty"`

	// ObservedGeneration is the last reconciled generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// CannonicalName is the name of the image repository with all the
	// implied bits made explicit; e.g., `docker.io/library/alpine`
	// rather than `alpine`.
	// +optional
	CanonicalImageName string `json:"canonicalImageName,omitempty"`

	// LastScanResult contains the number of fetched tags.
	// +optional
	LastScanResult ScanResult `json:"lastScanResult,omitempty"`
}

// SetImageRepositoryReadiness sets the ready condition with the given status, reason and message.
func SetImageRepositoryReadiness(ir ImageRepository, status corev1.ConditionStatus, reason, message string) ImageRepository {
	ir.Status.Conditions = []Condition{
		{
			Type:               ReadyCondition,
			Status:             status,
			LastTransitionTime: metav1.Now(),
			Reason:             reason,
			Message:            message,
		},
	}
	ir.Status.ObservedGeneration = ir.ObjectMeta.Generation
	return ir
}

func GetLastTransitionTime(ir ImageRepository) *metav1.Time {
	for _, condition := range ir.Status.Conditions {
		if condition.Type == ReadyCondition {
			return &condition.LastTransitionTime
		}
	}

	return nil
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Last scan",type=string,JSONPath=`.status.lastScanTime`
// +kubebuilder:printcolumn:name="Tags",type=string,JSONPath=`.status.lastScanResult.tagCount`

// ImageRepository is the Schema for the imagerepositories API
type ImageRepository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageRepositorySpec   `json:"spec,omitempty"`
	Status ImageRepositoryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ImageRepositoryList contains a list of ImageRepository
type ImageRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ImageRepository `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ImageRepository{}, &ImageRepositoryList{})
}
