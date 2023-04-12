package main

import (
	"encoding/json"
	"log"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type topologyMsg struct {
	Topology map[string][]string `json:"topology"`
}

func main() {
	n := maelstrom.NewNode()
	topology := make(map[string][]string)
	messages := make(map[int]struct{})

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		index := int(body["message"].(float64))
		_, ok := messages[index]
		if ok {
			return n.Reply(msg, map[string]any{
				"type": "broadcast_ok",
			})
		}

		messages[index] = struct{}{}
		for key := range topology {
			if n.ID() != key {
				if err := n.Send(key, body); err != nil {
					log.Printf("ERROR: %s", err)
					return err
				}
			}
		}

		return n.Reply(msg, map[string]any{
			"type": "broadcast_ok",
		})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		var message_list []int
		for key := range messages {
			message_list = append(message_list, key)
		}

		return n.Reply(msg, map[string]any{
			"type":     "read_ok",
			"messages": message_list,
		})
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body topologyMsg
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		topology = body.Topology

		return n.Reply(msg, map[string]any{
			"type": "topology_ok",
		})
	})

	// Execute the node's message loop. This will run until STDIN is closed.
	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
