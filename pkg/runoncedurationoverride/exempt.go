package runoncedurationoverride

import (
	corev1 "k8s.io/api/core/v1"
)

// Run-once pods are pods that have a RestartPolicy of Never or OnFailure
func isPodExempt(pod *corev1.Pod) bool {
	if pod.Spec.RestartPolicy == corev1.RestartPolicyOnFailure || pod.Spec.RestartPolicy == corev1.RestartPolicyNever {
		return false
	}
	return true
}
