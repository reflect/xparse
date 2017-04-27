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

package xtime

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	re1 = regexp.MustCompile("_")
	re2 = regexp.MustCompile("Z(.)") // Go doesn't support lookahead assertions.
)

func TestTimeFormats(t *testing.T) {
	for _, f := range TimeFormats {
		tx := re2.ReplaceAllString(re1.ReplaceAllString(f, " "), "+$1")
		expected, err := time.Parse(f, tx)
		require.NoError(t, err)
		actual, err := Parse(tx)
		require.NoError(t, err)
		require.Equal(t, expected.UnixNano(), actual.UnixNano())
	}
}

func TestTimeFormatsIsTime(t *testing.T) {
	for _, f := range TimeFormats {
		tx := re2.ReplaceAllString(re1.ReplaceAllString(f, " "), "+$1")
		require.True(t, IsTime(tx), "for date format %s and input %s", f, tx)
	}
}

func TestKnownTimes(t *testing.T) {
	times := []struct {
		In, Out string
	}{
		{"2017-03-01T23:16:37.986Z", "2017-03-01T23:16:37Z"},
		{"2017-03-01T23:16:37.986-07:00", "2017-03-01T23:16:37-07:00"},
		{"Tuesday, 03-Jan-06 15:04:05 -07:00", "2006-01-03T15:04:05-07:00"},
	}

	for _, tx := range times {
		parsed, err := Parse(tx.In)
		require.NoError(t, err, "for input %s", tx.In)
		require.Equal(t, tx.Out, parsed.Format(time.RFC3339), "for input %s", tx.In)
	}
}

func ExampleParse() {
	t1, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")
	t2, err := Parse("2006-01-02T15:04:05+07:00")
	if err != nil {
		fmt.Println(err)
	} else if t1.UnixNano() != t2.UnixNano() {
		fmt.Printf("%d != %d\n", t1.UnixNano(), t2.UnixNano())
	} else {
		fmt.Println(t2)
	}
}
