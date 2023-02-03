// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.
//

package databases

import (
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func Test_mClient_Fetch(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	type args struct {
		f *MongodbFilter
	}
	type base struct {
		want []*MongodbRecord
		mock func(mt *mtest.T)
	}
	type test struct {
		base    base
		name    string
		args    args
		wantErr bool
	}

	baseTest := base{
		want: []*MongodbRecord{
			&MongodbRecord{
				Key:        "testing",
				TotalCount: 3,
				CreatedAt:  "2023-01-02",
			},
			&MongodbRecord{
				Key:        "testing2",
				TotalCount: 2,
				CreatedAt:  "2023-01-02",
			},
		},
		mock: func(mt *mtest.T) {
			first := mtest.CreateCursorResponse(1, "test.recods", mtest.FirstBatch, bson.D{
				{"_id", "testing"},
				{"totalCount", 3},
				{"createdAt", "2023-01-02"},
			})
			getMore := mtest.CreateCursorResponse(1, "test.recods", mtest.NextBatch, bson.D{
				{"_id", "testing2"},
				{"totalCount", 2},
				{"createdAt", "2023-01-02"},
			})
			killCursors := mtest.CreateCursorResponse(0, "test.recods", mtest.NextBatch)
			mt.AddMockResponses(first, getMore, killCursors)

		},
	}

	tests := []test{
		{
			name: "full records",
			args: args{
				f: &MongodbFilter{},
			},
			base: baseTest,
		},
		{
			name: "full records / with min",
			args: args{
				f: &MongodbFilter{
					MinCount: func() *int {
						i := 1
						return &i
					}(),
				},
			},
			base: baseTest,
		},
		{
			name: "full records / with max",
			args: args{
				f: &MongodbFilter{
					MaxCount: func() *int {
						i := 10
						return &i
					}(),
				},
			},
			base: baseTest,
		},
		{
			name: "full records / with startDate",
			args: args{
				f: &MongodbFilter{
					StartDate: func() *Time {
						n := Time(time.Now())
						return &n
					}(),
				},
			},
			base: baseTest,
		},
		{
			name: "full records / with endDate",
			args: args{
				f: &MongodbFilter{
					EndDate: func() *Time {
						n := Time(time.Now())
						return &n
					}(),
				},
			},
			base: baseTest,
		},
	}
	for _, tt := range tests {
		mt.RunOpts(tt.name, mtest.NewOptions().DatabaseName("test").CollectionName("records"), func(mt *mtest.T) {
			tt.base.mock(mt)
			c := &mClient{
				client:     mt.Client,
				database:   mt.DB,
				collection: mt.Coll,
			}

			got, err := c.Fetch(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("mClient.Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.base.want) {
				t.Errorf("mClient.Fetch() = %v, want %v", got, tt.base.want)
			}
		})
	}
}

func TestTime_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "deneme",
			args:    args{b: []byte("2023-02-12")},
			wantErr: false,
		},
		{
			name:    "wrong date",
			args:    args{b: []byte("2023-02-30")},
			wantErr: true,
		},
		{
			name:    "wrong format",
			args:    args{b: []byte("2023-34-12")},
			wantErr: true,
		},
		{
			name:    "null string",
			args:    args{b: []byte("null")},
			wantErr: false,
		},
		{
			name:    "empty string",
			args:    args{b: []byte("")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Time{}
			if err := c.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Time.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTime_Time(t *testing.T) {
	tests := []struct {
		name string
		want time.Time
	}{
		{
			name: "empty",
			want: time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Time{}
			if got := c.Time(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Time.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		{
			name:    "empty",
			want:    []byte("\"0001-01-01\""),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Time{}
			got, err := c.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Time.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Time.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
