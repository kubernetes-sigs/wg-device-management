# Device Ecosystem for DRA

DRA is going to stable in 1.34 of Kubernetes. As more device vendors provide support for DRA, we want to have a tracking page where people can do to see who supports DRA.

## Contributing

We expect the status of the DRA drivers and their corresponds operators to change over time. Please feel free to add implementations as they graduate.

This document is not official supported by Kubernetes and is mainly used by #wg-device-management to track where device vendors are in their integrations for DRA.

## DRA Features

One area that we want to track is not just DRA availability but also support of various DRA features as they graduate.

### Table

| Vendor  | DRA Driver
|---------|------------
| Kubernetes Eaxmple Driver  | [dra-example-driver](https://github.com/kubernetes-sigs/dra-example-driver)
| NVIDIA  | [k8s-dra-driver-gpu](https://github.com/NVIDIA/k8s-dra-driver-gpu)
| AMD     | could not find
| Intel   | [intel-resource-drivers-for-kubernetes](https://github.com/intel/intel-resource-drivers-for-kubernetes)
| Google (TPU)  | could not find


### Example DRA Driver
The Kubernetes example DRA driver is a vendor-neutral reference implementation that demonstrates how to create a custom driver for managing specialized hardware in a Kubernetes cluster. It provides a foundational template that developers can adapt to support their specific resources, such as GPUs, FPGAs, or other accelerators. This driver showcases the core logic for resource discovery, allocation, and lifecycle management within the DRA framework.

The example-dra-driver was created by abstracting away NVIDIA-specific logic to serve as a vendor-neutral template for the community. To demonstrate its capabilities without requiring any specific physical hardware, it manages a set of simulated resources referred to as 'mock' GPUs. This allows developers to understand and test the core DRA control flow before integrating it with their actual devices.


### NVIDIA

The NVIDIA DRA driver is designed to manage NVIDIA GPUs in a Kubernetes cluster, offering more flexible and dynamic allocation of GPU resources to workloads. It moves beyond the limitations of the traditional, count based approach (e.g. nvidia.com/gpu).

A key feature is the introduction of ComputeDomains, an abstraction created to manage and secure Multi-Node NVLink (MNNVL) connectivity for  multi-node workloads. Architecturally, the driver is composed of two distinct kubelet plugins that can be enabled independently: the gpu-kubelet-plugin for core GPU allocation and the compute-domain-kubelet-plugin for handling the NVLink fabric.

The [demo section](https://github.com/NVIDIA/k8s-dra-driver-gpu/tree/main/demo/specs) provides several examples of how to allocate GPU resources:

- GPU Sharing: Demonstrates how multiple containers within the same pod can share access to a single GPU.

- MIG Allocation: Shows how to request and deploy pods to specific MIG profiles, partitioning a physical GPU for different workloads.

- IMEX for Multi-Node NVLink (MNNVL): Provides an advanced MPI-based example of how to create a ComputeDomain that spans multiple nodes, allowing pods to communicate directly over a secure NVLink fabric

The driver can be [installed via Helm](https://github.com/NVIDIA/k8s-dra-driver-gpu/blob/main/README.md#installation) and will be integrated into the NVIDIA GPU Operator in the future.

### Intel

#### Intel DRA Driver

Intel has [intel-resource-drivers-for-kubernetes](https://github.com/intel/intel-resource-drivers-for-kubernetes) which does seem to support DRA.

However, the device drivers are not yet GA according to [their github readme](https://github.com/intel/intel-resource-drivers-for-kubernetes?tab=readme-ov-file#intel-resource-drivers-for-kubernetes)

### AMD

#### AMD DRA Driver

Searching in [AMD ROCm](https://github.com/ROCm) I cannot seem to find a DRA implementation yet.

TODO: maybe some help on AMD on status on this.

### Google

#### Google TPU Driver

I wasn't able to find this. TODO: maybe some help on google on status on this.
