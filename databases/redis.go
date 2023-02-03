// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.
//

package databases

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

type RedisCommand struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RedisConnection struct {
	client *redis.Client
}

type Redis interface {
	Get(*RedisCommand) (*RedisCommand, error)
	Set(*RedisCommand) error
}

func InitializeRedis(cfg *Database) (*RedisConnection, error) {

	if cfg == nil {
		return nil, ErrConfigParameterMissing
	}

	// already initialized
	if rdb != nil {
		return &RedisConnection{client: rdb}, nil
	}
	url, err := redis.ParseURL(cfg.Conn)
	if err != nil {
		return nil, err
	}
	rdb = redis.NewClient(url)
	err = rdb.Ping(context.Background()).Err()

	if err != nil {
		return nil, err
	}

	return &RedisConnection{client: rdb}, nil
}

func (r *RedisConnection) Set(cmd *RedisCommand) error {
	if s := r.client.Set(context.Background(), cmd.Key, cmd.Value, 0); s.Err() != nil {
		return s.Err()
	}

	return nil
}

func (r *RedisConnection) Get(cmd *RedisCommand) (*RedisCommand, error) {
	s := r.client.Get(context.Background(), cmd.Key)
	if s.Err() != nil {
		return nil, s.Err()
	}

	return &RedisCommand{
		Key:   cmd.Key,
		Value: s.Val(),
	}, nil
}
