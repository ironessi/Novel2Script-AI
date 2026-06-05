# Novel2Script-AI 项目路线与进度文档

> 最后更新：2026-06-05

## 总体进度

| 步骤 | 内容 | 状态 | 完成日期 |
|------|------|------|----------|
| 第 0 步 | 文档 + 数据库脚本 | ✅ 已完成 | 2026-06-05 |
| 第 1 步 | GoFrame 后端基础 | ✅ 已完成 | 2026-06-05 |
| 第 2 步 | 文件上传与章节切分 | ✅ 已完成 | 2026-06-05 |
| 第 3 步 | Python AI 服务 | ⏳ 待开始 | - |
| 第 4 步 | 任务系统（接入 AI） | ⏳ 待开始 | - |
| 第 5 步 | 安全与幻觉检测 | ⏳ 待开始 | - |
| 第 6 步 | 前端页面 | ⏳ 待开始 | - |

---

## 第 0 步：文档 + 数据库脚本

**完成内容：**
- [x] `docs/ARCHITECTURE.md` — 系统架构设计
- [x] `docs/DATABASE.md` — 数据库设计文档（10 张表）
- [x] `docs/SECURITY-DESIGN.md` — 安全设计文档
- [x] `docs/AI-SAFETY-DESIGN.md` — AI 安全与幻觉控制
- [x] `docs/YAML-SCHEMA-DESIGN.md` — YAML Schema 设计
- [x] `deploy/mysql-init.sql` — 建表 SQL（10 张表）
- [x] `scripts/migrate.sql` — 迁移脚本
- [x] `scripts/seed.sql` — 测试数据
- [x] `deploy/docker-compose.yml` — MySQL + Redis 容器编排
- [x] `.env.example` / `.gitignore` — 配置模板

**遇到的问题：** 无

---

## 第 1 步：GoFrame 后端基础

**完成内容：**
- [x] GoFrame 项目初始化（go.mod、main.go、目录结构）
- [x] 用户注册/登录接口（bcrypt 密码哈希 + JWT 认证）
- [x] JWT 中间件（HS256、24 小时过期）
- [x] RBAC 权限中间件（admin/user 角色）
- [x] CORS 跨域中间件
- [x] 请求 ID 中间件（X-Request-Id）
- [x] 安全响应头中间件
- [x] 登录限流中间件（10 次/分钟）
- [x] 项目 CRUD 接口（创建/列表/详情/更新/删除）
- [x] 章节管理接口（列表/详情）
- [x] 任务管理接口（创建/状态查询）
- [x] 剧本管理接口（获取/编辑/校验/导出）
- [x] 审计日志接口
- [x] AI Service HTTP 客户端
- [x] 配置管理（.env 文件加载）
- [x] 数据库连接（GoFrame ORM + MySQL driver）
- [x] 完整 API 路由注册（20 个接口）

**遇到的问题：**

| 问题 | 原因 | 解决方案 |
|------|------|----------|
| `ghttp.Map` 未定义 | GoFrame v2 没有 `ghttp.Map` 类型 | 改用 `g.Map` |
| `g.DB().SetConfig` 不存在 | GoFrame v2 ORM 配置方式不同 | 改用 `gdb.New(node)` 创建实例 |
| `g.DB().GetCore().Ping` 不存在 | GoFrame v2 Core 没有 Ping 方法 | 改用 `db.GetOne(ctx, "SELECT 1")` 测试连接 |
| `s.GET()` 未定义 | GoFrame v2 Server 没有直接的 GET 方法 | 改用 `s.BindHandler("GET:/health", handler)` |
| `s.SetPort` 类型错误 | SetPort 接受 int，配置是 string | 用 `strconv.Atoi` 转换 |
| 数据库连接格式错误 | GoFrame 链接格式需要 `mysql:` 前缀 | 改为 `mysql:user:pass@tcp(host:port)/db` |
| 数据库驱动未导入 | GoFrame v2 驱动是独立的 contrib 包 | 导入 `github.com/gogf/gf/contrib/drivers/mysql/v2` |
| `Scan(&user)` 返回 "no rows" 错误 | GoFrame Scan 对空结果返回 error | 改用 `One()` + `record.IsEmpty()` + `record.Struct()` |
| `.env` 文件未加载 | GoFrame 不自动加载 .env | 自己实现 `loadEnvFile()` 函数 |
| `goproxy.cn` 下载超时 | 国内代理不稳定 | 改用 `GOPROXY=https://proxy.golang.org,direct` |

---

## 第 2 步：文件上传与章节切分

**完成内容：**
- [x] 文件上传接口 `POST /api/projects/:id/upload`
- [x] 文件校验工具 `utility/filecheck/`
  - 扩展名校验（.txt/.md/.docx）
  - MIME 类型校验
  - 文件大小限制（20MB）
  - 文件名清洗（防路径穿越）
  - 脚本注入检测
