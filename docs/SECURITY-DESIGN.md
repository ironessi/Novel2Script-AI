# Novel2Script-AI 安全设计文档

## 1. 概述

本文档定义 Novel2Script-AI 项目在代码层面的安全设计规范。虽然本项目为本地开发项目，但仍需遵循基本安全实践，培养安全编码习惯。

## 2. 认证安全

### 2.1 JWT 认证

- 使用 JWT 进行用户认证
- JWT 必须设置过期时间（建议 24 小时）
- JWT Secret 从环境变量读取，禁止硬编码
- 登录接口需要限制失败次数（防暴力破解）

### 2.2 密码安全

- 密码使用 bcrypt 或 argon2 哈希存储
- 禁止存储明文密码
- 密码最小长度 8 位
- 注册和修改密码时需要验证密码强度

### 2.3 环境变量管理

```bash
# .env 示例（不提交到 Git）
JWT_SECRET=change_me_to_a_long_random_string
MYSQL_PASSWORD=change_me
AI_SERVICE_TOKEN=change_me_internal_token
```

`.gitignore` 必须包含：
```
.env
.env.local
```

## 3. 权限控制

### 3.1 角色设计

| 角色 | 说明 |
|------|------|
| admin | 管理员，可访问所有项目 |
| user | 普通用户，只能访问自己的项目 |

### 3.2 项目级权限

所有项目相关接口必须校验权限：

```
当前用户 ID == project.owner_id OR 当前用户.role == "admin"
```

### 3.3 权限校验接口列表

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/projects/{id} | 获取项目详情 |
| PUT | /api/projects/{id} | 修改项目 |
| DELETE | /api/projects/{id} | 删除项目 |
| GET | /api/projects/{id}/chapters | 获取章节 |
| POST | /api/projects/{id}/upload | 上传小说 |
| POST | /api/projects/{id}/generate | 开始生成 |
| GET | /api/projects/{id}/script | 获取剧本 |
| PUT | /api/projects/{id}/script | 修改剧本 |
| POST | /api/projects/{id}/validate | 校验 YAML |
| GET | /api/projects/{id}/export | 导出剧本 |
| GET | /api/projects/{id}/audit | 审计日志 |

### 3.4 防止越权访问

**错误做法**（只按 project_id 查询）：
```sql
SELECT * FROM novel_project WHERE id = ?
```

**正确做法**（附带 owner_id 校验）：
```sql
SELECT * FROM novel_project WHERE id = ? AND owner_id = ?
```

## 4. SQL 注入防护

### 4.1 禁止字符串拼接

**错误写法**：
```go
query := "SELECT * FROM novel_project WHERE id = " + projectId
```

**正确写法**：
```go
db.Model("novel_project").Where("id = ? AND owner_id = ?", projectId, userId).One()
```

### 4.2 规则

1. 所有数据库查询必须使用 ORM 或参数化查询
2. 禁止字符串拼接 SQL
3. 所有 project_id 查询必须附带 owner_id 校验
4. 用户输入必须在应用层校验后再传入查询

## 5. 文件上传安全

### 5.1 文件类型限制

只允许以下扩展名：
- `.txt`
- `.md`
- `.docx`

### 5.2 文件大小限制

- 建议限制 10MB（可配置为 20MB）

### 5.3 MIME 类型检查

```go
allowedMIME := map[string]bool{
    "text/plain":              true,
    "text/markdown":           true,
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}
```

### 5.4 文件名安全

- 禁止使用原始文件名作为存储路径
- 使用 UUID 生成安全文件名
- 防止路径穿越（`../`）

**危险文件名示例**：
```
../../../../etc/passwd
evil.php
novel.txt<script>
```

**推荐存储路径**：
```
storage/uploads/{user_id}/{project_id}/{uuid}.txt
```

### 5.5 文件内容扫描

- 检查文件头部魔数（Magic Number）
- 扫描文件内容是否包含可执行代码
- 记录文件 SHA256 哈希用于去重和完整性校验

## 6. XSS 防护

### 6.1 前端输出转义

以下内容展示时必须转义：
- 小说原文
- AI 生成剧本
- YAML 内容
- Markdown 预览
- 错误提示

### 6.2 规则

1. 前端展示用户文本时必须 HTML 转义
2. Markdown 预览必须做 HTML sanitize
3. 禁止直接使用 `v-html` 渲染 AI 输出
4. YAML 编辑器只按纯文本展示
5. API 响应 Content-Type 设置为 `application/json`

## 7. 日志安全

### 7.1 禁止记录

- 用户密码
- JWT Token
- API Key
- 完整小说原文
- 完整 AI 响应
- 数据库连接密码
- 服务器绝对路径

### 7.2 允许记录

- request_id（请求追踪）
- user_id（用户标识）
- project_id（项目标识）
- task_id（任务标识）
- error_type（错误类型）
- 耗时（毫秒）
- HTTP 状态码

### 7.3 日志格式

```json
{
  "timestamp": "2024-01-01T00:00:00Z",
  "level": "error",
  "request_id": "req_abc123",
  "user_id": 1,
  "project_id": 42,
  "action": "script.generate",
  "error_type": "llm_timeout",
  "duration_ms": 30000,
  "message": "AI service request timed out"
}
```

## 8. API 安全

### 8.1 请求限流

- 登录接口：每 IP 每分钟最多 10 次
- 上传接口：每用户每小时最多 20 次
- 生成接口：每用户同时最多 3 个任务

### 8.2 请求 ID

每个请求必须生成唯一 request_id，用于日志追踪和问题排查。

### 8.3 安全响应头

```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
```

## 9. Docker 安全

### 9.1 容器安全

- 不在镜像中硬编码密码
- 使用 env_file 加载环境变量
- 不暴露不必要的端口到公网
- 上传目录挂载为 volume，不放入镜像

### 9.2 数据安全

- 不使用真实生产数据测试
- MySQL 数据使用 Docker volume 持久化
- 定期备份数据库（学习阶段可选）
