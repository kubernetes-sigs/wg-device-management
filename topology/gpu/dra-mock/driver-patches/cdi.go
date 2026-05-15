/*
Copyright The Kubernetes Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	nvdevice "github.com/NVIDIA/go-nvlib/pkg/nvlib/device"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/NVIDIA/nvidia-container-toolkit/pkg/nvcdi"
	"github.com/NVIDIA/nvidia-container-toolkit/pkg/nvcdi/spec"
	transformroot "github.com/NVIDIA/nvidia-container-toolkit/pkg/nvcdi/transform/root"
	"k8s.io/klog/v2"

	utilcache "k8s.io/apimachinery/pkg/util/cache"
	cdiapi "tags.cncf.io/container-device-interface/pkg/cdi"
	cdiparser "tags.cncf.io/container-device-interface/pkg/parser"
	cdispec "tags.cncf.io/container-device-interface/specs-go"

	"sigs.k8s.io/dra-driver-nvidia-gpu/internal/common"
)

const (
	cdiVendor      = "k8s." + DriverName
	cdiClaimClass  = "claim"
	defaultCDIRoot = "/var/run/cdi"
	procNvCapsPath = "/proc/driver/nvidia/capabilities"
)

type CDIHandler struct {
	logger            *logrus.Logger
	nvml              nvml.Interface
	nvdevice          nvdevice.Interface
	nvcdiClaim        nvcdi.Interface
	vfiocdi           *vfioCDIHandler
	driverRoot        string
	devRoot           string
	targetDriverRoot  string
	nvidiaCDIHookPath string

	specCache *utilcache.Expiring

	cdiRoot string
}

func NewCDIHandler(opts ...cdiOption) (*CDIHandler, error) {
	h := &CDIHandler{}
	for _, opt := range opts {
		opt(h)
	}

	if h.logger == nil {
		h.logger = logrus.New()
		h.logger.SetOutput(io.Discard)
	}
	if h.nvml == nil {
		h.nvml = nvml.New()
	}
	if h.cdiRoot == "" {
		h.cdiRoot = defaultCDIRoot
	}
	if h.nvdevice == nil {
		h.nvdevice = nvdevice.New(h.nvml)
	}

	if h.nvcdiClaim == nil {
		nvcdilib, err := nvcdi.New(
			nvcdi.WithDeviceLib(h.nvdevice),
			nvcdi.WithDriverRoot(h.driverRoot),
			nvcdi.WithDevRoot(h.devRoot),
			nvcdi.WithLogger(h.logger),
			nvcdi.WithNvmlLib(h.nvml),
			nvcdi.WithMode("nvml"),
			nvcdi.WithVendor(cdiVendor),
			nvcdi.WithClass(cdiClaimClass),
			nvcdi.WithNVIDIACDIHookPath(h.nvidiaCDIHookPath),
			nvcdi.WithFeatureFlags(nvcdi.FeatureDisableNvsandboxUtils),
		)
		if err != nil {
			return nil, fmt.Errorf("unable to create CDI library for claims: %w", err)
		}
		h.nvcdiClaim = nvcdilib
	}

	// The expiration time is defined upon key insert, not cache-globally.
	h.specCache = utilcache.NewExpiring()

	return h, nil
}

func (cdi *CDIHandler) GetCommonEditsCached() (*cdiapi.ContainerEdits, error) {
	key := "commonEdits"
	if v, ok := cdi.specCache.Get(key); ok {
		edits, ok := v.(*cdiapi.ContainerEdits)
		if !ok {
			return nil, fmt.Errorf("expected *cdiapi.ContainerEdits, got %T", v)
		}
		// Return a shallow copy so that cache entry consumer is less likely to
		// mutate the cache entry.
		clone := *edits
		return &clone, nil
	}

	t0 := time.Now()
	v, err := cdi.nvcdiClaim.GetCommonEdits()
	klog.V(7).Infof("t_cdi_get_common_edits %.3f s", time.Since(t0).Seconds())

	if err != nil {
		return nil, err
	}
	cdi.specCache.Set(key, v, time.Duration(5*time.Minute))
	// Return a shallow copy, see above.
	clone := *v
	return &clone, nil
}

func (cdi *CDIHandler) WarmupDevSpecCache(uuids []string) {
	for _, uuid := range uuids {
		_, err := cdi.GetDeviceSpecsByUUIDCached(uuid)
		if err != nil {
			klog.Warningf("Ignore error during cache warmup: GetDeviceSpecsByUUIDCached() failed: %s", err)
		}
	}
}

func (cdi *CDIHandler) GetDeviceSpecsByUUIDCached(uuid string) ([]cdispec.Device, error) {
	key := uuid
	if v, ok := cdi.specCache.Get(key); ok {
		devs, ok := v.([]cdispec.Device)
		if !ok {
			return nil, fmt.Errorf("expected []cdispec.Device, got %T", v)
		}
		clone := make([]cdispec.Device, len(devs))
		copy(clone, devs)
		return clone, nil
	}

	t0 := time.Now()
	devs, err := cdi.nvcdiClaim.GetDeviceSpecsByID(uuid)
	klog.V(1).Infof("GetDeviceSpecsByID() called for %s, t_cdi_get_device_specs_by_id %.3f s", uuid, time.Since(t0).Seconds())
	if err != nil {
		return nil, err
	}
	cdi.specCache.Set(key, devs, time.Duration(5*time.Minute))
	clone := make([]cdispec.Device, len(devs))
	copy(clone, devs)
	return clone, nil
}

// Note(JP): for a regular GPU, this canonical name is for example `gpu-0`, with
// the numerical suffix as of the time of writing reflecting the device minor.
// NVMLs' DeviceSetMigMode() is documented with 'This API may unbind or reset
// the device to activate the requested mode. Thus, the attributes associated
// with the device, such as minor number, might change. The caller of this API
// is expected to query such attributes again.' -- if the minor is indeed not
// necessarily stable, there may be problems associating this spec _long-term_
// with that name. That is an argument for always dynamically generating also
// full-GPU CDI spec during prepare() (or: to cache it, and re-generate it every
// now and then during this program's lifetime).
func (cdi *CDIHandler) CreateClaimSpecFile(claimUID string, preparedDevices PreparedDevices) error {
	specName := cdiapi.GenerateTransientSpecName(cdiVendor, cdiClaimClass, claimUID)
	filePath := filepath.Join(cdi.cdiRoot, specName+".yaml")

	content := fmt.Sprintf(`cdiVersion: "0.5.0"
kind: "%s/%s"
devices:
`, cdiVendor, cdiClaimClass)

	for _, group := range preparedDevices {
		for _, dev := range group.Devices {
			dname := fmt.Sprintf("%s-%s", claimUID, dev.CanonicalName())
			content += fmt.Sprintf(`  - name: "%s"
    containerEdits:
      env:
          - MOCK_DEVICE=%s
`, dname, dev.CanonicalName())
		}
	}

	content += fmt.Sprintf(`containerEdits:
    env:
        - NVIDIA_VISIBLE_DEVICES=void
`)

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write mock CDI spec: %w", err)
	}
	klog.V(6).Infof("Wrote mock CDI spec to %s", filePath)
	return nil
}

func (cdi *CDIHandler) DeleteClaimSpecFile(claimUID string) error {
	specName := cdiapi.GenerateTransientSpecName(cdiVendor, cdiClaimClass, claimUID)
	klog.V(6).Infof("Delete CDI spec file: '%s', claim '%s'", specName, claimUID)
	err := os.Remove(filepath.Join(cdi.cdiRoot, specName+".yaml"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Philosophy: all devices to be injected into a container are defined in a
// single, transient CDI spec. This function returns the fully qualified
// identifier for a device defined in that spec. Example:
// k8s.gpu.nvidia.com/claim=dab5ab50-d59a-42a6-af16-cfd4628c0f7a-gpu-0
// That identifier can be used elsewhere, and _points to the spec_.
func (cdi *CDIHandler) GetClaimDeviceName(claimUID string, device *AllocatableDevice, containerEdits *cdiapi.ContainerEdits) string {
	return cdiparser.QualifiedName(cdiVendor, cdiClaimClass, fmt.Sprintf("%s-%s", claimUID, device.CanonicalName()))
}

// Construct and return the CDI `deviceNodes` specification for the two
// character devices `/dev/nvidia-caps/nvidia-cap<CIm>` and
// `/dev/nvidia-caps/nvidia-cap<GIm>` for a specific MIG device.
//
// Context: for containerized workload to see and use a specific MIG device, it
// needs to be able to open three character device nodes:
//
// 1) `/dev/nvidia<Pm>`, with <Pm> referring to the parent's minor. This exists
// on the host.
//
// 2) /dev/nvidia-caps/nvidia-cap<CIm> and /dev/nvidia-caps/nvidia-cap<GIm>,
// with <GIm> and <CIm> referring to the MIG GPU instance's and Compute
// instance's minor, respectively. For the the latter two device nodes it is
// sufficient to create them in the container (with proper cgroups permissions),
// without actually requiring the same device nodes to be explicitly created on
// the host. That is what is achieved below with the structure created in
// cdiDevNodeFromNVCapDevInfo().
func (cdi *CDIHandler) GetDevNodesForMigDevice(mlt *MigLiveTuple) ([]*cdispec.DeviceNode, error) {
	gipath := fmt.Sprintf("%s/gpu%d/mig/gi%d/access", procNvCapsPath, mlt.ParentMinor, mlt.GIID)
	cipath := fmt.Sprintf("%s/gpu%d/mig/gi%d/ci%d/access", procNvCapsPath, mlt.ParentMinor, mlt.GIID, mlt.CIID)

	giCapsInfo, err := common.ParseNVCapDeviceInfo(gipath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GI capabilities file %s: %w", gipath, err)
	}

	ciCapsInfo, err := common.ParseNVCapDeviceInfo(cipath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CI capabilities file %s: %w", cipath, err)
	}

	devnodes := []*cdispec.DeviceNode{giCapsInfo.CDICharDevNode(), ciCapsInfo.CDICharDevNode()}
	return devnodes, nil
}

// Write CDI spec to the filesystem.
func (cdi *CDIHandler) writeSpec(spec spec.Interface, specName string) error {
	// Transform the spec to make it aware that it is running inside a container.
	err := transformroot.New(
		transformroot.WithRoot(cdi.driverRoot),
		transformroot.WithTargetRoot(cdi.targetDriverRoot),
		transformroot.WithRelativeTo("host"),
	).Transform(spec.Raw())
	if err != nil {
		return fmt.Errorf("failed to transform driver root in CDI spec: %w", err)
	}

	klog.V(7).Infof("Write CDI spec: %s", specName)
	return spec.Save(filepath.Join(cdi.cdiRoot, specName+".yaml"))
}
