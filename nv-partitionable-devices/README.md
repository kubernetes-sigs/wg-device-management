# nvk8s-resourcemodel

This repo is meant as a staging ground for defining extensions to the
Kubernetes structured parameters model as it pertains to NVIDIA hardware.
All of the running examples use a Mock DGXA100 server to feed its input, so the
results it generates are comparable to what we would see in real hardware.

## Proposed Extensions

Below is a summary of the currently proposed extensions. With these extensions
in place we are able to enable fully dynamic MIG with the possibility for the
scheduler to do intelligent allocation to avoid fragmentation.

```diff
 // NamedResourcesResources is used in ResourceModel.
 // +k8s:deepcopy-gen=true
 type NamedResourcesResources struct {
 	// The list of all individual resources instances currently available.
 	//
 	// +listType=atomic
 	Instances []NamedResourcesInstance `json:"instances" protobuf:"bytes,1,name=instances"`
 
+	// The list of all shared resources limits that are referenced by one or
+	// more instances.
+	//
+	// +listType=atomic
+	// +optional
+	SharedLimits []NamedResourcesSharedResourceGroup `json:"sharedLimits,omitempty" protobuf:"bytes,2,opt,name=sharedLimits"`
 }
 
 // NamedResourcesInstance represents one individual hardware instance that can be selected based
 // on its attributes.
 // +k8s:deepcopy-gen=true
 type NamedResourcesInstance struct {
 	// Name is unique identifier among all resource instances managed by
 	// the driver on the node. It must be a DNS subdomain.
 	Name string `json:"name" protobuf:"bytes,1,name=name"`
 
 	// Attributes defines the attributes of this resource instance.
 	// The name of each attribute must be unique.
 	//
 	// +listType=atomic
 	// +optional
 	Attributes []NamedResourcesAttribute `json:"attributes,omitempty" protobuf:"bytes,2,opt,name=attributes"`
 
+	// Resources defines the set of resources this instance consumes when
+	// allocated.
+	//
+	// +listType=atomic
+	// +optional
+	Resources []NamedResourcesSharedResourceGroup `json:"resources,omitempty" protobuf:"bytes,3,opt,name=resources"`
 }
+
+// NamedResourcesSharedResource represents a shared resource that is consumable by a top-level resource when allocated.
+// +k8s:deepcopy-gen=true
+type NamedResourcesSharedResource struct {
+	// Name is the name of the resource represented by this shared resource.
+	// It must be a DNS subdomain.
+	Name string `json:"name" protobuf:"bytes,1,name=name"`
+
+	// NamedResourcesAttributeValue is an embedded type representing the actual value of the shared resource.
+	NamedResourcesSharedResourceValue `json:",inline" protobuf:"bytes,2,opt,name=value"`
+}
+
+// NamedResourcesSharedResourceValue represents the value of a shared resource.
+// NamedResourcesSharedResourceValue must have one and only one field set.
+// +k8s:deepcopy-gen=true
+type NamedResourcesSharedResourceValue struct {
+	// QuantityValue is a quantity.
+	QuantityValue *resource.Quantity `json:"quantity,omitempty" protobuf:"bytes,1,opt,name=quantity"`
+
+	// IntRangeValue is a range of 64-bit integers.
+	IntRangeValue *intrange.IntRange `json:"intRange,omitempty" protobuf:"varint,2,rep,name=intRange"`
+}
+
+// NamedResourcesSharedResourceGroup represents a named group of shared resources.
+// +k8s:deepcopy-gen=true
+type NamedResourcesSharedResourceGroup struct {
+	// Name is unique identifier among all resource groups managed by
+	// the driver on the node. It must be a DNS subdomain.
+	Name string `json:"name" protobuf:"bytes,1,name=name"`
+
+	// Items represents the list of all resources in the shared resource group.
+	//
+	// +listType=atomic
+	// +optional
+	Items []NamedResourcesSharedResource `json:"items,omitempty" protobuf:"bytes,2,opt,name=items"`
+}
```

The basic idea is the following:

* Introduce a top-level `SharedLimits` field in the `NamedResources` struct.
  This field is a list of a new type called `NamedResourcesSharedResourceGroups`,
  each of which defines a named collection of resources and their limits. In
  the case of NVIDIA GPUs, there will be one `NamedResourcesSharedResourcesGroup`
  per GPU on the machine, whose name is `gpu-%d-shared-resources` and has
  limits set for all of the "sub-resources" that can be consumed by different
  MIG devices of that GPU.

