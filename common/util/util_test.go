package util

import (
	"testing"
)

func TestGenRandString(t *testing.T) {
	type args struct {
		times int
	}
	tests := []struct {
		name    string
		args    args
		wantlen int
	}{
		{
			name: "GenRandString(5) should return len 5 string",
			args: args{
				times: 5,
			},
			wantlen: 5,
		},
		{
			name: "GenRandString(10) should return len 10 string",
			args: args{
				times: 10,
			},
			wantlen: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCode := GenRandString(tt.args.times); len(gotCode) != tt.wantlen {
				t.Errorf("len(GenRandString(%v)) = %v, want %v", tt.args.times, len(gotCode), tt.wantlen)
			}
		})
	}
}
