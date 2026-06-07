-- ============================================================
-- Novel2Script-AI 数据库初始化脚本
-- 数据库：MySQL 8.0
-- 字符集：utf8mb4
-- ============================================================

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS novel2script
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE novel2script;

-- ============================================================
-- 1. 用户表 sys_user
-- ============================================================
CREATE TABLE IF NOT EXISTS sys_user (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(64) NOT NULL UNIQUE COMMENT '用户名',
    email VARCHAR(128) UNIQUE COMMENT '邮箱',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希（bcrypt）',
    role VARCHAR(32) NOT NULL DEFAULT 'user' COMMENT '角色：admin / user',
    status VARCHAR(32) NOT NULL DEFAULT 'active' COMMENT '状态：active / disabled',
    last_login_at DATETIME NULL COMMENT '最后登录时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统用户表';

-- ============================================================
-- 2. 小说项目表 novel_project
-- ============================================================
CREATE TABLE IF NOT EXISTS novel_project (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    owner_id BIGINT NOT NULL COMMENT '所属用户 ID',
    title VARCHAR(255) NOT NULL COMMENT '项目标题',
    description TEXT COMMENT '项目描述',
    adaptation_mode VARCHAR(64) NOT NULL DEFAULT 'screen_script' COMMENT '改编模式：screen_script / stage_play / short_video / radio_drama',
    visibility VARCHAR(32) NOT NULL DEFAULT 'private' COMMENT '可见性：private / public',
    status VARCHAR(32) NOT NULL DEFAULT 'created' COMMENT '状态：created / uploaded / processing / completed / archived',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME NULL COMMENT '软删除时间',
    INDEX idx_owner_id (owner_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='小说项目表';

-- ============================================================
-- 3. 原始文件表 novel_source_file
-- ============================================================
CREATE TABLE IF NOT EXISTS novel_source_file (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    project_id BIGINT NOT NULL COMMENT '所属项目 ID',
    owner_id BIGINT NOT NULL COMMENT '上传用户 ID',
    original_filename VARCHAR(255) NOT NULL COMMENT '原始文件名',
    storage_path VARCHAR(512) NOT NULL COMMENT '系统存储路径',
    file_hash CHAR(64) NOT NULL COMMENT '文件 SHA256 哈希',
    file_size BIGINT NOT NULL COMMENT '文件大小（字节）',
    mime_type VARCHAR(128) NOT NULL COMMENT 'MIME 类型',
    scan_status VARCHAR(32) NOT NULL DEFAULT 'pending' COMMENT '安全扫描状态：pending / clean / infected',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
    INDEX idx_project_id (project_id),
    INDEX idx_file_hash (file_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='原始上传文件表';

-- ============================================================
-- 4. 小说章节表 novel_chapter
-- ============================================================
CREATE TABLE IF NOT EXISTS novel_chapter (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    project_id BIGINT NOT NULL COMMENT '所属项目 ID',
    chapter_index INT NOT NULL COMMENT '章节序号',
    chapter_title VARCHAR(255) COMMENT '章节标题',
    content LONGTEXT NOT NULL COMMENT '章节内容',
    content_hash CHAR(64) NOT NULL COMMENT '内容 SHA256 哈希',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_project_chapter (project_id, chapter_index)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='小说章节表';

-- ============================================================
-- 5. AI 任务表 ai_task
-- ============================================================
CREATE TABLE IF NOT EXISTS ai_task (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    project_id BIGINT NOT NULL COMMENT '所属项目 ID',
    owner_id BIGINT NOT NULL COMMENT '发起用户 ID',
    task_type VARCHAR(64) NOT NULL COMMENT '任务类型：full_generate / character_extract / plot_extract / scene_split / validate',
    status VARCHAR(32) NOT NULL DEFAULT 'pending' COMMENT '状态：pending / running / validating / repairing / completed / failed / cancelled',
    progress INT NOT NULL DEFAULT 0 COMMENT '进度百分比 0-100',
    current_step VARCHAR(128) COMMENT '当前执行步骤',
    error_message TEXT COMMENT '错误信息',
    retry_count INT NOT NULL DEFAULT 0 COMMENT '重试次数',
    started_at DATETIME NULL COMMENT '开始时间',
    finished_at DATETIME NULL COMMENT '完成时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_project_id (project_id),
    INDEX idx_owner_id (owner_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI 任务表';

-- ============================================================
-- 6. 人物档案表 character_profile
-- ============================================================
CREATE TABLE IF NOT EXISTS character_profile (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    project_id BIGINT NOT NULL COMMENT '所属项目 ID',
    character_key VARCHAR(64) NOT NULL COMMENT '人物标识符',
    name VARCHAR(128) NOT NULL COMMENT '人物名称',
    aliases JSON COMMENT '别名列表',
    role_type VARCHAR(64) COMMENT '角色类型：protagonist / antagonist / supporting / minor',
    description TEXT COMMENT '人物描述',
    personality JSON COMMENT '性格特征',
    relationships JSON COMMENT '人物关系',
    source_refs JSON COMMENT '原文引用',
    confidence DECIMAL(5,2) DEFAULT 0.00 COMMENT '抽取置信度',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_project_character (project_id, character_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='人物档案表（Character Bible）';

-- ============================================================
-- 7. 剧情事件表 plot_event
-- ============================================================
CREATE TABLE IF NOT EXISTS plot_event (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    project_id BIGINT NOT NULL COMMENT '所属项目 ID',
    event_key VARCHAR(64) NOT NULL COMMENT '事件标识符',
    chapter_index INT NOT NULL COMMENT '所在章节序号',
    trigger_text TEXT COMMENT '触发描述',
    action_text TEXT COMMENT '行动描述',
    result_text TEXT COMMENT '结果描述',
    importance VARCHAR(32) COMMENT '重要程度：high / medium / low',
    source_refs JSON COMMENT '原文引用',
    confidence DECIMAL(5,2) DEFAULT 0.00 COMMENT '抽取置信度',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_project_id (project_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='剧情事件表（Plot Event Chain）';

-- ============================================================
-- 8. 剧本版本表 script_version
-- ============================================================
CREATE TABLE IF NOT EXISTS script_version (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    project_id BIGINT NOT NULL COMMENT '所属项目 ID',
    owner_id BIGINT NOT NULL COMMENT '创建用户 ID',
    version_no INT NOT NULL COMMENT '版本号',
    yaml_content LONGTEXT NOT NULL COMMENT 'YAML 剧本内容',
    yaml_hash CHAR(64) NOT NULL COMMENT 'YAML 内容哈希',
    validation_status VARCHAR(32) NOT NULL DEFAULT 'pending' COMMENT '校验状态：pending / valid / invalid / repaired',
    hallucination_risk VARCHAR(32) NOT NULL DEFAULT 'unknown' COMMENT '幻觉风险：unknown / low / medium / high',
    safety_risk VARCHAR(32) NOT NULL DEFAULT 'unknown' COMMENT '安全风险：unknown / low / medium / high',
    created_by VARCHAR(32) NOT NULL DEFAULT 'ai' COMMENT '创建来源：ai / user / system_repair',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_project_version (project_id, version_no)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='剧本版本表';

-- ============================================================
-- 9. 校验问题表 validation_issue
-- ============================================================
CREATE TABLE IF NOT EXISTS validation_issue (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    project_id BIGINT NOT NULL COMMENT '所属项目 ID',
    script_version_id BIGINT NOT NULL COMMENT '关联剧本版本 ID',
    issue_type VARCHAR(64) NOT NULL COMMENT '问题类型：schema_error / missing_field / invalid_reference / hallucination / safety_issue / format_error',
    severity VARCHAR(32) NOT NULL COMMENT '严重程度：high / medium / low',
    message TEXT NOT NULL COMMENT '问题描述',
    location_path VARCHAR(255) COMMENT '问题位置路径',
    suggestion TEXT COMMENT '修复建议',
    resolved TINYINT NOT NULL DEFAULT 0 COMMENT '是否已解决：0=否 1=是',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_script_version_id (script_version_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='校验问题表';

-- ============================================================
-- 10. 审计日志表 audit_log
-- ============================================================
CREATE TABLE IF NOT EXISTS audit_log (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT COMMENT '操作用户 ID',
    project_id BIGINT COMMENT '关联项目 ID',
    action VARCHAR(128) NOT NULL COMMENT '操作类型',
    resource_type VARCHAR(64) COMMENT '资源类型',
    resource_id BIGINT COMMENT '资源 ID',
    ip_address VARCHAR(64) COMMENT '客户端 IP',
    user_agent VARCHAR(512) COMMENT 'User-Agent',
    request_id VARCHAR(128) COMMENT '请求追踪 ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_user_id (user_id),
    INDEX idx_project_id (project_id),
    INDEX idx_action (action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计日志表';

-- ============================================================
-- 初始化管理员账号（密码：admin123456，仅用于开发环境）
-- bcrypt hash of "admin123456"
-- ============================================================
INSERT INTO sys_user (username, email, password_hash, role, status, created_at, updated_at)
VALUES (
    'admin',
    'admin@novel2script.local',
    '$2y$10$yLn5o5SS1.5SaEs02eUqve7yRPwfXsFcuirAT5qLQwqtJbmhDIT4W',
    'admin',
    'active',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE username=username;
