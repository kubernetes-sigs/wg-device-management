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

				// Generate attributes specific to this particular instance of a GPU.
				gpuSpecificAttributesName := toRFC1123Compliant(fmt.Sprintf("gpu-%d-attributes", d.Gpu.Index))
				mixinsMap[gpuSpecificAttributesName] = DeviceMixin{
					Name: gpuSpecificAttributesName,
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

				// Generate a mixin specific to this particular type of a GPU.
				specificGpuTypeName := toRFC1123Compliant(fmt.Sprintf("gpu-%s", d.Gpu.ProductName))
				mixinsMap[specificGpuTypeName] = DeviceMixin{
					Name: specificGpuTypeName,
					Partitionable: &PartitionableDeviceMixin{
						Includes: []DeviceMixinRef{
							{
								Name: systemAttributesName,
							},
							{
								Name: commonGpuAttributesName,
							},
						},
						Capacity: map[QualifiedName]DeviceCapacity{
							"memory": {
								Quantity: *resource.NewQuantity(int64(d.Gpu.MemoryBytes), resource.BinarySI),
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
				gpuSpecificMigAttributesName := toRFC1123Compliant(fmt.Sprintf("gpu-%d-mig-attributes", d.Mig.Parent.Index))
				mixinsMap[gpuSpecificMigAttributesName] = DeviceMixin{
					Name: gpuSpecificMigAttributesName,
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
				migProfileName := toRFC1123Compliant(fmt.Sprintf("mig-%s-%s", d.Mig.Profile, d.Mig.Parent.ProductName))
				mixinsMap[migProfileName] = DeviceMixin{
					Name: migProfileName,
					Partitionable: &PartitionableDeviceMixin{
						Includes: []DeviceMixinRef{
							{
								Name: systemAttributesName,
							},
							{
								Name: commonMigAttributesName,
							},
						},
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
				for k, v := range mixinsMap[migProfileName].Partitionable.Capacity {
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
								Name: gpuSpecificMigAttributesName,
							},
							{
								Name: migProfileName,
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
				specificGpuTypeName := toRFC1123Compliant(fmt.Sprintf("gpu-%s", d.Gpu.ProductName))
				gpuSpecificAttributesName := toRFC1123Compliant(fmt.Sprintf("gpu-%d-attributes", d.Gpu.Index))
				memorySlicesName := toRFC1123Compliant(fmt.Sprintf("memory-slices-%d-%d", 0, maxMemorySlice[d.Gpu.Index]))

				// Patch the specificGpuTypeName mixin with its capacities.
				for k, v := range maxCapacities[d.Gpu.Index] {
					mixinsMap[specificGpuTypeName].Partitionable.Capacity[k] = v
				}

				// Add each full GPU as a device in terms of its mixins.
				specificGpuName := toRFC1123Compliant(fmt.Sprintf("gpu-%d", d.Gpu.Index))
				devicesMap[specificGpuName] = Device{
					Name: specificGpuName,
					Partitionable: &PartitionableDevice{
						Includes: []DeviceMixinRef{
							{
								Name: specificGpuTypeName,
							},
							{
								Name: gpuSpecificAttributesName,
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
		visited := make(map[string]bool)
		stack := make(map[string]bool)

		deviceAsMixin := &DeviceMixin{
			Partitionable: &PartitionableDeviceMixin{
				Includes:   d.Partitionable.Includes,
				Attributes: d.Partitionable.Attributes,
				Capacity:   d.Partitionable.Capacity,
			},
		}

		flattened, err := flattenDeviceMixin(deviceAsMixin, mixinMap, visited, stack)
		if err != nil {
			return nil, err
		}

		d.Partitionable.Includes = flattened.Partitionable.Includes
		d.Partitionable.Attributes = flattened.Partitionable.Attributes
		d.Partitionable.Capacity = flattened.Partitionable.Capacity

		flattenedSpec.Devices = append(flattenedSpec.Devices, d)
	}

	return &flattenedSpec, nil
}

// flattenDeviceMixin flattens the mixin into a new PartitionableDeviceMixin, detecting errors along the way.
func flattenDeviceMixin(mixin *DeviceMixin, mixinMap map[string]*DeviceMixin, visited, stack map[string]bool) (*DeviceMixin, error) {
	// If the mixin is in the current recursion stack, a cycle is detected
	if stack[mixin.Name] {
		return nil, fmt.Errorf("cycle detected in mixin: %s", mixin.Name)
	}

	// Mark the mixin as visited and add it to the recursion stack
	visited[mixin.Name] = true
	stack[mixin.Name] = true

	// Create a new PartitionableDeviceMixin with merged attributes and capacities
	flattened := &DeviceMixin{
		Name: mixin.Name,
		Partitionable: &PartitionableDeviceMixin{
			Attributes: make(map[QualifiedName]DeviceAttribute),
			Capacity:   make(map[QualifiedName]DeviceCapacity),
		},
	}

	// If the mixin has a Partitionable definition, merge its attributes and capacities
	if mixin.Partitionable != nil {
		// Copy existing attributes and capacities to the flattened mixin
		for key, value := range mixin.Partitionable.Attributes {
			flattened.Partitionable.Attributes[key] = value
		}
		for key, value := range mixin.Partitionable.Capacity {
			flattened.Partitionable.Capacity[key] = value
		}

		// Traverse the included mixins and merge their attributes and capacities
		for _, includeRef := range mixin.Partitionable.Includes {
			includedMixin, exists := mixinMap[includeRef.Name]
			if !exists {
				return nil, fmt.Errorf("mixin %s not found", includeRef.Name)
			}

			// Recursively flatten the included mixin
			result, err := flattenDeviceMixin(includedMixin, mixinMap, visited, stack)
			if err != nil {
				return nil, err // Propagate any errors detected
			}

			// Merge attributes from the included mixin only if they do not already exist
			for key, value := range result.Partitionable.Attributes {
				if _, exists := flattened.Partitionable.Attributes[key]; !exists {
					flattened.Partitionable.Attributes[key] = value
				}
			}
			// Merge capacities from the included mixin only if they do not already exist
			for key, value := range result.Partitionable.Capacity {
				if _, exists := flattened.Partitionable.Capacity[key]; !exists {
					flattened.Partitionable.Capacity[key] = value
				}
			}
		}
	}

	// Remove the mixin from the recursion stack before backtracking
	stack[mixin.Name] = false
	return flattened, nil
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
