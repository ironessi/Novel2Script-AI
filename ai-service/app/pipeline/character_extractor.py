"""人物抽取 Pipeline"""

from app.core.llm_client import llm_client
from app.core.security import sanitize_novel_text
from app.core.logger import logger
from app.prompts.loader import load_prompt, get_system_guard


async def extract_characters(chapters: list[dict]) -> dict:
    """从小说章节中抽取人物信息

    Args:
        chapters: [{"chapter_index": 1, "title": "...", "content": "..."}]

    Returns:
        {"characters": [...]}
    """
    system_prompt = get_system_guard() + "\n\n" + load_prompt("character_extract")

    # 拼接章节内容
    text_parts = []
    for ch in chapters:
        text_parts.append(
            f"=== 第{ch['chapter_index']}章 {ch.get('title', '')} ===\n{ch['content']}"
        )
    novel_text = "\n\n".join(text_parts)
    novel_text = sanitize_novel_text(novel_text)

    # 限制输入长度（避免 token 超限）
    if len(novel_text) > 50000:
        novel_text = novel_text[:50000] + "\n\n[文本过长，已截断]"

    user_prompt = f"以下是小说内容：\n\n{novel_text}"

    logger.info(f"Extracting characters from {len(chapters)} chapters")

    result = await llm_client.chat_json(
        system_prompt=system_prompt,
        user_prompt=user_prompt,
        temperature=0.1,
    )

    # 确保返回格式正确
    if "characters" not in result:
        raise ValueError("人物抽取结果缺少 characters 字段")

    logger.info(f"Extracted {len(result['characters'])} characters")
    return result
