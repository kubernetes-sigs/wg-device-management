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

func dgxa100Pool(nodeName, poolName string) (*api.DevicePool, error) {
	// Instantiate an instance of a mock dgxa100 server and build a nvDeviceLib
	// from it. The nvDeviceLib is then used to populate the list of allocatable
	// devices from this mock server using standard NVML calls.
	l := nvdevicelib.New(dgxa100.New())

	// Get the full list of allocatable devices from GPU 0 on the server.
	allocatable, err := l.GetPerGpuAllocatableDevices(0)
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

	var devices []api.Device
	for _, instance := range model.NamedResources.Instances {
		devices = append(devices, instanceToDevice(instance))
	}

	return &api.DevicePool{
		TypeMeta: metav1.TypeMeta{
			APIVersion: api.DevMgmtAPIVersion,
			Kind:       "DevicePool",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: nodeName + "-" + poolName,
		},
		Spec: api.DevicePoolSpec{
			NodeName: &nodeName,
			Driver:   "gpu.nvidia.com/dra",
			Attributes: []api.Attribute{
				{Name: "vendor", StringValue: ptr("nvidia")},
				{Name: "model", StringValue: ptr("dgxa100")},
			},
			SharedResources: sharedGroupToResources(model.NamedResources.SharedLimits[0], false),
			Devices:         devices,
		},
	}, nil
}

func instanceToDevice(instance newresourceapi.NamedResourcesInstance) api.Device {
	device := api.Device{
		Name:       instance.Name,
		Attributes: attributesToAttributes(instance.Attributes),
	}

	if len(instance.Resources) > 0 {
		device.SharedResourcesConsumed = sharedGroupToRequests(instance.Resources[0])
		device.ClaimResourcesProvided = sharedGroupToResources(instance.Resources[0], true)
	}

	return device
}

func attributesToAttributes(attrs []resourceapi.NamedResourcesAttribute) []api.Attribute {
	var attributes []api.Attribute

	for _, a := range attrs {
		if a.QuantityValue != nil {
			attributes = append(attributes, api.Attribute{
				Name:          a.Name,
				QuantityValue: a.QuantityValue,
			})
		} else if a.StringValue != nil {
			attributes = append(attributes, api.Attribute{
				Name:        a.Name,
				StringValue: a.StringValue,
			})
		} else if a.IntValue != nil {
			attributes = append(attributes, api.Attribute{
				Name:     a.Name,
				IntValue: ptr(int(*a.IntValue)),
			})
		} else if a.VersionValue != nil {
			attributes = append(attributes, api.Attribute{
				Name:        a.Name,
				SemVerValue: ptr(api.SemVer(*a.VersionValue)),
			})
		}
		// don't convert ones not supported in the prototype
	}

	return attributes
}

func sharedGroupToResources(group newresourceapi.NamedResourcesSharedResourceGroup, userOnly bool) []api.ResourceCapacity {
	var resources []api.ResourceCapacity

	for _, item := range group.Items {
		// as an example, make it so that only memory is a user-facing
		// allocatable resource, whereas other resources are just
		// about the shared pool
		if userOnly && item.Name != "memory" {
			continue
		}
		if item.QuantityValue != nil {
			resources = append(resources, api.ResourceCapacity{
				Name:     item.Name,
				Capacity: *item.QuantityValue,
			})
		} else if item.IntRangeValue != nil {
			// sorry, unrolling these intranges to avoid additional types
			// beyond Quantity. Assumes max range 0-7
			for i := 0; i < 8; i++ {
				single := intrange.NewIntRange(int64(i), 1)
				if item.IntRangeValue.Contains(single) {
					resources = append(resources, api.ResourceCapacity{
						Name:     fmt.Sprintf("%s-%d", item.Name, i),
						Capacity: resource.MustParse("1"),
					})
				}
			}
		}
	}

	return resources
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
