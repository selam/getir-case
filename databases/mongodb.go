// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.
//

package databases

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Time time.Time

func (c *Time) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", value) //parse time
	if err != nil {
		return err
	}
	*c = Time(t) //set result using the pointer
	return nil
}

func (c *Time) Time() time.Time {
	return time.Time(*c)
}

func (c Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format("2006-01-02") + `"`), nil
}

type MongodbRecord struct {
	Key        string `bson:"_id" json:"key"`
	CreatedAt  string `bson:"createdAt" json:"createdAt"`
	TotalCount int    `bson:"totalCount" json:"totalCount"`
}

type MongodbFilter struct {
	StartDate *Time `json:"startDate"`
	EndDate   *Time `json:"endDate"`
	MinCount  *int  `json:"minCount"`
	MaxCount  *int  `json:"maxCount"`
}

type MongoClient interface {
	Fetch(*MongodbFilter) ([]*MongodbRecord, error)
}

// wrap mongo client to write more easy tests
type mClient struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

var client *mongo.Client
var database *mongo.Database

func InitializeMongodb(cfg *Database) (*mClient, error) {
	if cfg == nil {
		return nil, ErrConfigParameterMissing
	}

	// already initialized
	if client != nil {
		return &mClient{client: client, database: database}, nil
	}
	cOptions := options.Client().ApplyURI(cfg.Conn)
	client, err := mongo.NewClient(cOptions)
	if err != nil {
		return nil, err
	}
	if err := client.Connect(context.Background()); err != nil {
		return nil, err
	}

	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		return nil, err
	}
	database := client.Database(cfg.Name)
	return &mClient{client: client, database: database, collection: database.Collection("records")}, nil
}

func (c *mClient) Fetch(f *MongodbFilter) ([]*MongodbRecord, error) {
	createdAtFilter := bson.M{}
	minMaxFilter := bson.M{}
	if f.EndDate != nil || f.StartDate != nil {
		m := make(map[string]interface{})

		if f.StartDate != nil {
			m["$gte"] = primitive.NewDateTimeFromTime(f.StartDate.Time())
		}
		if f.EndDate != nil {
			m["$lte"] = primitive.NewDateTimeFromTime(f.EndDate.Time())
		}
		createdAtFilter = bson.M{"created_at": m}
	}
	if f.MinCount != nil || f.MaxCount != nil {
		m := make(map[string]interface{})
		if f.MinCount != nil {
			m["$gte"] = *f.MinCount
		}
		if f.MaxCount != nil {
			m["$lte"] = *f.MaxCount
		}
		minMaxFilter = bson.M{"count": m}
	}

	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$key"},
			{"totalCount", bson.D{{"$sum", "$count"}}},
		},
		},
	}
	matchStage := bson.D{
		{"$match", bson.D{
			{"$and", []bson.M{createdAtFilter, minMaxFilter}}}},
	}
	cursor, err := c.collection.Aggregate(context.TODO(), mongo.Pipeline{
		matchStage,
		groupStage})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var records []*MongodbRecord
	if err := cursor.All(context.TODO(), &records); err != nil {
		return nil, err
	}
	return records, nil
}
