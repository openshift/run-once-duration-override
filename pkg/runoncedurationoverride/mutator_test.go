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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestMutator_Mutate(t *testing.T) {
	config := &Config{
		ActiveDeadlineSeconds: 10,
	}
	mutator, err := NewMutator(config)
	require.NoError(t, err)
	require.NotNil(t, mutator)

	pod := &corev1.Pod{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "name",
					Image: "image",
				},
			},
		},
	}

	podGot, errGot := mutator.Mutate(pod)

	seconds := int64(10)

	assert.NoError(t, errGot)
	assert.NotNil(t, podGot)
	assert.Equal(t, &seconds, podGot.Spec.ActiveDeadlineSeconds)
}
