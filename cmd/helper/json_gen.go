package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

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
		JsonGenerator(os.Stdout)
		time.Sleep(time.Second)
	}
}

func JsonGenerator(writer io.Writer) {
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
		fmt.Fprintln(writer, string(b))
		time.Sleep(time.Second)
	}
}
