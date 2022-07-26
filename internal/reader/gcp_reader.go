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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	logging "cloud.google.com/go/logging/apiv2"
	"github.com/rivo/tview"
	"google.golang.org/api/iterator"
	loggingpb "google.golang.org/genproto/googleapis/logging/v2"
)

type gcpStream struct {
	reader
	projectID string
	filter    string
	freshness string
	isTail    bool
	stop      bool
}

var scopes = []string{
	"https://www.googleapis.com/auth/logging.read",
	"https://www.googleapis.com/auth/cloud-platform.read-only",
}

func MakeGCPReader(project, filter, freshness string) *gcpStream {
	return &gcpStream{
		projectID: project,
		filter:    filter,
		freshness: freshness,
		isTail:    freshness == "tail",
	}
}

func (s *gcpStream) StreamInto(strChan chan<- string) error {
	ctx := context.Background()
	err := s.checkAuth(ctx)
	if err != nil {
		log.Fatal(err)
	}

	c, err := logging.NewClient(ctx)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			r := recover()
			log.Fatal(r)
		}()
		defer c.Close()
		if s.isTail {
			s.streamTail(ctx, c, strChan)
		} else {
			s.streamFrom(ctx, c, strChan)
			//fallback to tail if from returns
			s.streamTail(ctx, c, strChan)
		}
	}()
	return nil
}

func (s *gcpStream) streamFrom(ctx context.Context, c *logging.Client, strChan chan<- string) {
	lastTime := s.freshness
	lastFilter := ""
	for !s.stop {
		//if !lastTS.IsZero()
		filter := fmt.Sprintf(`timestamp > "%s"`, lastTime)
		if filter == lastFilter {
			// wait
			_ = lastFilter
			//t, _ := time.Parse(time.RFC3339, lastTime)
			//filter = fmt.Sprintf(`timestamp > "%s"`, t.Add(time.Second).Format(time.RFC3339))
			return
		}
		lastFilter = filter
		if len(s.filter) > 0 {
			filter = fmt.Sprintf(`timestamp > "%s" AND (%s)`, lastTime, s.filter)
		}

		it := c.ListLogEntries(ctx, &loggingpb.ListLogEntriesRequest{
			ResourceNames: []string{"projects/" + s.projectID},
			Filter:        filter,
			PageSize:      100,
		})
		for {
			resp, err := it.Next()
			if iterator.Done == err {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			var b []byte
			b, lastTime = massageEntryLog(resp)
			strChan <- string(b)
		}
	}
}

func (s *gcpStream) streamTail(ctx context.Context, c *logging.Client, strChan chan<- string) {
	stream, err := c.TailLogEntries(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.CloseSend()

	req := &loggingpb.TailLogEntriesRequest{
		ResourceNames: []string{"projects/" + s.projectID},
		Filter:        s.filter,
	}
	if err := stream.Send(req); err != nil {
		log.Fatal(err)
	}

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, resp := range chunk.Entries {
			if iterator.Done == err {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			var b []byte
			b, _ = massageEntryLog(resp)
			strChan <- string(b)
		}
	}
}

func massageEntryLog(resp *loggingpb.LogEntry) ([]byte, string) {
	lastTime := resp.GetTimestamp().AsTime().Local().Format(time.RFC3339)
	severity := resp.GetSeverity().String()
	b, _ := json.Marshal(resp)
	m := make(map[string]interface{})
	_ = json.Unmarshal(b, &m)
	m["severity"] = severity
	m["timestamp"] = lastTime
	if resp.GetJsonPayload() != nil {
		m["jsonPayload"] = m["Payload"].(map[string]interface{})["JsonPayload"]
	} else if len(resp.GetTextPayload()) > 0 {
		m["textPayload"] = m["Payload"].(map[string]interface{})["TextPayload"]
	}
	delete(m, "Payload")
	delete(m, "receive_timestamp")
	b, _ = json.Marshal(m)
	return b, lastTime
}

func (s *gcpStream) Close() {
	s.stop = true
	close(s.strChan)
}

func (s *gcpStream) checkAuth(ctx context.Context) error {
	c, err := logging.NewClient(ctx)
	if err == nil {
		it := c.ListLogs(ctx, &loggingpb.ListLogsRequest{
			ResourceNames: []string{"projects/" + s.projectID},
			PageSize:      1,
		})
		_, err = it.Next()
	}
	if err != nil {
		app := tview.NewApplication()
		modal := tview.NewModal().
			SetText("Authenticating with gcloud... \nRedirecting to your browser.")
		go func() {
			cmd := exec.Command("gcloud", "auth", "application-default", "login")
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
			app.Stop()
		}()
		if err := app.SetRoot(modal, false).EnableMouse(true).Run(); err != nil {
			panic(err)
		}
	}

	return nil
}

func ParseFrom(str string) string {
	str = strings.TrimSpace(str)
	if str == "tail" {
		return "tail"
	}
	regF := regexp.MustCompile(`^\d+(s|m|d|h)$`)
	regD := regexp.MustCompile(`^\d{4}(-\d{2}){2}T(\d{2}:){2}\d{2}$`)
	if regF.Match([]byte(str)) {
		numb, _ := strconv.ParseInt(str[0:len(str)-1], 10, 64)
		unit := str[len(str)-1:]
		var duration time.Duration
		switch unit {
		case "s":
			duration = time.Second * time.Duration(numb)
		case "m":
			duration = time.Minute * time.Duration(numb)
		case "h":
			duration = time.Hour * time.Duration(numb)
		case "d":
			duration = time.Hour * time.Duration(numb) * 24
		default:
		}
		return time.Now().Add(-1 * time.Duration(duration)).Format(time.RFC3339)
	} else if regD.Match([]byte(str)) {
		t, err := time.Parse(`2006-01-02T15:04:05`, str)
		if err != nil {
			log.Fatal("Invalid parameter for 'from' flag - bad format: ", err)
		}
		return t.Format(time.RFC3339)
	} else {
		log.Fatal("Invalid parameter for 'from' flag.")
	}
	return ""
}
