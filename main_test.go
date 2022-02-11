// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

func run(t *testing.T, name string, args ...string) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, out)
	}
}

func tempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "license")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestInitial(t *testing.T) {
	if os.Getenv("RUNME") != "" {
		main()
		return
	}

	tmp := tempDir(t)
	t.Logf("tmp dir: %s", tmp)
	run(t, "cp", "-r", "testdata/initial", tmp)

	// run at least 2 times to ensure the program is idempotent
	for i := 0; i < 2; i++ {
		t.Logf("run #%d", i)
		targs := []string{"-test.run=TestInitial"}
		cargs := []string{"-l", "apache", "-c", "Bhojpur Consulting Private Limited, India", "-y", "2018", tmp}
		c := exec.Command(os.Args[0], append(targs, cargs...)...)
		c.Env = []string{"RUNME=1"}
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("%v\n%s", err, out)
		}

		run(t, "diff", "-r", filepath.Join(tmp, "initial"), "testdata/expected")
	}
}

func TestMultiyear(t *testing.T) {
	if os.Getenv("RUNME") != "" {
		main()
		return
	}

	tmp := tempDir(t)
	t.Logf("tmp dir: %s", tmp)
	samplefile := filepath.Join(tmp, "file.c")
	const sampleLicensed = "testdata/multiyear_file.c"

	run(t, "cp", "testdata/initial/file.c", samplefile)
	cmd := exec.Command(os.Args[0],
		"-test.run=TestMultiyear",
		"-l", "bsd", "-c", "Bhojpur Consulting Private Limited, India.",
		"-y", "2005-2008,2018", samplefile,
	)
	cmd.Env = []string{"RUNME=1"}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	run(t, "diff", samplefile, sampleLicensed)
}

func TestWriteErrors(t *testing.T) {
	if os.Getenv("RUNME") != "" {
		main()
		return
	}

	tmp := tempDir(t)
	t.Logf("tmp dir: %s", tmp)
	samplefile := filepath.Join(tmp, "file.c")

	run(t, "cp", "testdata/initial/file.c", samplefile)
	run(t, "chmod", "0444", samplefile)
	cmd := exec.Command(os.Args[0],
		"-test.run=TestWriteErrors",
		"-l", "apache", "-c", "Bhojpur Consulting Private Limited, India.", "-y", "2018",
		samplefile,
	)
	cmd.Env = []string{"RUNME=1"}
	out, err := cmd.CombinedOutput()
	if err == nil {
		run(t, "chmod", "0644", samplefile)
		t.Fatalf("TestWriteErrors exited with a zero exit code.\n%s", out)
	}
	run(t, "chmod", "0644", samplefile)
}

func TestReadErrors(t *testing.T) {
	if os.Getenv("RUNME") != "" {
		main()
		return
	}

	tmp := tempDir(t)
	t.Logf("tmp dir: %s", tmp)
	samplefile := filepath.Join(tmp, "file.c")

	run(t, "cp", "testdata/initial/file.c", samplefile)
	run(t, "chmod", "a-r", samplefile)
	cmd := exec.Command(os.Args[0],
		"-test.run=TestReadErrors",
		"-l", "apache", "-c", "Bhojpur Consulting Private Limited, India.", "-y", "2018",
		samplefile,
	)
	cmd.Env = []string{"RUNME=1"}
	out, err := cmd.CombinedOutput()
	if err == nil {
		run(t, "chmod", "0644", samplefile)
		t.Fatalf("TestWriteErrors exited with a zero exit code.\n%s", out)
	}
	run(t, "chmod", "0644", samplefile)
}

func TestCheckSuccess(t *testing.T) {
	if os.Getenv("RUNME") != "" {
		main()
		return
	}

	tmp := tempDir(t)
	t.Logf("tmp dir: %s", tmp)
	samplefile := filepath.Join(tmp, "file.c")

	run(t, "cp", "testdata/expected/file.c", samplefile)
	cmd := exec.Command(os.Args[0],
		"-test.run=TestCheckSuccess",
		"-l", "apache", "-c", "Bhojpur Consulting Private Limited, India.", "-y", "2018",
		"-check", samplefile,
	)
	cmd.Env = []string{"RUNME=1"}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
}

func TestCheckFail(t *testing.T) {
	if os.Getenv("RUNME") != "" {
		main()
		return
	}

	tmp := tempDir(t)
	t.Logf("tmp dir: %s", tmp)
	samplefile := filepath.Join(tmp, "file.c")

	run(t, "cp", "testdata/initial/file.c", samplefile)
	cmd := exec.Command(os.Args[0],
		"-test.run=TestCheckFail",
		"-l", "apache", "-c", "Bhojpur Consulting Private Limited, India.", "-y", "2018",
		"-check", samplefile,
	)
	cmd.Env = []string{"RUNME=1"}
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("TestCheckFail exited with a zero exit code.\n%s", out)
	}
}

