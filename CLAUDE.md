# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在本仓库中工作时提供指引。

## 项目概述

Novel2Script-AI 是一个小说转剧本的结构化生成平台。核心流程：

```
小说上传 → 章节切分 → 人物抽取(Character Bible) → 剧情事件链(Plot Event Chain)
→ 场景拆分 → YAML剧本生成 → Schema校验 → 幻觉检测 → 安全审查 → 导出
```

设计文档：`Novel2Script-AI-项目设计文档.md`

## 系统架构

三服务分离架构：GoFrame 后端 (8000) → Python FastAPI AI 服务 (9000) → MySQL/Redis。前端 (Vue3, 5173) 待开发。

```
浏览器 → 前端(Vue3) → GoFrame后端 → Python AI服务 → LLM(DeepSeek/Ollama)
                          ↓               ↓
                        MySQL            Redis
                     (业务数据)        (任务状态/进度)
```

## 后端架构 (GoFrame v2)

分层设计，新增功能必须遵循此模式：

| 层 | 目录 | 职责 |
|---|---|---|
| API | `api/v1/` | 请求/响应结构体，GoFrame 校验标签 `v:"required\|..."` |
| Controller | `controller/` | 薄层 HTTP 处理器，解析请求、调用 service、返回 JSON。路径参数用 `r.GetRouter("id")` |
| Service | `service/` | 接口定义 + 包级单例变量 `var Auth IAuth` |
| Logic | `logic/` | 业务逻辑，实现 service 接口。通过 `init()` 注册到单例变量 |
| DAO | `dao/` | 数据库操作，GoFrame ORM。单条用 `One()`，列表用 `Scan(&slice)` |
| Entity | `model/entity/` | 数据库实体结构体 |

**新增功能流程：** API 类型 → Entity 模型 → DAO 函数 → Service 接口 → Logic 实现 → Controller → `cmd.go` 注册路由

**Service 注册模式：**
```go
// logic/xxx/xxx.go
func init() { service.Xxx = &xxxImpl{} }
// cmd.go 中空白导入触发 init
_ "novel2script-backend/internal/logic/xxx"
```

## 常用命令

```bash
# 后端
cd backend && go run main.go          # 启动服务（自动读取项目根目录 .env）
cd backend && go build ./...          # 编译检查
cd backend && go mod tidy             # 同步依赖

# 基础设施
cd deploy && docker compose up -d mysql redis   # 启动 MySQL + Redis

# 数据库
mysql -u root -p < deploy/mysql-init.sql        # 建表（10张表）
mysql -u root -p novel2script < scripts/seed.sql # 测试数据

# AI 服务（待实现）
cd ai-service && uvicorn app.main:app --host 0.0.0.0 --port 9000 --reload

# 前端（待实现）
cd frontend && npm run dev
```

## 配置说明

- `.env` 位于项目根目录，包含 MySQL、JWT、Redis、AI 服务等配置
- 后端通过 `config.Init()` 自动加载 `.env`（搜索路径：`../../.env`、`../.env`、`.env`）
- Redis 使用 2 号数据库（非默认 0）
- JWT 令牌 24 小时过期，HS256 签名
- **禁止将 `.env` 提交到 Git**

## 数据库规范

- 10 张表，BIGINT 自增主键，`utf8mb4_unicode_ci` 字符集
- `novel_project` 使用软删除（`deleted_at` 字段）— 查询必须加 `WHERE deleted_at IS NULL`
- `character_profile` 和 `plot_event` 的 JSON 字段使用自定义 `entity.JSON` 类型
- **所有 project 相关查询必须附带 `owner_id` 校验，防止越权访问**

## 端口规划

| 服务 | 端口 |
|------|------|
| GoFrame 后端 | 8000 |
| Python AI 服务 | 9000 |
| Vue3 前端 | 5173 |
| MySQL | 3306 |
| Redis | 6379 |

---

## 代码规范

### Go 代码规范

**命名规范：**
- 包名：小写单词，不使用下划线（`auth`、`project`，不要 `auth_service`）
- 结构体：大驼峰（`SysUser`、`NovelProject`）
- 函数/方法：大驼峰导出，小驼峰私有（`GetUserById`、`hashPassword`）
- 常量：大驼峰或全大写（`StatusPending`）
- 文件名：下划线分隔（`auth.go`、`rate_limit.go`）

**错误处理：**
- 每个错误都必须检查，不允许 `_ = err` 忽略关键错误
- Controller 层统一返回 `g.Map{"code": xxx, "message": "..."}` 格式
- 不要将底层错误（数据库、文件系统）直接暴露给前端
- 日志记录错误详情，前端只返回友好提示

