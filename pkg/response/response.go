/*
Copyright 2023 The run-once-duration-operator Authors.

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

package response

import (
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func WithInternalServerError(request *admissionv1.AdmissionRequest, err error) *admissionv1.AdmissionResponse {
	return &admissionv1.AdmissionResponse{
		UID:     request.UID,
		Allowed: false,
		Result: &metav1.Status{
			Status:  metav1.StatusFailure,
			Code:    http.StatusInternalServerError,
			Reason:  metav1.StatusReasonInternalError,
			Message: err.Error(),
		},
	}
}

func WithAllowed(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	return &admissionv1.AdmissionResponse{
		UID:     request.UID,
		Allowed: true,
	}
}

func WithBadRequest(request *admissionv1.AdmissionRequest, err error) *admissionv1.AdmissionResponse {
	return &admissionv1.AdmissionResponse{
		UID:     request.UID,
		Allowed: false,
		Result: &metav1.Status{
			Status:  metav1.StatusFailure,
			Code:    http.StatusBadRequest,
			Reason:  metav1.StatusReasonBadRequest,
			Message: err.Error(),
		},
	}
}

func WithPatch(request *admissionv1.AdmissionRequest, patch []byte) *admissionv1.AdmissionResponse {
	response := &admissionv1.AdmissionResponse{
		UID:     request.UID,
		Allowed: true,
	}

	if patch != nil {
		patchType := admissionv1.PatchTypeJSONPatch

		response.Patch = patch
		response.PatchType = &patchType
	}

	return response
}
