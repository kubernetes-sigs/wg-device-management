# dra-evolution

The [k8srm-prototype](../k8srm-prototype/README.md) is an attempt to derive a
new API for device management from scratch. The API in this directory is taking
the opposite approach: it incorporates ideas from the prototype into the 1.30
DRA API. For some problems it picks a different approach. The following
comparison is provided for those who already know one or the other
approach. Everyone else should probably read the proposals first and then come
back here. The last column explains why dra-evolution takes this approach.

| Use case, problem | DRA 1.30 | k8srm-prototype | dra-evolution | rationale |
| --- | --- | --- | --- | --- |
Classes | required, provide admin-level config and the driver name | DeviceClass: required, selects vendor driver and one device | ResourceClass: optional, can be vendor-independent, adds configuration, selection criteria and default parameters. | This avoids a two-level selection mechanism for devices (first the class, then the instance). Making a class potentially as descriptive as a claim enables additional use cases, like a pre-defined set of different devices from different vendors.
Custom APIs with CRDs | Vendors convert CRDs into class or claim parameters. | CRDs only provide configuration, content gets copied during allocation by scheduler. | As in 1.30 minus class CRDs. Claim parameters usually get specified directly in the claim. The ResourceClaimSpecification (= former ResourceClaimParameters) type is only used when a CRD reference is involved. | It is unclear whether any approach that depends on core Kubernetes reading vendor CRDs will pass reviews. Once this is clarified, this aspect can be revisited.
Management access | only in "classic DRA" | Field for device, not in class, checked via v1.ResourceQuota during admission. | Field for device, can be set in class, checked via resource.k8s.io ResourceQuota during allocation. | Checking at admission time is too limited. Eventually we will need a quota system that is based on device attributes.
Pod-level claims | Flat list with each entry mapping to a claim or claim template. | undecided ? | Flat list with each entry mapping to a claim, claim template or class as short-hand for "create claim for this class". | Adding the short-hand simplifies usage in simple cases.
Container-level claim references | name from list | two-level (claim + device in claim) ? | one level (all devices in a claim), two-level (specific device in claim) | The two-level case is needed when using a single claim to do matching between different devices and then wanting a container to use only one of the devices.
Matching between devices | only in "classic DRA" | MatchAttributes in claim | MatchAttributes in claim | This solves a sub-set of the matching problem. A more general solution would be a CEL expression, but that needs more thought and would be harder to use, so providing a "simple" solution seems worthwhile. Matching across claims is not supported by either proposal. This can only be done by putting fields whose semantic might still need to evolve into a v1 API. After GA?
Alternative device sets ("give me X, otherwise Y and Z") | only in "classic DRA" | oneOf, allOf | oneOf | "oneOf" seems to be a common requirement that might warrant special treatment to provide a simple API. "allOf" can be handled by replicating requests at the claim level.
Scoring | only in "classic DRA" | none | none | Like matching, this needs to be defined for a claim, with all devices of a potential solution as input. This is a tough problem that already occurs for a single device (pick "smallest" GPU or "biggest"?) and can easily lead to combinatorial explosion.
Claim status | Only allocation | Allocation, per-plugin status | Only allocation | Kubelet writing data provided by plugins leads to the [version skew problem](https://github.com/kubernetes/kubernetes/issues/123699). This becomes even worse when that data is likely to change when new status fields get added. This needs more thought before we put anything into the API that depends on sorting out this implementation challenge.
Claim template | Separate type | Re-uses claim + object meta in pod spec | Separate type | Defining claims that will never be used as claims "feels" weird. They also show up in `kubectl get resourceclaims -A` as "unallocated", which could be confusing.
"Resource" vs. "device" | resource | device | resource at top level, device inside | Only some of the semantic defined in the prototype is specific to devices. Other parts (like creating claims from templates, deallocation) are generic. If we ever need to add support for some other kind of resource, we would have to duplicate the entire outer API and copy-and-paste the generic code (Go generics don't support accessing "common" fields unless we define interfaces for everything, also typed client-go, etc.).
Resource model | one, potentially others | only one | one, potentially others, but with simpler YAML structure | The API should be as simple and natural as possible, but we need to keep the ability to add future extensions.
Driver handling allocation | in "classic DRA" | none | in "classic DRA" | We are not going to handle all the advanced scheduling use cases that people have solved with custom DRA control plane controllers, not now and perhaps never. It's too early to drop "classic DRA".
Vendor configuration for multiple devices | vendor parameters in claim and class | none ? | vendor parameters in claim and class | Storing configuration that isn't specific to one device under one device feels like a workaround. In a "oneOf", that same configuration would have to be repeated for each device.
Partioning | only in "classic DRA" | SharedResources | not added yet, still uses "named resources" | For the sake of simplicity, the current proposal doesn't attempt to modify how instances are described.

## Overall Model

As a refresher (see the KEPs for the details), the scope of the overall DRA /
Device Management effort is to select, configure, and allocate devices, and then
attach them to pods and containers. "Devices" here typically means on-node
devices but there are use cases for networked devices as well as devices that
can be attached/detached at runtime.

The high-level model here is:
- Drivers, typically running on the node, publish information about the devices
  they manage to the control plane.
- A user can make "claims" in their `PodSpec`, requesting one or more devices
  based on their needs.
- The scheduler looks at available capacity and selects the possible options
  that can meet the user's needs, scores them, and allocates them.
- The allocation information, along with the appropriate configuration
  information, is sent to kubelet along with the other pod information, and
  kubelet passes it on to the appropriate on-node drivers.
- The drivers perform the necessary (usually privileged) on-node actions, and
  write the status back to the control plane via kubelet.

The scope *of this prototype* is to quickly iterate on possible APIs to meet the
needs of workload authors, device vendors, Kubernetes vendors, platform
administrators, and higher level components such as autoscalers and ecosystem
projects.

## Open Questions

The next few sections of this document describe a proposed model. Note that this
is really a brainstorming exercise and under active development. See the [open
questions](open-questions.md) document for some of the still under discussion
items.

We are also looking at how we might extend the existing 1.30 DRA model with some
of these ideas, rather than changing it out for these specific types.

## Pod Spec

This prototype changes the `PodSpec` a little from how it is in DRA in 1.30.

In 1.30, the `PodSpec` has a list of named sources. The sources are structs that
could contain either a claim name or a template name. The names are used to
associate individual claims with containers. The example below allocates a
single "foozer" device to the container in the pod.

```yaml
apiVersion: resource.k8s.io/v1alpha1
kind: ResourceClaimTemplate
metadata:
  name: foozer
  namespace: default
spec:
  spec:
    resourceClassName: example.com-foozer
---
apiVersion: v1
kind: Pod
metadata:
  name: foozer
  namespace: default
spec:
  containers:
  - image: registry.k8s.io/pause:3.6
    name: my-container
    resources:
      requests:
        cpu: 10m
        memory: 10Mi
      claims:
      - name: gpu
  resourceClaims:
  - name: gpu
    source:
      resourceClaimTemplate: foozer
```

In the prototype model, we are adding `matchAttributes` constraints to control
consistency within a selection of devices. In particular, we want to be able to
specify a `matchAttributes` constraint across two separate named sources, so
that we can ensure for example, a GPU chosen for one container is the same model
as one chosen for another container. This would imply we need `matchAttributes`
that apply across the list present in `PodSpec`. However, we don't want to put
things like `matchAttributes` into `PodSpec`, since it is already `v1`.
Therefore matching is limited to devices within a claim. This limitation may be
removed once matching is stable enough to be included in the `PodSpec`.

To support selecting a specific device from a claim for a container, a
`resources.devices` list gets added:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: foozer
  namespace: default
spec:
  containers:
  - image: registry.k8s.io/pause:3.6
    name: my-container
    resources:
      requests:
        cpu: 10m
        memory: 10Mi
      devices:
      - claimName: gpu
        deviceName: gpu-one
  - image: registry.k8s.io/pause:3.6
    name: my-container
    resources:
      requests:
        cpu: 10m
        memory: 10Mi
      devices:
      - claimName: gpu
        deviceName: gpu-two
  resourceClaims:
  - name: gpu
    source:
      resourceClaimTemplate: two-foozers
```

Resource classes are capable of describing everything that a user might put
into a claim. Therefore a simple claim or claim template might contain nothing
but a resource class name. For this simple case, a new `claimWithClassName` gets
added which creates such a claim. Here object meta is supported:

```yaml
  resourceClaims:
  - name: gpu
    source:
      forClass:
        className: two-foozers-class
        metdadata:
          labels:
            foo: bar
```

How devices are named inside this class needs to be part of the class
documentation if users are meant to have the ability to select specific devices
for their containers.

These `PodSpec` Go types can be seen in [podspec.go](testdata/podspec.go). This
is not the complete `PodSpec` but just the relevant parts of the 1.30 and
proposed versions.

## Types

Types are divided into "claim" types, which form the UX, "capacity" types which
are populated by drivers, and "allocation types" which are used to capture the
results of scheduling. Allocation types are really just the status types of the
claim types.

Claim and allocation types are found in [claim_types.go](pkg/api/claim_types.go);
individual types and fields are described in detail there in the comments.
Capacity types are in [capacity_types.go](pkg/api/capacity_types.go). A quota
mechanism is defined in [quota_types.go](pkg/api/quota_types.go).

Vendors and administrators create `ResourceClass` resources to pre-configure
various options for claims. Depending on what gets set in a class, users can:
- Ask for exactly the set of devices pre-defined in a class.
- Add additional configuration to their claim. This configuration is
  passed down to the driver as coming from an admin, so it may control
  options that normal users must not set themselves.
- Restrict the choice of devices via additional constraints.

Classes are not necessarily associated with a single vendor. Whether they are
depends on how the constraints in them are defined.

Example classes are in [classes.yaml](testdata/classes.yaml).

Example pod definitions can be found in the `pod-*.yaml` and `two-pods-*.yaml`
files in [testdata](testdata).

Drivers publish capacity via `ResourcePool` objects. Examples may be found in
the `pools-*.yaml` files in [testdata](testdata).

## Building

Soon we will add back in scheduling algorithms so people can see how these would
work. But right now, the actual code that runs is just for generating sample
capacity data.

Just run `make`, it will build everything.

```console
dra-evolution$ make
gofmt -s -w .
go test ./...
?   	github.com/kubernetes-sigs/wg-device-management/dra-evolution/cmd/mock-apiserver	[no test files]
?   	github.com/kubernetes-sigs/wg-device-management/dra-evolution/cmd/schedule	[no test files]
?   	github.com/kubernetes-sigs/wg-device-management/dra-evolution/pkg/api	[no test files]
?   	github.com/kubernetes-sigs/wg-device-management/dra-evolution/pkg/gen	[no test files]
ok  	github.com/kubernetes-sigs/wg-device-management/dra-evolution/pkg/schedule	(cached)
cd cmd/schedule && go build
cd cmd/mock-apiserver && go build
```

## Mock APIServer

This repo includes a crude mock API server that can be loaded with the examples
and used to try out scheduling (WIP). It will spit out some errors but you can
ignore them.

```console
dra-evolution$ ./cmd/mock-apiserver/mock-apiserver
W0422 13:20:21.238440 2062725 memorystorage.go:93] type info not known for apiextensions.k8s.io/v1, Kind=CustomResourceDefinition
W0422 13:20:21.238598 2062725 memorystorage.go:93] type info not known for apiregistration.k8s.io/v1, Kind=APIService
W0422 13:20:21.238639 2062725 memorystorage.go:267] type info not known for foozer.example.com/v1alpha1, Kind=FoozerConfig
W0422 13:20:21.238666 2062725 memorystorage.go:267] type info not known for devmgmtproto.k8s.io/v1alpha1, Kind=DeviceDriver
W0422 13:20:21.238685 2062725 memorystorage.go:267] type info not known for devmgmtproto.k8s.io/v1alpha1, Kind=DeviceClass
W0422 13:20:21.238700 2062725 memorystorage.go:267] type info not known for devmgmtproto.k8s.io/v1alpha1, Kind=DeviceClaim
W0422 13:20:21.238712 2062725 memorystorage.go:267] type info not known for devmgmtproto.k8s.io/v1alpha1, Kind=DevicePrivilegedClaim
W0422 13:20:21.238723 2062725 memorystorage.go:267] type info not known for devmgmtproto.k8s.io/v1alpha1, Kind=DevicePool
2024/04/22 13:20:21 addr =  [::]:55441
```

The included `kubeconfig` will access that server. For example:

```console
dra-evolution$ kubectl --kubeconfig kubeconfig apply -f testdata/drivers.yaml
devicedriver.devmgmtproto.k8s.io/example.com-foozer created
devicedriver.devmgmtproto.k8s.io/example.com-barzer created
devicedriver.devmgmtproto.k8s.io/sriov-nic created
devicedriver.devmgmtproto.k8s.io/vlan created
dra-evolution$ kubectl --kubeconfig kubeconfig get devicedrivers
NAME                 AGE
example.com-foozer   2y112d
example.com-barzer   2y112d
sriov-nic            2y112d
vlan                 2y112d
dra-evolution$
```

## `schedule` CLI

This is CLI that represents what the scheduler and/or other controllers will do
in a real system. That is, it will take a pod and a list of nodes and schedule
the pod to the node, taking into account the device claims and writing the
results to the various status fields. This doesn't work right now, it needs to
be updated for the most recent changes.

