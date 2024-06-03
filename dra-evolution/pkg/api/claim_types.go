package api

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// DeviceClass is a vendor or admin-provided resource that contains
// device configuration and requirements. It can be referenced in
// the device requests of a claim to apply these presets.
// Cluster scoped.
type DeviceClass struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Config defines configuration parameters that apply to each device which is claimed via this class.
	// Some classses may potentially be satisfied by multiple drivers, so each instance of a vendor
	// configuration applies to exactly one driver.
	//
	// They are passed to the driver, but are not consider while allocating the claim.
	//
	// +optional
	// +listType=atomic
	Config []ConfigurationParameters `json:"config,omitempty"`

	// Requirements must be satisfied by devices. Applies to all devices of
	// a claim when the claim references the class and only to the devices
	// in a request when referenced there.
	//
	// +optional
	// +listType=atomic
	Requirements []Requirement `json:"requirements,omitempty" protobuf:"bytes,4,opt,name=requirements"`
}

// ConfigurationParameters must have one and only one field set.
type ConfigurationParameters struct {
	Opaque *OpaqueConfigurationParameters `json:"opaque,omitempty" protobuf:"bytes,1,opt,name=opaque"`
}

// OpaqueConfigurationParameters contains configuration parameters for a driver.
type OpaqueConfigurationParameters struct {
	// DriverName is used to determine which kubelet plugin needs
	// to be passed these configuration parameters.
	//
	// An admission webhook provided by the driver developer could use this
	// to decide whether it needs to validate them.
	//
	// Must be a DNS subdomain and should end with a DNS domain owned by the
	// vendor of the driver.
	DriverName string `json:"driverName" protobuf:"bytes,1,name=driverName"`

	// Parameters can contain arbitrary data. It is the responsibility of
	// the driver developer to handle validation and versioning. Typically this
	// includes self-identification and a version ("kind" + "apiVersion" for
	// Kubernetes types), with conversion between different versions.
	Parameters runtime.RawExtension `json:"parameters,omitempty" protobuf:"bytes,2,opt,name=parameters"`
}

// Requirement must have one and only one field set.
type Requirement struct {
	// This CEL expression which must evaluate to true if a
	// device is suitable. This covers qualitative aspects of
	// device selection.
	//
	// The language is as defined in
	// https://kubernetes.io/docs/reference/using-api/cel/
	// with several additions that are specific to device selectors.
	//
	// For each attribute type there is a
	// `device.<type>Attributes` map that resolves to the corresponding
	// value of the instance under evaluation. The type of those map
	// entries are known at compile time, which makes it easier to
	// detect errors like string to int comparisons.
	//
	// In cases where such type safety is not needed or not desired,
	// `device.attributes` can be used instead. The type of the entries
	// then only gets checked at runtime.
	//
	// Unknown keys are not an error. Instead, `device.<type>Attributes`
	// returns a default value for each type:
	// - empty string
	// - false for a boolean
	// - zero quantity
	// - 0.0.0 for a version
	//
	// `device.attributes` returns nil.
	//
	// The `device.driverName` string variable can be used to check for a specific
	// driver explicitly in a filter that is meant to work for devices from
	// different vendors. It is provided by Kubernetes and matches the
	// `driverName` from the ResourcePool which provides the device.
	//
	// The CEL expression is applied to *all* available devices from any driver.
	// Because of the defaults, it is safe to reference and use attribute values
	// without checking first whether they are set. For this to work without
	// ambiguity, attribute names have to be fully-qualified.
	//
	// Some examples:
	//    device.quantityAttributes["memory.dra.example.com"].isGreaterThan(quantity("1Gi")) # >= 1Gi
	//    "memory.dra.example.com" in device.attributes # attribute is set
	//
	// +optional
	DeviceSelector *string `json:"deviceSelector,omitempty" protobuf:"bytes,1,opt,name=deviceSelector"`

	// TODO for 1.31: define how to request a "partioned device"
}

// Namespace scoped.

