package api

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResourcePool represents a collection of devices managed by a given driver. How
// devices are divided into pools is driver-specific, but typically the
// expectation would a be a pool per identical collection of devices, per node.
// It is fine to have more than one pool for a given node, for the same driver.
type ResourcePool struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec ResourcePoolSpec `json:"spec"`

	// Stretch goal for 1.31: define status
}

type ResourcePoolSpec struct {
	ResourcePoolAccessability `json:",inline"`

	// DriverName identifies the DRA driver providing the capacity information.
	// A field selector can be used to list only ResourceSlice
	// objects with a certain driver name.
	DriverName string `json:"driverName" protobuf:"bytes,3,name=driverName"`

	// If this field is true, devices must be reserved before they can be
	// allocated for a claim. The protocol for that has not been defined
	// yet, so currently setting it to true is not valid. Clients must
	// check nonetheless because that might change in the future. If
	// a client does not support distributed allocation, it must ignore
	// the pool.
	DistributedAllocation *bool `json:"distributedAllocation"`

	// Devices lists all available devices in this pool.
	Devices []Device `json:"devices,omitempty"`

	// FUTURE EXTENSION: some other kind of list, should we ever need it.
	// Old clients seeing an empty Devices field can safely ignore the (to
	// them) empty pool.
}

// Exactly one field must be set. Clients which see an empty
// field must ignore the pool when looking for devices.
type ResourcePoolAccessability struct {
	// NodeName identifies the node which provides the devices
	// if (and only if) they are local to a node. Support for other
	// kind of devices may get added in the future, in which case
	// this field will be empty.
	//
	// +optional
	NodeName *string `json:"nodeName,omitempty"`

	// FUTURE EXTENSION:
	// The devices in the pool are not local to any particular
	// node. They are accessible from any node matching
	// this selector. The empty selector matches all nodes.
	// NodeSelector *v1.NodeSelector
}

// Device represents one individual hardware instance that can be selected based
// on its attributes.
type Device struct {
	// Name is unique identifier among all devices managed by
	// the driver on the node. It must be a DNS subdomain.
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	// Attributes defines the attributes of this device.
	// The name of each attribute must be unique.
	//
	// +listType=atomic
	// +optional
	Attributes []DeviceAttribute `json:"attributes,omitempty" protobuf:"bytes,3,opt,name=attributes"`

	// TODO for 1.31: define how to support partitionable devices
}

// DeviceAttribute is a combination of an attribute name and its value.
type DeviceAttribute struct {
	// Name is a unique identifier across all possible attributes of devices.
	//
	// If this is a DNS subdomain (no dot), then the driver name gets added
	// when looking up attributes. This avoids name collisions with attributes
	// used by other drivers.
	//
	// If this is a full DNS domain, then the meaning of the attribute driver-independent.
	// For example, Kubernetes will use `*.k8s.io` names when defining attributes that
	// drivers from different vendors are supposed to use.
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	DeviceAttributeValue `json:",inline" protobuf:"bytes,2,opt,name=attributeValue"`
}

// The Go field names below have a Value suffix to avoid a conflict between the
// field "String" and the corresponding method. That method is required.
// The Kubernetes API is defined without that suffix to keep it more natural.

// DeviceAttributeValue must have one and only one field set.
type DeviceAttributeValue struct {
	// QuantityValue is a quantity.
	QuantityValue *resource.Quantity `json:"quantity,omitempty" protobuf:"bytes,6,opt,name=quantity"`
	// BoolValue is a true/false value.
	BoolValue *bool `json:"bool,omitempty" protobuf:"bytes,2,opt,name=bool"`
	// StringValue is a string.
	StringValue *string `json:"string,omitempty" protobuf:"bytes,5,opt,name=string"`
	// VersionValue is a semantic version according to semver.org spec 2.0.0.
	VersionValue *string `json:"version,omitempty" protobuf:"bytes,10,opt,name=version"`
}
