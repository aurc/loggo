/*
Copyright © 2022 Aurelio Calegari

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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/aurc/loggo/internal/uitest/helper"
	"github.com/google/uuid"
)

const datelayout = "2006-01-02T15:04:05-0700"

func main() {
	b, err := ioutil.ReadFile("testdata/test1.json")
	if err != nil {
		panic(err)
	}
	jm := make(map[string]interface{})
	_ = json.Unmarshal(b, &jm)

	for {
		uid1 := uuid.New().String()
		uid2 := uuid.New().String()
		id3 := rand.Intn(30)
		jm["insertId"] = uid1
		jm["trace"] = uid2
		jm["spanId"] = fmt.Sprintf("%d", id3)
		jm["timestamp"] = time.Now().Format(datelayout)
		b, _ = json.Marshal(jm)
		helper.JsonGenerator(os.Stdout)
		time.Sleep(time.Second)
	}
}
