package main

// create a websocket client that connects to the server
// and sends a message every second

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	//"time"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

var (
	logger      zerolog.Logger
	ctx, cancel = context.WithCancel(context.Background())
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Str("service", "bigmac.com/channel_test").
		Caller().
		Logger()
}

func main() {

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		s := <-sig
		logger.Info().Str("signal", s.String()).Msg("received signal")
		cancel()
	}()

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	fmt.Printf("connecting to %s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				return
			}
			fmt.Printf("recv: %s\n", message)
		}
	}()

	// ticker := time.NewTicker(time.Second)
	// defer ticker.Stop()

	// for {
	// 	select {
	// 	case <-done:
	// 		return
	// 	case t := <-ticker.C:
	// 		err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
	// 		if err != nil {
	// 			fmt.Println("write:", err)
	// 			return
	// 		}
	// 	}
	// }

	<-ctx.Done()
}
