package main

import (
	"reflect"
	"testing"
)

func Test_stringSubtraction(t *testing.T) {
	type args struct {
		s      string
		chunks []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"memcached",
			args{
				s:      "memcached-1.4.33.tar.gz",
				chunks: []string{"memcached-", ".", ".", ".tar.gz"},
			},
			[]string{"1", "4", "33"},
			false,
		}, {
			"no-suffix",
			args{
				s:      "package-1.4.33",
				chunks: []string{"package-", ".", "."},
			},
			[]string{"1", "4", "33"},
			false,
		}, {
			"no-prefix",
			args{
				s:      "1.4.33.tar.gz",
				chunks: []string{".", ".", ".tar.gz"},
			},
			[]string{"1", "4", "33"},
			false,
		}, {
			"out-of-order",
			args{
				s:      "memcached-1.4.33.tar.gz",
				chunks: []string{".", ".tar.gz", "."},
			},
			[]string{},
			true,
		}, {
			"missing",
			args{
				s:      "memcached-1.4.33.tar.gz",
				chunks: []string{".", ".tar.gz", "."},
			},
			[]string{},
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringSubtraction(tt.args.s, tt.args.chunks)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringSubtraction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringSubtraction() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