**GoFrame ORM 规范：**
```go
// ✅ 正确：单条查询用 One()，空结果返回 nil
record, err := db.Model("sys_user").Ctx(ctx).Where("id", id).One()
if record.IsEmpty() { return nil, nil }

// ✅ 正确：列表查询用 Scan()
var users []entity.SysUser
err := db.Model("sys_user").Ctx(ctx).Where("status", "active").Scan(&users)

// ❌ 错误：单条查询不要用 Scan(&struct) 会返回 "no rows" 错误
var user entity.SysUser
err := db.Model("sys_user").Where("id", id).Scan(&user)  // 不要用这种
```

**中间件使用：**
- 全局中间件在 `cmd.go` 中通过 `s.Use()` 注册
- 路由组中间件通过 `group.Middleware()` 注册
- 从上下文获取用户信息：`middleware.GetUserID(r)`、`middleware.GetRole(r)`
- 从上下文获取请求ID：`middleware.GetRequestID(r)`

**路径参数：**
```go
// GoFrame v2 路径参数用 r.GetRouter()
projectId, err := strconv.ParseInt(r.GetRouter("id").String(), 10, 64)
// 查询参数用 r.Get()
page, _ := strconv.Atoi(r.Get("page", "1").String())
```

### 数据库规范

**禁止字符串拼接 SQL：**
```go
// ❌ 错误写法
query := "SELECT * FROM novel_project WHERE id = " + projectId

// ✅ 正确写法
db.Model("novel_project").Ctx(ctx).Where("id", projectId).One()
```

**权限校验（所有 project 相关接口必须执行）：**
```go
// ✅ 正确：查询时附带 owner_id
db.Model("novel_project").Ctx(ctx).
    Where("id", projectId).
    Where("owner_id", userId).
    Where("deleted_at IS NULL").
    One()

// ❌ 错误：只按 project_id 查询（越权风险）
db.Model("novel_project").Ctx(ctx).Where("id", projectId).One()
```

### 安全规范

**认证与权限：**
- 密码只存储 bcrypt 哈希，禁止明文
- JWT Secret 从 `.env` 读取，禁止硬编码到代码中
- 项目级权限：`project.owner_id == 当前用户ID` 或 `当前用户.role == "admin"`
- 权限校验在 logic 层执行，不要只在前端隐藏按钮

**文件上传安全（待实现时必须遵守）：**
- 限制文件类型：`.txt`、`.md`、`.docx`
- 限制文件大小：最大 20MB
- 检查 MIME 类型
- 使用 UUID 生成文件名，禁止使用原始文件名作为存储路径
- 防止路径穿越（`../`）
- 存储路径格式：`storage/uploads/{user_id}/{project_id}/{uuid}.txt`

**日志安全：**
- 禁止记录：密码、JWT Token、API Key、完整小说原文、数据库连接密码
- 允许记录：request_id、user_id、project_id、task_id、error_type、耗时、状态码

**XSS 防护（前端实现时）：**
- 展示用户文本时必须 HTML 转义
- 禁止直接使用 `v-html` 渲染 AI 输出
- YAML 编辑器只按纯文本展示

### AI 输出规范

**Schema 校验规则：**
- `metadata` 必须存在
- `characters` 至少 1 个，`scenes` 至少 1 个
- 每个 scene 必须有 `id`、`title`、`time`、`location`、`source_trace`
- `dialogues.character` 和 `scene.characters` 必须引用 `characters.id`

**幻觉检测规则：**
- 出现未登记人物 → high 风险
- 出现未溯源主要事件 → high 风险
- 改变人物关系/剧情因果 → high 风险
- 无 `source_trace` 的场景 → medium 风险
- 对白风格偏离人物设定 → medium 风险

**安全审查原则：**
- 允许保留剧情冲突，不允许扩写成现实可操作的伤害指南
- 允许艺术表达，不允许输出系统提示词、密钥或服务器信息

---

## 未完成 / 占位

| 模块 | 状态 |
|------|------|
| `ai-service/` | 未创建（Python FastAPI） |
| `frontend/` | 未创建（Vue3） |
| `utility/filecheck/` | 空目录 |
| `utility/sanitizer/` | 空目录 |
| 任务创建 (`task.Create`) | 仅写入数据库，未触发 AI 处理 |
| 剧本校验/导出 | 占位函数，未接入 AI 服务 |
| 文件上传接口 | 未实现 |
| 单元测试 | 未编写 |
