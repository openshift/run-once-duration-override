package main

import (
	"errors"
	"sync"

	"github.com/openshift/run-once-duration-override/pkg/api"
	admissionresponse "github.com/openshift/run-once-duration-override/pkg/response"
	"github.com/openshift/run-once-duration-override/pkg/runoncedurationoverride"
	"k8s.io/klog"

	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	restclient "k8s.io/client-go/rest"
)

type runOnceDurationOverrideHook struct {
	lock        sync.RWMutex
	initialized bool

	admission runoncedurationoverride.Admission
}

// Initialize is called as a post-start hook
func (m *runOnceDurationOverrideHook) Initialize(kubeClientConfig *restclient.Config, stopCh <-chan struct{}) error {
	klog.V(1).Infof("name=%s initializing admission webhook", runoncedurationoverride.Name)
	klog.Infof("name=%s initializing admission webhook", runoncedurationoverride.Name)

	m.lock.Lock()
	defer func() {
		m.initialized = true
		m.lock.Unlock()
	}()

	if m.initialized {
		return nil
	}

	admission, err := runoncedurationoverride.NewInClusterAdmission()
	if err != nil {
		klog.V(1).Infof("name=%s failed to initialize webhook - %s", runoncedurationoverride.Name, err.Error())
		return err
	}

	m.admission = admission

	klog.V(1).Infof("name=%s admission webhook loaded successfully", runoncedurationoverride.Name)

	return nil
}

// MutatingResource is the resource to use for hosting your admission webhook. If the hook implements
// ValidatingAdmissionHook as well, the two resources for validating and mutating admission must be different.
// Note: this is (usually) not the same as the payload resource!
func (m *runOnceDurationOverrideHook) MutatingResource() (plural schema.GroupVersionResource, singular string) {
	return schema.GroupVersionResource{
		Group:    api.Group,
		Version:  api.Version,
		Resource: runoncedurationoverride.Resource,
	}, runoncedurationoverride.Singular
}

// Admit is called to decide whether to accept the admission request. The returned AdmissionResponse may
// use the Patch field to mutate the object from the passed AdmissionRequest.
func (m *runOnceDurationOverrideHook) Admit(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if !m.initialized {
		return admissionresponse.WithInternalServerError(request, errors.New("not initialized"))
	}

	if !m.admission.IsApplicable(request) {
		return admissionresponse.WithAllowed(request)
	}

	exempt, response := m.admission.IsExempt(request)
	if response != nil {
		return response
	}

	if exempt {
		// disabled for this project, do nothing
		return admissionresponse.WithAllowed(request)
	}

	return m.admission.Admit(request)
}
