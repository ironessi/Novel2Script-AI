"""安全审查 Pipeline"""

import json
import yaml

from app.core.llm_client import llm_client
from app.core.security import check_output_safety
from app.core.logger import logger
from app.prompts.loader import load_prompt, get_system_guard


async def check_safety(yaml_content: str) -> dict:
    """对 AI 生成的剧本进行安全审查

    Returns:
        {"safe": bool, "risk_level": "...", "issues": [...]}
    """
    # 1. 规则检查（快速）
    rule_check = check_output_safety(yaml_content)
    if rule_check:
        return {
            "safe": False,
            "risk_level": "high",
            "issues": [
                {
                    "type": "info_leak",
                    "severity": "high",
                    "location": "script",
                    "message": rule_check,
                    "suggestion": "删除敏感信息后重新生成",
                }
            ],
        }

    # 2. LLM 深度审查
    system_prompt = get_system_guard() + "\n\n" + load_prompt("safety_check")

    try:
        script_data = yaml.safe_load(yaml_content)
        script_info = json.dumps(script_data, ensure_ascii=False, indent=2)
    except yaml.YAMLError:
        script_info = yaml_content

    if len(script_info) > 30000:
        script_info = script_info[:30000] + "\n\n[内容过长，已截断]"

    user_prompt = f"## 待审查的剧本内容\n{script_info}"

    logger.info("Running safety check")

    result = await llm_client.chat_json(
        system_prompt=system_prompt,
        user_prompt=user_prompt,
        temperature=0.1,
    )

    safe = result.get("safe", True)
    risk_level = result.get("risk_level", "unknown")
    issue_count = len(result.get("issues", []))
    logger.info(f"Safety check: safe={safe}, risk={risk_level}, issues={issue_count}")

    return result