// ResourceClaim describes which resources (typically one or more devices)
// are needed by a claim consumer.
// Its status tracks whether the claim has been allocated and what the
// resulting attributes are.
//
// This is an alpha type and requires enabling the DynamicResourceAllocation
// feature gate.
type ResourceClaim struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines what to allocated and how to configure it.
	Spec ResourceClaimSpec `json:"spec"`

	// Status describes whether the claim is ready for use.
	// +optional
	Status ResourceClaimStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type ResourceClaimSpec struct {
	// This field holds configuration for multiple potential drivers which
	// could satisfy requests in this claim. The configuration applies to
	// the entire claim. It is ignored while allocating the claim.
	//
	// +optional
	// +listType=atomic
	Config []ConfigurationParameters `json:"config,omitempty" protobuf:"bytes,4,opt,name=config"`

	// These constraints must be satisfied by the set of devices that get
	// allocated for the claim.
	Constraints []Constraint `json:"constraints"`

	// Requests are individual requests for separate resources for the claim.
	// An empty list is valid and means that the claim can always be allocated
	// without needing anything. A class can be referenced to use the default
	// requests from that class.
	//
	// +listType=atomic
	Requests []Request `json:"requests,omitempty" protobuf:"bytes,5,name=requests"`

	// Future extension, ignored by older schedulers. This is fine because scoring
	// allows users to define a preference, without making it a hard requirement.
	//
	//
	// Score *SomeScoringStruct

	// Shareable indicates whether the allocated claim is meant to be shareable
	// by multiple consumers at the same time.
	// +optional
	Shareable bool `json:"shareable,omitempty" protobuf:"bytes,3,opt,name=shareable"`
}

// Constraint must have one and only one field set.
type Constraint struct {
	// All devices must have this attribute and its value must be the same.
	//
	// For example, if you specified "numa.dra.example.com" (a hypothetical example!),
	// then only devices in the same NUMA node will be chosen.
	//
	// +optional
	// +listType=atomic
	MatchAttribute *string `json:"matchAttribute,omitempty"`

	// Future extension, not part of the current design:
	// A CEL expression which compares different devices and returns
	// true if they match.
	//
	// Because it would be part of a one-of, old schedulers will not
	// accidentally ignore this additional, for them unknown match
	// criteria.
	//
	// matcher string
}

// Request is a request for one of many resources required for a claim.
// This is typically a request for a single resource like a device, but can
// also ask for one of several different alternatives.
type Request struct {
	// Name can be set to enable referencing this request in a pod.spec.containers[].resources.devices
	// entry, if that is desired.
	//
	// Must be a DNS label.
	//
	// +optional
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	*ResourceRequestDetail `json:",inline,omitempty"`

	// FUTURE EXTENSION:
	//
	// OneOf contains a list of requests, only one of which must be satisfied.
	// Requests are listed in order of priority.
	//
	// +optional
	// +listType=atomic
	// OneOf []ResourceRequestDetail `json:"oneOf,omitempty"` // candidate for a separate KEP in 1.32, not required for 1.31
}

type ResourceRequestDetail struct {
	// When referencing a DeviceClass, a request inherits additional
	// configuration and requirements.
	//
	// +optional
	DeviceClassName *string `json:"deviceClassName,omitempty"`

	// Config defines configuration parameters that apply to the requested resource(s).
	// They are ignored while allocating the claim.
	//
	// +optional
	// +listType=atomic
	Config []ConfigurationParameters `json:"config,omitempty" protobuf:"bytes,1,opt,name=config"`

	// AdminAccess indicates that this is a claim for administrative access
	// to the device(s). Claims with AdminAccess are expected to be used for
	// monitoring or other management services for a device.  They ignore
	// all ordinary claims to the device with respect to access modes and
	// any resource allocations. Ability to request this kind of access is
	// controlled via ResourceQuota in the resource.k8s.io API.
	//
	// Can be combined with a range to ask for access to all devices
	// on a node which match the requrirements.
	//
	// Default is false.
	//
	// +optional
	AdminAccess *bool `json:"adminAccess,omitempty"`

	// FUTURE EXTENSION:
	// - Count with min, max (including "all that are available").
	// - Constraints for the set of devices belonging to this request.

	// Requirements describe additional contraints that all must be met by a device
	// to satisfy the request.
	//
	// +optional
	// +listType=atomic
	Requirements []Requirement `json:"requirements,omitempty" protobuf:"bytes,4,opt,name=requirements"`
}

