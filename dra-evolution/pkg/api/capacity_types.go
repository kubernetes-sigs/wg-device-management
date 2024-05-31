package api

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResourcePool represents a collection of devices managed by a given driver. How
// devices are divided into pools is driver-specific, but typically the
// expectation would a be a pool per identical collection of devices, per node.
// It is fine to have more than one pool for a given node, for the same driver.
//
// Where a device gets published may change over time. The unique identifier
// for a device is the tuple `<driver name>/<node name>/<device name>`. Each
// of these names is a DNS label or domain, so it is okay to concatenate them
// like this in a string with a slash as separator.
type ResourcePool struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec ResourcePoolSpec `json:"spec"`

	// Stretch goal for 1.31: define status
	//
	// To discuss:
	// - Who writes that status?
	//   After https://github.com/kubernetes/kubernetes/pull/125163 (implements
	//   https://github.com/kubernetes/enhancements/pull/4667), kubelet is not
	//   involved with publishing ResourcePools, they come directly from the driver.
	// - What information should be in it?
	// - Does it need to be a sub-resource? This would make it harder
	//   for a driver to publish a new device (spec) and its corresponding
	//   health information (status).
}

type ResourcePoolSpec struct {
	// NodeName identifies the node which provides the devices. All devices
	// are local to that node.
	//
	// This is currently required, but this might get relaxed in the future.
	NodeName string `json:"nodeName"`

	// POTENTIAL FUTURE EXTENSION: NodeSelector *v1.NodeSelector

	// DriverName identifies the DRA driver providing the capacity information.
	// A field selector can be used to list only ResourcePool
	// objects with a certain driver name.
	//
	// Must be a DNS subdomain and should end with a DNS domain owned by the
	// vendor of the driver.
	DriverName string `json:"driverName" protobuf:"bytes,3,name=driverName"`

	// Devices lists all available devices in this pool.
	Devices []Device `json:"devices,omitempty"`

	// FUTURE EXTENSION: some other kind of list, should we ever need it.
	// Old clients seeing an empty Devices field can safely ignore the (to
	// them) empty pool.
}

// Device represents one individual hardware instance that can be selected based
// on its attributes.
type Device struct {
	// Name is unique identifier among all devices managed by
	// the driver on the node. It must be a DNS label.
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
	// It must be a DNS subdomain, with one additional restriction: the
	// first part must not contain a hyphen and not start with a digit.
	//
	// If this is a DNS label (no dot), then the driver name gets added
	// when looking up attributes. This avoids name collisions with attributes
	// used by other drivers.
	//
	// If this is a full DNS subdomain, then the meaning of the attribute is driver-independent.
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
