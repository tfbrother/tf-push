package main

import (
	"github.com/gorilla/websocket"
	"github.com/tfbrother/tf-push/connection"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	//生成唯一ID
	serverId = uint64(time.Now().Unix())
	//初始化连接管理器

)

var mgr = connection.InitConnMgr(100)

func wsHande(w http.ResponseWriter, r *http.Request) {
	var (
		wsConn *websocket.Conn
		err    error
		conn   *connection.Connection
		data   []byte
	)
	wsConn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// 连接唯一标识
	connId := atomic.AddUint64(&serverId, 1)

	if conn, err = connection.InitConnection(connId, wsConn); err != nil {
		goto ERR
	}

	mgr.AddConn(conn)

	//心跳
	go func() {
		var err error

		for {
			if err = conn.WriteMessage([]byte("heatbeat")); err != nil {
				return
			}

			time.Sleep(1 * time.Second)
		}
	}()

	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}

		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}
ERR:
	wsConn.Close()
}

//推送接口，基于HTTP协议
//curl  http://localhost:7777/pushAll -d 'msg="msg hello"'
func pushAll(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg string
	)
	if err = r.ParseForm(); err != nil {
		return
	}

	msg = r.PostForm.Get("msg")

	mgr.PushAll([]byte(msg))
	w.Write([]byte(msg))
}

func main() {
	http.HandleFunc("/connect", wsHande)
	http.HandleFunc("/pushAll", pushAll)

	http.ListenAndServe("0.0.0.0:7777", nil)
}
