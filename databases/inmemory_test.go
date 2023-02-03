// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.

package databases

import (
	"reflect"
	"sync"
	"testing"
)

func tearUp() {
	inmemory = &sync.Map{}
}

func tearDown() {
	inmemory = nil
}

func Test_sS_Get(t *testing.T) {
	type args struct {
		cmd *InmemoryCommand
	}
	tests := []struct {
		name    string
		args    args
		want    *InmemoryCommand
		wantErr bool
		setup   func()
	}{
		{
			name:    "get / failed",
			args:    args{cmd: &InmemoryCommand{Key: "test"}},
			want:    nil,
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "get / success",
			args: args{cmd: &InmemoryCommand{Key: "test"}},
			want: &InmemoryCommand{
				Key:   "test",
				Value: "testinmemory",
			},
			setup: func() {
				inmemory.LoadOrStore("test", "testinmemory")
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tearUp()
			defer tearDown()
			tt.setup()
			s := &sS{}
			got, err := s.Get(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("sS.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sS.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sS_Set(t *testing.T) {
	type args struct {
		cmd *InmemoryCommand
	}
	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr bool
	}{
		{
			name:    "set / failed",
			wantErr: true,
			setup:   func() {},
			args: args{
				cmd: &InmemoryCommand{Key: "test", Value: "testinmemory"},
			},
		},
		{
			name:    "set / success",
			wantErr: false,
			setup: func() {
				tearUp()
			},
			args: args{
				cmd: &InmemoryCommand{Key: "test", Value: "testinmemory"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sS{}
			tt.setup()
			if err := s.Set(tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("sS.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitializeInmemory(t *testing.T) {
	type args struct {
		cfg *Database
	}
	tests := []struct {
		name    string
		args    args
		want    Inmemory
		wantErr bool
		setup   func()
	}{
		{
			name:    "initalize / failed",
			args:    args{cfg: nil},
			want:    nil,
			setup:   tearDown,
			wantErr: true,
		},
		{
			name:    "initalize / success / no setup",
			args:    args{cfg: &Database{}},
			want:    &sS{},
			setup:   tearDown,
			wantErr: false,
		},
		{
			name:    "initalize / success / setup",
			args:    args{cfg: &Database{}},
			want:    &sS{},
			setup:   tearUp,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got, err := InitializeInmemory(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitializeInmemory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeInmemory() = %v, want %v", got, tt.want)
			}
		})
	}
}
