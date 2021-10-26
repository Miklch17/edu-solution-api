package batcher

import (
	"github.com/gammazero/workerpool"
	"time"
)
type Batch interface {
	AddItem(item func())
	RunPool()
}

type BatchStruct struct {
	batch  []func()
	size int
	workerPool *workerpool.WorkerPool
}

func NewBatch(size int, workerPool *workerpool.WorkerPool) *BatchStruct {
	return &BatchStruct{
		batch: make([]func(), size),
		size: size,
		workerPool: workerPool,
	}
}
func (b *BatchStruct) RunPool(){
	for _, item := range b.batch {
		b.workerPool.Submit(item)
	}
	for b.workerPool.WaitingQueueSize() != 0 {
		time.Sleep(time.Microsecond)
	}
}
func (b *BatchStruct) AddItem(item func()) {
	if len(b.batch) == b.size {
		b.RunPool()
		b.batch = make([]func(), b.size)
	}
	b.batch = append(b.batch, item)
}
