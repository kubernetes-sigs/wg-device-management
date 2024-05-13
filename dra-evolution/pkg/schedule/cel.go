package schedule

import (
	"fmt"
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/kubernetes-sigs/wg-device-management/dra-evolution/pkg/api"
)

const (
	DeviceVarName = "device"
)

func MeetsConstraints(constraints *string, attrs []api.Attribute) (bool, error) {
	if constraints == nil || *constraints == "" {
		return true, nil
	}

	inputs := make(map[string]interface{})
	inputs[DeviceVarName] = attributesToInputs(attrs)

	return evalExpr(*constraints, inputs)
}

func attributesToInputs(attributes []api.Attribute) map[string]interface{} {
	result := make(map[string]interface{}, len(attributes))

	for _, a := range attributes {
		if a.StringValue != nil {
			result[a.Name] = *a.StringValue
		} else if a.IntValue != nil {
			result[a.Name] = *a.IntValue
		} else if a.QuantityValue != nil {
			result[a.Name] = *a.QuantityValue
		} else if a.SemVerValue != nil {
			result[a.Name] = *a.SemVerValue
		}
	}

	return result
}

func evalExpr(expr string, inputs map[string]interface{}) (bool, error) {
	prog, err := compileExpr(expr)
	if err != nil {
		return false, err
	}

	val, _, err := prog.Eval(inputs)
	if err != nil {
		return false, err
	}

	result, err := val.ConvertToNative(reflect.TypeOf(true))
	if err != nil {
		return false, err
	}

	s, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("expression returned non-string value: %v", result)
	}

	return s, nil
}

// compileExpr returns a compiled CEL expression.
func compileExpr(expr string) (cel.Program, error) {
	var opts []cel.EnvOption
	opts = append(opts, cel.HomogeneousAggregateLiterals())
	opts = append(opts, cel.EagerlyValidateDeclarations(true), cel.DefaultUTCTimeZone(true))
	opts = append(opts, cel.Variable(DeviceVarName, cel.DynType))

	env, err := cel.NewEnv(opts...)
	if err != nil {
		return nil, err
	}

	ast, issues := env.Compile(expr)
	if issues != nil {
		return nil, issues.Err()
	}

	_, err = cel.AstToCheckedExpr(ast)
	if err != nil {
		return nil, err
	}
	return env.Program(ast,
		cel.EvalOptions(cel.OptOptimize),
	)
}
