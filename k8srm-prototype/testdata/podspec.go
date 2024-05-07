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
        ResourceClaimName *string
        ResourceClaimTemplateName *string
}

// This proposal instead has:
//

type PodSpecProposed struct {
	// lots of stuff, ending with:

	DeviceClaims *PodDeviceClaim
}

type PodDeviceClaim struct {

	// Exactly one must be populated

	Devices []DirectDeviceClaim

	ClaimName *string

	Template *DeviceClaimTemplate
}

type DirectDeviceClaim struct {
	Name string
	Class string
}

type DeviceClaimTemplate struct {
	ObjectMeta `json:",inline"`

	// +required
	ClaimName `json:"claimName"`
}

