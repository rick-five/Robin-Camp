# 电影评分 API 项目设计文档

## 1. 项目概述

本项目需要实现一个电影评分 API 服务，遵循指定的 OpenAPI 规范，提供电影信息管理、评分功能以及票房数据集成。服务需要支持基本的 CRUD 操作、用户评分、数据聚合、搜索和分页功能。

## 2. 技术选型

### 2.1 后端服务选型

**语言：Go**

选择 Go 语言的原因：
- 高性能和高并发处理能力
- 编译型语言，部署简单
- 丰富的标准库和第三方库支持
- 适合构建微服务和 API 服务
- Docker 镜像体积小，启动快

**框架：Gin**

选择 Gin 框架的原因：
- 高性能 HTTP 路由框架
- 简洁易用的 API 设计
- 中间件支持完善
- 内置 JSON 序列化/反序列化
- 社区活跃，文档完善

### 2.2 数据库选型

**数据库：PostgreSQL**

选择 PostgreSQL 的原因：
- 功能强大的开源关系型数据库
- 支持复杂查询和事务处理
- 具有良好的数据一致性和可靠性
- 支持 JSON 类型，便于扩展
- 与 Docker 集成良好
- 支持索引优化，适合搜索和分页场景

### 2.3 数据库访问层

**ORM：GORM**

选择 GORM 的原因：
- Go 语言中功能最完整的 ORM 库之一
- 支持 PostgreSQL
- 提供迁移功能
- 支持关联查询和预加载
- 简化数据库操作代码

## 3. 系统架构设计

### 3.1 整体架构

```
┌─────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│   客户端请求    │────▶│  电影评分 API    │────▶│ 票房 Mock API    │
└─────────────────┘     │    服务          │     │                  │
                        │                  │     │                  │
                        │  ┌─────────────┐ │     └──────────────────┘
                        │  │   Gin路由   │ │              
                        │  └─────────────┘ │     ┌──────────────────┐
                        │  ┌─────────────┐ │     │   PostgreSQL     │
                        │  │  控制器层   │ │◀───▶│     数据库       │
                        │  └─────────────┘ │     │                  │
                        │  ┌─────────────┐ │     └──────────────────┘
                        │  │   服务层    │ │              
                        │  └─────────────┘ │              
                        │  ┌─────────────┐ │              
                        │  │   数据访问  │ │              
                        │  │     层      │ │              
                        │  └─────────────┘ │              
                        └──────────────────┘              
```

### 3.2 数据库设计

#### 3.2.1 电影表 (movies)
```sql
CREATE TABLE movies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) UNIQUE NOT NULL,
    release_date DATE NOT NULL,
    genre VARCHAR(100),
    distributor VARCHAR(255),
    budget BIGINT,
    mpa_rating VARCHAR(10),
    box_office JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### 3.2.2 评分表 (ratings)
```sql
CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_title VARCHAR(255) NOT NULL REFERENCES movies(title) ON DELETE CASCADE,
    rater_id VARCHAR(255) NOT NULL,
    rating DECIMAL(2,1) NOT NULL CHECK (rating >= 0.5 AND rating <= 5.0 AND rating * 2 = ROUND(rating * 2)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(movie_title, rater_id)
);
```

## 4. 核心功能实现

### 4.1 电影创建与票房数据集成

当用户创建电影时：
1. 验证请求数据
2. 创建电影记录（不含票房数据）
3. 异步调用票房 API 获取票房数据
4. 更新电影记录的票房字段
5. 如果票房 API 调用失败，不影响电影创建，票房字段保持为空

### 4.2 评分系统

评分功能采用 Upsert 语义：
1. 用户通过 [X-Rater-Id] 头标识身份
2. 对同一电影的重复评分会更新已有评分
3. 评分值限定在 {0.5, 1.0, ..., 5.0} 范围内

### 4.3 评分聚合

提供电影评分统计接口：
1. 计算平均分并四舍五入保留一位小数
2. 统计评分总数

### 4.4 搜索与分页

支持以下搜索参数：
- q: 关键词搜索（电影标题）
- year: 年份筛选
- genre: 类型筛选
- limit: 每页数量
- cursor: 分页游标

## 5. 安全与认证

### 5.1 写操作认证

所有写操作（创建电影、提交评分）需要通过 Bearer Token 认证：
- 请求头需包含 `Authorization: Bearer {AUTH_TOKEN}`
- Token 从环境变量中读取

### 5.2 评分身份标识

评分提交需要通过 `X-Rater-Id` 头标识评分用户身份。

## 6. 错误处理

严格按照 OpenAPI 规范返回相应状态码：
- 400: 请求数据错误
- 401: 未提供认证信息
- 403: 认证失败
- 404: 资源未找到
- 422: 数据验证失败
- 500: 服务器内部错误

## 7. Docker 部署

### 7.1 Dockerfile 设计

采用多阶段构建：
1. 构建阶段：使用 golang 镜像编译应用
2. 运行阶段：使用 alpine 镜像运行应用，非 root 用户运行

### 7.2 docker-compose.yml 设计

包含两个服务：
1. app: 应用服务，监听 8080 端口
2. db: PostgreSQL 数据库服务，带健康检查

## 8. 数据库迁移

使用 GORM 的自动迁移功能，在应用启动时自动创建表结构。

## 9. 健康检查

提供 `/healthz` 端点，检查应用和数据库连接状态。

## 10. 环境变量配置

所有配置通过环境变量注入：
- PORT: 服务端口
- AUTH_TOKEN: 写操作认证 Token
- DB_URL: 数据库连接字符串
- BOXOFFICE_URL: 票房 API 地址
- BOXOFFICE_API_KEY: 票房 API 认证 Key

## 11. Makefile 目标

提供以下 Make 目标：
- `make docker-up`: 构建并启动服务
- `make docker-down`: 停止并清理服务
- `make test-e2e`: 运行端到端测试

## 12. 可优化内容

1. **缓存机制**：对票房 API 调用结果进行缓存，减少重复请求
2. **重试机制**：对票房 API 调用增加重试逻辑，提高系统健壮性
3. **限流功能**：增加 API 限流，防止恶意请求
4. **监控指标**：添加 Prometheus 指标，便于监控系统状态
5. **日志系统**：增强日志功能，便于问题排查
6. **测试覆盖**：增加单元测试和集成测试，提高代码质量