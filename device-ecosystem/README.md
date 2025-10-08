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
| NVIDIA  | [dra-driver](https://github.com/NVIDIA/k8s-dra-driver-gpu)
| AMD     | [dra-driver](https://github.com/ROCm/k8s-gpu-dra-driver)
| Intel   | [intel-resource-drivers-for-kubernetes](https://github.com/intel/intel-resource-drivers-for-kubernetes)
| Google (TPU)  | could not find
| FuriosaAI | could not find


### NVIDIA

#### NVIDIA DRA Driver

NVIDIA has the [dra-driver](https://github.com/NVIDIA/k8s-dra-driver-gpu).

As of August 2025, one can install the DRA-driver via a helm chart to support NVIDIA + DRA.

#### GPU Operator

NVIDIA has plans to bundle the dra-driver in the GPU Operator. https://github.com/NVIDIA/gpu-operator/pull/1541 is the best I can find for tracking this work. 

### Intel

#### Intel DRA Driver

Intel has [intel-resource-drivers-for-kubernetes](https://github.com/intel/intel-resource-drivers-for-kubernetes) which does seem to support DRA.

However, the device drivers are not yet GA according to [their github readme](https://github.com/intel/intel-resource-drivers-for-kubernetes?tab=readme-ov-file#intel-resource-drivers-for-kubernetes)

### AMD

#### AMD DRA Driver

AMD's GPU DRA driver is available here: https://github.com/ROCm/k8s-gpu-dra-driver

- Status: Experimental (alpha).
- For installation, requirements, demos, and examples, refer to the repository:
	- Installation & Developer Guide: https://github.com/ROCm/k8s-gpu-dra-driver/blob/main/docs/installation.md
	- Examples: https://github.com/ROCm/k8s-gpu-dra-driver/tree/main/example
	- Demo Guide: https://github.com/ROCm/k8s-gpu-dra-driver/blob/main/docs/demo.md

Planned: Integration with the [AMD GPU Operator](https://github.com/ROCm/gpu-operator) coming soon.

### Google

#### Google TPU Driver

I wasn't able to find this. TODO: maybe some help on google on status on this.

### FuriosaAI

#### FuriosaAI DRA Driver
FuriosaAI's DRA Driver has not been officially released yet.
It is expected to be released and open-sourced in 2025 Q4.
Engineering version and it's guide is available:
- https://hub.docker.com/r/furiosaai/furiosa-dra-driver
- https://github.com/furiosa-ai/furiosa-dra-driver-guide

#### FuriosaAI NPU Operator
It is currently under development and expected to be released in 2025 Q4.
