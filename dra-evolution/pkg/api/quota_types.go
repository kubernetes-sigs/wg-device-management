package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Quota controls whether a ResourceClaim may get allocated.
// Quota is namespaced and applies to claims within the same namespace.
type Quota struct {
	metav1.TypeMeta
	// Standard object metadata.
	metav1.ObjectMeta

	Spec QuotaSpec
}

type QuotaSpec struct {
	// AllowManagementAccess controls whether claims with ManagementAccess
	// may be allocated. If multiple quota objects exist and at least one
	// has a true value, access will be allowed. The default to deny such access.
	//
	// +optional
	AllowManagementAccess bool `json:"allowManagementAccess,omitempty"`

	// Stretch goals for 1.31:
	//
	// - maximum number of devices matching a selector
	// - maximum sum of a certain quantity attribute
	//
	// These are additional requirements when checking whether a device
	// instance can satisfy a request. Creating a claim is always allowed,
	// but allocating it fails when the quota is currently exceeded.

	// Other useful future extensions (>= 1.32):

	// DeviceLimits is a CEL expression which take the currently allocated
	// devices and their attributes and some new allocations as input and
	// returns false if those allocations together are not permitted in the
	// namespace.
	//
	// DeviceLimits string

	// A class listed in DeviceClassDenyList must not be used in this
	// namespace. This can be useful for classes which contain
	// configuration pieces that a user in this namespace should not have
	// access to.
	//
	// DeviceClassDenyList []string

	// A class listed in ResourceClassAllowList may be used in this namespace
	// even when that class is marked as "privileged". Normally classes
	// are not privileged and using them does not require explicit listing
	// here, but some classes may contain more sensitive configurations
	// that not every user should have access to.
	//
	// DeviceClassAllowList []string
}
