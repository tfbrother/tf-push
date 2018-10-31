package connection

import (
	"encoding/json"
	"sync/atomic"
)

type Stats struct {
	ConnNum    int64 `json:"ConnNum"`    //当前连接总数
	BucketNum  int64 `json:"BucketNum"`  //bucket的数量
	RevMsgNum  int64 `json:"RevMsgNum"`  //接收消息总数
	PushMsgNum int64 `json:"PushMsgNum"` //推送消息总数
}

var (
	G_stats *Stats
)

func InitStats() (err error) {
	G_stats = &Stats{
		ConnNum:    0,
		BucketNum:  0,
		RevMsgNum:  0,
		PushMsgNum: 0,
	}

	return
}

func (s *Stats) IncrConnNum() {
	atomic.AddInt64(&s.ConnNum, 1)
}

func (s *Stats) DescConnNum() {
	atomic.AddInt64(&s.ConnNum, -1)
}

func (s *Stats) IncrBucketNum() {
	atomic.AddInt64(&s.BucketNum, 1)
}

func (s *Stats) DescBucketNum() {
	atomic.AddInt64(&s.BucketNum, -1)
}

func (s *Stats) IncrRevMsgNum() {
	atomic.AddInt64(&s.RevMsgNum, 1)
}

func (s *Stats) IncrPushMsgNum() {
	atomic.AddInt64(&s.PushMsgNum, 1)
}

func (s *Stats) Dump() (data []byte, err error) {
	return json.Marshal(s)
}
