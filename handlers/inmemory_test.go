// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.

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

type mockInmemory struct {
	g func(*databases.InmemoryCommand) (*databases.InmemoryCommand, error)
	s func(*databases.InmemoryCommand) error
}

func (m *mockInmemory) Get(cmd *databases.InmemoryCommand) (*databases.InmemoryCommand, error) {
	return m.g(cmd)
}
func (m *mockInmemory) Set(cmd *databases.InmemoryCommand) error {
	return m.s(cmd)
}

func TestInmemoryHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		client databases.Inmemory
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
			name: "inmemory patch",
			args: args{
				method:      http.MethodPatch,
				path:        "/inmemory",
				contentType: "application/json",
			},
			want: `{"error": "method not allowed"}`,
			fields: fields{client: &mockInmemory{
				g: func(ic *databases.InmemoryCommand) (*databases.InmemoryCommand, error) {
					return nil, errors.New("redis: nil")
				},
			}},
		},
		{
			name: "inmemory get / not exists",
			args: args{
				method:      http.MethodGet,
				path:        "/inmemory?key=not-exists",
				contentType: "application/json",
			},
			want: `{"error": "inmemory: nil"}`,
			fields: fields{client: &mockInmemory{
				g: func(ic *databases.InmemoryCommand) (*databases.InmemoryCommand, error) {
					return nil, errors.New("inmemory: nil")
				},
			}},
		},
		{
			name: "inmemory get / exists",
			args: args{
				method:      http.MethodGet,
				path:        "/inmemory?key=exists",
				contentType: "application/json",
			},
			want: `{"key":"exists","value":"exists"}`,
			fields: fields{client: &mockInmemory{
				g: func(ic *databases.InmemoryCommand) (*databases.InmemoryCommand, error) {
					return &databases.InmemoryCommand{Key: "exists", Value: "exists"}, nil
				},
			}},
		},
		{
			name: "inmemory post / wrong content type",
			args: args{
				method:      http.MethodPost,
				path:        "/inmemory",
				contentType: "text/html",
			},
			want: `{"error": "invalid content-type"}`,
			fields: fields{client: &mockInmemory{
				g: func(ic *databases.InmemoryCommand) (*databases.InmemoryCommand, error) {
					return nil, nil
				},
				s: func(ic *databases.InmemoryCommand) error {
					return nil
				},
			}},
		},
		{
			name: "inmemory post / empty body",
			args: args{
				method:      http.MethodPost,
				path:        "/inmemory",
				contentType: "application/json",
			},
			want: `{"error": "invalid json input"}`,
			fields: fields{client: &mockInmemory{
				g: func(ic *databases.InmemoryCommand) (*databases.InmemoryCommand, error) {
					return nil, nil
				},
				s: func(ic *databases.InmemoryCommand) error {
					return nil
				},
			}},
		},
		{
			name: "inmemory post / valid body",
			args: args{
				method:      http.MethodPost,
				path:        "/inmemory",
				contentType: "application/json",
				body:        bytes.NewBufferString(`{"key": "test","value":"test"}`),
			},
			want: `{"key":"test","value":"test"}`,
			fields: fields{client: &mockInmemory{
				g: func(ic *databases.InmemoryCommand) (*databases.InmemoryCommand, error) {
					return ic, nil
				},
				s: func(ic *databases.InmemoryCommand) error {
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
			h := NewInmemoryHandler(tt.fields.client)
			h.ServeHTTP(rw, r)
			if rw.Body.String() != tt.want {
				t.Errorf("ServeHTTP() = %s, want %s", rw.Body.String(), tt.want)
			}
		})
	}
}
