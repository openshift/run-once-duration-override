// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

// ProjectAccessApplyConfiguration represents a declarative configuration of the ProjectAccess type for use
// with apply.
type ProjectAccessApplyConfiguration struct {
	AvailableClusterRoles []string `json:"availableClusterRoles,omitempty"`
}

// ProjectAccessApplyConfiguration constructs a declarative configuration of the ProjectAccess type for use with
// apply.
func ProjectAccess() *ProjectAccessApplyConfiguration {
	return &ProjectAccessApplyConfiguration{}
}

// WithAvailableClusterRoles adds the given value to the AvailableClusterRoles field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the AvailableClusterRoles field.
func (b *ProjectAccessApplyConfiguration) WithAvailableClusterRoles(values ...string) *ProjectAccessApplyConfiguration {
	for i := range values {
		b.AvailableClusterRoles = append(b.AvailableClusterRoles, values[i])
	}
	return b
}