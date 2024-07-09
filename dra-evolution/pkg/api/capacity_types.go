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

	// DeviceShape defines the common shape of all devices in this pool.
	//
	// +required
	DeviceShape DeviceShape `json:"deviceShape"`

	// Devices lists all available devices in this pool.
	//
	// Must not have more than 128 entries.
	Devices []Device `json:"devices,omitempty"`

	// FUTURE EXTENSION: some other kind of list, should we ever need it.
	// Old clients seeing an empty Devices field can safely ignore the (to
	// them) empty pool.
}

const ResourcePoolMaxSharedCapacity = 128
const ResourcePoolMaxDevices = 128

type DeviceShape struct {
	// Attributes defines the attributes of this device shape.
	// The name of each attribute must be unique.
	//
	// Must not have more than 32 entries.
	//
	// +listType=atomic
	// +optional
	Attributes []DeviceAttribute `json:"attributes,omitempty" protobuf:"bytes,3,opt,name=attributes"`

	// Partitions defines the set of partitions into which this device shape
	// may be allocated. If not populated, then the device shape is always
	// consumed in its entirety.
	//
	// +listType=atomic
	// +optional
	Partitions []DevicePartition `json:"partitions,omitempty"`

	// SharedCapacity defines the set of shared capacity consumable by
	// partitions in this DeviceShape. Not meaninful for non-partitioned
	// devices.
	//
	// Must not have more than 128 entries.
	//
	// +listType=atomic
	// +optional
	SharedCapacity []SharedCapacity `json:"sharedCapacity,omitempty"`
}

// StringOrExpression contains either an explicit string Value or
// a CEL expression that will return a string.
type StringOrExpression struct {
	Value      *string `json:"value,omitempty"`
	Expression *string `json:"expression,omitempty"`
}

// QuantityOrExpression contains either an explicit resource.Quantity Value
// or a CEL expression that results in a resource.Quantity (or a string that parses
// to one).
type QuantityOrExpression struct {
	Value      *resource.Quantity `json:"value,omitempty"`
	Expression *string            `json:"expression,omitempty"`
}

// QuantityOrExpression contains either an explicit bool Value
// or a CEL expression that results in a bool
type BoolOrExpression struct {
	Value      *bool   `json:"value,omitempty"`
	Expression *string `json:"expression,omitempty"`
}

// Device represents a format for a partition, and a count. The actual partitions of
// the device are generated in-memory by evaluating the format for the index values
// 0..Count.
type DevicePartition struct {

	// Count identifies the number of partitions using this format.
	//
	// +required
	Count int `json:"count"`

	// Name is unique identifier among all partitions for this device. The
	// device name as recorded in the allocation will be the concatenation
	// of the device name and the partition name with a '-' separator.
	//
	// NOTE: may need a better naming scheme
	//
	// It must be a DNS label.
	//
	// +required
	Name StringOrExpression `json:"name" protobuf:"bytes,1,name=name"`

	// Attributes defines the attributes of this partition.
	// The name of each attribute must be unique. The values
	// in here are overlayed on top of the values in the device
	// shape (overwriting them if the names are the same).
	//
	// NOTE: probably can get away with fewer
	//
	// Must not have more than 32 entries.
	//
	// +listType=atomic
	// +optional
	Attributes []DeviceAttributeFormat `json:"attributes,omitempty" protobuf:"bytes,3,opt,name=attributes"`

	// SharedCapacityConsumed defines the set of shared capacity consumed by
	// this partition.
	//
	// Must not have more than 32 entries.
	//
	// +listType=atomic
	// +optional
	SharedCapacityConsumed []SharedCapacityFormat `json:"sharedCapacityConsumed,omitempty"`
}

type Device struct {
	// Name is unique identifier among all devices managed by
	// the driver on the node. It must be a DNS label.
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	// Attributes defines the attributes of this device.
	// The name of each attribute must be unique. The values
	// in here are overlayed on top of the values in the device
	// shape (overwriting them if the names are the same).
	//
	// Must not have more than 32 entries.
	//
	// NOTE: probably can get away with fewer
	//
	// +listType=atomic
	// +optional
	Attributes []DeviceAttribute `json:"attributes,omitempty" protobuf:"bytes,3,opt,name=attributes"`
}

const ResourcePoolMaxAttributesPerDevice = 32
const ResourcePoolMaxSharedCapacityConsumedPerDevice = 32

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

type SharedCapacity struct {
	// Name is a unique identifier among all shared capacities managed by the
	// driver in the pool.
	//
	// It is referenced both when defining the total amount of shared capacity
	// that is available, as well as by individual devices when declaring
	// how much of this shared capacity they consume.
	//
	// SharedCapacity names must be a C-style identifier (e.g. "the_name") with
	// a maximum length of 32.
	//
	// By limiting these names to a C-style identifier, the same validation can
	// be used for both these names and the identifier portion of a
	// DeviceAttribute name.
	//
	// +required
	Name string `json:"name"`

	// Capacity is the total capacity of the named resource.
	// This can either represent the total *available* capacity, or the total
	// capacity *consumed*, depending on the context where it is referenced.
	//
	// +required
	Capacity resource.Quantity `json:"capacity"`
}

type DeviceAttributeFormat struct {
	Name StringOrExpression `json:"name"`

	QuantityValue *QuantityOrExpression `json:"quantity,omitempty"`
	BoolValue     *BoolOrExpression     `json:"bool,omitempty"`
	StringValue   *StringOrExpression   `json:"string,omitempty"`
	VersionValue  *StringOrExpression   `json:"version,omitempty"`
}

type SharedCapacityFormat struct {
	Name StringOrExpression `json:"name"`

	Capacity *QuantityOrExpression `json:"capacity,omitempty"`
}

// CStyleIdentifierMaxLength is the maximum length of a c-style identifier used for naming.
const CStyleIdentifierMaxLength = 32

// DeviceAttributeMaxIDLength is the maximum length of the identifier in a device attribute name (`<domain>/<ID>`).
const DeviceAttributeMaxIDLength = CStyleIdentifierMaxLength

// DeviceAttributeMaxValueLength is the maximum length of a string or version attribute value.
const DeviceAttributeMaxValueLength = 64

// SharedCapacityMaxNameLength is the maximum length of a shared capacity name.
const SharedCapacityMaxNameLength = CStyleIdentifierMaxLength
