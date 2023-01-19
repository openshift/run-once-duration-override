package main

import (
	"sync"

	"github.com/openshift/run-once-duration-override/pkg/runoncedurationoverride"
	"k8s.io/klog"

	restclient "k8s.io/client-go/rest"
)

type runOnceDurationOverrideHook struct {
	lock        sync.RWMutex
	initialized bool
}

// Initialize is called as a post-start hook
func (m *runOnceDurationOverrideHook) Initialize(kubeClientConfig *restclient.Config, stopCh <-chan struct{}) error {
	klog.V(1).Infof("name=%s initializing admission webhook", runoncedurationoverride.Name)

	m.lock.Lock()
	defer func() {
		m.initialized = true
		m.lock.Unlock()
	}()

	if m.initialized {
		return nil
	}

	return nil
}
