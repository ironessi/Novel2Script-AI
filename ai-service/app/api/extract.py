"""人物抽取接口"""

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from app.pipeline.character_extractor import extract_characters

router = APIRouter(prefix="/ai", tags=["extract"])


class ChapterItem(BaseModel):
    chapter_index: int
    title: str = ""
    content: str


class ExtractRequest(BaseModel):
    project_id: int
    chapters: list[ChapterItem]


@router.post("/extract-characters")
async def extract_characters_api(req: ExtractRequest):
    """从小说章节中抽取人物信息"""
    try:
        chapters = [ch.model_dump() for ch in req.chapters]
        result = await extract_characters(chapters)
        return {"status": "completed", "data": result}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
