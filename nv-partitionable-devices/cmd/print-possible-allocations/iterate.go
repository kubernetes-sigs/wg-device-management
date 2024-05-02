package main

import (
	newresourceapi "github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/resource/new"
)

// LoopControl is a variable to track Continues and Breaks from a recursive function.
type LoopControl int

const (
	// Continue indicates the "loop" of the recursive function should continue.
	Continue LoopControl = iota
	// Break indicates the "loop" of the recursive function should break.
	Break
)

// NamedResourcesInstance is an alias of newresourceapi.NamedResourcesInstance
type NamedResourcesInstance = newresourceapi.NamedResourcesInstance

// NamedResourcesInstances is an alias of []newresourceapi.NamedResourcesInstance
type NamedResourcesInstances []newresourceapi.NamedResourcesInstance

// GetNames returns a list of all instances in NamedResourcesInstances.
func (instances NamedResourcesInstances) GetNames() []string {
	var names []string
	for _, instance := range instances {
		names = append(names, instance.Name)
	}
	return names
}

// IterateCombinations iterates through all combinations of a given NamedResourcesInstances.
func (instances NamedResourcesInstances) IterateCombinations(f func(NamedResourcesInstances) LoopControl) {
	var iterate func(i int, accum NamedResourcesInstances) LoopControl

	iterate = func(i int, accum NamedResourcesInstances) LoopControl {
		accum = append(accum, instances[i])
		control := f(accum)
		if control == Break {
			return Break
		}
		for j := i + 1; j < len(instances); j++ {
			iterate(j, accum)
		}
		return Continue
	}

	for i := 0; i < len(instances); i++ {
		iterate(i, NamedResourcesInstances{})
	}
}
