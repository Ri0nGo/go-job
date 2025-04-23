package notify

import (
	"context"
	"fmt"
	"go-job/internal/model"
	"go-job/internal/pkg/email"
	"log/slog"
	"sync"
	"time"
)

var (
	onceMemory          sync.Once
	memoryNotifyStore   INotifyStore
	defaultWorkerNumber = 5
)

type MemoryNotifyStore struct {
	mux       sync.RWMutex
	notifyMap map[int]NotifyConfig // 任务启用通知列表
	queue     chan NotifyUnit      // 需要的队列
	workerNum int                  // worker数量

	emailSvc email.IEmailService // 邮件发送服务
}

// Set 添加一个任务的通知配置
func (m *MemoryNotifyStore) Set(ctx context.Context, jobId int, config NotifyConfig) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.notifyMap[jobId] = config
	return nil
}

// Get 获取一个任务的通知配置
func (m *MemoryNotifyStore) Get(ctx context.Context, jobId int) (nc NotifyConfig, ok bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	notify, ok := m.notifyMap[jobId]
	if !ok {
		return NotifyConfig{}, false
	}
	return notify, true
}

// Delete 删除一个任务的通知配置
func (m *MemoryNotifyStore) Delete(ctx context.Context, jobId int) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	delete(m.notifyMap, jobId)
	return nil
}

// PushNotifyUnit 推送通知单元到队列中
func (m *MemoryNotifyStore) PushNotifyUnit(ctx context.Context, jobId int, unit NotifyUnit) error {
	if m.needNotify(unit) {
		select {
		case m.queue <- unit:
		default:
			// TODO 队列满了则通知会被丢弃，后续可以实现一个补偿机制，可以通过一张表来记录丢弃的通知
			slog.Info("memory notify queue is full", "jobId", jobId)
		}
		return nil
	}
	return nil
}

// startWorker 启动worker
func (m *MemoryNotifyStore) startWorker() {
	for i := 0; i < m.workerNum; i++ {
		go m.worker()
	}
}

// worker 处理队列中的通知单元
func (m *MemoryNotifyStore) worker() {
	for unit := range m.queue {
		if err := m.dispatch(unit); err != nil {
			slog.Error("notification failed", "jobId", unit.JobID, "err", err)
		}
	}
}

// needNotify 是否需要发送通知
func (m *MemoryNotifyStore) needNotify(unit NotifyUnit) bool {
	switch unit.NotifyStrategy {
	case model.NotifyAfterSuccess:
		return unit.Status == model.Success
	case model.NotifyAfterFailed:
		return unit.Status == model.Failed
	case model.NotifyAlways:
		return true
	default:
		return false
	}
}

func (m *MemoryNotifyStore) dispatch(unit NotifyUnit) error {
	switch unit.NotifyType {
	case model.NotifyTypeEmail:
		tpl := email.GetEmailTpl(email.EmailJobNotifyTpl)
		subject := fmt.Sprintf(tpl.Subject, unit.Name)
		content := fmt.Sprintf(tpl.Content, unit.Name, unit.Status.String(),
			unit.StartExecTime.Format(time.DateTime), unit.Duration,
			unit.Output, unit.Error)
		return m.emailSvc.Send(context.Background(), []string{unit.NotifyMark}, subject, content)
	default:
		return nil
	}
}

func (m *MemoryNotifyStore) generateSubject(name, status string) string {
	return fmt.Sprintf("任务: %s, 状态: %s", name, status)
}

func newMemoryNotifyStore(workerNum int, emailSvc email.IEmailService) INotifyStore {
	return &MemoryNotifyStore{
		notifyMap: make(map[int]NotifyConfig),
		queue:     make(chan NotifyUnit, 1024),
		workerNum: workerNum,
		emailSvc:  emailSvc,
	}
}

func InitMemoryNotifyStore(emailSvc email.IEmailService) INotifyStore {
	onceMemory.Do(func() {
		if emailSvc == nil {
			slog.Error("emailSvc is nil")
		}
		// TODO 这里可以通过配置文件来设置worker数量
		memoryNotifyStore = newMemoryNotifyStore(defaultWorkerNumber, emailSvc)
		if store, ok := memoryNotifyStore.(*MemoryNotifyStore); ok {
			store.startWorker()
		}
	})
	return memoryNotifyStore
}
