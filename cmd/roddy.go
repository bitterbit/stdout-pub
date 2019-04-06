package main

import (
	"flag"
	"log"
        "github.com/bitterbit/piper-roddy/server"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var index []byte

func main() {
	flag.Parse()
	log.SetFlags(0)
        roddy := server.NewRoddyServer("static/index.html")
        log.Fatal(roddy.Start(*addr))
}

