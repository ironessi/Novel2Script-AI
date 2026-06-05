# Novel2Script-AI 数据库设计文档

## 1. 概述

- 数据库：MySQL 8.0
- 字符集：utf8mb4
- 排序规则：utf8mb4_unicode_ci
- 存储引擎：InnoDB
- 所有表使用 BIGINT 自增主键
- 时间字段使用 DATETIME 类型
- 软删除使用 deleted_at 字段

## 2. 表结构总览

| 表名 | 说明 |
|------|------|
| sys_user | 系统用户表 |
| novel_project | 小说项目表 |
| novel_source_file | 原始文件表 |
| novel_chapter | 小说章节表 |
| ai_task | AI 任务表 |
| character_profile | 人物档案表 |
| plot_event | 剧情事件表 |
| script_version | 剧本版本表 |
| validation_issue | 校验问题表 |
| audit_log | 审计日志表 |

## 3. 表关系

```
sys_user (1) ──▶ (N) novel_project
novel_project (1) ──▶ (N) novel_source_file
novel_project (1) ──▶ (N) novel_chapter
novel_project (1) ──▶ (N) ai_task
novel_project (1) ──▶ (N) character_profile
novel_project (1) ──▶ (N) plot_event
novel_project (1) ──▶ (N) script_version
script_version (1) ──▶ (N) validation_issue
```

## 4. 表结构详情

### 4.1 sys_user

用户表，存储系统用户信息。

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| id | BIGINT | PK | AUTO_INCREMENT | 主键 |
| username | VARCHAR(64) | YES | - | 用户名，唯一 |
| email | VARCHAR(128) | NO | NULL | 邮箱，唯一 |
| password_hash | VARCHAR(255) | YES | - | 密码哈希（bcrypt） |
| role | VARCHAR(32) | YES | 'user' | 角色：admin / user |
| status | VARCHAR(32) | YES | 'active' | 状态：active / disabled |
| last_login_at | DATETIME | NO | NULL | 最后登录时间 |
| created_at | DATETIME | YES | - | 创建时间 |
| updated_at | DATETIME | YES | - | 更新时间 |

索引：
- PRIMARY KEY (id)
- UNIQUE KEY uk_username (username)
- UNIQUE KEY uk_email (email)

### 4.2 novel_project

小说项目表。

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| id | BIGINT | PK | AUTO_INCREMENT | 主键 |
| owner_id | BIGINT | YES | - | 所属用户 ID |
| title | VARCHAR(255) | YES | - | 项目标题 |
| description | TEXT | NO | NULL | 项目描述 |
| adaptation_mode | VARCHAR(64) | YES | 'screen_script' | 改编模式 |
| visibility | VARCHAR(32) | YES | 'private' | 可见性 |
| status | VARCHAR(32) | YES | 'created' | 状态 |
| created_at | DATETIME | YES | - | 创建时间 |
| updated_at | DATETIME | YES | - | 更新时间 |
| deleted_at | DATETIME | NO | NULL | 软删除时间 |

adaptation_mode 枚举：
- `screen_script` - 影视剧本
- `stage_play` - 舞台剧
- `short_video` - 短视频分镜
- `radio_drama` - 广播剧

visibility 枚举：
- `private` - 私有
- `public` - 公开

status 枚举：
- `created` - 已创建
- `uploaded` - 已上传小说
- `processing` - AI 处理中
- `completed` - 已完成
- `archived` - 已归档

索引：
- PRIMARY KEY (id)
- INDEX idx_owner_id (owner_id)

### 4.3 novel_source_file

原始上传文件表。

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| id | BIGINT | PK | AUTO_INCREMENT | 主键 |
| project_id | BIGINT | YES | - | 所属项目 ID |
| owner_id | BIGINT | YES | - | 上传用户 ID |
| original_filename | VARCHAR(255) | YES | - | 原始文件名 |
| storage_path | VARCHAR(512) | YES | - | 系统存储路径 |
| file_hash | CHAR(64) | YES | - | 文件 SHA256 哈希 |
| file_size | BIGINT | YES | - | 文件大小（字节） |
| mime_type | VARCHAR(128) | YES | - | MIME 类型 |
| scan_status | VARCHAR(32) | YES | 'pending' | 安全扫描状态 |
| created_at | DATETIME | YES | - | 上传时间 |

scan_status 枚举：
- `pending` - 待扫描
- `clean` - 安全
- `infected` - 检测到风险

