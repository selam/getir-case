// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.

package handlers

import (
	"bytes"
	"getircase/databases"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockMongo struct {
	f func(*databases.MongodbFilter) ([]*databases.MongodbRecord, error)
}

func (m *mockMongo) Fetch(f *databases.MongodbFilter) ([]*databases.MongodbRecord, error) {
	return m.f(f)
}

func Test_mongodbHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		client databases.MongoClient
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
			name: "mongo patch",
			args: args{
				method:      http.MethodPatch,
				path:        "/mongodb/recods",
				contentType: "application/json",
			},
			want: `{"code":1,"msg":"method not allowed"}`,
			fields: fields{client: &mockMongo{
				f: func(ic *databases.MongodbFilter) ([]*databases.MongodbRecord, error) {
					return nil, nil
				},
			}},
		},
		{
			name: "mongo post / wrong content type",
			args: args{
				method:      http.MethodPost,
				path:        "/mongodb/recods",
				contentType: "text/html",
			},
			want: `{"code":1,"msg":"invalid content-type"}`,
			fields: fields{client: &mockMongo{
				f: func(ic *databases.MongodbFilter) ([]*databases.MongodbRecord, error) {
					return nil, nil
				},
			}},
		},
		{
			name: "mongo post / empty body",
			args: args{
				method:      http.MethodPost,
				path:        "/mongodb/records",
				contentType: "application/json",
			},
			want: `{"code":0,"msg":"success","records":[{"key":"a","createdAt":"","totalCount":1}]}`,
			fields: fields{client: &mockMongo{
				f: func(ic *databases.MongodbFilter) ([]*databases.MongodbRecord, error) {
					return []*databases.MongodbRecord{
						&databases.MongodbRecord{
							Key:        "a",
							CreatedAt:  "",
							TotalCount: 1,
						},
					}, nil
				},
			}},
		},
		{
			name: "mongo post / invalid body / start date / parse error / 3333-33-33",
			args: args{
				method:      http.MethodPost,
				path:        "/mongodb/records",
				contentType: "application/json",
				body:        bytes.NewBufferString(`{"startDate":"3333-33-33"}`),
			},
			want: `{"code":1,"msg":"json: marshal"}`,
			fields: fields{client: &mockMongo{
				f: func(ic *databases.MongodbFilter) ([]*databases.MongodbRecord, error) {
					return []*databases.MongodbRecord{
						&databases.MongodbRecord{
							Key:        "a",
							CreatedAt:  "",
							TotalCount: 1,
						},
					}, nil
				},
			}},
		},
		{
			name: "mongo post / invalid body / start date / parse error / 33-3333-33",
			args: args{
				method:      http.MethodPost,
				path:        "/mongodb/records",
				contentType: "application/json",
				body:        bytes.NewBufferString(`{"startDate":"33-3333-33"}`),
			},
			want: `{"code":1,"msg":"json: marshal"}`,
			fields: fields{client: &mockMongo{
				f: func(ic *databases.MongodbFilter) ([]*databases.MongodbRecord, error) {
					return []*databases.MongodbRecord{
						&databases.MongodbRecord{
							Key:        "a",
							CreatedAt:  "",
							TotalCount: 1,
						},
					}, nil
				},
			}},
		},
		{
			name: "mongo post / start date 2022-01-22",
			args: args{
				method:      http.MethodPost,
				path:        "/mongodb/records",
				contentType: "application/json",
				body:        bytes.NewBufferString(`{"startDate":"2022-01-22"}`),
			},
			want: `{"code":0,"msg":"success","records":[{"key":"a","createdAt":"","totalCount":1}]}`,
			fields: fields{client: &mockMongo{
				f: func(ic *databases.MongodbFilter) ([]*databases.MongodbRecord, error) {
					return []*databases.MongodbRecord{
						&databases.MongodbRecord{
							Key:        "a",
							CreatedAt:  "",
							TotalCount: 1,
						},
					}, nil
				},
			}},
		},
		{
			name: "mongo post / invalid body / minCount",
			args: args{
				method:      http.MethodPost,
				path:        "/mongodb/records",
				contentType: "application/json",
				body:        bytes.NewBufferString(`{"minCount":"1a"}`),
			},
			want: `{"code":1,"msg":"json: marshal"}`,
			fields: fields{client: &mockMongo{
				f: func(ic *databases.MongodbFilter) ([]*databases.MongodbRecord, error) {
					return []*databases.MongodbRecord{
						&databases.MongodbRecord{
							Key:        "a",
							CreatedAt:  "",
							TotalCount: 1,
						},
					}, nil
				},
			}},
		},
		{
			name: "mongo post / invalid body / maxCount",
			args: args{
				method:      http.MethodPost,
				path:        "/mongodb/records",
				contentType: "application/json",
				body:        bytes.NewBufferString(`{"maxCount":"1a"}`),
			},
			want: `{"code":1,"msg":"json: marshal"}`,
			fields: fields{client: &mockMongo{
				f: func(ic *databases.MongodbFilter) ([]*databases.MongodbRecord, error) {
					return []*databases.MongodbRecord{
						&databases.MongodbRecord{
							Key:        "a",
							CreatedAt:  "",
							TotalCount: 1,
						},
					}, nil
				},
			}},
		},
		{
			name: "mongo post / minCount",
			args: args{
				method:      http.MethodPost,
				path:        "/mongodb/records",
				contentType: "application/json",
				body:        bytes.NewBufferString(`{"minCount":1}`),
			},
			want: `{"code":0,"msg":"success","records":[{"key":"a","createdAt":"","totalCount":1}]}`,
			fields: fields{client: &mockMongo{
				f: func(ic *databases.MongodbFilter) ([]*databases.MongodbRecord, error) {
					return []*databases.MongodbRecord{
						&databases.MongodbRecord{
							Key:        "a",
							CreatedAt:  "",
							TotalCount: 1,
						},
					}, nil
				},
			}},
		},
		{
			name: "mongo post / maxCount",
			args: args{
				method:      http.MethodPost,
				path:        "/mongodb/records",
				contentType: "application/json",
				body:        bytes.NewBufferString(`{"maxCount":1}`),
			},
			want: `{"code":0,"msg":"success","records":[{"key":"a","createdAt":"","totalCount":1}]}`,
			fields: fields{client: &mockMongo{
				f: func(ic *databases.MongodbFilter) ([]*databases.MongodbRecord, error) {
					return []*databases.MongodbRecord{
						&databases.MongodbRecord{
							Key:        "a",
							CreatedAt:  "",
							TotalCount: 1,
						},
					}, nil
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.args.method, tt.args.path, tt.args.body)
			r.Header.Add("Content-Type", tt.args.contentType)
			rw := httptest.NewRecorder()
			h := NewMongodbHandler(tt.fields.client)
			h.ServeHTTP(rw, r)
			if rw.Body.String() != tt.want {
				t.Errorf("ServeHTTP() = %s, want %s", rw.Body.String(), tt.want)
			}

		})
	}
}
