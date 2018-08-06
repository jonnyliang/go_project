//go 语言实战中学习到的包
package work

import (
	"sync"
)

//worker 接口
type Worker interface {
	Task()
}

//地址池用于完成提交的Worker任务
type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

//New 创建一个新的工作池
func New(maxGoroutines int) *Pool {
	p := Pool{
		work: make(chan Worker),
	}
	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.work {
				w.Task()
			}
			p.wg.Done()
		}()
	}
	return &p
}

//Run 提交工作到工作池
func (p *Pool) Run(w Worker) {
	p.work <- w
}

// Shutdown 等待所有gorountine 停止工作
func (p *Pool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
