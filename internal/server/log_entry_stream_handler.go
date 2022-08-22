package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type sockSession struct {
	id        string
	startFrom chan int
	server    *server
}

func (s *server) StreamHandler(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	ss := &sockSession{
		id:        uuid.New().String(),
		startFrom: make(chan int, 1),
		server:    s,
	}

	go ss.reader(ws)
	ss.writer(ws)
}

func (ss *sockSession) writer(ws *websocket.Conn) {
	defer func() {
		ws.Close()
	}()
	curr := <-ss.startFrom
	s := ss.server
	for {
		if len(s.inSlice) > curr {
			payload, err := json.Marshal(s.inSlice[curr])
			if err != nil {
				fmt.Println("failed to marshal log entry ", err)
			}
			resp, err := json.Marshal(&LogEntryStreamResponse{
				Position: curr,
				Payload:  string(payload),
			})
			if err != nil {
				fmt.Println("failed to marshal response ", err)
			}
			w, err := ws.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Println("failed to acquire writer ", err)
			}
			if _, err := w.Write(resp); err != nil {
				fmt.Println("fail to dispatch message ", err)
			}
			curr++
		} else {
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func (ss *sockSession) reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.NextReader()
		if err != nil {
			log.Println(err)
			return
		}
		switch messageType {
		case websocket.TextMessage:
			b, err := ioutil.ReadAll(p)
			if err != nil {
				log.Println(err)
				return
			}
			var req LogEntryStreamRequest
			if err := json.Unmarshal(b, &req); err != nil {
				log.Println(err)
				return
			}
			ss.startFrom <- req.PositionFrom
			log.Println("Received position ", req.PositionFrom)
		}
	}
}
