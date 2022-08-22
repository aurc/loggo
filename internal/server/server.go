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

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aurc/loggo/internal/config"
	"github.com/aurc/loggo/internal/reader"
	"github.com/aurc/loggo/internal/server/handler"
	"github.com/gorilla/mux"
)

type server struct {
	settings *Settings
	inSlice  []map[string]any
}

type Settings struct {
	Reader reader.Reader
	Config *config.Config
	Port   int64
}

func Run(settings *Settings) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	httpAddr := fmt.Sprintf(`:%d`, settings.Port)

	ls := &server{
		settings: settings,
		inSlice:  make([]map[string]any, 0),
	}
	ls.buffer()

	r := mux.NewRouter()
	r.HandleFunc("/", handler.RootHandler)
	r.HandleFunc("/stream", ls.StreamHandler)
	r.HandleFunc("/template", handler.TemplateHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    httpAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()

}

func (s *server) buffer() {
	go func() {
		l := s.settings.Reader
		err := l.StreamInto()
		if err != nil {
			panic(err)
		}
		for {
			t := <-l.ChanReader()
			if len(t) > 0 {
				m := make(map[string]interface{})
				err := json.Unmarshal([]byte(t), &m)
				if err != nil {
					m[config.ParseErr] = err.Error()
					m[config.TextPayload] = t
				}
				s.inSlice = append(s.inSlice, m)
				//fmt.Println(m)
			}
		}
	}()
}
