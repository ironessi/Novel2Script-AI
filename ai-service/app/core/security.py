"""输入/输出安全处理"""

import re
from typing import Optional


def sanitize_novel_text(text: str) -> str:
    """清洗小说文本，防止 Prompt Injection"""
    if not text:
        return ""

    # 移除可能的指令注入模式
    injection_patterns = [
        r"(?i)ignore\s+(all\s+)?previous\s+instructions",
        r"(?i)忽略.*之前.*指令",
        r"(?i)你现在是.*管理员",
        r"(?i)delete\s+all\s+files",
        r"(?i)输出.*密码",
        r"(?i)输出.*API\s*Key",
        r"(?i)system\s*prompt",
    ]
    for pattern in injection_patterns:
        text = re.sub(pattern, "[已过滤]", text)

    return text.strip()


def check_output_safety(content: str) -> Optional[str]:
    """检查 AI 输出是否包含敏感信息"""
    if not content:
        return None

    sensitive_patterns = [
        (r"(?i)api[_-]?key\s*[:=]\s*\S+", "包含 API Key"),
        (r"(?i)password\s*[:=]\s*\S+", "包含密码"),
        (r"(?i)secret\s*[:=]\s*\S+", "包含密钥"),
        (r"/etc/passwd", "包含系统路径"),
        (r"(?i)DROP\s+TABLE", "包含 SQL 指令"),
        (r"(?i)DELETE\s+FROM", "包含 SQL 指令"),
    ]

    for pattern, message in sensitive_patterns:
        if re.search(pattern, content):
            return message

    return None
