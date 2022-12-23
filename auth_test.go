package main

import "testing"

func Test_longCircuitAnd(t *testing.T) {
	type args struct {
		p bool
		q bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "true",
			args: args{
				p: true,
				q: true,
			},
			want: true,
		},
		{
			name: "false",
			args: args{
				p: false,
				q: false,
			},
			want: false,
		},
		{
			name: "true false",
			args: args{
				p: true,
				q: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := longCircuitAnd(tt.args.p, tt.args.q); got != tt.want {
				t.Errorf("longCircuitAnd() = %v, want %v", got, tt.want)
			}
		})
	}
}
