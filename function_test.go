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

func Test_stringsToInts(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			"single",
			args{[]string{"3"}},
			[]int{3},
			false,
		}, {
			"double",
			args{[]string{"1", "5"}},
			[]int{1, 5},
			false,
		}, {
			"single invalid",
			args{[]string{"a"}},
			[]int{},
			true,
		}, {
			"double both invalid",
			args{[]string{"a", "b"}},
			[]int{},
			true,
		}, {
			"double first invalid",
			args{[]string{"a", "2"}},
			[]int{},
			true,
		}, {
			"double second invalid",
			args{[]string{"2", "j"}},
			[]int{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringsToInts(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringsToInts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringsToInts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prefixArray(t *testing.T) {
	type args struct {
		arr    []string
		prefix string
	}
	tests := []struct {
		name    string
		args    args
		wantOut []string
	}{
		{
			"prefixDemo",
			args{
				arr: []string{
					"banana",
					"berry",
					"blackberry",
					"blood orange",
					"blueberry",
					"boysenberry",
					"breadfruit",
				},
				prefix: "test-",
			},
			[]string{
				"test-banana",
				"test-berry",
				"test-blackberry",
				"test-blood orange",
				"test-blueberry",
				"test-boysenberry",
				"test-breadfruit",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := prefixArray(tt.args.arr, tt.args.prefix)
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("prefixArray() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_trimEmpty(t *testing.T) {
	type args struct {
		parts []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimEmpty(tt.args.parts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trimEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