索引：
- PRIMARY KEY (id)
- INDEX idx_project_id (project_id)
- INDEX idx_file_hash (file_hash)

### 4.4 novel_chapter

小说章节表。

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| id | BIGINT | PK | AUTO_INCREMENT | 主键 |
| project_id | BIGINT | YES | - | 所属项目 ID |
| chapter_index | INT | YES | - | 章节序号 |
| chapter_title | VARCHAR(255) | NO | NULL | 章节标题 |
| content | LONGTEXT | YES | - | 章节内容 |
| content_hash | CHAR(64) | YES | - | 内容 SHA256 哈希 |
| created_at | DATETIME | YES | - | 创建时间 |

索引：
- PRIMARY KEY (id)
- INDEX idx_project_chapter (project_id, chapter_index)

### 4.5 ai_task

AI 生成任务表。

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| id | BIGINT | PK | AUTO_INCREMENT | 主键 |
| project_id | BIGINT | YES | - | 所属项目 ID |
| owner_id | BIGINT | YES | - | 发起用户 ID |
| task_type | VARCHAR(64) | YES | - | 任务类型 |
| status | VARCHAR(32) | YES | 'pending' | 任务状态 |
| progress | INT | YES | 0 | 进度百分比 |
| current_step | VARCHAR(128) | NO | NULL | 当前步骤 |
| error_message | TEXT | NO | NULL | 错误信息 |
| retry_count | INT | YES | 0 | 重试次数 |
| started_at | DATETIME | NO | NULL | 开始时间 |
| finished_at | DATETIME | NO | NULL | 完成时间 |
| created_at | DATETIME | YES | - | 创建时间 |
| updated_at | DATETIME | YES | - | 更新时间 |

task_type 枚举：
- `full_generate` - 完整生成
- `character_extract` - 仅人物抽取
- `plot_extract` - 仅剧情提取
- `scene_split` - 仅场景拆分
- `validate` - 仅校验

status 枚举：
- `pending` - 等待执行
- `running` - 执行中
- `validating` - 校验中
- `repairing` - 修复中
- `completed` - 已完成
- `failed` - 失败
- `cancelled` - 已取消

索引：
- PRIMARY KEY (id)
- INDEX idx_project_id (project_id)
- INDEX idx_owner_id (owner_id)
- INDEX idx_status (status)

### 4.6 character_profile

人物档案表（Character Bible）。

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| id | BIGINT | PK | AUTO_INCREMENT | 主键 |
| project_id | BIGINT | YES | - | 所属项目 ID |
| character_key | VARCHAR(64) | YES | - | 人物标识符 |
| name | VARCHAR(128) | YES | - | 人物名称 |
| aliases | JSON | NO | NULL | 别名列表 |
| role_type | VARCHAR(64) | NO | NULL | 角色类型 |
| description | TEXT | NO | NULL | 人物描述 |
| personality | JSON | NO | NULL | 性格特征 |
| relationships | JSON | NO | NULL | 人物关系 |
| source_refs | JSON | NO | NULL | 原文引用 |
| confidence | DECIMAL(5,2) | YES | 0.00 | 抽取置信度 |
| created_at | DATETIME | YES | - | 创建时间 |
| updated_at | DATETIME | YES | - | 更新时间 |

role_type 枚举：
- `protagonist` - 主角
- `antagonist` - 反派
- `supporting` - 配角
- `minor` - 龙套

索引：
- PRIMARY KEY (id)
- UNIQUE KEY uk_project_character (project_id, character_key)

### 4.7 plot_event

剧情事件表（Plot Event Chain）。

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| id | BIGINT | PK | AUTO_INCREMENT | 主键 |
| project_id | BIGINT | YES | - | 所属项目 ID |
| event_key | VARCHAR(64) | YES | - | 事件标识符 |
| chapter_index | INT | YES | - | 所在章节序号 |
| trigger_text | TEXT | NO | NULL | 触发描述 |
| action_text | TEXT | NO | NULL | 行动描述 |
| result_text | TEXT | NO | NULL | 结果描述 |
| importance | VARCHAR(32) | NO | NULL | 重要程度 |
| source_refs | JSON | NO | NULL | 原文引用 |
| confidence | DECIMAL(5,2) | YES | 0.00 | 抽取置信度 |
| created_at | DATETIME | YES | - | 创建时间 |

importance 枚举：
- `high` - 关键事件
- `medium` - 重要事件
- `low` - 次要事件

