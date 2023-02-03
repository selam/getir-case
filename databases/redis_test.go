// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.
//

package databases

import (
	"reflect"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
)

func TestRedisConnection_Get(t *testing.T) {
	type fields struct {
		client *redis.Client
	}
	type args struct {
		cmd *RedisCommand
	}
	tests := []struct {
		name    string
		fields  *fields
		args    args
		want    *RedisCommand
		wantErr bool
	}{
		{
			name: "get / failed",
			fields: func() *fields {
				db, mock := redismock.NewClientMock()
				mock.ExpectGet("testresdis").RedisNil()
				return &fields{
					client: db,
				}
			}(),
			args: args{
				cmd: &RedisCommand{Key: "testredis"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "get / success",
			fields: func() *fields {
				db, mock := redismock.NewClientMock()
				mock.ExpectGet("testredis").SetVal("testvalue")
				return &fields{
					client: db,
				}
			}(),
			args: args{
				cmd: &RedisCommand{Key: "testredis", Value: "testvalue"},
			},
			want: &RedisCommand{
				Key:   "testredis",
				Value: "testvalue",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisConnection{
				client: tt.fields.client,
			}
			got, err := r.Get(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisConnection.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisConnection.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisConnection_Set(t *testing.T) {
	type fields struct {
		client *redis.Client
	}
	type args struct {
		cmd *RedisCommand
	}
	tests := []struct {
		name    string
		fields  *fields
		args    args
		want    *RedisCommand
		wantErr bool
	}{
		{
			name: "set / failed",
			fields: func() *fields {
				db, mock := redismock.NewClientMock()
				mock.ExpectSet("testresdis", "testredis", 0).RedisNil()
				return &fields{
					client: db,
				}
			}(),
			args: args{
				cmd: &RedisCommand{Key: "testredis", Value: "testredis"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "set / success",
			fields: func() *fields {
				db, mock := redismock.NewClientMock()
				mock.ExpectSet("testredis", "testvalue", 0).SetVal("")
				return &fields{
					client: db,
				}
			}(),
			args: args{
				cmd: &RedisCommand{Key: "testredis", Value: "testvalue"},
			},
			want: &RedisCommand{
				Key:   "testredis",
				Value: "testvalue",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisConnection{
				client: tt.fields.client,
			}
			if err := r.Set(tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("RedisConnection.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitializeRedis(t *testing.T) {
	type args struct {
		cfg *Database
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "failed / no cfg",
			args:    args{cfg: nil},
			want:    false,
			wantErr: true,
		},
		{
			name: "failed / with cfg",
			args: args{cfg: &Database{
				Conn: "notexists://",
			}},
			want:    false,
			wantErr: true,
		},
		{
			name: "falied / with cfg",
			args: args{cfg: &Database{
				Conn: "redis://127.0.0.99/0",
			}},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitializeRedis(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitializeRedis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.want {
				t.Errorf("InitializeRedis() = %v, want %v", got, tt.want)
			}
		})
	}
}
