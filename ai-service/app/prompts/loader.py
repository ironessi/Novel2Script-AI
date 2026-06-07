"""Prompt 模板加载器"""

import os

_PROMPT_DIR = os.path.dirname(__file__)

_cache: dict[str, str] = {}


def load_prompt(name: str) -> str:
    """加载 prompt 模板文件"""
    if name in _cache:
        return _cache[name]

    path = os.path.join(_PROMPT_DIR, f"{name}.md")
    if not os.path.exists(path):
        raise FileNotFoundError(f"Prompt template not found: {name}")

    with open(path, "r", encoding="utf-8") as f:
        content = f.read()

    _cache[name] = content
    return content


def get_system_guard() -> str:
    return load_prompt("system_guard")
