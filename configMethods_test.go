package main

import "testing"

func Test_roundFloat(t *testing.T) {
	type args struct {
		val       float64
		precision uint
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"10 to 2", args{10.1234567891, 2}, 10.12},
		{"3 to 1", args{99.999, 1}, 100.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := roundFloat(tt.args.val, tt.args.precision); got != tt.want {
				t.Errorf("roundFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
