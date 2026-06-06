"""安全审查接口"""

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from app.pipeline.safety_checker import check_safety

router = APIRouter(prefix="/ai", tags=["safety"])


class SafetyRequest(BaseModel):
    yaml_content: str


@router.post("/check-safety")
async def check_safety_api(req: SafetyRequest):
    """对 AI 生成的剧本进行安全审查"""
    try:
        result = await check_safety(req.yaml_content)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
