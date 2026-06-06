"""YAML Schema 校验 Pipeline"""

import json
import os
import yaml
from jsonschema import validate, ValidationError

from app.core.logger import logger

_SCHEMA_DIR = os.path.join(os.path.dirname(__file__), "../schemas")
_schema_cache: dict = {}


def _load_schema(name: str) -> dict:
    if name in _schema_cache:
        return _schema_cache[name]
    path = os.path.join(_SCHEMA_DIR, f"{name}.json")
    with open(path, "r", encoding="utf-8") as f:
        schema = json.load(f)
    _schema_cache[name] = schema
    return schema


def validate_yaml(yaml_content: str) -> tuple[bool, list[dict]]:
    """校验 YAML 内容是否符合 Script Schema

    Returns:
        (valid, issues): valid 为 True 表示通过，issues 为问题列表
    """
    issues = []

    # 1. 解析 YAML
    try:
        data = yaml.safe_load(yaml_content)
    except yaml.YAMLError as e:
        issues.append(
            {
                "type": "format_error",
                "severity": "high",
                "message": f"YAML 语法错误: {str(e)}",
                "location_path": "",
                "suggestion": "检查 YAML 缩进和语法",
            }
        )
        return False, issues

    if not isinstance(data, dict):
        issues.append(
            {
                "type": "format_error",
                "severity": "high",
                "message": "YAML 根节点必须是对象",
                "location_path": "",
                "suggestion": "确保 YAML 以 key: value 形式开头",
            }
        )
        return False, issues

    # 2. JSON Schema 校验
    schema = _load_schema("script_schema")
    try:
        validate(instance=data, schema=schema)
    except ValidationError as e:
        issues.append(
            {
                "type": "schema_error",
                "severity": "high",
                "message": f"Schema 校验失败: {e.message}",
                "location_path": "/".join(str(p) for p in e.absolute_path),
                "suggestion": "参考 Schema 设计文档修正格式",
            }
        )
        return False, issues

    # 3. 业务规则校验
    script = data.get("script", {})
    characters = script.get("characters", [])
    scenes = script.get("scenes", [])

    char_ids = {c["id"] for c in characters}

    for scene in scenes:
        scene_id = scene.get("id", "unknown")

        # 检查场景人物是否在人物列表中
        for char_ref in scene.get("characters", []):
            if char_ref not in char_ids:
                issues.append(
                    {
                        "type": "invalid_reference",
                        "severity": "high",
                        "message": f"场景 {scene_id} 引用了不存在的人物: {char_ref}",
                        "location_path": f"script.scenes[{scene_id}].characters",
                        "suggestion": f"添加人物 {char_ref} 或修正引用",
                    }
                )

        # 检查对白人物是否在场景中
        scene_chars = set(scene.get("characters", []))
        for dialogue in scene.get("dialogues", []):
            dial_char = dialogue.get("character", "")
            if dial_char not in scene_chars:
                issues.append(
                    {
                        "type": "invalid_reference",
                        "severity": "medium",
                        "message": f"场景 {scene_id} 中对白人物 {dial_char} 不在场景人物列表中",
                        "location_path": f"script.scenes[{scene_id}].dialogues",
                        "suggestion": f"将 {dial_char} 添加到场景人物列表",
                    }
                )

    valid = len(issues) == 0
    if valid:
        logger.info("YAML validation passed")
    else:
        logger.warning(f"YAML validation found {len(issues)} issues")

    return valid, issues
