-- ============================================================
-- Novel2Script-AI 数据库迁移脚本
-- 用于版本升级时的表结构变更
-- ============================================================

USE novel2script;

-- ============================================================
-- Migration 001: 初始版本
-- 对应 deploy/mysql-init.sql 中的表结构
-- 版本号：1.0.0
-- ============================================================

-- 记录迁移版本
CREATE TABLE IF NOT EXISTS schema_migrations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    version VARCHAR(32) NOT NULL UNIQUE COMMENT '迁移版本号',
    description VARCHAR(255) COMMENT '迁移描述',
    applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '执行时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据库迁移记录表';

INSERT INTO schema_migrations (version, description)
VALUES ('1.0.0', '初始表结构：sys_user, novel_project, novel_source_file, novel_chapter, ai_task, character_profile, plot_event, script_version, validation_issue, audit_log')
ON DUPLICATE KEY UPDATE version=version;

-- ============================================================
-- Migration 002: 预留扩展字段
-- 为未来功能预留的字段变更
-- 版本号：1.1.0
-- ============================================================

-- 示例：为 novel_project 添加标签字段
-- ALTER TABLE novel_project ADD COLUMN tags JSON COMMENT '项目标签' AFTER description;

-- 示例：为 character_profile 添加头像字段
-- ALTER TABLE character_profile ADD COLUMN avatar_url VARCHAR(512) COMMENT '人物头像URL' AFTER description;

-- 示例：为 ai_task 添加优先级字段
-- ALTER TABLE ai_task ADD COLUMN priority INT NOT NULL DEFAULT 0 COMMENT '任务优先级' AFTER task_type;

-- INSERT INTO schema_migrations (version, description)
-- VALUES ('1.1.0', '添加扩展字段：project.tags, character.avatar_url, task.priority');
