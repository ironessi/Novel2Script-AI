"""幻觉检测 Pipeline"""

import json
import yaml

from app.core.llm_client import llm_client
from app.core.logger import logger
from app.prompts.loader import load_prompt, get_system_guard


async def check_hallucination(
    yaml_content: str, chapters: list[dict], characters: list[dict]
) -> dict:
    """检测 AI 生成的剧本是否存在幻觉

    Returns:
        {"risk_level": "low|medium|high", "issues": [...]}
    """
    system_prompt = get_system_guard() + "\n\n" + load_prompt("hallucination_check")

    # 解析剧本
    try:
        script_data = yaml.safe_load(yaml_content)
    except yaml.YAMLError:
        return {
            "risk_level": "high",
            "issues": [
                {
                    "type": "format_error",
                    "severity": "high",
                    "location": "root",
                    "message": "YAML 格式错误，无法解析",
                    "suggestion": "修复 YAML 语法错误",
                }
            ],
        }

    # 拼接原文
    text_parts = []
    for ch in chapters[:10]:  # 限制章节数量
        text_parts.append(
            f"=== 第{ch['chapter_index']}章 {ch.get('title', '')} ===\n{ch['content'][:3000]}"
        )
    novel_text = "\n\n".join(text_parts)

    # 格式化人物信息
    char_info = json.dumps(characters, ensure_ascii=False, indent=2)
    script_info = json.dumps(script_data, ensure_ascii=False, indent=2)

    if len(script_info) > 30000:
        script_info = script_info[:30000] + "\n\n[内容过长，已截断]"

    user_prompt = (
        f"## 人物档案（Character Bible）\n{char_info}\n\n"
        f"## 原文内容\n{novel_text}\n\n"
        f"## AI 生成的剧本\n{script_info}"
    )

    logger.info("Checking hallucination")

    result = await llm_client.chat_json(
        system_prompt=system_prompt,
        user_prompt=user_prompt,
        temperature=0.1,
    )

    risk_level = result.get("risk_level", "unknown")
    issue_count = len(result.get("issues", []))
    logger.info(f"Hallucination check: risk={risk_level}, issues={issue_count}")

    return result
