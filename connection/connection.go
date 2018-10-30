package connection

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	wsConn    *websocket.Conn
	inChan    chan []byte //接收队列
	outChan   chan []byte //发送队列
	closeChan chan byte   //是否关闭chan
	isClosed  bool        //是否关闭
	mux       sync.Mutex
	connId    uint64
}

func InitConnection(connId uint64, ws *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:    ws,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
		isClosed:  false,
		connId:    connId,
	}

	//启动读协程
	go conn.readLoop()
	//启动写协程
	go conn.writeLoop()
	return
}

func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("conn is closed")
	}
	return
}

func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("conn is closed")
	}

	return
}

func (conn *Connection) Close() {
	//线程安全
	conn.wsConn.Close()

	//只能执行一次，且要保证线程安全
	conn.mux.Lock()
	if !conn.isClosed {
		conn.isClosed = true
		close(conn.closeChan)
	}
	conn.mux.Unlock()
}

func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)

	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			goto ERR
		}
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			goto ERR
		}
	}
ERR:
	conn.Close()
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)

	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			goto ERR
		}
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.Close()
}
