package main

import (
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml/mock/dgxa100"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"

	nvdevicelib "github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/nvdevice"
	resourceapi "github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/resource"
)

// Main queries the list of allocatable devices and prints them as a kubernetes
// resource slice.
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

	// Print the original resource slice spec.
	fmt.Printf("Original spec:\n")
	if err := printResourceSliceSpec(spec); err != nil {
		klog.Fatalf("Error printResourceSliceSpec: %v", err)
	}
	fmt.Printf("\n")

	// Flatten the original spec into a new one.
	spec, err = spec.Flatten()
	if err != nil {
		klog.Fatalf("Error spec.Flatten: %v", err)
	}

	// Print the flattened resource slice spec.
	fmt.Printf("Flattened spec:\n")
	if err := printResourceSliceSpec(spec); err != nil {
		klog.Fatalf("Error printResourceSliceSpec: %v", err)
	}
}

// printResourcesSliceSpec prints the resource slice spec as yaml.
func printResourceSliceSpec(spec *resourceapi.ResourceSliceSpec) error {
	specYaml, err := yaml.Marshal(spec)
	if err != nil {
		klog.Fatalf("Error marshaling resource spec spec to yaml: %v", err)
	}
	fmt.Printf("%v", string(specYaml))
	return nil
}