索引：
- PRIMARY KEY (id)
- INDEX idx_project_id (project_id)

### 4.8 script_version

剧本版本表。

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| id | BIGINT | PK | AUTO_INCREMENT | 主键 |
| project_id | BIGINT | YES | - | 所属项目 ID |
| owner_id | BIGINT | YES | - | 创建用户 ID |
| version_no | INT | YES | - | 版本号 |
| yaml_content | LONGTEXT | YES | - | YAML 剧本内容 |
| yaml_hash | CHAR(64) | YES | - | YAML 内容哈希 |
| validation_status | VARCHAR(32) | YES | 'pending' | 校验状态 |
| hallucination_risk | VARCHAR(32) | YES | 'unknown' | 幻觉风险等级 |
| safety_risk | VARCHAR(32) | YES | 'unknown' | 安全风险等级 |
| created_by | VARCHAR(32) | YES | 'ai' | 创建来源 |
| created_at | DATETIME | YES | - | 创建时间 |

validation_status 枚举：
- `pending` - 待校验
- `valid` - 校验通过
- `invalid` - 校验不通过
- `repaired` - 已修复

hallucination_risk / safety_risk 枚举：
- `unknown` - 未检测
- `low` - 低风险
- `medium` - 中风险
- `high` - 高风险

created_by 枚举：
- `ai` - AI 生成
- `user` - 用户编辑
- `system_repair` - 系统自动修复

索引：
- PRIMARY KEY (id)
- UNIQUE KEY uk_project_version (project_id, version_no)

### 4.9 validation_issue

校验问题表。

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| id | BIGINT | PK | AUTO_INCREMENT | 主键 |
| project_id | BIGINT | YES | - | 所属项目 ID |
| script_version_id | BIGINT | YES | - | 关联剧本版本 ID |
| issue_type | VARCHAR(64) | YES | - | 问题类型 |
| severity | VARCHAR(32) | YES | - | 严重程度 |
| message | TEXT | YES | - | 问题描述 |
| location_path | VARCHAR(255) | NO | NULL | 问题位置路径 |
| suggestion | TEXT | NO | NULL | 修复建议 |
| resolved | TINYINT | YES | 0 | 是否已解决 |
| created_at | DATETIME | YES | - | 创建时间 |

issue_type 枚举：
- `schema_error` - Schema 校验错误
- `missing_field` - 缺少必填字段
- `invalid_reference` - 无效引用
- `hallucination` - 幻觉问题
- `safety_issue` - 安全问题
- `format_error` - 格式错误

severity 枚举：
- `high` - 高
- `medium` - 中
- `low` - 低

索引：
- PRIMARY KEY (id)
- INDEX idx_script_version_id (script_version_id)

### 4.10 audit_log

审计日志表。

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| id | BIGINT | PK | AUTO_INCREMENT | 主键 |
| user_id | BIGINT | NO | NULL | 操作用户 ID |
| project_id | BIGINT | NO | NULL | 关联项目 ID |
| action | VARCHAR(128) | YES | - | 操作类型 |
| resource_type | VARCHAR(64) | NO | NULL | 资源类型 |
| resource_id | BIGINT | NO | NULL | 资源 ID |
| ip_address | VARCHAR(64) | NO | NULL | 客户端 IP |
| user_agent | VARCHAR(512) | NO | NULL | User-Agent |
| request_id | VARCHAR(128) | NO | NULL | 请求追踪 ID |
| created_at | DATETIME | YES | - | 创建时间 |

action 示例：
- `user.register` - 用户注册
- `user.login` - 用户登录
- `project.create` - 创建项目
- `project.delete` - 删除项目
- `file.upload` - 上传文件
- `task.create` - 创建任务
- `script.generate` - 生成剧本
- `script.edit` - 编辑剧本
- `script.export` - 导出剧本

索引：
- PRIMARY KEY (id)
- INDEX idx_user_id (user_id)
- INDEX idx_project_id (project_id)
- INDEX idx_action (action)

## 5. 安全约束

1. **密码存储**：只存储 bcrypt/argon2 哈希，禁止明文
2. **参数化查询**：所有 SQL 必须使用 ORM 或参数化查询
3. **权限校验**：所有 project 相关查询必须附带 owner_id 校验
4. **软删除**：项目使用 deleted_at 进行软删除
5. **审计追踪**：关键操作记录到 audit_log
