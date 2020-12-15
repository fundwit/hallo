package util

import (
	"testing"
)

func TestHashSha1Hex(t *testing.T) {
	type args struct {
		input []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "case1", args: args{input: []byte("Sha1测试")}, want: "56e3f4ad0b053132a8a76fd62c942f4ec57de368"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashSha1Hex(tt.args.input); got != tt.want {
				t.Errorf("HashSha1Hex() = %v, want %v", got, tt.want)
			}
		})
	}
}
