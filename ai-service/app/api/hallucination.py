"""幻觉检测接口"""

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from app.pipeline.hallucination_checker import check_hallucination

router = APIRouter(prefix="/ai", tags=["hallucination"])


class ChapterItem(BaseModel):
    chapter_index: int
    title: str = ""
    content: str


class HallucinationRequest(BaseModel):
    yaml_content: str
    chapters: list[ChapterItem]
    characters: list[dict]


@router.post("/check-hallucination")
async def check_hallucination_api(req: HallucinationRequest):
    """检测 AI 生成的剧本是否存在幻觉"""
    try:
        chapters = [ch.model_dump() for ch in req.chapters]
        result = await check_hallucination(req.yaml_content, chapters, req.characters)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
