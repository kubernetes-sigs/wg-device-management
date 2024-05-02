package resource

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"

	"github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/intrange"
	nvdevicelib "github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/nvdevice"
	currentresourceapi "github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/resource/current"
)

// PerGpuAllocatableDevices is an alias of nvdevicelib.PerGpuAllocatableDevices
type PerGpuAllocatableDevices nvdevicelib.PerGpuAllocatableDevices

// AllocatableDevices is an alias of nvdevicelib.AllocatableDevices
type AllocatableDevices nvdevicelib.AllocatableDevices

// GpuInfo is an alias of nvdevicelib.GpuInfo
type GpuInfo nvdevicelib.GpuInfo

// MigInfo is an alias of nvdevicelib.MigInfo
type MigInfo nvdevicelib.MigInfo

// ToNamedResourcesResourceModel converts a list of PerGpuAllocatableDevices to a NamedResources ResourceModel.
func (devices PerGpuAllocatableDevices) ToNamedResourcesResourceModel() ResourceModel {
	instances := devices.ToNamedResourceInstances()
	sharedLimits := devices.ToSharedLimits()
	model := ResourceModel{
		NamedResources: &NamedResourcesResources{
			Instances:    instances,
			SharedLimits: sharedLimits,
		},
	}
	return model
}

// ToNamedResourcesResourceModel converts a list of AllocatableDevices to a NamedResources ResourceModel.
func (devices AllocatableDevices) ToNamedResourcesResourceModel() ResourceModel {
	instances := devices.ToNamedResourceInstances()
	sharedLimits := devices.ToSharedLimits()
	model := ResourceModel{
		NamedResources: &NamedResourcesResources{
			Instances:    instances,
			SharedLimits: []NamedResourcesSharedResourceGroup{sharedLimits},
		},
	}
	return model
}

// ToNamedResourceInstances converts a list of PerGpuAllocatableDevices to a list of NamedResourcesInstances.
func (devices PerGpuAllocatableDevices) ToNamedResourceInstances() []NamedResourcesInstance {
	var instances []NamedResourcesInstance
	for _, perGpuDevices := range devices {
		instances = slices.Concat(instances, AllocatableDevices(perGpuDevices).ToNamedResourceInstances())
	}
	return instances
}

// ToNamedResourceInstances converts a list of AllocatableDevices to a list of NamedResourcesInstances.
func (devices AllocatableDevices) ToNamedResourceInstances() []NamedResourcesInstance {
	var instances []NamedResourcesInstance
	for _, device := range devices {
		var instance *NamedResourcesInstance
		if device.Mig != nil {
			instance = (*MigInfo)(device.Mig).ToNamedResourcesInstance()
		}
		if device.Gpu != nil {
			instance = (*GpuInfo)(device.Gpu).ToNamedResourcesInstance()
		}
		if instance != nil {
			instances = append(instances, *instance)
		}
	}
	return instances
}

// ToSharedLimits converts a list of PerGpuAllocatableDevices to a list of NamedResourcesSharedResourceGroups shared resource limits.
func (devices PerGpuAllocatableDevices) ToSharedLimits() []NamedResourcesSharedResourceGroup {
	var limits []NamedResourcesSharedResourceGroup
	for _, perGpuDevices := range devices {
		limits = append(limits, AllocatableDevices(perGpuDevices).ToSharedLimits())
	}
	return limits
}

// ToSharedLimits converts a list of AllocatableDevices to a NamedResourcesSharedResourceGroup of shared resource limits.
func (devices AllocatableDevices) ToSharedLimits() NamedResourcesSharedResourceGroup {
	var limits NamedResourcesSharedResourceGroup
	for _, device := range devices {
		var resources NamedResourcesSharedResourceGroup
		if device.Gpu != nil {
			resources = (*GpuInfo)(device.Gpu).getResources()
		}
		if device.Mig != nil {
			resources = (*MigInfo)(device.Mig).getResources()
		}
		if len(limits.Name) == 0 {
			limits.Name = resources.Name
		}
		for _, item := range resources.Items {
			if item.QuantityValue != nil {
				limits.addOrReplaceQuantityIfLarger(&item)
			}
			if item.IntRangeValue != nil {
				limits.addOrReplaceIntRangeIfLarger(&item)
			}
		}
	}
	return limits
}

// ToNamedResourcesInstance converts a GpuInfo object to a NamedResourcesInstance.
func (gpu *GpuInfo) ToNamedResourcesInstance() *NamedResourcesInstance {
	currentInstance := (*currentresourceapi.GpuInfo)(gpu).ToNamedResourcesInstance()

	var attributes []NamedResourcesAttribute
	for _, attribute := range currentInstance.Attributes {
		switch attribute.Name {
		case "memory":
			break
		default:
			attributes = append(attributes, attribute)
		}
	}

	resources := []NamedResourcesSharedResourceGroup{
		(*GpuInfo)(gpu).getResources(),
	}

	newInstance := &NamedResourcesInstance{
		Name:       currentInstance.Name,
		Attributes: attributes,
		Resources:  resources,
	}

	return newInstance
}

