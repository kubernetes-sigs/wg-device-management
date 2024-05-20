package api

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DevMgmtAPIVersion = "devmgmtproto.k8s.io/v1alpha1"
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

	// NodeName identifies the node which provides the resources
	// if they are local to a node.
	//
	// A field selector can be used to list only ResourceSlice
	// objects with a certain node name.
	//
	// +optional
	NodeName string `json:"nodeName,omitempty" protobuf:"bytes,2,opt,name=nodeName"`

	// DriverName identifies the DRA driver providing the capacity information.
	// A field selector can be used to list only ResourceSlice
	// objects with a certain driver name.
	DriverName string `json:"driverName" protobuf:"bytes,3,name=driverName"`

	ResourceModel `json:",inline" protobuf:"bytes,4,name=resourceModel"`

	// To be decided: We could split ResourcePool into spec and status or
	// extend NamedDevice below so that it reflects the current state of
	// that device. Like the spec that is data that is coming from a driver,
	// so adding more information is not a conceptual change. Could be
	// added in any future release as an extension.
	//
	// Much harder is tracking "allocated for". That would be useful to
	// enable running multiple schedulers for different sets of nodes where
	// all instances share ownership of some ResourcePool with network-attached
	// devices. It also would enable a DRA driver's controller to co-exist
	// with structured parameters in the scheduler.
	//
	// If we had transactions in the apiserver, we could combine a status
	// update of the ResourcePool with a status update of the claim.  But
	// we don't and adding it [is
	// hard](https://kubernetes.slack.com/archives/C0EG7JC6T/p1714373064352099).
	// Without transactions, there will be a risk of leaking resources
	// and/or races around freeing leaked resources. This would be great to
	// have in 1.31 because requiring schedulers to use the ResourcePool
	// status for allocation will be a change of behavior, but it'll be
	// hard to design and implement in time.
}

// ResourceModel must have one and only one field set.
type ResourceModel struct {
	// NamedDevices describes available devices by listing them.
	//
	// +optional
	NamedDevices *NamedDevices `json:"namedDevices,omitempty" protobuf:"bytes,1,opt,name=namedResources"`
}

// NamedDevices is used in ResourceModel.
type NamedDevices struct {
	// CommonAttributes defines the set of attributes common to all devices in this ResourcePool.
	// The name of each attribute must be unique.
	//
	// +listType=atomic
	// +optional
	CommonAttributes []DeviceAttribute `json:"commonAttributes" protobuf:"bytes,2,name=commonAttributes"`

	// A set of shared resources that get consumed as devices in this
	// ResourcePool get allocated out.
	SharedResources []DeviceResource `json:"sharedResources" protobuf:"bytes,3,name=sharedResources"`

	// The list of all devices currently available.
	//
	// +listType=atomic
	Devices []Device `json:"devices" protobuf:"bytes,1,name=instances"`
}

// Device represents one individual hardware instance that can be selected based
// on its attributes.
type Device struct {
	// Name is unique identifier among all devices managed by
	// the driver on the node. It must be a DNS subdomain.
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	// The set of resources this device consumes from various sources.
	//
	// +listType=atomic
	// +optional
	Resources *DeviceResourceSources `json:"resources" protobuf:"bytes,2,name=resources"`

	// Attributes defines the attributes of this device.
	// The name of each attribute must be unique.
	//
	// +listType=atomic
	// +optional
	Attributes []DeviceAttribute `json:"attributes,omitempty" protobuf:"bytes,3,opt,name=attributes"`
}

// DeviceResource represents a resource that is consumed when a device gets allocated.
type DeviceResource struct {
	// Name is the name of the resource represented by this device resource.
	// It must be a DNS subdomain.
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	// DeviceResourceValue is an embedded type representing the actual value of the device resource.
	DeviceResourceValue `json:",inline" protobuf:"bytes,2,opt,name=value"`
}

// DeviceResourceValue represents the value of a device resource.
// DeviceResourceValue must have one and only one field set.
type DeviceResourceValue struct {
	// QuantityValue is a quantity.
	QuantityValue *resource.Quantity `json:"quantity,omitempty" protobuf:"bytes,1,opt,name=quantity"`

	// IntRangeValue is a range of 64-bit integers.
	IntRangeValue *IntRange `json:"intRange,omitempty" protobuf:"varint,2,rep,name=intRange"`
}

// DeviceResourceSources represents a set of DeviceResources from many possible sources.
type DeviceResourceSources struct {
	// FromSharedPool holds the list of DeviceResources that are consumed from the shared pool.
	FromSharedPool []DeviceResource `json:"fromSharedPool,omitempty" protobuf:"bytes,1,opt,name=fromSharedPool"`

	// Local holds the list of DeviceResources that are local to just this device.
	Local []DeviceResource `json:"local,omitempty" protobuf:"bytes,2,opt,name=local"`
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

	// DeviceAttributeValue is an embedded type representing the actual value of the device attribute.
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
	// IntValue is a 64-bit integer.
	IntValue *int64 `json:"int,omitempty" protobuf:"varint,7,opt,name=int"`
	// IntSliceValue is an array of 64-bit integers.
	IntSliceValue *IntSlice `json:"intSlice,omitempty" protobuf:"varint,8,rep,name=intSlice"`
	// StringValue is a string.
	StringValue *string `json:"string,omitempty" protobuf:"bytes,5,opt,name=string"`
	// StringSliceValue is an array of strings.
	StringSliceValue *StringSlice `json:"stringSlice,omitempty" protobuf:"bytes,9,rep,name=stringSlice"`
	// VersionValue is a semantic version according to semver.org spec 2.0.0.
	VersionValue *string `json:"version,omitempty" protobuf:"bytes,10,opt,name=version"`
}

// IntSlice contains a slice of 64-bit integers.
type IntSlice struct {
	// Ints is the slice of 64-bit integers.
	//
	// +listType=atomic
	Ints []int64 `json:"ints" protobuf:"bytes,1,opt,name=ints"`
}

// StringSlice contains a slice of strings.
type StringSlice struct {
	// Strings is the slice of strings.
	//
	// +listType=atomic
	Strings []string `json:"strings" protobuf:"bytes,1,opt,name=strings"`
}
