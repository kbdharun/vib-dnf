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

// Configuration for a DNF module
type DnfModule struct {
	Name    string       `json:"name"`
	Type    string       `json:"type"`
	Options DnfOptions   `json:"options"`
	Sources []api.Source `json:"sources"`
}

// Options for DNF package management
type DnfOptions struct {
	AllowErasing    bool     `json:"allowerasing"`
	SkipBroken      bool     `json:"skip_broken"`
	SkipUnavailable bool     `json:"skip_unavailable"`
	AllowDowngrade  bool     `json:"allow_downgrade"`
	DownloadOnly    bool     `json:"downloadonly"`
	Security        bool     `json:"security"`
	Bugfix          bool     `json:"bugfix"`
	Enhancement     bool     `json:"enhancement"`
	ExtraFlags      []string `json:"extra_flags"`
}

// Provide plugin information as a JSON string
//
//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "dnf", Type: api.BuildPlugin, UseContainerCmds: false}
	pluginjson, err := json.Marshal(plugininfo)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	return C.CString(string(pluginjson))
}

// Generate a dnf install command from the provided module and recipe.
// Handle package installation and apply appropriate options.
//
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

	args := ""
	if module.Options.AllowErasing {
		args += "--allowerasing "
	}
	if module.Options.SkipBroken {
		args += "--skip-broken "
	}
	if module.Options.SkipUnavailable {
		args += "--skip-unavailable "
	}
	if module.Options.AllowDowngrade {
		args += "--allow-downgrade "
	}
	if module.Options.DownloadOnly {
		args += "--downloadonly "
	}
	if module.Options.Security {
		args += "--security "
	}
	if module.Options.Bugfix {
		args += "--bugfix "
	}
	if module.Options.Enhancement {
		args += "--enhancement "
	}
	if len(module.Options.ExtraFlags) > 0 {
		args += strings.Join(module.Options.ExtraFlags, " ") + " "
	}

	packages := ""
	for _, source := range module.Sources {
		if len(source.Packages) > 0 {
			for _, pkg := range source.Packages {
				packages += pkg + " "
			}
		}

		if len(strings.TrimSpace(source.Path)) > 0 {
			filePath := source.Path
			if !filepath.IsAbs(source.Path) && recipe != nil && recipe.ParentPath != "" {
				filePath = filepath.Join(recipe.ParentPath, source.Path)
			}

			fileInfo, err := os.Stat(filePath)
			if err != nil {
				return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
			}
			if !fileInfo.Mode().IsRegular() {
				continue
			}
			file, err := os.Open(filePath)
			if err != nil {
				return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" && !strings.HasPrefix(line, "#") {
					packages += line + " "
				}
			}

			if err := scanner.Err(); err != nil {
				return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
			}
		}
	}

	if len(packages) >= 1 {
		cmd := fmt.Sprintf("dnf install -y %s %s && dnf clean all", args, packages)
		return C.CString(cmd)
	}

	return C.CString("ERROR: no packages or paths specified")
}

func main() {}
