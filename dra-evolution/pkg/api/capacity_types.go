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
	// The list of all devices currently available.
	//
	// +listType=atomic
	Devices []NamedDevice `json:"devices" protobuf:"bytes,1,name=instances"`
}

// NamedDevice represents one individual hardware instance that can be selected based
// on its attributes.
type NamedDevice struct {
	// Name is unique identifier among all devices managed by
	// the driver on the node. It must be a DNS subdomain.
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	// Attributes defines the attributes of this device.
	// The name of each attribute must be unique.
	//
	// +listType=atomic
	// +optional
	Attributes []DeviceAttribute `json:"attributes,omitempty" protobuf:"bytes,2,opt,name=attributes"`
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
