package log

import "github.com/duanchi/min/abstract"

type Log struct {
	abstract.Service

	ConcurrentLines int64 `value:"Gateway.Log.ConcurrentLines"`

	queue *LogQueue
}

func (this *Log) Init() {
	this.queue = NewQueue(1024 * 1024)
}

func (this *Log) Put(val interface{}) (ok bool, quantity uint32) {
	return this.queue.Put(val)
}

func (this *Log) Record() {
	var n int64
	logList := []map[string]string{}
	for n = 0; n < this.ConcurrentLines; n++ {
		value, ok, _ := this.queue.Get()
		if ok {
			logList = append(logList, value.(map[string]string))
		} else {
			break
		}
	}
}
