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
