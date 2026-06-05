# Novel2Script-AI 架构设计文档

## 1. 系统架构概览

Novel2Script-AI 采用前后端分离的三层架构：

```
┌─────────────┐     ┌─────────────────┐     ┌──────────────────┐
│   Frontend   │────▶│  GoFrame Backend │────▶│ Python AI Service │
│  Vue3/React  │     │  REST API        │     │  FastAPI + LLM    │
└─────────────┘     └────────┬────────┘     └────────┬─────────┘
                             │                       │
                    ┌────────┴────────┐     ┌────────┴─────────┐
                    │     MySQL       │     │      Redis        │
                    │   业务数据       │     │   任务状态/缓存    │
                    └─────────────────┘     └──────────────────┘
```

## 2. 技术栈

| 层级 | 技术 | 用途 |
|------|------|------|
| 前端 | Vue3 + Vite | 项目管理、上传、YAML 编辑、结果展示 |
| 后端 | Go (GoFrame) | 用户认证、项目 CRUD、任务调度、权限控制、审计日志 |
| AI 服务 | Python + FastAPI | 文本处理、LLM 调用、人物抽取、剧本生成、校验 |
| 数据库 | MySQL 8.0 | 持久化存储业务数据 |
| 缓存 | Redis 7 | 任务状态、生成进度、会话缓存 |
| 文件存储 | 本地文件系统 | 上传小说文件、导出剧本文件 |
| 容器化 | Docker Compose | 本地开发环境编排 |

## 3. 模块职责

### 3.1 Frontend

- 用户登录/注册页面
- 项目列表与详情
- 小说文本上传
- 章节列表查看
- AI 生成任务进度展示
- YAML 剧本在线编辑器
- 校验问题面板
- 幻觉/安全风险报告
- 剧本导出（YAML / Markdown）
- 审计日志查看

### 3.2 GoFrame Backend

- **用户认证**: 注册、登录、JWT 签发与校验
- **项目管理**: 项目 CRUD、软删除
- **文件上传**: 接收小说文件、安全校验、存储管理
- **章节管理**: 章节入库、列表查询
- **任务调度**: 创建 AI 任务、状态流转、失败重试
- **剧本管理**: 剧本版本管理、在线编辑保存
- **权限控制**: 基于 owner_id 的项目级权限校验
- **审计日志**: 关键操作记录

### 3.3 Python AI Service

- **文本清洗**: 去除乱码、标准化空白字符
- **章节切分**: 基于正则和 LLM 的章节识别
- **人物抽取 (Character Bible)**: 提取人物档案
- **剧情事件链 (Plot Event Chain)**: 提取触发-行动-结果事件
- **场景拆分 (Scene Plan)**: 基于地点/时间/人物变化拆分场景
- **剧本生成**: 分场景生成 YAML 格式剧本
- **Schema 校验**: 验证 YAML 结构合规性
- **幻觉检测**: 检查 AI 生成内容是否脱离原文
- **安全审查**: 检查输出是否包含危险内容
- **YAML 修复**: 自动修复格式错误

## 4. 数据流

### 4.1 小说上传流程

```
用户上传 .txt/.md/.docx
  → 后端文件安全校验（类型、大小、路径）
  → 存储到 storage/uploads/{user_id}/{project_id}/{uuid}.ext
  → 记录文件元数据到 novel_source_file 表
  → 触发章节切分
  → 章节内容入库 novel_chapter 表
```

### 4.2 AI 生成流程

```
用户发起生成请求
  → 后端创建 ai_task 记录（status=pending）
  → 后端异步调用 Python AI Service
  → AI Service 多阶段流水线处理：
      1. 文本清洗
      2. 章节切分确认
      3. 人物抽取 → character_profile 表
      4. 剧情事件链构建 → plot_event 表
      5. 场景拆分
      6. 分场景剧本生成
      7. YAML Schema 校验
      8. 幻觉检测
      9. 安全审查
      10. 自动修复（如需要）
  → 生成结果存入 script_version 表
  → 校验问题存入 validation_issue 表
  → 更新 ai_task 状态
  → 前端轮询任务状态并展示结果
```

## 5. 服务间通信

- **Frontend → Backend**: HTTP REST API（JSON）
- **Backend → AI Service**: HTTP REST API（JSON），内部 Token 认证
- **Backend → MySQL**: GoFrame ORM
- **Backend → Redis**: go-redis
- **AI Service → LLM**: HTTP API（OpenAI 兼容格式）

## 6. 端口规划

| 服务 | 端口 |
|------|------|
| Frontend | 5173 |
| Backend | 8000 |
| AI Service | 9000 |
| MySQL | 3306 |
| Redis | 6379 |

## 7. 开发模式

### 本地开发

```bash
# 启动基础设施
docker compose up -d mysql redis

# 启动后端
cd backend && go run main.go

# 启动 AI 服务
cd ai-service && uvicorn app.main:app --host 0.0.0.0 --port 9000 --reload

# 启动前端
cd frontend && npm run dev
```

### Docker Compose 一键启动

```bash
docker compose up -d
```
