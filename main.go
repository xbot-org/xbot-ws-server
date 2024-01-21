package main

import (
	"github.com/gorilla/websocket"
	"github.com/hanson/go-toolbox/utils"
	"github.com/patrickmn/go-cache"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connMap = make(map[string]*websocket.Conn)

var Cache *cache.Cache

func main() {
	Cache = cache.New(20*time.Second, 10*time.Minute)

	loadConfig()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("err: %+v", err)
			return
		}

		defer c.Close()

		log.Println("已连接：", r.RemoteAddr, r.URL.Query().Get("key"))

		connMap[r.URL.Query().Get("key")] = c

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)

			traceId := gjson.Get(string(message), "trace_id").String()

			if traceId != "" {
				Cache.Set(traceId, message, 15*time.Second)
			}
		}
	})
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("err: %+v", err)
			return
		}

		if c, ok := connMap[r.URL.Query().Get("key")]; ok {

			traceId := utils.RandStr(16, 0)

			b, _ = sjson.SetBytes(b, "trace_id", traceId)

			err = c.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				log.Printf("err: %+v", err)
				return
			}

			for i := 0; i < 15; i++ {
				resp, exists := Cache.Get(traceId)
				if exists {
					w.Write(resp.([]byte))
					return
				}
				time.Sleep(time.Second)
			}
			w.Write([]byte(`{"code":500,"msg":"timeout"}`))

		} else {
			log.Printf("err: key of [%s] not found", r.URL.Query().Get("key"))
			return
		}
	})
	err := http.ListenAndServe("localhost:"+Cfg.Port, nil)
	if err != nil {
		log.Printf("err: %+v", err)
		return
	}
}
