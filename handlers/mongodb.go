// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.

package handlers

import (
	"encoding/json"
	"getircase/databases"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Code    int                        `json:"code"`
	Msg     string                     `json:"msg"`
	Records []*databases.MongodbRecord `json:"records,omitempty"`
}

type mongodbHandler struct {
	client databases.MongoClient
}

func NewMongodbHandler(c databases.MongoClient) *mongodbHandler {
	return &mongodbHandler{
		client: c,
	}
}

func (h *mongodbHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		createFailResponse(rw, http.StatusMethodNotAllowed, ErrInvalidRequestMethod)
		return
	}
	if r.Header.Get("content-type") != "application/json" {
		createFailResponse(rw, http.StatusUnsupportedMediaType, ErrInvalidContentType)
		return
	}

	h.Retrieve(rw, r)
}

func (h *mongodbHandler) Retrieve(rw http.ResponseWriter, r *http.Request) {
	var filter = &databases.MongodbFilter{}
	if r.ContentLength != 0 {
		f, err := ioutil.ReadAll(r.Body)
		if err != nil {
			createFailResponse(rw, http.StatusInternalServerError, err)
			return
		}

		if err := json.Unmarshal(f, filter); err != nil {
			createFailResponse(rw, http.StatusInternalServerError, ErrMarshalError)
			return
		}
	}

	records, err := h.client.Fetch(filter)
	if err != nil {
		createFailResponse(rw, http.StatusInternalServerError, ErrFetchError)
		return
	}
	resp := createSuccessResponse(records)
	d, err := json.Marshal(resp)
	if err != nil {
		createFailResponse(rw, http.StatusInternalServerError, ErrMarshalError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(d)
}

func createSuccessResponse(records []*databases.MongodbRecord) *Response {
	return &Response{
		Code:    0,
		Msg:     "success",
		Records: records,
	}
}

func createFailResponse(rw http.ResponseWriter, statusCode int, msg error) {
	rsp := &Response{
		Code:    1,
		Msg:     msg.Error(),
		Records: nil,
	}
	d, err := json.Marshal(rsp)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(statusCode)
	rw.Write(d)
}
