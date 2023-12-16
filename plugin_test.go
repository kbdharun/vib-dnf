// Copyright 2023 - 2023, axtlos <axtlos@disroot.org>
// Copyright 2023 - present, K.B.Dharun Krishna <mail@kbdharun.dev>
// SPDX-License-Identifier: GPL-3.0-ONLY

package main

import (
	"testing"

	"github.com/vanilla-os/vib/api"
)

type testModule struct {
	Name       string
	Type       string
	ExtraFlags []string
	Packages   []string
}

type testCases struct {
	module   interface{}
	expected string
}

var test = []testCases{
	{testModule{"Single Package, Single Flag", "dnf", []string{"--verbose"}, []string{"bash"}}, "dnf install --verbose bash"},
	{testModule{"Single Package, No Flag", "dnf", []string{""}, []string{"bash"}}, "dnf install bash"},
	{testModule{"Multiple Packages, No Flag", "dnf", []string{""}, []string{"bash", "fish"}}, "dnf install bash fish"},
	{testModule{"Multiple Packages, Multiple Flags", "dnf", []string{"--verbose", "-y"}, []string{"bash", "fish"}}, "dnf install --verbose -y bash fish"},
}

func TestBuildModule(t *testing.T) {
	for _, testCase := range test {
		if output, _ := BuildModule(testCase.module, &api.Recipe{}); output != testCase.expected {
			t.Errorf("Output %q not equivalent to expected %q", output, testCase.expected)
		}
	}
}
