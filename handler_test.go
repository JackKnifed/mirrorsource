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

func Test_versionType_possibleUpgrades(t *testing.T) {
	type fields struct {
		fmt  string
		cur  string
		past []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name: "memcached-vanilla",
			fields: fields{
				fmt:  "memcached-%v.%v.%v.tar.gz",
				cur:  "memcached-1.4.33.tar.gz",
				past: []string{},
			},
			want: []string{
				"memcached-2.0.0.tar.gz",
				"memcached-1.5.0.tar.gz",
				"memcached-1.4.34.tar.gz",
			},
			wantErr: false,
		}, {
			name: "memcached-vanilla",
			fields: fields{
				fmt:  "memcached-%v.%v.%v",
				cur:  "memcached-1.4.33",
				past: []string{},
			},
			want: []string{
				"memcached-2.0.0",
				"memcached-1.5.0",
				"memcached-1.4.34",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := versionType{
				fmt:  tt.fields.fmt,
				cur:  tt.fields.cur,
				past: tt.fields.past,
			}
			got, err := v.possibleUpgrades()
			if (err != nil) != tt.wantErr {
				t.Errorf("versionType.possibleUpgrades() error\nhave%v\nwantErr %v",
					err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("versionType.possibleUpgrades()\nhave %v\nwant %v", got, tt.want)
			}
		})
	}
}
