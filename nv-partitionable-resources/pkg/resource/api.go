package resource

import (
	"k8s.io/apimachinery/pkg/api/resource"
)

// QualifiedName is the name of a device attribute or capacity.
//
// Attributes and capacities are defined either by the owner of the specific
// driver (usually the vendor) or by some 3rd party (e.g. the Kubernetes
// project). Because they are sometimes compared across devices, a given name
// is expected to mean the same thing and have the same type on all devices.
//
// Names must be either a C identifier (e.g. "theName") or a DNS subdomain
// followed by a slash ("/") followed by a C identifier
// (e.g. "dra.example.com/theName"). Names which do not include the
// domain prefix are assumed to be part of the driver's domain. Attributes
// or capacities defined by 3rd parties must include the domain prefix.
//
// The maximum length for the DNS subdomain is 63 characters (same as
// for driver names) and the maximum length of the C identifier
// is 32.
type QualifiedName string

// ResourceSliceSpec contains the information published by the driver in one ResourceSlice.
// +k8s:deepcopy-gen=true
type ResourceSliceSpec struct {
	// Devices lists some or all of the devices in this pool.
	//
	// Must not have more than 128 entries.
	//
	// +optional
	// +listType=atomic
	Devices []Device `json:"devices,omitempty"`

	// DeviceMixins represents a list of device mixins, i.e. a collection of
	// shared attributes and capacities that an actual device can "include"
	// to extend the set of attributes and capacities it already defines.
	//
	// The main purposes of these mixins is to reduce the memory footprint
	// of devices since they can reference the mixins provided here rather
	// than duplicate them.
	//
	// Must not have more than 128 entries.
	//
	// +optional
	// +listType=atomic
	DeviceMixins []DeviceMixin `json:"deviceMixins,omitempty"`
}

// DeviceMixin defines a specific device mixin for each device type.
// Besides the name, exactly one field must be set.
// +k8s:deepcopy-gen=true
type DeviceMixin struct {
	// Name is a unique identifier among all mixins managed by the driver
	// in the pool. It must be a DNS label.
	//
	// +required
	Name string `json:"name"`

	// Partitionable defines a mixin usable by a partitionable device.
	//
	// +optional
	// +oneOf=deviceMixinType
	Partitionable *PartitionableDeviceMixin `json:"partitionable,omitempty"`
}

// PartitionableDeviceMixin defines a mixin that a partitionable device can include.
// +k8s:deepcopy-gen=true
type PartitionableDeviceMixin struct {
	// Attributes defines the set of attributes for this mixin.
	// The name of each attribute must be unique in that set.
	//
	// To ensure this uniqueness, attributes defined by the vendor
	// must be listed without the driver name as domain prefix in
	// their name. All others must be listed with their domain prefix.
	//
	// Conflicting attributes from those provided via other mixins are
	// overwritten by the ones provided here.
	//
	// The maximum number of attributes and capacities combined is 32.
	//
	// +optional
	Attributes map[QualifiedName]DeviceAttribute `json:"attributes,omitempty"`

	// Capacity defines the set of capacities for this mixin.
	// The name of each capacity must be unique in that set.
	//
	// To ensure this uniqueness, capacities defined by the vendor
	// must be listed without the driver name as domain prefix in
	// their name. All others must be listed with their domain prefix.
	//
	// Conflicting capacities from those provided via other mixins are
	// overwritten by the ones provided here.
	//
	// The maximum number of attributes and capacities combined is 32.
	//
	// +optional
	Capacity map[QualifiedName]DeviceCapacity `json:"capacity,omitempty"`
}

// Device represents one individual hardware instance that can be selected based
// on its attributes. Besides the name, exactly one field must be set.
// +k8s:deepcopy-gen=true
type Device struct {
	// Name is unique identifier among all devices managed by
	// the driver in the pool. It must be a DNS label.
	//
	// +required
	Name string `json:"name"`

	// Paartitionable defines one partitionable device instance.
	//
	// +optional
	// +oneOf=deviceType
	Partitionable *PartitionableDevice `json:"partitionable,omitempty"`
}

