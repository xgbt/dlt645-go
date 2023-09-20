package dlt645

import (
	"reflect"
	"testing"
)

func Test_dataBlock(t *testing.T) {
	type args struct {
		value []uint16
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{
			args: args{
				value: []uint16{0x1122, 0x3344},
			},
			want: []byte{0x22, 0x11},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dataBlock(tt.args.value...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
