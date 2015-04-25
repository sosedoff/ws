package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: ws URL")
		os.Exit(1)
	}

	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ws, _, err := dialer.Dial(os.Args[1], nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	bio := bufio.NewReader(os.Stdin)

	go func() {
		for {
			line, _, err := bio.ReadLine()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			ws.WriteMessage(websocket.TextMessage, []byte(line))
		}
	}()

	for {
		mtype, mdata, err := ws.ReadMessage()
		ts := time.Now()

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		var prettyJSON bytes.Buffer

		err = json.Indent(&prettyJSON, mdata, "", "  ")
		if err == nil {
			fmt.Printf("%s\n%s\n", ts.Format("20060102-15:04:05.000"), prettyJSON.Bytes())
		} else {
			fmt.Printf("%s\n%s\n", ts.Format("20060102-15:04:05.000"), mtype, mdata)
		}

		fmt.Println("-----------------------------------------------------------")
	}
}
