package cvmath

import (
	"fmt"
	"image"
	"reflect"
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

func TestStepSelect(t *testing.T) {
	type args struct {
		step   int
		offset int
		data   []float64
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				step:   3,
				offset: 0,
				data:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			want: []float64{1, 4, 7, 10},
		},
		{
			name: "",
			args: args{
				step:   3,
				offset: 1,
				data:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			want: []float64{2, 5, 8},
		},
		{
			name: "",
			args: args{
				step:   2,
				offset: 1,
				data:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			want: []float64{2, 4, 6, 8, 10},
		},
		{
			name: "",
			args: args{
				step:   4,
				offset: 3,
				data:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			want: []float64{4, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StepSelect(tt.args.step, tt.args.offset, tt.args.data...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StepSelect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModeVec4(t *testing.T) {
	type args struct {
		data []float64
	}
	tests := []struct {
		name string
		args args
		want [4]float64
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				data: []float64{
					1, 1, 3, 5,
					1, 1, 3, 4,
					1, 3, 3, 4,
					1, 2, 3, 4,
					1, 2, 3, 4,
					1, 2.0, 3, 4,
					1, 3.0, 3, 4,
				},
			},
			want: [4]float64{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ModeVec4(tt.args.data...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModeVec4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModeStep4(t *testing.T) {
	type args struct {
		data []float64
	}
	tests := []struct {
		name string
		args args
		want [4]float64
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				data: []float64{
					1, 1, 3, 5,
					1, 1, 3, 4,
					1, 3, 3, 4,
					1, 2, 3, 4,
					1, 2, 3, 4,
					1, 2.0, 3, 4,
					1, 3.0, 3, 4,
				},
			},
			want: [4]float64{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ModeStep4(tt.args.data...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModeStep4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeanStep4(t *testing.T) {
	type args struct {
		data []float64
	}
	tests := []struct {
		name     string
		args     args
		wantMean [4]float64
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				data: []float64{
					1, 2, 3, 4,
					1, 3, 3, 4,
					1, 3, 3, 4,
					1, 3, 3, 4,
					1, 2, 3, 4,
					1, 2.0, 3, 4,
					1, 3.0, 3, 4,
				},
			},
			wantMean: [4]float64{1, 3, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMean := MeanStep4(tt.args.data...); !reflect.DeepEqual(gotMean, tt.wantMean) {
				t.Errorf("MeanVec4() = %v, want %v", gotMean, tt.wantMean)
			}
		})
	}
}

func TestMedianStep4(t *testing.T) {
	type args struct {
		data []float64
	}
	tests := []struct {
		name string
		args args
		want [4]float64
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				data: []float64{
					1, 2, 3, 4,
					1, 3, 3, 3,
					1, 1, 3, 3,
					1, 3, 3, 3,
					1, 4, 3, 3,
					1, 5.0, 3, 4,
					1, 3.0, 3, 4,
				},
			},
			want: [4]float64{1, 3, 3, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MedianStep4(tt.args.data...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MedianVec4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYaXb(t *testing.T) {
	type args struct {
		pt1 image.Point
		pt2 image.Point
	}
	tests := []struct {
		name           string
		args           args
		wantA          float64
		wantB          float64
		wantVertical   bool
		wantVX         float64
		wantHorizontal bool
		wantHY         float64
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				pt1: image.Pt(1, 1),
				pt2: image.Pt(0, 0),
			},
			wantA:          1,
			wantB:          0,
			wantVertical:   false,
			wantHorizontal: false,
		},
		{
			name: "",
			args: args{
				pt1: image.Pt(1, 1),
				pt2: image.Pt(1, 0),
			},
			wantA:        1,
			wantB:        0,
			wantVertical: true,
			wantVX:       1,
		},
		{
			name: "",
			args: args{
				pt1: image.Pt(1, 1),
				pt2: image.Pt(0, 1),
			},
			wantA:          0,
			wantB:          1,
			wantHorizontal: true,
			wantHY:         1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotA, gotB, gotVertical, gotVX, gotHorizontal, gotHY := YaXb(tt.args.pt1, tt.args.pt2)
			if gotA != tt.wantA {
				t.Errorf("YaXb() gotA = %v, want %v", gotA, tt.wantA)
			}
			if gotB != tt.wantB {
				t.Errorf("YaXb() gotB = %v, want %v", gotB, tt.wantB)
			}
			if gotVertical != tt.wantVertical {
				t.Errorf("YaXb() gotVertical = %v, want %v", gotVertical, tt.wantVertical)
			}
			if gotVX != tt.wantVX {
				t.Errorf("YaXb() gotVX = %v, want %v", gotVX, tt.wantVX)
			}
			if gotHorizontal != tt.wantHorizontal {
				t.Errorf("YaXb() gotHorizontal = %v, want %v", gotHorizontal, tt.wantHorizontal)
			}
			if gotHY != tt.wantHY {
				t.Errorf("YaXb() gotHY = %v, want %v", gotHY, tt.wantHY)
			}
		})
	}
}

func TestInterpolate(t *testing.T) {
	type args struct {
		pt1  image.Point
		pt2  image.Point
		step float64
	}
	tests := []struct {
		name    string
		args    args
		wantPts []image.Point
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				pt1:  image.Pt(10, 0),
				pt2:  image.Pt(100, 35),
				step: 0.5,
			},
		},
		// {
		// 	name: "",
		// 	args: args{
		// 		pt1:  image.Pt(100, 0),
		// 		pt2:  image.Pt(100, 35),
		// 		step: 5,
		// 	},
		// },
		// {
		// 	name: "",
		// 	args: args{
		// 		pt1:  image.Pt(0, 100),
		// 		pt2:  image.Pt(35, 100),
		// 		step: 0.5,
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if gotPts := InterpolateLine(tt.args.pt1, tt.args.pt2, tt.args.step); !reflect.DeepEqual(gotPts, tt.wantPts) {
			// 	t.Errorf("InterpolateLine() = %v, want %v", gotPts, tt.wantPts)
			// }

			gotPts := interpolate(tt.args.pt1, tt.args.pt2, tt.args.step)
			for _, pt := range gotPts {
				fmt.Println(pt)
			}
		})
	}
}

func TestRandomPoint(t *testing.T) {
	type args struct {
		pts []image.Point
	}
	tests := []struct {
		name string
		args args
		want image.Point
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				pts: []image.Point{
					{X: 2, Y: 1},
					{X: 3, Y: 2},
					{X: 4, Y: 3},
					{X: 5, Y: 4},
				},
			},
			want: image.Pt(200, 100),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				got := RandomPoint(tt.args.pts...)
				switch got {
				case tt.args.pts[0]:
					fmt.Println("0")
				case tt.args.pts[1]:
					fmt.Println("1")
				case tt.args.pts[2]:
					fmt.Println("2")
				case tt.args.pts[3]:
					fmt.Println("3")
				}

				// fmt.Println(got)
			}
		})
	}
}
