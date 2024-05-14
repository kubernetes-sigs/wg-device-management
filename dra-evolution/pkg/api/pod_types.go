/*
Copyright 20224 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Today in podspec, we have:
//

// type PodSpec130 struct {
// 	// lots of stuff, ending with:

// 	ResourceClaims []PodResourceClaim
// }

// type PodResourceClaim struct {
// 	// Name uniquely identifies this resource claim inside the pod.
// 	// This must be a DNS_LABEL.
// 	Name string

// 	// Source describes where to find the ResourceClaim.
// 	Source ClaimSource
// }

// type ClaimSource struct {
// 	ResourceClaimName         *string
// 	ResourceClaimTemplateName *string
// }

// // Claims at the pod level get referenced in containers with:
// type ResourceRequirements struct {
// 	// lots of stuff, ending with:

// 	// Claims lists the names of resources, defined in spec.resourceClaims,
// 	// that are used by this container.
// 	//
// 	// This is an alpha field and requires enabling the
// 	// DynamicResourceAllocation feature gate.
// 	//
// 	// This field is immutable. It can only be set for containers.
// 	//
// 	// +featureGate=DynamicResourceAllocation
// 	// +optional
// 	Claims []ResourceClaim
// }

// // ResourceClaim references one entry in PodSpec.ResourceClaims.
// type ResourceClaim struct {
// 	// Name must match the name of one entry in pod.spec.resourceClaims of
// 	// the Pod where this field is used. It makes that resource available
// 	// inside a container.
// 	Name string
// }

// This proposal instead has:
//

// Pod is a collection of containers that can run on a host. This resource is created
// by clients and scheduled onto hosts.
type Pod struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec PodSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type PodSpec struct {
	// stuff...

	// List of containers belonging to the pod.
	// Containers cannot currently be added or removed.
	// There must be at least one container in a Pod.
	// Cannot be updated.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	Containers []Container `json:"containers" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=containers"`

	// ResourceClaims defines which ResourceClaims must be allocated
	// and reserved before the Pod is allowed to start. The resources
	// will be made available to those containers which consume them
	// by name.
	//
	// This is an alpha field and requires enabling the
	// DynamicResourceAllocation feature gate.
	//
	// This field is immutable.
	//
	// +patchMergeKey=name
	// +patchStrategy=merge,retainKeys
	// +listType=map
	// +listMapKey=name
	// +featureGate=DynamicResourceAllocation
	// +optional
	ResourceClaims []PodResourceClaim `json:"resourceClaims,omitempty" patchStrategy:"merge,retainKeys" patchMergeKey:"name" protobuf:"bytes,39,rep,name=resourceClaims"`
}

// A single application container that you want to run within a pod.
type Container struct {
	// Name of the container specified as a DNS_LABEL.
	// Each container in a pod must have a unique name (DNS_LABEL).
	// Cannot be updated.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Container image name.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// This field is optional to allow higher level config management to default or override
	// container images in workload controllers like Deployments and StatefulSets.
	// +optional
	Image string `json:"image,omitempty" protobuf:"bytes,2,opt,name=image"`
	// Entrypoint array. Not executed within a shell.
	// The container image's ENTRYPOINT is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced
	// to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will
	// produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless
	// of whether the variable exists or not. Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	// +optional
	// +listType=atomic
	Command []string `json:"command,omitempty" protobuf:"bytes,3,rep,name=command"`

	// stuff...

	// Compute Resources required by this container.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Resources ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,8,opt,name=resources"`
}

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements struct {
	// Limits describes the maximum amount of compute resources allowed.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Limits ResourceList `json:"limits,omitempty" protobuf:"bytes,1,rep,name=limits,casttype=ResourceList,castkey=ResourceName"`

	// Requests describes the minimum amount of compute resources required.
	// If Requests is omitted for a container, it defaults to Limits if that is explicitly specified,
	// otherwise to an implementation-defined value. Requests cannot exceed Limits.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Requests ResourceList `json:"requests,omitempty" protobuf:"bytes,2,rep,name=requests,casttype=ResourceList,castkey=ResourceName"`

	// Claims lists the names of resources, defined in spec.resourceClaims,
	// that are used by this container.
	//
	// This is an alpha field and requires enabling the
	// DynamicResourceAllocation feature gate.
	//
	// This field is immutable. It can only be set for containers.
	//
	// +listType=map
	// +listMapKey=name
	// +featureGate=DynamicResourceAllocation
	// +optional
	Claims []ResourceClaim `json:"claims,omitempty" protobuf:"bytes,3,opt,name=claims"`
}

// ResourceName is the name identifying various resources in a ResourceList.
type ResourceName string

// ResourceList is a set of (resource name, quantity) pairs.
type ResourceList map[ResourceName]resource.Quantity

type PodResourceClaim struct {
	// Name uniquely identifies this resource claim inside the pod.
	// This must be a DNS_LABEL.
	Name string

	// Source describes where to find the ResourceClaim.
	Source ClaimSource
}

type ClaimSource struct {
	ForClass                  *ResourceClaimForClass
	ResourceClaimName         *string
	ResourceClaimTemplateName *string
}

type ResourceClaimForClass struct {
	ClassName string

	// ObjectMeta may contain labels and annotations that will be copied into the PVC
	// when creating it. No other fields are allowed and will be rejected during
	// validation.
	// +optional
	metav1.ObjectMeta
}

// ResourceClaimDevice references specific devices inside a claim. If the named request
// is satisfied by allocating multiple devices, then all of those are matched.
type ResourceClaimDevice struct {
	// pod.spec.resourceClaims
	ClaimName string
	// claim.spec.requests
	RequestName string
}
