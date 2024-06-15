package gen

import (
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml/mock/dgxa100"
	"github.com/kubernetes-sigs/wg-device-management/dra-evolution/pkg/api"
	"github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/intrange"
	nvdevicelib "github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/nvdevice"
	newresourceapi "github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/resource/new"

	resourceapi "k8s.io/api/resource/v1alpha2"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func dgxa100Pool(nodeName, poolName string, gpus int) (*api.ResourcePool, error) {
	// Instantiate an instance of a mock dgxa100 server and build a nvDeviceLib
	// from it. The nvDeviceLib is then used to populate the list of allocatable
	// devices from this mock server using standard NVML calls.
	l := nvdevicelib.New(dgxa100.New())

	var devices []api.Device
	var shared []api.SharedCapacityGroup
	for gpu := 0; gpu < gpus; gpu++ {
		// Get the full list of allocatable devices from this GPU on the server
		allocatable, err := l.GetPerGpuAllocatableDevices(gpu)
		if err != nil {
			return nil, err
		}

		model := newresourceapi.PerGpuAllocatableDevices(allocatable).ToNamedResourcesResourceModel()

		if model.NamedResources == nil {
			return nil, fmt.Errorf("found no named resource object in the model")
		}

		if len(model.NamedResources.SharedLimits) != 1 {
			return nil, fmt.Errorf("found %d shared limit groups in the resources", len(model.NamedResources.SharedLimits))
		}

		shared = append(shared, sharedGroupToResources(model.NamedResources.SharedLimits[0], gpu))
		for _, instance := range model.NamedResources.Instances {
			devices = append(devices, instanceToDevice(instance, gpu))
		}
	}

	return &api.ResourcePool{
		TypeMeta: metav1.TypeMeta{
			APIVersion: DevMgmtAPIVersion,
			Kind:       "ResourcePool",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: nodeName + "-" + poolName,
		},
		Spec: api.ResourcePoolSpec{
			NodeName:       nodeName,
			DriverName:     "gpu.nvidia.com/dra",
			SharedCapacity: shared,
			Devices:        devices,
		},
	}, nil
}

func instanceToDevice(instance newresourceapi.NamedResourcesInstance, gpu int) api.Device {
	device := api.Device{
		Name:       instance.Name,
		Attributes: attributesToDeviceAttributes(instance.Attributes),
	}

	if len(instance.Resources) > 0 {
		device.SharedCapacityConsumed = []api.SharedCapacityGroup{sharedGroupToResources(instance.Resources[0], gpu)}
	}

	return device
}

func attributesToDeviceAttributes(attrs []resourceapi.NamedResourcesAttribute) []api.DeviceAttribute {
	var attributes []api.DeviceAttribute

	for _, a := range attrs {
		if a.QuantityValue != nil {
			attributes = append(attributes, api.DeviceAttribute{
				Name:          a.Name,
				QuantityValue: a.QuantityValue,
			})
		} else if a.BoolValue != nil {
			attributes = append(attributes, api.DeviceAttribute{
				Name:      a.Name,
				BoolValue: a.BoolValue,
			})
		} else if a.StringValue != nil {
			attributes = append(attributes, api.DeviceAttribute{
				Name:        a.Name,
				StringValue: a.StringValue,
			})
		} else if a.VersionValue != nil {
			attributes = append(attributes, api.DeviceAttribute{
				Name:         a.Name,
				VersionValue: a.VersionValue,
			})
		}
		// don't convert ones not supported in the prototype
	}

	return attributes
}

func sharedGroupToResources(group newresourceapi.NamedResourcesSharedResourceGroup, gpu int) api.SharedCapacityGroup {
	var newGroup api.SharedCapacityGroup

	newGroup.Name = group.Name
	for _, item := range group.Items {
		if item.QuantityValue != nil && !item.QuantityValue.IsZero() {
			newGroup.SharedCapacity = append(newGroup.SharedCapacity, api.SharedCapacity{
				Name:     item.Name,
				Capacity: *item.QuantityValue,
			})
		} else if item.IntRangeValue != nil {
			// sorry, unrolling these intranges to avoid additional types
			// beyond Quantity. Assumes max range 0-7
			for i := 0; i < 8; i++ {
				single := intrange.NewIntRange(int64(i), 1)
				if item.IntRangeValue.Contains(single) {
					newGroup.SharedCapacity = append(newGroup.SharedCapacity, api.SharedCapacity{
						Name:     fmt.Sprintf("%s-%d", item.Name, i),
						Capacity: resource.MustParse("1"),
					})
				}
			}
		}
	}

	return newGroup
}

func sharedGroupToRequests(group newresourceapi.NamedResourcesSharedResourceGroup) map[string]resource.Quantity {
	requests := make(map[string]resource.Quantity)

	for _, item := range group.Items {
		if item.QuantityValue != nil {
			if item.QuantityValue.IsZero() {
				continue
			}
			requests[item.Name] = *item.QuantityValue
		} else if item.IntRangeValue != nil {
			// sorry, unrolling these intranges to avoid additional types
			// beyond Quantity. Assumes max range 0-15
			for i := 0; i < 16; i++ {
				single := intrange.NewIntRange(int64(i), 1)
				if item.IntRangeValue.Contains(single) {
					requests[fmt.Sprintf("%s-%02d", item.Name, i)] = resource.MustParse("1")
				}
			}
		}
	}

	return requests
}
