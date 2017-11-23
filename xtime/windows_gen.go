// +build ignore

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/xml"
	"flag"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"

	"github.com/reflect/xflag"
)

var (
	Source     *url.URL
	OutputPath string
)

type MapZoneTypes []string

func (mzt *MapZoneTypes) UnmarshalXMLAttr(attr xml.Attr) error {
	*mzt = strings.Fields(attr.Value)
	return nil
}

type MapZone struct {
	Other     string       `xml:"other,attr"`
	Territory string       `xml:"territory,attr"`
	Types     MapZoneTypes `xml:"type,attr"`
}

type MapTimezones struct {
	OtherVersion string    `xml:"otherVersion,attr"`
	TypeVersion  string    `xml:"typeVersion,attr"`
	MapZones     []MapZone `xml:"mapZone"`
}

type WindowsZones struct {
	MapTimezones MapTimezones `xml:"windowsZones>mapTimezones"`
}

func init() {
	xflag.URLVar(&Source, "source", "https://unicode.org/repos/cldr/trunk/common/supplemental/windowsZones.xml", "the URL to load source CLDR data from", xflag.ValidateURLSchemes("http", "https"))
	flag.StringVar(&OutputPath, "output-path", "-", "the path to write output to")
}

var program = `
// Generated by windows_gen.go. Do not modify; instead, edit the template in
// windows_gen.go.

// Copyright © 1991-2017 Unicode, Inc. All rights reserved.
// Distributed under the Terms of Use in http://www.unicode.org/copyright.html.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of the Unicode data files and any associated documentation
// (the "Data Files") or Unicode software and any associated documentation
// (the "Software") to deal in the Data Files or Software
// without restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, and/or sell copies of
// the Data Files or Software, and to permit persons to whom the Data Files
// or Software are furnished to do so, provided that either
// (a) this copyright and permission notice appear with all copies
// of the Data Files or Software, or
// (b) this copyright and permission notice appear in associated
// Documentation.
//
// THE DATA FILES AND SOFTWARE ARE PROVIDED "AS IS", WITHOUT WARRANTY OF
// ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT OF THIRD PARTY RIGHTS.
// IN NO EVENT SHALL THE COPYRIGHT HOLDER OR HOLDERS INCLUDED IN THIS
// NOTICE BE LIABLE FOR ANY CLAIM, OR ANY SPECIAL INDIRECT OR CONSEQUENTIAL
// DAMAGES, OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE,
// DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
// TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
// PERFORMANCE OF THE DATA FILES OR SOFTWARE.
//
// Except as contained in this notice, the name of a copyright holder
// shall not be used in advertising or otherwise to promote the sale,
// use or other dealings in these Data Files or Software without prior
// written authorization of the copyright holder.

// Source:                   {{.URL}}
// Source hash (SHA-256):    {{.URLHash}}
// tzdata version:           {{.IANAVersion}}
// Windows database version: {{.WindowsVersion}}

package xtime

var tzdataWindowsMapping = map[string]string{
{{- range .Zones -}}
	{{- if ne .Territory "001" -}}
		{{- $zone := . -}}
		{{- range .Types}}
			"{{.}}": "{{$zone.Other}}",
		{{- end -}}
	{{- end -}}
{{- end}}
}
`

var compiledProgram = template.Must(template.New("program").Parse(program))

func main() {
	flag.Parse()

	res, err := http.Get(Source.String())
	if err != nil {
		log.Fatalf("Could not download source: %+v", err)
	}
	defer res.Body.Close()

	// Write to our hashing function as well.
	h := sha256.New()
	r := io.TeeReader(res.Body, h)

	var zones WindowsZones
	if err := xml.NewDecoder(r).Decode(&zones); err != nil {
		log.Fatalf("Could not parse source as XML: %+v", err)
	}

	data := struct {
		URL            *url.URL
		URLHash        string
		IANAVersion    string
		WindowsVersion string
		Zones          []MapZone
	}{
		URL:            Source,
		URLHash:        fmt.Sprintf("%x", h.Sum(nil)),
		IANAVersion:    zones.MapTimezones.TypeVersion,
		WindowsVersion: zones.MapTimezones.OtherVersion,
		Zones:          zones.MapTimezones.MapZones,
	}

	var buf bytes.Buffer
	if err := compiledProgram.Execute(&buf, data); err != nil {
		log.Fatalf("Could not generate program: %+v", err)
	}

	output, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("Could not validate program code: %+v", err)
	}

	if len(OutputPath) > 0 && OutputPath != "-" {
		err = ioutil.WriteFile(OutputPath, output, 0644)
	} else {
		_, err = os.Stdout.Write(output)
	}

	if err != nil {
		log.Fatalf("Could not write program: %+v", err)
	}
}
