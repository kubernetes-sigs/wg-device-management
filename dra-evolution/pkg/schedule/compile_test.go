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
		driverName         string
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
			expression:         `device.quantityAttributes["no-such-attr"]`,
			expectCompileError: "must evaluate to bool",
		},
		"runtime-error": {
			expression:       `device.quantityAttributes["no-such-attr"].isGreaterThan(quantity("0"))`,
			expectMatchError: "no such key: no-such-attr",
		},
		"quantity": {
			expression:  `device.quantityAttributes["name"].isGreaterThan(quantity("0"))`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{QuantityValue: ptr.To(resource.MustParse("1"))}}},
			expectMatch: true,
		},
		"bool": {
			expression:  `device.boolAttributes["name"]`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{BoolValue: ptr.To(true)}}},
			expectMatch: true,
		},
		"int": {
			expression:  `device.intAttributes["name"] > 0`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{IntValue: ptr.To(int64(1))}}},
			expectMatch: true,
		},
		"intUntyped": {
			expression:  `device.attributes["name"] > 0`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{IntValue: ptr.To(int64(1))}}},
			expectMatch: true,
		},
		"intslice": {
			expression:  `device.intsliceAttributes["name"].isSorted() && device.intsliceAttributes["name"].indexOf(3) == 2`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{IntSliceValue: &api.IntSlice{Ints: []int64{1, 2, 3}}}}},
			expectMatch: true,
		},
		"empty-intslice": {
			expression:  `size(device.intsliceAttributes["name"]) == 0`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{IntSliceValue: &api.IntSlice{}}}},
			expectMatch: true,
		},
		"string": {
			expression:  `device.stringAttributes["name"] == "fish"`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{StringValue: ptr.To("fish")}}},
			expectMatch: true,
		},
		"stringslice": {
			expression:  `device.stringsliceAttributes["name"].isSorted() && device.stringsliceAttributes["name"].indexOf("a") == 0`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{StringSliceValue: &api.StringSlice{Strings: []string{"a", "b", "c"}}}}},
			expectMatch: true,
		},
		"empty-stringslice": {
			expression:  `size(device.stringsliceAttributes["name"]) == 0`,
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{StringSliceValue: &api.StringSlice{}}}},
			expectMatch: true,
		},
		"all": {
			expression: `device.quantityAttributes["quantity"].isGreaterThan(quantity("0")) &&
device.boolAttributes["bool"] &&
device.intAttributes["int"] > 0 &&
device.intsliceAttributes["intslice"].isSorted() &&
device.stringAttributes["string"] == "fish" &&
device.stringsliceAttributes["stringslice"].isSorted()`,
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
			expression: `device.boolAttributes["a"] && device.boolAttributes["b"] && device.boolAttributes["c"]`,
			attributes: []api.DeviceAttribute{
				{Name: "a", DeviceAttributeValue: api.DeviceAttributeValue{BoolValue: ptr.To(true)}},
				{Name: "b", DeviceAttributeValue: api.DeviceAttributeValue{BoolValue: ptr.To(true)}},
				{Name: "c", DeviceAttributeValue: api.DeviceAttributeValue{BoolValue: ptr.To(true)}},
			},
			expectMatch: true,
		},
		"driverNameSuffix": {
			expression:  `device.boolAttributes["name.dra.example.com"]`,
			driverName:  "dra.example.com",
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{BoolValue: ptr.To(true)}}},
			expectMatch: true,
		},
		// TODO: negative test cases and <key> in `device.bool`.
		// Code doesn't support that properly yet either.
		"celSuffix": {
			expression:  `device.boolAttributes["name"]`,
			celSuffix:   "dra.example.com",
			attributes:  []api.DeviceAttribute{{Name: "name.dra.example.com", DeviceAttributeValue: api.DeviceAttributeValue{BoolValue: ptr.To(true)}}},
			expectMatch: true,
		},
		"bothSuffix": {
			expression:  `device.boolAttributes["name"]`,
			celSuffix:   "dra.example.com",
			driverName:  "dra.example.com",
			attributes:  []api.DeviceAttribute{{Name: "name", DeviceAttributeValue: api.DeviceAttributeValue{BoolValue: ptr.To(true)}}},
			expectMatch: true,
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
			input := DeviceAttributes{
				DriverName: scenario.driverName,
				Attributes: scenario.attributes,
				CELSuffix:  scenario.celSuffix,
			}
			match, err := result.DeviceMatches(ctx, input)
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
