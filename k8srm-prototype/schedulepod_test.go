package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
	"testing"
)

func dumpTestCase(tn string, claim PodCapacityClaim) {
	if os.Getenv("DUMP_TEST_CASES") != "y" {
		return
	}

	cleanup := func(r rune) rune {
		if r < 'a' || r > 'z' {
			return '-'
		}
		return r
	}
	file := "testdata/pod-" + strings.Map(cleanup, strings.ToLower(tn)) + ".yaml"

	b, _ := yaml.Marshal(claim)
	err := os.WriteFile(file, b, 0644)
	if err != nil {
		fmt.Printf("error saving file %q: %s\n", file, err)
	}
}

func TestSchedulePodForCore(t *testing.T) {
	testCases := map[string]struct {
		claim         PodCapacityClaim
		expectSuccess bool
	}{
		"single pod, single container": {
			claim: PodCapacityClaim{
				PodClaim: CapacityClaim{
					Name:   "my-pod",
					Claims: []ResourceClaim{genClaimPod()},
				},
				ContainerClaims: []CapacityClaim{
					{
						Name:   "my-container",
						Claims: []ResourceClaim{genClaimContainer("", "")},
					},
				},
			},
			expectSuccess: true,
		},
		"single pod, single container, with CPU and memory, enough": {
			claim: PodCapacityClaim{
				PodClaim: CapacityClaim{
					Name:   "my-pod",
					Claims: []ResourceClaim{genClaimPod()},
				},
				ContainerClaims: []CapacityClaim{
					{
						Name:   "my-container",
						Claims: []ResourceClaim{genClaimContainer("7127m", "8Gi")},
					},
				},
			},
			expectSuccess: true,
		},
		"single pod, single container, with CPU and memory, insufficient CPU": {
			claim: PodCapacityClaim{
				PodClaim: CapacityClaim{
					Name:   "my-pod",
					Claims: []ResourceClaim{genClaimPod()},
				},
				ContainerClaims: []CapacityClaim{
					{
						Name:   "my-container",
						Claims: []ResourceClaim{genClaimContainer("64", "8Gi")},
					},
				},
			},
			expectSuccess: false,
		},
		"single pod, single container, with CPU and memory, insufficient memory": {
			claim: PodCapacityClaim{
				PodClaim: CapacityClaim{
					Name:   "my-pod",
					Claims: []ResourceClaim{genClaimPod()},
				},
				ContainerClaims: []CapacityClaim{
					{
						Name:   "my-container",
						Claims: []ResourceClaim{genClaimContainer("4", "256Gi")},
					},
				},
			},
			expectSuccess: false,
		},
		"single pod, multiple containers, with CPU and memory, enough": {
			claim: PodCapacityClaim{
				PodClaim: CapacityClaim{
					Name:   "my-pod",
					Claims: []ResourceClaim{genClaimPod()},
				},
				ContainerClaims: []CapacityClaim{
					{
						Name:   "my-container-1",
						Claims: []ResourceClaim{genClaimContainer("7127m", "8Gi")},
					},
					{
						Name:   "my-container-2",
						Claims: []ResourceClaim{genClaimContainer("200m", "8Gi")},
					},
					{
						Name:   "my-container-3",
						Claims: []ResourceClaim{genClaimContainer("200m", "8Gi")},
					},
				},
			},
			expectSuccess: true,
		},
		"single pod, multiple containers, with CPU and memory, not enough": {
			claim: PodCapacityClaim{
				PodClaim: CapacityClaim{
					Name:   "my-pod",
					Claims: []ResourceClaim{genClaimPod()},
				},
				ContainerClaims: []CapacityClaim{
					{
						Name:   "my-container-1",
						Claims: []ResourceClaim{genClaimContainer("7127m", "8Gi")},
					},
					{
						Name:   "my-container-2",
						Claims: []ResourceClaim{genClaimContainer("8", "8Gi")},
					},
					{
						Name:   "my-container-3",
						Claims: []ResourceClaim{genClaimContainer("4", "8Gi")},
					},
				},
			},
			expectSuccess: false,
		},
		"no resources for driver": {
			claim: PodCapacityClaim{
				PodClaim: CapacityClaim{
					Name: "my-foozer-pod",
					Claims: []ResourceClaim{
						genClaimPod(),
						genClaimFoozer("foozer", "1m", "2Gi", 1),
					},
				},
				ContainerClaims: []CapacityClaim{
					{
						Name:   "my-container",
						Claims: []ResourceClaim{genClaimContainer("7127m", "8Gi")},
					},
				},
			},
			expectSuccess: false,
		},
	}

	for tn, tc := range testCases {
		capacity := genCapShapeZero(2)
		t.Run(tn, func(t *testing.T) {
			dumpTestCase(tn, tc.claim)
			allocation := SchedulePod(capacity, tc.claim)
			require.Equal(t, tc.expectSuccess, allocation != nil)
		})
	}
}

