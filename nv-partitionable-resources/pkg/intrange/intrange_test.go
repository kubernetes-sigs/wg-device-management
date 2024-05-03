/**
# Copyright 2024 NVIDIA CORPORATION
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
**/

package intrange

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntRangeDifference(t *testing.T) {
	testCases := []struct {
		description string
		original    *IntRange
		subtract    *IntRange
		expected    *IntRange
	}{
		{
			description: "",
			original:    MustParse("0"),
			subtract:    MustParse("0"),
			expected:    &IntRange{},
		},
		{
			description: "",
			original:    MustParse("0-100"),
			subtract:    MustParse("0-100"),
			expected:    &IntRange{},
		},
		{
			description: "",
			original:    MustParse("25-75"),
			subtract:    MustParse("0-100"),
			expected:    &IntRange{},
		},
		{
			description: "",
			original:    MustParse("0-100"),
			subtract:    MustParse("25-75"),
			expected:    MustParse("0-24,76-100"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			diff := tc.original.Difference(tc.subtract)
			require.Equal(t, tc.expected, diff)
		})
	}
}

func TestIntRangeUnion(t *testing.T) {
	testCases := []struct {
		description string
		original    *IntRange
		add         *IntRange
		expected    *IntRange
	}{
		{
			description: "",
			original:    MustParse("0"),
			add:         MustParse("0"),
			expected:    MustParse("0"),
		},
		{
			description: "",
			original:    MustParse("0-24,76-100"),
			add:         MustParse("25-75"),
			expected:    MustParse("0-100"),
		},
		{
			description: "",
			original:    MustParse("0-24,76-100"),
			add:         MustParse("20-80"),
			expected:    MustParse("0-100"),
		},
		{
			description: "",
			original:    MustParse("76-100"),
			add:         MustParse("50-80"),
			expected:    MustParse("50-100"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			sum := tc.original.Union(tc.add)
			require.Equal(t, tc.expected, sum)
		})
	}
}
