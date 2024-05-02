package resource

import (
	resourceapi "k8s.io/api/resource/v1alpha2"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/intrange"
)

// NamedResourcesAttribute is an alias of resourceapi.NamedResourcesAttribute
type NamedResourcesAttribute = resourceapi.NamedResourcesAttribute

// NamedResourcesAttributeValue is an alias of resourceapi.NamedResourcesAttributeValue
type NamedResourcesAttributeValue = resourceapi.NamedResourcesAttributeValue

// NamedResourcesSharedResource represents a shared resource that is consumable by a top-level resource when allocated.
// +k8s:deepcopy-gen=true
type NamedResourcesSharedResource struct {
	// Name is the name of the resource represented by this shared resource.
	// It must be a DNS subdomain.
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	// NamedResourcesAttributeValue is an embedded type representing the actual value of the shared resource.
	NamedResourcesSharedResourceValue `json:",inline" protobuf:"bytes,2,opt,name=value"`
}

// NamedResourcesSharedResourceValue represents the value of a shared resource.
// NamedResourcesSharedResourceValue must have one and only one field set.
// +k8s:deepcopy-gen=true
type NamedResourcesSharedResourceValue struct {
	// QuantityValue is a quantity.
	QuantityValue *resource.Quantity `json:"quantity,omitempty" protobuf:"bytes,1,opt,name=quantity"`

	// IntRangeValue is a range of 64-bit integers.
	IntRangeValue *intrange.IntRange `json:"intRange,omitempty" protobuf:"varint,2,rep,name=intRange"`
}

// NamedResourcesSharedResourceGroup represents a named group of shared resources.
// +k8s:deepcopy-gen=true
type NamedResourcesSharedResourceGroup struct {
	// Name is unique identifier among all resource groups managed by
	// the driver on the node. It must be a DNS subdomain.
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	// Items represents the list of all resources in the shared resource group.
	//
	// +listType=atomic
	// +optional
	Items []NamedResourcesSharedResource `json:"items,omitempty" protobuf:"bytes,2,opt,name=items"`
}

// ResourceModel must have one and only one field set.
// +k8s:deepcopy-gen=true
type ResourceModel struct {
	// NamedResources describes available resources using the named resources model.
	//
	// +optional
	NamedResources *NamedResourcesResources `json:"namedResources,omitempty" protobuf:"bytes,1,opt,name=namedResources"`
}

// NamedResourcesResources is used in ResourceModel.
// +k8s:deepcopy-gen=true
type NamedResourcesResources struct {
	// The list of all individual resources instances currently available.
	//
	// +listType=atomic
	Instances []NamedResourcesInstance `json:"instances" protobuf:"bytes,1,name=instances"`

	// The list of all shared resources limits that are referenced by one or
	// more instances.
	//
	// +listType=atomic
	// +optional
	SharedLimits []NamedResourcesSharedResourceGroup `json:"sharedLimits,omitempty" protobuf:"bytes,2,opt,name=sharedLimits"`
}

// NamedResourcesInstance represents one individual hardware instance that can be selected based
// on its attributes.
// +k8s:deepcopy-gen=true
type NamedResourcesInstance struct {
	// Name is unique identifier among all resource instances managed by
	// the driver on the node. It must be a DNS subdomain.
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	// Attributes defines the attributes of this resource instance.
	// The name of each attribute must be unique.
	//
	// +listType=atomic
	// +optional
	Attributes []NamedResourcesAttribute `json:"attributes,omitempty" protobuf:"bytes,2,opt,name=attributes"`

	// Resources defines the set of resources this instance consumes when
	// allocated.
	//
	// +listType=atomic
	// +optional
	Resources []NamedResourcesSharedResourceGroup `json:"resources,omitempty" protobuf:"bytes,3,opt,name=resources"`
}