func TestMPL(t *testing.T) {
	if os.Getenv("RUNME") != "" {
		main()
		return
	}

	tmp := tempDir(t)
	t.Logf("tmp dir: %s", tmp)
	samplefile := filepath.Join(tmp, "file.c")

	run(t, "cp", "testdata/expected/file.c", samplefile)
	cmd := exec.Command(os.Args[0],
		"-test.run=TestMPL",
		"-l", "mpl", "-c", "Bhojpur Consulting Private Limited, India.", "-y", "2018",
		"-check", samplefile,
	)
	cmd.Env = []string{"RUNME=1"}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
}

func createTempFile(contents string, pattern string) (*os.File, error) {
	f, err := ioutil.TempFile("", pattern)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(f.Name(), []byte(contents), 0644); err != nil {
		return nil, err
	}

	return f, nil
}

func TestLicense(t *testing.T) {
	tmpl := template.Must(template.New("").Parse("{{.Holder}}{{.Year}}{{.SPDXID}}"))
	data := LicenseData{Holder: "H", Year: "Y", SPDXID: "S"}

	tests := []struct {
		contents     string
		wantContents string
		wantUpdated  bool
	}{
		{"", "// HYS\n\n", true},
		{"content", "// HYS\n\ncontent", true},

		// various headers that should be left intact. Many don't make
		// sense for our temp file extension, but that doesn't matter.
		{"#!/bin/bash\ncontent", "#!/bin/bash\n// HYS\n\ncontent", true},
		{"<?xml version='1.0'?>\ncontent", "<?xml version='1.0'?>\n// HYS\n\ncontent", true},
		{"<!doctype html>\ncontent", "<!doctype html>\n// HYS\n\ncontent", true},
		{"<!DOCTYPE HTML>\ncontent", "<!DOCTYPE HTML>\n// HYS\n\ncontent", true},
		{"# encoding: UTF-8\ncontent", "# encoding: UTF-8\n// HYS\n\ncontent", true},
		{"# frozen_string_literal: true\ncontent", "# frozen_string_literal: true\n// HYS\n\ncontent", true},
		{"<?php\ncontent", "<?php\n// HYS\n\ncontent", true},
		{"# escape: `\ncontent", "# escape: `\n// HYS\n\ncontent", true},
		{"# syntax: docker/dockerfile:1.3\ncontent", "# syntax: docker/dockerfile:1.3\n// HYS\n\ncontent", true},

		// ensure files with existing license or generated files are
		// skipped. No need to test all permutations of these, since
		// there are specific tests below.
		{"// Copyright 2018 Bhojpur Consulting\ncontent", "// Copyright 2018 Bhojpur Consulting\ncontent", false},
		{"// Code generated by Bhojpur License; DO NOT EDIT.\ncontent", "// Code generated by Bhojpur License; DO NOT EDIT.\ncontent", false},
	}

	for _, tt := range tests {
		// create temp file with contents
		f, err := createTempFile(tt.contents, "*.go")
		if err != nil {
			t.Error(err)
		}
		fi, err := f.Stat()
		if err != nil {
			t.Error(err)
		}

		// run Bhojpur License
		updated, err := BhojpurLicense(f.Name(), fi.Mode(), tmpl, data)
		if err != nil {
			t.Error(err)
		}

		// check results
		if updated != tt.wantUpdated {
			t.Errorf("Bhojpur License with contents %q returned updated: %t, want %t", tt.contents, updated, tt.wantUpdated)
		}
		gotContents, err := ioutil.ReadFile(f.Name())
		if err != nil {
			t.Error(err)
		}
		if got := string(gotContents); got != tt.wantContents {
			t.Errorf("Bhojpur License with contents %q returned contents: %q, want %q", tt.contents, got, tt.wantContents)
		}

		// if all tests passed, cleanup temp file
		if !t.Failed() {
			_ = os.Remove(f.Name())
		}
	}
}

