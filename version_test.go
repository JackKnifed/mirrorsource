package main

import (
	"reflect"
	"testing"
)

func Test_versionType_possibleUpgrades(t *testing.T) {
	type fields struct {
		Fmt   string
		Cur   string
		Local []string
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
				Fmt:   "memcached-%v.%v.%v.tar.gz",
				Cur:   "memcached-1.4.33.tar.gz",
				Local: []string{},
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
				Fmt:   "memcached-%v.%v.%v",
				Cur:   "memcached-1.4.33",
				Local: []string{},
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
				Fmt:   "package-%v.%v",
				Cur:   "package-1.2",
				Local: []string{},
			},
			want: []string{
				"package-2.0",
				"package-1.3",
			},
			wantErr: false,
		}, {
			name: "simple-noprefix",
			fields: fields{
				Fmt:   "%v.%v.tar.gz",
				Cur:   "1.2.tar.gz",
				Local: []string{},
			},
			want: []string{
				"2.0.tar.gz",
				"1.3.tar.gz",
			},
			wantErr: false,
		}, {
			name: "fourpart-normal",
			fields: fields{
				Fmt:   "fourpart-%v.%v.%v.%v.tar.gz",
				Cur:   "fourpart-5.3.1.2.tar.gz",
				Local: []string{},
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
				Fmt:   "%v.%v.tar.gz",
				Cur:   "blah-12targz",
				Local: []string{},
			},
			want:    []string{},
			wantErr: true,
		}, {
			name: "simple-stringConvErr",
			fields: fields{
				Fmt:   "blah-%v.%v.tar.gz",
				Cur:   "blah-t.a.tar.gz",
				Local: []string{},
			},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := versionType{
				Fmt:   tt.fields.Fmt,
				Cur:   tt.fields.Cur,
				Local: tt.fields.Local,
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
