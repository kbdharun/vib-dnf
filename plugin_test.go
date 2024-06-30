package main

import (
	"encoding/json"
	"testing"

	"github.com/vanilla-os/vib/api"
)

type testCases struct {
	module   interface{}
	expected string
}

var test = []testCases{
	// Existing test cases
	{DnfModule{Name: "Single Package, Single Flag", Type: "dnf", Options: DnfOptions{ExtraFlags: []string{"--verbose"}}, Source: api.Source{Packages: []string{"bash"}}}, "dnf install -y --verbose bash && dnf clean all"},
	{DnfModule{Name: "Single Package, No Flag", Type: "dnf", Options: DnfOptions{ExtraFlags: []string{}}, Source: api.Source{Packages: []string{"bash"}}}, "dnf install -y  bash && dnf clean all"},
	{DnfModule{Name: "Multiple Packages, No Flag", Type: "dnf", Options: DnfOptions{ExtraFlags: []string{}}, Source: api.Source{Packages: []string{"bash", "fish"}}}, "dnf install -y  bash fish && dnf clean all"},
	{DnfModule{Name: "Multiple Packages, Multiple Flags", Type: "dnf", Options: DnfOptions{ExtraFlags: []string{"--verbose", "--best"}}, Source: api.Source{Packages: []string{"bash", "fish"}}}, "dnf install -y --verbose --best bash fish && dnf clean all"},

	// New test case for path-based installation
	{DnfModule{Name: "Path-based Installation", Type: "dnf", Options: DnfOptions{ExtraFlags: []string{"--assumeyes"}}, Source: api.Source{Paths: []string{"test.inst"}}}, "dnf install -y --assumeyes package1 package2 package3 && dnf clean all"},
}

func TestBuildModule(t *testing.T) {
	for _, testCase := range test {
		moduleInterface, err := json.Marshal(testCase.module)
		if err != nil {
			t.Errorf("Error in json %s", err.Error())
		}
		if output := BuildModule(convertToCString(string(moduleInterface)), convertToCString("")); convertToGoString(output) != testCase.expected {
			t.Errorf("Output %s not equivalent to expected %s", convertToGoString(output), testCase.expected)
		}
	}
}
