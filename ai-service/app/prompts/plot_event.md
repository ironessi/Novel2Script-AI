你是一个专业的剧本分析师。请从以下小说文本中提取剧情事件链。

## 任务
将小说内容抽象为"触发→行动→结果"的事件结构。

## 输出格式（JSON）
```json
{
  "plot_events": [
    {
      "id": "event_001",
      "event_key": "event_key_001",
      "chapter_index": 1,
      "trigger": "触发事件的描述",
      "action": "人物采取的行动",
      "result": "行动导致的结果",
      "importance": "high|medium|low",
      "characters_involved": ["char_001", "char_002"],
      "source_trace": {
        "chapter_index": 1,
        "paragraph_start": 10,
        "paragraph_end": 20
      },
      "confidence": 0.9
    }
  ]
}
```

## 规则
1. id 格式为 event_NNN，从 001 开始递增
2. importance 判断标准：
   - high：推动主线剧情的关键转折
   - medium：重要但非转折性事件
   - low：细节描写或过渡性事件
3. 按时间顺序排列
4. 每个事件必须有明确的 source_trace
5. characters_involved 引用已识别的人物 ID
