package api

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// DeviceClass is a vendor or admin-provided resource that contains
// contraint and configuration information. Essentially, it is a re-usable
// collection of predefined data that device claims may use.
// Cluster scoped.
type ResourceClass struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// ControllerName defines the name of the dynamic resource driver that is
	// used for allocation of a ResourceClaim that uses this class. If empty,
	// structured parameters are used for allocating claims using this class.
	//
	// Resource drivers have a unique name in forward domain order
	// (acme.example.com).
	//
	// This is an alpha field and requires enabling the DRAControlPlaneController
	// feature gate.
	//
	// +optional
	ControllerName string `json:"controllername,omitempty"`

	// Only nodes matching the selector will be considered by the scheduler
	// when trying to find a Node that fits a Pod when that Pod uses
	// a ResourceClaim that has not been allocated yet.
	//
	// Setting this field is optional. If null, all nodes are candidates.
	// +optional
	SuitableNodes *v1.NodeSelector `json:"suitableNodes,omitempty" protobuf:"bytes,4,opt,name=suitableNodes"`

	// ClaimConfig defines configuration parameters that apply to each claim using this class.
	// They are ignored while allocating the claim.
	//
	// +optional
	ClaimConfig *ConfigurationParameters `json:"claimConfig,omitempty" protobuf:"bytes,3,opt,name=config"`

	// RequestConfig defines configuration parameters that apply to each request in a claim using this class.
	// They are ignored while allocating the claim.
	//
	// +optional
	RequestConfig *ConfigurationParameters `json:"requestConfig,omitempty" protobuf:"bytes,3,opt,name=config"`

	// Filters describes additional contraints that must be met when using the class.
	//
	// +optional
	Filter *ResourceFilterModel `json:"filter,omitempty" protobuf:"bytes,4,opt,name=filter"`

	// DefaultRequests are individual requests for separate resources for a
	// claim using this class. In contrast to config and filter, these
	// requests are only used if the claim does not specify its own list of
	// requests.
	//
	// +listType=atomic
	DefaultRequests []ResourceRequest `json:"defaultRequests" protobuf:"bytes,5,name=requests"`
}

// ConfigurationParameters must have one and only one field set.
type ConfigurationParameters struct {
	// +listType=atomic
	Vendor []VendorConfigurationParameters `json:"vendor,omitempty" protobuf:"bytes,1,opt,name=vendor"`
}

// VendorConfigurationParameters contains configuration parameters for a driver.
type VendorConfigurationParameters struct {
	// DriverName is used to determine which kubelet plugin needs
	// to be passed these configuration parameters.
	//
	// An admission webhook provided by the vendor could use this
	// to decide whether it needs to validate them.
	DriverName string `json:"driverName,omitempty" protobuf:"bytes,1,opt,name=driverName"`

	// Parameters can contain arbitrary data. It is the
	// responsibility of the vendor to handle validation and
	// versioning.
	Parameters runtime.RawExtension `json:"parameters,omitempty" protobuf:"bytes,2,opt,name=parameters"`
}

// ResourceFilterModel must have one and only one field set.
type ResourceFilterModel struct {
	// Devices describes a filter based on device attributes.
	//
	// +optional
	Devices *DeviceFilter `json:"devices,omitempty"`
}

type DeviceFilter struct {
	// DriverName, if set, excludes any device not provided by this driver.
	//
	// +optional
	DriverName string `json:"driverName,omitempty" protobuf:"bytes,1,opt,name=driverName"`

	// Selector is a CEL expression which must evaluate to true if a
	// resource instance is suitable. The language is as defined in
	// https://kubernetes.io/docs/reference/using-api/cel/
	//
	// In addition, for each type in NamedResourcesAttributeValue there is a map that
	// resolves to the corresponding value of the instance under evaluation. Unknown
	// names cause a runtime error. Note that the CEL expression is applied to
	// *all* available resource instances by default, regardless of which driver provides it.
	// In that case. the CEL expression must first check that the instance has certain
	// attributes before using them.
	//
	// For example:
	//    attributes.quantity.has("a.dra.example.com") &&
	//    attributes.quantity["a.dra.example.com"].isGreaterThan(quantity("0")) &&
	//    # No separate check, b.dra.example.com is set whenever a.dra.example.com is,
	//    attributes.stringslice["b.dra.example.com"].isSorted()
	//
	// If a driver name is set, then such a check is not be needed if all instances
	// are known to have the attribute. Attributes names don't have to have
	// the driver name suffix.
	//
	// For example:
	//    attributes.quantity["a"].isGreaterThan(quantity("0")) &&
	//    attributes.stringslice["b"].isSorted()
	//
	// If empty, the selector matches any device.
	//
	// +optional
	Selector string `json:"selector" protobuf:"bytes,2,name=selector"`
}