// PartitionableDevice defines one device instance.
// +k8s:deepcopy-gen=true
type PartitionableDevice struct {
	// Includes defines the set of device mixins that this device includes.
	//
	// The propertes of each included mixin are applied to this device in
	// order. Conflicting properties from multiple mixins are taken from the
	// last mixin listed that contains them.
	//
	// The maximum number of mixins that can be included is 8.
	//
	// +optional
	Includes []DeviceMixinRef `json:"includes,omitempty"`

	// ConsumesCapacityFrom defines the set of devices where any capacity
	// consumed by this device should be pulled from. This applies recursively.
	// In cases where the dvice names itself as its source, the recursion is
	// halted.
	//
	// Conflicting capacities from multiple devices are taken from the
	// last device listed that contains them.
	//
	// The maximum number of devices that can be referenced is 8.
	//
	// +optional
	ConsumesCapacityFrom []DeviceRef `json:"consumesCapacityFrom,omitempty"`

	// Attributes defines the set of attributes for this device.
	// The name of each attribute must be unique in that set.
	//
	// To ensure this uniqueness, attributes defined by the vendor
	// must be listed without the driver name as domain prefix in
	// their name. All others must be listed with their domain prefix.
	//
	// Conflicting attributes from those provided via mixins are
	// overwritten by the ones provided here.
	//
	// The maximum number of attributes and capacities combined is 32.
	//
	// +optional
	Attributes map[QualifiedName]DeviceAttribute `json:"attributes,omitempty"`

	// Capacity defines the set of capacities for this device.
	// The name of each capacity must be unique in that set.
	//
	// To ensure this uniqueness, capacities defined by the vendor
	// must be listed without the driver name as domain prefix in
	// their name. All others must be listed with their domain prefix.
	//
	// Conflicting capacities from those provided via mixins are
	// overwritten by the ones provided here.
	//
	// The maximum number of attributes and capacities combined is 32.
	//
	// +optional
	Capacity map[QualifiedName]DeviceCapacity `json:"capacity,omitempty"`
}

// DeviceMixinRef defines a reference to a device mixin.
// +k8s:deepcopy-gen=true
type DeviceMixinRef struct {
	// Name refers to the name of a device mixin in the pool.
	//
	// +required
	Name string `json:"name"`
}

// DeviceRef defines a reference to a device.
// +k8s:deepcopy-gen=true
type DeviceRef struct {
	// Name refers to the name of a device in the pool.
	//
	// +required
	Name string `json:"name"`
}

// DeviceAttribute must have exactly one field set.
// +k8s:deepcopy-gen=true
type DeviceAttribute struct {
	// The Go field names below have a Value suffix to avoid a conflict between the
	// field "String" and the corresponding method. That method is required.
	// The Kubernetes API is defined without that suffix to keep it more natural.

	// IntValue is a number.
	//
	// +optional
	// +oneOf=ValueType
	IntValue *int64 `json:"int,omitempty"`

	// BoolValue is a true/false value.
	//
	// +optional
	// +oneOf=ValueType
	BoolValue *bool `json:"bool,omitempty"`

	// StringValue is a string. Must not be longer than 64 characters.
	//
	// +optional
	// +oneOf=ValueType
	StringValue *string `json:"string,omitempty"`

	// VersionValue is a semantic version according to semver.org spec 2.0.0.
	// Must not be longer than 64 characters.
	//
	// +optional
	// +oneOf=ValueType
	VersionValue *string `json:"version,omitempty"`
}

// DeviceCapacity defines consumable capacity of a device.
// +k8s:deepcopy-gen=true
type DeviceCapacity struct {
	// Quantity defines how much of a certain device capacity is available.
	Quantity resource.Quantity `json:"quantity,omitempty"`

	// potential future addition: fields which define how to "consume"
	// capacity (= share a single device between different consumers).
}
