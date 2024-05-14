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

package podspec

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apiserver/pkg/cel/environment"

	"github.com/kubernetes-sigs/wg-device-management/dra-evolution/pkg/api"
	cel "github.com/kubernetes-sigs/wg-device-management/dra-evolution/pkg/schedule"
)

//go:embed *.yaml
var yamls embed.FS

func TestParse(t *testing.T) {
	files, err := yamls.ReadDir(".")
	require.NoError(t, err)

	scheme := runtime.NewScheme()
	require.NoError(t, api.AddToScheme(scheme))
	serializer := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme, scheme, json.SerializerOptions{Yaml: true, Pretty: true, Strict: true})

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			fh, err := yamls.Open(file.Name())
			require.NoError(t, err)
			content, err := io.ReadAll(fh)
			require.NoError(t, err)

			// Split at the "---" separator before working on
			// individual item. Only works for .yaml.
			items := bytes.Split(content, []byte("\n---"))
			if len(items) > 1 {
				for i, item := range items {
					if len(item) > 0 {
						t.Run(fmt.Sprintf("item_%d", i), func(t *testing.T) {
							testDecode(t, serializer, item)
						})
					}
				}
			} else {
				testDecode(t, serializer, content)
			}
		})
	}
}

func testDecode(t *testing.T, serializer *json.Serializer, content []byte) {
	obj, gvk, err := serializer.Decode(content, nil, nil)
	require.NoError(t, err)
	t.Logf("Got object %T = %s", obj, gvk)

	switch obj := obj.(type) {
	case *api.ResourceClass:
		if obj.Filter != nil {
			validateDeviceFilter(t, obj.Filter.Devices, "class.filter.devices")
		}
		validateRequests(t, obj.DefaultRequests, "class.defaultRequests")
	case *api.ResourceClaim:
		if obj.Spec != nil {
			validateResourceClaimSpec(t, *obj.Spec, "claim.spec")
		}
	case *api.ResourceClaimTemplate:
		validateResourceClaimSpec(t, obj.Spec.Spec, "claimTemplate.spec.spec")
	}
}

func validateDeviceFilter(t *testing.T, filter *api.DeviceFilter, path string) {
	if !assert.NotNil(t, filter, path) {
		return
	}
	if filter.Selector == "" {
		return
	}
	result := cel.Compiler.CompileCELExpression(filter.Selector, environment.StoredExpressions)
	assert.Nil(t, result.Error, path+".selector parse error")
}

func validateRequests(t *testing.T, requests []api.ResourceRequest, path string) {
	for i, request := range requests {
		if request.ResourceRequestDetail != nil &&
			len(request.OneOf) > 0 {
			t.Errorf("%s[%d]: requesting one device and oneOf are mutually exclusive", path, i)
		}
		if request.ResourceRequestDetail == nil &&
			len(request.OneOf) == 0 {
			t.Errorf("%s[%d]: must request one device or oneOf", path, i)
			continue
		}
		if request.ResourceRequestDetail != nil {
			validateRequest(t, request.ResourceRequestDetail, fmt.Sprintf("%s[%d]", path, i))
		}
		for e, request := range request.OneOf {
			validateRequest(t, &request, fmt.Sprintf("%s[%d].oneOf[%d]", path, i, e))
		}
	}
}

func validateRequest(t *testing.T, request *api.ResourceRequestDetail, path string) {
	validateDeviceRequest(t, request.Device, path+".device")
}

func validateDeviceRequest(t *testing.T, request *api.DeviceRequest, path string) {
	if !assert.NotNil(t, request, path) {
		return
	}
	if request.Selector == "" {
		return
	}
	result := cel.Compiler.CompileCELExpression(request.Selector, environment.StoredExpressions)
	assert.Nil(t, result.Error, path+".selector parse error")
}

func validateResourceClaimSpec(t *testing.T, claimSpec api.ResourceClaimSpec, path string) {
	validateRequests(t, claimSpec.Requests, path+".requests")
}
