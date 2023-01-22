package runoncedurationoverride

import (
	"errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type CPUMemory struct {
	CPU    *resource.Quantity
	Memory *resource.Quantity
}

func NewMutator(config *Config) (mutator *podMutator, err error) {
	if config == nil {
		err = errors.New("NewMutator: invalid input")
		return
	}

	return &podMutator{activeDeadlineSeconds: config.ActiveDeadlineSeconds}, nil
}

type podMutator struct {
	activeDeadlineSeconds int64
}

func (m *podMutator) Mutate(in *corev1.Pod) (out *corev1.Pod, err error) {
	current := in.DeepCopy()

	activeDeadlineSeconds := m.activeDeadlineSeconds
	current.Spec.ActiveDeadlineSeconds = &activeDeadlineSeconds

	return current, nil
}
