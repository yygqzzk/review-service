// pkg/pool/pool.go
package pool

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/panjf2000/ants/v2"
)

// WorkerPool 协程池封装
type WorkerPool struct {
	pool *ants.Pool
	log  *log.Helper
}

// Task 任务接口
type Task interface {
	Execute(ctx context.Context) error
	GetID() string
}

// Config 协程池配置
type Config struct {
	Size         int               // 协程池大小
	ExpiryTime   time.Duration     // 空闲协程过期时间
	PreAlloc     bool              // 是否预分配
	MaxBlocking  int               // 最大阻塞等待数
	PanicHandler func(interface{}) // panic处理函数
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Size:        runtime.NumCPU() * 2, // 默认为CPU核数的2倍
		ExpiryTime:  time.Minute * 10,     // 10分钟过期
		PreAlloc:    false,
		MaxBlocking: 100,
		PanicHandler: func(p interface{}) {
			log.Errorf("goroutine panic: %v", p)
		},
	}
}

// NewWorkerPool 创建协程池
func NewWorkerPool(config *Config, logger log.Logger) (*WorkerPool, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 创建ants池
	pool, err := ants.NewPool(
		config.Size,
		ants.WithExpiryDuration(config.ExpiryTime),
		ants.WithPreAlloc(config.PreAlloc),
		ants.WithMaxBlockingTasks(config.MaxBlocking),
		ants.WithPanicHandler(config.PanicHandler),
	)
	if err != nil {
		return nil, err
	}

	return &WorkerPool{
		pool: pool,
		log:  log.NewHelper(logger),
	}, nil
}

// Submit 提交任务
func (w *WorkerPool) Submit(task Task) error {
	return w.pool.Submit(func() {
		ctx := context.Background()
		if err := task.Execute(ctx); err != nil {
			w.log.WithContext(ctx).Errorf("task %s execute failed: %v", task.GetID(), err)
		}
	})
}

// SubmitWithContext 带上下文提交任务
func (w *WorkerPool) SubmitWithContext(ctx context.Context, task Task) error {
	return w.pool.Submit(func() {
		if err := task.Execute(ctx); err != nil {
			w.log.WithContext(ctx).Errorf("task %s execute failed: %v", task.GetID(), err)
		}
	})
}

// SubmitFunc 提交函数
func (w *WorkerPool) SubmitFunc(fn func()) error {
	return w.pool.Submit(fn)
}

// BatchSubmit 批量提交任务
func (w *WorkerPool) BatchSubmit(tasks []Task) []error {
	var wg sync.WaitGroup
	errors := make([]error, len(tasks))

	for i, task := range tasks {
		wg.Add(1)
		i, task := i, task // 避免闭包问题

		err := w.pool.Submit(func() {
			defer wg.Done()
			ctx := context.Background()
			if err := task.Execute(ctx); err != nil {
				errors[i] = err
				w.log.WithContext(ctx).Errorf("batch task %s execute failed: %v", task.GetID(), err)
			}
		})

		if err != nil {
			errors[i] = err
			wg.Done()
		}
	}

	wg.Wait()
	return errors
}

// Stats 获取池状态
func (w *WorkerPool) Stats() (running, free, capacity int) {
	return w.pool.Running(), w.pool.Free(), w.pool.Cap()
}

// Close 关闭协程池
func (w *WorkerPool) Close() {
	w.pool.Release()
}
