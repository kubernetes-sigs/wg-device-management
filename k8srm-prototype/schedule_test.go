package main

import (
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

func TestCapacityReduce(t *testing.T) {
	testCases := map[string]struct {
		capacity Capacity
		request  CapacityRequest
		result   Capacity
	}{
		"counter": {
			capacity: Capacity{
				Name:    "counter-test",
				Counter: &ResourceCounter{Capacity: 10},
			},
			request: CapacityRequest{
				Capacity: "counter-test",
				Counter:  &ResourceCounterRequest{Request: 4},
			},
			result: Capacity{
				Name:    "counter-test",
				Counter: &ResourceCounter{Capacity: 6},
			},
		},
		"quantity": {
			capacity: Capacity{
				Name:     "quantity-test",
				Quantity: &ResourceQuantity{Capacity: resource.MustParse("10M")},
			},
			request: CapacityRequest{
				Capacity: "quantity-test",
				Quantity: &ResourceQuantityRequest{Request: resource.MustParse("1M")},
			},
			result: Capacity{
				Name:     "quantity-test",
				Quantity: &ResourceQuantity{Capacity: resource.MustParse("9M")},
			},
		},
		"block": {
			capacity: Capacity{
				Name: "block-test",
				Block: &ResourceBlock{
					Capacity: resource.MustParse("10M"),
					Size:     resource.MustParse("1M"),
				},
			},
			request: CapacityRequest{
				Capacity: "block-test",
				Quantity: &ResourceQuantityRequest{Request: resource.MustParse("1M")},
			},
			result: Capacity{
				Name: "block-test",
				Block: &ResourceBlock{
					Capacity: resource.MustParse("9M"),
					Size:     resource.MustParse("1M"),
				},
			},
		},
		"accessMode-readonlyshared": {
			capacity: Capacity{
				Name: "access-test",
				AccessMode: &ResourceAccessMode{
					AllowReadOnlyShared: true,
					ReadOnlyShared:      3,
				},
			},
			request: CapacityRequest{
				Capacity:   "access-test",
				AccessMode: &ResourceAccessModeRequest{Request: ReadOnlyShared},
			},
			result: Capacity{
				Name: "access-test",
				AccessMode: &ResourceAccessMode{
					AllowReadOnlyShared: true,
					ReadOnlyShared:      4,
				},
			},
		},
		"accessMode-readwriteshared": {
			capacity: Capacity{
				Name: "access-test",
				AccessMode: &ResourceAccessMode{
					AllowReadWriteShared: true,
					ReadWriteShared:      3,
				},
			},
			request: CapacityRequest{
				Capacity:   "access-test",
				AccessMode: &ResourceAccessModeRequest{Request: ReadWriteShared},
			},
			result: Capacity{
				Name: "access-test",
				AccessMode: &ResourceAccessMode{
					AllowReadWriteShared: true,
					ReadWriteShared:      4,
				},
			},
		},
	}
	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			result, err := tc.capacity.reduce(tc.request)
			require.NoError(t, err)
			require.Equal(t, tc.result, result)
		})
	}

}

