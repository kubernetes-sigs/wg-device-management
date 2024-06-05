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
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/util/validation"
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
	if runtime.IsNotRegisteredError(err) {
		t.Skipf("YAML file has not been updated yet: %v", err)
	}
	require.NoError(t, err)
	t.Logf("Got object %T = %s", obj, gvk)

	switch obj := obj.(type) {
	case *api.DeviceClass:
		validateRequestRequirements(t, obj.Requirements, "class.requirements")
	case *api.ResourceClaim:
		validateResourceClaimSpec(t, obj.Spec, "claim.spec")
	case *api.ResourceClaimTemplate:
		validateResourceClaimSpec(t, obj.Spec.Spec, "claimTemplate.spec.spec")
	case *api.ResourcePool:
		for i, device := range obj.Spec.Devices {
			for e, attribute := range device.Attributes {
				validateAttributeName(t, &attribute.Name, fmt.Sprintf("resourcePool.devices[%d].attributes[%d]", i, e))
			}
		}
	}
}

func validateRequestRequirements(t *testing.T, requirements []api.Requirement, path string) {
	for i, requirement := range requirements {
		validateDeviceSelector(t, requirement.DeviceSelector, fmt.Sprintf("%s[%d].deviceSelector", path, i))
	}
}

func validateClaimConstraints(t *testing.T, requirements []api.Constraint, path string) {
	for i, requirement := range requirements {
		validateAttributeName(t, requirement.MatchAttribute, fmt.Sprintf("%s[%d].matchAttribute", path, i))
	}
}

func validateDeviceSelector(t *testing.T, deviceSelector *string, path string) {
	if !assert.NotNil(t, deviceSelector, path) {
		return
	}
	result := cel.Compiler.CompileCELExpression(*deviceSelector, environment.StoredExpressions)
	assert.Nil(t, result.Error, path+".selector parse error")
}

func validateRequests(t *testing.T, requests []api.Request, path string) {
	for i, request := range requests {
		// if request.ResourceRequestDetail != nil &&
		// 	len(request.OneOf) > 0 {
		// 	t.Errorf("%s[%d]: requesting one device and oneOf are mutually exclusive", path, i)
		// }
		if request.ResourceRequestDetail == nil /* && len(request.OneOf) == 0 */ {
			t.Errorf("%s[%d]: must request one device or oneOf", path, i)
			continue
		}
		if request.ResourceRequestDetail != nil {
			validateRequest(t, request.ResourceRequestDetail, fmt.Sprintf("%s[%d]", path, i))
		}
		// for e, request := range request.OneOf {
		// 	validateRequest(t, &request, fmt.Sprintf("%s[%d].oneOf[%d]", path, i, e))
		// }
	}
}

func validateRequest(t *testing.T, request *api.ResourceRequestDetail, path string) {
	validateRequestRequirements(t, request.Requirements, path+".requirements")
}

func validateResourceClaimSpec(t *testing.T, claimSpec api.ResourceClaimSpec, path string) {
	validateClaimConstraints(t, claimSpec.Constraints, path+".constraints")
	validateRequests(t, claimSpec.Requests, path+".requests")
}

const qnameCharFmt string = "[A-Za-z0-9]"
const qnameExtCharFmt string = "[A-Za-z0-9_]"
const qualifiedNameFmt string = "(" + qnameCharFmt + qnameExtCharFmt + "*)?" + qnameCharFmt
const qualifiedNameErrMsg string = "must consist of alphanumeric characters or '_', and must start and end with an alphanumeric character"
const attributeNameIdentifierMaxLength int = 63

var attributeNameIdentifierRegexp = regexp.MustCompile("^" + qualifiedNameFmt + "$")

func validateAttributeName(t *testing.T, namePtr *string, path string) {
	if !assert.NotNil(t, namePtr, path) {
		return
	}
	name := *namePtr

	parts := strings.Split(name, "/")
	var domain, identifier string
	switch len(parts) {
	case 1:
		identifier = parts[0]
	case 2:
		domain, identifier = parts[0], parts[1]
		if len(domain) == 0 {
			t.Errorf("%s: domain part %s", path, validation.EmptyError())
		} else {
			validateDriverName(t, domain, path)
		}
	default:
		t.Errorf("%s: an attribute name must consist of a C-style identifier with an optional DNS subdomain prefix and '/' (e.g. 'example.com/MyName')", path)
	}

	if len(identifier) == 0 {
		t.Errorf("%s: identifier part %s", path, validation.EmptyError())
	} else if len(identifier) > attributeNameIdentifierMaxLength {
		t.Errorf("%s: identifier part %s", path, validation.MaxLenError(attributeNameIdentifierMaxLength))
	}
	if !attributeNameIdentifierRegexp.MatchString(identifier) {
		t.Errorf("%s: identifier part %s", path, validation.RegexError(qualifiedNameErrMsg, qualifiedNameFmt, "MyName", "my.name", "123-abc"))
	}
}

func validateDriverName(t *testing.T, driverName string, path string) {
	if len(driverName) == 0 {
		t.Errorf("%s: required", path)
	}

	if len(driverName) > 63 {
		t.Errorf("%s: too long", path)
	}

	for _, msg := range validation.IsDNS1123Subdomain(strings.ToLower(driverName)) {
		t.Errorf("%s: %s", path, msg)
	}
}
