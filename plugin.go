// Copyright 2023 - 2023, axtlos <axtlos@disroot.org>
// Copyright 2023 - present, K.B.Dharun Krishna <mail@kbdharun.dev>
// SPDX-License-Identifier: GPL-3.0-ONLY

package main

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/vanilla-os/vib/api"
)

type DnfModule struct {
	Name string `json:"name"`
	Type string `json:"type"`

	ExtraFlags []string
	Packages   []string
}

func BuildModule(moduleInterface interface{}, _ *api.Recipe) (string, error) {
	var module DnfModule
	err := mapstructure.Decode(moduleInterface, &module)
	if err != nil {
		return "", err
	}

	cmd := fmt.Sprintf("dnf install -y %s %s", strings.Join(module.ExtraFlags, " "), strings.Join(module.Packages, " "))

	return cmd, nil
}
