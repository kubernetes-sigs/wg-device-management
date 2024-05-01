package api

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeviceClass is a vendor or admin-provided resource that contains
// contraint and configuration information. Essentially, it is a re-usable
// collection of predefined data that device claims may use.
// Cluster scoped.
type DeviceClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DeviceClassSpec `json:"spec,omitempty"`
}

// DeviceClassSpec provides the details of the DeviceClass.
type DeviceClassSpec struct {
	// Driver specifies the driver that should handle this class of devices.
	// When a DeviceClaim uses this class, only devices published by the
	// specified driver will be considered.
	// +required
	Driver string `json:driver,omitempty`

	// DeviceType is a driver-independent classification of the device.  In
	// claims, this may be used instead of specifying the class
	// explicitly, so that we do not aribtrarily limit claims to a
	// particular vendor's devices.
	//
	// Alternatively, we may want to consider a DeviceCapabilities vector,
	// or use device attributes or individual resource types supported by a
	// device to indicate device functions.
	//
	// +required
	DeviceType string `json:deviceType,omitempty`

	// Constraints is a CEL expression that operates on device attributes,
	// and must evaluate to true for a device to be considered. It will be
	// ANDed with any Constraints field in the DeviceClaim using this class.
	//
	// +optional
	Constraints *string `json:"constraints,omitempty"`

	// AdminAccess indicates that this class provides administrative access
	// to the devices. Claims using a class with AdminAccess are expected
	// to be used for monitoring or other management services for a device.
	// They ignore all ordinary claims to the device with respect to access
	// modes and any resource allocations. Access to these classes must be
	// controlled via ResourceQuota. Default is false.
	//
	// +optional
	AdminAccess *bool `json:"adminAccess,omitempty"`

	// DeviceConfigs contains references to arbitrary vendor device configuration
	// objects that will be attached to the device allocation.
	//
	// +optional
	Configs []DeviceClassConfigReference `json:"configs,omitempty"`
}

// DeviceClaim is used to specify a request for a set of devices.
// Namespace scoped.
type DeviceClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeviceClaimSpec   `json:"spec,omitempty"`
	Status DeviceClaimStatus `json:"status,omitempty"`
}

// DeviceClaimSpec details the requirements that devices chosen
// to satisfy this claim must meet.
type DeviceClaimSpec struct {
	// MatchAttributes allows specifying a constraint that will apply
	// across all of the claims. For example, if you specified "numa", then
	// this overall claim could only be successfully fulfilled if all of
	// the included claims could be fulfilled by pools with the same "numa"
	// attribute value. If we simply set matchAttributes in each claim
	// separately, then they could be consistent within claims, but
	// inconsistent across claims. Therefore, we need this additional
	// constraint.
	//
	// +optional
	MatchAttributes []string `json:"matchAttributes,omitempty"`

	// Claims contains the actual claim details, arranged into groups
	// containing claims which must all be satsified, or for which only
	// one needs to be satisfied.
	//
	// +required
	Claims []DeviceClaimInstance `json:"claims,omitempty"`
}

// DeviceClaimInstance captures a claim which must be satisfied,
// or a group for which one must be sastisfied.
type DeviceClaimInstance struct {
	// At least one of AllOf and OneOf must be populated.

	// If fields of DeviceClaimDetail are populated, OneOf should
	// be empty.
	DeviceClaimDetail `json:",inline"`

	// OneOf contains a list of claims, only one of which must be satisfied.
	// Claims are listed in order of priority.
	//
	// +optional
	OneOf []DeviceClaimDetail `json:"oneOf,omitempty"`
}

