package main

import (
	"os"
	"log"
	"time"
	"encoding/json"
	"strings"
	"strconv"
	"github.com/google/uuid"
  maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)


func makeTimestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}

func main() {
 n:=maelstrom.NewNode()
 n.Handle("generate", func(msg maelstrom.Message) error {
	var body map[string]any

	if err:=json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	body["type"] = "generate_ok"
	id := uuid.New()
	timestamp:= makeTimestamp()

	var randValue  strings.Builder

	randValue.WriteString(id.String())
	randValue.WriteString(strconv.FormatInt(timestamp, 10))

	body["id"] = id


	return n.Reply(msg, body)
 })

 if err := n.Run(); err != nil {
	log.Printf("ERROR: %s", err)
	os.Exit(1)
 }
}
