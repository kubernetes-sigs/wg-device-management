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

// DeviceClassSpec provides either a selector to find other
// classes, or the specific constraints defining this class.
type DeviceClassSpec struct {

	// One of selector and classCriteria must be populated.

	// Selector specifies a label selector used to identify classes which
	// may be considered equivalent to this class. In other words, this
	// allows the adminstrator to define groups of classes which may be
	// referred to as a single class in the claim.
	//
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	// ClassCriteria defines the criteria to determine if a device is part
	// of this class.
	ClassCriteria *DeviceClassDetail `json:"classCriteria,omitempty"`
}

// DeviceClassDetail defines the subset of devices to consider part of this
// class, along with device-specific configuration information for those
// devices.
type DeviceClassDetail struct {
	// Driver specifies the driver that should handle this class of devices.
	// When a DeviceClaim uses this class, only devices published by the
	// specified driver will be considered.
	// +required
	Driver string `json:driver,omitempty`

	// Constraints is a CEL expression that operates on device attributes,
	// and must evaluate to true for a device to be considered. It will be
	// ANDed with any Constraints field in the DeviceClaim using this class.
	//
	// +optional
	Constraints *string `json:"constraints,omitempty"`

	// DeviceConfigs contains references to arbitrary vendor device configuration
	// objects that will be attached to the device allocation.
	//
	// +optional
	// +listType=atomic
	Configs []DeviceClassConfigReference `json:"configs,omitempty"`
}

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
	// +listType=atomic
	MatchAttributes []string `json:"matchAttributes,omitempty"`

	// Claims contains the actual claim details, arranged into groups
	// containing claims which must all be satsified, or for which only
	// one needs to be satisfied.
	//
	// +required
	// +listType=atomic
	Claims []DeviceClaimInstance `json:"claims,omitempty"`
}

// DeviceClaimInstance captures a claim which must be satisfied,
// or a group for which one must be sastisfied.
type DeviceClaimInstance struct {
	// At least one of AllOf and OneOf must be populated.

	// If fields of DeviceClaimDetail are populated, OneOf should
	// be empty.
	*DeviceClaimDetail `json:",inline"`

	// OneOf contains a list of claims, only one of which must be satisfied.
	// Claims are listed in order of priority.
	//
	// +optional
	// +listType=atomic
	OneOf []DeviceClaimDetail `json:"oneOf,omitempty"`
}

// DeviceClaimDetail contains the details of how to fulfill a specific
// request for devices.
type DeviceClaimDetail struct {
	// DeviceClass is the name of the DeviceClass to which the requested
	// devices must belong.
	//
	// +required
	DeviceClass string `json:"deviceClass"`

	// AdminAccess indicates that this is a claim for administrative access
	// to the devices. Claims with AdminAccess are expected to be used for
	// monitoring or other management services for a device.  They ignore
	// all ordinary claims to the device with respect to access modes and
	// any resource allocations. Ability to create these claims is
	// controlled via ResourceQuota.
	//
	// Default is false. If true, a DeviceClass must be specified so that
	// the Driver is known, and only the Constraints fields will be taken
	// into account for device selection. All devices meeting the
	// Constraints expressions will be made available as part of the claim
	// and may be assigned to a container.
	//
	// NOTE: This cannot appear in class because the quota code does not do
	// an indirection. Still searching for a better way to handle this use
	// case, without creating a new top-level type for it.
	//
	// +optional
	AdminAccess *bool `json:"adminAccess,omitempty"`

	// Constraints is a CEL expression that operates on device attributes.
	// In order for a device to be considered, this CEL expression and the
	// Constraints expression from the DeviceClass must both be true.
	//
	// +optional
	Constraints *string `json:"constraints,omitempty"`

	// Device claims may be satisfied by choosing multiple devices instead
	// of just a single device.

	// MaxDevices allows the user to specify a maximum number of devices
	// that may be allocated to sastisfy this claim. The default is equal
	// to requests.devices. Thus, without specifying MaxDevices, the user
	// will get exactly the number of devices specified in
	// requests.devices.
	//
	// +optional
	MaxDevices *int `json:"maxDevices,omitempty"`

	// Requests allows the user to specify the minimum resource
	// requirements that must be satisfied across all devices allocated for
	// this claim.
	//
	// All drivers can the support "devices", resource, which represents
	// the minimum number of devices to allocate.
	//
	// Other resource types are driver-specific. The values in `requests`
	// represent the aggregate across all devices, not the per-device
	// values.
	//
	// If not set, requests.devices will default to 1.
	//
	// +optional
	Requests map[string]resource.Quantity `json:"requests,omitempty"`

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
	// +listType=atomic
	MatchAttributes []string `json:"matchAttributes,omitempty"`

	// Configs contains references to arbitrary vendor device configuration
	// objects that will be attached to the device allocation.
	// +optional
	// +listType=atomic
	Configs []DeviceConfigReference `json:"configs,omitempty"`
}