* Introduce a `Resources` field in the `NamedResourcesInstance` struct.
  This field is also a list of `NamedResourcesSharedResourcesGroups`, but
  instead of defining a _limit_ of available resources as before, it defines
  the actual _set_ of resources that would be consumed by this instance if it
  were to be allocated.  In the case of NVIDI GPUs, we declare one instance for
  every possible MIG device or full GPU that could be allocated from the
  machine.  In turn, each of these instances declare resources for the total
  amount of memory they consume, the discrete set of memory slices they consume
  (to help avoid fragmentation later), as well as the number of Jpeg, Ofa,
  Decoder engines, etc. that will be consumed if they get allocated.

With these simple additions in place, we have everyting we need to support
fully dynamic partitoining of GPUs with MIG. In essence, the scheduler is now
able to track the consumption of any shared (i.e. overlapping) resources and
ensure that their limits are not exceeded when making a scheduling decision.
In the case of NVIDIA GPUs, this means that overlapping MIG devices (as well as
the full GPUs they are part of) can be considered by the scheduler
independently, without it needing to understand the exact device hierarchy of
the hardware itself.

## Examples / Proof of Concept

There are two different commands available:

**`print-model`**:
    This command will print out two resource models for GPU 0 on the Mock
    DGXA100 server. The first is defined using the `NamedResources` model in
    Kubernetes v1.30. The second is with a set of proposed extensions to
    support dynamic partitioning of MIG devices.

**`print-possible-allocations`**:
    This command will print out all possible allocations of the resources
    declared for GPU 0 on the Mock DGXA100 server using the proposed extensions
    to the `NamedResources` model. This is meant to demonstrate how the
    scheduler can use these extensions to dynamically allocate non-overlapping
    devices partitioned from the same piece of hardware.

To run these commands, invoke `make <cmd-name>`.

The output of running each command can be seen below:

