package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// IntRange represents a range of integer values.
// +k8s:deepcopy-gen=true
type IntRange struct {
	segments []segment
}

// segment represents a contiguous range of integers.
type segment struct {
	start int64
	end   int64
}

// result creates a new IntRange object with the specified bounds.
func NewIntRange(start, size int64) *IntRange {
	segments := []segment{
		{start: start, end: start + size - 1},
	}
	return &IntRange{segments: segments}
}

// String returns a string representation of the IntRange in the format "start-end,start-end".
func (r *IntRange) String() string {
	var parts []string
	for _, s := range r.segments {
		if s.start == s.end {
			parts = append(parts, fmt.Sprintf("%d", s.start))
		} else {
			parts = append(parts, fmt.Sprintf("%d-%d", s.start, s.end))
		}
	}
	return strings.Join(parts, ",")
}

// Parse parses a string representation of a range and returns an IntRange object.
func Parse(s string) (*IntRange, error) {
	parts := strings.Split(s, ",")
	var result IntRange
	for _, part := range parts {
		bounds := strings.Split(part, "-")
		if len(bounds) == 1 {
			val, err := strconv.ParseInt(bounds[0], 10, 64)
			if err != nil {
				return nil, errors.New("invalid range format")
			}
			result.segments = append(result.segments, segment{start: val, end: val})
			continue
		}
		if len(bounds) == 2 {
			start, err := strconv.ParseInt(bounds[0], 10, 64)
			if err != nil {
				return nil, errors.New("invalid start index")
			}
			end, err := strconv.ParseInt(bounds[1], 10, 64)
			if err != nil {
				return nil, errors.New("invalid end index")
			}
			if end < start {
				return nil, errors.New("end index must be greater than or equal to start index")
			}
			result.segments = append(result.segments, segment{start: start, end: end})
			continue
		}
		return nil, errors.New("invalid range format")
	}
	sort.Slice(result.segments, func(i, j int) bool {
		return result.segments[i].start < result.segments[j].start
	})
	return &result, nil
}

// MustParse parses a string representation of a range and returns an IntRange object.
func MustParse(s string) *IntRange {
	result, err := Parse(s)
	if err != nil {
		return nil
	}
	return result
}

// MarshalJSON is a custom marshaler for the IntRange type.
func (r *IntRange) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// UnmarshalJSON is a custom unmarshaler for the IntRange type.
func (r *IntRange) UnmarshalJSON(data []byte) error {
	var dataString string
	if err := json.Unmarshal(data, &dataString); err != nil {
		return err
	}
	newR, err := Parse(dataString)
	if err != nil {
		return fmt.Errorf("unable to parse %v as an IntRange: %w", dataString, err)
	}
	*r = *newR
	return nil
}