// DeviceClaimDetail contains the details of how to fulfill a specific
// request for devices.
type DeviceClaimDetail struct {
	// DeviceType may be specified to get a device from any class that
	// supports this type of device. For example, the user can request
	// an 'sriov-nic', and any class that can provide that type will be
	// considered for fulfillment of the claim.
	//
	// +optional
	DeviceType *string `json:"deviceType"`

	// DeviceClass is the name of the DeviceClass containing the basic information
	// about the device being requested.
	//
	// +optional
	DeviceClass *string `json:"deviceClass"`

	// Constraints is a CEL expression that operates on device attributes.
	// In order for a device to be considered, this CEL expression and the
	// Constraints expression from the DeviceClass must both be true.
	//
	// +optional
	Constraints *string `json:"constraints,omitempty"`

	// Device claims may be satisfied by choosing multiple devices instead
	// of just a single device. How and when that happens depends on the
	// requests and limits specified.

	// Requests allows the user to specify the minimum requirements that
	// must be satisfied across all devices allocated for this claim.  All
	// drivers can support "count" for requests, but other resource types
	// are driver-specific. The default value for requests is a single
	// entry for "count" with value of 1.
	Requests map[string]resource.Quantity `json:"requests,omitempty"`

	// Limits allows the user to control the maximum count of devices
	// that is allocated to satisfy the claim. Depending on the driver
	// and device other resource limits may or may not be enforceable.
	Limits map[string]resource.Quantity `json:"limits,omitempty"`

	// MatchAttributes allows specifying a constraint within a set of
	// chosen devices, without having to explicitly specify the value of
	// the constraint.  For example, this allows constraints like "all
	// devices must be the same model", without having to specify the exact
	// model. We may be able to use this for some basic topology
	// constraints too, by representing the topology as device attributes.
	//
	// Currently, these are just strings. However, we could make them
	// structs, and include required vs preferred matches. Required matches
	// would fail if not met, where as preferred would lower the score if
	// not met. We could even allow low/medium/high priority and adjust the
	// score differently for each.
	//
	// +optional
	MatchAttributes []string `json:"matchAttributes,omitempty"`

	// Configs contains references to arbitrary vendor device configuration
	// objects that will be attached to the device allocation.
	// +optional
	Configs []DeviceConfigReference `json:"configs,omitempty"`
}

// DeviceClaimStatus contains the results of the claim allocation.
type DeviceClaimStatus struct {
	// ClassConfigs contains the entire set of dereferenced vendor
	// configurations from the DeviceClass, as of the time of allocation.
	// +optional
	ClassConfigs []runtime.RawExtension `json:"classConfigs,omitempty"`

	// ClaimConfigs contains the entire set of dereferenced vendor
	// configurations from the DeviceClaim, as of the time of allocation.
	// +optional
	ClaimConfigs []runtime.RawExtension `json:"claimConfigs,omitempty"`

	// Allocations contains the list of device allocations needed to
	// satisfy the claim, one per pool from which devices were allocated.
	//
	// Note that the "current capacity" of the cluster is the result of
	// applying all such allocations to the published DevicePools. This
	// means storing these allocations only in claim status fields is likely
	// to scale poorly, and we will need a different strategy in the real
	// code. For example, we may need to accumulate these in the DevicePool
	// status fields themselves, and just reference them from here.
	//
	// This field is owned by the scheduler, whereas the Devices field
	// is owned by the driver.
	//
	// +optional
	Allocations []DeviceAllocation `json:"allocations,omitempty"`

	// Devices contains the status of each device assigned to this
	// claim, as reported by the driver. This can include driver-specific
	// information. Entries are owned by their respective drivers.
	// TODO: How can we do that?
	DeviceStatuses []DeviceStatus `json:"deviceStatuses,omitempty"`

	// PodNames contains the names of all Pods using this claim.
	// TODO: Can we just use ownerRefs instead?
	// +optional
	PodNames []string `json:"podNames,omitempty"`
}

// NOTE: We no longer have DeviceClaimTemplate. Instead, the PodSpec will
// directly contain a either a DeviceClaimName (to enable multiple pods to
// refer to a pre-provisioned claim), or an embedded struct that includes
// ObjectMeta and a list of the *unrolled fields* of DeviceClaimSpec.  The
// DeviceClaimSpec type itself will not be embedded, but instead its fields
// duplicated, allowing them to evolve independently (and have independent
// validation and avoid Go cyclical dependencies).
//
// NOTE: Feedback on this plan has been negative; the complexity of claims may
// be unmanageble for ordinary users. We may want to be able to embed that
// complexity in classes. We may need some namespaced version of classes, which
// may be a template?

