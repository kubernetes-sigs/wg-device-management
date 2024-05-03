package api

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DevMgmtAPIVersion = "devmgmtproto.k8s.io/v1alpha1"
)

// DevicePool represents a collection of devices managed by a given driver. How
// devices are divided into pools is driver-specific, but typically the
// expectation would a be a pool per identical collection of devices, per node.
// It is fine to have more than one pool for a given node, for the same driver.
type DevicePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevicePoolSpec   `json:"spec,omitempty"`
	Status DevicePoolStatus `json:"status,omitempty"`
}

// DevicePoolSpec identifies the driver and contains the data for the pool
// prior to any allocations.
// NOTE: It's not clear that spec/status is the right model for this data.
type DevicePoolSpec struct {
	// NodeName is the name of the node containing the devices in the pool.
	// For network attached devices, this may be empty.
	// +optional
	NodeName *string `json:"nodeName,omitempty"`

	// Driver is the name of the DeviceDriver that created this object and
	// owns the data in it.
	//
	// +required
	DriverName string `json:"driverName"`

	// Attributes contains device attributes that are common to all devices
	// in the pool.
	//
	// +optional
	Attributes []Attribute `json:"attributes,omitempty"`

	// SharedResources are resources that are shared by all devices in the
	// pool. This is typically used when representing a partitionable
	// device, and need not be populated otherwise.
	//
	// +optional
	SharedResources []DeviceResource `json:"sharedResources,omitempty"`

	// Devices contains the individual devices in the pool.
	//
	// +required
	Devices []Device `json:"devices"`
}

// DevicePoolStatus contains the state of the pool as last reported by the
// driver. Note that this will not include the allocations that have been made
// by the scheduler but not yet seen by the driver. Thus, it is NOT sufficient
// to make future scheduling decisions.
type DevicePoolStatus struct {
	AllocatedDevices []AllocatedDevice `json:"allocatedDevices,omitempty"`
}

// AllocatedDevice represents a device that has been allocated from the pool.
type AllocatedDevice struct {
	Name string
	ClaimUID types.UID
}

// Device is used to track individual devices in a pool.
type Device struct {
	// Name is a driver-specific identifier for the device.
	//
	// +required
	Name string `json:"name"`

	// Attributes contain additional metadata that can be used in
	// constraints. If an attribute name overlaps with the pool attribute,
	// the device attribute takes precedence.
	//
	// +optional
	Attributes []Attribute `json:"attributes,omitempty"`

	// SharedResourceRequests contains requests for some subset of the shared
	// resources available in the overall pool. They represent the set of
	// shared resources consumed by this device when it gets allocated.
	//
	// +optional
	SharedResourceRequests []DeviceResourceRequest `json:"sharedResourceRequests,omitempty"`
}

// DeviceResource represents a named resource that is consumable by a device.
// We indirect through DeviceResourceValue to enforce a one-of-many for the
// different types of resources this struct can represent.
type DeviceResource struct {
	Name string `json:"name"`
	*DeviceResourceValue `json:",inline"`
}

// DeviceResourceValue holds the actual value of the named device resource.
// Only one of the fields in this struct can be set for any given instance.
//
// +kubebuilder:validation:MaxProperties=1
type DeviceResourceValue struct {
	Quantity *DeviceResourceQuantity `json:"quantity"`
	IntRange *DeviceResourceIntRange `json:"intrange"`
}

// DeviceResourceRequest represents a request to consume a named resource
// from a device. We indirect through DeviceResourcRequestValue to enforce a
// one-of-many for the different types of resources this struct can represent.
type DeviceResourceRequest struct {
	Name string `json:"name"`
	*DeviceResourceRequestValue `json:",inline"`
}

// DeviceResourceRequestValue holds the actual value of the request for the
// named device resource. Only one of the fields in this struct can be set for
// any given instance.
//
// +kubebuilder:validation:MaxProperties=1
type DeviceResourceRequestValue struct {
	Quantity *resource.Quantity `json:"quantity"`
	IntRange *resource.IntRange `json:"intrange"`
}

// DeviceResourceQuantity represents a consumable resource on a device as a
// quantity that is allocatable with a given block size.
type DeviceResourceQuantity struct {
	// Value represents the actual quantity that this resource holds.
	//
	// +required
	Value resource.Quantity `json:"value"`

	// BlockSize is the increments in which quantity is consumed. For
	// example, if you can only allocate memory in 4k pages, then the
	// block size should be "4Ki". Default is 1.
	//
	// If the resource is consumable in a fractional way, then the
	// default of 1 should not be used; instead this should be a fractional
	// amount corresponding the increment size. We may also need a minimum
	// value, if the minimum is larger than the block size (as is the case
	// for standard Kubernetes CPU resources).
	//
	// +optional
	BlockSize *resource.Quantity `json:"blockSize,omitempty"`
}

// DeviceResourceIntRange represents a consumable resource on a device as a
// discrete range of integers.
type DeviceResourceIntRange struct {
	// Value represents the actual range of integers that this resource holds.
	//
	// +required
	Value resource.IntRange `json:"value"`
}

// DeviceAttribute capture the name, value, and type of a device attribute.
// We indirect through DeviceAttributeValue to enforce a one-of-many for the
// different types of attributes this struct can represent.
type DeviceAttribute struct {
	Name string `json:"name"`
	*DeviceAttributeValue `json:",inline"`
}

// DeviceAttributeValue holds the actual value of the attribute for the device.
// Only one of the fields in this struct can be set for any given instance.
//
// +kubebuilder:validation:MaxProperties=1
type DeviceAttributeValue struct {
	StringValue   *string            `json:"string,omitempty"`
	IntValue      *int               `json:"int,omitempty"`
	QuantityValue *resource.Quantity `json:"quantity,omitempty"`
	SemVerValue   *SemVer            `json:"semver,omitempty"`
}
