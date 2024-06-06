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

	// Only nodes matching the selector will be considered by the scheduler
	// when trying to find a Node that fits a Pod when that Pod uses
	// a claim that has not been allocated yet *and* that claim
	// gets allocated through a control plane controller. It is ignored
	// when the claim does not use a control plane controller
	// for allocation.
	//
	// Setting this field is optional. If unset, all nodes are candidates.
	//
	// This is an alpha field and requires enabling the DRAControlPlaneController
	// feature gate.
	//
	// +optional
	SuitableNodes *v1.NodeSelector `json:"suitableNodes,omitempty" protobuf:"bytes,5,opt,name=suitableNodes"`
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

	// ControllerName defines the name of the DRA driver that is meant
	// to handle allocation of this claim. If empty, allocation is handled
	// by the scheduler while scheduling a pod.
	//
	// Must be a DNS subdomain and should end with a DNS domain owned by the
	// vendor of the driver.
	//
	// This is an alpha field and requires enabling the DRAControlPlaneController
	// feature gate.
	//
	// +optional
	ControllerName string `json:"controllerName,omitempty"`
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
// This is typically a request for a single resource like a device, but may get
// extended to also ask for one of several different alternatives.
type Request struct {
	// The name can be used to reference this request in a pod.spec.containers[].resources.claims
	// entry.
	//
	// Must be a DNS label.
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

// ResourceClaimStatus tracks whether the claim has been allocated and what
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

	// Indicates that a claim is to be deallocated. While this is set,
	// no new consumers may be added to ReservedFor.
	//
	// This is only used if the claim needs to be deallocated by a DRA driver.
	// That driver then must deallocate this claim and reset the field
	// together with clearing the Allocation field.
	//
	// This is an alpha field and requires enabling the DRAControlPlaneController
	// feature gate.
	//
	// +optional
	DeallocationRequested bool `json:"deallocationRequested,omitempty" protobuf:"varint,4,opt,name=deallocationRequested"`
}

// AllocationResult contains attributes of an allocated resource.
type AllocationResult struct {
	// ControllerName is the name of the DRA driver which handled the
	// allocation. That driver is also responsible for deallocating the
	// claim. It is empty when the claim can be deallocated without
	// involving a driver.
	//
	// A driver may allocate devices provided by other drivers, so this
	// driver name here can be different from the driver names listed in
	// DriverData.
	//
	// This is an alpha field and requires enabling the DRAControlPlaneController
	// feature gate.
	//
	// +optional
	ControllerName string `json:"controllerName,omitempty"`

	// This field holds configuration for all drivers which
	// satisfied requests in this claim. The configuration applies to
	// the entire claim.
	//
	// +optional
	// +listType=atomic
	Config []ConfigurationParameters `json:"config,omitempty" protobuf:"bytes,4,opt,name=config"`

	// Results lists all allocated devices.
	//
	// +listType=atomic
	Results []RequestAllocationResult `json:"results" protobuf:"bytes,4,name=results"`

	// Setting this field is optional. If unset, the allocated devices are available everywhere.
	//
	// +optional
	AvailableOnNodes *v1.NodeSelector `json:"availableOnNodes,omitempty" protobuf:"bytes,2,opt,name=availableOnNodes"`

	// Shareable determines whether the claim supports more than one consumer at a time.
	//
	// +optional
	Shareable bool `json:"shareable,omitempty" protobuf:"varint,3,opt,name=shareable"`
}

// DriverData holds information for processing by a specific kubelet plugin.
type DriverData struct {
	// DriverName specifies the name of the DRA driver whose kubelet
	// plugin should be invoked to process the allocation once the claim is
	// needed on a node.
	//
	// Must be a DNS subdomain and should end with a DNS domain owned by the
	// vendor of the driver.
	DriverName string `json:"driverName" protobuf:"bytes,1,name=driverName"`

	// Config contains all the configuration pieces that apply to the entire claim
	// and that were meant for the driver which handles these devices.
	// They get collected during the allocation and stored here
	// to ensure that they remain available while the claim is allocated.
	//
	// Entries are listed in the same order as in claim.config.
	//
	// +optional
	Config []DriverConfigurationParameters `json:"config,omitempty"`
}

// RequestAllocationResult contains configuration and the allocation result for
// one request.
type RequestAllocationResult struct {
	// DriverName specifies the name of the DRA driver whose kubelet
	// plugin should be invoked to process the allocation once the claim is
	// needed on a node.
	//
	// Must be a DNS subdomain and should end with a DNS domain owned by the
	// vendor of the driver.
	DriverName string `json:"driverName" protobuf:"bytes,1,name=driverName"`

	// Config contains all the configuration pieces that apply to the request
	// and that were meant for the driver which handles these devices.
	// They get collected during the allocation and stored here
	// to ensure that they remain available while the claim is allocated.
	//
	// Entries are list in the same order as in class.config and claim.config,
	// with class.config entries first.
	//
	// +optional
	Config []DeviceConfiguration `json:"config,omitempty"`

	// RequestName identifies the request in the claim which caused this
	// device to be allocated.
	RequestName string `json:"requestName"`

	// This node name together with the driver name and
	// the device name field identify which device was allocated.
	NodeName string `json:"nodeName"`

	// DeviceName references one device instance via its name in the driver's
	// resource pool.
	DeviceName string `json:"deviceName"`
}

// DeviceConfiguration is one entry in a list of configuration pieces for a device.
type DeviceConfiguration struct {
	// Admins is true if the source of the piece was a class and thus
	// not something that a normal user would have been able to set.
	Admin bool `json:"admin,omnitempty"`

	DriverConfigurationParameters `json:",inline"`
}

// DriverConfigurationParameters must have one and only one one field set.
//
// In contrast to ConfigurationParameters, the driver name is
// not included and has to be infered from the context.
type DriverConfigurationParameters struct {
	Opaque *runtime.RawExtension `json:"opaque,omitempty" protobuf:"bytes,1,opt,name=opaque"`
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
