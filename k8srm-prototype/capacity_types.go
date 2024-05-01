package main

import (
	"k8s.io/apimachinery/pkg/api/resource"
)

// This prototype models nodes as a collection of resource
// pools, each populated by resources, which in turn hold
// capacities.
type NodeResources struct {
	Name string `json:"name"`

	Pools []ResourcePool `json:"pools"`
}

type ResourcePool struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`

	// attributes for constraints at the pool level
	Attributes []Attribute `json:"attributes,omitempty"`

	Resources []Resource `json:"resources,omitempty"`
}

type Resource struct {
	Name string `json:"name"`

	// attributes for constraints
	Attributes []Attribute `json:"attributes,omitempty"`

	// capacities that can be allocated
	Capacities []Capacity `json:"capacities,omitempty"`
}

type Capacity struct {
	Name string `json:"name"`

	Topologies []Topology `json:"topologies,omitempty"`

	// exactly one of the fields should be populated
	// examples implemented:
	//  - counter: integer capacity decremented by integers
	//  - quantity: resource.Quantity capacity decremented by quantities
	//  - block:  resource.Quantity capacity decremented in blocks
	//  - accessMode: allows various controlled access:
	//       - readonly-shared: allowed with other consumers using *-shared, write-exclusive
	//       - readwrite-shared: allowed with other consumers using *-shared
	//       - write-exclusive: allowed other consumers using readonly-shared
	//       - readwrite-exclusive: no other consumers allowed

	// +optional
	Counter *ResourceCounter `json:"counter,omitempty"`

	// +optional
	Quantity *ResourceQuantity `json:"quantity,omitempty"`

	// +optional
	Block *ResourceBlock `json:"block,omitempty"`

	// +optional
	AccessMode *ResourceAccessMode `json:"accessMode,omitempty"`
}

type Attribute struct {
	Name string `json:"name"`

	// One of the following:
	StringValue   *string            `json:"stringValue,omitempty"`
	IntValue      *int               `json:"intValue,omitempty"`
	QuantityValue *resource.Quantity `json:"quantityValue,omitempty"`
	SemVerValue   *SemVer            `json:"semVerValue,omitempty"`
}

type SemVer string

// This prototype does not address limitations of actuation. We
// would need to do that in the real deal. For example, today
// topology acuation is controlled at the node level, so it is not
// something we can just arbitrarily assign for any node. Instead,
// we need to look at the static topology policy of the node, and evaluate
// if that node assignment can meet the topology constraint in the request
// based upon that policy.
type Topology struct {
	Name string `json:"name"`
	Type string `json:"type"`

	// GroupInResource allows a claim to be satisfied by capacities from
	// different topologies, but in the same resource.
	GroupInResource bool `json:"groupInResource"`
}

type ResourceCounter struct {
	Capacity int64 `json:"capacity"`
}

type ResourceQuantity struct {
	Capacity resource.Quantity `json:"capacity"`
}

type ResourceBlock struct {
	Size     resource.Quantity `json:"size"`
	Capacity resource.Quantity `json:"capacity"`
}

type ResourceAccessMode struct {
	// if not allowed, any requests for that access mode
	// will be converted to a request for the next highest
	// allowed access mode.
	AllowReadOnlyShared  bool `json:"allowReadOnlyShared"`
	AllowReadWriteShared bool `json:"allowReadWriteShared"`
	AllowWriteExclusive  bool `json:"allowWriteExclusive"`

	// tracks reference counts for each access mode
	ReadOnlyShared     int `json:"readOnlyShared"`
	ReadWriteShared    int `json:"readWriteShared"`
	WriteExclusive     int `json:"writeExclusive"`
	ReadWriteExclusive int `json:"readWriteExclusive"`
}
