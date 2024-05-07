package api

// ClaimSource describes a reference to a ResourceClaim.
//
// Exactly one of these fields should be set.  Consumers of this type must
// treat an empty object as if it has an unknown value.
type ClaimSource struct {
	// DeviceClaimName is the name of a DeviceClaim object in the same
	// namespace as this pod.
	DeviceClaimName *string

	TemplateClass *ClaimSourceTemplateClass
	TemplateClaim *ClaimSourceTemplateClaim
	Generator     *ClaimSourceGenerator
}

type ClaimSourceTemplateClass struct {
	ObjectMeta `json:",inline"`

	// +required
	ClassName string `json:"className"`
}

type ClaimSourceTemplateClaim struct {
	ObjectMeta `json:",inline"`

	// +required
	ClaimName string `json:"className"`
}

type ClaimSourceGenerator struct {
	ObjectMeta `json:",inline"`

	// +required
	Generator GeneratorReference `json:"generatorRef"`
}

type GeneratorReference struct {
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
