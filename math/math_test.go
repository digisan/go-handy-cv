package cvmath

import (
	"image"
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

func TestDisPt(t *testing.T) {
	type args struct {
		pt1 image.Point
		pt2 image.Point
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				pt1: image.Pt(-4, 0),
				pt2: image.Pt(0, -3),
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DisPt(tt.args.pt1, tt.args.pt2); got != tt.want {
				t.Errorf("DisPt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMedian(t *testing.T) {
	type args struct {
		data []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				data: []float64{4, 6, 2, 3, 5},
			},
			want: 4,
		},
		{
			name: "",
			args: args{
				data: []float64{4, 6, 2, 3, 5, 7},
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Median(tt.args.data...); got != tt.want {
				t.Errorf("Median() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMode(t *testing.T) {
	type args struct {
		data []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				data: []float64{4, 1, 1, 6, 2, 3, 3, 4, 5, 3, 3, 7, 3},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mode(tt.args.data...); got != tt.want {
				t.Errorf("Mode() = %v, want %v", got, tt.want)
			}
		})
	}
}
