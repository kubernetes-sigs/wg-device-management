package main

// This prototype demonstrates allocating capacity from nodes,
// adhering to the claim constraints and requests.
// Currently, allocations are for a pod, and on a single node. However,
// the general framework should be extensible across multi-pod workloads and
// multi-node capacity.

// NodeAllocationResult contains the results of an attempt to satisfy a
// set of CapacityClaims (e.g., for a pod) against a node.
type NodeAllocationResult struct {
	NodeName             string                `json:"nodeName"`
	CapacityClaimResults []CapacityClaimResult `json:"capacityClaimResults"`
}

// CapacityClaimResult contains the results of an attempt to satisfy a
// CapacityClaim against a collection of pools (typically a node)
type CapacityClaimResult struct {
	ClaimName            string                `json:"claimName"`
	ResourceClaimResults []ResourceClaimResult `json:"resourceClaimResults,omitempty"`
}

// ResourceClaimResult contains the results of an attempt to satisfy a
// ResourceClaim against a collection of pools (typically a node)
type ResourceClaimResult struct {
	ClaimName     string       `json:"claimName"`
	PoolResults   []PoolResult `json:"poolResults"`
	Best          int          `json:"best"`
	FailureReason string       `json:"failureReason,omitempty"`
}

// PoolResult contains the results of an attempt to satisfy a
// resource claim against a specific resource pool.
type PoolResult struct {
	PoolName        string           `json:"poolName"`
	ResourceResults []ResourceResult `json:"resourceResults"`
	Best            int              `json:"best"`
	FailureReason   string           `json:"failureReason,omitempty"`
}

// ResourceResult containst the results of an attempt to satisfy a
// resource claim against a specific resource.

type ResourceResult struct {
	ResourceName    string           `json:"resourceName"`
	CapacityResults []CapacityResult `json:"capacityResults"`

	// Score is a number from 0 to 100, with 0 meaning that
	// the resource was unable to satisfy the request, and
	// 100 meaning the request was satisfied optimally. Numbers
	// in between indicate that the request can be met, but
	// sub-optimally - for example, by splitting requests
	// across topologies, or without all the preferred topological
	// alignments.
	Score         int    `json:"score"`
	FailureReason string `json:"failureReason,omitempty"`
}

// CapacityResult is the result of an attempt to allocate capacity from
// a resource. If successful, it includes the amount allocated and the
// specific topological assignment.
type CapacityResult struct {
	// If successful, CapacityRequest contains the allocated amount (which may
	// be different than the original, requested amount. If unsuccessful, it
	// contains the original requested amount.
	CapacityRequest `json:",inline"`

	// Topologies contains the topology assignments of the request allocation. Note
	// that exactly one of each topology type from the original Capacity must be in
	// this list. It is possible for the same requested capacity type, we split the
	// request across multiple topologies. This is the case, for example, if a
	// single memory request cannot be satisfied by a single NUMA node.
	Topologies []TopologyAssignment `json:"topologies,omitempty"`

	// FailureReason is non-empty if the allocation was unsuccessful
	FailureReason string `json:"failureReason,omitempty"`
}

// TopologyAssignment contains the specific topology from which a capacity is drawn.
type TopologyAssignment struct {
	Type string `json:"type"`
	Name string `json:"name"`
}
