# dra-evolution

The [k8srm-prototype](../k8srm-prototype/README.md) is an attempt to derive a
new API for device management from scratch. The API in this directory is taking
the opposite approach: it incorporates ideas from the prototype into the 1.30
DRA API. For some problems it picks a different approach.
To compare YAML files, something like this can be used:
```
diff -C2 ../k8srm-prototype/testdata/classes.yaml <(sed -e 's;resource.k8s.io/v1alpha2;devmgmtproto.k8s.io/v1alpha1;' -e 's/ResourceClass/DeviceClass/' testdata/classes.yaml)
```

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
is really a brainstorming exercise and under active development.

## Pod Spec

This prototype changes the `PodSpec` a little from how it is in DRA in 1.30.

As 1.30, the `PodSpec` has a list of named sources. The sources are structs that
could contain either a claim name or a template name. The names are used to
associate individual claims with containers.

Each claim may contain multiple request for different devices. Containers can
also be associated with individual requests inside a claim.

Allocating multiple devices per claim allows specifying constraints for a set
of devices, like "some attribute has to be the same". Long-term, it would be
good to allow such constraints also across claims when a pod references more
than one, but that would imply extending the `PodSpec` with complex fields
where we are not sure yet what they need to look like. Therefore these
constraints are currently limited to claims. This limitation may be
removed once constraints are stable enough to be included in the `PodSpec`.

These `PodSpec` Go types can be seen in [pod_types.go](pkg/api/pod_types.go).

## Types

Types are divided into "claim" types, which form the UX, "capacity" types which
are populated by drivers, and "allocation types" which are used to capture the
results of scheduling. Allocation types are really just the status types of the
claim types.

Claim and allocation types are found in [claim_types.go](pkg/api/claim_types.go);
individual types and fields are described in detail there in the comments.
Capacity types are in [capacity_types.go](pkg/api/capacity_types.go). A quota
mechanism is defined in [quota_types.go](pkg/api/quota_types.go).

Vendors and administrators create `DeviceClass` resources to pre-configure
various options for requests in claims. Such a class contains:
- configuration for a device, potentially including options that
  only an administrator may set
- device requirements which select device instances that match the intended
  semantic of the class ("give me a GPU")

Classes are not necessarily associated with a single vendor. Whether they are
depends on how the requirements in them are defined.

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

