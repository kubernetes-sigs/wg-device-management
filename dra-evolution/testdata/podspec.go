package podspec

// Today in podspec, we have:
//

type PodSpec130 struct {
	// lots of stuff, ending with:

	ResourceClaims []PodResourceClaim
}

type PodResourceClaim struct {
	// Name uniquely identifies this resource claim inside the pod.
	// This must be a DNS_LABEL.
	Name string

	// Source describes where to find the ResourceClaim.
	Source ClaimSource
}

type ClaimSource struct {
	ResourceClaimName         *string
	ResourceClaimTemplateName *string
}

// Claims at the pod level get referenced in containers with:
type ResourceRequirements struct {
	// lots of stuff, ending with:

	// Claims lists the names of resources, defined in spec.resourceClaims,
	// that are used by this container.
	//
	// This is an alpha field and requires enabling the
	// DynamicResourceAllocation feature gate.
	//
	// This field is immutable. It can only be set for containers.
	//
	// +featureGate=DynamicResourceAllocation
	// +optional
	Claims []ResourceClaim
}

// ResourceClaim references one entry in PodSpec.ResourceClaims.
type ResourceClaim struct {
	// Name must match the name of one entry in pod.spec.resourceClaims of
	// the Pod where this field is used. It makes that resource available
	// inside a container.
	Name string
}

// This proposal instead has:
//

type PodSpecProposed struct {
	// lots of stuff, ending with:

	ResourceClaims []PodResourceClaim
}

type PodResourceClaim struct {
	// Name uniquely identifies this resource claim inside the pod.
	// This must be a DNS_LABEL.
	Name string

	// Source describes where to find the ResourceClaim.
	Source ClaimSource
}

type ClaimSource struct {
	ForClass                  *ResourceClaimForClass
	ResourceClaimName         *string
	ResourceClaimTemplateName *string
}

type ResourceClaimForClass struct {
	ClassName string

	// ObjectMeta may contain labels and annotations that will be copied into the PVC
	// when creating it. No other fields are allowed and will be rejected during
	// validation.
	// +optional
	metav1.ObjectMeta
}

type ResourceRequirements struct {
	// lots of stuff, ending with:

	// Claims lists the names of resources, defined in spec.resourceClaims,
	// that are used by this container. All devices allocated through each
	// resource claim are made available to the container.
	//
	// This is an alpha field and requires enabling the
	// DynamicResourceAllocation feature gate.
	//
	// This field is immutable. It can only be set for containers.
	//
	// +featureGate=DynamicResourceAllocation
	// +optional
	Claims []ResourceClaim

	// Devices lists names of resources, defined in spec.resourceClaims,
	// and the name of individual device requests, defined in resourceClaim.requests,
	// that are used by this container.
	Devices []ResourceClaimDevice
}

// ResourceClaimDevice references specific devices inside a claim. If the named request
// is satisfied by allocating multiple devices, then all of those are matched.
type ResourceClaimDevice struct {
	// pod.spec.resourceClaims
	ClaimName string
	// claim.spec.requests
	RequestName string
}
