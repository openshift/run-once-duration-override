package runoncedurationoverride

import (
	"errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/integer"
)

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
	current.Spec.ActiveDeadlineSeconds = int64MinP(&activeDeadlineSeconds, current.Spec.ActiveDeadlineSeconds)

	return current, nil
}

func int64MinP(a, b *int64) *int64 {
	switch {
	case a == nil:
		return b
	case b == nil:
		return a
	default:
		c := integer.Int64Min(*a, *b)
		return &c
	}
}
