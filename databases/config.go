// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.

package databases

type Database struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Conn string `json:"connection_string"`
}
