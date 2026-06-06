"""场景拆分 Pipeline"""

import json

from app.core.llm_client import llm_client
from app.core.security import sanitize_novel_text
from app.core.logger import logger
from app.prompts.loader import load_prompt, get_system_guard


async def split_scenes(
    chapters: list[dict], characters: list[dict], plot_events: list[dict]
) -> dict:
    """根据小说内容、人物和事件拆分场景

    Args:
        chapters: 章节列表
        characters: 人物档案列表
        plot_events: 剧情事件列表

    Returns:
        {"scenes": [...]}
    """
    system_prompt = get_system_guard() + "\n\n" + load_prompt("scene_split")

    # 拼接章节内容
    text_parts = []
    for ch in chapters:
        text_parts.append(
            f"=== 第{ch['chapter_index']}章 {ch.get('title', '')} ===\n{ch['content']}"
        )
    novel_text = "\n\n".join(text_parts)
    novel_text = sanitize_novel_text(novel_text)

    if len(novel_text) > 50000:
        novel_text = novel_text[:50000] + "\n\n[文本过长，已截断]"

    char_info = json.dumps(characters, ensure_ascii=False, indent=2)
    event_info = json.dumps(plot_events, ensure_ascii=False, indent=2)

    user_prompt = (
        f"## 人物档案\n{char_info}\n\n"
        f"## 剧情事件链\n{event_info}\n\n"
        f"## 小说内容\n{novel_text}"
    )

    logger.info("Splitting scenes")

    result = await llm_client.chat_json(
        system_prompt=system_prompt,
        user_prompt=user_prompt,
        temperature=0.2,
    )

    if "scenes" not in result:
        raise ValueError("场景拆分结果缺少 scenes 字段")

    logger.info(f"Split into {len(result['scenes'])} scenes")
    return result
