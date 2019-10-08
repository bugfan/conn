package main

import (
	"flag"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"conn/ws"

	"github.com/gorilla/websocket"
)

var (
	Port       string
	MyUpgrader = websocket.Upgrader{
		// Allow cross domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	MyConn Conn // all connection
)

func init() {
	// get cmd args
	flag.StringVar(&Port, "p", "9000", "listen port") // if -p=8080 then listen port is 8080 not 9000
	flag.Parse()
	// init data
	MyConn.Data = make(map[string]*ws.Connection)
}

type Conn struct {
	Data map[string]*ws.Connection
	sync.RWMutex
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		wsConn *websocket.Conn
		err    error
		data   []byte
		conn   *ws.Connection
	)
	// 1.judge auth

	// 2.get uuid
	uuid := "test"
	if wsConn, err = MyUpgrader.Upgrade(w, r, nil); err != nil {
		log.Println("Upgrade fail:", err)
		return
	}
	if conn, err = ws.InitConnection(wsConn); err != nil {
		log.Println("init ws conn fail:", err)
		return
	}
	MyConn.Lock()
	MyConn.Data[uuid] = conn // save conn
	MyConn.Unlock()
	defer func() { // release conn
		wsConn.Close()
		MyConn.Lock()
		delete(MyConn.Data, uuid) //delete conn
		MyConn.Unlock()
	}()
	log.Println("Connect Ok:", uuid)
	// send heartbeat
	go func() {
		var err error
		for {
			time.Sleep(2e9)
			if err = conn.WriteMessage([]byte(`heartbeat`)); err != nil {
				return
			}
		}
	}()
	for {
		if data, err = conn.ReadMessage(); err != nil {
			log.Println("Receive fail:", err)
			conn.Close()
			break
		}
		log.Println("Receive:", string(data))
	}
}

// send testing message to all online user
func TimingFeedback() {
	for {
		time.Sleep(5e9)
		MyConn.RLock()
		for _, v := range MyConn.Data {
			m := []byte(`{"status": "200"}`)
			if err := v.WriteMessage(m); err != nil {
				log.Println("Send fail:", err)
			}
		}
		MyConn.RUnlock()
	}
}
func TestWs(*testing.T) {
	go TimingFeedback()
	http.HandleFunc("/r/v1/ws", wsHandler)
	http.HandleFunc("/api/v1/ws", wsHandler)
	http.ListenAndServe(":"+Port, nil)

	// you can use https
	// http.ListenAndServeTLS(":"+Port, "./cert.pem", "./key.pem", nil)
}
