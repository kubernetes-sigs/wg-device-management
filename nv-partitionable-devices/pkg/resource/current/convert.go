package resource

import (
	"fmt"

	"github.com/Masterminds/semver"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"

	nvdevicelib "github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources/pkg/nvdevice"
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
	model := ResourceModel{
		NamedResources: &NamedResourcesResources{Instances: instances},
	}
	return model
}

// ToNamedResourcesResourceModel converts a list of AllocatableDevices to a NamedResources ResourceModel.
func (devices AllocatableDevices) ToNamedResourcesResourceModel() ResourceModel {
	instances := devices.ToNamedResourceInstances()
	model := ResourceModel{
		NamedResources: &NamedResourcesResources{Instances: instances},
	}
	return model
}

// ToNamedResourceInstances converts a list of PerGpuAllocatableDevices to a list of NamedResourcesInstances.
func (devices PerGpuAllocatableDevices) ToNamedResourceInstances() []NamedResourcesInstance {
	var instances []NamedResourcesInstance
	for _, perGpuDevices := range devices {
		instances = append(instances, AllocatableDevices(perGpuDevices).ToNamedResourceInstances()...)
	}
	return instances
}

// ToNamedResourceInstances converts a list of AllocatableDevices to a list of NamedResourcesInstances.
func (devices AllocatableDevices) ToNamedResourceInstances() []NamedResourcesInstance {
	var instances []NamedResourcesInstance
	for _, device := range devices {
		var instance *NamedResourcesInstance
		if device.Mig != nil {
			continue
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

// ToNamedResourcesInstance converts a GpuInfo object to a NamedResourcesInstance.
func (gpu *GpuInfo) ToNamedResourcesInstance() *NamedResourcesInstance {
	instance := &NamedResourcesInstance{
		Name: fmt.Sprintf("gpu-%v", gpu.Index),
		Attributes: []NamedResourcesAttribute{
			{
				Name: "minor",
				NamedResourcesAttributeValue: NamedResourcesAttributeValue{
					IntValue: ptr.To(int64(gpu.Minor)),
				},
			},
			{
				Name: "index",
				NamedResourcesAttributeValue: NamedResourcesAttributeValue{
					IntValue: ptr.To(int64(gpu.Index)),
				},
			},
			{
				Name: "uuid",
				NamedResourcesAttributeValue: NamedResourcesAttributeValue{
					StringValue: &gpu.UUID,
				},
			},
			{
				Name: "memory",
				NamedResourcesAttributeValue: NamedResourcesAttributeValue{
					QuantityValue: resource.NewQuantity(int64(gpu.MemoryBytes), resource.BinarySI),
				},
			},
			{
				Name: "product-name",
				NamedResourcesAttributeValue: NamedResourcesAttributeValue{
					StringValue: &gpu.ProductName,
				},
			},
			{
				Name: "brand",
				NamedResourcesAttributeValue: NamedResourcesAttributeValue{
					StringValue: &gpu.Brand,
				},
			},
			{
				Name: "architecture",
				NamedResourcesAttributeValue: NamedResourcesAttributeValue{
					StringValue: &gpu.Architecture,
				},
			},
			{
				Name: "cuda-compute-capability",
				NamedResourcesAttributeValue: NamedResourcesAttributeValue{
					VersionValue: ptr.To(semver.MustParse(gpu.CudaComputeCapability).String()),
				},
			},
			{
				Name: "driver-version",
				NamedResourcesAttributeValue: NamedResourcesAttributeValue{
					VersionValue: ptr.To(semver.MustParse(gpu.DriverVersion).String()),
				},
			},
			{
				Name: "cuda-driver-version",
				NamedResourcesAttributeValue: NamedResourcesAttributeValue{
					VersionValue: ptr.To(semver.MustParse(gpu.CudaDriverVersion).String()),
				},
			},
			{
				Name: "mig-capable",
				NamedResourcesAttributeValue: NamedResourcesAttributeValue{
					BoolValue: &gpu.MigCapable,
				},
			},
			{
				Name: "mig-enabled",
				NamedResourcesAttributeValue: NamedResourcesAttributeValue{
					BoolValue: &gpu.MigEnabled,
				},
			},
		},
	}

	return instance
}
