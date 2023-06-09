package cron

import (
	"github.com/hq2005001/modules/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"sync"
)

var (
	// Cron cron客户端
	Cron Crontab
	once sync.Once
)

// Crontab 定时任务
type Crontab struct {
	client   *cron.Cron
	handlers map[string]CrontabItem
	mutex    *sync.Mutex
	logger   *logger.Logger
}

// CrontabItem 计划任务item
type CrontabItem interface {
	Name() string
	Rule() string
	Handle()
}

// Register 注册
func (c *Crontab) Register(item CrontabItem) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handlers[item.Name()] = item
}

// Handle 处理计划任务
func (c *Crontab) Handle() {
	for name, item := range c.handlers {
		id, err := c.client.AddFunc(item.Rule(), item.Handle)
		if err != nil {
			c.logger.Error("添加计划任务失败", zap.Int("id", int(id)), zap.Error(err))
			continue
		}
		c.logger.Debug("添加计划任务成功, ", zap.String("name", name))
	}
	c.client.Run()
}

func (c *Crontab) Get(name string) (CrontabItem, bool) {
	crontab, exists := c.handlers[name]
	return crontab, exists
}

// New 新建一个cron
func New(logf *logger.Logger) *Crontab {
	once.Do(func() {
		Cron = Crontab{
			client:   cron.New(cron.WithSeconds()),
			handlers: make(map[string]CrontabItem),
			mutex:    &sync.Mutex{},
			logger:   logf,
		}
	})
	return &Cron
}
