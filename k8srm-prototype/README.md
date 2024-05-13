# k8srm-prototype

For more background, please see this document, though it is not yet up to date
with the latest in this repo:
- [Revisiting Kubernetes Resource
  Model](https://docs.google.com/document/d/1Xy8HpGATxgA2S5tuFWNtaarw5KT8D2mj1F4AP1wg6dM/edit?usp=sharing).

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

See [dra-evolution](../dra-evolution/README.md) for an attempt to integrate these ideas
into DRA.

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

So, we tweak the `PodSpec` a bit from 1.30, such that, instead of a list of
named sources, with each source being a oneOf, we instead have a single
`DeviceClaims` oneOf in the `PodSpec`. This oneOf could be:
- A list of named sources, where sources are limited to a simple "class" name
  (ie, not a list of oneOfs, just a list of simple structs).
- A template struct, which consists of ObjectMeta + a claim name.
- A claim name.

Additionally we move the container association from
`spec.containers[*].resources.claims` to `spec.containers[*].devices`.

The first form of the of the `DeviceClaims` oneOf allows for our simplest of use
cases to be very simple to express, without creating a secondary object to which
we must then refer. So, the equivalent of the 1.30 YAML above would be:

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
    - name: gpu
  deviceClaims:
    devices:
    - name: gpu
      class: example.com-foozer
```

Each entry in `spec.deviceClaims.devices` is just a name/class pair, but in fact
serves as a template to generate claims that exist with the lifecycle of the
pod. We may want to add `ObjectMeta` here as well, since it is behaving as a
template, to allow setting labels, etc.

The second form of `DeviceClaims` is a single struct with an ObjectMeta, and a
claim name. The key with this form is that it is not *list* of named objects.
Instead, it is a reference to a single claim object, and the named entries are
*inside* the referenced object. This is to avoid a two-key mount in the
`spec.containers[*].devices` entry. If that's not important, then we can tweak
this a bit. In any case, this form allows claims which follow the lifecycle of
the pod, similar to the first form. Since a top-level API claim spec can can
contain multiple claim instances, this should be equally as expressive as if we
included `matchAttributes` in the `PodSpec`, without having to do so.

The third form of `DeviceClaims` is just a string; it is a claim name and allows
the user to share a pre-provisioned claim between pods.

Given that the first and second forms both have a template-like structure, we
may want to combine them and use two-key indexing in the mounts. If we do so, we
still want the direct specification of the class, so that the most common case
does not need separate object just to reference a class.

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

Vendors and administrators create `DeviceClass` resources to pre-configure
various options for claims. DeviceClass resources come in two varieties:
- Ordinary or "leaf" classes that represent devices managed by a specific
  driver, along with some optional selection constraints and configuration.
- "Meta" or "Group" or "Aggregate" or "Composition" classes that use a label
  selector to identify a *set* of leaf classes. This allows a claim to be
  satistfied by one of many classes.

Example classes are in [classes.yaml](testdata/classes.yaml).

Example pod definitions can be found in the `pod-*.yaml` and `two-pods-*.yaml`
files in [testdata](testdata).

Drivers publish capacity via `DevicePool` resources. Examples may be found in
the `pools-*.yaml` files in [testdata](testdata).

## Building

Soon we will add back in scheduling algorithms so people can see how these would
work. But right now, the actual code that runs is just for generating sample
capacity data.

Just run `make`, it will build everything.

```console
k8srm-prototype$ make
gofmt -s -w .
go test ./...
?   	github.com/kubernetes-sigs/wg-device-management/k8srm-prototype/cmd/mock-apiserver	[no test files]
?   	github.com/kubernetes-sigs/wg-device-management/k8srm-prototype/cmd/schedule	[no test files]
?   	github.com/kubernetes-sigs/wg-device-management/k8srm-prototype/pkg/api	[no test files]
?   	github.com/kubernetes-sigs/wg-device-management/k8srm-prototype/pkg/gen	[no test files]
ok  	github.com/kubernetes-sigs/wg-device-management/k8srm-prototype/pkg/schedule	(cached)
cd cmd/schedule && go build
cd cmd/mock-apiserver && go build
```

## Mock APIServer

This repo includes a crude mock API server that can be loaded with the examples
and used to try out scheduling (WIP). It will spit out some errors but you can
ignore them.

```console
k8srm-prototype$ ./cmd/mock-apiserver/mock-apiserver
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
k8srm-prototype$ kubectl --kubeconfig kubeconfig apply -f testdata/drivers.yaml
devicedriver.devmgmtproto.k8s.io/example.com-foozer created
devicedriver.devmgmtproto.k8s.io/example.com-barzer created
devicedriver.devmgmtproto.k8s.io/sriov-nic created
devicedriver.devmgmtproto.k8s.io/vlan created
k8srm-prototype$ kubectl --kubeconfig kubeconfig get devicedrivers
NAME                 AGE
example.com-foozer   2y112d
example.com-barzer   2y112d
sriov-nic            2y112d
vlan                 2y112d
k8srm-prototype$
```

## `schedule` CLI

This is CLI that represents what the scheduler and/or other controllers will do
in a real system. That is, it will take a pod and a list of nodes and schedule
the pod to the node, taking into account the device claims and writing the
results to the various status fields. This doesn't work right now, it needs to
be updated for the most recent changes.

