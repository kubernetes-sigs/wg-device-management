package main

import (
	"slices"

	resourceapi "github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/resource"
)

// LoopControl is a variable to track Continues and Breaks from a recursive function.
type LoopControl int

const (
	// Continue indicates the "loop" of the recursive function should continue.
	Continue LoopControl = iota
	// Break indicates the "loop" of the recursive function should break.
	Break
)

// Device is an alias of resourceapi.Device
type Device = resourceapi.Device

// Devices is an alias of []resourceapi.Devices
type Devices []resourceapi.Device

// GetNames returns a list of all device names in Devices.
func (devices Devices) GetNames() []string {
	var names []string
	for _, device := range devices {
		names = append(names, device.Name)
	}
	slices.Sort(names)
	return names
}

// IterateCombinations iterates through all combinations of the provided devices.
func (devices Devices) IterateCombinations(f func(Devices) LoopControl) {
	slices.SortFunc(devices, func(a, b resourceapi.Device) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})

	var iterate func(i int, accum Devices) LoopControl
	iterate = func(i int, accum Devices) LoopControl {
		accum = append(accum, devices[i])
		control := f(accum)
		if control == Break {
			return Break
		}
		for j := i + 1; j < len(devices); j++ {
			iterate(j, accum)
		}
		return Continue
	}

	for i := 0; i < len(devices); i++ {
		iterate(i, Devices{})
	}
}
