// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.
//

package config

import (
	"getircase/databases"
)

type Configuration struct {
	Application Application           `json:"application"`
	Databases   []*databases.Database `json:"databases"`
}

type Application struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}
