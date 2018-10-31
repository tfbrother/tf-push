package connection

import (
	"encoding/json"
	"sync/atomic"
)

type Stats struct {
	connNum    int64 //当前连接总数
	bucketNum  int64 //bucket的数量
	revMsgNum  int64 //接收消息总数
	pushMsgNum int64 //推送消息总数
}

var (
	G_stats *Stats
)

func InitStats() (err error) {
	G_stats = &Stats{}

	return
}

func (s *Stats) IncrConnNum() {
	atomic.AddInt64(&s.connNum, 1)
}

func (s *Stats) DescConnNum() {
	atomic.AddInt64(&s.connNum, -1)
}

func (s *Stats) IncrBucketNum() {
	atomic.AddInt64(&s.bucketNum, 1)
}

func (s *Stats) DescBucketNum() {
	atomic.AddInt64(&s.bucketNum, -1)
}

func (s *Stats) IncrRevMsgNum() {
	atomic.AddInt64(&s.revMsgNum, 1)
}

func (s *Stats) IncrPushMsgNum() {
	atomic.AddInt64(&s.pushMsgNum, 1)
}

func (s *Stats) Dump() (data []byte, err error) {
	return json.Marshal(s)
}
