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
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"

	"k8s.io/apimachinery/pkg/util/version"
	celconfig "k8s.io/apiserver/pkg/apis/cel"
	apiservercel "k8s.io/apiserver/pkg/cel"
	"k8s.io/apiserver/pkg/cel/environment"

	"github.com/kubernetes-sigs/wg-device-management/dra-evolution/pkg/api"
)

const (
	attributesVarPrefix = "device."
)

var (
	Compiler = newCompiler()
)

// CompilationResult represents a compiled expression.
type CompilationResult struct {
	Program     cel.Program
	Error       *apiservercel.Error
	Expression  string
	OutputType  *cel.Type
	Environment *cel.Env
}

// Input defines the input values for a CEL expression.
type Input struct {
	// attributeSuffix gets appended to any attribute which does not already have
	// a fully qualified name.
	AttributeSuffix string
	Attributes      []api.DeviceAttribute

	// CELSuffix gets appended to all attributes name in a CEL expression which are
	// not already fully qualified.
	CELSuffix string
}

type compiler struct {
	envset *environment.EnvSet
}

func newCompiler() *compiler {
	return &compiler{envset: mustBuildEnv()}
}

// CompileCELExpression returns a compiled CEL expression. It evaluates to bool.
func (c compiler) CompileCELExpression(expression string, envType environment.Type) CompilationResult {
	resultError := func(errorString string, errType apiservercel.ErrorType) CompilationResult {
		return CompilationResult{
			Error: &apiservercel.Error{
				Type:   errType,
				Detail: errorString,
			},
			Expression: expression,
		}
	}

	env, err := c.envset.Env(envType)
	if err != nil {
		return resultError(fmt.Sprintf("unexpected error loading CEL environment: %v", err), apiservercel.ErrorTypeInternal)
	}

	ast, issues := env.Compile(expression)
	if issues != nil {
		return resultError("compilation failed: "+issues.String(), apiservercel.ErrorTypeInvalid)
	}
	expectedReturnType := cel.BoolType
	if ast.OutputType() != expectedReturnType {
		return resultError(fmt.Sprintf("must evaluate to %v", expectedReturnType.String()), apiservercel.ErrorTypeInvalid)
	}
	_, err = cel.AstToCheckedExpr(ast)
	if err != nil {
		// should be impossible since env.Compile returned no issues
		return resultError("unexpected compilation error: "+err.Error(), apiservercel.ErrorTypeInternal)
	}
	prog, err := env.Program(ast,
		cel.InterruptCheckFrequency(celconfig.CheckFrequency),
	)
	if err != nil {
		return resultError("program instantiation failed: "+err.Error(), apiservercel.ErrorTypeInternal)
	}
	return CompilationResult{
		Program:     prog,
		Expression:  expression,
		OutputType:  ast.OutputType(),
		Environment: env,
	}
}

var valueTypes = map[string]struct {
	celType *cel.Type
	// get returns nil if the attribute doesn't have the type, otherwise
	// the value of that type.
	get func(attr api.DeviceAttribute) (any, error)
}{
	"quantity": {apiservercel.QuantityType, func(attr api.DeviceAttribute) (any, error) {
		if attr.QuantityValue == nil {
			return nil, nil
		}
		return apiservercel.Quantity{Quantity: attr.QuantityValue}, nil
	}},
	"bool": {cel.BoolType, func(attr api.DeviceAttribute) (any, error) {
		if attr.BoolValue == nil {
			return nil, nil
		}
		return *attr.BoolValue, nil
	}},
	"int": {cel.IntType, func(attr api.DeviceAttribute) (any, error) {
		if attr.IntValue == nil {
			return nil, nil
		}
		return *attr.IntValue, nil
	}},
	"intslice": {types.NewListType(cel.IntType), func(attr api.DeviceAttribute) (any, error) {
		if attr.IntSliceValue == nil {
			return nil, nil
		}
		return attr.IntSliceValue.Ints, nil
	}},
	"string": {cel.StringType, func(attr api.DeviceAttribute) (any, error) {
		if attr.StringValue == nil {
			return nil, nil
		}
		return *attr.StringValue, nil
	}},
	"stringslice": {types.NewListType(cel.StringType), func(attr api.DeviceAttribute) (any, error) {
		if attr.StringSliceValue == nil {
			return nil, nil
		}
		return attr.StringSliceValue.Strings, nil
	}},
	"version": {SemverType, func(attr api.DeviceAttribute) (any, error) {
		if attr.VersionValue == nil {
			return nil, nil
		}
		v, err := semver.Parse(*attr.VersionValue)
		if err != nil {
			return nil, fmt.Errorf("parse semantic version: %v", err)
		}

		return Semver{Version: v}, nil
	}},
}

