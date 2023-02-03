// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.
//

package config

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		filePath *string
	}
	tests := []struct {
		name    string
		args    args
		want    *Configuration
		wantErr bool
	}{
		{
			name: "empty param/nil check",
			args: args{
				filePath: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrongname.json",
			args: args{
				filePath: func() *string {
					s := "wrongname.json"
					return &s
				}(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "config.json",
			args: args{
				filePath: func() *string {
					s := "data/test.json"
					return &s
				}(),
			},
			want: &Configuration{
				Application: Application{
					Host: "1.1.1.1",
					Port: 8080,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
