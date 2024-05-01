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

	// DeviceCount contains the total number of devices in the pool.
	// +required
	DeviceCount int `json:"count,omitempty"`
}

// DevicePoolStatus contains the state of the pool as last reported by the
// driver. Note that this will not include the allocations that have been made
// by the scheduler but not yet seen by the driver. Thus, it is NOT sufficient
// to make future scheduling decisions.
type DevicePoolStatus struct {
	AvailableDevices int `json:"availableDevices,omitempty"`
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
