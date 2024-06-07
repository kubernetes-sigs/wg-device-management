package api

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResourcePool represents a collection of devices managed by a given driver. How
// devices are divided into pools is driver-specific, but typically the
// expectation would a be a pool per identical collection of devices, per node.
// It is fine to have more than one pool for a given node, for the same driver.
//
// Where a device gets published may change over time. The unique identifier
// for a node-local device is the tuple `<driver name>/<node name>/<device name>`. Each
// of these names is a DNS label or domain, so it is okay to concatenate them
// like this in a string with a slash as separator.
//
// For non-local devices, the driver can either make all device names globally
// unique (`<driver name>/<device name>`) or provide a device pool name
// (`<driver name>/<device pool name>/<device name>`).
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
	// Node name and device pool are mutually exclusive. At least one must be set.
	NodeName string `json:"nodeName, omitempty"`

	// DevicePool can be used for non-instead of a NodeName to define where the devices
	// are available.
	//
	// Node name and device pool are mutually exclusive. At least one must be set.
	DevicePool *DevicePool `json:"devicePool,omitempty"`

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

	// FUTURE EXTENSION: some other kind of list, should we ever need it.
	// Old clients seeing an empty Devices field can safely ignore the (to
	// them) empty pool.
}

const ResourcePoolMaxDevices = 128

// DevicePool is a more general description for a set of devices.
// A single node is a special case of this.
type DevicePool struct {
	// This name together with the driver name and the device name
	// identify a device.
	//
	// Must not be longer than 253 characters and may contain one or more
	// DNS sub-domains separated by slashes.
	//
	// +optional
	Name string `json:"name,omitempty"`

	// This identifies nodes where the devices may be used. If unset,
	// all nodes have access.
	//
	// +optional
	NodeSelector *v1.NodeSelector `json:"nodeSelector,omitempty"`
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
	// Must not have more than 32 entries.
	//
	// +listType=atomic
	// +optional
	Attributes []DeviceAttribute `json:"attributes,omitempty" protobuf:"bytes,3,opt,name=attributes"`

	// TODO for 1.31: define how to support partitionable devices
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
