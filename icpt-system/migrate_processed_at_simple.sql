-- 简化的数据库迁移脚本
-- 添加 processed_at 字段到 images 表

-- 添加 processed_at 字段
ALTER TABLE images 
ADD COLUMN processed_at TIMESTAMP NULL DEFAULT NULL AFTER file_size;

-- 创建索引以提高查询性能
CREATE INDEX idx_images_processed_at ON images(processed_at);

-- 为已完成和失败的图像设置 processed_at 时间
-- 使用 created_at 作为近似的处理完成时间（仅对历史数据）
UPDATE images 
SET processed_at = created_at 
WHERE status IN ('completed', 'failed') 
AND processed_at IS NULL;

-- 查看更新结果
SELECT 
    status,
    COUNT(*) as count,
    COUNT(processed_at) as with_processed_at
FROM images 
GROUP BY status; 