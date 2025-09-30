# 等待时间缓存系统

## 概述

该系统将原本的"每次请求实时抓取"的方式改进为"定时缓存+快速响应"的架构，显著提升了系统性能和用户体验。

## 系统架构

### 1. 缓存管理器 (`cache_manager.go`)
- 提供线程安全的内存缓存
- 支持多景区数据并发管理
- 提供缓存状态查询和清理功能

### 2. 定时调度器 (`scheduler.go`)
- 每天 7:00-22:00 时间段内运行
- 每 5 分钟自动抓取一次数据
- 支持手动触发和状态查询

### 3. 服务层改进 (`internal/service/wait_time.go`)
- 优先从缓存获取数据
- 缓存缺失时自动回退到直接抓取
- 提供详细的日志记录

## 主要功能

### 自动定时抓取
- **时间范围**: 每天 7:00 - 22:00
- **频率**: 每 5 分钟一次
- **覆盖**: 所有支持等待时间的景区

### 快速响应
- 从内存缓存直接返回数据
- 响应时间从秒级降低到毫秒级
- 减少对第三方 API 的依赖

### 容错机制
- 缓存缺失时自动回退到直接抓取
- 抓取失败不影响已有缓存数据
- 详细的错误日志和状态监控

## API 接口

### 用户接口 (原有接口保持不变)
- `GET /api/tripMap/getWaitTimeScenicAreaList` - 获取支持等待时间的景区列表
- `GET /api/tripMap/getScenicAreaWaitTimeList?cId=景区ID` - 获取景区等待时间

### 管理接口 (新增)
- `GET /api/tripMap/manage/waitTime/cacheStatus` - 查看缓存状态
- `POST /api/tripMap/manage/waitTime/forceRefresh` - 强制刷新缓存
- `DELETE /api/tripMap/manage/waitTime/clearCache?scenicAreaId=景区ID` - 清空缓存
- `GET /api/tripMap/manage/waitTime/direct?cId=景区ID` - 直接获取数据 (测试用)

## 使用示例

### 查看缓存状态
```bash
GET /api/tripMap/manage/waitTime/cacheStatus
```

响应示例:
```json
{
  "code": 200,
  "data": {
    "scheduler_running": true,
    "total_cached_areas": 1,
    "cached_areas": {
      "7469417": {
        "updated_at": "2025-09-29 14:30:05",
        "scenic_spots": 19
      }
    }
  }
}
```

### 强制刷新缓存
```bash
POST /api/tripMap/manage/waitTime/forceRefresh
```

### 清空特定景区缓存
```bash
DELETE /api/tripMap/manage/waitTime/clearCache?scenicAreaId=7469417
```

### 清空所有缓存
```bash
DELETE /api/tripMap/manage/waitTime/clearCache
```

## 部署说明

系统在应用启动时会自动初始化：

1. 在 `main.go` 中调用 `wait_time.StartGlobalScheduler()`
2. 调度器自动开始工作
3. 首次启动会立即执行一次数据抓取

## 监控和维护

### 日志监控
系统会产生以下关键日志：
- 调度器启动/停止
- 数据抓取成功/失败
- 缓存更新状态
- API 调用情况

### 性能指标
- 缓存命中率: 正常运行时应接近 100%
- 响应时间: 从缓存获取数据 < 50ms
- 数据新鲜度: 最长 5 分钟延迟

### 故障处理
1. **调度器异常**: 检查系统资源和网络连接
2. **数据抓取失败**: 检查第三方 API 可用性
3. **缓存缺失**: 使用管理接口强制刷新

## 扩展支持

要添加新的景区等待时间支持：

1. 实现 `WaitTimeHandler` 接口
2. 在 `ScenicAreaWaitTimeHandlerMap` 中注册
3. 系统会自动包含到定时抓取中

## 配置建议

### 生产环境
- 建议配置监控告警
- 定期检查缓存状态
- 监控第三方 API 调用频率

### 开发环境
- 可使用管理接口进行测试
- 支持手动清空缓存重新抓取
- 提供直接获取接口用于调试
