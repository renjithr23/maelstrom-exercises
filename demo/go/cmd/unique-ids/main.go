package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	// Register a handler for the "generate" message that responds with an "generate_ok" with a unique id.
	n.Handle("generate", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type.
		body["type"] = "generate_ok"
		newUUID, err := exec.Command("uuidgen").Output()
		if err != nil {
			log.Printf("ERROR: %s", err)
			os.Exit(1)
		}
		body["id"] = newUUID

		// Reply to the original message with the updated message type.
		return n.Reply(msg, body)
	})

	// Execute the node's message loop. This will run until STDIN is closed.
	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
