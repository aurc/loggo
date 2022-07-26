/*
Copyright © 2022 Aurelio Calegari, et al.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package reader

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseFrom(t *testing.T) {
	tests := []struct {
		name       string
		givenValue string
		wantsValue string
	}{
		{
			name:       "Test tail",
			givenValue: "tail",
			wantsValue: "tail",
		},
		{
			name:       "Test Relative Second",
			givenValue: "86400s",
			wantsValue: time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		},
		{
			name:       "Test Relative Minute",
			givenValue: "1440m",
			wantsValue: time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		},
		{
			name:       "Test Relative Hour",
			givenValue: "24h",
			wantsValue: time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		},
		{
			name:       "Test Relative Day",
			givenValue: "1d",
			wantsValue: time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		},
		{
			name:       "Test Fixed Time",
			givenValue: "2021-01-30T15:00:00",
			wantsValue: func() string {
				tv, _ := time.Parse("2006-01-02T15:04:05", "2021-01-30T15:00:00")
				return tv.Format(time.RFC3339)
			}(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := ParseFrom(test.givenValue)
			fmt.Println(v)
			assert.Equal(t, test.wantsValue, v)
		})
	}
}
