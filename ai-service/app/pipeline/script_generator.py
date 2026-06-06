"""剧本生成 Pipeline"""

import json
import yaml

from app.core.llm_client import llm_client
from app.core.logger import logger
from app.prompts.loader import load_prompt, get_system_guard


async def generate_script(
    chapters: list[dict],
    characters: list[dict],
    plot_events: list[dict],
    scenes: list[dict],
    metadata: dict,
) -> str:
    """为每个场景生成剧本内容，组合为完整 YAML 剧本

    Returns:
        完整的 YAML 剧本字符串
    """
    system_prompt = get_system_guard() + "\n\n" + load_prompt("script_generate")

    script_scenes = []

    for scene in scenes:
        # 准备场景相关的章节内容
        scene_text_parts = []
        for trace in scene.get("source_trace", []):
            ch_idx = trace.get("chapter_index", 0)
            for ch in chapters:
                if ch["chapter_index"] == ch_idx:
                    paragraphs = ch["content"].split("\n")
                    start = trace.get("paragraph_start", 1) - 1
                    end = trace.get("paragraph_end", len(paragraphs))
                    scene_text = "\n".join(paragraphs[start:end])
                    scene_text_parts.append(scene_text)
                    break

        scene_text = "\n".join(scene_text_parts)

        # 获取场景相关的人物信息
        scene_chars = [
            c for c in characters if c["id"] in scene.get("characters", [])
        ]

        # 获取场景相关的事件
        scene_events = [
            e for e in plot_events if e["id"] in scene.get("related_events", [])
        ]

        char_info = json.dumps(scene_chars, ensure_ascii=False, indent=2)
        event_info = json.dumps(scene_events, ensure_ascii=False, indent=2)

        user_prompt = (
            f"## 场景信息\n"
            f"- 标题：{scene['title']}\n"
            f"- 时间：{scene.get('time', '未指定')}\n"
            f"- 地点：{scene.get('location', '未指定')}\n\n"
            f"## 相关人物\n{char_info}\n\n"
            f"## 相关事件\n{event_info}\n\n"
            f"## 原文内容\n{scene_text}"
        )

        logger.info(f"Generating script for scene: {scene['id']}")

        try:
            result = await llm_client.chat_json(
                system_prompt=system_prompt,
                user_prompt=user_prompt,
                temperature=0.3,
            )

            # 合并场景元数据和生成内容
            script_scene = {
                "id": scene["id"],
                "title": scene["title"],
                "order": scene.get("order", 0),
                "time": scene.get("time", ""),
                "location": scene.get("location", ""),
                "source_trace": scene.get("source_trace", []),
                "characters": scene.get("characters", []),
                "summary": scene.get("summary", ""),
                "actions": result.get("actions", []),
                "dialogues": result.get("dialogues", []),
                "adaptation_notes": result.get("adaptation_notes", []),
            }
            script_scenes.append(script_scene)

        except Exception as e:
            logger.error(f"Failed to generate script for scene {scene['id']}: {e}")
            # 使用场景摘要作为降级方案
            script_scenes.append(
                {
                    "id": scene["id"],
                    "title": scene["title"],
                    "order": scene.get("order", 0),
                    "time": scene.get("time", ""),
                    "location": scene.get("location", ""),
                    "source_trace": scene.get("source_trace", []),
                    "characters": scene.get("characters", []),
                    "summary": scene.get("summary", ""),
                    "actions": [],
                    "dialogues": [],
                    "adaptation_notes": [],
                }
            )

    # 组装完整剧本
    script = {
        "script": {
            "metadata": metadata,
            "characters": characters,
            "scenes": script_scenes,
            "adaptation_summary": {
                "total_chapters": len(chapters),
                "total_scenes": len(script_scenes),
                "main_conflict": "",
                "hallucination_risk": "unknown",
            },
        }
    }

    # 转换为 YAML
    yaml_content = yaml.dump(script, allow_unicode=True, default_flow_style=False, sort_keys=False)

    logger.info(f"Generated script with {len(script_scenes)} scenes")
    return yaml_content
