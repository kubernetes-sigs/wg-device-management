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
- name: common-gpu-mock-nvidia-a100-sxm4-40gb-attributes
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
      type:
        string: gpu
- name: common-gpu-mock-nvidia-a100-sxm4-40gb-capacities
  partitionable:
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
        quantity: 40Gi
      multiprocessors:
        quantity: "98"
      ofa-engines:
        quantity: "1"
- name: common-mig-1g.10gb-mock-nvidia-a100-sxm4-40gb
  partitionable:
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
        quantity: 9856Mi
      multiprocessors:
        quantity: "14"
      ofa-engines:
        quantity: "0"
- name: common-mig-1g.5gb-me-mock-nvidia-a100-sxm4-40gb
  partitionable:
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
- name: common-mig-1g.5gb-mock-nvidia-a100-sxm4-40gb
  partitionable:
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
- name: common-mig-2g.10gb-mock-nvidia-a100-sxm4-40gb
  partitionable:
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
        quantity: 9856Mi
      multiprocessors:
        quantity: "28"
      ofa-engines:
        quantity: "0"
- name: common-mig-3g.20gb-mock-nvidia-a100-sxm4-40gb
  partitionable:
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
        quantity: 19968Mi
      multiprocessors:
        quantity: "42"
      ofa-engines:
        quantity: "0"
- name: common-mig-4g.20gb-mock-nvidia-a100-sxm4-40gb
  partitionable:
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
        quantity: 19968Mi
      multiprocessors:
        quantity: "56"
      ofa-engines:
        quantity: "0"
- name: common-mig-7g.40gb-mock-nvidia-a100-sxm4-40gb
  partitionable:
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
        quantity: 40192Mi
      multiprocessors:
        quantity: "98"
      ofa-engines:
        quantity: "1"
- name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
      type:
        string: mig
- name: memory-slices-0
  partitionable:
    capacity:
      memorySlice0:
        quantity: "1"
- name: memory-slices-0-1
  partitionable:
    capacity:
      memorySlice0:
        quantity: "1"
      memorySlice1:
        quantity: "1"
- name: memory-slices-0-3
  partitionable:
    capacity:
      memorySlice0:
        quantity: "1"
      memorySlice1:
        quantity: "1"
      memorySlice2:
        quantity: "1"
      memorySlice3:
        quantity: "1"
- name: memory-slices-0-7
  partitionable:
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
- name: memory-slices-1
  partitionable:
    capacity:
      memorySlice1:
        quantity: "1"
- name: memory-slices-2
  partitionable:
    capacity:
      memorySlice2:
        quantity: "1"
- name: memory-slices-2-3
  partitionable:
    capacity:
      memorySlice2:
        quantity: "1"
      memorySlice3:
        quantity: "1"
- name: memory-slices-3
  partitionable:
    capacity:
      memorySlice3:
        quantity: "1"
- name: memory-slices-4
  partitionable:
    capacity:
      memorySlice4:
        quantity: "1"
- name: memory-slices-4-5
  partitionable:
    capacity:
      memorySlice4:
        quantity: "1"
      memorySlice5:
        quantity: "1"
- name: memory-slices-4-7
  partitionable:
    capacity:
      memorySlice4:
        quantity: "1"
      memorySlice5:
        quantity: "1"
      memorySlice6:
        quantity: "1"
      memorySlice7:
        quantity: "1"
- name: memory-slices-5
  partitionable:
    capacity:
      memorySlice5:
        quantity: "1"
- name: memory-slices-6
  partitionable:
    capacity:
      memorySlice6:
        quantity: "1"
- name: memory-slices-6-7
  partitionable:
    capacity:
      memorySlice6:
        quantity: "1"
      memorySlice7:
        quantity: "1"
- name: specific-gpu-0-attributes
  partitionable:
    attributes:
      index:
        int: 0
      minor:
        int: 0
      uuid:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
- name: specific-gpu-0-mig-attributes
  partitionable:
    attributes:
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
- name: system-attributes
  partitionable:
    attributes:
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
devices:
- name: gpu-0
  partitionable:
    includes:
    - name: system-attributes
    - name: common-gpu-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-gpu-mock-nvidia-a100-sxm4-40gb-capacities
    - name: specific-gpu-0-attributes
    - name: memory-slices-0-7
