package wait_time

import (
	"log"
	"sync"
	"time"
)

// WaitTimeCache 等待时间缓存结构
type WaitTimeCache struct {
	Data      map[string]ScenicSpotWaitTime `json:"data"`
	UpdatedAt time.Time                     `json:"updated_at"`
}

// CacheManager 缓存管理器
type CacheManager struct {
	cache map[string]*WaitTimeCache // key: 景区ID, value: 缓存数据
	mutex sync.RWMutex
}

var (
	// 全局缓存管理器实例
	GlobalCacheManager *CacheManager
	once               sync.Once
)

// GetCacheManager 获取全局缓存管理器实例 (单例模式)
func GetCacheManager() *CacheManager {
	once.Do(func() {
		GlobalCacheManager = &CacheManager{
			cache: make(map[string]*WaitTimeCache),
		}
	})
	return GlobalCacheManager
}

// SetWaitTimeData 设置等待时间数据到缓存
func (cm *CacheManager) SetWaitTimeData(scenicAreaId string, data map[string]ScenicSpotWaitTime) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.cache[scenicAreaId] = &WaitTimeCache{
		Data:      data,
		UpdatedAt: time.Now(),
	}

	log.Printf("缓存已更新: 景区ID=%s, 景点数量=%d, 更新时间=%s",
		scenicAreaId, len(data), time.Now().Format("2006-01-02 15:04:05"))
}

// GetWaitTimeData 从缓存获取等待时间数据
func (cm *CacheManager) GetWaitTimeData(scenicAreaId string) (map[string]ScenicSpotWaitTime, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	cache, exists := cm.cache[scenicAreaId]
	if !exists {
		return nil, false
	}

	return cache.Data, true
}

// GetCacheInfo 获取缓存信息 (用于调试)
func (cm *CacheManager) GetCacheInfo(scenicAreaId string) (*WaitTimeCache, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	cache, exists := cm.cache[scenicAreaId]
	if !exists {
		return nil, false
	}

	return cache, true
}

// ClearCache 清空指定景区的缓存
func (cm *CacheManager) ClearCache(scenicAreaId string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	delete(cm.cache, scenicAreaId)
	log.Printf("已清空缓存: 景区ID=%s", scenicAreaId)
}

// ClearAllCache 清空所有缓存
func (cm *CacheManager) ClearAllCache() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.cache = make(map[string]*WaitTimeCache)
	log.Printf("已清空所有缓存")
}

// GetAllCachedScenicAreas 获取所有已缓存的景区ID列表
func (cm *CacheManager) GetAllCachedScenicAreas() []string {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	areas := make([]string, 0, len(cm.cache))
	for scenicAreaId := range cm.cache {
		areas = append(areas, scenicAreaId)
	}
	return areas
}