// ResourceClaimStatus tracks whether the resource has been allocated and what
// the result of that was.
type ResourceClaimStatus struct {
	// Allocation is set once the claim has been allocated successfully.
	// +optional
	Allocation *AllocationResult `json:"allocation,omitempty" protobuf:"bytes,2,opt,name=allocation"`

	// ReservedFor indicates which entities are currently allowed to use
	// the claim. A Pod which references a ResourceClaim which is not
	// reserved for that Pod will not be started.
	//
	// There can be at most 32 such reservations. This may get increased in
	// the future, but not reduced.
	//
	// +listType=map
	// +listMapKey=uid
	// +patchStrategy=merge
	// +patchMergeKey=uid
	// +optional
	ReservedFor []ResourceClaimConsumerReference `json:"reservedFor,omitempty" protobuf:"bytes,3,opt,name=reservedFor" patchStrategy:"merge" patchMergeKey:"uid"`
}

// AllocationResult contains attributes of an allocated resource.
type AllocationResult struct {
	// DriverData contains the state associated with an allocation that
	// should be maintained throughout the lifetime of a claim. Each
	// entry contains data that should be passed to a specific kubelet
	// plugin once the claim lands on a node.
	//
	// Setting this field is optional. It has a maximum size of 32 entries.
	// If empty, nothing was allocated for the claim and kubelet does not
	// need to prepare anything for it.
	//
	// +listType=atomic
	// +optional
	DriverData []DriverData `json:"driverData,omitempty" protobuf:"bytes,1,opt,name=driverData"`

	// This field defines where Pods can be scheduled which reference
	// an allocated claim.
	//
	// Setting this field is optional. If null, the resource is available
	// everywhere.
	// +optional
	AvailableOnNodes *v1.NodeSelector `json:"availableOnNodes,omitempty" protobuf:"bytes,2,opt,name=availableOnNodes"`

	// Shareable determines whether the resource supports more
	// than one consumer at a time.
	// +optional
	Shareable bool `json:"shareable,omitempty" protobuf:"varint,3,opt,name=shareable"`
}

// DriverData holds information for processing by a specific kubelet plugin.
type DriverData struct {
	// DriverName specifies the name of the resource driver whose kubelet
	// plugin should be invoked to process the allocation once the claim is
	// needed on a node.
	//
	// Must be a DNS subdomain and should end with a DNS domain owned by the
	// vendor of the driver.
	DriverName string `json:"driverName" protobuf:"bytes,1,name=driverName"`

	// Data contains all information about the allocation that the kubelet
	// plugin will need.
	//
	// +optional
	Data *StructuredDriverData `json:"data,omitempty" protobuf:"bytes,2,opt,name=data"`

	// Alternative to "Data", not describe further here for the sake of brevity:
	// OpaqueData *string // Set by control plane controller when using "classic DRA".
}

// StructuredDriverData is the in-tree representation of the allocation result.
type StructuredDriverData struct {
	// Config contains all the configuration pieces that apply to the entire claim
	// and that were meant for the driver which handles these resources.
	// They get collected during the allocation and stored here
	// to ensure that they remain available while the claim is allocated.
	//
	// Entries are sorted from "most specific" first to "least specific" last:
	// - claim
	// - claim class reference
	//
	// +optional
	Config []DriverConfiguration `json:"config,omitempty"`

	// NodeName is the name of the node providing the necessary resources
	// if the resources are local to a node.
	//
	// +optional
	NodeName string `json:"nodeName,omitempty" protobuf:"bytes,3,name=nodeName"`

	// Results lists all allocated driver resources.
	//
	// +listType=atomic
	Results []RequestAllocationResult `json:"results" protobuf:"bytes,4,name=results"`
}

// DriverConfiguration is one entry in a list of configuration pieces.
type DriverConfiguration struct {
	// Admins is true if the source of the piece was a class and thus
	// not something that a normal user would have been able to set.
	Admin bool `json:"admin,omnitempty"`

	DriverConfigurationAlternatives `json:",inline"`
}

