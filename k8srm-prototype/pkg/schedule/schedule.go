package schedule

import (
	"fmt"

	"github.com/kubernetes-sigs/wg-device-management/k8srm-prototype/pkg/api"

	"gonum.org/v1/gonum/stat/combin"
)

// SelectNode will select the node that can best satisfy all the claims.
//
// Prior to passing in the list of pools, the caller should apply any existing
// allocations to those pools. That is, the device counts in the pools should
// be reduced to only the available devices. The algorithm here will consider
// allocations made only by claims passed to this function.
//
// The first returned value is an array of the allocations from each pool that
// are needed to satisfy all the claims. In the event no node can be selected,
// this will be empty. The second returned value is an array of the results of
// evaluating each node.

func SelectNode(claims []api.DeviceClaim, pools []api.DevicePool) ([]api.DevicePoolAllocation, []NodeResult) {
	// Collect the pools by node
	poolsByNode := make(map[string][]api.DevicePool)
	for _, p := range pools {
		// ignore pools not associated with a node
		if p.Spec.NodeName == nil || *p.Spec.NodeName == "" {
			continue
		}

		poolsByNode[*p.Spec.NodeName] = append(poolsByNode[*p.Spec.NodeName], p)
	}

	var results []NodeResult
	i := -1
	best := -1
	// Evaluate each node against the claims
	for node, nodeDevPools := range poolsByNode {
		nr := evaluateNode(node, claims, nodeDevPools)
		results = append(results, nr)
		i += 1

		if best > -1 && nr.Score() > results[best].Score() {
			best = i
			continue
		}

		if best == -1 && nr.Score() > 0 {
			best = i
		}
	}

	if best == -1 {
		return nil, results
	}

	return results[best].Allocations(), results
}

func evaluateNode(node string, claims []api.DeviceClaim, pools []api.DevicePool) NodeResult {
	nr := NodeResult{
		NodeName: node,
	}

	// Theoretically, when there are multiple claims, the order in which
	// they are considered may make a difference. In general without
	// partitioning and without per-device resources, it should be rare
	// for this to be an issue. The `MatchAttributes` functionality, along
	// with multiple claims for the similar devices, could conceivably
	// result in one order being solvable, and another not being solvable.
	//
	// Regardless, for the prototype we will not worry about this, and will
	// just evaluate the claims in the order presented.
	for _, c := range claims {
		dcr := evaluateNodeForClaim(c, pools)
		nr.DeviceClaimResults = append(nr.DeviceClaimResults, dcr)

		// TODO: apply the allocations to the underlying pools, so that
		// subsequent, overlapping claims do not double-allocate
	}

	return nr
}

