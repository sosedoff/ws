package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var opts struct {
	raw bool
}

func init() {
	flag.BoolVar(&opts.raw, "raw", false, "Print raw data")
	flag.Parse()
}

func main() {
	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	if flag.Arg(0) == "" {
		fmt.Println("Usage: ws URL")
		os.Exit(1)
	}

	ws, _, err := dialer.Dial(flag.Arg(0), nil)
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

		if opts.raw {
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				continue
			}
			fmt.Printf("%s\n", prettyJSON.Bytes())
		} else {
			if err == nil {
				fmt.Printf("%s\n%s\n", ts.Format("20060102-15:04:05.000"), prettyJSON.Bytes())
			} else {
				fmt.Printf("%s\n%s\n", ts.Format("20060102-15:04:05.000"), mtype, mdata)
			}

			fmt.Println("-----------------------------------------------------------")
		}
	}
}