```console
$ make print-model

######## NamedResourceModel v1.30 ########
namedResources:
  instances:
  - attributes:
    - int: 0
      name: minor
    - int: 0
      name: index
    - name: uuid
      string: GPU-46a3c1c9-604b-4acd-8fed-35724140fc17
    - name: memory
      quantity: 40Gi
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    - bool: true
      name: mig-capable
    - bool: false
      name: mig-enabled
    name: gpu-0

######## Proposed NamedResourceModel v1.31 ########
namedResources:
  instances:
  - attributes:
    - int: 0
      name: minor
    - int: 0
      name: index
    - name: uuid
      string: GPU-46a3c1c9-604b-4acd-8fed-35724140fc17
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    - bool: true
      name: mig-capable
    - bool: false
      name: mig-enabled
    name: gpu-0
    resources:
    - items:
      - name: memory
        quantity: 40Gi
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.5gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.5gb-0
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "0"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 4864Mi
      - intRange: "0"
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.5gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.5gb-1
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "0"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 4864Mi
      - intRange: "1"
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.5gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.5gb-2
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "0"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 4864Mi
      - intRange: "2"
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.5gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.5gb-3
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "0"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 4864Mi
      - intRange: "3"
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.5gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.5gb-4
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "0"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 4864Mi
      - intRange: "4"
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.5gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.5gb-5
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "0"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 4864Mi
      - intRange: "5"
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.5gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.5gb-6
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "0"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 4864Mi
      - intRange: "6"
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 2g.10gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-2g.10gb-0
    resources:
    - items:
      - name: multiprocessors
        quantity: "28"
      - name: copy-engines
        quantity: "2"
      - name: decoders
        quantity: "1"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 9856Mi
      - intRange: 0-1
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 2g.10gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-2g.10gb-2
    resources:
    - items:
      - name: multiprocessors
        quantity: "28"
      - name: copy-engines
        quantity: "2"
      - name: decoders
        quantity: "1"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 9856Mi
      - intRange: 2-3
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 2g.10gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-2g.10gb-4
    resources:
    - items:
      - name: multiprocessors
        quantity: "28"
      - name: copy-engines
        quantity: "2"
      - name: decoders
        quantity: "1"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 9856Mi
      - intRange: 4-5
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 3g.20gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-3g.20gb-0
    resources:
    - items:
      - name: multiprocessors
        quantity: "42"
      - name: copy-engines
        quantity: "3"
      - name: decoders
        quantity: "2"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 19968Mi
      - intRange: 0-3
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 3g.20gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-3g.20gb-4
    resources:
    - items:
      - name: multiprocessors
        quantity: "42"
      - name: copy-engines
        quantity: "3"
      - name: decoders
        quantity: "2"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 19968Mi
      - intRange: 4-7
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 4g.20gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-4g.20gb-0
    resources:
    - items:
      - name: multiprocessors
        quantity: "56"
      - name: copy-engines
        quantity: "4"
      - name: decoders
        quantity: "2"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 19968Mi
      - intRange: 0-3
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 7g.40gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-7g.40gb-0
    resources:
    - items:
      - name: multiprocessors
        quantity: "98"
      - name: copy-engines
        quantity: "7"
      - name: decoders
        quantity: "5"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "1"
      - name: ofa-engines
        quantity: "1"
      - name: memory
        quantity: 40192Mi
      - intRange: 0-7
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.5gb+me
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.5gb-me-0
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "1"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "1"
      - name: ofa-engines
        quantity: "1"
      - name: memory
        quantity: 4864Mi
      - intRange: "0"
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.5gb+me
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.5gb-me-1
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "1"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "1"
      - name: ofa-engines
        quantity: "1"
      - name: memory
        quantity: 4864Mi
      - intRange: "1"
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.5gb+me
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.5gb-me-2
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "1"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "1"
      - name: ofa-engines
        quantity: "1"
      - name: memory
        quantity: 4864Mi
      - intRange: "2"
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.5gb+me
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.5gb-me-3
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "1"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "1"
      - name: ofa-engines
        quantity: "1"
      - name: memory
        quantity: 4864Mi
      - intRange: "3"
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.5gb+me
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.5gb-me-4
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "1"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "1"
      - name: ofa-engines
        quantity: "1"
      - name: memory
        quantity: 4864Mi
      - intRange: "4"
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.5gb+me
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.5gb-me-5
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "1"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "1"
      - name: ofa-engines
        quantity: "1"
      - name: memory
        quantity: 4864Mi
      - intRange: "5"
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.5gb+me
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.5gb-me-6
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "1"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "1"
      - name: ofa-engines
        quantity: "1"
      - name: memory
        quantity: 4864Mi
      - intRange: "6"
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.10gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.10gb-0
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "1"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 9856Mi
      - intRange: 0-1
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.10gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.10gb-2
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "1"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 9856Mi
      - intRange: 2-3
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.10gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.10gb-4
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "1"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 9856Mi
      - intRange: 4-5
        name: memory-slices
      name: gpu-0-shared-resources
  - attributes:
    - name: mig-profile
      string: 1g.10gb
    - name: product-name
      string: Mock NVIDIA A100-SXM4-40GB
    - name: brand
      string: Nvidia
    - name: architecture
      string: Ampere
    - name: cuda-compute-capability
      version: 8.0.0
    - name: driver-version
      version: 550.54.15
    - name: cuda-driver-version
      version: 12.4.0
    name: gpu-0-mig-1g.10gb-6
    resources:
    - items:
      - name: multiprocessors
        quantity: "14"
      - name: copy-engines
        quantity: "1"
      - name: decoders
        quantity: "1"
      - name: encoders
        quantity: "0"
      - name: jpeg-engines
        quantity: "0"
      - name: ofa-engines
        quantity: "0"
      - name: memory
        quantity: 9856Mi
      - intRange: 6-7
        name: memory-slices
      name: gpu-0-shared-resources
  sharedLimits:
  - items:
    - name: memory
      quantity: 40Gi
    - name: multiprocessors
      quantity: "98"
    - name: copy-engines
      quantity: "7"
    - name: decoders
      quantity: "5"
    - name: encoders
      quantity: "0"
    - name: jpeg-engines
      quantity: "1"
    - name: ofa-engines
      quantity: "1"
    - intRange: 0-7
      name: memory-slices
    name: gpu-0-shared-resources
```

```console
$ make print-possible-allocations

[gpu-0]
[gpu-0-mig-1g.5gb-0]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-3g.20gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-4g.20gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-3g.20gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-4g.20gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-3g.20gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-4g.20gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-3g.20gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-4g.20gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-3g.20gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-4g.20gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-3g.20gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-4g.20gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-6]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-3g.20gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-4g.20gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-3g.20gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-4g.20gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-2g.10gb-0]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-4]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-3g.20gb-0]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-4g.20gb-0]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-2g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-3g.20gb-0]
[gpu-0-mig-3g.20gb-0 gpu-0-mig-3g.20gb-4]
[gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-3g.20gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-3g.20gb-4]
[gpu-0-mig-3g.20gb-4 gpu-0-mig-4g.20gb-0]
[gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-3g.20gb-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-4g.20gb-0]
[gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-4g.20gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-7g.40gb-0]
[gpu-0-mig-1g.5gb-me-0]
[gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-1]
[gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-1 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-2]
[gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-3]
[gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-3 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-4]
[gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-5]
[gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-5 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.5gb-me-6]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.5gb-me-6 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.10gb-0]
[gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.10gb-0 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.10gb-2]
[gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.10gb-2 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.10gb-4]
[gpu-0-mig-1g.10gb-4 gpu-0-mig-1g.10gb-6]
[gpu-0-mig-1g.10gb-6]
```