// Namespace scoped.

// ResourceClaim describes which resources (typically one or more devices)
// are needed by a resource consumer.
// Its status tracks whether the resource has been allocated and what the
// resulting attributes are.
//
// This is an alpha type and requires enabling the DynamicResourceAllocation
// feature gate.
type ResourceClaim struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	ResourceClaimSpecAlternatives `json:",inline"` // Inlined with "spec" defined in the ResourceClaimSpec.

	// Status describes whether the claim is ready for use.
	// +optional
	Status ResourceClaimStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// ResourceClaimSpecAlternatives defines how a claim is to be allocated.
type ResourceClaimSpecAlternatives struct {
	// Spec define what to allocated and how to configure it.
	// Spec and SpecRef are mutually exclusive.
	// +optional
	Spec *ResourceClaimSpec `json:"spec"`

	// SpecRef references a separate object with the specification of the claim
	// that will be used by the driver when allocating a resource for the
	// claim. Parameters and ParametersRef are mutually exclusive.
	//
	// SpecRef is typically used to reference a vendor CR.
	// A vendor controller then converts that CR into a ResourceClaimSpecification
	// object and that is then used to allocate the claim.
	//
	// The object must be in the same namespace as the ResourceClaim.
	// +optional
	SpecRef *ResourceClaimSpecReference `json:"specRef,omitempty"`
}

// Used inside a ResourceClaimSpecAlternatives or a ResourceClaimSpecification object.
type ResourceClaimSpec struct {
	// ResourceClassName references additional configuration and filters
	// that apply to the whole claim and all requests in it. If the class
	// contains default requests, then those are used if (and only if)
	// the claim does not provide those itself.
	//
	// Filters in the class must match in addition to the filters in the claim
	// parameters.
	//
	// +optional
	ResourceClassName string `json:"resourceClassName,omitempty" protobuf:"bytes,1,name=resourceClassName"`

	// To be decided: does it make sense to allow referencing multiple
	// classes?  As it stands now, there could be a need to create a large
	// variety of different classes where each class is one combination of
	// different options.
	//
	// When allowing different classes, each class could describe one aspect:
	// - "devices from vendor foo"
	// - "more than 10 GiB of RAM"
	//
	// Probably need better use cases for classes. Stuff for a new KEP in >= 1.32...

	// Config defines configuration parameters that apply to the entire claim.
	// They are ignored while allocating the claim.
	//
	// +optional
	Config *ConfigurationParameters `json:"config,omitempty" protobuf:"bytes,4,opt,name=config"`

	// Requests are individual requests for separate resources for the claim.
	// An empty list is valid and means that the claim can always be allocated
	// without needing anything. A class can be referenced to use the default
	// requests from that class.
	//
	// +listType=atomic
	Requests []ResourceRequest `json:"requests,omitempty" protobuf:"bytes,5,name=requests"`

	// MatchAttributes allows specifying a constraint that will apply
	// across all of devices that need to be allocated for this claim.
	//
	// For example, if you specified "numa.dra.example.com" (a hypothetical example!),
	// then only devices which have that attribute with the same value will
	// be considered.
	//
	// +optional
	// +listType=atomic
	MatchAttributes []string `json:"matchAttributes,omitempty"`

	// Shareable indicates whether the allocated claim is meant to be shareable
	// by multiple consumers at the same time.
	// +optional
	Shareable bool `json:"shareable,omitempty" protobuf:"bytes,3,opt,name=shareable"`
}

