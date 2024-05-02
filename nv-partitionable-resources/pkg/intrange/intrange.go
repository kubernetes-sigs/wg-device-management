package intrange

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
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

// Contains checks if a specified range is present within the IntRange.
func (r *IntRange) Contains(o *IntRange) bool {
	matched := 0
	for _, s1 := range o.segments {
		for _, s2 := range r.segments {
			if s2.contains(&s1) {
				matched++
				break
			}
		}
	}
	if matched == len(o.segments) {
		return true
	}
	return false
}

// available checks if a specified range is available within the segment.
func (s segment) contains(o *segment) bool {
	if o.end < o.start {
		return false
	}
	if o.start < s.start {
		return false
	}
	if o.end > s.end {
		return false
	}
	return true
}

// Difference returns the difference of the original intrange from the other.
func (r *IntRange) Difference(o *IntRange) *IntRange {
	other := o.DeepCopy()
	other.compress()
	result := &IntRange{}
	for _, segR := range r.segments {
		for _, segOther := range other.segments {
			newSegments := segR.difference(&segOther)
			result.segments = slices.Concat(result.segments, newSegments)
		}
	}
	result.compress()
	return result
}

// difference returns the difference of the original segment from the other (as a list os segments).
func (s *segment) difference(other *segment) []segment {
	if s.start > s.end {
		return []segment{*s}
	}
	if s.end < other.start {
		return []segment{*s}
	}
	if s.start >= other.start && s.end <= other.end {
		return []segment{}
	}
	if s.start >= other.start {
		return []segment{{start: other.end + 1, end: s.end}}
	}
	if s.end <= other.end {
		return []segment{{start: s.start, end: other.start - 1}}
	}
	return []segment{
		{start: s.start, end: other.start - 1},
		{start: other.end + 1, end: s.end},
	}
}

// Add adds an IntRange to the original IntRange and returns an error if its already present.
func (r *IntRange) Union(o *IntRange) *IntRange {
	other := o.DeepCopy()
	other.compress()
	result := &IntRange{}
	result.segments = slices.Concat(r.segments, other.segments)
	result.compress()
	return result
}

// compress compresses the available segments into contiguous ranges.
func (r *IntRange) compress() {
	if len(r.segments) < 2 {
		return
	}

	sort.Slice(r.segments, func(i, j int) bool {
		return r.segments[i].end < r.segments[j].end
	})
	sort.Slice(r.segments, func(i, j int) bool {
		return r.segments[i].start < r.segments[j].start
	})

	var newSegments []segment
	for _, s := range r.segments {
		mergeOrAppendSegment(&newSegments, s)

	}
	r.segments = newSegments
}

// mergeSegments either updates the last segment in the original segment list
// (if a merge is possible) or appends a new segment to the list.
func mergeOrAppendSegment(original *[]segment, additional segment) {
	if len(*original) == 0 {
		*original = append(*original, additional)
		return
	}
	lastItem := &((*original)[len(*original)-1])
	if (lastItem.end + 1) >= additional.start {
		lastItem.end = additional.end
		return
	}
	*original = append(*original, additional)
}
