package connection

import "github.com/gorilla/websocket"
import "sync"
import "errors"

type Connection struct {
	wsConn    *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte
	isClosed  bool
	mux       sync.Mutex
}

func InitConnection(ws *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:    ws,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
		isClosed:  false,
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