// ResourceRequest is a request for one of many resources required for a claim.
// This is typically a request for a single resource like a device, but can
// also ask for one of several different alternatives.
type ResourceRequest struct {
	// Name can be set to enable referencing this request in a pod.spec.containers[].resources.devices
	// entry, if that is desired.
	//
	// +optional
	Name string

	*ResourceRequestDetail `json:",inline,omitempty"`

	// OneOf contains a list of requests, only one of which must be satisfied.
	// Requests are listed in order of priority.
	//
	// +optional
	// +listType=atomic
	OneOf []ResourceRequestDetail `json:"oneOf,omitempty"` // candidate for a separate KEP in 1.32, not required for 1.31
}

type ResourceRequestDetail struct {
	// ResourceClassName references additional configuration and filters that apply
	// to the request.
	//
	// Filters in the class must match in addition to the filters in the claim
	// parameters.
	ResourceClassName string `json:"resourceClassName" protobuf:"bytes,1,name=resourceClassName"`

	// Config defines configuration parameters that apply to the requested resource(s).
	// They are ignored while allocating the claim.
	Config *ConfigurationParameters `json:"config,omitempty" protobuf:"bytes,1,opt,name=config"`

	// AdminAccess indicates that this is a claim for administrative access
	// to the device(s). Claims with AdminAccess are expected to be used for
	// monitoring or other management services for a device.  They ignore
	// all ordinary claims to the device with respect to access modes and
	// any resource allocations. Ability to request this kind of access is
	// controlled via ResourceQuota in the resource.k8s.io API.
	//
	// Can be combined with a range to ask for access to all devices
	// on a node which match the filter.
	//
	// Default is false.
	//
	// +optional
	AdminAccess *bool `json:"adminAccess,omitempty"`

	// Count defines how many instances are desired. If unset, exactly one
	// instance must be available. When a range is set, it is possible to
	// ask for:
	// - x >= minimum instances up to all that are available
	// - 0 <= x <= maximum (up to a certain number, with zero instances acceptable)
	// - minimum <= 0 <= maximum (within a certain range)
	// +optional
	Count *IntRange `json:"count,omitempty"`

	ResourceRequestModel `json:",inline" protobuf:"bytes,2,name=resourceRequestModel"`
}

// IntRange defines how many instances are desired.
type IntRange struct {
	// Minimum defines the lower limit. At least this many instances
	// must be available (x >= minimum). The default if unset is one.
	Minimum *int `json:"miminum"`

	// Maximum defines the upper limit. At most this many instances
	// may be allocated (x <= maximum). The default if unset is unlimited.
	Maximum *int `json:"maximum"`
}

// ResourceRequestModel must have one and only one field set.
type ResourceRequestModel struct {
	// Device describes a request for a specific device.
	//
	// +optional
	Device *DeviceRequest `json:"device,omitempty"`
}

// DeviceRequest is used in ResourceRequestModel.
type DeviceRequest struct {
	// DriverName excludes any named resource not provided by this driver.
	//
	// +optional
	DriverName *string `json:"driverName,omitempty" protobuf:"bytes,1,opt,name=driverName"`

	// Selector is a CEL expression which must evaluate to true if a
	// resource instance is suitable. The language is as defined in
	// https://kubernetes.io/docs/reference/using-api/cel/
	//
	// In addition, for each type in NamedResourcesAttributeValue there is a map that
	// resolves to the corresponding value of the instance under evaluation. Unknown
	// names cause a runtime error. Note that the CEL expression is applied to
	// *all* available resource instances by default, regardless of which driver provides it.
	// In that case. the CEL expression must first check that the instance has certain
	// attributes before using them.
	//
	// For example:
	//    attributes.quantity.has("a.dra.example.com") &&
	//    attributes.quantity["a.dra.example.com"].isGreaterThan(quantity("0")) &&
	//    # No separate check, b.dra.example.com is set whenever a.dra.example.com is,
	//    attributes.stringslice["b.dra.example.com"].isSorted()
	//
	// If a driver name is set, then such a check is not be needed if all instances
	// are known to have the attribute. Attributes names don't have to have
	// the driver name suffix.
	//
	// For example:
	//    attributes.quantity["a"].isGreaterThan(quantity("0")) &&
	//    attributes.stringslice["b"].isSorted()
	//
	// If empty, any device matches.
	//
	// +optional
	Selector string `json:"selector" protobuf:"bytes,2,name=selector"`
}