- name: gpu-0-mig-1g.10gb-0-1
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.10gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-0-1
- name: gpu-0-mig-1g.10gb-2-3
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.10gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-2-3
- name: gpu-0-mig-1g.10gb-4-5
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.10gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-4-5
- name: gpu-0-mig-1g.10gb-6-7
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.10gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-6-7
- name: gpu-0-mig-1g.5gb-0
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.5gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-0
- name: gpu-0-mig-1g.5gb-1
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.5gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-1
- name: gpu-0-mig-1g.5gb-2
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.5gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-2
- name: gpu-0-mig-1g.5gb-3
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.5gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-3
- name: gpu-0-mig-1g.5gb-4
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.5gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-4
- name: gpu-0-mig-1g.5gb-5
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.5gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-5
- name: gpu-0-mig-1g.5gb-6
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.5gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-6
- name: gpu-0-mig-1g.5gb-me-0
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.5gb-me-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-0
- name: gpu-0-mig-1g.5gb-me-1
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.5gb-me-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-1
- name: gpu-0-mig-1g.5gb-me-2
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.5gb-me-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-2
- name: gpu-0-mig-1g.5gb-me-3
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.5gb-me-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-3
- name: gpu-0-mig-1g.5gb-me-4
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.5gb-me-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-4
- name: gpu-0-mig-1g.5gb-me-5
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.5gb-me-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-5
- name: gpu-0-mig-1g.5gb-me-6
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-1g.5gb-me-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-6
- name: gpu-0-mig-2g.10gb-0-1
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-2g.10gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-0-1
- name: gpu-0-mig-2g.10gb-2-3
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-2g.10gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-2-3
- name: gpu-0-mig-2g.10gb-4-5
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-2g.10gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-4-5
- name: gpu-0-mig-3g.20gb-0-3
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-3g.20gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-0-3
- name: gpu-0-mig-3g.20gb-4-7
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-3g.20gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-4-7
- name: gpu-0-mig-4g.20gb-0-3
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-4g.20gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-0-3
- name: gpu-0-mig-7g.40gb-0-7
  partitionable:
    consumesCapacityFrom:
    - name: gpu-0
    includes:
    - name: system-attributes
    - name: common-mig-mock-nvidia-a100-sxm4-40gb-attributes
    - name: common-mig-7g.40gb-mock-nvidia-a100-sxm4-40gb
    - name: specific-gpu-0-mig-attributes
    - name: memory-slices-0-7

Flattened spec:
devices:
- name: gpu-0
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      index:
        int: 0
      minor:
        int: 0
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
      type:
        string: gpu
      uuid:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
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
        quantity: 40Gi
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
- name: gpu-0-mig-1g.10gb-0-1
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
        quantity: 9856Mi
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
- name: gpu-0-mig-1g.10gb-2-3
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
        quantity: 9856Mi
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
- name: gpu-0-mig-1g.10gb-4-5
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
        quantity: 9856Mi
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
- name: gpu-0-mig-1g.10gb-6-7
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
        quantity: 9856Mi
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
- name: gpu-0-mig-1g.5gb-0
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
- name: gpu-0-mig-1g.5gb-1
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
- name: gpu-0-mig-1g.5gb-2
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
- name: gpu-0-mig-1g.5gb-3
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
- name: gpu-0-mig-1g.5gb-4
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
- name: gpu-0-mig-1g.5gb-5
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
- name: gpu-0-mig-1g.5gb-6
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
- name: gpu-0-mig-1g.5gb-me-0
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
- name: gpu-0-mig-1g.5gb-me-1
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
- name: gpu-0-mig-1g.5gb-me-2
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
- name: gpu-0-mig-1g.5gb-me-3
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
- name: gpu-0-mig-1g.5gb-me-4
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
- name: gpu-0-mig-1g.5gb-me-5
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
- name: gpu-0-mig-1g.5gb-me-6
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
- name: gpu-0-mig-2g.10gb-0-1
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
        quantity: 9856Mi
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
- name: gpu-0-mig-2g.10gb-2-3
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
        quantity: 9856Mi
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
- name: gpu-0-mig-2g.10gb-4-5
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
        quantity: 9856Mi
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
- name: gpu-0-mig-3g.20gb-0-3
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
        quantity: 19968Mi
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
- name: gpu-0-mig-3g.20gb-4-7
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
        quantity: 19968Mi
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
- name: gpu-0-mig-4g.20gb-0-3
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
        quantity: 19968Mi
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
- name: gpu-0-mig-7g.40gb-0-7
  partitionable:
    attributes:
      architecture:
        string: Ampere
      brand:
        string: Nvidia
      cudaComputeCapability:
        string: "8.0"
      cudaDriverVersion:
        version: "12.4"
      driverVersion:
        version: 550.54.15
      parentIndex:
        int: 0
      parentMinor:
        int: 0
      parentUUID:
        string: GPU-ea456537-2782-4f6a-8907-fbd43b999173
      productName:
        string: Mock NVIDIA A100-SXM4-40GB
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
        quantity: 40192Mi
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
