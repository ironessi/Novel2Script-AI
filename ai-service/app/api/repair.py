"""YAML 修复接口"""

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from app.pipeline.yaml_repairer import repair_yaml

router = APIRouter(prefix="/ai", tags=["repair"])


class RepairRequest(BaseModel):
    yaml_content: str
    issues: list[dict]


@router.post("/repair-yaml")
async def repair_yaml_api(req: RepairRequest):
    """尝试自动修复 YAML 中的问题"""
    try:
        result = await repair_yaml(req.yaml_content, req.issues)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
