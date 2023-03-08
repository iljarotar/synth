package file

import (
	"reflect"
	"testing"
)

func Test_intToBytes(t *testing.T) {
	type args struct {
		n   int
		num int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "0",
			args: args{
				n:   0,
				num: 4,
			},
			want: []byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			name: "1",
			args: args{
				n:   1,
				num: 4,
			},
			want: []byte{0x01, 0x00, 0x00, 0x00},
		},
		{
			name: "256",
			args: args{
				n:   256,
				num: 4,
			},
			want: []byte{0x00, 0x01, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intToBytes(tt.args.n, tt.args.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("intToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
