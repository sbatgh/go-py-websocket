package main

// create a web server that listens on port 8080
// and responds with "Hello, World!" to every request

import (
	"fmt"
	//"io"
	"net/http"
	"time"

	ws "github.com/sbatgh/go-py-websocket/websocket"
)

var chancounter = make(chan int, 1)

func counter() {
	i := 0
	for {
		i++
		time.Sleep(1 * time.Second)
		chancounter <- i
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := ws.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}
	go ws.Writer(ws)
	ws.Reader(ws)
}

func setupRouteWs() {
	http.HandleFunc("/ws", serveWs)
}

func main() {
	go counter()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	setupRouteWs()
	http.ListenAndServe(":8080", nil)
}
