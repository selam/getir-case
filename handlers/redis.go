// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.
//

package handlers

import (
	"encoding/json"
	"fmt"
	"getircase/databases"
	"io/ioutil"
	"net/http"
)

type RedisHandler struct {
	client databases.Redis
}

func NewRedisHandler(client databases.Redis) *RedisHandler {
	return &RedisHandler{client: client}
}

func (h *RedisHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		writeError(rw, http.StatusMethodNotAllowed, ErrInvalidRequestMethod)
		return
	}

	if r.Method == http.MethodPost {
		if r.Header.Get("Content-Type") != "application/json" {
			writeError(rw, http.StatusUnsupportedMediaType, ErrInvalidContentType)
			return
		}

		h.CreateOrUpdate(rw, r)
		return
	}

	h.Get(rw, r)
}

func writeError(rw http.ResponseWriter, status int, err error) {
	rw.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err.Error())))
}

func (h *RedisHandler) CreateOrUpdate(rw http.ResponseWriter, r *http.Request) {
	command := &databases.RedisCommand{}
	if r.ContentLength != 0 {
		f, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(rw, http.StatusBadRequest, err)
			return
		}

		if err := json.Unmarshal(f, command); err != nil {
			writeError(rw, http.StatusBadRequest, err)
			return
		}
	}
	if command.Key == "" || command.Value == "" {
		writeError(rw, http.StatusBadRequest, ErrInvalidInput)
		return
	}
	if err := h.client.Set(command); err != nil {
		writeError(rw, http.StatusBadRequest, err)
		return
	}

	cmd, err := h.client.Get(command)
	if err != nil {
		writeError(rw, http.StatusBadRequest, err)
		return
	}

	b, err := json.Marshal(cmd)
	if err != nil {
		writeError(rw, http.StatusBadRequest, err)
		return
	}

	rw.Write(b)
}

func (h *RedisHandler) Get(rw http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		writeError(rw, http.StatusBadRequest, ErrKeyEmpty)
		return
	}
	command := &databases.RedisCommand{
		Key: key,
	}

	cmd, err := h.client.Get(command)
	if err != nil {
		writeError(rw, http.StatusBadRequest, err)
		return
	}

	b, err := json.Marshal(cmd)
	if err != nil {
		writeError(rw, http.StatusBadRequest, err)
		return
	}

	rw.Write(b)
}
