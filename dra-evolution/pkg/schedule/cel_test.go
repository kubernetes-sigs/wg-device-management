package schedule

import (
	"testing"

	"github.com/kubernetes-sigs/wg-device-management/dra-evolution/pkg/api"
	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/api/resource"
)

func ptr[T any](val T) *T {
	var v T = val
	return &v
}

func TestMeetsConstraints(t *testing.T) {
	testCases := map[string]struct {
		constraints *string
		attrs       []api.Attribute
		expErr      string
		result      bool
	}{
		"nil constraint": {
			constraints: nil,
			result:      true,
		},
		"empty constraint": {
			constraints: ptr(""),
			result:      true,
		},
		"simple constraint met": {
			constraints: ptr("device.vendor == 'example.com'"),
			attrs: []api.Attribute{
				{
					Name:        "vendor",
					StringValue: ptr("example.com"),
				},
			},
			result: true,
		},
		"simple constraint failed": {
			constraints: ptr("device.vendor == 'example.com'"),
			attrs: []api.Attribute{
				{
					Name:        "vendor",
					StringValue: ptr("example.org"),
				},
			},
			result: false,
		},
		"multi-attribute constraint met": {
			constraints: ptr("device.vendor == 'example.com' && device.model == 'foozer-1000'"),
			attrs: []api.Attribute{
				{
					Name:        "vendor",
					StringValue: ptr("example.com"),
				},
				{
					Name:        "model",
					StringValue: ptr("foozer-1000"),
				},
			},
			result: true,
		},
		"simple device and pool constraint failed": {
			constraints: ptr("device.vendor == 'example.com' && device.model == 'foozer-1000'"),
			attrs: []api.Attribute{
				{
					Name:        "vendor",
					StringValue: ptr("example.org"),
				},
				{
					Name:        "model",
					StringValue: ptr("foozer-1000"),
				},
			},
			result: false,
		},
		//TODO: add CEL type conversion for resource.Quantity so this test below can be enabled
		"quantity constraint met": {
			constraints: ptr("device.memory >= '10Gi'"),
			attrs: []api.Attribute{
				{
					Name:          "memory",
					QuantityValue: ptr(resource.MustParse("10Gi")),
				},
			},
			expErr: "unsupported conversion to ref.Val: (resource.Quantity){{10737418240 0} {<nil>} 10Gi BinarySI}",
			result: true,
		},
	}
	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			result, err := MeetsConstraints(tc.constraints, tc.attrs)
			if tc.expErr == "" {
				require.NoError(t, err)
				require.Equal(t, tc.result, result)
			} else {
				require.EqualError(t, err, tc.expErr)
			}
		})
	}
}
