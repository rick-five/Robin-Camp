# 电影评分 API

本项目是一个电影评分 API 服务，遵循指定的 OpenAPI 规范，提供电影信息管理、评分功能以及票房数据集成。

## 功能特性

1. 电影创建与管理
2. 电影评分提交（支持更新）
3. 评分统计（平均分和评分数量）
4. 电影搜索与分页
5. 与票房 API 集成
6. 健康检查端点

## 技术栈

- **语言**: Go
- **框架**: Gin (Web框架), GORM (ORM库)
- **数据库**: PostgreSQL
- **容器化**: Docker + Docker Compose

## 数据库设计

### 电影表 (movies)
```sql
CREATE TABLE movies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) UNIQUE NOT NULL,
    release_date DATE NOT NULL,
    genre VARCHAR(100) NOT NULL,
    distributor VARCHAR(255),
    budget BIGINT,
    mpa_rating VARCHAR(10),
    box_office JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 评分表 (ratings)
```sql
CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    movie_title VARCHAR(255) NOT NULL,
    rater_id VARCHAR(255) NOT NULL,
    rating DECIMAL(2,1) NOT NULL CHECK (rating >= 0.5 AND rating <= 5.0 AND rating * 2 = ROUND(rating * 2)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(movie_title, rater_id),
    FOREIGN KEY (movie_title) REFERENCES movies(title) ON DELETE CASCADE
);
```

## API 接口

### 电影相关接口

1. 创建电影: `POST /movies`
2. 列出电影: `GET /movies`

### 评分相关接口

1. 提交评分: `POST /movies/{title}/ratings`
2. 获取评分统计: `GET /movies/{title}/rating`

### 健康检查

1. 健康检查: `GET /healthz`

## 环境变量

- `PORT`: 服务端口，默认 8080
- `AUTH_TOKEN`: 写操作认证 Token
- `DB_URL`: 数据库连接字符串
- `BOXOFFICE_URL`: 票房 API 地址
- `BOXOFFICE_API_KEY`: 票房 API 认证 Key

## 部署

使用 Docker Compose 一键部署：

```bash
# 1. 复制并配置环境变量
cp .env.example .env
# 编辑 .env 文件设置你的环境变量

# 2. 启动服务
make docker-up

# 3. 运行测试
make test-e2e

# 4. 停止服务
make docker-down
```

## 设计思路

### 数据库选型和设计

选择 PostgreSQL 作为数据库，因为它：
1. 是一个功能强大的开源关系型数据库
2. 支持 JSONB 类型，方便存储票房数据
3. 支持复杂查询和事务处理
4. 具有良好的数据一致性和可靠性

电影表设计包含了所有必要的字段，并为常用查询字段添加了索引。评分表通过外键关联到电影表，并确保同一用户对同一电影只能有一条评分记录。

### 后端服务选型和设计

选择 Go 语言和 Gin 框架：
1. Go 语言具有高性能和高并发处理能力
2. Gin 是一个轻量级的 Web 框架，性能优秀
3. 代码简洁易维护
4. 部署简单，Docker 镜像体积小

使用 GORM 作为 ORM 库，它提供了自动迁移、关联查询等功能，简化了数据库操作。

### 可优化内容

1. **缓存机制**：对票房 API 调用结果进行缓存，减少重复请求
2. **重试机制**：对票房 API 调用增加重试逻辑，提高系统健壮性
3. **限流功能**：增加 API 限流，防止恶意请求
4. **监控指标**：添加 Prometheus 指标，便于监控系统状态
5. **日志系统**：增强日志功能，便于问题排查
6. **测试覆盖**：增加单元测试和集成测试，提高代码质量