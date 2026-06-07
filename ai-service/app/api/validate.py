"""YAML 校验接口"""

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from app.pipeline.schema_validator import validate_yaml

router = APIRouter(prefix="/ai", tags=["validate"])


class ValidateRequest(BaseModel):
    yaml_content: str


@router.post("/validate-yaml")
async def validate_yaml_api(req: ValidateRequest):
    """校验 YAML 内容是否符合 Script Schema"""
    try:
        valid, issues = validate_yaml(req.yaml_content)
        return {"valid": valid, "issues": issues}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
