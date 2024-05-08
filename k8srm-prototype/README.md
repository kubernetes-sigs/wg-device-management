# k8srm-prototype

For more background, please see this document, though it is not yet up to date with the latest in this repo:
- [Revisiting Kubernetes Resource Model](https://docs.google.com/document/d/1Xy8HpGATxgA2S5tuFWNtaarw5KT8D2mj1F4AP1wg6dM/edit?usp=sharing).

## Overall Model

As a refresher (see the KEPs for the details), the scope of the overall DRA /
Device Management effort is to select, configure, and allocated devices, and
attach them to pods and containers. "Devices" here typically means on-node
devices but there are use cases for networked devices as well as devices that
can be attached/detached at runtime.

The high-level model here is:
- Drivers, typically running on the node, publish information about the devices
  they manage to the control plane.
- A user can make "claims" in their PodSpec, requesting one or more devices
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

## Types

Note: this is really a brainstorming exercise and under active development. See
the [notes and open questions](notes-and-open-questions.md) document for some of
the still under discussion items.

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

Users can make claims in their PodSpec in a few different ways, see
[podspec.go](testdata/podspec.go) for a description of the various ways.

Example pod definitions can be found in the `pod-*.yaml` files in
[testdata](testdata).

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

