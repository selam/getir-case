// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.
package databases

import "errors"

var ErrConfigParameterMissing = errors.New("configuration value is nil")
var ErrInmemoryKeyNotFound = errors.New("inmemory: nil")
var ErrInmemoryOperationFailed = errors.New("inmemory: operation failed")
var ErrInmemoryInitializeFirst = errors.New("inmemory: initialize inmemory first")
