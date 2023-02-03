// Copyright (C) 2023 Timu Eren
//
// This file is part of getir-case.
//

package config

import (
	"encoding/json"
	"io/ioutil"
)

// Parse Parse configuration json file
func Parse(filePath *string) (*Configuration, error) {
	if filePath == nil {
		return nil, ErrConfigfileParameteMissing
	}
	content, err := ioutil.ReadFile(*filePath)
	if err != nil {
		return nil, err
	}
	cfg := &Configuration{}
	err = json.Unmarshal(content, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil

}
