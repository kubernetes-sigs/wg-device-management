package cel

import (
	"testing"

	"github.com/kubernetes-sigs/wg-device-management/dra-evolution/pkg/api"
	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"
)

func TestMeetsConstraints(t *testing.T) {
	testCases := map[string]struct {
		constraints *string
		attrs       []api.DeviceAttribute
		expErr      string
		result      bool
	}{
		"nil constraint": {
			constraints: nil,
			result:      true,
		},
		"empty constraint": {
			constraints: ptr.To(""),
			result:      true,
		},
		"simple constraint met": {
			constraints: ptr.To("device.string['vendor'] == 'example.com'"),
			attrs: []api.DeviceAttribute{
				{
					Name:                 "vendor",
					DeviceAttributeValue: api.DeviceAttributeValue{StringValue: ptr.To("example.com")},
				},
			},
			result: true,
		},
		"simple constraint failed": {
			constraints: ptr.To("device.string['vendor'] == 'example.com'"),
			attrs: []api.DeviceAttribute{
				{
					Name:                 "vendor",
					DeviceAttributeValue: api.DeviceAttributeValue{StringValue: ptr.To("example.org")},
				},
			},
			result: false,
		},
		"multi-attribute constraint met": {
			constraints: ptr.To("device.string['vendor'] == 'example.com' && device.string['model'] == 'foozer-1000'"),
			attrs: []api.DeviceAttribute{
				{
					Name:                 "vendor",
					DeviceAttributeValue: api.DeviceAttributeValue{StringValue: ptr.To("example.com")},
				},
				{
					Name:                 "model",
					DeviceAttributeValue: api.DeviceAttributeValue{StringValue: ptr.To("foozer-1000")},
				},
			},
			result: true,
		},
		"simple device and pool constraint failed": {
			constraints: ptr.To("device.string['vendor'] == 'example.com' && device.string['model'] == 'foozer-1000'"),
			attrs: []api.DeviceAttribute{
				{
					Name:                 "vendor",
					DeviceAttributeValue: api.DeviceAttributeValue{StringValue: ptr.To("example.org")},
				},
				{
					Name:                 "model",
					DeviceAttributeValue: api.DeviceAttributeValue{StringValue: ptr.To("foozer-1000")},
				},
			},
			result: false,
		},
		// CEL type conversion for resource.Quantity is intentionally not
		// defined in Kubernetes because it would make runtime cost evaluation
		// very hard.
		"quantity constraint met": {
			constraints: ptr.To("device.quantity['memory'].compareTo(quantity('10Gi')) >= 0"),
			attrs: []api.DeviceAttribute{
				{
					Name:                 "memory",
					DeviceAttributeValue: api.DeviceAttributeValue{QuantityValue: ptr.To(resource.MustParse("10Gi"))},
				},
			},
			result: true,
		},
		"fully-qualified-name": {
			constraints: ptr.To("device.string['vendor.dra.example.com'] == 'example.com'"),
			attrs: []api.DeviceAttribute{
				{
					Name:                 "vendor.dra.example.com",
					DeviceAttributeValue: api.DeviceAttributeValue{StringValue: ptr.To("example.com")},
				},
			},
			result: true,
		},
	}
	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			result, err := MeetsConstraints(tc.constraints, DeviceAttributes{Attributes: tc.attrs})
			if tc.expErr == "" {
				require.NoError(t, err)
				require.Equal(t, tc.result, result)
			} else {
				require.EqualError(t, err, tc.expErr)
			}
		})
	}
}
