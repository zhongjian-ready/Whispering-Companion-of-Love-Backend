# 微信小程序后台接口文档

## 基础信息

- **Base URL**: `http://localhost:8080` (本地开发)
- **Content-Type**: `application/json`

---

## 1. 用户相关 (User)

### 1.1 获取用户列表

获取所有用户的基本信息。

- **URL**: `/users`
- **Method**: `GET`
- **Response**:
  ```json
  {
    "count": 10,
    "users": [
      {
        "username": "user1",
        "nickname": "Nick",
        "gender": 1
      }
    ]
  }
  ```

### 1.2 获取用户设置

获取当前用户的个性化设置，包括每日目标、提醒配置等。

- **URL**: `/settings`
- **Method**: `GET`
- **Query Params**:
  - `user_id` (可选，目前 Mock 为 1): 指定用户 ID
- **Response**:
  ```json
  {
    "daily_goal": 2000,
    "reminder_enabled": true,
    "reminder_interval": 60,
    "reminder_start_time": "08:00",
    "reminder_end_time": "22:00",
    "quick_add_presets": [200, 300, 500, 800],
    "reminder_settings": null,
    "quick_add_settings": null
  }
  ```

### 1.3 更新用户设置

更新用户的个性化设置。

- **URL**: `/settings`
- **Method**: `PUT`
- **Body**:
  ```json
  {
    "daily_goal": 2500,
    "reminder_enabled": true,
    "reminder_interval": 90,
    "reminder_start_time": "09:00",
    "reminder_end_time": "23:00"
  }
  ```
- **Response**:
  ```json
  {
    "message": "Settings updated"
  }
  ```

---

## 2. 饮水记录 (Intake)

### 2.1 添加饮水记录

记录一次饮水数据。

- **URL**: `/intake`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "amount": 300
  }
  ```
- **Response**:
  ```json
  {
    "message": "Success",
    "data": {
      "id": 1,
      "user_id": 1,
      "amount": 300,
      "recorded_at": "2025-12-22T16:00:00+08:00",
      "date": "2025-12-22"
    }
  }
  ```

### 2.2 获取今日记录

获取今日的所有饮水记录、总摄入量以及目标完成百分比。

- **URL**: `/intake/today`
- **Method**: `GET`
- **Response**:
  ```json
  {
    "records": [
      {
        "id": 1,
        "amount": 300,
        "recorded_at": "..."
      }
    ],
    "total": 1200,
    "goal": 2000,
    "percentage": 60
  }
  ```

### 2.3 删除饮水记录

删除指定的一条饮水记录。

- **URL**: `/intake/:id`
- **Method**: `DELETE`
- **URL Params**:
  - `id`: 记录 ID
- **Response**:
  ```json
  {
    "message": "Deleted"
  }
  ```

### 2.4 获取周统计数据

获取最近 7 天的每日饮水总量统计。

- **URL**: `/intake/stats/weekly`
- **Method**: `GET`
- **Response**:
  ```json
  {
    "data": [
      {
        "date": "2025-12-16",
        "total": 1500
      },
      {
        "date": "2025-12-17",
        "total": 2100
      },
      ...
    ]
  }
  ```

---

## 3. 成就系统 (Achievements)

### 3.1 获取成就列表

获取所有成就列表，并标记当前用户是否已解锁。

- **URL**: `/achievements`
- **Method**: `GET`
- **Response**:
  ```json
  {
    "data": [
      {
        "id": 1,
        "name": "初次尝试",
        "description": "完成第一次饮水记录",
        "icon_url": "http://...",
        "condition_type": "first_intake",
        "condition_val": 1,
        "is_unlocked": true
      },
      {
        "id": 2,
        "name": "持之以恒",
        "description": "连续打卡7天",
        "icon_url": "http://...",
        "condition_type": "streak",
        "condition_val": 7,
        "is_unlocked": false
      }
    ]
  }
  ```

---

## 4. 健康检查

### 4.1 Ping

服务健康检查接口。

- **URL**: `/ping`
- **Method**: `GET`
- **Response**:
  ```json
  {
    "message": "pong"
  }
  ```
