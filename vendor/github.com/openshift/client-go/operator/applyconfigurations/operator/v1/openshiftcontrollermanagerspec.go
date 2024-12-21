// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	operatorv1 "github.com/openshift/api/operator/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// OpenShiftControllerManagerSpecApplyConfiguration represents a declarative configuration of the OpenShiftControllerManagerSpec type for use
// with apply.
type OpenShiftControllerManagerSpecApplyConfiguration struct {
	OperatorSpecApplyConfiguration `json:",inline"`
}

// OpenShiftControllerManagerSpecApplyConfiguration constructs a declarative configuration of the OpenShiftControllerManagerSpec type for use with
// apply.
func OpenShiftControllerManagerSpec() *OpenShiftControllerManagerSpecApplyConfiguration {
	return &OpenShiftControllerManagerSpecApplyConfiguration{}
}

// WithManagementState sets the ManagementState field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ManagementState field is set to the value of the last call.
func (b *OpenShiftControllerManagerSpecApplyConfiguration) WithManagementState(value operatorv1.ManagementState) *OpenShiftControllerManagerSpecApplyConfiguration {
	b.ManagementState = &value
	return b
}

// WithLogLevel sets the LogLevel field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LogLevel field is set to the value of the last call.
func (b *OpenShiftControllerManagerSpecApplyConfiguration) WithLogLevel(value operatorv1.LogLevel) *OpenShiftControllerManagerSpecApplyConfiguration {
	b.LogLevel = &value
	return b
}

// WithOperatorLogLevel sets the OperatorLogLevel field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the OperatorLogLevel field is set to the value of the last call.
func (b *OpenShiftControllerManagerSpecApplyConfiguration) WithOperatorLogLevel(value operatorv1.LogLevel) *OpenShiftControllerManagerSpecApplyConfiguration {
	b.OperatorLogLevel = &value
	return b
}

// WithUnsupportedConfigOverrides sets the UnsupportedConfigOverrides field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UnsupportedConfigOverrides field is set to the value of the last call.
func (b *OpenShiftControllerManagerSpecApplyConfiguration) WithUnsupportedConfigOverrides(value runtime.RawExtension) *OpenShiftControllerManagerSpecApplyConfiguration {
	b.UnsupportedConfigOverrides = &value
	return b
}

// WithObservedConfig sets the ObservedConfig field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ObservedConfig field is set to the value of the last call.
func (b *OpenShiftControllerManagerSpecApplyConfiguration) WithObservedConfig(value runtime.RawExtension) *OpenShiftControllerManagerSpecApplyConfiguration {
	b.ObservedConfig = &value
	return b
}