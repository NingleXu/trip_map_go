# 北京环球影城等待时间接入指南

## 概述

已成功将北京环球影城接入等待时间缓存系统。系统会自动定时抓取北京环球影城的景点等待时间数据，并提供快速的API响应。

## 当前状态

✅ **已完成的工作:**
- 创建北京环球影城数据抓取器 (`core/wait_time/beijing_universal_studios.go`)
- 注册到等待时间处理系统
- 实现定时数据抓取
- 提供管理接口用于查看原始数据
- 系统编译通过，准备就绪

⚠️ **待完善的工作:**
- 景点ID映射 (需要根据实际数据完善)

## API 数据结构

### 北京环球影城 API 返回数据示例

```json
{
    "ret": 0,
    "msg": "OK",
    "data": {
        "list": [
            {
                "id": "66d6d4e493335431b63da613",
                "title": "81号公馆",
                "subtitle": "变形金刚基地",
                "waiting_time": -1,
                "gems_copywriting": "关闭",
                "is_closed": 0,
                "service_time": {
                    "open": "15:00",
                    "close": "22:00"
                }
            }
        ]
    }
}
```

### 等待时间状态解析

| `waiting_time` | `is_closed` | `gems_copywriting` | 解析结果 |
|----------------|-------------|-------------------|----------|
| -1             | 0           | "关闭"            | 关闭状态 |
| -1             | 1           | -                 | 关闭状态 |
| 0              | 0           | -                 | 开放，无需等待 |
| >0             | 0           | -                 | 开放，等待X分钟 |
| -               | -           | "维护"            | 维护状态 |

## 使用方法

### 1. 查看原始数据 (用于ID映射)

```bash
GET /api/tripMap/manage/waitTime/rawData?cId=BEIJING_UNIVERSAL
```

这个接口会返回北京环球影城的原始API数据，包含所有景点信息，可用于分析和完善ID映射。

### 2. 测试直接数据抓取

```bash
GET /api/tripMap/manage/waitTime/direct?cId=BEIJING_UNIVERSAL
```

直接调用数据抓取逻辑，用于测试连接性和数据解析。

### 3. 查看缓存状态

```bash
GET /api/tripMap/manage/waitTime/cacheStatus
```

查看北京环球影城是否已成功缓存数据。

## ID 映射完善

### 当前状态
- 系统已预留ID映射机制
- 需要根据实际景点数据完善映射关系

### 完善步骤

1. **获取原始数据**
   ```bash
   GET /api/tripMap/manage/waitTime/rawData?cId=BEIJING_UNIVERSAL
   ```

2. **分析景点列表**
   从返回的 `data.list` 中获取所有景点的：
   - `id`: 北京环球影城的景点ID
   - `title`: 景点名称
   - `subtitle`: 所属区域

3. **建立映射关系**
   在 `core/wait_time/beijing_universal_studios.go` 中的 `IDMapping` 添加映射：
   ```go
   IDMapping: map[string]string{
       "66d6d4e493335431b63da613": "系统内部景点ID1", // 81号公馆
       "68b006084e43b17a297f8f33": "系统内部景点ID2", // 画皮幽宅
       // ... 更多映射
   }
   ```

### 主要景点列表 (基于提供的数据)

| 环球影城ID | 景点名称 | 所属区域 | 系统ID (待映射) |
|------------|----------|----------|----------------|
| 66d6d4e493335431b63da613 | 81号公馆 | 变形金刚基地 | 待分配 |
| 68b006084e43b17a297f8f33 | 画皮幽宅 | 小黄人乐园 | 待分配 |
| 5fa247615ba7f1491f6289a2 | 变形金刚：火种源争夺战 | 变形金刚基地 | 待分配 |
| 5fa22624dcb5e53c6c72ab42 | 哈利·波特与禁忌之旅 | 哈利·波特的魔法世界™ | 待分配 |
| 5f914afc26509774c642cf42 | 神偷奶爸小黄人闹翻天 | 小黄人乐园 | 待分配 |
| 5f913de126509774c642cf36 | 侏罗纪世界大冒险 | 侏罗纪世界努布拉岛 | 待分配 |
| 5f912fd8d72a471ef6359d02 | 霸天虎过山车 | 变形金刚基地 | 待分配 |

## 景区配置

### 当前景区ID
- 系统中使用: `BEIJING_UNIVERSAL`
- 建议: 后续可以改为实际的景区ID (如数据库中的景区ID)

### 修改景区ID
如需修改景区ID，需要同时更新：
1. `core/wait_time/config.go` 中的 `ScenicAreaWaitTimeHandlerMap`
2. 相关的API调用参数

## 数据抓取频率

- **抓取时间**: 每天 7:00-22:00
- **抓取频率**: 每 5 分钟一次
- **自动启动**: 应用启动时自动开始

## 故障排查

### 1. 数据抓取失败
检查日志中的错误信息：
```
景区 BEIJING_UNIVERSAL 等待时间数据抓取失败: [具体错误]
```

可能原因：
- 网络连接问题
- API接口变更
- 请求头参数过期

### 2. 缓存为空
使用管理接口强制刷新：
```bash
POST /api/tripMap/manage/waitTime/forceRefresh
```

### 3. ID映射错误
日志会显示：
```
北京环球影城景点等待时间item:[景点数据],找不到对应的景点id
```

需要完善 `IDMapping` 中的映射关系。

## 测试建议

1. **启动应用后立即测试**
   ```bash
   GET /api/tripMap/manage/waitTime/cacheStatus
   ```

2. **查看原始数据确认连接性**
   ```bash
   GET /api/tripMap/manage/waitTime/rawData?cId=BEIJING_UNIVERSAL
   ```

3. **测试数据解析**
   ```bash
   GET /api/tripMap/manage/waitTime/direct?cId=BEIJING_UNIVERSAL
   ```

## 后续改进

1. **动态令牌处理**: 如果API需要动态令牌，可以考虑添加令牌获取逻辑
2. **错误重试机制**: 增强网络请求的重试逻辑
3. **数据验证**: 添加对返回数据的有效性验证
4. **性能监控**: 添加更详细的性能监控指标

## 联系方式

如有问题或需要技术支持，请查看系统日志或使用管理接口进行故障排查。
