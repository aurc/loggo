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

package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func JsonGenerator(writer io.Writer) {
	b, err := ioutil.ReadFile("internal/testdata/test1.json")
	if err != nil {
		panic(err)
	}
	jm := make(map[string]interface{})
	_ = json.Unmarshal(b, &jm)
	i := 0
	for {
		//if i != 0 && i%(rand.Intn(27)+1) == 0 {
		//	_, _ = fmt.Fprintln(writer, "bad json")
		//	i++
		//	continue
		//}
		i++
		uid1 := uuid.New().String()
		uid2 := uuid.New().String()
		id3 := rand.Intn(30)
		jm["severity"] = []string{"INFO", "ERROR", "DEBUG", "WARN"}[rand.Intn(4)]
		jm["insertId"] = uid1
		jm["trace"] = uid2
		jm["spanId"] = fmt.Sprintf("%d", id3)
		jm["timestamp"] = time.Now().Format("2006-01-02T15:04:05-0700")
		b, _ = json.Marshal(jm)
		_, _ = fmt.Fprintln(writer, string(b))
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(800)))
		if i%(rand.Intn(5)+1) == 0 {
			time.Sleep(time.Second * time.Duration(rand.Intn(5)))
		}

	}
}