func TestResourceReduceCapacity(t *testing.T) {
	testCases := map[string]struct {
		resource    Resource
		allocations []CapacityResult
		result      Resource
		expErr      string
	}{
		"missing capacity name for allocation": {
			resource: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name:    "counter-test",
						Counter: &ResourceCounter{Capacity: 10},
					},
				},
			},
			allocations: []CapacityResult{
				{
					CapacityRequest: CapacityRequest{
						Capacity: "invalid-counter-test",
						Counter:  &ResourceCounterRequest{Request: 4},
					},
				},
			},
			expErr: `allocated capacity "invalid-counter-test" not found in resource capacities`,
		},
		"missing capacity topology for allocation": {
			resource: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name:    "counter-test",
						Counter: &ResourceCounter{Capacity: 10},
					},
				},
			},
			allocations: []CapacityResult{
				{
					CapacityRequest: CapacityRequest{
						Capacity: "counter-test",
						Counter:  &ResourceCounterRequest{Request: 4},
					},
					Topologies: []TopologyAssignment{
						{
							Type: "numa",
							Name: "numa-0",
						},
					},
				},
			},
			expErr: `allocated capacity "counter-test;numa=numa-0" not found in resource capacities`,
		},
		"single counter": {
			resource: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name:    "counter-test",
						Counter: &ResourceCounter{Capacity: 10},
					},
				},
			},
			allocations: []CapacityResult{
				{
					CapacityRequest: CapacityRequest{
						Capacity: "counter-test",
						Counter:  &ResourceCounterRequest{Request: 4},
					},
				},
			},
			result: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name:    "counter-test",
						Counter: &ResourceCounter{Capacity: 6},
					},
				},
			},
		},
		"single quantity": {
			resource: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name:     "quantity-test",
						Quantity: &ResourceQuantity{Capacity: resource.MustParse("10M")},
					},
				},
			},
			allocations: []CapacityResult{
				{
					CapacityRequest: CapacityRequest{
						Capacity: "quantity-test",
						Quantity: &ResourceQuantityRequest{Request: resource.MustParse("1M")},
					},
				},
			},
			result: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name:     "quantity-test",
						Quantity: &ResourceQuantity{Capacity: resource.MustParse("9M")},
					},
				},
			},
		},
		"single block": {
			resource: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name: "block-test",
						Block: &ResourceBlock{
							Capacity: resource.MustParse("10M"),
							Size:     resource.MustParse("1M"),
						},
					},
				},
			},
			allocations: []CapacityResult{
				{
					CapacityRequest: CapacityRequest{
						Capacity: "block-test",
						Quantity: &ResourceQuantityRequest{Request: resource.MustParse("1M")},
					},
				},
			},
			result: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name: "block-test",
						Block: &ResourceBlock{
							Capacity: resource.MustParse("9M"),
							Size:     resource.MustParse("1M"),
						},
					},
				},
			},
		},
		"multiple capacities, one allocation": {
			resource: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name:    "counter-test",
						Counter: &ResourceCounter{Capacity: 10},
					},
					{
						Name:    "counter-test-two",
						Counter: &ResourceCounter{Capacity: 10},
					},
				},
			},
			allocations: []CapacityResult{
				{
					CapacityRequest: CapacityRequest{
						Capacity: "counter-test",
						Counter:  &ResourceCounterRequest{Request: 4},
					},
				},
			},
			result: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name:    "counter-test",
						Counter: &ResourceCounter{Capacity: 6},
					},
					{
						Name:    "counter-test-two",
						Counter: &ResourceCounter{Capacity: 10},
					},
				},
			},
		},
		"multiple capacities, multiple allocations": {
			resource: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name:    "counter-test",
						Counter: &ResourceCounter{Capacity: 10},
					},
					{
						Name:    "counter-test-two",
						Counter: &ResourceCounter{Capacity: 10},
					},
				},
			},
			allocations: []CapacityResult{
				{
					CapacityRequest: CapacityRequest{
						Capacity: "counter-test",
						Counter:  &ResourceCounterRequest{Request: 4},
					},
				},
				{
					CapacityRequest: CapacityRequest{
						Capacity: "counter-test-two",
						Counter:  &ResourceCounterRequest{Request: 1},
					},
				},
			},
			result: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name:    "counter-test",
						Counter: &ResourceCounter{Capacity: 6},
					},
					{
						Name:    "counter-test-two",
						Counter: &ResourceCounter{Capacity: 9},
					},
				},
			},
		},
		"single capacity with single topology": {
			resource: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name: "counter-test",
						Topologies: []Topology{
							{
								Type:            "numa",
								Name:            "numa-0",
								GroupInResource: true,
							},
						},
						Counter: &ResourceCounter{Capacity: 10},
					},
				},
			},
			allocations: []CapacityResult{
				{
					CapacityRequest: CapacityRequest{
						Capacity: "counter-test",
						Counter:  &ResourceCounterRequest{Request: 4},
					},
					Topologies: []TopologyAssignment{
						{
							Type: "numa",
							Name: "numa-0",
						},
					},
				},
			},
			result: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name: "counter-test",
						Topologies: []Topology{
							{
								Type:            "numa",
								Name:            "numa-0",
								GroupInResource: true,
							},
						},
						Counter: &ResourceCounter{Capacity: 6},
					},
				},
			},
		},
		"single capacity, single topology type, multiple topologies": {
			resource: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name: "counter-test",
						Topologies: []Topology{
							{
								Type:            "numa",
								Name:            "numa-0",
								GroupInResource: true,
							},
						},
						Counter: &ResourceCounter{Capacity: 10},
					},
					{
						Name: "counter-test",
						Topologies: []Topology{
							{
								Type:            "numa",
								Name:            "numa-1",
								GroupInResource: true,
							},
						},
						Counter: &ResourceCounter{Capacity: 10},
					},
				},
			},
			allocations: []CapacityResult{
				{
					CapacityRequest: CapacityRequest{
						Capacity: "counter-test",
						Counter:  &ResourceCounterRequest{Request: 4},
					},
					Topologies: []TopologyAssignment{
						{
							Type: "numa",
							Name: "numa-1",
						},
					},
				},
			},
			result: Resource{
				Name: "test",
				Capacities: []Capacity{
					{
						Name: "counter-test",
						Topologies: []Topology{
							{
								Type:            "numa",
								Name:            "numa-0",
								GroupInResource: true,
							},
						},
						Counter: &ResourceCounter{Capacity: 10},
					},
					{
						Name: "counter-test",
						Topologies: []Topology{
							{
								Type:            "numa",
								Name:            "numa-1",
								GroupInResource: true,
							},
						},
						Counter: &ResourceCounter{Capacity: 6},
					},
				},
			},
		},
	}
	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			result := tc.result
			err := result.ReduceCapacity(tc.allocations)
			if tc.expErr != "" {
				require.EqualError(t, err, tc.expErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.result, result)
			}
		})
	}
}

func TestPoolReduceCapacity(t *testing.T) {
	basePool := ResourcePool{
		Name:   "primary",
		Driver: "kubelet",
		Resources: []Resource{
			{
				Name: "primary",
				Capacities: []Capacity{
					{
						Name:    "pods",
						Counter: &ResourceCounter{100},
					},
					{
						Name:    "containers",
						Counter: &ResourceCounter{1000},
					},
				},
			},
		},
	}

	singleAllocPool := basePool
	singleAllocPool.Resources[0].Capacities[0].Counter.Capacity = 96

	testCases := map[string]struct {
		pool       ResourcePool
		allocation PoolResult
		result     ResourcePool
		expErr     string
	}{
		"single allocation": {
			pool: basePool,
			allocation: PoolResult{
				PoolName: "primary",
				ResourceResults: []ResourceResult{
					{
						ResourceName: "primary",
						CapacityResults: []CapacityResult{
							{
								CapacityRequest: CapacityRequest{
									Capacity: "pods",
									Counter:  &ResourceCounterRequest{Request: 4},
								},
							},
						},
					},
				},
				Best: 0,
			},
			result: singleAllocPool,
		},
	}
	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			result := tc.result
			err := result.ReduceCapacity(tc.allocation)
			if tc.expErr != "" {
				require.EqualError(t, err, tc.expErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.result, result)
			}
		})
	}
}
