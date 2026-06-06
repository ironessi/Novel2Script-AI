# Novel2Script-AI

基于大模型的小说转剧本结构化生成与安全校验工具。

## 项目简介

Novel2Script-AI 帮助小说作者将小说文本自动转换为结构化 YAML 剧本。系统通过多阶段 AI 流水线，依次完成人物抽取、剧情事件链构建、场景拆分、剧本生成，并提供 Schema 校验、幻觉检测和安全审查能力。

### 核心创新点

1. **Character Bible 人物一致性机制** — 先抽取人物档案再生成剧本，防止人物改名/性格漂移/关系矛盾
2. **Plot Event Chain 剧情事件链** — 将小说抽象为"触发→行动→结果"结构，保证剧本不偏离主线
3. **Source Trace 原文溯源** — 每个场景记录来源章节和段落，降低幻觉、增强可解释性
4. **Schema-Guided Generation** — 通过 JSON Schema 约束 AI 输出格式，配合自动修复机制
5. **Hallucination Checker** — 自动检测未登记人物、未溯源事件、剧情冲突等幻觉问题

## 系统架构

```
浏览器 → Vue3 前端(5173) → GoFrame 后端(8000) → Python AI 服务(9000) → LLM
                            ↓                       ↓
                          MySQL                   Redis
                       (业务数据)              (任务状态/进度)
```

## 技术栈

| 模块 | 技术 | 说明 |
|------|------|------|
| 前端 | Vue 3 + Vite + Vue Router + Axios | 项目管理、上传、YAML 编辑、结果展示 |
| 后端 | Go + GoFrame v2 | 用户认证、项目管理、任务调度、权限控制 |
| AI 服务 | Python + FastAPI + OpenAI SDK | 多阶段剧本生成流水线、幻觉检测、安全审查 |
| 数据库 | MySQL 8.0 | 10 张业务表，BIGINT 主键，utf8mb4 |
| 缓存 | Redis 7 | 任务进度跟踪、会话缓存 |
| 容器化 | Docker Compose | MySQL + Redis 本地开发编排 |
| LLM | DeepSeek / OpenAI / Ollama | 支持远程 API 和本地模型 |

### 第三方依赖

**Go 后端：**
- `github.com/gogf/gf/v2` — GoFrame v2 Web 框架
- `github.com/golang-jwt/jwt/v5` — JWT 认证
- `github.com/redis/go-redis/v9` — Redis 客户端
- `golang.org/x/crypto` — bcrypt 密码哈希

**Python AI 服务：**
- `fastapi` — Web 框架
- `uvicorn` — ASGI 服务器
- `openai` — LLM API 客户端（兼容 DeepSeek/Ollama）
- `pyyaml` — YAML 解析
- `jsonschema` — JSON Schema 校验
- `httpx` — HTTP 客户端

**Vue3 前端：**
- `vue` + `vue-router` — 框架与路由
- `axios` — HTTP 请求
- `js-yaml` — YAML 解析

## 快速启动

### 1. 启动基础设施

```bash
cd deploy && docker compose up -d mysql redis
```

### 2. 初始化数据库

```bash
mysql -u root -p < deploy/mysql-init.sql
mysql -u root -p novel2script < scripts/seed.sql
```

