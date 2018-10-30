//连接管理器
//把所有的连接放入管理器，才能进行推送

package connection

import (
	"errors"
	"strconv"
	"sync"
)

type ConnMgr struct {
	rwMux   sync.RWMutex
	num     uint64 //全局的连接数量
	buckets []*Bucket
	inChan  chan []byte
}

//添加连接
func (connMgr *ConnMgr) AddConn(conn *Connection) (err error) {
	connMgr.rwMux.Lock()
	defer connMgr.rwMux.Unlock()
	connMgr.num++
	bucket := connMgr.GetBucket(conn)
	bucket.AddConn(conn)
	return
}

//删除连接
func (connMgr *ConnMgr) DelConn(conn *Connection) (err error) {
	connMgr.rwMux.Lock()
	defer connMgr.rwMux.Unlock()
	connMgr.num--
	bucket := connMgr.GetBucket(conn)
	bucket.DelConn(conn.connId)

	return
}

//给所有的在线连接推送消息
func (connMgr *ConnMgr) PushAll(msg []byte) (err error) {
	select {
	case connMgr.inChan <- msg:
	default:
		err = errors.New("inChan is full")
	}

	return
}

//给所有的在线连接推送消息
func (connMgr *ConnMgr) pushAll() (err error) {
	for {
		select {
		case msg := <-connMgr.inChan: //读取信息
			for _, bucket := range connMgr.buckets {
				bucket.PushAll([]byte(string(msg) + strconv.Itoa(bucket.bucketId)))
			}
		}
	}

	return
}

func (connMgr *ConnMgr) GetBucket(conn *Connection) (bucket *Bucket) {
	bucket = connMgr.buckets[conn.connId%uint64(len(connMgr.buckets))]
	return
}

//初始化连接管理器
func InitConnMgr(bucketLen int) (mgr *ConnMgr) {
	mgr = &ConnMgr{
		buckets: make([]*Bucket, bucketLen),
		num:     0,
		inChan:  make(chan []byte, 1000), //1000个缓冲
	}

	for bucketIdx, _ := range mgr.buckets {
		mgr.buckets[bucketIdx] = InitBucket(bucketIdx) // 初始化Bucket
	}

	go mgr.pushAll()

	return
}
