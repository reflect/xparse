// Copyright (c) 2014 Dataence, LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// xtime is a time parser that parses the time without knowning the exact format.
package xtime

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrUnknownTimeFormat = errors.New("Unknown time format")
)

// TimeFormats is a list of commonly seen time formats from log messages
var TimeFormats []string = []string{
	"Mon Jan _2 15:04:05 2006",
	"Mon Jan _2 15:04:05 MST 2006",
	"Mon Jan 02 15:04:05 -0700 2006",
	"02 Jan 06 15:04 MST",
	"02 Jan 06 15:04 -0700",
	"Monday, 02-Jan-06 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"2006-01-02T00:00:00",
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05.999999999Z07:00",
	"2006-01-02T15:04:05.999999999Z",
	"2006-01-02T15:04:05.999Z07:00",
	"2006-01-02T15:04:05.999Z",
	"Jan _2 15:04:05",
	"Jan _2 15:04:05.000",
	"Jan _2 15:04:05.000000",
	"Jan _2 15:04:05.000000000",
	"_2/Jan/2006:15:04:05 -0700",
	"Jan 2, 2006 3:04:05 PM",
	"Jan 2 2006 15:04:05",
	"Jan 2 15:04:05 2006",
	"Jan 2 15:04:05 -0700",
	"2006-01-02 15:04:05,000 -0700",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04:05-0700",
	"2006-01-02 15:04:05,000",
	"2006-01-02 15:04:05 MST",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"2006/01/02",
	"2006-01-02",
	"01/02/2006",
	"01/02/2006 15:04:05",
	"06-01-02 15:04:05,000 -0700",
	"06-01-02 15:04:05,000",
	"06-01-02 15:04:05",
	"06/01/02 15:04:05",
	"15:04:05,000",
	"1/2/2006 3:04:05 PM",
	"1/2/06 3:04:05.000 PM",
	"1/2/2006 15:04",
	"1/2/2006",
	"2006/1/2",
}

type TimeTree struct {
	formats []string
	root    *timeNode
}

func (tt *TimeTree) Parse(t string) (time.Time, error) {
	tx := strings.ToLower(t)
	cur := tt.root

	for i, r := range tx {
		typ := tnType(r)

		for _, n := range cur.children {
			if (n.ntype == timeNodeDigitOrSpace && (typ == timeNodeDigit || typ == timeNodeSpace)) ||
				(n.ntype == typ && (typ != timeNodeLiteral || (typ == timeNodeLiteral && rune(n.value) == r))) {

				cur = n
				break
			}
		}

		if cur.final && i == len(tx)-1 {
			return time.Parse(tt.formats[cur.subtype], t)
		}
	}

	return time.Time{}, ErrUnknownTimeFormat
}

func (tt *TimeTree) IsTime(t string) bool {
	_, err := tt.Parse(t)
	return err == nil
}

func Root() *TimeTree {
	return timeTreeRoot
}

func Compile(formats []string) *TimeTree {
	return &TimeTree{
		formats: formats,
		root:    buildTimeTree(formats),
	}
}

func Parse(t string) (time.Time, error) {
	return timeTreeRoot.Parse(t)
}

func IsTime(t string) bool {
	return timeTreeRoot.IsTime(t)
}

type timeNodeType int

type timeNode struct {
	ntype    timeNodeType
	value    rune
	final    bool
	subtype  int
	children []*timeNode
}

const (
	timeNodeRoot timeNodeType = iota
	timeNodeDigit
	timeNodeLetter
	timeNodeLiteral
	timeNodeSpace
	timeNodeDigitOrSpace
	timeNodePlusOrMinus
)

var (
	timeTreeRoot *TimeTree
)

func init() {
	timeTreeRoot = Compile(TimeFormats)
}

func buildTimeTree(formats []string) *timeNode {
	root := &timeNode{ntype: timeNodeRoot}

	for i, f := range formats {
		tf := strings.ToLower(f)
		parent := root

		for _, r := range tf {
			typ := tnType(r)

			hasChild := false
			var child *timeNode

			for _, child = range parent.children {
				if (child.ntype == typ && (typ != timeNodeLiteral || (typ == timeNodeLiteral && child.value == r))) ||
					(child.ntype == timeNodeDigitOrSpace && (typ == timeNodeDigit || typ == timeNodeSpace)) {
					hasChild = true
					break
				} else if child.ntype == timeNodeDigit && typ == timeNodeDigitOrSpace {
					child.ntype = timeNodeDigitOrSpace
					hasChild = true
					break
				}
			}

			if !hasChild {
				child = &timeNode{ntype: typ, value: r}
				parent.children = append(parent.children, child)
			}

			parent = child
		}

		parent.final = true
		parent.subtype = i
	}

	return root
}

func tnType(r rune) timeNodeType {
	switch {
	case r >= '0' && r <= '9':
		return timeNodeDigit
	case r >= 'a' && r <= 'y':
		return timeNodeLetter
	case r == ' ':
		return timeNodeSpace
	case r == '_':
		return timeNodeDigitOrSpace
	case r == '+' || r == '-':
		return timeNodePlusOrMinus
	case r >= 'A' && r <= 'Y':
		return timeNodeLetter
	case r == 'z' || r == 'Z':
		return timeNodePlusOrMinus
	}

	return timeNodeLiteral
}
