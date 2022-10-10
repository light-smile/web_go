package provider

import (
	"flag"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "127.0.0.1:3000", "http service address")

func TestInitWebSocket(t *testing.T) {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/2"}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go timeWriter(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		fmt.Printf("received: %s\n", message)
		// fmt.Println(message)
	}
}

func timeWriter(conn *websocket.Conn) {
	for {
		time.Sleep(time.Second * 5)
		conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
	}
}
