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

package runoncedurationoverride

import (
	"encoding/json"

	jsonpatch "gomodules.xyz/jsonpatch/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Patch takes 2 byte arrays and returns a new response with json patch.
// The original object should be passed in as raw bytes to avoid the roundtripping problem
// described in https://github.com/kubernetes-sigs/kubebuilder/issues/510.
func Patch(original runtime.RawExtension, mutated *corev1.Pod) (patches []byte, err error) {
	current, marshalErr := json.Marshal(mutated)
	if marshalErr != nil {
		err = marshalErr
		return
	}

	operations, patchErr := jsonpatch.CreatePatch(original.Raw, current)
	if patchErr != nil {
		err = patchErr
		return
	}

	patchBytes, marshalErr := json.Marshal(operations)
	if marshalErr != nil {
		err = marshalErr
		return
	}

	patches = patchBytes
	return
}
