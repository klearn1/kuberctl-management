package winstats

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGROUP_AFFINITY_Processors(t *testing.T) {
	tests := []struct {
		name  string
		Mask  uintptr
		Group uint16
		want  []int
	}{
		{
			name:  "empty",
			Mask:  0,
			Group: 0,
			want:  []int{},
		},
		{
			name:  "empty group 2",
			Mask:  0,
			Group: 1,
			want:  []int{},
		},
		{
			name:  "cpu 1 Group 0",
			Mask:  1,
			Group: 0,
			want:  []int{0},
		},
		{
			name:  "cpu 64 Group 0",
			Mask:  1 << 63,
			Group: 0,
			want:  []int{63},
		},
		{
			name:  "cpu 128 Group 1",
			Mask:  1 << 63,
			Group: 1,
			want:  []int{127},
		},
		{
			name:  "cpu 128 (Group 1)",
			Mask:  1 << 63,
			Group: 1,
			want:  []int{127},
		},
		{
			name:  "Mask 1 Group 2",
			Mask:  1,
			Group: 2,
			want:  []int{128},
		},
		{
			name:  "64 cpus group 0",
			Mask:  0xffffffffffffffff,
			Group: 0,
			want:  makeRange(0, 63),
		},
		{
			name:  "64 cpus group 1",
			Mask:  0xffffffffffffffff,
			Group: 1,
			want:  makeRange(64, 127),
		},
		{
			name:  "64 cpus group 1",
			Mask:  0xffffffffffffffff,
			Group: 1,
			want:  makeRange(64, 127),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := GROUP_AFFINITY{
				Mask:  tt.Mask,
				Group: tt.Group,
			}
			assert.Equalf(t, tt.want, a.Processors(), "Processors()")
		})
	}
}

// https://stackoverflow.com/a/39868255/697126
func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func TestCpusToGroupAffinity(t *testing.T) {
	tests := []struct {
		name string
		cpus []int
		want map[int]*GROUP_AFFINITY
	}{
		{
			name: "empty",
			want: map[int]*GROUP_AFFINITY{},
		},
		{
			name: "single cpu group 0",
			cpus: []int{0},
			want: map[int]*GROUP_AFFINITY{
				0: {
					Mask:  1,
					Group: 0,
				},
			},
		},
		{
			name: "single cpu group 0",
			cpus: []int{63},
			want: map[int]*GROUP_AFFINITY{
				0: {
					Mask:  1 << 63,
					Group: 0,
				},
			},
		},
		{
			name: "single cpu group 1",
			cpus: []int{64},
			want: map[int]*GROUP_AFFINITY{
				1: {
					Mask:  1,
					Group: 1,
				},
			},
		},
		{
			name: "multiple cpus same group",
			cpus: []int{0, 1, 2},
			want: map[int]*GROUP_AFFINITY{
				0: {
					Mask:  1 | 2 | 4, // Binary OR to combine the masks
					Group: 0,
				},
			},
		},
		{
			name: "multiple cpus different groups",
			cpus: []int{0, 64},
			want: map[int]*GROUP_AFFINITY{
				0: {
					Mask:  1,
					Group: 0,
				},
				1: {
					Mask:  1,
					Group: 1,
				},
			},
		},
		{
			name: "multiple cpus different groups",
			cpus: []int{0, 1, 2, 64, 65, 66},
			want: map[int]*GROUP_AFFINITY{
				0: {
					Mask:  1 | 2 | 4,
					Group: 0,
				},
				1: {
					Mask:  1 | 2 | 4,
					Group: 1,
				},
			},
		},
		{
			name: "64 cpus group 0",
			cpus: makeRange(0, 63),
			want: map[int]*GROUP_AFFINITY{
				0: {
					Mask:  0xffffffffffffffff, // All 64 bits set
					Group: 0,
				},
			},
		},
		{
			name: "64 cpus group 1",
			cpus: makeRange(64, 127),
			want: map[int]*GROUP_AFFINITY{
				1: {
					Mask:  0xffffffffffffffff, // All 64 bits set
					Group: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CpusToGroupAffinity(tt.cpus), "CpusToGroupAffinity(%v)", tt.cpus)
		})
	}
}
