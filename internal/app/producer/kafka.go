package producer

import (
	"context"
	"github.com/ozonmp/edu-solution-api/internal/batcher"
	"github.com/ozonmp/edu-solution-api/internal/model"
	"github.com/ozonmp/edu-solution-api/internal/sender"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
)

type Producer interface {
	Start(ctx context.Context)
	Close()
}

type producer struct {
	n       uint64
	timeout time.Duration

	sender sender.EventSender
	events <-chan model.SolutionEvent

	workerPool *workerpool.WorkerPool

	wg   *sync.WaitGroup
}

// todo for students: add repo
func NewKafkaProducer(
	n uint64,
	sender sender.EventSender,
	events <-chan model.SolutionEvent,
	workerPool *workerpool.WorkerPool,
) Producer {

	wg := &sync.WaitGroup{}

	return &producer{
		n:          n,
		sender:     sender,
		events:     events,
		workerPool: workerPool,
		wg:         wg,
	}
}

type Batch interface {
	AddItem(item func())
	RunPool()
}

func (p *producer) Start(ctx context.Context) {
	for i := uint64(0); i < p.n; i++ {
		p.wg.Add(1)
		go func() {
			batchSize := 100
			var batchPool Batch
			batchPool = batcher.NewBatch(batchSize, p.workerPool)
			defer p.wg.Done()
			doneEvent := false
			for {
				select {
				case event := <-p.events:
					if doneEvent {
						batchPool.AddItem(func() {
							event.SetStatus(model.Created, model.Deferred)
							//event.Unlock()???
						})
						continue
					}
					if err := p.sender.Send(&event); err != nil {
						batchPool.AddItem(func() {
							event.SetStatus(model.Removed, model.Processed)
							//event.Unlock()???
						})
					} else {
						batchPool.AddItem(func() {
							event.SetStatus(model.Created, model.Deferred)
							//event.Unlock()???
						})
					}
				case <-ctx.Done():
					doneEvent = true
					continue
				default:
						if doneEvent {
							batchPool.RunPool()
							return
						}
				}
			}
		}()
	}
}

func (p *producer) Close() {
	p.wg.Wait()
}


/*
func main() {
	var allTask []*workerpool.Task
	for i := 1; i <= 100; i++ {
		task := workerpool.NewTask(func(data interface{}) error {
			taskID := data.(int)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Task %d processed\n", taskID)
			return nil
		}, i)
		allTask = append(allTask, task)
	}

	pool := workerpool.NewPool(allTask, 5)
	pool.Run()
}
*/