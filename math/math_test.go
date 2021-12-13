package cvmath

import (
	"testing"
)

func TestMaxMinAbs(t *testing.T) {
	type args struct {
		data []float64
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want1 float64
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				data: []float64{1, -4, 3, 2, 0.5, -0.01, 7},
			},
			want:  7,
			want1: 0.01,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := MaxMinAbs(tt.args.data...)
			if got != tt.want {
				t.Errorf("MaxMinAbs() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MaxMinAbs() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
