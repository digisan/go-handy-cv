package blob

import (
	"reflect"
	"testing"
)

func TestParseIntPair(t *testing.T) {
	type args struct {
		pair string
	}
	tests := []struct {
		name     string
		args     args
		wantData [2]int
	}{
		// TODO: Add test cases.
		{
			name: "ParseIntPair",
			args: args{
				pair: "[1,2]",
			},
			wantData: [2]int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotData := ParseIntPair(tt.args.pair); !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ParseIntPair() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}
