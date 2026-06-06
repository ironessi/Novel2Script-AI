"""YAML 自动修复 Pipeline"""

import json

from app.core.llm_client import llm_client
from app.core.logger import logger
from app.prompts.loader import load_prompt, get_system_guard

from app.pipeline.schema_validator import validate_yaml


async def repair_yaml(yaml_content: str, issues: list[dict]) -> dict:
    """尝试自动修复 YAML 中的问题

    Args:
        yaml_content: 原始 YAML 内容
        issues: 校验发现的问题列表

    Returns:
        {"fixed_yaml": "...", "changes": [...], "success": bool}
    """
    system_prompt = get_system_guard() + "\n\n" + load_prompt("yaml_repair")

    issues_info = json.dumps(issues, ensure_ascii=False, indent=2)

    user_prompt = (
        f"## 发现的问题\n{issues_info}\n\n"
        f"## 待修复的 YAML\n```yaml\n{yaml_content}\n```"
    )

    logger.info(f"Attempting YAML repair, {len(issues)} issues to fix")

    result = await llm_client.chat_json(
        system_prompt=system_prompt,
        user_prompt=user_prompt,
        temperature=0.1,
    )

    fixed_yaml = result.get("fixed_yaml", "")
    success = result.get("success", False)

    # 验证修复结果
    if fixed_yaml and success:
        valid, new_issues = validate_yaml(fixed_yaml)
        if valid:
            logger.info("YAML repair successful, validation passed")
        else:
            logger.warning(f"YAML repair partially successful, {len(new_issues)} issues remain")
            success = False
            result["remaining_issues"] = new_issues

    return result
