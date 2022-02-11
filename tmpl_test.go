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
	"errors"
	"os"
	"testing"
	"text/template"
)

func init() {
	// ensure that pre-defined templates must parse
	template.Must(template.New("").Parse(tmplApache))
	template.Must(template.New("").Parse(tmplMIT))
	template.Must(template.New("").Parse(tmplBSD))
	template.Must(template.New("").Parse(tmplMPL))
}

func TestFetchTemplate(t *testing.T) {
	tests := []struct {
		description  string   // test case description
		license      string   // license passed to fetchTemplate
		templateFile string   // templatefile passed to fetchTemplate
		spdx         spdxFlag // spdx value passed to fetchTemplate
		wantTemplate string   // expected returned template
		wantErr      error    // expected returned error
	}{
		// custom template files
		{
			"non-existent template file",
			"",
			"/does/not/exist",
			spdxOff,
			"",
			os.ErrNotExist,
		},
		{
			"custom template file",
			"",
			"testdata/custom.tpl",
			spdxOff,
			"Copyright {{.Year}} {{.Holder}}\n\nCustom License Template\n",
			nil,
		},

		{
			"unknown license",
			"unknown",
			"",
			spdxOff,
			"",
			errors.New(`unknown license: "unknown". Include the '-s' flag to request SPDX style headers using this license`),
		},

		// pre-defined license templates, no SPDX
		{
			"apache license template",
			"Apache-2.0",
			"",
			spdxOff,
			tmplApache,
			nil,
		},
		{
			"mit license template",
			"MIT",
			"",
			spdxOff,
			tmplMIT,
			nil,
		},
		{
			"bsd license template",
			"bsd",
			"",
			spdxOff,
			tmplBSD,
			nil,
		},
		{
			"mpl license template",
			"MPL-2.0",
			"",
			spdxOff,
			tmplMPL,
			nil,
		},

		// SPDX variants
		{
			"apache license template with SPDX added",
			"Apache-2.0",
			"",
			spdxOn,
			tmplApache + spdxSuffix,
			nil,
		},
		{
			"apache license template with SPDX only",
			"Apache-2.0",
			"",
			spdxOnly,
			tmplSPDX,
			nil,
		},
		{
			"unknown license with SPDX only",
			"unknown",
			"",
			spdxOnly,
			tmplSPDX,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tpl, err := fetchTemplate(tt.license, tt.templateFile, tt.spdx)
			if tt.wantErr != nil && (err == nil || (!errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error())) {
				t.Fatalf("fetchTemplate(%q, %q) returned error: %#v, want %#v", tt.license, tt.templateFile, err, tt.wantErr)
			}
			if tpl != tt.wantTemplate {
				t.Errorf("fetchTemplate(%q, %q) returned template: %q, want %q", tt.license, tt.templateFile, tpl, tt.wantTemplate)
			}
		})
	}
}

func TestExecuteTemplate(t *testing.T) {
	tests := []struct {
		template      string
		data          LicenseData
		top, mid, bot string
		want          string
	}{
		{
			"",
			LicenseData{},
			"", "", "",
			"\n",
		},
		{
			"{{.Holder}}{{.Year}}{{.SPDXID}}",
			LicenseData{Holder: "H", Year: "Y", SPDXID: "S"},
			"", "", "",
			"HYS\n\n",
		},
		{
			"{{.Holder}}{{.Year}}{{.SPDXID}}",
			LicenseData{Holder: "H", Year: "Y", SPDXID: "S"},
			"", "// ", "",
			"// HYS\n\n",
		},
		{
			"{{.Holder}}{{.Year}}{{.SPDXID}}",
			LicenseData{Holder: "H", Year: "Y", SPDXID: "S"},
			"/*", " * ", "*/",
			"/*\n * HYS\n*/\n\n",
		},

		// ensure we don't escape HTML characters by using the wrong template package
		{
			"{{.Holder}}",
			LicenseData{Holder: "A&Z"},
			"", "", "",
			"A&Z\n\n",
		},
	}

	for _, tt := range tests {
		tpl, err := template.New("").Parse(tt.template)
		if err != nil {
			t.Errorf("error parsing template: %v", err)
		}
		got, err := ExecuteTemplate(tpl, tt.data, tt.top, tt.mid, tt.bot)
		if err != nil {
			t.Errorf("ExecuteTemplate(%q, %v, %q, %q, %q) returned error: %v", tt.template, tt.data, tt.top, tt.mid, tt.bot, err)
		}
		if string(got) != tt.want {
			t.Errorf("ExecuteTemplate(%q, %v, %q, %q, %q) returned %q, want: %q", tt.template, tt.data, tt.top, tt.mid, tt.bot, string(got), tt.want)
		}
	}
}
