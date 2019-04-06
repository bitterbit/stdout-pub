// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"
	"net/url"
        "bufio"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func pipeStdin(ws *websocket.Conn, done chan int){
        stdin := bufio.NewReaderSize(os.Stdin, 1024)

        defer ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

        for  {
            line, err := stdin.ReadBytes('\n')
            if err != nil {
                log.Printf("err: %v", err)
                done <- 0
                return
            }

            ws.WriteMessage(websocket.TextMessage, line)
        }
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/piper"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer c.Close()
	done := make(chan int)

        go pipeStdin(c, done)

	for {
		select {
		case <-done:
                    log.Printf("Done chan")
                    return

		case <-interrupt:
                    log.Println("interrupt")

                    // Cleanly close the connection by sending a close message and then
                    // waiting (with timeout) for the server to close the connection.
                    err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
                    if err != nil {
                            log.Println("write close:", err)
                            return
                    }
                    select {
                    case <-done:
                    case <-time.After(time.Second):
                    }
                    return
		}
	}
}
