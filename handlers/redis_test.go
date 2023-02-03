// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.
//

package handlers

import (
	"bytes"
	"errors"
	"getircase/databases"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockRedis struct {
	g func(*databases.RedisCommand) (*databases.RedisCommand, error)
	s func(*databases.RedisCommand) error
}

func (m *mockRedis) Get(cmd *databases.RedisCommand) (*databases.RedisCommand, error) {
	return m.g(cmd)
}
func (m *mockRedis) Set(cmd *databases.RedisCommand) error {
	return m.s(cmd)
}
func TestRedisHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		client databases.Redis
	}
	type args struct {
		method      string
		path        string
		contentType string
		body        io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "redis patch",
			args: args{
				method:      http.MethodPatch,
				path:        "/inmemory",
				contentType: "application/json",
			},
			want: `{"error": "method not allowed"}`,
			fields: fields{client: &mockRedis{
				g: func(ic *databases.RedisCommand) (*databases.RedisCommand, error) {
					return nil, errors.New("redis: nil")
				},
			}},
		},
		{
			name: "redis get / not exists",
			args: args{
				method:      http.MethodGet,
				path:        "/redis?key=not-exists",
				contentType: "application/json",
			},
			want: `{"error": "redis: nil"}`,
			fields: fields{client: &mockRedis{
				g: func(ic *databases.RedisCommand) (*databases.RedisCommand, error) {
					return nil, errors.New("redis: nil")
				},
			}},
		},
		{
			name: "redis get / exists",
			args: args{
				method:      http.MethodGet,
				path:        "/redis?key=exists",
				contentType: "application/json",
			},
			want: `{"key":"exists","value":"exists"}`,
			fields: fields{client: &mockRedis{
				g: func(ic *databases.RedisCommand) (*databases.RedisCommand, error) {
					return &databases.RedisCommand{Key: "exists", Value: "exists"}, nil
				},
			}},
		},
		{
			name: "redis post / empty body",
			args: args{
				method:      http.MethodPost,
				path:        "/redis",
				contentType: "application/json",
			},
			want: `{"error": "invalid json input"}`,
			fields: fields{client: &mockRedis{
				g: func(ic *databases.RedisCommand) (*databases.RedisCommand, error) {
					return nil, nil
					//return &databases.InmemoryCommand{Key: "exists", Value: "exists"}, nil
				},
				s: func(ic *databases.RedisCommand) error {
					return nil
				},
			}},
		},
		{
			name: "redis post / wrong content type",
			args: args{
				method:      http.MethodPost,
				path:        "/redis",
				contentType: "text/html",
			},
			want: `{"error": "invalid content-type"}`,
			fields: fields{client: &mockRedis{
				g: func(ic *databases.RedisCommand) (*databases.RedisCommand, error) {
					return nil, nil
				},
				s: func(ic *databases.RedisCommand) error {
					return nil
				},
			}},
		},
		{
			name: "redis post / valid body",
			args: args{
				method:      http.MethodPost,
				path:        "/redis",
				contentType: "application/json",
				body:        bytes.NewBufferString(`{"key": "test","value":"test"}`),
			},
			want: `{"key":"test","value":"test"}`,
			fields: fields{client: &mockRedis{
				g: func(ic *databases.RedisCommand) (*databases.RedisCommand, error) {
					return ic, nil
				},
				s: func(ic *databases.RedisCommand) error {
					return nil
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.args.method, tt.args.path, tt.args.body)
			r.Header.Add("Content-Type", tt.args.contentType)
			rw := httptest.NewRecorder()
			h := NewRedisHandler(tt.fields.client)
			h.ServeHTTP(rw, r)
			if rw.Body.String() != tt.want {
				t.Errorf("ServeHTTP() = %s, want %s", rw.Body.String(), tt.want)
			}
		})
	}
}
