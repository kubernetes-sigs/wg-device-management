package main

import (
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml/mock/dgxa100"
	"k8s.io/klog/v2"

	nvdevicelib "github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/nvdevice"
	newresourceapi "github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/resource/new"
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

	// Print the list of possible allocations with the new model.
	if err := printNewResourceModelPossibleAllocations(allocatable); err != nil {
		klog.Fatalf("Error printNewResourceModel: %v", err)
	}
}

// printNewResourceModelPossibleAllocations prints the possible allocations of the new structured resource model.
func printNewResourceModelPossibleAllocations(allocatable nvdevicelib.PerGpuAllocatableDevices) error {
	model := newresourceapi.PerGpuAllocatableDevices(allocatable).ToNamedResourcesResourceModel()
	sharedLimits := model.NamedResources.SharedLimits
	instances := model.NamedResources.Instances

	NamedResourcesInstances(instances).IterateCombinations(
		func(instanceCombo NamedResourcesInstances) LoopControl {
			for _, limits := range sharedLimits {
				limitsCopy := limits.DeepCopy()
				for _, instance := range instanceCombo {
					for _, resources := range instance.Resources {
						if limits.Name != resources.Name {
							continue
						}
						success, err := limitsCopy.Sub(&resources)
						if err != nil {
							klog.Errorf("error subtracting resources from %+v: %v\n", instance, err)
						}
						if !success {
							return Break
						}
					}
				}
			}
			fmt.Printf("%v\n", instanceCombo.GetNames())
			return Continue
		},
	)
	return nil
}
