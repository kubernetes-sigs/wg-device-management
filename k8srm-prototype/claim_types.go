package main

import (
	"k8s.io/apimachinery/pkg/api/resource"
)

// This prototype models requests for capacity as ResourceClaims, which
// are structured based on the workload structure. For example, we can
// request capacity required to run a pod. This includes resource claims
// for the pod itself (for example, we have a counter for number of pods
// allowed on a node), as well as resource claims for each container in
// the pod. Claims may include CEL-based constraints, as well as topological
// constraints. Those topological constraints may apply to the whole pod
// (equivalent to topology manager scope=pod), or to individual containers
// in the pod (equivalent to topology manager scope=container).

type PodCapacityClaim struct {
	// PodClaim contains the resource claims needed for the pod
	// level, such as the pod capacity needed to run pods,
	// or devices that may be attached to a container later
	// +required
	PodClaim CapacityClaim `json:"podClaim"`

	// ContainerClaims contains the resource claims needed on a
	// per-container level, such as CPU and memory
	// +required
	ContainerClaims []CapacityClaim `json:"containerClaims"`

	// Topologies specifies the topological alignment and preferences
	// across all containers and resources in the pod
	// +optional
	Topologies []TopologyConstraint `json:"topologies,omitempty"`
}

type CapacityClaim struct {
	// Name is used to identify the capacity claim to help in troubleshooting
	// unschedulable claims.
	// +required
	Name string `json:"name"`

	// Claims contains the set of resource claims that are part of
	// this capacity claim
	// +required
	Claims []ResourceClaim `json:"claims"`

	// Topologies specifies the topological alignment and preferences
	// across all resources in this capacity claim
	// +optional
	Topologies []TopologyConstraint `json:"topologies,omitempty"`
}

type ResourceClaim struct {
	// Name is used to identify the resource claim to help in troubleshooting
	// unschedulable claims.
	// +required
	Name string `json:"name"`

	// Driver will limit the scope of resources considered
	// to only those published by the specified driver
	// +optional
	Driver string `json:"driver,omitempty"`

	// Constraints is a CEL expression that operates on
	// node and resource attributes, and must evaluate to true
	// for a resource to be considered
	// +optional
	Constraints string `json:"constraints,omitempty"`

	// Topologies specifies topological alignment constraints and
	// preferences for the allocated capacities. These constraints
	// apply across the capacities within the resource.
	// +optional
	Topologies []TopologyConstraint `json:"topologies,omitempty"`

	// Capacities specifies the individual allocations needed
	// from the capacities provided by the resource
	// +required
	Capacities []CapacityRequest `json:"capacities"`
}

type TopologyConstraint struct {
	// Type identifies the type of topology to constrain
	// +required
	Type string `json:"type"`

	// Policy defines the specific constraint. All types support 'prefer'
	// and 'require', with 'prefer' being the default. 'Prefer' means
	// that allocations will be made according to the topology when
	// possible, but the allocation will not fail if the constraint cannot
	// be met. 'Require' will fail the allocation if the constraint is not
	// met. Types may add additional policies.
	// +optional
	Policy string `json:"policy,omitempty"`
}

type CapacityRequest struct {
	// +required
	Capacity string `json:"capacity"`

	// one of these must be populated
	// note that we only need three different type of capacity requests
	// even though we have four different types of capacity models
	// the ResourceQuantity and ResourceBlock capacity models both
	// are drawn down on via the ResourceQuantityRequest type
	Counter    *ResourceCounterRequest    `json:"counter,omitempty"`
	Quantity   *ResourceQuantityRequest   `json:"quantity,omitempty"`
	AccessMode *ResourceAccessModeRequest `json:"accessMode,omitempty"`
}

type ResourceCounterRequest struct {
	// +required
	Request int64 `json:"request"`
}

type ResourceQuantityRequest struct {
	// +required
	Request resource.Quantity `json:"request"`
}

type CapacityAccessMode string

const (
	ReadOnlyShared     = "ReadOnlyShared"
	ReadWriteShared    = "ReadWriteShared"
	WriteExclusive     = "WriteExclusive"
	ReadWriteExclusive = "ReadWriteExclusive"
)

type ResourceAccessModeRequest struct {
	// +required
	Request CapacityAccessMode `json:"request"`
}
