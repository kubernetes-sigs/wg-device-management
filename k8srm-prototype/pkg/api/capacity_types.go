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
	// +required
	Driver string `json:"driver,omitempty"`

	// Attributes contains device attributes that are common to all devices
	// in the pool.
	// +optional
	Attributes []Attribute `json:"attributes,omitempty"`

	// Resources are pooled resources that are shared by all devices in the
	// pool. This is typically used when representing a partitionable
	// device, and need not be populated otherwise.
	//
	// +optional
	Resources []ResourceCapacity `json:"resources,omitempty"`

	// Devices contains the individual devices in the pool.
	//
	// +required
	Devices []Device `json:"devices,omitempty"`
}

// DevicePoolStatus contains the state of the pool as last reported by the
// driver. Note that this will not include the allocations that have been made
// by the scheduler but not yet seen by the driver. Thus, it is NOT sufficient
// to make future scheduling decisions.
type DevicePoolStatus struct {
	AvailableDevices int `json:"availableDevices,omitempty"`
}

// Device is used to track individual devices in a pool.
type Device struct {
	// Name is a driver-specific identifier for the device.
	// +required
	Name string `json:"name"`

	// Attributes contain additional metadata that can be used in
	// constraints. If an attribute name overlaps with the pool attribute,
	// the device attribute takes precedence.
	//
	// +optional
	Attributes []Attribute `json:"attributes,omitempty"`

	// Requests contains the pooled resources that are consumed when
	// this device is allocated.
	//
	// +optional
	Requests map[string]resource.Quantity `json:"requests,omitempty"`

	// Resources allows the definition of per-device resources that can
	// be allocated in a manner similar to standard Kubernetes resources.
	// NOTE: We may not need to implement this right away.
	//
	// +optional
	Resources []ResourceCapacity `json:"resources,omitempty"`
}

type ResourceCapacity struct {
	// Name is the resource name/type.
	Name string `json:"name"`

	// Capacity is the total capacity of the named resource.
	// +required
	Capacity resource.Quantity `json:"capacity"`

	// BlockSize is the increments in which capacity is consumed. For
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

	// If we ever need per-resource topology, then topology associations
	// would also go in here.
}

// Attribute capture the name, value, and type of an device attribute.
type Attribute struct {
	Name string `json:"name"`

	// One of the following:
	StringValue   *string            `json:"stringValue,omitempty"`
	IntValue      *int               `json:"intValue,omitempty"`
	QuantityValue *resource.Quantity `json:"quantityValue,omitempty"`
	SemVerValue   *SemVer            `json:"semVerValue,omitempty"`
}

func (a Attribute) Equal(b Attribute) bool {
	if a.Name != b.Name {
		return false
	}

	return a.EqualValue(b)
}

func (a Attribute) EqualValue(b Attribute) bool {
	if a.StringValue != nil && b.StringValue != nil && *a.StringValue == *b.StringValue {
		return true
	}

	if a.IntValue != nil && b.IntValue != nil && *a.IntValue == *b.IntValue {
		return true
	}

	if a.QuantityValue != nil && b.QuantityValue != nil && (*a.QuantityValue).Equal(*b.QuantityValue) {
		return true
	}

	if a.SemVerValue != nil && b.SemVerValue != nil && *a.SemVerValue == *b.SemVerValue {
		return true
	}

	return false
}

// SemVer represents a semantic version value. In this prototype it is just a
// string.
type SemVer string
