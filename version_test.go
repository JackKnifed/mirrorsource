package main

import (
	"reflect"
	"testing"
)

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
		}, {
			name: "simple-nosuffix",
			fields: fields{
				fmt:  "package-%v.%v",
				cur:  "package-1.2",
				past: []string{},
			},
			want: []string{
				"package-2.0",
				"package-1.3",
			},
			wantErr: false,
		}, {
			name: "simple-noprefix",
			fields: fields{
				fmt:  "%v.%v.tar.gz",
				cur:  "1.2.tar.gz",
				past: []string{},
			},
			want: []string{
				"2.0.tar.gz",
				"1.3.tar.gz",
			},
			wantErr: false,
		}, {
			name: "fourpart-normal",
			fields: fields{
				fmt:  "fourpart-%v.%v.%v.%v.tar.gz",
				cur:  "fourpart-5.3.1.2.tar.gz",
				past: []string{},
			},
			want: []string{
				"fourpart-6.0.0.0.tar.gz",
				"fourpart-5.4.0.0.tar.gz",
				"fourpart-5.3.2.0.tar.gz",
				"fourpart-5.3.1.3.tar.gz",
			},
			wantErr: false,
		}, {
			name: "simple-stringSubErr",
			fields: fields{
				fmt:  "%v.%v.tar.gz",
				cur:  "blah-12targz",
				past: []string{},
			},
			want:    []string{},
			wantErr: true,
		}, {
			name: "simple-stringConvErr",
			fields: fields{
				fmt:  "blah-%v.%v.tar.gz",
				cur:  "blah-t.a.tar.gz",
				past: []string{},
			},
			want:    []string{},
			wantErr: true,
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
