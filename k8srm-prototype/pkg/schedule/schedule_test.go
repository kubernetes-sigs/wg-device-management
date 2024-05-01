package schedule

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/johnbelamaric/k8srm-prototype/pkg/api"
	"github.com/johnbelamaric/k8srm-prototype/pkg/gen"
	"github.com/stretchr/testify/require"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/yaml"
)

func dumpTestClaims(tn string, claims []api.DeviceClaim) {
	if os.Getenv("DUMP_TEST_CASES") != "y" {
		return
	}

	cleanup := func(r rune) rune {
		if r < 'a' || r > 'z' {
			return '-'
		}
		return r
	}
	file := "testdata/claims-" + strings.Map(cleanup, strings.ToLower(tn)) + ".yaml"

	b, _ := yaml.Marshal(claims)
	err := os.WriteFile(file, b, 0644)
	if err != nil {
		fmt.Printf("error saving file %q: %s\n", file, err)
	}
}

func TestSelectNode(t *testing.T) {
	mixedPools := append(gen.GenShapeZero(2), gen.GenShapeThree(2)...)
	testCases := map[string]struct {
		claims        []api.DeviceClaim
		pools         []api.DevicePool
		expectSuccess bool
	}{
		"single by driver": {
			claims: []api.DeviceClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myclaim",
						Namespace: "default",
					},
					Spec: api.DeviceClaimSpec{
						DeviceClass: "not implemented yet",
						Driver:      ptr("example.com-foozer"),
					},
				},
			},
			pools:         mixedPools,
			expectSuccess: true,
		},
		"single with constraint met": {
			claims: []api.DeviceClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myclaim",
						Namespace: "default",
					},
					Spec: api.DeviceClaimSpec{
						DeviceClass: "not implemented yet",
						Constraints: ptr("device.model == 'foozer-1000'"),
					},
				},
			},
			pools:         mixedPools,
			expectSuccess: true,
		},
		"single with constraint not met": {
			claims: []api.DeviceClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myclaim",
						Namespace: "default",
					},
					Spec: api.DeviceClaimSpec{
						DeviceClass: "not implemented yet",
						Constraints: ptr("device.model == 'foozer-8000'"),
					},
				},
			},
			pools:         mixedPools,
			expectSuccess: false,
		},
		"multiple single pool": {
			claims: []api.DeviceClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myclaim",
						Namespace: "default",
					},
					Spec: api.DeviceClaimSpec{
						DeviceClass:    "not implemented yet",
						Driver:         ptr("example.com-foozer"),
						MinDeviceCount: ptr(2),
					},
				},
			},
			pools:         gen.GenShapeZero(2),
			expectSuccess: true,
		},
		"split across pools": {
			claims: []api.DeviceClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myclaim",
						Namespace: "default",
					},
					Spec: api.DeviceClaimSpec{
						DeviceClass:    "not implemented yet",
						Driver:         ptr("example.com-foozer"),
						MinDeviceCount: ptr(4),
					},
				},
			},
			pools:         gen.GenShapeOne(2),
			expectSuccess: true,
		},
		"two claims not satisfiable by a single node": {
			claims: []api.DeviceClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foozer-claim",
						Namespace: "default",
					},
					Spec: api.DeviceClaimSpec{
						DeviceClass:    "not implemented yet",
						Driver:         ptr("example.com-foozer"),
						MinDeviceCount: ptr(2),
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "barzer-claim",
						Namespace: "default",
					},
					Spec: api.DeviceClaimSpec{
						DeviceClass:    "not implemented yet",
						Driver:         ptr("example.com-barzer"),
						MinDeviceCount: ptr(2),
					},
				},
			},
			pools:         mixedPools,
			expectSuccess: false,
		},
		"two claims satisfiable by a single node": {
			claims: []api.DeviceClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foozer-claim",
						Namespace: "default",
					},
					Spec: api.DeviceClaimSpec{
						DeviceClass:    "not implemented yet",
						Driver:         ptr("example.com-foozer"),
						MinDeviceCount: ptr(2),
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "barzer-claim",
						Namespace: "default",
					},
					Spec: api.DeviceClaimSpec{
						DeviceClass:    "not implemented yet",
						Driver:         ptr("example.com-barzer"),
						MinDeviceCount: ptr(2),
					},
				},
			},
			pools:         gen.GenFoozerBarzerNodes(2),
			expectSuccess: true,
		},
		"claim cannot be met due to NUMA MatchAttribute": {
			claims: []api.DeviceClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myclaim",
						Namespace: "default",
					},
					Spec: api.DeviceClaimSpec{
						DeviceClass:     "not implemented yet",
						Driver:          ptr("example.com-foozer"),
						MinDeviceCount:  ptr(4),
						MatchAttributes: []string{"numa"},
					},
				},
			},
			pools:         gen.GenShapeOne(2),
			expectSuccess: false,
		},
	}

	for tn, tc := range testCases {
		verbose := os.Getenv("VERBOSE") == "y"

		t.Run(tn, func(t *testing.T) {
			dumpTestClaims(tn, tc.claims)
			allocations, results := SelectNode(tc.claims, tc.pools)
			b, _ := yaml.Marshal(allocations)
			fmt.Println()
			fmt.Println("=== TEST " + tn)
			fmt.Println()
			fmt.Println("ALLOCATIONS")
			fmt.Println("-----------")
			fmt.Println(string(b))
			fmt.Println("NODE RESULTS")
			fmt.Println("------------")
			if verbose {
				b, _ = yaml.Marshal(results)
				fmt.Println(string(b))
			} else {
				for _, nr := range results {
					fmt.Println(nr.Summary())
				}
			}
			fmt.Println()
			fmt.Println("=== DONE " + tn)
			fmt.Println()

			require.Equal(t, tc.expectSuccess, allocations != nil)
		})
	}
}