// DeviceClassConfigReference is used to refer to arbitrary configuration
// objects from the class. Since it is the class, and therefore is created by
// the administrator, it allows referencing objects in any namespace.
type DeviceClassConfigReference struct {
	// API version of the referent.
	// +required
	APIVersion string `json:"apiVersion"`

	// Kind of the referent.
	// +required
	Kind string `json:"kind"`

	// Namespace of the referent.
	// +required
	Namespace string `json:"namespace"`

	// Name of the referent.
	// +required
	Name string `json:"name"`
}

// DeviceConfigReference is used to refer to arbitrary configuration object
// from the claim. Since it is created by the end user, the referenced objects
// are restricted to the same namespace as the DeviceClaim.
type DeviceConfigReference struct {
	// API version of the referent.
	// +required
	APIVersion string `json:"apiVersion"`

	// Kind of the referent.
	// +required
	Kind string `json:"kind"`

	// Name of the referent.
	// +required
	Name string `json:"name"`
}

// DeviceAllocation contains an individual device allocation result, including
// per-device resource allocations, when applicable.
type DeviceAllocation struct {
	// DevicePoolName is the name of the DevicePool to which this
	// device belongs.
	// +required
	DevicePoolName string `json:"devicePoolName"`

	// DeviceName contains the name of the allocated Device.
	// +required
	DeviceName string `json:"deviceName,omitempty"`

	// Allocations contain the resource allocations from this device,
	// for the claim. Note that this may only satisfy part of the claim.
	// Also, because devices may allocate some resources in blocks, this
	// may even be larger than the requests or limits in the claim.
	// +optional
	Allocations []ResourceAllocation `json:"allocations,omitempty"`
}

// ResourceAllocation contains the per-device resource allocations.
type ResourceAllocation struct {
	// Name is the resource name/string for this allocation.
	// +required
	Name string `json:"name"`

	// Amount is the amount of resource allocated for this claim.
	// +required
	Allocation resource.Quantity `json:"allocation,inline"`

	// If we ever need to support intra-device topology (for example,
	// if standard node compute becomes a device), then we may also
	// include topology assignments in here.
}

type DeviceIP struct {
	// IP is the IP address assigned to the device
	IP string `json:"ip,omitempty"`
}

// DeviceStatus contains the status of an allocated result, if the driver
// chooses to report it. This may include driver-specific information.
type DeviceStatus struct {
	// DevicePoolName is the name of the DevicePool to which this
	// device belongs. The driver for that device pool owns this
	// entry.
	// +required
	DevicePoolName string `json:"devicePoolName"`

	// DeviceName contains the name of the allocated Device.
	// +required
	DeviceName string `json:"deviceName,omitempty"`

	// Conditions contains the latest observation of the device's state.
	Conditions []metav1.Condition `json:"conditions"`

	// DeviceIP contains the IP allocated for the device, if appropriate.
	// +optional
	DeviceIP *string `json:"deviceIP,omitempty"`

	// DeviceIPs contains all of the IPs allocated for a device, if any.
	// If populated, the zero'th entry must match DeviceIP.
	// +optional
	DeviceIPs []DeviceIP `json:"deviceIPs,omitempty"`

	// Arbitrary driver-specific data.
	// +optional
	DeviceInfo []runtime.RawExtension `json:"deviceInfo,omitempty"`

	// Allocations contain the resource allocations from this device,
	// for the claim. Note that this may only satisfy part of the claim.
	// Also, because devices may allocate some resources in blocks, this
	// may even be larger than the requests or limits in the claim.
	// +optional
	Allocations []ResourceAllocation `json:"allocations,omitempty"`
}
