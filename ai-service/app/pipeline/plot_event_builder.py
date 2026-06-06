"""剧情事件链构建 Pipeline"""

from app.core.llm_client import llm_client
from app.core.security import sanitize_novel_text
from app.core.logger import logger
from app.prompts.loader import load_prompt, get_system_guard


async def build_plot_events(chapters: list[dict], characters: list[dict]) -> dict:
    """从小说章节中提取剧情事件链

    Args:
        chapters: [{"chapter_index": 1, "title": "...", "content": "..."}]
        characters: 人物档案列表

    Returns:
        {"plot_events": [...]}
    """
    system_prompt = get_system_guard() + "\n\n" + load_prompt("plot_event")

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

    # 格式化人物信息
    import json
    char_info = json.dumps(characters, ensure_ascii=False, indent=2)

    user_prompt = f"## 人物档案\n{char_info}\n\n## 小说内容\n{novel_text}"

    logger.info(f"Building plot events from {len(chapters)} chapters")

    result = await llm_client.chat_json(
        system_prompt=system_prompt,
        user_prompt=user_prompt,
        temperature=0.1,
    )

    if "plot_events" not in result:
        raise ValueError("剧情事件链结果缺少 plot_events 字段")

    logger.info(f"Built {len(result['plot_events'])} plot events")
    return result