### 3. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env，填写 MySQL 密码、JWT Secret、LLM API Key
```

### 4. 启动后端

```bash
cd backend && go run main.go
# 服务运行在 http://localhost:8000
```

### 5. 启动 AI 服务

```bash
cd ai-service && pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 9000
# 服务运行在 http://localhost:9000
# API 文档：http://localhost:9000/docs
```

### 6. 启动前端

```bash
cd frontend && npm install && npm run dev
# 服务运行在 http://localhost:5173
```

## 目录结构

```
Novel2Script-AI/
├── backend/                  # GoFrame 后端
│   ├── api/v1/              # 请求/响应结构体
│   ├── internal/
│   │   ├── cmd/             # 启动与路由注册
│   │   ├── controller/      # HTTP 控制器
│   │   ├── service/         # 服务接口
│   │   ├── logic/           # 业务逻辑
│   │   ├── dao/             # 数据库操作
│   │   ├── model/entity/    # 实体模型
│   │   ├── middleware/      # 中间件（JWT/RBAC/限流）
│   │   ├── client/          # AI 服务客户端
│   │   ├── redis/           # Redis 客户端
│   │   └── runner/          # 异步任务执行器
│   └── utility/             # 工具包（密码/JWT/文件校验）
├── ai-service/              # Python AI 服务
│   ├── app/
│   │   ├── api/             # 9 个 API 接口
│   │   ├── core/            # 配置/LLM客户端/安全/日志
│   │   ├── pipeline/        # 8 个处理模块
│   │   ├── prompts/         # 8 个 Prompt 模板
│   │   └── schemas/         # JSON Schema 校验文件
│   ├── requirements.txt
│   └── Dockerfile
├── frontend/                # Vue3 前端
│   ├── src/
│   │   ├── api/             # API 封装
│   │   ├── pages/           # 7 个页面
│   │   ├── layouts/         # 布局组件
│   │   ├── router/          # 路由配置
│   │   └── styles/          # 全局样式
│   └── package.json
├── deploy/                  # 部署配置
│   ├── docker-compose.yml   # MySQL + Redis
│   ├── mysql-init.sql       # 建表脚本
│   └── redis.conf           # Redis 配置
├── docs/                    # 设计文档
│   ├── ARCHITECTURE.md
│   ├── DATABASE.md
│   ├── SECURITY-DESIGN.md
│   ├── AI-SAFETY-DESIGN.md
│   ├── YAML-SCHEMA-DESIGN.md
│   └── PROGRESS.md
├── scripts/                 # 数据库脚本
│   ├── migrate.sql
│   └── seed.sql
├── examples/                # 示例文件
│   └── novel_sample.txt
└── .env.example             # 环境变量模板
```

## API 接口

### 后端接口 (GoFrame, :8000)

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/auth/register | 用户注册 |
| POST | /api/auth/login | 用户登录 |
| GET | /api/auth/me | 当前用户信息 |
| POST | /api/projects | 创建项目 |
| GET | /api/projects | 项目列表 |
| GET | /api/projects/:id | 项目详情 |
| PUT | /api/projects/:id | 更新项目 |
| DELETE | /api/projects/:id | 删除项目 |
| POST | /api/projects/:id/upload | 上传小说 |
| GET | /api/projects/:id/chapters | 章节列表 |
| POST | /api/projects/:id/generate | 创建 AI 任务 |
| GET | /api/tasks/:id/status | 任务状态 |
| GET | /api/projects/:id/script | 获取剧本 |
| PUT | /api/projects/:id/script | 编辑剧本 |
| POST | /api/projects/:id/validate | YAML 校验 |
| POST | /api/projects/:id/check-hallucination | 幻觉检测 |
| POST | /api/projects/:id/check-safety | 安全审查 |
| GET | /api/projects/:id/export | 导出剧本 |
| GET | /api/projects/:id/audit | 审计日志 |

### AI 服务接口 (FastAPI, :9000)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /health | 健康检查 |
| POST | /ai/analyze | 完整分析流水线 |
| POST | /ai/extract-characters | 人物抽取 |
| POST | /ai/build-plot-events | 剧情事件链 |
| POST | /ai/split-scenes | 场景拆分 |
| POST | /ai/generate-script | 剧本生成 |
| POST | /ai/validate-yaml | YAML 校验 |
| POST | /ai/check-hallucination | 幻觉检测 |
| POST | /ai/check-safety | 安全审查 |
| POST | /ai/repair-yaml | YAML 修复 |

## 原创功能说明

本项目的核心业务逻辑均为原创实现，包括：

1. **多阶段剧本生成流水线** — 自研的 8 阶段处理流程（人物抽取→事件链→场景拆分→剧本生成→校验→修复）
2. **Character Bible 机制** — 基于人物档案约束剧本生成的一致性方案
3. **Plot Event Chain 机制** — 触发-行动-结果的事件链抽象
4. **Source Trace 溯源机制** — 场景到原文段落的映射关系
5. **幻觉检测规则引擎** — 结合规则检查和 LLM 审查的双重检测
6. **Prompt Injection 防护** — 小说文本中的恶意指令过滤
7. **YAML 自动修复** — 校验失败后自动尝试修复并重新校验

## 端口规划

| 服务 | 端口 |
|------|------|
| Vue3 前端 | 5173 |
| GoFrame 后端 | 8000 |
| Python AI 服务 | 9000 |
| MySQL | 3306 |
| Redis | 6379 |

## License

MIT