func TestSchedulePodForFoozer(t *testing.T) {
	testCases := map[string]struct {
		claim         PodCapacityClaim
		expectSuccess bool
	}{
		"single pod, container, cpu/mem, and foozer": {
			claim: PodCapacityClaim{
				PodClaim: CapacityClaim{
					Name: "my-foozer-pod",
					Claims: []ResourceClaim{
						genClaimPod(),
						genClaimFoozer("foozer", "1", "2Gi", 0),
					},
				},
				ContainerClaims: []CapacityClaim{
					{
						Name:   "my-container",
						Claims: []ResourceClaim{genClaimContainer("1", "4Gi")},
					},
				},
			},
			expectSuccess: true,
		},
		"no foozer big enough": {
			claim: PodCapacityClaim{
				PodClaim: CapacityClaim{
					Name: "my-foozer-pod",
					Claims: []ResourceClaim{
						genClaimPod(),
						genClaimFoozer("foozer", "16", "32Gi", 0),
					},
				},
				ContainerClaims: []CapacityClaim{
					{
						Name:   "my-container",
						Claims: []ResourceClaim{genClaimContainer("1", "4Gi")},
					},
				},
			},
			expectSuccess: false,
		},
	}

	for tn, tc := range testCases {
		capacity := genCapShapeOne(2)
		t.Run(tn, func(t *testing.T) {
			dumpTestCase(tn, tc.claim)
			allocation := SchedulePod(capacity, tc.claim)
			require.Equal(t, tc.expectSuccess, allocation != nil)
		})
	}
}

func TestSchedulePodForBigFoozer(t *testing.T) {
	testCases := map[string]struct {
		claim         PodCapacityClaim
		expectSuccess bool
	}{
		"single pod, container, cpu/mem, and foozer": {
			claim: PodCapacityClaim{
				PodClaim: CapacityClaim{
					Name: "my-foozer-pod",
					Claims: []ResourceClaim{
						genClaimPod(),
						genClaimFoozer("foozer", "1", "2Gi", 0),
					},
				},
				ContainerClaims: []CapacityClaim{
					{
						Name:   "my-container",
						Claims: []ResourceClaim{genClaimContainer("1", "4Gi")},
					},
				},
			},
			expectSuccess: true,
		},
		"no foozer big enough": {
			claim: PodCapacityClaim{
				PodClaim: CapacityClaim{
					Name: "my-foozer-pod",
					Claims: []ResourceClaim{
						genClaimPod(),
						genClaimFoozer("foozer", "16", "32Gi", 0),
					},
				},
				ContainerClaims: []CapacityClaim{
					{
						Name:   "my-container",
						Claims: []ResourceClaim{genClaimContainer("1", "4Gi")},
					},
				},
			},
			expectSuccess: true,
		},
	}

	for tn, tc := range testCases {
		capacity := genCapShapeTwo(2, 4)
		t.Run(tn, func(t *testing.T) {
			dumpTestCase(tn, tc.claim)
			allocation := SchedulePod(capacity, tc.claim)
			require.Equal(t, tc.expectSuccess, allocation != nil)
		})
	}
}