// DeviceClaimStatus contains the results of the claim allocation.
type DeviceClaimStatus struct {
	// ClassConfigs contains the entire set of dereferenced vendor
	// configurations from the DeviceClass, as of the time of allocation.
	// +optional
	// +listType=atomic
	ClassConfigs []runtime.RawExtension `json:"classConfigs,omitempty"`

	// ClaimConfigs contains the entire set of dereferenced vendor
	// configurations from the DeviceClaim, as of the time of allocation.
	// +optional
	// +listType=atomic
	ClaimConfigs []runtime.RawExtension `json:"claimConfigs,omitempty"`

	// Allocations contains the list of device allocations needed to
	// satisfy the claim, one per pool from which devices were allocated.
	//
	// Note that the "current capacity" of the cluster is the result of
	// applying all such allocations to the published DevicePools. This
	// means storing these allocations only in claim status fields is
	// likely to scale poorly, and we will need a different strategy in the
	// real code. For example, we may need to accumulate these in the
	// DevicePool status fields themselves, and just reference them from
	// here.
	//
	// This field is owned by the scheduler, whereas the DeviceStatuses
	// field is owned by the drivers.
	//
	// +optional
	// +listType=atomic
	Allocations []DeviceAllocation `json:"allocations,omitempty"`

	// DeviceStatuses contains the status of each device allocated for this
	// claim, as reported by the driver. This can include driver-specific
	// information. Entries are owned by their respective drivers.
	//
	// +optional
	// +listType=map
	// +listMapKey=devicePoolName
	// +listMapKey=deviceName
	DeviceStatuses []AllocatedDeviceStatus `json:"deviceStatuses,omitempty"`

	// PodNames contains the names of all Pods using this claim.
	// TODO: Can we just use ownerRefs instead?
	//
	// +optional
	// +listType=set
	PodNames []string `json:"podNames,omitempty"`
}

// NOTE: We no longer have DeviceClaimTemplate. Instead, the PodSpec will
// directly contain a either a DeviceClaimName (to enable multiple pods to
// refer to a pre-provisioned claim), or an embedded struct that includes
// ObjectMeta and a DeviceClaimName. In this case, the named DeviceClaim will
// be treated as a template; that is, its spec will be copied to create a new
// claim, based on the new metadata. Re-using claim-as-a-template avoids
// another, nearly identical top-level API object. But it may be confusing, we
// need feedback.
//
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
	//
	// +optional
	// +listType=atomic
	Allocations []ResourceAllocation `json:"allocations,omitempty"`
}

// ResourceAllocation contains the per-device resource allocations.
type ResourceAllocation struct {
	// Name is the resource name/string for this allocation.
	// +required
	Name string `json:"name"`

	// Amount is the amount of resource allocated for this claim.
	// +required
	Allocation resource.Quantity `json:"allocation,omitempty"`
}

type DeviceIP struct {
	// IP is the IP address assigned to the device
	IP string `json:"ip,omitempty"`
}

// AllocatedDeviceStatus contains the status of an allocated device, if the
// driver chooses to report it. This may include driver-specific information.
type AllocatedDeviceStatus struct {
	// DevicePoolName is the name of the DevicePool to which this
	// device belongs. The driver for that device pool owns this
	// entry.
	//
	// +required
	DevicePoolName string `json:"devicePoolName"`

	// DeviceName contains the name of the allocated Device.
	//
	// +required
	DeviceName string `json:"deviceName,omitempty"`

	// Conditions contains the latest observation of the device's state.
	// If the device has been configured according to the class and claim
	// config references, the `Ready` condition should be True.
	//
	// +optional
	// +listType=atomic
	Conditions []metav1.Condition `json:"conditions"`

	// DeviceIPs contains all of the IPs allocated for a device, if any.
	//
	// +optional
	// +listType=atomic
	DeviceIPs []DeviceIP `json:"deviceIPs,omitempty"`

	// Arbitrary driver-specific data.
	//
	// +optional
	// +listType=atomic
	DeviceInfo []runtime.RawExtension `json:"deviceInfo,omitempty"`
}
