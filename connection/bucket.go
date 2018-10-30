package connection

import "sync"

type Bucket struct {
	bucketId int
	id2conn  map[uint64]*Connection //所有的连接列表
	rwmux    sync.RWMutex
	num      uint64 //管理的连接数量
}

func (b *Bucket) AddConn(conn *Connection) (err error) {
	b.rwmux.Lock()
	defer b.rwmux.Unlock()
	b.id2conn[conn.connId] = conn
	b.num++
	return
}

func (b *Bucket) DelConn(connId uint64) (err error) {
	b.rwmux.Lock()
	defer b.rwmux.Unlock()
	delete(b.id2conn, connId)
	b.num--
	return
}

//给bucket内所有在线连接推送消息
func (b *Bucket) PushAll(msg []byte) (err error) {
	b.rwmux.Lock()
	defer b.rwmux.Unlock()
	for _, conn := range b.id2conn {
		conn.WriteMessage(msg)
	}

	return
}

//初始化bucket
func InitBucket(bucketId int) (b *Bucket) {
	b = &Bucket{
		id2conn:  make(map[uint64]*Connection),
		num:      0,
		bucketId: bucketId,
	}

	return
}
