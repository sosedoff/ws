package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/gorilla/websocket"
)

var (
	dialer = websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func userInput(ws *websocket.Conn) {
	bio := bufio.NewReader(os.Stdin)
	line, _, err := bio.ReadLine()

	if err != nil && err != io.EOF {
		fmt.Println("Read error:", err)
		return
	}

	ws.WriteMessage(websocket.TextMessage, []byte(line))
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "" {
		fmt.Println("Usage: ws URL")
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Connecting...")
	ws, _, err := dialer.Dial(flag.Arg(0), nil)
	if err != nil {
		fmt.Println(os.Stderr, "Error:", err)
		os.Exit(2)
	}

	fmt.Fprintln(os.Stderr, "Connected. Listening for messages...")
	userInput(ws)

	for {
		_, mdata, err := ws.ReadMessage()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			break
		}

		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, mdata, "", "  ")

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			break
		}

		fmt.Printf("%s\n", prettyJSON.Bytes())
	}
}