var boolType = reflect.TypeOf(true)

func (c CompilationResult) Evaluate(ctx context.Context, input Input) (bool, error) {
	variables := make(map[string]any, len(valueTypes))
	for name, valueType := range valueTypes {
		m, err := buildValueMapper(c.Environment.CELTypeAdapter(), input, valueType.get)
		if err != nil {
			return false, fmt.Errorf("extract attributes with type %s: %v", name, err)
		}
		variables[attributesVarPrefix+name] = m
	}
	result, _, err := c.Program.ContextEval(ctx, variables)
	if err != nil {
		return false, err
	}
	resultAny, err := result.ConvertToNative(boolType)
	if err != nil {
		return false, fmt.Errorf("CEL result of type %s could not be converted to bool: %w", result.Type().TypeName(), err)
	}
	resultBool, ok := resultAny.(bool)
	if !ok {
		return false, fmt.Errorf("CEL native result value should have been a bool, got instead: %T", resultAny)
	}
	return resultBool, nil
}

func mustBuildEnv() *environment.EnvSet {
	envset := environment.MustBaseEnvSet(environment.DefaultCompatibilityVersion())
	versioned := []environment.VersionedOptions{
		{
			// Feature epoch was actually 1.30, but we artificially set it to 1.0 because these
			// options should always be present.
			//
			// TODO (https://github.com/kubernetes/kubernetes/issues/123687): set this
			// version properly before going to beta.
			IntroducedVersion: version.MajorMinor(1, 0),
			EnvOptions: append(buildVersionedAttributes(),
				SemverLib(),
			),
		},
	}
	envset, err := envset.Extend(versioned...)
	if err != nil {
		panic(fmt.Errorf("internal error building CEL environment: %w", err))
	}
	return envset
}

func buildVersionedAttributes() []cel.EnvOption {
	options := make([]cel.EnvOption, 0, len(valueTypes))
	for name, valueType := range valueTypes {
		options = append(options, cel.Variable(attributesVarPrefix+name, types.NewMapType(cel.StringType, valueType.celType)))
	}
	return options
}

func buildValueMapper(adapter types.Adapter, input Input, get func(api.DeviceAttribute) (any, error)) (traits.Mapper, error) {
	// This implementation constructs a map and then let's cel handle the
	// lookup and iteration. This is done for the sake of simplicity.
	// Whether it's faster than writing a custom mapper depends on
	// real-world attribute sets and CEL expressions and would have to be
	// benchmarked.
	valueMap := make(map[string]any)
	for _, attribute := range input.Attributes {
		value, err := get(attribute)
		if err != nil {
			return nil, fmt.Errorf("attribute %q: %v", attribute.Name, err)
		}
		if value != nil {
			name := qualifyAttributeName(attribute.Name, input.AttributeSuffix)
			valueMap[name] = value
		}
	}
	return &stringInterfaceMap{
		Mapper: types.NewStringInterfaceMap(adapter, valueMap),
		suffix: input.CELSuffix,
	}, nil
}

func qualifyAttributeName(name, suffix string) string {
	if suffix == "" || strings.Contains(name, ".") {
		return name
	}
	return name + "." + suffix
}

type stringInterfaceMap struct {
	traits.Mapper
	suffix string
}

var stringType = reflect.TypeOf("")

// Find adds the suffix to string keys. Unknown key types are passed
// down to the mapper which then will produce an error.
func (m *stringInterfaceMap) Find(key ref.Val) (ref.Val, bool) {
	if m.suffix == "" {
		return m.Mapper.Find(key)
	}
	anyKey, err := key.ConvertToNative(stringType)
	if err != nil {
		return m.Mapper.Find(key)
	}
	name, ok := anyKey.(string)
	if !ok {
		return m.Mapper.Find(key)
	}
	name = qualifyAttributeName(name, m.suffix)
	return m.Mapper.Find(types.String(name))
}
