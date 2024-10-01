package resource

import (
	"fmt"
	"regexp"
	"strings"

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

// ToResourceSliceSpec converts a list of PerGpuAllocatableDevices to a ResourceSliceSpec
func (pgads PerGpuAllocatableDevices) ToResourceSliceSpec() *ResourceSliceSpec {
	mixinsMap := make(map[string]DeviceMixin)
	devicesMap := make(map[string]Device)

	// Generate a device mixin for all common system attributes.
	systemAttributesName := toRFC1123Compliant("system-attributes")
	mixinsMap[systemAttributesName] = DeviceMixin{
		Name: "system-attributes",
		Partitionable: &PartitionableDeviceMixin{
			Attributes: map[QualifiedName]DeviceAttribute{
				"driverVersion": {
					VersionValue: ptr.To(nvdevicelib.PerGpuAllocatableDevices(pgads).SystemInfo.DriverVersion),
				},
				"cudaDriverVersion": {
					VersionValue: ptr.To(nvdevicelib.PerGpuAllocatableDevices(pgads).SystemInfo.CudaDriverVersion),
				},
			},
		},
	}

	// Track the biggest memory slice discovered per GPU so we can include that
	// in the device specs we eventually generate.
	maxMemorySlice := make(map[int]uint32)

	// Track the maximum capacity of all MIG resources (on a per gpu basis) so
	// that we can apply them to each full GPU later on.
	maxCapacities := make(map[int]map[QualifiedName]DeviceCapacity)

	// Loop through all discovered devices and generate mixins and devices from them.
	for i := range pgads.Devices {
		for _, d := range pgads.Devices[i] {
			// If this is a full GPU ...
			if d.Gpu != nil {
				// Generate common attributes for all instances of the current GPU type.
				commonGpuAttributesName := toRFC1123Compliant(fmt.Sprintf("common-gpu-%s-attributes", d.Gpu.ProductName))
				mixinsMap[commonGpuAttributesName] = DeviceMixin{
					Name: commonGpuAttributesName,
					Partitionable: &PartitionableDeviceMixin{
						Attributes: map[QualifiedName]DeviceAttribute{
							"type": {
								StringValue: ptr.To("gpu"),
							},
							"architecture": {
								StringValue: ptr.To(d.Gpu.Architecture),
							},
							"brand": {
								StringValue: ptr.To(d.Gpu.Brand),
							},
							"productName": {
								StringValue: ptr.To(d.Gpu.ProductName),
							},
							"cudaComputeCapability": {
								StringValue: ptr.To(d.Gpu.CudaComputeCapability),
							},
						},
					},
				}

				// Generate common capacities for all instance of the current GPU type.
				// NOTE: These will be patched up after computing the max GMI capacities required.
				commonGpuCapacitiesName := toRFC1123Compliant(fmt.Sprintf("common-gpu-%s-capacities", d.Gpu.ProductName))
				mixinsMap[commonGpuCapacitiesName] = DeviceMixin{
					Name: commonGpuCapacitiesName,
					Partitionable: &PartitionableDeviceMixin{
						Capacity: map[QualifiedName]DeviceCapacity{
							"memory": {
								Quantity: *resource.NewQuantity(int64(d.Gpu.MemoryBytes), resource.BinarySI),
							},
						},
					},
				}

				// Generate attributes specific to this particular instance of a GPU.
				specificGpuAttributesName := toRFC1123Compliant(fmt.Sprintf("specific-gpu-%d-attributes", d.Gpu.Index))
				mixinsMap[specificGpuAttributesName] = DeviceMixin{
					Name: specificGpuAttributesName,
					Partitionable: &PartitionableDeviceMixin{
						Attributes: map[QualifiedName]DeviceAttribute{
							"index": {
								IntValue: ptr.To(int64(d.Gpu.Index)),
							},
							"minor": {
								IntValue: ptr.To(int64(d.Gpu.Minor)),
							},
							"uuid": {
								StringValue: ptr.To(d.Gpu.UUID),
							},
						},
					},
				}
			}

			// If this is a MIG device ...
			if d.Mig != nil {
				// Generate common attributes for all MIG instances of the current GPU type.
				commonMigAttributesName := toRFC1123Compliant(fmt.Sprintf("common-mig-%s-attributes", d.Mig.Parent.ProductName))
				mixinsMap[commonMigAttributesName] = DeviceMixin{
					Name: commonMigAttributesName,
					Partitionable: &PartitionableDeviceMixin{
						Attributes: map[QualifiedName]DeviceAttribute{
							"type": {
								StringValue: ptr.To("mig"),
							},
							"architecture": {
								StringValue: ptr.To(d.Mig.Parent.Architecture),
							},
							"brand": {
								StringValue: ptr.To(d.Mig.Parent.Brand),
							},
							"productName": {
								StringValue: ptr.To(d.Mig.Parent.ProductName),
							},
							"cudaComputeCapability": {
								StringValue: ptr.To(d.Mig.Parent.CudaComputeCapability),
							},
						},
					},
				}

				// Generate MIG attributes specific to this particular instance of a GPU.
				specificGpuMigAttributesName := toRFC1123Compliant(fmt.Sprintf("specific-gpu-%d-mig-attributes", d.Mig.Parent.Index))
				mixinsMap[specificGpuMigAttributesName] = DeviceMixin{
					Name: specificGpuMigAttributesName,
					Partitionable: &PartitionableDeviceMixin{
						Attributes: map[QualifiedName]DeviceAttribute{
							"parentIndex": {
								IntValue: ptr.To(int64(d.Mig.Parent.Index)),
							},
							"parentMinor": {
								IntValue: ptr.To(int64(d.Mig.Parent.Minor)),
							},
							"parentUUID": {
								StringValue: ptr.To(d.Mig.Parent.UUID),
							},
						},
					},
				}

				// Generate a mixin for current the MIG profile specific to this particular type of GPU.
				info := d.Mig.GIProfileInfo
				commonMigMixinName := toRFC1123Compliant(fmt.Sprintf("common-mig-%s-%s", d.Mig.Profile, d.Mig.Parent.ProductName))
				mixinsMap[commonMigMixinName] = DeviceMixin{
					Name: commonMigMixinName,
					Partitionable: &PartitionableDeviceMixin{
						Attributes: map[QualifiedName]DeviceAttribute{
							"profile": {
								StringValue: ptr.To(d.Mig.Profile.String()),
							},
						},
						Capacity: map[QualifiedName]DeviceCapacity{
							"multiprocessors": {
								Quantity: *resource.NewQuantity(int64(info.MultiprocessorCount), resource.BinarySI),
							},
							"copy-engines": {
								Quantity: *resource.NewQuantity(int64(info.CopyEngineCount), resource.BinarySI),
							},
							"decoders": {
								Quantity: *resource.NewQuantity(int64(info.DecoderCount), resource.BinarySI),
							},
							"encoders": {
								Quantity: *resource.NewQuantity(int64(info.EncoderCount), resource.BinarySI),
							},
							"jpeg-engines": {
								Quantity: *resource.NewQuantity(int64(info.JpegCount), resource.BinarySI),
							},
							"ofa-engines": {
								Quantity: *resource.NewQuantity(int64(info.OfaCount), resource.BinarySI),
							},
							"memory": {
								Quantity: *resource.NewQuantity(int64(info.MemorySizeMB*1024*1024), resource.BinarySI),
							},
						},
					},
				}

				// Track the maxCapacities of all capacity consumed by any MIG
				// device on the current full GPU so that we can apply these
				// max capacities to the full GPU later on.
				if maxCapacities[d.Mig.Parent.Index] == nil {
					maxCapacities[d.Mig.Parent.Index] = make(map[QualifiedName]DeviceCapacity)
				}
				for k, v := range mixinsMap[commonMigMixinName].Partitionable.Capacity {
					if k == "memory" {
						continue
					}
					if ptr.To(maxCapacities[d.Mig.Parent.Index][k].Quantity).Cmp(v.Quantity) <= 0 {
						maxCapacities[d.Mig.Parent.Index][k] = DeviceCapacity{
							Quantity: v.Quantity,
						}
					}
				}

				// Generate a mixin to represent all of the memory slices consumed by the current MIG device.
				placement := d.Mig.MemorySlices
				memorySlicesSuffix := fmt.Sprintf("%d", placement.Start)
				if placement.Size > 1 {
					memorySlicesSuffix = fmt.Sprintf("%s-%d", memorySlicesSuffix, placement.Start+placement.Size-1)
				}
				memorySlicesName := toRFC1123Compliant(fmt.Sprintf("memory-slices-%s", memorySlicesSuffix))
				mixinsMap[memorySlicesName] = DeviceMixin{
					Name: memorySlicesName,
					Partitionable: &PartitionableDeviceMixin{
						Capacity: map[QualifiedName]DeviceCapacity{},
					},
				}
				for i := placement.Start; i < placement.Start+placement.Size; i++ {
					sliceName := QualifiedName(fmt.Sprintf("memorySlice%d", i))
					mixinsMap[memorySlicesName].Partitionable.Capacity[sliceName] = DeviceCapacity{
						Quantity: *resource.NewQuantity(1, resource.BinarySI),
					}
				}

				// Track the max memory slice consumed by any MIG device on the
				// current full GPU so that we can apply it to the full GPU
				// later on.
				maxMemorySlice[d.Mig.Parent.Index] = max(maxMemorySlice[d.Mig.Parent.Index], placement.Start+placement.Size-1)

				// Generate the actual MIG device spec (not a mixin).
				migDeviceName := toRFC1123Compliant(fmt.Sprintf("gpu-%d-mig-%s-%s", d.Mig.Parent.Index, d.Mig.Profile, memorySlicesSuffix))
				devicesMap[migDeviceName] = Device{
					Name: migDeviceName,
					Partitionable: &PartitionableDevice{
						Includes: []DeviceMixinRef{
							{
								Name: systemAttributesName,
							},
							{
								Name: commonMigAttributesName,
							},
							{
								Name: commonMigMixinName,
							},
							{
								Name: specificGpuMigAttributesName,
							},
							{
								Name: memorySlicesName,
							},
						},
						ConsumesCapacityFrom: []DeviceRef{
							{
								Name: fmt.Sprintf("gpu-%d", d.Mig.Parent.Index),
							},
						},
					},
				}
			}
		}
	}

	// Loop through all full GPus again to add device specs for them.
	// We need to do this in a second loop to make use of the maxMemorySlice and
	// maxCapacities objects constructed in the first loop.
	for i := range pgads.Devices {
		for _, d := range pgads.Devices[i] {
			if d.Gpu != nil {
				// Recreate the names of the mixins we need to include in each concrete device.
				commonGpuAttributesName := toRFC1123Compliant(fmt.Sprintf("common-gpu-%s-attributes", d.Gpu.ProductName))
				commonGpuCapacitiesName := toRFC1123Compliant(fmt.Sprintf("common-gpu-%s-capacities", d.Gpu.ProductName))
				specificGpuAttributesName := toRFC1123Compliant(fmt.Sprintf("specific-gpu-%d-attributes", d.Gpu.Index))
				memorySlicesName := toRFC1123Compliant(fmt.Sprintf("memory-slices-%d-%d", 0, maxMemorySlice[d.Gpu.Index]))

				// Patch the specificGpuTypeName mixin with its capacities.
				for k, v := range maxCapacities[d.Gpu.Index] {
					mixinsMap[commonGpuCapacitiesName].Partitionable.Capacity[k] = v
				}

				// Add each full GPU as a device in terms of its mixins.
				specificGpuName := toRFC1123Compliant(fmt.Sprintf("gpu-%d", d.Gpu.Index))
				devicesMap[specificGpuName] = Device{
					Name: specificGpuName,
					Partitionable: &PartitionableDevice{
						Includes: []DeviceMixinRef{
							{
								Name: systemAttributesName,
							},
							{
								Name: commonGpuAttributesName,
							},
							{
								Name: commonGpuCapacitiesName,
							},
							{
								Name: specificGpuAttributesName,
							},
							{
								Name: memorySlicesName,
							},
						},
					},
				}
			}
		}
	}

	var mixins []DeviceMixin
	for _, v := range mixinsMap {
		mixins = append(mixins, v)
	}

	var devices []Device
	for _, v := range devicesMap {
		devices = append(devices, v)
	}

	spec := ResourceSliceSpec{
		DeviceMixins: mixins,
		Devices:      devices,
	}

	return &spec
}

// Flatten unrolls a spec into a list of flat devices resolving any mixins into
// inline attributes and capacities.
func (s *ResourceSliceSpec) Flatten() (*ResourceSliceSpec, error) {
	// Declare a new resource slice spec to hold the flattened result.
	var flattenedSpec ResourceSliceSpec

	// Build a map of all the mixins by name. Return an error if any names are duplicated.
	mixinMap := make(map[string]*DeviceMixin)
	for _, m := range s.DeviceMixins {
		if _, exists := mixinMap[m.Name]; exists {
			return nil, fmt.Errorf("duplicate mixin name detected: %s", m.Name)
		}
		mixinMap[m.Name] = &m
	}

	// Flatten each device and add it back to the flattened spec.
	for _, d := range s.Devices {
		flattened := &DeviceMixin{
			Partitionable: &PartitionableDeviceMixin{
				Attributes: make(map[QualifiedName]DeviceAttribute),
				Capacity:   make(map[QualifiedName]DeviceCapacity),
			},
		}
		for _, m := range d.Partitionable.Includes {
			for k, v := range mixinMap[m.Name].Partitionable.Attributes {
				flattened.Partitionable.Attributes[k] = v
			}
			for k, v := range mixinMap[m.Name].Partitionable.Capacity {
				flattened.Partitionable.Capacity[k] = v
			}
		}
		for k, v := range d.Partitionable.Attributes {
			flattened.Partitionable.Attributes[k] = v
		}
		for k, v := range d.Partitionable.Capacity {
			flattened.Partitionable.Capacity[k] = v
		}

		d.Partitionable.Includes = nil
		d.Partitionable.Attributes = flattened.Partitionable.Attributes
		d.Partitionable.Capacity = flattened.Partitionable.Capacity

		flattenedSpec.Devices = append(flattenedSpec.Devices, d)
	}

	return &flattenedSpec, nil
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
