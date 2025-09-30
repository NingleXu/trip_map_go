package schedule

import (
	"context"
	"log"
	"time"
	"trip-map/core/wait_time"
	"trip-map/internal/service"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	cacheManager *wait_time.CacheManager
	ctx          context.Context
	cancel       context.CancelFunc
	running      bool
}

// NewScheduler 创建新的调度器
func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		cacheManager: wait_time.GetCacheManager(),
		ctx:          ctx,
		cancel:       cancel,
		running:      false,
	}
}

// Start 启动定时任务调度器
func (s *Scheduler) Start() {
	if s.running {
		log.Printf("等待时间调度器已经在运行中")
		return
	}

	s.running = true
	log.Printf("等待时间调度器启动")

	// 启动定时任务
	go s.scheduleLoop()
}

// Stop 停止定时任务调度器
func (s *Scheduler) Stop() {
	if !s.running {
		return
	}

	s.cancel()
	s.running = false
	log.Printf("等待时间调度器已停止")
}

// scheduleLoop 调度循环，每到整5分钟执行一次
func (s *Scheduler) scheduleLoop() {
	for {
		// 计算距离下一个整5分钟的时间
		now := time.Now()
		// 下一个5分钟点的分钟数
		next := now.Truncate(5 * time.Minute).Add(5 * time.Minute)
		sleepDuration := time.Until(next)

		select {
		case <-s.ctx.Done():
			log.Printf("调度器上下文取消，退出调度循环")
			return
		case <-time.After(sleepDuration):
			// 检查是否在工作时间范围内 (7:00-22:00)
			if s.isWorkingHours() {
				s.executeTask()
			} else {
				log.Printf("当前时间不在工作时间范围内，跳过数据抓取")
			}
		}
	}
}

// isWorkingHours 检查当前是否在工作时间 (7:00-22:00)
func (s *Scheduler) isWorkingHours() bool {
	now := time.Now()
	hour := now.Hour()
	return hour >= 7 && hour < 22
}

// executeTask 执行数据抓取任务
func (s *Scheduler) executeTask() {
	log.Printf("开始执行等待时间数据抓取任务 - %s", time.Now().Format("2006-01-02 15:04:05"))

	// 遍历所有支持等待时间的景区
	for scenicAreaId, handler := range wait_time.ScenicAreaWaitTimeHandlerMap {
		s.fetchAndCacheWaitTime(scenicAreaId, handler)
	}

	log.Printf("等待时间数据抓取任务完成 - %s", time.Now().Format("2006-01-02 15:04:05"))
}

// fetchAndCacheWaitTime 抓取并缓存单个景区的等待时间数据
func (s *Scheduler) fetchAndCacheWaitTime(scenicAreaId string, handler wait_time.WaitTimeHandler) {
	startTime := time.Now()

	waitTimeData, err := handler.GetWaitTime()
	if err != nil {
		log.Printf("抓取景区 %s 等待时间数据失败: %v", scenicAreaId, err)
		return
	}

	// 将数据存入缓存
	s.cacheManager.SetWaitTimeData(scenicAreaId, waitTimeData)

	// 存入数据库表
	err = service.BatchRecord[wait_time.ScenicSpotWaitTime](service.BizTypeScenic, waitTimeData)
	if err != nil {
		log.Printf("景区的排队时间持久化失败！%v", err)
	}

	duration := time.Since(startTime)
	log.Printf("景区 %s 等待时间数据抓取成功，耗时: %v, 景点数量: %d",
		scenicAreaId, duration, len(waitTimeData))
}

// GetStatus 获取调度器状态
func (s *Scheduler) GetStatus() bool {
	return s.running
}

// ForceExecute 强制执行一次数据抓取任务 (用于测试或手动触发)
func (s *Scheduler) ForceExecute() {
	log.Printf("强制执行等待时间数据抓取任务")
	s.executeTask()
}

// 全局调度器实例
var globalScheduler *Scheduler

// StartGlobalScheduler 启动全局调度器
func StartGlobalScheduler() {
	if globalScheduler == nil {
		globalScheduler = NewScheduler()
	}
	globalScheduler.Start()
}

// StopGlobalScheduler 停止全局调度器
func StopGlobalScheduler() {
	if globalScheduler != nil {
		globalScheduler.Stop()
	}
}

// GetGlobalScheduler 获取全局调度器实例
func GetGlobalScheduler() *Scheduler {
	return globalScheduler
}
