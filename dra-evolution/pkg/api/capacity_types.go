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
//
// Consumers should be prepared to handle situations where the same device is
// listed in different pools, for example because the producer already added it
// to a new pool before removing it from an old one. Should this occur, then
// there is still only one such device instance. If the two device definitions
// disagree in any way, the one found in the newest ResourcePool, as determined
// by creationTimestamp, is preferred.
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
	//
	// Must not have more than 128 entries.
	Devices []Device `json:"devices,omitempty"`

	// SharedCapacity are pooled resources that are shared by all
	// devices in the pool. This is typically used when representing a
	// partitionable device, and need not be populated otherwise.
	//
	// Must not have more than 32 entries.
	//
	// +optional
	// +listType=atomic
	SharedCapacity []SharedCapacity `json:"sharedCapacity,omitempty"`

	// FUTURE EXTENSION: some other kind of list, should we ever need it.
	// Old clients seeing an empty Devices field can safely ignore the (to
	// them) empty pool.
}

type ResourceCapacity struct {
	// Name is the resource name/type.
	// +required
	Name string `json:"name"`

	// Capacity is the total capacity of the named resource.
	// +required
	Capacity resource.Quantity `json:"capacity"`
}

const ResourcePoolMaxDevices = 128
const ResourcePoolMaxSharedCapacity = 32

// Device represents one individual hardware instance that can be selected based
// on its attributes.
type Device struct {
	// Name is unique identifier among all devices managed by
	// the driver on the node. It must be a DNS label.
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	// Attributes defines the attributes of this device.
	// The name of each attribute must be unique.
	//
	// Must not have more than 32 entries.
	//
	// +listType=atomic
	// +optional
	Attributes []DeviceAttribute `json:"attributes,omitempty" protobuf:"bytes,3,opt,name=attributes"`

	// SharedCapacityConsumed contains the pooled allocatable resources
	// that are consumed when this device is allocated.
	//
	// Must not have more than 32 entries.
	//
	// +optional
	SharedCapacityConsumed []SharedCapacityRequest `json:"sharedCapacityConsumed,omitempty"`
}

const ResourcePoolMaxAttributesPerDevice = 32

// ResourcePoolMaxDevices and ResourcePoolMaxAttributesPerDevice where chosen
// so that with the maximum attribute length of 96 characters the total size of
// the ResourcePool object is around 420KB.

// DeviceAttribute is a combination of an attribute name and its value.
// Exactly one value must be set.
type DeviceAttribute struct {
	// Name is a unique identifier for this attribute, which will be
	// referenced when selecting devices.
	//
	// Attributes are defined either by the owner of the specific driver
	// (usually the vendor) or by some 3rd party (e.g. the Kubernetes
	// project). Because attributes are sometimes compared across devices,
	// a given name is expected to mean the same thing and have the same
	// type on all devices.
	//
	// Attribute names must be either a C-style identifier
	// (e.g. "the_name") or a DNS subdomain followed by a slash ("/")
	// followed by a C-style identifier
	// (e.g. "example.com/the_name"). Attributes whose name does not
	// include the domain prefix are assumed to be part of the driver's
	// domain. Attributes defined by 3rd parties must include the domain
	// prefix.
	//
	// The maximum length for the DNS subdomain is 63 characters (same as
	// for driver names) and the maximum length of the C-style identifier
	// is 32.
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	// The Go field names below have a Value suffix to avoid a conflict between the
	// field "String" and the corresponding method. That method is required.
	// The Kubernetes API is defined without that suffix to keep it more natural.

	// QuantityValue is a quantity.
	QuantityValue *resource.Quantity `json:"quantity,omitempty" protobuf:"bytes,2,opt,name=quantity"`
	// BoolValue is a true/false value.
	BoolValue *bool `json:"bool,omitempty" protobuf:"bytes,3,opt,name=bool"`
	// StringValue is a string. Must not be longer than 64 characters.
	StringValue *string `json:"string,omitempty" protobuf:"bytes,4,opt,name=string"`
	// VersionValue is a semantic version according to semver.org spec 2.0.0.
	// Must not be longer than 64 characters.
	VersionValue *string `json:"version,omitempty" protobuf:"bytes,5,opt,name=version"`
}

const DeviceAttributeMaxIDLength = 32
const DeviceAttributeMaxValueLength = 64

// SharedCapacity is a set of resources which can be drawn down upon.
type SharedCapacity struct {
	// Name is the name of this set of shared capacity, which is unique
	// resource pool.
	// +required
	Name string `json:"name"`

	// Resources are the list of resources provided in this set.
	Resources []ResourceCapacity `json:"resources,omitempty"`
}

// SharedCapacityRequest is a request to draw down resources from a particular
// SharedCapacity.
type SharedCapacityRequest struct {

	// Name is the name of the SharedCapacity from which to draw the
	// resources.
	//
	// +required
	Name string `json:"name"`

	// Resources are the list of resources and the amount of resources
	// to draw down.
	Resources []ResourceCapacity `json:"resources,omitempty"`
}
