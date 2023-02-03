// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.

package databases

import "sync"

var inmemory *sync.Map

type InmemoryCommand struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Inmemory interface {
	Get(*InmemoryCommand) (*InmemoryCommand, error)
	Set(*InmemoryCommand) error
}

type sS struct{}

func (s *sS) Get(cmd *InmemoryCommand) (*InmemoryCommand, error) {
	val, ok := inmemory.Load(cmd.Key)
	if !ok {
		return nil, ErrInmemoryKeyNotFound
	}
	return &InmemoryCommand{
		Key:   cmd.Key,
		Value: val.(string),
	}, nil
}

func (s *sS) Set(cmd *InmemoryCommand) error {
	if inmemory == nil {
		return ErrInmemoryInitializeFirst
	}
	inmemory.Store(cmd.Key, cmd.Value)

	return nil
}

func InitializeInmemory(cfg *Database) (Inmemory, error) {
	if cfg == nil {
		return nil, ErrConfigParameterMissing
	}
	// already initialized
	if inmemory != nil {
		return &sS{}, nil
	}

	inmemory = &sync.Map{}

	return &sS{}, nil
}
