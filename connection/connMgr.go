//连接管理器
//把所有的连接放入管理器，才能进行推送

package connection

import (
	"sync"
)

type ConnMgr struct {
	rwMux   sync.RWMutex
	id2conn map[uint64]*Connection //所有的连接列表
	num     uint64                 //管理的连接数量
}

//添加连接
func (connMgr *ConnMgr) AddConn(conn *Connection) (err error) {
	connMgr.rwMux.Lock()
	defer connMgr.rwMux.Unlock()
	connMgr.id2conn[conn.connId] = conn
	connMgr.num++

	return
}

//删除连接
func (connMgr *ConnMgr) DelConn(conn *Connection) (err error) {
	connMgr.rwMux.Lock()
	defer connMgr.rwMux.Unlock()
	delete(connMgr.id2conn, conn.connId)
	connMgr.num--

	return
}

//给所有的在线连接推送消息
// 存在的问题，连接随时上下线，加锁又影响性能。
func (connMgr *ConnMgr) PushAll(msg []byte) (err error) {
	connMgr.rwMux.Lock()
	defer connMgr.rwMux.Unlock()
	for _, conn := range connMgr.id2conn {
		conn.WriteMessage(msg)
	}

	return
}

//初始化连接管理器
func InitConnMgr() (mgr *ConnMgr) {
	mgr = &ConnMgr{
		id2conn: make(map[uint64]*Connection),
		num:     0,
	}

	return
}