func evaluateNodeForClaim(claim api.DeviceClaim, pools []api.DevicePool) DeviceClaimResult {
	dcr := DeviceClaimResult{
		ClaimName: claim.Name,
		Best:      -1,
	}

	// This function implements an algorithm which assumes that allocating
	// multiple devices out of the same pool is better than allocating them
	// out of different pools. This allows device drivers to organize their
	// pools by some coarse topology like NUMA. We may want to break this
	// out into an interface and experiment with a couple algorithms. For
	// example, we could "flatten" all pools into just a set of devices,
	// and treat Node as a MatchAttributes field.

	// Pool-Based Algorithm Assumptions
	//
	// A brute force approach would enumerate every permutation of 1 or
	// more pools, and then score the claim against each. We can do some
	// early pruning by following these principles:
	// - Any pools that do not match the constraints can be
	//   removed from consideration, reducing the combinatorial effect.
	// - Any pools with no available devices can similarly be removed
	//   (recall that the passed in pools have their availability reduced
	//   according to any existing allocations).
	// - Ordering of the pools does not matter when evaluating single claim
	//   against a list of pools. Thus we can evaluate combinations (sets)
	//   rather than permutations of pools.
	// - A solution with fewer pools is always better. This means we can
	//   start with sets of one pool, only proceeding to two pool
	//   combinations if no single pool works, and so on.

	// First, eliminate any non-matching or fully committed pools.
	var goodPools []api.DevicePool
	for _, p := range pools {
		if p.Spec.DeviceCount <= 0 {
			dcr.IgnoredPools = append(dcr.IgnoredPools, PoolResult{
				PoolName:      p.Name,
				FailureReason: "no available devices",
			})
			continue
		}

		if claim.Spec.Driver != nil && *claim.Spec.Driver != p.Spec.Driver {
			dcr.IgnoredPools = append(dcr.IgnoredPools, PoolResult{
				PoolName:      p.Name,
				FailureReason: "claim and pool driver do not match",
			})
			continue
		}

		// TODO: Consider *class* contraints
		meets, err := MeetsConstraints(claim.Spec.Constraints, p.Spec.Attributes)
		if err != nil {
			dcr.IgnoredPools = append(dcr.IgnoredPools, PoolResult{
				PoolName:      p.Name,
				FailureReason: fmt.Sprintf("error evaluating constraints: %s", err.Error()),
			})
			continue
		}
		if !meets {
			dcr.IgnoredPools = append(dcr.IgnoredPools, PoolResult{
				PoolName:      p.Name,
				FailureReason: "constraints not met",
			})
			continue
		}

		goodPools = append(goodPools, p)
	}

	// Now, iterate through the possible lengths of different combination sets,
	// scoring each set. If we find one or more successful sets at a given size,
	// we do not need to continue evaluating the next size, based on the principle
	// stated above.
	for setSize := 1; setSize <= len(goodPools); setSize++ {
		combinations := combin.Combinations(len(goodPools), setSize)
		for _, combo := range combinations {
			psr := evaluatePoolSetForClaim(claim, poolSet(goodPools, combo))
			dcr.PoolSetResults = append(dcr.PoolSetResults, psr)
			if psr.Score > 0 {
				if dcr.Best == -1 || psr.Score > dcr.PoolSetResults[dcr.Best].Score {
					dcr.Best = len(dcr.PoolSetResults) - 1
				}

			}
		}
		// if we found at least one that works, we do not need to check
		// the next setSize, and we are done
		if dcr.Best > -1 {
			return dcr
		}
	}

	return dcr
}

func poolSet(pools []api.DevicePool, combo []int) []api.DevicePool {
	var result []api.DevicePool

	for _, index := range combo {
		result = append(result, pools[index])
	}

	return result
}

// attempts to satisfy the claim using the specified pools
// Assumptions (very important!):
//   - All the passed pools meet the class and claim constraints
//   - No pool is fully exhausted (corollary of the first assumption)
//   - No subset will satisfy the claim. That is, we must use ALL the
//     passed pools. This is important otherwise we need to consider
//     MatchAttributes across permutations, not combinations.
func evaluatePoolSetForClaim(claim api.DeviceClaim, pools []api.DevicePool) PoolSetResult {
	required := 1
	if claim.Spec.MinDeviceCount != nil {
		required = *claim.Spec.MinDeviceCount
	}

	origRequired := required

	psr := PoolSetResult{}

	// TODO: Include both claim and class MatchAttributes
	matchAttrs := make(map[string]api.Attribute)

	for i, p := range pools {
		pr := PoolResult{
			PoolName: p.Name,
		}

		// For the first pool, grab the values of the MatchAttributes.
		// All subsequent pools must have the same values.
		for _, attrName := range claim.Spec.MatchAttributes {
			for _, attr := range p.Spec.Attributes {
				if attrName == attr.Name {
					if i == 0 {
						matchAttrs[attrName] = attr
					} else if !matchAttrs[attrName].Equal(attr) {
						pr.FailureReason = "claim MatchAttributes constraint failed"
					}
				}
			}
		}

		if pr.FailureReason != "" {
			psr.PoolResults = append(psr.PoolResults, pr)
			continue
		}

		if p.Spec.DeviceCount <= required {
			pr.DeviceCount = p.Spec.DeviceCount
			required -= p.Spec.DeviceCount
		} else {
			pr.DeviceCount = required
			required = 0
		}

		psr.PoolResults = append(psr.PoolResults, pr)
	}

	if required > 0 {
		psr.Score = 0
		psr.FailureReason = fmt.Sprintf("unable to satisfy %d of %d device requests", required, origRequired)
	} else {
		psr.Score = 100
	}
	return psr
}
