package schedule

import (
	"fmt"
	"strings"

	"github.com/kubernetes-sigs/wg-device-management/k8srm-prototype/pkg/api"
)

type NodeResult struct {
	NodeName           string
	DeviceClaimResults []DeviceClaimResult
}

// DeviceClaimResult contains the results of an attempt to satisfy a
// DeviceClaim against a collection of pools (typically a node)
type DeviceClaimResult struct {
	ClaimName      string          `json:"claimName"`
	PoolSetResults []PoolSetResult `json:"poolSetResults"`
	Best           int             `json:"best"`

	IgnoredPools []PoolResult `json:"ignoredPools,omitempty"`
}

// PoolSetResult contains the results of an attempt to satisfy a
// DeviceClaim with a given set of DevicePools.
type PoolSetResult struct {
	PoolResults []PoolResult `json:"poolResults"`
	Score       int          `json:"score"`

	FailureReason string `json:"failureReason,omitempty"`
}

// PoolResult contains the results of an attempt to satisfy a
// device claim against a specific device pool.
type PoolResult struct {
	PoolName    string `json:"poolName"`
	DeviceCount int    `json:"deviceCount"`

	FailureReason string `json:"failureReason,omitempty"`
}

// NodeResult methods

func (nr *NodeResult) Score() int {
	// The score for this node is zero if any
	// claim could not be satisfied, and the average
	// score for all claims otherwise.
	sum := 0
	for _, dcr := range nr.DeviceClaimResults {
		if dcr.Score() == 0 {
			return 0
		}

		sum += dcr.Score()
	}

	return sum / len(nr.DeviceClaimResults)
}

func (nr *NodeResult) Allocations() []api.DevicePoolAllocation {
	if nr.Score() == 0 {
		return nil
	}

	var allocations []api.DevicePoolAllocation
	for _, dcr := range nr.DeviceClaimResults {
		allocations = append(allocations, dcr.Allocations()...)
	}

	return allocations
}

func (nr *NodeResult) Summary() string {
	if nr.Score() > 0 {
		return fmt.Sprintf("%s: satisfied all claims with score %d", nr.NodeName, nr.Score())
	}

	var unsatisfied []string
	for _, dcr := range nr.DeviceClaimResults {
		if dcr.Score() == 0 {
			unsatisfied = append(unsatisfied, dcr.ClaimName)
		}
	}

	return fmt.Sprintf("%s: could not satisfy these claims: %s", nr.NodeName, strings.Join(unsatisfied, ", "))
}

// DeviceClaimResult methods

func (dcr *DeviceClaimResult) Score() int {
	if dcr.Best == -1 {
		return 0
	}

	return dcr.PoolSetResults[dcr.Best].Score
}

func (dcr *DeviceClaimResult) Allocations() []api.DevicePoolAllocation {
	if dcr.Best == -1 {
		return nil
	}

	return dcr.PoolSetResults[dcr.Best].Allocations()
}

// PoolSetResult methods

func (psr *PoolSetResult) Allocations() []api.DevicePoolAllocation {
	if psr.Score == 0 {
		return nil
	}

	var results []api.DevicePoolAllocation
	for _, pr := range psr.PoolResults {
		results = append(results, pr.DevicePoolAllocation())
	}

	return results
}

// PoolResult methods
func (pr *PoolResult) DevicePoolAllocation() api.DevicePoolAllocation {
	return api.DevicePoolAllocation{
		DevicePoolName: pr.PoolName,
		DeviceCount:    pr.DeviceCount,
	}
}
