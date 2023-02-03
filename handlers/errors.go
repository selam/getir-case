// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.
package handlers

import (
	"errors"
	"net/http"
	"strings"
)

var ErrInvalidDateFormat = errors.New("invalid date format")
var ErrKeyEmpty = errors.New("key can not be empty")
var ErrInvalidInput = errors.New("invalid json input")
var ErrInvalidContentType = errors.New("invalid content-type")
var ErrInvalidRequestMethod = errors.New(strings.ToLower(http.StatusText(http.StatusMethodNotAllowed)))
var ErrFetchError = errors.New("mongodb: fetch error")
var ErrMarshalError = errors.New("json: marshal")
