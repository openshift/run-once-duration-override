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