// DriverConfigurationAlternatives must have one and only one one field set.
//
// In contrast to VendorConfigurationParameters, the driver name is
// not included and has to be infered from the context.
type DriverConfigurationAlternatives struct {
	Opaque *runtime.RawExtension `json:"opaque,omitempty" protobuf:"bytes,1,opt,name=opaque"`
}

// RequestAllocationResult contains configuration and the allocation result for
// one request.
type RequestAllocationResult struct {
	// Config contains all the configuration pieces that apply to the request
	// and that were meant for the driver which handles these resources.
	// They get collected during the allocation and stored here
	// to ensure that they remain available while the claim is allocated.
	//
	// Entries are sorted from "most specific" first to "least specific" last:
	// - claim request
	// - claim request class reference
	//
	// +optional
	Config []DriverConfiguration `json:"config,omitempty"`

	// RequestName identifies the request in the claim which caused this
	// resource to be allocated.
	//
	// +optional
	RequestName string `json:"requestName,omitempty"`

	AllocationResultModel `json:",inline" protobuf:"bytes,2,name=allocationResultModel"`
}

// AllocationResultModel must have one and only one field set.
type AllocationResultModel struct {
	// Device references one device instance.
	//
	// +optional
	Device *NamedDeviceAllocationResult `json:"device,omitempty" protobuf:"bytes,1,opt,name=device"`
}

// NamedDeviceAllocationResult is used in AllocationResultModel.
type NamedDeviceAllocationResult struct {
	// Name is the name of the selected device instance.
	Name string `json:"name" protobuf:"bytes,1,name=name"`
}

// ResourceClaimConsumerReference contains enough information to let you
// locate the consumer of a ResourceClaim. The user must be a resource in the same
// namespace as the ResourceClaim.
type ResourceClaimConsumerReference struct {
	// APIGroup is the group for the resource being referenced. It is
	// empty for the core API. This matches the group in the APIVersion
	// that is used when creating the resources.
	// +optional
	APIGroup string `json:"apiGroup,omitempty" protobuf:"bytes,1,opt,name=apiGroup"`
	// Resource is the type of resource being referenced, for example "pods".
	Resource string `json:"resource" protobuf:"bytes,3,name=resource"`
	// Name is the name of resource being referenced.
	Name string `json:"name" protobuf:"bytes,4,name=name"`
	// UID identifies exactly one incarnation of the resource.
	UID types.UID `json:"uid" protobuf:"bytes,5,name=uid"`
}

// ResourceClaimSpecReference contains enough information to let you
// locate the specification for a ResourceClaim. The object must be in the same
// namespace as the ResourceClaim.
type ResourceClaimSpecReference struct {
	// APIGroup is the group for the resource being referenced. It is
	// empty for the core API. This matches the group in the APIVersion
	// that is used when creating the resources.
	// +optional
	APIGroup string `json:"apiGroup,omitempty" protobuf:"bytes,1,opt,name=apiGroup"`
	// Kind is the type of resource being referenced. This is the same
	// value as in the parameter object's metadata, for example "ConfigMap".
	Kind string `json:"kind" protobuf:"bytes,2,name=kind"`
	// Name is the name of resource being referenced.
	Name string `json:"name" protobuf:"bytes,3,name=name"`
}

// ResourceClaimTemplate is used to produce ResourceClaim objects.
type ResourceClaimTemplate struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Describes the ResourceClaim that is to be generated.
	//
	// This field is immutable. A ResourceClaim will get created by the
	// control plane for a Pod when needed and then not get updated
	// anymore.
	Spec ResourceClaimTemplateSpec `json:"spec" protobuf:"bytes,2,name=spec"`
}

// ResourceClaimTemplateSpec contains the metadata and fields for a ResourceClaim.
type ResourceClaimTemplateSpec struct {
	// ObjectMeta may contain labels and annotations that will be copied into the PVC
	// when creating it. No other fields are allowed and will be rejected during
	// validation.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec for the ResourceClaim. The entire content is copied unchanged
	// into the ResourceClaim that gets created from this template. The
	// same fields as in a ResourceClaim are also valid here.
	Spec ResourceClaimSpec `json:"spec" protobuf:"bytes,2,name=spec"`
}