// ToNamedResourcesInstance converts a MigInfo object to a NamedResourcesInstance.
func (mig *MigInfo) ToNamedResourcesInstance() *NamedResourcesInstance {
	parentInstance := (*currentresourceapi.GpuInfo)(mig.Parent).ToNamedResourcesInstance()

	attributes := []NamedResourcesAttribute{
		{
			Name: "mig-profile",
			NamedResourcesAttributeValue: NamedResourcesAttributeValue{
				StringValue: ptr.To(mig.Profile.String()),
			},
		},
	}
	for _, attribute := range parentInstance.Attributes {
		switch attribute.Name {
		case "product-name", "brand", "architecture", "cuda-compute-capability", "driver-version", "cuda-driver-version":
			attributes = append(attributes, attribute)
		}
	}

	resources := []NamedResourcesSharedResourceGroup{
		(*MigInfo)(mig).getResources(),
	}

	name := fmt.Sprintf("%s-mig-%s-%d", parentInstance.Name, mig.Profile, mig.MemorySlices.Start)

	instance := &NamedResourcesInstance{
		Name:       toRFC1123Compliant(name),
		Attributes: attributes,
		Resources:  resources,
	}

	return instance
}

// getResources gets the set of shared resources consumed by the GPU.
func (gpu *GpuInfo) getResources() NamedResourcesSharedResourceGroup {
	name := fmt.Sprintf("gpu-%v-shared-resources", gpu.Index)

	resources := []NamedResourcesSharedResource{
		{
			Name: "memory",
			NamedResourcesSharedResourceValue: NamedResourcesSharedResourceValue{
				QuantityValue: resource.NewQuantity(int64(gpu.MemoryBytes), resource.BinarySI),
			},
		},
	}

	group := NamedResourcesSharedResourceGroup{
		Name:  name,
		Items: resources,
	}

	return group
}

// getResources gets the set of shared resources consumed by the MIG device.
func (mig *MigInfo) getResources() NamedResourcesSharedResourceGroup {
	name := fmt.Sprintf("gpu-%v-shared-resources", mig.Parent.Index)

	info := mig.GIProfileInfo
	resources := []NamedResourcesSharedResource{
		{
			Name: "multiprocessors",
			NamedResourcesSharedResourceValue: NamedResourcesSharedResourceValue{
				QuantityValue: resource.NewQuantity(int64(info.MultiprocessorCount), resource.BinarySI),
			},
		},
		{
			Name: "copy-engines",
			NamedResourcesSharedResourceValue: NamedResourcesSharedResourceValue{
				QuantityValue: resource.NewQuantity(int64(info.CopyEngineCount), resource.BinarySI),
			},
		},
		{
			Name: "decoders",
			NamedResourcesSharedResourceValue: NamedResourcesSharedResourceValue{
				QuantityValue: resource.NewQuantity(int64(info.DecoderCount), resource.BinarySI),
			},
		},
		{
			Name: "encoders",
			NamedResourcesSharedResourceValue: NamedResourcesSharedResourceValue{
				QuantityValue: resource.NewQuantity(int64(info.EncoderCount), resource.BinarySI),
			},
		},
		{
			Name: "jpeg-engines",
			NamedResourcesSharedResourceValue: NamedResourcesSharedResourceValue{
				QuantityValue: resource.NewQuantity(int64(info.JpegCount), resource.BinarySI),
			},
		},
		{
			Name: "ofa-engines",
			NamedResourcesSharedResourceValue: NamedResourcesSharedResourceValue{
				QuantityValue: resource.NewQuantity(int64(info.OfaCount), resource.BinarySI),
			},
		},
		{
			Name: "memory",
			NamedResourcesSharedResourceValue: NamedResourcesSharedResourceValue{
				QuantityValue: resource.NewQuantity(int64(info.MemorySizeMB*1024*1024), resource.BinarySI),
			},
		},
		{
			Name: "memory-slices",
			NamedResourcesSharedResourceValue: NamedResourcesSharedResourceValue{
				IntRangeValue: intrange.NewIntRange(int64(mig.MemorySlices.Start), int64(mig.MemorySlices.Size)),
			},
		},
	}

	group := NamedResourcesSharedResourceGroup{
		Name:  name,
		Items: resources,
	}

	return group
}

// toRFC1123Compliant converts the incoming string to a valid RFC1123 DNS domain name.
func toRFC1123Compliant(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace invalid characters with '-'
	re := regexp.MustCompile(`[^a-z0-9-.]`)
	name = re.ReplaceAllString(name, "-")

	// Trim leading/trailing '-'
	name = strings.Trim(name, "-")

	// Trim trailing '.'
	name = strings.TrimSuffix(name, ".")

	// Truncate to 253 characters
	if len(name) > 253 {
		name = name[:253]
	}

	return name
}
