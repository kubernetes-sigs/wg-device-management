package main

import (
	"fmt"
	"maps"

	"github.com/NVIDIA/go-nvml/pkg/nvml/mock/dgxa100"
	"k8s.io/klog/v2"

	nvdevicelib "github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/nvdevice"
	resourceapi "github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/resource"
)

// Main queries the list of allocatable devices and prints them as a kubernetes
// structure resource model.
func main() {
	// Instantiate an instance of a mock dgxa100 server and build a nvDeviceLib
	// from it. The nvDeviceLib is then used to populate the list of allocatable
	// devices from this mock server using standard NVML calls.
	l := nvdevicelib.New(dgxa100.New())

	// Get the full list of allocatable devices from GPU 0 on the server.
	allocatable, err := l.GetPerGpuAllocatableDevices(0)
	if err != nil {
		klog.Fatalf("Error getAllocatableDevices: %v", err)
	}

	// Generate a resource slice spec from it:
	spec := (*resourceapi.PerGpuAllocatableDevices)(allocatable).ToResourceSliceSpec()

	// Flatten the original spec into a new one.
	spec, err = spec.Flatten()
	if err != nil {
		klog.Fatalf("Error spec.Flatten: %v", err)
	}

	// Print the list of possible allocations from the set of allocatable devices.
	if err := printPossibleAllocations(spec); err != nil {
		klog.Fatalf("Error printPossibleAllocations: %v", err)
	}
}

// printPossibleAllocations prints the possible allocations from the set of allocatable devices passed in.
func printPossibleAllocations(spec *resourceapi.ResourceSliceSpec) error {
	deviceCapacityMap := make(map[string]map[resourceapi.QualifiedName]resourceapi.DeviceCapacity)
	for _, d := range spec.Devices {
		if len(d.Partitionable.Capacity) > 0 {
			deviceCapacityMap[d.Name] = maps.Clone(d.Partitionable.Capacity)
		}
	}

	var err error
	Devices(spec.Devices).IterateCombinations(
		func(deviceCombo Devices) LoopControl {
			localDeviceCapacityMap := make(map[string]map[resourceapi.QualifiedName]resourceapi.DeviceCapacity)
			for _, device := range deviceCombo {
				if len(device.Partitionable.ConsumesCapacityFrom) == 0 {
					device.Partitionable.ConsumesCapacityFrom = append(device.Partitionable.ConsumesCapacityFrom, resourceapi.DeviceRef{device.Name})
				}
				for _, source := range device.Partitionable.ConsumesCapacityFrom {
					if _, exists := deviceCapacityMap[source.Name]; !exists {
						err = fmt.Errorf("device does not exist: %s", source)
						return Break
					}
					if _, exists := localDeviceCapacityMap[source.Name]; !exists {
						localDeviceCapacityMap[source.Name] = maps.Clone(deviceCapacityMap[source.Name])
					}
				}
			}
			for _, device := range deviceCombo {
				availableCapacity := make(map[resourceapi.QualifiedName]resourceapi.DeviceCapacity)
				for _, source := range device.Partitionable.ConsumesCapacityFrom {
					for k, v := range localDeviceCapacityMap[source.Name] {
						newc := availableCapacity[k]
						(&newc.Quantity).Add(v.Quantity)
						availableCapacity[k] = newc
					}
				}
				for k, v := range device.Partitionable.Capacity {
					if _, exists := availableCapacity[k]; !exists {
						fmt.Printf("device %+v, %v\n", device.Partitionable, availableCapacity)
						err = fmt.Errorf("missing capacity in sources: %s", k)
						return Break
					}
					if (&v.Quantity).Cmp(availableCapacity[k].Quantity) > 0 {
						return Break
					}
				}
				for k, v := range device.Partitionable.Capacity {
					remaining := v.Quantity
					for _, source := range device.Partitionable.ConsumesCapacityFrom {
						if _, exists := localDeviceCapacityMap[source.Name][k]; !exists {
							continue
						}
						if (&remaining).Cmp(localDeviceCapacityMap[source.Name][k].Quantity) > 0 {
							(&remaining).Sub(localDeviceCapacityMap[source.Name][k].Quantity)
							newc := localDeviceCapacityMap[source.Name][k]
							(&newc.Quantity).Reset()
							localDeviceCapacityMap[source.Name][k] = newc
						} else {
							newc := localDeviceCapacityMap[source.Name][k]
							(&newc.Quantity).Sub(v.Quantity)
							localDeviceCapacityMap[source.Name][k] = newc
						}
						if remaining.IsZero() {
							break
						}
					}
				}
			}
			fmt.Printf("%v\n", deviceCombo.GetNames())
			return Continue
		},
	)
	return err
}
