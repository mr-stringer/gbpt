package main

import "testing"

func Test_getPssdFromSize(t *testing.T) {
	type args struct {
		sz uint
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"P1", args{3}, "P1"},
		{"P2", args{8}, "P2"},
		{"P3", args{9}, "P3"},
		{"P4", args{30}, "P4"},
		{"P6", args{64}, "P6"},
		{"P10", args{65}, "P10"},
		{"P15", args{256}, "P15"},
		{"P20", args{500}, "P20"},
		{"P30", args{1024}, "P30"},
		{"P40", args{1543}, "P40"},
		{"P50", args{4096}, "P50"},
		{"P60", args{8192}, "P60"},
		{"P70", args{16384}, "P70"},
		{"P80", args{32768}, "P80"},
		{"Error", args{32999}, "error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPssdFromSize(tt.args.sz); got != tt.want {
				t.Errorf("getPssdFromSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
