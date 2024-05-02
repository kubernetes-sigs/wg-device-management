package nvdevice

import (
	"fmt"
	"slices"

	"k8s.io/klog/v2"

	nvdev "github.com/NVIDIA/go-nvlib/pkg/nvlib/device"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

// PerGpuAllocatableDevices holds the list of allocatable devices per GPU.
type PerGpuAllocatableDevices map[int]AllocatableDevices

// AllocatableDevices holds the list of allocatable devices.
type AllocatableDevices []AllocatableDevice

// AllocatableDevice represents an individual device that can be allocated.
// This can either be a full GPU or MIG device, but not both.
type AllocatableDevice struct {
	Gpu *GpuInfo
	Mig *MigInfo
}

// GpuInfo holds all of the relevant information about a GPU.
type GpuInfo struct {
	Minor                 int
	Index                 int
	UUID                  string
	MemoryBytes           uint64
	ProductName           string
	Brand                 string
	Architecture          string
	CudaComputeCapability string
	DriverVersion         string
	CudaDriverVersion     string
	MigCapable            bool
	MigEnabled            bool
}

// MigInfo holds all of the relevant information about a MIG device.
type MigInfo struct {
	Parent        *GpuInfo
	Profile       nvdev.MigProfile
	GIProfileInfo nvml.GpuInstanceProfileInfo
	MemorySlices  nvml.GpuInstancePlacement
}

// NVDeviceLib encapsulates the set of libraries and methods required to query
// info about an NVIDIA device.
type NVDeviceLib struct {
	nvdev nvdev.Interface
	nvml  nvml.Interface
}

// New creates a new instance of an NVDeviceLib given an nvml.Interface to work
// from.
func New(nvmllib nvml.Interface) *NVDeviceLib {
	nvdevlib := nvdev.New(
		nvdev.WithNvml(nvmllib),
	)
	return &NVDeviceLib{
		nvml:  nvmllib,
		nvdev: nvdevlib,
	}
}

// Init initializes an NVDeviceLib for use.
func (l NVDeviceLib) Init() error {
	ret := l.nvml.Init()
	if ret != nvml.SUCCESS {
		return fmt.Errorf("error initializing NVML: %w", ret)
	}
	return nil
}

// AlwaysShutdown unconditionally shuts down an NVDeviceLib logging any errors.
func (l NVDeviceLib) AlwaysShutdown() {
	ret := l.nvml.Shutdown()
	if ret != nvml.SUCCESS {
		klog.Warningf("error shutting down NVML: %v", ret)
	}
}

// GetPerGpuAllocatableDevices gets the set of allocatable devices using
// NVDeviceLib.  A list of GPU indices can be optionally provided to limit the
// set of allocatable devices to just those GPUs. If no indices are provided,
// the full set of allocatable devices across all GPUs are returned.
// NOTE: Both full GPUs and MIG devices are returned as part of this call.
func (l NVDeviceLib) GetPerGpuAllocatableDevices(indices ...int) (PerGpuAllocatableDevices, error) {
	if err := l.Init(); err != nil {
		return nil, err
	}
	defer l.AlwaysShutdown()

	allocatable := make(PerGpuAllocatableDevices)
	err := l.nvdev.VisitDevices(func(i int, d nvdev.Device) error {
		if indices != nil && !slices.Contains(indices, i) {
			return nil
		}

		gpuInfo, err := l.getGpuInfo(i, d)
		if err != nil {
			return fmt.Errorf("error getting info for GPU %v: %w", i, err)
		}
		gpuDevice := AllocatableDevice{
			Gpu: gpuInfo,
		}
		allocatable[gpuInfo.Index] = append(allocatable[gpuInfo.Index], gpuDevice)

		migInfos, err := l.getMigInfos(gpuInfo, d)
		if err != nil {
			return fmt.Errorf("error getting MIG info for GPU %v: %w", i, err)
		}
		for _, migInfo := range migInfos {
			migDevice := AllocatableDevice{
				Mig: migInfo,
			}
			allocatable[gpuInfo.Index] = append(allocatable[gpuInfo.Index], migDevice)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error visiting devices: %w", err)
	}

	return allocatable, nil
}

// getGpuInfo returns info about the GPU at the provided index.
func (l NVDeviceLib) getGpuInfo(index int, device nvdev.Device) (*GpuInfo, error) {
	minor, ret := device.GetMinorNumber()
	if ret != nvml.SUCCESS {
		return nil, fmt.Errorf("error getting minor number for device %d: %w", index, ret)
	}
	uuid, ret := device.GetUUID()
	if ret != nvml.SUCCESS {
		return nil, fmt.Errorf("error getting UUID for device %d: %w", index, ret)
	}
	memory, ret := device.GetMemoryInfo()
	if ret != nvml.SUCCESS {
		return nil, fmt.Errorf("error getting memory info for device %d: %w", index, ret)
	}
	productName, ret := device.GetName()
	if ret != nvml.SUCCESS {
		return nil, fmt.Errorf("error getting product name for device %d: %w", index, ret)
	}
	architecture, err := device.GetArchitectureAsString()
	if err != nil {
		return nil, fmt.Errorf("error getting architecture for device %d: %w", index, err)
	}
	brand, err := device.GetBrandAsString()
	if err != nil {
		return nil, fmt.Errorf("error getting brand for device %d: %w", index, err)
	}
	cudaComputeCapability, err := device.GetCudaComputeCapabilityAsString()
	if err != nil {
		return nil, fmt.Errorf("error getting CUDA compute capability for device %d: %w", index, err)
	}
	driverVersion, ret := l.nvml.SystemGetDriverVersion()
	if ret != nvml.SUCCESS {
		return nil, fmt.Errorf("error getting driver version: %w", err)
	}
	cudaDriverVersion, ret := l.nvml.SystemGetCudaDriverVersion()
	if ret != nvml.SUCCESS {
		return nil, fmt.Errorf("error getting CUDA driver version: %w", err)
	}
	migCapable, err := device.IsMigCapable()
	if err != nil {
		return nil, fmt.Errorf("error checking if MIG capable for device %d: %w", index, err)
	}
	migEnabled, err := device.IsMigEnabled()
	if err != nil {
		return nil, fmt.Errorf("error checking if MIG mode enabled for device %d: %w", index, err)
	}

	gpuInfo := &GpuInfo{
		Minor:                 minor,
		Index:                 index,
		UUID:                  uuid,
		MemoryBytes:           memory.Total,
		ProductName:           productName,
		Brand:                 brand,
		Architecture:          architecture,
		CudaComputeCapability: cudaComputeCapability,
		DriverVersion:         driverVersion,
		CudaDriverVersion:     fmt.Sprintf("%v.%v", cudaDriverVersion/1000, (cudaDriverVersion%1000)/10),
		MigCapable:            migCapable,
		MigEnabled:            migEnabled,
	}

	return gpuInfo, nil
}

// getMigInfos returns a list of MigInfos for the GPU represented by device.
func (l NVDeviceLib) getMigInfos(gpuInfo *GpuInfo, device nvdev.Device) ([]*MigInfo, error) {
	var migInfos []*MigInfo
	err := device.VisitMigProfiles(func(migProfile nvdev.MigProfile) error {
		if migProfile.GetInfo().C != migProfile.GetInfo().G {
			return nil
		}

		if migProfile.GetInfo().CIProfileID == nvml.COMPUTE_INSTANCE_PROFILE_1_SLICE_REV1 {
			return nil
		}

		giProfileInfo, ret := device.GetGpuInstanceProfileInfo(migProfile.GetInfo().GIProfileID)
		if ret == nvml.ERROR_NOT_SUPPORTED {
			return nil
		}
		if ret == nvml.ERROR_INVALID_ARGUMENT {
			return nil
		}
		if ret != nvml.SUCCESS {
			return fmt.Errorf("error getting GI Profile info for MIG profile %v: %w", migProfile, ret)
		}

		giPlacements, ret := device.GetGpuInstancePossiblePlacements(&giProfileInfo)
		if ret == nvml.ERROR_NOT_SUPPORTED {
			return nil
		}
		if ret == nvml.ERROR_INVALID_ARGUMENT {
			return nil
		}
		if ret != nvml.SUCCESS {
			return fmt.Errorf("error getting GI possible placements for MIG profile %v: %w", migProfile, ret)
		}

		for _, giPlacement := range giPlacements {
			migInfo := &MigInfo{
				Parent:        gpuInfo,
				Profile:       migProfile,
				GIProfileInfo: giProfileInfo,
				MemorySlices:  giPlacement,
			}
			migInfos = append(migInfos, migInfo)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error visiting MIG profiles: %w", err)
	}

	return migInfos, nil
}