// Test that Bhojpur License headers are added using the appropriate prefix for
// different filenames and extensions.
func TestLicenseHeader(t *testing.T) {
	tpl := template.Must(template.New("").Parse("{{.Holder}}{{.Year}}{{.SPDXID}}"))
	data := LicenseData{Holder: "H", Year: "Y", SPDXID: "S"}

	tests := []struct {
		paths []string // paths passed to licenseHeader
		want  string   // expected result of executing template
	}{
		{
			[]string{"f.unknown"},
			"",
		},
		{
			[]string{"f.c", "f.h", "f.gv", "f.java", "f.scala", "f.kt", "f.kts"},
			"/*\n * HYS\n */\n\n",
		},
		{
			[]string{"f.js", "f.mjs", "f.cjs", "f.jsx", "f.tsx", "f.css", "f.scss", "f.sass", "f.tf", "f.ts"},
			"/**\n * HYS\n */\n\n",
		},
		{
			[]string{"f.cc", "f.cpp", "f.cs", "f.go", "f.hcl", "f.hh", "f.hpp", "f.m", "f.mm", "f.proto",
				"f.rs", "f.swift", "f.dart", "f.groovy", "f.v", "f.sv", "f.php"},
			"// HYS\n\n",
		},
		{
			[]string{"f.py", "f.sh", "f.yaml", "f.yml", "f.dockerfile", "dockerfile", "f.rb", "gemfile", "f.tcl", "f.bzl", "f.pl"},
			"# HYS\n\n",
		},
		{
			[]string{"f.el", "f.lisp"},
			";; HYS\n\n",
		},
		{
			[]string{"f.erl"},
			"% HYS\n\n",
		},
		{
			[]string{"f.hs", "f.sql", "f.sdl"},
			"-- HYS\n\n",
		},
		{
			[]string{"f.html", "f.xml", "f.vue", "f.wxi", "f.wxl", "f.wxs"},
			"<!--\n HYS\n-->\n\n",
		},
		{
			[]string{"f.ml", "f.mli", "f.mll", "f.mly"},
			"(**\n   HYS\n*)\n\n",
		},
		{
			[]string{"cmakelists.txt", "f.cmake", "f.cmake.in"},
			"# HYS\n\n",
		},

		// ensure matches are case insenstive
		{
			[]string{"F.PY", "DoCkErFiLe"},
			"# HYS\n\n",
		},
	}

	for _, tt := range tests {
		for _, path := range tt.paths {
			header, _ := licenseHeader(path, tpl, data)
			if got := string(header); got != tt.want {
				t.Errorf("licenseHeader(%q) returned: %q, want: %q", path, got, tt.want)
			}
		}
	}
}

// Test that generated files are properly recognized.
func TestIsGenerated(t *testing.T) {
	tests := []struct {
		content string
		want    bool
	}{
		{"", false},
		{"Generated", false},
		{"// Code generated by Bhojpur License; DO NOT EDIT.", true},
		{"/*\n* Code generated by Bhojpur License; DO NOT EDIT.\n*/\n", true},
		{"DO NOT EDIT! Replaced on runs of cargo-raze", true},
	}

	for _, tt := range tests {
		b := []byte(tt.content)
		if got := isGenerated(b); got != tt.want {
			t.Errorf("isGenerated(%q) returned %v, want %v", tt.content, got, tt.want)
		}
	}
}

// Test that existing license headers are identified.
func TestHasLicense(t *testing.T) {
	tests := []struct {
		content string
		want    bool
	}{
		{"", false},
		{"This is my license", false},
		{"This code is released into the public domain.", false},
		{"SPDX: MIT", false},

		{"Copyright 2018", true},
		{"CoPyRiGhT 2018", true},
		{"Subject to the terms of the Mozilla Public License", true},
		{"SPDX-License-Identifier: MIT", true},
		{"spdx-license-identifier: MIT", true},
	}

	for _, tt := range tests {
		b := []byte(tt.content)
		if got := hasLicense(b); got != tt.want {
			t.Errorf("hasLicense(%q) returned %v, want %v", tt.content, got, tt.want)
		}
	}
}

func TestFileMatches(t *testing.T) {
	tests := []struct {
		pattern   string
		path      string
		wantMatch bool
	}{
		// basic single directory patterns
		{"", "file.c", false},
		{"*.c", "file.h", false},
		{"*.c", "file.c", true},

		// subdirectory patterns
		{"*.c", "vendor/file.c", false},
		{"**/*.c", "vendor/file.c", true},
		{"vendor/**", "vendor/file.c", true},
		{"vendor/**/*.c", "vendor/file.c", true},
		{"vendor/**/*.c", "vendor/a/b/file.c", true},

		// single character "?" match
		{"*.?", "file.c", true},
		{"*.?", "file.go", false},
		{"*.??", "file.c", false},
		{"*.??", "file.go", true},

		// character classes - sets and ranges
		{"*.[ch]", "file.c", true},
		{"*.[ch]", "file.h", true},
		{"*.[ch]", "file.ch", false},
		{"*.[a-z]", "file.c", true},
		{"*.[a-z]", "file.h", true},
		{"*.[a-z]", "file.go", false},
		{"*.[a-z]", "file.R", false},

		// character classes - negations
		{"*.[^ch]", "file.c", false},
		{"*.[^ch]", "file.h", false},
		{"*.[^ch]", "file.R", true},
		{"*.[!ch]", "file.c", false},
		{"*.[!ch]", "file.h", false},
		{"*.[!ch]", "file.R", true},

		// comma-separated alternative matches
		{"*.{c,go}", "file.c", true},
		{"*.{c,go}", "file.go", true},
		{"*.{c,go}", "file.h", false},

		// negating alternative matches
		{"*.[^{c,go}]", "file.c", false},
		{"*.[^{c,go}]", "file.go", false},
		{"*.[^{c,go}]", "file.h", true},
	}

	for _, tt := range tests {
		patterns := []string{tt.pattern}
		if got := fileMatches(tt.path, patterns); got != tt.wantMatch {
			t.Errorf("fileMatches(%q, %q) returned %v, want %v", tt.path, patterns, got, tt.wantMatch)
		}
	}
}
