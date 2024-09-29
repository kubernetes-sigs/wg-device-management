# nvk8s-resourcemodel

This repo is meant as a staging ground for defining extensions to the
Kubernetes structured parameters model as it pertains to NVIDIA hardware.
All of the running examples use a Mock DGXA100 server to feed its input, so the
results it generates are comparable to what we would see in real hardware.

## Proposed Extensions

Tese extensions included in this repo allow us to enable fully dynamic MIG with
the possibility for the scheduler to do intelligent allocation to avoid
fragmentation.

## Examples / Proof of Concept

There are two different commands available:

**`print-model`**:
    This command will print out two ResourceSliceSpecs for GPU 0 on the Mock
	DGXA100 server. The first is the original spec and the second is a
	flattened version of it with all "mixins" resolved.

**`print-possible-allocations`**:
    This command will print out all possible allocations of the resources
	declared for GPU 0 on the Mock DGXA100 server using the proposed
	extensions. This is meant to demonstrate how the scheduler can use these
	extensions to dynamically allocate non-overlapping devices partitioned from
	the same piece of hardware.

To run these commands, invoke `make <cmd-name>`.

The output of running each command can be seen below:

```console
$ make print-model

Original spec:
deviceMixins:
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      productName:
        string: NVIDIA A100-SXM4-40GB
      type:
        string: gpu
  name: common-gpu-nvidia-a100-sxm4-40gb-attributes
- basic:
    attributes:
      index:
        int: 0
      minor:
        int: 0
      uuid:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
  name: gpu-0-attributes
- basic:
    capacity:
      memorySlice5:
        quantity: "1"
  name: memory-slices-5
- basic:
    attributes:
      profile:
        string: 1g.5gb+me
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "1"
      memory:
        quantity: 4864Mi
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "1"
    includes:
    - name: system-attributes
    - name: common-mig-nvidia-a100-sxm4-40gb-attributes
  name: mig-1g.5gb-me-nvidia-a100-sxm4-40gb
- basic:
    attributes:
      profile:
        string: 1g.10gb
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 9984Mi
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "0"
    includes:
    - name: system-attributes
    - name: common-mig-nvidia-a100-sxm4-40gb-attributes
  name: mig-1g.10gb-nvidia-a100-sxm4-40gb
- basic:
    attributes:
      profile:
        string: 7g.40gb
    capacity:
      copy-engines:
        quantity: "7"
      decoders:
        quantity: "5"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "1"
      memory:
        quantity: 40320Mi
      multiprocessors:
        quantity: "98"
      ofa-engines:
        quantity: "1"
    includes:
    - name: system-attributes
    - name: common-mig-nvidia-a100-sxm4-40gb-attributes
  name: mig-7g.40gb-nvidia-a100-sxm4-40gb
- basic:
    attributes:
      profile:
        string: 1g.5gb
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "0"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 4864Mi
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "0"
    includes:
    - name: system-attributes
    - name: common-mig-nvidia-a100-sxm4-40gb-attributes
  name: mig-1g.5gb-nvidia-a100-sxm4-40gb
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      productName:
        string: NVIDIA A100-SXM4-40GB
      type:
        string: mig
  name: common-mig-nvidia-a100-sxm4-40gb-attributes
- basic:
    capacity:
      memorySlice4:
        quantity: "1"
  name: memory-slices-4
- basic:
    capacity:
      memorySlice6:
        quantity: "1"
  name: memory-slices-6
- basic:
    attributes:
      profile:
        string: 2g.10gb
    capacity:
      copy-engines:
        quantity: "2"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 9984Mi
      multiprocessors:
        quantity: "28"
      ofa-engines:
        quantity: "0"
    includes:
    - name: system-attributes
    - name: common-mig-nvidia-a100-sxm4-40gb-attributes
  name: mig-2g.10gb-nvidia-a100-sxm4-40gb
- basic:
    capacity:
      memorySlice0:
        quantity: "1"
      memorySlice1:
        quantity: "1"
      memorySlice2:
        quantity: "1"
      memorySlice3:
        quantity: "1"
      memorySlice4:
        quantity: "1"
      memorySlice5:
        quantity: "1"
      memorySlice6:
        quantity: "1"
      memorySlice7:
        quantity: "1"
  name: memory-slices-0-7
- basic:
    capacity:
      copy-engines:
        quantity: "7"
      decoders:
        quantity: "5"
      jpeg-engines:
        quantity: "1"
      multiprocessors:
        quantity: "98"
      ofa-engines:
        quantity: "1"
    includes:
    - name: system-attributes
    - name: common-gpu-nvidia-a100-sxm4-40gb-attributes
  name: gpu-nvidia-a100-sxm4-40gb
- basic:
    attributes:
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
  name: system-attributes
- basic:
    capacity:
      memorySlice3:
        quantity: "1"
  name: memory-slices-3
- basic:
    capacity:
      memorySlice4:
        quantity: "1"
      memorySlice5:
        quantity: "1"
  name: memory-slices-4-5
- basic:
    attributes:
      profile:
        string: 3g.20gb
    capacity:
      copy-engines:
        quantity: "3"
      decoders:
        quantity: "2"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 20096Mi
      multiprocessors:
        quantity: "42"
      ofa-engines:
        quantity: "0"
    includes:
    - name: system-attributes
    - name: common-mig-nvidia-a100-sxm4-40gb-attributes
  name: mig-3g.20gb-nvidia-a100-sxm4-40gb
- basic:
    capacity:
      memorySlice0:
        quantity: "1"
      memorySlice1:
        quantity: "1"
      memorySlice2:
        quantity: "1"
      memorySlice3:
        quantity: "1"
  name: memory-slices-0-3
- basic:
    capacity:
      memorySlice6:
        quantity: "1"
      memorySlice7:
        quantity: "1"
  name: memory-slices-6-7
- basic:
    attributes:
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
  name: gpu-0-mig-attributes
- basic:
    capacity:
      memorySlice2:
        quantity: "1"
  name: memory-slices-2
- basic:
    capacity:
      memorySlice0:
        quantity: "1"
      memorySlice1:
        quantity: "1"
  name: memory-slices-0-1
- basic:
    attributes:
      profile:
        string: 4g.20gb
    capacity:
      copy-engines:
        quantity: "4"
      decoders:
        quantity: "2"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 20096Mi
      multiprocessors:
        quantity: "56"
      ofa-engines:
        quantity: "0"
    includes:
    - name: system-attributes
    - name: common-mig-nvidia-a100-sxm4-40gb-attributes
  name: mig-4g.20gb-nvidia-a100-sxm4-40gb
- basic:
    capacity:
      memorySlice2:
        quantity: "1"
      memorySlice3:
        quantity: "1"
  name: memory-slices-2-3
- basic:
    capacity:
      memorySlice4:
        quantity: "1"
      memorySlice5:
        quantity: "1"
      memorySlice6:
        quantity: "1"
      memorySlice7:
        quantity: "1"
  name: memory-slices-4-7
- basic:
    capacity:
      memorySlice0:
        quantity: "1"
  name: memory-slices-0
- basic:
    capacity:
      memorySlice1:
        quantity: "1"
  name: memory-slices-1
devices:
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-3g.20gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-0-3
  name: gpu-0-mig-3g.20gb-0-3
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-3g.20gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-4-7
  name: gpu-0-mig-3g.20gb-4-7
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-7g.40gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-0-7
  name: gpu-0-mig-7g.40gb-0-7
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.5gb-me-nvidia-a100-sxm4-40gb
    - name: memory-slices-4
  name: gpu-0-mig-1g.5gb-me-4
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.5gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-0
  name: gpu-0-mig-1g.5gb-0
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.5gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-2
  name: gpu-0-mig-1g.5gb-2
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.5gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-4
  name: gpu-0-mig-1g.5gb-4
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.5gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-5
  name: gpu-0-mig-1g.5gb-5
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.10gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-4-5
  name: gpu-0-mig-1g.10gb-4-5
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.10gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-6-7
  name: gpu-0-mig-1g.10gb-6-7
- basic:
    includes:
    - name: gpu-nvidia-a100-sxm4-40gb
    - name: gpu-0-attributes
    - name: memory-slices-0-7
  name: gpu-0
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.5gb-me-nvidia-a100-sxm4-40gb
    - name: memory-slices-3
  name: gpu-0-mig-1g.5gb-me-3
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.5gb-me-nvidia-a100-sxm4-40gb
    - name: memory-slices-6
  name: gpu-0-mig-1g.5gb-me-6
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.5gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-1
  name: gpu-0-mig-1g.5gb-1
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.5gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-6
  name: gpu-0-mig-1g.5gb-6
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-2g.10gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-0-1
  name: gpu-0-mig-2g.10gb-0-1
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-2g.10gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-4-5
  name: gpu-0-mig-2g.10gb-4-5
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.5gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-3
  name: gpu-0-mig-1g.5gb-3
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-2g.10gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-2-3
  name: gpu-0-mig-2g.10gb-2-3
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.5gb-me-nvidia-a100-sxm4-40gb
    - name: memory-slices-2
  name: gpu-0-mig-1g.5gb-me-2
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.10gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-0-1
  name: gpu-0-mig-1g.10gb-0-1
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.10gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-2-3
  name: gpu-0-mig-1g.10gb-2-3
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-4g.20gb-nvidia-a100-sxm4-40gb
    - name: memory-slices-0-3
  name: gpu-0-mig-4g.20gb-0-3
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.5gb-me-nvidia-a100-sxm4-40gb
    - name: memory-slices-0
  name: gpu-0-mig-1g.5gb-me-0
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.5gb-me-nvidia-a100-sxm4-40gb
    - name: memory-slices-1
  name: gpu-0-mig-1g.5gb-me-1
- basic:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: gpu-0-mig-attributes
    - name: mig-1g.5gb-me-nvidia-a100-sxm4-40gb
    - name: memory-slices-5
  name: gpu-0-mig-1g.5gb-me-5

Flattened spec:
devices:
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 3g.20gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "3"
      decoders:
        quantity: "2"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 20096Mi
      memorySlice0:
        quantity: "1"
      memorySlice1:
        quantity: "1"
      memorySlice2:
        quantity: "1"
      memorySlice3:
        quantity: "1"
      multiprocessors:
        quantity: "42"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-3g.20gb-0-3
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 3g.20gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "3"
      decoders:
        quantity: "2"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 20096Mi
      memorySlice4:
        quantity: "1"
      memorySlice5:
        quantity: "1"
      memorySlice6:
        quantity: "1"
      memorySlice7:
        quantity: "1"
      multiprocessors:
        quantity: "42"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-3g.20gb-4-7
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 7g.40gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "7"
      decoders:
        quantity: "5"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "1"
      memory:
        quantity: 40320Mi
      memorySlice0:
        quantity: "1"
      memorySlice1:
        quantity: "1"
      memorySlice2:
        quantity: "1"
      memorySlice3:
        quantity: "1"
      memorySlice4:
        quantity: "1"
      memorySlice5:
        quantity: "1"
      memorySlice6:
        quantity: "1"
      memorySlice7:
        quantity: "1"
      multiprocessors:
        quantity: "98"
      ofa-engines:
        quantity: "1"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-7g.40gb-0-7
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.5gb+me
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "1"
      memory:
        quantity: 4864Mi
      memorySlice4:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "1"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.5gb-me-4
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.5gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "0"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 4864Mi
      memorySlice0:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.5gb-0
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.5gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "0"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 4864Mi
      memorySlice2:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.5gb-2
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.5gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "0"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 4864Mi
      memorySlice4:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.5gb-4
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.5gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "0"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 4864Mi
      memorySlice5:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.5gb-5
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.10gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 9984Mi
      memorySlice4:
        quantity: "1"
      memorySlice5:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.10gb-4-5
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.10gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 9984Mi
      memorySlice6:
        quantity: "1"
      memorySlice7:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.10gb-6-7
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      index:
        int: 0
      minor:
        int: 0
      productName:
        string: NVIDIA A100-SXM4-40GB
      type:
        string: gpu
      uuid:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
    capacity:
      copy-engines:
        quantity: "7"
      decoders:
        quantity: "5"
      jpeg-engines:
        quantity: "1"
      memorySlice0:
        quantity: "1"
      memorySlice1:
        quantity: "1"
      memorySlice2:
        quantity: "1"
      memorySlice3:
        quantity: "1"
      memorySlice4:
        quantity: "1"
      memorySlice5:
        quantity: "1"
      memorySlice6:
        quantity: "1"
      memorySlice7:
        quantity: "1"
      multiprocessors:
        quantity: "98"
      ofa-engines:
        quantity: "1"
  name: gpu-0
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.5gb+me
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "1"
      memory:
        quantity: 4864Mi
      memorySlice3:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "1"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.5gb-me-3
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.5gb+me
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "1"
      memory:
        quantity: 4864Mi
      memorySlice6:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "1"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.5gb-me-6
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.5gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "0"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 4864Mi
      memorySlice1:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.5gb-1
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.5gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "0"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 4864Mi
      memorySlice6:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.5gb-6
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 2g.10gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "2"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 9984Mi
      memorySlice0:
        quantity: "1"
      memorySlice1:
        quantity: "1"
      multiprocessors:
        quantity: "28"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-2g.10gb-0-1
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 2g.10gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "2"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 9984Mi
      memorySlice4:
        quantity: "1"
      memorySlice5:
        quantity: "1"
      multiprocessors:
        quantity: "28"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-2g.10gb-4-5
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.5gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "0"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 4864Mi
      memorySlice3:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.5gb-3
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 2g.10gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "2"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 9984Mi
      memorySlice2:
        quantity: "1"
      memorySlice3:
        quantity: "1"
      multiprocessors:
        quantity: "28"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-2g.10gb-2-3
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.5gb+me
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "1"
      memory:
        quantity: 4864Mi
      memorySlice2:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "1"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.5gb-me-2
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.10gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 9984Mi
      memorySlice0:
        quantity: "1"
      memorySlice1:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.10gb-0-1
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.10gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 9984Mi
      memorySlice2:
        quantity: "1"
      memorySlice3:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.10gb-2-3
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 4g.20gb
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "4"
      decoders:
        quantity: "2"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "0"
      memory:
        quantity: 20096Mi
      memorySlice0:
        quantity: "1"
      memorySlice1:
        quantity: "1"
      memorySlice2:
        quantity: "1"
      memorySlice3:
        quantity: "1"
      multiprocessors:
        quantity: "56"
      ofa-engines:
        quantity: "0"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-4g.20gb-0-3
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.5gb+me
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "1"
      memory:
        quantity: 4864Mi
      memorySlice0:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "1"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.5gb-me-0
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.5gb+me
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "1"
      memory:
        quantity: 4864Mi
      memorySlice1:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "1"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.5gb-me-1
- basic:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.6"
      driverVersion:
        version: 560.35.03
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-4cf8db2d-06c0-7d70-1a51-e59b25b2c16c
      productName:
        string: NVIDIA A100-SXM4-40GB
      profile:
        string: 1g.5gb+me
      type:
        string: mig
    capacity:
      copy-engines:
        quantity: "1"
      decoders:
        quantity: "1"
      encoders:
        quantity: "0"
      jpeg-engines:
        quantity: "1"
      memory:
        quantity: 4864Mi
      memorySlice5:
        quantity: "1"
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "1"
    consumesCapacityFrom:
    - name: gpu-0
  name: gpu-0-mig-1g.5gb-me-5
```

```console
$ make print-possible-allocations

[gpu-0]
[gpu-0-mig-1g.10gb-0-1]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-2-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-0-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-2-3]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-4-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-0 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-2-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.10gb-4-5]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.10gb-6-7 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.10gb-4-5 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-4 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-5 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-4-5 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-2g.10gb-4-5 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.10gb-6-7 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-0]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4-5 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-me-0 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-me-0 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-me-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-me-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-me-2 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-me-2 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-me-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-me-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-me-4 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-me-4 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-me-4 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-me-5 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-me-5 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-me-5 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-2g.10gb-4-5 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-2g.10gb-0-1]
[gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-2-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-2g.10gb-0-1 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-2g.10gb-0-1 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-2g.10gb-2-3]
[gpu-0-mig-2g.10gb-2-3 gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-2g.10gb-2-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-2g.10gb-4-5]
[gpu-0-mig-2g.10gb-4-5 gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-2g.10gb-4-5 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-3g.20gb-0-3]
[gpu-0-mig-3g.20gb-0-3 gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-3g.20gb-4-7]
[gpu-0-mig-3g.20gb-4-7 gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-4g.20gb-0-3]
[gpu-0-mig-7g.40gb-0-7]
```
