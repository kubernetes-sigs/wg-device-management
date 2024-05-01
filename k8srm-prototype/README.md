# k8srm-prototype

For more background, please see this document, though it is not yet up to date with the latest in this repo:
- [Revisiting Kubernetes Resource Model](https://docs.google.com/document/d/1Xy8HpGATxgA2S5tuFWNtaarw5KT8D2mj1F4AP1wg6dM/edit?usp=sharing).


## Building

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

## Types

Types are divided into "claim" types, which form the UX, "capacity" types which
are populated by drivers, and "allocation types" which are used to capture the
results of scheduling.

Claim types are found in [claim_types.go](pkg/api/claim_types.go);
individual types and fields are described in detail there in the comments.

When making a claim for a device (or set of devices), the user may either
specify a device managed by a specific driver, or they may specify an arbitrary
"device type"; for example, "sriov-nic". Individual drivers register with the
control plan and publish the device types which they handle using the
cluster-scoped `DeviceDriver` resource. Examples:
[drivers.yaml](testdata/drivers.yaml).

Vendors and administrators create `DeviceClass` resources to pre-configure
various options for claims. DeviceClass resources must refer to a specific
DeviceType, and may refer to a specific DeviceDriver. Examples:
[classes.yaml](testdata/classes.yaml).

Users create `DeviceClaim` resources, which must refer to a specific
DeviceClass resource. The rest of the DeviceClaim spec can be used to further
specify configuration and selection criteria for the set of desired devices.

DeviceClaim resources are embedded or referenced from the PodSpec, much like
volumes. We should discuss whether we need a separate `DeviceClaimTemplate`
class or if we can simply refer to a DeviceClaim as if it were a temlate.
Probably the separate resource type is cleaner. Examples may be found in the
`testdata` directory in files starting with `pod-`; e.g.,
[pod-template-foozer-single.yaml](testdata/pod-template-foozer-single.yaml).

## Examples

There are some examples in [schedule_test.go](pkg/schedule/schedule_test.go). If
you run the schedule test you can also get some output as to what it is doing:

```console
k8srm-prototype$ cd pkg/schedule/
schedule$ go test

=== TEST single by driver

ALLOCATIONS
-----------
- deviceCount: 1
  devicePoolName: shape-zero-00-foozer-00

NODE RESULTS
------------
shape-three-01: could not satisfy these claims: myclaim
shape-zero-00: satisfied all claims with score 100
shape-zero-01: satisfied all claims with score 100
shape-three-00: could not satisfy these claims: myclaim

=== DONE single by driver

...snipped...
```

Or for even more details, including all options considered and their various
scores:

```console
schedule$ VERBOSE=y go test

=== TEST single by driver

ALLOCATIONS
-----------
- deviceCount: 1
  devicePoolName: shape-zero-00-foozer-00

NODE RESULTS
------------
- DeviceClaimResults:
  - best: -1
    claimName: myclaim
    ignoredPools:
    - deviceCount: 0
      failureReason: claim and pool driver do not match
      poolName: shape-three-00-barzer-00
    - deviceCount: 0
      failureReason: claim and pool driver do not match
      poolName: shape-three-00-barzer-01
    - deviceCount: 0
      failureReason: claim and pool driver do not match
      poolName: shape-three-00-barzer-02
    - deviceCount: 0
      failureReason: claim and pool driver do not match
      poolName: shape-three-00-barzer-03
    poolSetResults: null
  NodeName: shape-three-00
- DeviceClaimResults:
  - best: -1
    claimName: myclaim
    ignoredPools:
    - deviceCount: 0
      failureReason: claim and pool driver do not match
      poolName: shape-three-01-barzer-00
    - deviceCount: 0
      failureReason: claim and pool driver do not match
      poolName: shape-three-01-barzer-01
    - deviceCount: 0
      failureReason: claim and pool driver do not match
      poolName: shape-three-01-barzer-02
    - deviceCount: 0
      failureReason: claim and pool driver do not match
      poolName: shape-three-01-barzer-03
    poolSetResults: null
  NodeName: shape-three-01
- DeviceClaimResults:
  - best: 0
    claimName: myclaim
    poolSetResults:
    - poolResults:
      - deviceCount: 1
        poolName: shape-zero-00-foozer-00
      score: 100
  NodeName: shape-zero-00
- DeviceClaimResults:
  - best: 0
    claimName: myclaim
    poolSetResults:
    - poolResults:
      - deviceCount: 1
        poolName: shape-zero-01-foozer-00
      score: 100
  NodeName: shape-zero-01


=== DONE single by driver

...snipped...
```

