package runoncedurationoverride

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestIsNamespaceExempt(t *testing.T) {
	tests := []struct {
		name   string
		pod    *corev1.Pod
		exempt bool
	}{
		{
			name:   "Empty pod spec",
			pod:    &corev1.Pod{},
			exempt: true,
		},
		{
			name: "RestartPolicy Always",
			pod: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
				},
			},
			exempt: true,
		},
		{
			name: "RestartPolicy Never",
			pod: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			exempt: false,
		},
		{
			name: "RestartPolicy OnFailure",
			pod: &corev1.Pod{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
				},
			},
			exempt: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isPodExempt(tt.pod)
			assert.Equal(t, tt.exempt, got)
		})
	}
}