- [x] 文本清洗工具 `utility/sanitizer/`
  - 控制字符去除
  - 连续空白/换行规范化
  - 章节标题正则识别（中文/英文）
  - 章节标题提取
- [x] 文件存储（UUID 重命名 + SHA256 哈希）
- [x] 章节切分逻辑（正则识别章节标题）
- [x] DAO 层扩展（CreateSourceFile、BatchCreateChapters）
- [x] 上传响应返回章节列表

**测试结果：**
```
上传 novel_sample.txt (1252 字节)
→ 正确切分为 3 章：
  - 第一章 重逢
  - 第二章 暗流
  - 第三章 雨夜
```

**遇到的问题：**

| 问题 | 原因 | 解决方案 |
|------|------|----------|
| IDE 误报 `io` import 未使用 | IDE 缓存过期 | 实际编译通过，忽略 IDE 诊断 |
| `go build ./...` 匹配不到包 | 在错误的目录执行 | 确保在 `backend/` 目录下执行 |
| 上传文件存储路径在 backend/ 下 | StoragePath 是相对路径 `./storage` | 符合预期，从 backend 目录运行 |

---

## 第 3 步：Python AI 服务（待开始）

**计划内容：**
- [ ] FastAPI 项目结构（`ai-service/`）
- [ ] LLM Client（支持 DeepSeek / Ollama）
- [ ] Prompt 模板管理（`prompts/` 目录）
- [ ] 人物抽取接口 `POST /ai/extract-characters`
- [ ] 剧情事件链接口 `POST /ai/build-plot-events`
- [ ] 场景拆分接口 `POST /ai/split-scenes`
- [ ] 剧本生成接口 `POST /ai/generate-script`
- [ ] YAML 校验接口 `POST /ai/validate-yaml`
- [ ] 幻觉检测接口 `POST /ai/check-hallucination`
- [ ] 安全审查接口 `POST /ai/check-safety`
- [ ] YAML 修复接口 `POST /ai/repair-yaml`
- [ ] Dockerfile

**技术选型：**
- Python 3.11+ / FastAPI
- PyYAML + jsonschema（YAML 校验）
- OpenAI SDK（兼容 DeepSeek / Ollama）

---

## 第 4 步：任务系统（待开始）

**计划内容：**
- [ ] Go 后端创建 AI 任务时实际调用 AI Service
- [ ] Redis 保存任务进度
- [ ] 异步任务处理（goroutine）
- [ ] 失败重试机制
- [ ] 任务状态实时更新
- [ ] 前端轮询任务进度

---

## 第 5 步：安全与幻觉检测（待开始）

**计划内容：**
- [ ] Schema Validator（JSON Schema 校验）
- [ ] Hallucination Checker（幻觉检测）
- [ ] Safety Checker（安全审查）
- [ ] YAML Repairer（自动修复）
- [ ] Validation Issue 入库
- [ ] 风险报告展示

---

## 第 6 步：前端页面（待开始）

**计划内容：**
- [ ] Vue3 + Vite 项目初始化
- [ ] 登录/注册页面
- [ ] 项目列表页面
- [ ] 文件上传页面
- [ ] 任务进度页面
- [ ] YAML 在线编辑器
- [ ] 校验问题面板
- [ ] 导出功能（YAML / Markdown）
- [ ] 审计日志页面

---

## 技术债务与已知问题

| 编号 | 问题 | 优先级 | 状态 |
|------|------|--------|------|
| TD-1 | `utility/filecheck/` 缺少 docx MIME 检测 | 低 | 待处理 |
| TD-2 | `utility/sanitizer/` docx 解析是简化版 | 中 | 待处理 |
| TD-3 | 任务创建未触发 AI 处理（第 4 步解决） | 高 | 待处理 |
| TD-4 | 剧本校验/导出是占位函数 | 中 | 待处理 |
| TD-5 | 没有单元测试 | 中 | 待处理 |
| TD-6 | go.mod 中有未使用的 websocket 依赖 | 低 | 待处理 |
| TD-7 | `gorilla/websocket` 可以从 go.mod 移除 | 低 | 待处理 |

---

## Git 提交记录

```
d43e5ed feat: 文件上传与章节切分 - 第2步
1e23be5 docs: CLAUDE.md 改为中文，补充完整代码规范
8d18770 docs: add CLAUDE.md for Claude Code guidance
698ca0a feat: GoFrame 后端基础搭建 - 第1步
d502b75 init: 项目初始化 - 第0步文档 + 数据库脚本
```

---

## 端口规划

| 服务 | 端口 | 状态 |
|------|------|------|
| GoFrame 后端 | 8000 | ✅ 运行中 |
| Python AI 服务 | 9000 | ⏳ 待开发 |
| Vue3 前端 | 5173 | ⏳ 待开发 |
| MySQL | 3306 | ✅ 运行中 |
| Redis | 6379 | ✅ 运行中 |
