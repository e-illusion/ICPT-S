# ICPT 系统数据库设计文档

<div align="center">

![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1?style=for-the-badge&logo=mysql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-6.0+-DC382D?style=for-the-badge&logo=redis&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-ORM-00ADD8?style=for-the-badge)

**ICPT 高性能图像处理系统数据库设计**

</div>

## 📋 目录

- [数据库概述](#数据库概述)
- [系统架构](#系统架构)
- [数据库配置](#数据库配置)
- [表结构设计](#表结构设计)
- [关系设计](#关系设计)
- [索引策略](#索引策略)
- [数据迁移](#数据迁移)
- [性能优化](#性能优化)
- [安全策略](#安全策略)
- [备份策略](#备份策略)

## 🎯 数据库概述

ICPT 系统采用 **MySQL + Redis** 混合存储架构，其中：

- **MySQL 8.0+**: 主数据存储，负责用户数据、图像元数据的持久化
- **Redis 6.0+**: 缓存和任务队列，负责会话管理、任务调度、实时数据缓存
- **GORM**: Go语言ORM框架，提供类型安全的数据库操作

### 🌟 设计原则

1. **数据一致性**: 严格的外键约束和事务控制
2. **性能优化**: 合理的索引设计和查询优化  
3. **扩展性**: 支持水平扩展的表结构设计
4. **安全性**: 密码加密、SQL注入防护
5. **可维护性**: 清晰的表命名和字段注释

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                    应用层 (Go + GORM)                        │
├─────────────────────┬───────────────────┬─────────────────────┤
│                     │                   │                     │
│   ┌─────────────────▼─────────────────┐ │ ┌─────────────────▼───┐ │
│   │         MySQL 8.0+              │ │ │      Redis 6.0+     │ │
│   │  ┌─────────────┬─────────────┐   │ │ │ ┌─────────────────┐ │ │
│   │  │    users    │   images    │   │ │ │ │   任务队列      │ │ │
│   │  │    表       │    表       │   │ │ │ │   用户会话      │ │ │
│   │  └─────────────┴─────────────┘   │ │ │ │   缓存数据      │ │ │
│   │                                  │ │ │ └─────────────────┘ │ │
│   └──────────────────────────────────┘ │ └─────────────────────┘ │
│            持久化存储                    │        临时存储           │
└─────────────────────┬───────────────────┴─────────────────────────┤
                      │                                             │
                 ┌────▼────┐                                        │
                 │  备份   │                                        │
                 │  策略   │                                        │
                 └─────────┘                                        │
                                                                    │
├─────────────────────────────────────────────────────────────────┤
│                      监控和日志                                   │
│  - 慢查询日志    - 连接池监控    - 缓存命中率    - 错误追踪        │
└─────────────────────────────────────────────────────────────────┘
```

## ⚙️ 数据库配置

### MySQL 配置

```yaml
# config.yaml - 数据库配置
database:
  host: "127.0.0.1"         # 数据库主机
  port: 3306                # 数据库端口
  user: "icpt_user"         # 数据库用户
  password: "123"           # 数据库密码 (生产环境请使用强密码)
  dbname: "ICPT_System"     # 数据库名称
```

### 连接参数

```go
// DSN连接字符串配置
dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
    c.User,           // 用户名
    c.Password,       // 密码
    c.Host,          // 主机地址
    c.Port,          // 端口
    c.DBName,        // 数据库名
)

// GORM配置
DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
    DisableForeignKeyConstraintWhenMigrating: false, // 启用外键约束
    NamingStrategy: schema.NamingStrategy{
        SingularTable: false, // 使用复数表名
    },
})
```

### Redis 配置

```yaml
# config.yaml - Redis配置  
redis:
  addr: "127.0.0.1:6379"    # Redis地址
  password: ""              # Redis密码 (如无密码留空)
  db: 0                     # 数据库编号
```

## 📊 表结构设计

### 1. 用户表 (users)

用户表存储系统用户的基本信息和认证数据。

```sql
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户唯一标识',
  `username` varchar(255) NOT NULL COMMENT '用户名',
  `email` varchar(255) NOT NULL COMMENT '邮箱地址',
  `password_hash` varchar(255) NOT NULL COMMENT '密码哈希值',
  `status` varchar(50) NOT NULL DEFAULT 'active' COMMENT '用户状态: active, inactive, banned',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`),
  UNIQUE KEY `idx_users_email` (`email`),
  KEY `idx_users_status` (`status`),
  KEY `idx_users_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
```

#### 字段说明

| 字段名 | 类型 | 约束 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `id` | bigint unsigned | PK, AUTO_INCREMENT | - | 用户唯一标识 |
| `username` | varchar(255) | NOT NULL, UNIQUE | - | 用户名，3-20字符 |
| `email` | varchar(255) | NOT NULL, UNIQUE | - | 邮箱地址，用于登录 |
| `password_hash` | varchar(255) | NOT NULL | - | bcrypt加密的密码哈希 |
| `status` | varchar(50) | NOT NULL | 'active' | 用户状态枚举 |
| `created_at` | datetime(3) | NOT NULL | NOW() | 账户创建时间 |
| `updated_at` | datetime(3) | NOT NULL | NOW() | 最后更新时间 |

#### Go 模型定义

```go
// User 结构体对应 'users' 表
type User struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    Username     string    `gorm:"type:varchar(255);not null;unique" json:"username"`
    Email        string    `gorm:"type:varchar(255);not null;unique" json:"email"`
    PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"` // 不在JSON中显示
    Status       string    `gorm:"type:varchar(50);not null;default:'active'" json:"status"`
    CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
    return "users"
}
```

### 2. 图像表 (images)

图像表存储用户上传的图像文件信息和处理状态。

```sql
CREATE TABLE `images` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '图像唯一标识',
  `user_id` bigint unsigned NOT NULL COMMENT '所属用户ID',
  `original_filename` varchar(255) NOT NULL COMMENT '原始文件名',
  `storage_path` varchar(1024) NOT NULL COMMENT '原始文件存储路径',
  `thumbnail_path` varchar(1024) DEFAULT NULL COMMENT '缩略图存储路径',
  `status` varchar(50) NOT NULL DEFAULT 'processing' COMMENT '处理状态: processing, completed, failed',
  `error_info` text DEFAULT NULL COMMENT '错误信息详情',
  `file_size` bigint NOT NULL DEFAULT 0 COMMENT '文件大小(字节)',
  `processed_at` timestamp NULL DEFAULT NULL COMMENT '处理完成时间',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_images_user_id` (`user_id`),
  KEY `idx_images_status` (`status`),
  KEY `idx_images_created_at` (`created_at`),
  KEY `idx_images_processed_at` (`processed_at`),
  CONSTRAINT `fk_images_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图像表';
```

#### 字段说明

| 字段名 | 类型 | 约束 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `id` | bigint unsigned | PK, AUTO_INCREMENT | - | 图像唯一标识 |
| `user_id` | bigint unsigned | NOT NULL, FK | - | 关联用户ID |
| `original_filename` | varchar(255) | NOT NULL | - | 用户上传的原始文件名 |
| `storage_path` | varchar(1024) | NOT NULL | - | 服务器存储路径 |
| `thumbnail_path` | varchar(1024) | NULL | - | 缩略图存储路径 |
| `status` | varchar(50) | NOT NULL | 'processing' | 处理状态枚举 |
| `error_info` | text | NULL | - | 处理失败时的错误信息 |
| `file_size` | bigint | NOT NULL | 0 | 文件大小(字节) |
| `processed_at` | timestamp | NULL | - | 处理完成时间戳 |
| `created_at` | datetime(3) | NOT NULL | NOW() | 记录创建时间 |

#### 状态枚举说明

| 状态值 | 说明 | 业务含义 |
|--------|------|---------|
| `processing` | 处理中 | 图像已上传，正在队列中等待或正在处理 |
| `completed` | 已完成 | 图像处理成功，缩略图生成完毕 |
| `failed` | 处理失败 | 图像处理过程中出现错误 |

#### Go 模型定义

```go
// Image 结构体对应 'images' 表
type Image struct {
    ID               uint       `gorm:"primaryKey"`
    UserID           uint       `gorm:"index"`
    OriginalFilename string     `gorm:"type:varchar(255);not null"`
    StoragePath      string     `gorm:"type:varchar(1024);not null"`
    ThumbnailPath    string     `gorm:"type:varchar(1024)"`
    Status           string     `gorm:"type:varchar(50);not null;default:'processing'"`
    ErrorInfo        string     `gorm:"type:text"`
    FileSize         int64      `gorm:"type:bigint;default:0"`
    ProcessedAt      *time.Time `gorm:"index"`
    CreatedAt        time.Time  `gorm:"autoCreateTime"`
}

// TableName 指定表名
func (Image) TableName() string {
    return "images"
}
```

## 🔗 关系设计

### 表间关系

```
users (1) ──────────── (N) images
  │                        │
  │ id                     │ user_id (FK)
  │                        │
  └── 一个用户可以上传多张图像
      用户删除时级联删除所有图像
```

### 外键约束

```sql
-- 图像表外键约束
ALTER TABLE `images` 
ADD CONSTRAINT `fk_images_user_id` 
FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) 
ON DELETE CASCADE ON UPDATE CASCADE;
```

- **ON DELETE CASCADE**: 用户删除时，自动删除该用户的所有图像记录
- **ON UPDATE CASCADE**: 用户ID更新时，自动更新图像表中的关联ID

## 📈 索引策略

### 主键索引

```sql
-- 自动创建的主键索引
PRIMARY KEY (`id`) -- 聚簇索引，查询效率最高
```

### 唯一索引

```sql
-- users表唯一索引
UNIQUE KEY `idx_users_username` (`username`)  -- 用户名唯一约束
UNIQUE KEY `idx_users_email` (`email`)        -- 邮箱唯一约束
```

### 普通索引

```sql
-- users表索引
KEY `idx_users_status` (`status`)           -- 按状态查询用户
KEY `idx_users_created_at` (`created_at`)   -- 按注册时间排序

-- images表索引  
KEY `idx_images_user_id` (`user_id`)        -- 按用户查询图像(最常用)
KEY `idx_images_status` (`status`)          -- 按状态筛选图像
KEY `idx_images_created_at` (`created_at`)  -- 按上传时间排序
KEY `idx_images_processed_at` (`processed_at`) -- 按处理时间排序
```

### 组合索引设计

```sql
-- 高频查询的组合索引
CREATE INDEX `idx_user_status_created` ON `images` (`user_id`, `status`, `created_at`);
-- 优化查询: SELECT * FROM images WHERE user_id = ? AND status = ? ORDER BY created_at DESC

CREATE INDEX `idx_status_processed` ON `images` (`status`, `processed_at`);  
-- 优化查询: SELECT * FROM images WHERE status = 'completed' ORDER BY processed_at DESC
```

### 索引使用建议

| 查询场景 | 推荐索引 | 说明 |
|----------|----------|------|
| 用户登录 | `idx_users_username`, `idx_users_email` | 支持用户名或邮箱登录 |
| 用户图像列表 | `idx_user_status_created` | 支持分页和状态筛选 |
| 图像状态查询 | `idx_images_status` | 快速筛选不同状态图像 |
| 处理时间统计 | `idx_status_processed` | 计算平均处理时间 |

## 🔄 数据迁移

### 自动迁移

系统使用 GORM 的 AutoMigrate 功能进行自动迁移：

```go
// 自动迁移 - 在应用启动时执行
err = DB.AutoMigrate(&models.Image{}, &models.User{})
if err != nil {
    log.Fatalf("错误: 数据库迁移失败: %v", err)
}
```

### 手动迁移脚本

#### 初始化脚本

```sql
-- 创建数据库
CREATE DATABASE `ICPT_System` 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

-- 创建专用用户
CREATE USER 'icpt_user'@'localhost' IDENTIFIED BY '123';
GRANT ALL PRIVILEGES ON ICPT_System.* TO 'icpt_user'@'localhost';
FLUSH PRIVILEGES;
```

#### 字段迁移脚本

```sql
-- migrate_processed_at_simple.sql
-- 添加 processed_at 字段到 images 表

-- 添加字段
ALTER TABLE images 
ADD COLUMN processed_at TIMESTAMP NULL DEFAULT NULL AFTER file_size;

-- 创建索引
CREATE INDEX idx_images_processed_at ON images(processed_at);

-- 更新历史数据
UPDATE images 
SET processed_at = created_at 
WHERE status IN ('completed', 'failed') 
AND processed_at IS NULL;
```

### 迁移最佳实践

1. **备份数据**: 迁移前务必备份数据库
2. **测试环境验证**: 先在测试环境执行迁移
3. **版本控制**: 将迁移脚本纳入版本控制
4. **回滚计划**: 准备回滚脚本应对异常情况

## ⚡ 性能优化

### 查询优化

#### 分页查询优化

```go
// 优化前 - 可能产生性能问题
var images []models.Image
db.Where("user_id = ?", userID).
   Offset((page - 1) * pageSize).
   Limit(pageSize).
   Find(&images)

// 优化后 - 使用游标分页
var images []models.Image
db.Where("user_id = ? AND id < ?", userID, lastID).
   Order("id DESC").
   Limit(pageSize).
   Find(&images)
```

#### 统计查询优化

```go
// 仪表盘统计 - 使用聚合查询
type StatsResult struct {
    TotalImages    int64   `json:"total_images"`
    CompletedCount int64   `json:"completed_count"`
    FailedCount    int64   `json:"failed_count"`
    AvgProcessTime float64 `json:"avg_process_time"`
}

// 一次查询获取所有统计数据
db.Model(&models.Image{}).
   Select(`
       COUNT(*) as total_images,
       SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed_count,
       SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed_count,
       AVG(TIMESTAMPDIFF(MICROSECOND, created_at, processed_at)) as avg_process_time
   `).
   Where("user_id = ?", userID).
   Scan(&result)
```

### 连接池配置

```go
// 配置数据库连接池
sqlDB, err := DB.DB()
if err != nil {
    log.Fatal(err)
}

// 设置最大空闲连接数
sqlDB.SetMaxIdleConns(10)

// 设置最大打开连接数  
sqlDB.SetMaxOpenConns(100)

// 设置连接最大生存时间
sqlDB.SetConnMaxLifetime(time.Hour)
```

### 缓存策略

#### Redis 缓存设计

```go
// 用户信息缓存
userKey := fmt.Sprintf("user:%d", userID)
Rdb.Set(ctx, userKey, userJSON, 30*time.Minute)

// 图像列表缓存 (短期缓存)
listKey := fmt.Sprintf("images:user:%d:page:%d", userID, page)
Rdb.Set(ctx, listKey, imagesJSON, 5*time.Minute)

// 图像状态缓存
statusKey := fmt.Sprintf("image:%d:status", imageID) 
Rdb.Set(ctx, statusKey, status, 10*time.Minute)
```

## 🔒 安全策略

### 密码安全

```go
// 使用 bcrypt 进行密码哈希
func (u *User) HashPassword(password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(password), 
        bcrypt.DefaultCost, // 默认强度 10
    )
    if err != nil {
        return err
    }
    u.PasswordHash = string(hashedPassword)
    return nil
}

// 验证密码
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword(
        []byte(u.PasswordHash), 
        []byte(password),
    )
    return err == nil
}
```

### SQL 注入防护

```go
// GORM 自动防护 SQL 注入
// ✅ 安全 - 使用参数绑定
db.Where("username = ? OR email = ?", username, email).First(&user)

// ❌ 危险 - 字符串拼接 (已避免)
// db.Where(fmt.Sprintf("username = '%s'", username)).First(&user)
```

### 数据访问控制

```go
// 用户只能访问自己的数据
func GetUserImages(userID uint, page int) ([]models.Image, error) {
    var images []models.Image
    err := DB.Where("user_id = ?", userID). // 强制用户隔离
             Order("created_at DESC").
             Offset((page - 1) * pageSize).
             Limit(pageSize).
             Find(&images).Error
    return images, err
}
```

## 💾 备份策略

### 定期备份

```bash
#!/bin/bash
# backup_database.sh - 数据库备份脚本

# 配置
DB_HOST="127.0.0.1"
DB_USER="icpt_user"  
DB_PASS="123"
DB_NAME="ICPT_System"
BACKUP_DIR="/backup/mysql"
DATE=$(date +%Y%m%d_%H%M%S)

# 创建备份目录
mkdir -p ${BACKUP_DIR}

# 执行备份
mysqldump -h${DB_HOST} -u${DB_USER} -p${DB_PASS} \
    --single-transaction \
    --routines \
    --triggers \
    ${DB_NAME} > ${BACKUP_DIR}/icpt_backup_${DATE}.sql

# 压缩备份文件
gzip ${BACKUP_DIR}/icpt_backup_${DATE}.sql

# 删除7天前的备份
find ${BACKUP_DIR} -name "*.sql.gz" -mtime +7 -delete

echo "备份完成: icpt_backup_${DATE}.sql.gz"
```

### 增量备份

```bash
# 启用 binlog 进行增量备份
# my.cnf 配置
[mysqld]
log-bin=mysql-bin
server-id=1
binlog-format=ROW
expire_logs_days=7
```

### Redis 持久化

```bash
# redis.conf 配置
save 900 1      # 900秒内至少1个key变化时保存
save 300 10     # 300秒内至少10个key变化时保存  
save 60 10000   # 60秒内至少10000个key变化时保存

# AOF 持久化
appendonly yes
appendfsync everysec
```

## 📊 监控和诊断

### 性能监控

```sql
-- 慢查询监控
SHOW VARIABLES LIKE 'slow_query_log';
SHOW VARIABLES LIKE 'long_query_time';

-- 查看慢查询
SELECT * FROM mysql.slow_log ORDER BY start_time DESC LIMIT 10;
```

### 连接监控

```sql
-- 查看当前连接
SHOW PROCESSLIST;

-- 查看连接统计
SHOW STATUS LIKE 'Connections';
SHOW STATUS LIKE 'Threads_connected';
```

### 索引分析

```sql
-- 检查索引使用情况
EXPLAIN SELECT * FROM images WHERE user_id = 1 AND status = 'completed';

-- 查看索引统计
SELECT 
    TABLE_NAME,
    INDEX_NAME,
    SEQ_IN_INDEX,
    COLUMN_NAME,
    CARDINALITY
FROM INFORMATION_SCHEMA.STATISTICS 
WHERE TABLE_SCHEMA = 'ICPT_System';
```

## 🚀 扩展规划

### 水平扩展

1. **读写分离**: 主库写入，从库读取
2. **分库分表**: 按用户ID或时间分片
3. **缓存集群**: Redis Cluster 部署

### 新功能表设计

```sql
-- 图像标签表 (计划中)
CREATE TABLE `image_tags` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `image_id` bigint unsigned NOT NULL,
  `tag_name` varchar(100) NOT NULL,
  `confidence` decimal(5,4) DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_image_tags_image_id` (`image_id`),
  KEY `idx_image_tags_tag_name` (`tag_name`),
  CONSTRAINT `fk_image_tags_image_id` FOREIGN KEY (`image_id`) REFERENCES `images` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图像标签表';

-- 处理日志表 (计划中)  
CREATE TABLE `processing_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `image_id` bigint unsigned NOT NULL,
  `worker_id` varchar(100) NOT NULL,
  `start_time` datetime(3) NOT NULL,
  `end_time` datetime(3) DEFAULT NULL,
  `status` varchar(50) NOT NULL,
  `error_message` text,
  `processing_time_ms` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_processing_logs_image_id` (`image_id`),
  KEY `idx_processing_logs_worker_id` (`worker_id`),
  CONSTRAINT `fk_processing_logs_image_id` FOREIGN KEY (`image_id`) REFERENCES `images` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图像处理日志表';
```

## 📋 总结

### 当前数据库特点

- ✅ **2个核心表**: users、images
- ✅ **完整的约束**: 主键、外键、唯一约束
- ✅ **合理的索引**: 覆盖主要查询场景
- ✅ **安全设计**: 密码加密、SQL注入防护
- ✅ **性能优化**: 连接池、缓存策略

### 技术优势

1. **类型安全**: GORM ORM 提供编译时类型检查
2. **自动迁移**: 减少手动数据库操作
3. **混合存储**: MySQL + Redis 发挥各自优势
4. **扩展性好**: 支持水平扩展和功能扩展

### 建议改进

1. **监控完善**: 增加更详细的性能监控
2. **备份自动化**: 实现自动备份和恢复测试
3. **分库分表**: 为大规模用户做准备
4. **读写分离**: 提高查询性能

---

<div align="center">

**🗄️ 企业级数据库设计，支撑高性能图像处理系统！🗄️**

Made with ❤️ by ICPT Team

</div> 