// Copyright 2023 - 2023, axtlos <axtlos@disroot.org>
// Copyright 2023 - present, K.B.Dharun Krishna <mail@kbdharun.dev>
// SPDX-License-Identifier: GPL-3.0-ONLY

package main

import (
	"C"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vanilla-os/vib/api"
)

type DnfModule struct {
	Name    string     `json:"name"`
	Type    string     `json:"type"`
	Options DnfOptions `json:"options"`
	Source  api.Source `json:"source"`
}

type DnfOptions struct {
	ExtraFlags []string `json:"extra_flags"`
}

//export BuildModule
func BuildModule(moduleInterface *C.char, recipeInterface *C.char) *C.char {
	var module *DnfModule
	var recipe *api.Recipe

	err := json.Unmarshal([]byte(C.GoString(moduleInterface)), &module)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	err = json.Unmarshal([]byte(C.GoString(recipeInterface)), &recipe)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	extraFlags := strings.Join(module.Options.ExtraFlags, " ")

	if len(module.Source.Packages) > 0 {
		packages := strings.Join(module.Source.Packages, " ")
		return C.CString(fmt.Sprintf("dnf install -y %s %s && dnf clean all", extraFlags, packages))
	}

	if len(module.Source.Paths) > 0 {
		cmd := ""
		for i, path := range module.Source.Paths {
			instPath := filepath.Join(recipe.ParentPath, path+".inst")
			packages := ""
			file, err := os.Open(instPath)
			if err != nil {
				return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				packages += scanner.Text() + " "
			}
			if err := scanner.Err(); err != nil {
				return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
			}
			cmd += fmt.Sprintf("dnf install -y %s %s", extraFlags, packages)
			if i != len(module.Source.Paths)-1 {
				cmd += "&& "
			} else {
				cmd += "&& dnf clean all"
			}
		}
		return C.CString(cmd)
	}

	return C.CString("ERROR: no packages or paths specified")
}

func main() {}
