/*
Copyright 2022 The Kubernetes Authors.

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

package cel

import (
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apiserver/pkg/cel/environment"
	"k8s.io/klog/v2/ktesting"
	"k8s.io/utils/ptr"

	"github.com/kubernetes-sigs/wg-device-management/dra-evolution/pkg/api"
)

func TestCompile(t *testing.T) {
	for name, scenario := range map[string]struct {
		expression         string
		attributeSuffix    string
		attributes         []api.DeviceAttribute
		celSuffix          string
		expectCompileError string
		expectMatchError   string
		expectMatch        bool
	}{
		"true": {
			expression:  "true",
			expectMatch: true,
		},
		"false": {
			expression:  "false",
			expectMatch: false,
		},
		"syntax-error": {
			expression:         "?!",
			expectCompileError: "Syntax error",
		},
		"type-error": {
			expression:         `device.quantity["no-such-attr"]`,
			expectCompileError: "must evaluate to bool",
		},
		"runtime-error": {
			expression:       `device.quantity["no-such-attr"].isGreaterThan(quantity("0"))`,
			expectMatchError: "no such key: no-such-attr",
		},
		"quantity": {
			expression:  `device.quantity["name"].isGreaterThan(quantity("0"))`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{QuantityValue: ptr.To(resource.MustParse("1"))}}},
			expectMatch: true,
		},
		"bool": {
			expression:  `device.bool["name"]`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{BoolValue: ptr.To(true)}}},
			expectMatch: true,
		},
		"int": {
			expression:  `device.int["name"] > 0`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{IntValue: ptr.To(int64(1))}}},
			expectMatch: true,
		},
		"intslice": {
			expression:  `device.intslice["name"].isSorted() && device.intslice["name"].indexOf(3) == 2`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{IntSliceValue: &api.IntSlice{Ints: []int64{1, 2, 3}}}}},
			expectMatch: true,
		},
		"empty-intslice": {
			expression:  `size(device.intslice["name"]) == 0`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{IntSliceValue: &api.IntSlice{}}}},
			expectMatch: true,
		},
		"string": {
			expression:  `device.string["name"] == "fish"`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{StringValue: ptr.To("fish")}}},
			expectMatch: true,
		},
		"stringslice": {
			expression:  `device.stringslice["name"].isSorted() && device.stringslice["name"].indexOf("a") == 0`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{StringSliceValue: &api.StringSlice{Strings: []string{"a", "b", "c"}}}}},
			expectMatch: true,
		},
		"empty-stringslice": {
			expression:  `size(device.stringslice["name"]) == 0`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{StringSliceValue: &api.StringSlice{}}}},
			expectMatch: true,
		},
		"all": {
			expression: `device.quantity["quantity"].isGreaterThan(quantity("0")) &&
device.bool["bool"] &&
device.int["int"] > 0 &&
device.intslice["intslice"].isSorted() &&
device.string["string"] == "fish" &&
device.stringslice["stringslice"].isSorted()`,
			attributes: []api.DeviceAttribute{
				{Name: "quantity", DeviceAttributeValue: api.DeviceAttributeValue{QuantityValue: ptr.To(resource.MustParse("1"))}},
				{Name: "bool", DeviceAttributeValue: api.DeviceAttributeValue{BoolValue: ptr.To(true)}},
				{Name: "int", DeviceAttributeValue: api.DeviceAttributeValue{IntValue: ptr.To(int64(1))}},
				{Name: "intslice", DeviceAttributeValue: api.DeviceAttributeValue{IntSliceValue: &api.IntSlice{Ints: []int64{1, 2, 3}}}},
				{Name: "string", DeviceAttributeValue: api.DeviceAttributeValue{StringValue: ptr.To("fish")}},
				{Name: "stringslice", DeviceAttributeValue: api.DeviceAttributeValue{StringSliceValue: &api.StringSlice{Strings: []string{"a", "b", "c"}}}},
			},
			expectMatch: true,
		},
		"many": {
			expression: `device.bool["a"] && device.bool["b"] && device.bool["c"]`,
			attributes: []api.DeviceAttribute{
				{Name: "a", DeviceAttributeValue: api.DeviceAttributeValue{BoolValue: ptr.To(true)}},
				{Name: "b", DeviceAttributeValue: api.DeviceAttributeValue{BoolValue: ptr.To(true)}},
				{Name: "c", DeviceAttributeValue: api.DeviceAttributeValue{BoolValue: ptr.To(true)}},
			},
			expectMatch: true,
		},
		"attributeSuffix": {
			expression:      `device.bool["name.k8s.io"]`,
			attributeSuffix: "k8s.io",
			attributes:      []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{BoolValue: ptr.To(true)}}},
			expectMatch:     true,
		},
		"celSuffix": {
			expression:  `device.bool["name"]`,
			celSuffix:   "k8s.io",
			attributes:  []api.DeviceAttribute{{Name: "name.k8s.io", DeviceAttributeValue: api.DeviceAttributeValue{BoolValue: ptr.To(true)}}},
			expectMatch: true,
		},
		"bothSuffix": {
			expression:      `device.bool["name"]`,
			celSuffix:       "k8s.io",
			attributeSuffix: "k8s.io",
			attributes:      []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{BoolValue: ptr.To(true)}}},
			expectMatch:     true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, ctx := ktesting.NewTestContext(t)
			result := Compiler.CompileCELExpression(scenario.expression, environment.StoredExpressions)
			if scenario.expectCompileError != "" && result.Error == nil {
				t.Fatalf("expected compile error %q, got none", scenario.expectCompileError)
			}
			if result.Error != nil {
				if scenario.expectCompileError == "" {
					t.Fatalf("unexpected compile error: %v", result.Error)
				}
				if !strings.Contains(result.Error.Error(), scenario.expectCompileError) {
					t.Fatalf("expected compile error to contain %q, but got instead: %v", scenario.expectCompileError, result.Error)
				}
				return
			}
			input := Input{
				AttributeSuffix: scenario.attributeSuffix,
				Attributes:      scenario.attributes,
				CELSuffix:       scenario.celSuffix,
			}
			match, err := result.Evaluate(ctx, input)
			if err != nil {
				if scenario.expectMatchError == "" {
					t.Fatalf("unexpected evaluation error: %v", err)
				}
				if !strings.Contains(err.Error(), scenario.expectMatchError) {
					t.Fatalf("expected evaluation error to contain %q, but got instead: %v", scenario.expectMatchError, err)
				}
				return
			}
			if match != scenario.expectMatch {
				t.Fatalf("expected result %v, got %v", scenario.expectMatch, match)
			}
		})
	}
}