// ResourceClaimStatus tracks whether the resource has been allocated and what
// the result of that was.
type ResourceClaimStatus struct {
	// ControllerName is a copy of the driver name from the ResourceClass at
	// the time when allocation started. It is empty when the claim was
	// allocated through structured parameters,
	//
	// This is an alpha field and requires enabling the DRAControlPlaneController
	// feature gate.
	//
	// +optional
	ControllerName string `json:"controllerName"`

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

	// DeallocationRequested indicates that a ResourceClaim is to be
	// deallocated.
	//
	// The driver then must deallocate this claim and reset the field
	// together with clearing the Allocation field.
	//
	// While DeallocationRequested is set, no new consumers may be added to
	// ReservedFor.
	//
	// This is an alpha field and requires enabling the DRAControlPlaneController
	// feature gate.
	//
	// +optional
	DeallocationRequested bool `json:"deallocationRequested,omitempty" protobuf:"varint,4,opt,name=deallocationRequested"`
}

// AllocationResult contains attributes of an allocated resource.
type AllocationResult struct {
	// ResourceHandles contain the state associated with an allocation that
	// should be maintained throughout the lifetime of a claim. Each
	// ResourceHandle contains data that should be passed to a specific kubelet
	// plugin once it lands on a node. This data is returned by the driver
	// after a successful allocation and is opaque to Kubernetes. Driver
	// documentation may explain to users how to interpret this data if needed.
	//
	// Setting this field is optional. It has a maximum size of 32 entries.
	// If null (or empty), it is assumed this allocation will be processed by a
	// single kubelet plugin with no ResourceHandle data attached. The name of
	// the kubelet plugin invoked will match the DriverName set in the
	// ResourceClaimStatus this AllocationResult is embedded in.
	//
	// +listType=atomic
	// +optional
	ResourceHandles []ResourceHandle `json:"resourceHandles,omitempty" protobuf:"bytes,1,opt,name=resourceHandles"`

	// This field will get set by the resource driver after it has allocated
	// the resource to inform the scheduler where it can schedule Pods using
	// the ResourceClaim.
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

// ResourceHandle holds opaque resource data for processing by a specific kubelet plugin.
type ResourceHandle struct {
	// DriverName specifies the name of the resource driver whose kubelet
	// plugin should be invoked to process this ResourceHandle's data once it
	// lands on a node.
	DriverName string `json:"driverName" protobuf:"bytes,1,name=driverName"`

	// Data contains the opaque data associated with this ResourceHandle. It is
	// set by the controller component of the resource driver whose name
	// matches the DriverName set in the ResourceClaimStatus this
	// ResourceHandle is embedded in. It is set at allocation time and is
	// intended for processing by the kubelet plugin whose name matches
	// the DriverName set in this ResourceHandle.
	//
	// The maximum size of this field is 16KiB. This may get increased in the
	// future, but not reduced.
	//
	// This is an alpha field and requires enabling the DRAControlPlaneController feature gate.
	//
	// +optional
	Data string `json:"data,omitempty" protobuf:"bytes,2,opt,name=data"`

	// If StructuredData is set, then it needs to be used instead of Data.
	//
	// +optional
	StructuredData *StructuredResourceHandle `json:"structuredData,omitempty" protobuf:"bytes,5,opt,name=structuredData"`
}

// ResourceHandleDataMaxSize represents the maximum size of resourceHandle.data.
const ResourceHandleDataMaxSize = 16 * 1024

// StructuredResourceHandle is the in-tree representation of the allocation result.
type StructuredResourceHandle struct {
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
	Vendor *runtime.RawExtension `json:"vendor,omitempty" protobuf:"bytes,1,opt,name=vendor"`
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
	Device *NamedDeviceAllocationResult `json:"namedResources,omitempty" protobuf:"bytes,1,opt,name=namedResources"`
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

// ResourceClaimSpecification contains the specification for a ResourceClaim in an
// in-tree format understood by Kubernetes.
type ResourceClaimSpecification struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// If this object was created from some other resource, then this links
	// back to that resource. This field is used to find the in-tree representation
	// of the claim parameters when the parameter reference of the claim refers
	// to some unknown type.
	// +optional
	GeneratedFrom *ResourceClaimSpecReference `json:"generatedFrom,omitempty" protobuf:"bytes,2,opt,name=generatedFrom"`

	ResourceClaimSpec // inline
}
