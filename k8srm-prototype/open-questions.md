# Open Questions

- Will the inlined pointer be OK for the DeviceClaimInstance one-of?
- The way management pods are handled is clunky at best. A separate top-level
  type still seems cleaner.
- Should classes be just (or nearly) as expressive as claims? For example,
  should the `oneOf` logic be something that can be encapsulated in a class?
- Do we need BlockSize now? What about IntRange?
- We need to standardize some attribute names, otherwise `matchAttributes` will
  not be useful cross-driver.
- Is there a use case for `distinctAttributes`? If so, what should we name
  `match` and `distinct` attributes? `match` and `mismatch` are used in
  `PodAffinityTerm` but I am not sure they mean the same thing as here.
- A claim template is just an ObjectMeta + ClaimSpec. Can we just re-use
  DeviceClaim for that, or do we need a separate DeviceClaimTemplate type? If we
  re-use claim, how do we differentiate between a "real" claim and a "template"
  claim?
- We may be able to keep `PodSpec.ResourceClaims` as it is named, but add some
  values to the oneOf for DeviceClaim, DeviceClaimTemplate. We would then leave
  `container.resources.claims` as is, too.
- We need limits on the size of the ClassConfigs and ClaimConfigs lists.
- We need to consider how this interacts with local volumes, for example NVME.
  See PersistentVolumeMode. For example, "I need a GPU (one-of: [1x A100, 2x L4])
  and an NVME SSD on the same PCIE.
- More generally we need to exercise various multi-claim, all-of, one-of type of
  use cases and make sure this is sufficiently expressive.
- Do we need PodNames in the claim status to show what Pods are using the claim,
  or can we use ownerRefs?
- Is DevicePool the right concept? It seems that we have either simple devices
  or we have complex devices with shared resources (like MIG devices). How often
  will it be the latter? SR-IOV NICs with multiple PFs look a lot like
  partitioned devices. Do we really need to normalize to pools or can we just
  make each entry a separate device, which may be a simple device or may be a
  complex, partitioned device? Related discussions and some points pulled from
  them:
  - https://github.com/kubernetes-sigs/wg-device-management/pull/5/files#r1587245528
  - https://github.com/kubernetes-sigs/wg-device-management/pull/5#discussion_r1591620848
  - Are most cases going to use some form of `SharedResources` or do we think
    there are lots of cases where it's just "whole PCIe" card?
  - Are we trying to come up with a new abstraction that is already part of the
    PCIe structure of functions? Can we model based on that instead?
  - If we have to do things like "pf-0-vfs" and "memory-slice-0", are we missing
    some fundamental part of the model around shared resources - like groups of
    shared resources a la Kevin's named model extension (I thought that was too
    complex...)?
  - Where it gets weird though is when a vendor has a mix of card types to
    support.
    * Where does one draw the boundaries around each `DevicePool`?
    * Do they put all of their simple devices in a single `DevicePool` alongside
      all partitionable devices in their own individual `DevicePool`s?
    * Why is there a "pool" boundary at all if there is no real meaning tied to
      it?
    * In the case of simple devices, is it only there so as to avoid having lots
      of separate `Device` objects in the API server?
    * If we do want to support both simple and partitionable devices in a single
      `DevicePool`, do we need to add one level of embedding as @thockin
      proposed? An alternative of this is to go back to my concept of having a
      list of named `SharedResourceGroups`  (as opposed to a single
      `SharedResource` section), forcing individual devices to refer back to the
      name of the `SharedResourceGroup` they pull a given shared resource from.
    * What happens as the size of these device pools grow and we hit the limit
      of a single API server object? Where do we draw the boundary then?
- Do we need to deal with cross-pod (and cross-node) "linked" claims ala this
  [discussion](https://github.com/kubernetes-sigs/wg-device-management/pull/5#pullrequestreview-2035165945)?
  Or can that be handled by a higher-level workload controller that understands
  the workload, generates a claim, and then re-uses that claim name between
  pods?
- Does
  [this](https://github.com/kubernetes-sigs/wg-device-management/pull/5#discussion_r1591489552)
  work better for allOf/oneOf construct in the claims?
- Should we put in well-known status fields for certain types of devices (like
  IP for network devices). See this
  [discussion](https://github.com/kubernetes-sigs/wg-device-management/pull/5#discussion_r1589739514).
- Better document "user-mediated sharing" versus "platform-mediated" sharing.
  See this
  [discussion](https://github.com/kubernetes-sigs/wg-device-management/pull/5#discussion_r1591621730).
- Capacity model naming around shared resource consumption and claim resources
  provided is not great.
  - ~See
    [here](https://github.com/kubernetes-sigs/wg-device-management/pull/5#discussion_r1591623761)~
    (addressed)
  - See
    [here](https://github.com/kubernetes-sigs/wg-device-management/pull/5#discussion_r1591623874)
- Having DevicePoolName in the claim status field makes pools immutable which is
  bad. We need a different solution, maybe device UUIDs? See this
  [discussion](https://github.com/kubernetes-sigs/wg-device-management/pull/5#discussion_r1591664614).
- We handle compound devices in this model by making the user select both
  devices with a `matchAttributes` constraint. It would be better if we could
  handle that directly on the capacity side. This may be possible if we have a
  hook in our driver framework to allow a second driver to contribute to the
  device object, perhaps adding something into the SharedResources? The idea
  would be that if the cloud provider (or other node provider) understands that
  Device A goes with Device B, then they could leverage this hook to
  "subordinate" device B under device A, even if those are handled by different
  drivers.
