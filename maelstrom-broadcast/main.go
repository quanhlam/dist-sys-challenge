package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"
	"sync"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type MsgResp struct {
	Type	string	`json:"type"`
	MsgId 	int		`json:"msg_id"`
	InReplyTo int	`json:"in_reply_to"`
}


type MsgRespRead struct {
	Type	string	`json:"type"`
	MsgId 	int		`json:"msg_id"`
	InReplyTo int	`json:"in_reply_to"`
	Messages	[]float64	`json:"messages"`
}

type MsgReqBroadcast struct {
	Message	float64	`json:"message"`
	MsgId 	int		`json:"msg_id"`
	InReplyTo int	`json:"in_reply_to"`
}

type MsgReqTopology struct {
	MsgId 	int		`json:"msg_id"`
	InReplyTo int	`json:"in_reply_to"`
	Topology map[string]interface{} `json:"topology"`
}

type MsgReqRead struct {
	MsgId 	int		`json:"msg_id"`
	InReplyTo int	`json:"in_reply_to"`
}

type server struct {
	n	*maelstrom.Node

	messages	map[float64]struct{}
	mLock sync.RWMutex

	topology	map[string]interface{}
}



func contains(s map[float64]struct{}, e float64) bool {
    for key, _ := range s {
        if key == e {
            return true
        }
    }
    return false
}

func (s *server) broadcastAndWait(dest string, body MsgReqBroadcast) {
	maxRetries := 20
	for i := 0; i < maxRetries; i++ {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Duration(2*(i+1))*time.Second)
		defer cancel()

		_, err := s.n.SyncRPC(ctx, dest, body)
		if err == nil {
			return
		}
	}
}


func (s *server) broadcastHandler(msg maelstrom.Message) error {
	var body MsgReqBroadcast
	reply := new(MsgResp)

	if err:=json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	if err:=json.Unmarshal(msg.Body, &reply); err != nil {
		return err
	}
	
	message:=body.Message
	
	s.mLock.Lock()
	isExisted := contains(s.messages, message)
	s.mLock.Unlock()

	if (isExisted) {
		return s.n.Reply(msg, reply)
	}

	s.mLock.Lock()
	s.messages[body.Message] = struct{}{}
	s.mLock.Unlock()

	var neighborhoods []interface{} = s.topology[s.n.ID()].([]interface{})
	for _, element := range neighborhoods {
		s.broadcastAndWait(element.(string), body)
	}

	reply.Type = "broadcast_ok"

	return s.n.Reply(msg, reply)
}

func (s *server) readHandler(msg maelstrom.Message) error {
	body := new(MsgReqRead)
	reply := new(MsgRespRead)

	if err:=json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	var responseMessages []float64
	reply.Type = "read_ok"
	s.mLock.Lock()
	for key := range s.messages {
		responseMessages = append(responseMessages, key)
	}
	reply.Messages = responseMessages
	s.mLock.Unlock()

	return s.n.Reply(msg, reply)
}

func (s *server) topologyHandler(msg maelstrom.Message) error {

	body := new(MsgReqTopology)
	reply := new(MsgResp)

	if err:=json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}
	
	if err:=json.Unmarshal(msg.Body, &reply); err != nil {
		return err
	}

	s.topology = body.Topology
	reply.Type = "topology_ok"

	return s.n.Reply(msg, reply)
}


func main() {

 n:=maelstrom.NewNode()
 s :=&server{n:n}
 
 n.Handle("broadcast", s.broadcastHandler)
 n.Handle("read", s.readHandler)
 n.Handle("topology", s.topologyHandler)
 
 if err := n.Run(); err != nil {
	log.Printf("ERROR: %s", err)
	os.Exit(1)
 }
}
