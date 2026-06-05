# Novel2Script-AI YAML Schema 设计文档

## 1. 概述

本项目使用 YAML 作为剧本的结构化输出格式。所有 AI 生成的剧本必须符合预定义的 YAML Schema，通过校验后才能保存和导出。

## 2. 设计原则

1. **结构化**：剧本必须是机器可解析的结构化数据
2. **可溯源**：每个场景必须有原文来源（source_trace）
3. **可校验**：提供 JSON Schema 用于自动校验
4. **可编辑**：支持人工在线编辑
5. **可导出**：支持导出为 YAML 和 Markdown

## 3. 完整 YAML 结构

```yaml
script:
  metadata:
    title: 剧本标题
    source_title: 原著小说标题
    adaptation_mode: screen_script  # screen_script | stage_play | short_video | radio_drama
    language: zh-CN
    version: 1
    author: 作者名（可选）
    created_at: "2024-01-01T00:00:00Z"

  characters:
    - id: char_001
      name: 林舟
      aliases:
        - 阿舟
        - 小林
      role: protagonist  # protagonist | antagonist | supporting | minor
      description: 性格冷静的男主角，大学教授
      personality:
        - 冷静
        - 克制
        - 理性
      relationships:
        - target: char_002
          relation: 青梅竹马
        - target: char_003
          relation: 师生
      source_trace:
        - chapter_index: 1
          paragraph_start: 2
          paragraph_end: 5

  scenes:
    - id: scene_001
      title: 旧教室重逢
      order: 1
      time: 傍晚
      location: 旧教学楼
      source_trace:
        - chapter_index: 1
          paragraph_start: 12
          paragraph_end: 18
      characters:
        - char_001
        - char_002
      summary: 林舟在旧教室与苏晚重逢。
      actions:
        - character: char_001
          description: 林舟推开教室门，脚步停在门口。
        - character: char_002
          description: 苏晚坐在窗边，夕阳落在她的肩上。
      dialogues:
        - character: char_002
          line: 你终于来了。
          emotion: 平静中带着一丝等待
        - character: char_001
          line: 对不起，我来晚了。
          emotion: 愧疚
      adaptation_notes:
        - source: 第一章第 14 段
          change: 将心理描写改写为人物动作
          reason: 剧本需要通过可视化动作表达人物情绪

  adaptation_summary:
    total_chapters: 3
    total_scenes: 8
    main_conflict: 主角与旧友之间的误会逐渐揭开
    hallucination_risk: low  # low | medium | high
```

## 4. 字段说明

### 4.1 metadata

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | YES | 剧本标题 |
| source_title | string | YES | 原著小说标题 |
| adaptation_mode | enum | YES | 改编模式 |
| language | string | YES | 语言代码 |
| version | integer | YES | 版本号 |
| author | string | NO | 作者 |
| created_at | datetime | NO | 创建时间 |

### 4.2 characters[]

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | YES | 唯一标识符，格式 char_NNN |
| name | string | YES | 人物名称 |
| aliases | string[] | NO | 别名列表 |
| role | enum | YES | 角色类型 |
| description | string | YES | 人物描述 |
| personality | string[] | YES | 性格特征 |
| relationships | object[] | NO | 人物关系 |
| source_trace | object[] | YES | 原文溯源 |

### 4.3 scenes[]

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | YES | 唯一标识符，格式 scene_NNN |
| title | string | YES | 场景标题 |
| order | integer | YES | 场景顺序 |
| time | string | YES | 时间 |
| location | string | YES | 地点 |
| source_trace | object[] | YES | 原文溯源 |
| characters | string[] | YES | 场景中的人物 ID |
| summary | string | YES | 场景概述 |
| actions | object[] | NO | 动作描述 |
| dialogues | object[] | NO | 对白 |
| adaptation_notes | object[] | NO | 改编说明 |

### 4.4 dialogues[]

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| character | string | YES | 人物 ID，必须引用 characters.id |
| line | string | YES | 对白内容 |
| emotion | string | NO | 情绪/语气 |

### 4.5 actions[]

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| character | string | NO | 人物 ID（可选） |
| description | string | YES | 动作描述 |

### 4.6 source_trace[]

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| chapter_index | integer | YES | 章节序号 |
| paragraph_start | integer | YES | 起始段落 |
| paragraph_end | integer | YES | 结束段落 |

### 4.7 adaptation_notes[]

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| source | string | YES | 原文位置 |
| change | string | YES | 改动内容 |
| reason | string | YES | 改动原因 |

## 5. JSON Schema 校验规则

### 5.1 必填字段校验

- `script.metadata` 必须存在
- `script.characters` 至少 1 个
- `script.scenes` 至少 1 个
- 每个 scene 必须有 `id`、`title`、`time`、`location`、`source_trace`

### 5.2 引用完整性校验

- `dialogues.character` 必须引用 `characters.id`
- `scene.characters` 必须引用 `characters.id`
- `relationships.target` 必须引用 `characters.id`

### 5.3 格式校验

- `id` 格式：`char_NNN` 或 `scene_NNN`
- `source_trace.chapter_index` 必须大于 0
- `source_trace.paragraph_start` <= `paragraph_end`
- YAML 语法正确

### 5.4 唯一性校验

- `characters.id` 必须唯一
- `scenes.id` 必须唯一
- `scenes.order` 必须唯一

## 6. 导出格式

### 6.1 YAML 导出

直接导出完整的 YAML 文件。

### 6.2 Markdown 导出

```markdown
# 剧本标题

> 原著：原著小说标题
> 改编模式：影视剧本
> 版本：1

## 人物表

| ID | 姓名 | 角色 | 描述 |
|----|------|------|------|
| char_001 | 林舟 | 主角 | 性格冷静的男主角 |

---

## 场景 1：旧教室重逢

**时间**：傍晚
**地点**：旧教学楼
**人物**：林舟、苏晚

### 动作

- 林舟推开教室门，脚步停在门口。
- 苏晚坐在窗边，夕阳落在她的肩上。

### 对白

**苏晚**（平静中带着一丝等待）：你终于来了。

**林舟**（愧疚）：对不起，我来晚了。

---
```
