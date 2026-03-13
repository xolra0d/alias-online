package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var pingTimeout = sync.OnceValue(func() time.Duration {
	pingT, err := strconv.ParseUint(os.Getenv("PING_TIMEOUT"), 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Duration(pingT) * time.Second
})

//type Room struct {
//	id string
//
//	lock sync.RWMutex
//}

type Rooms struct {
}

func (r *Rooms) ServeHTTP(writer http.ResponseWriter, reader *http.Request, roomId string) error {
	c, err := upgrader.Upgrade(writer, reader, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	var wLock sync.Mutex

	ping := make(chan struct{})
	done := make(chan struct{})

	c.SetPingHandler(func(appData string) error {
		ping <- struct{}{}
		wLock.Lock()
		defer wLock.Unlock()
		return c.WriteMessage(websocket.PongMessage, []byte{})
	})

	go func() {
		t := time.NewTicker(pingTimeout())
		for {
			select {
			case <-t.C:
				done <- struct{}{}
				c.Close()
				return
			case <-ping:
				t.Reset(pingTimeout())
			}
		}
	}()

	for {
		mt, msg, err := c.ReadMessage()
		fmt.Println("SERVER: GOT MESSAGE:", mt, msg, err)
	}

}
