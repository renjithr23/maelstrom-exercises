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
	var messages []int

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		messages = append(messages, int(body["message"].(float64)))
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

		return n.Reply(msg, map[string]any{
			"type":     "read_ok",
			"messages": messages,
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
