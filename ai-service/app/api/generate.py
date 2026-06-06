"""剧本生成接口"""

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from app.pipeline.script_generator import generate_script

router = APIRouter(prefix="/ai", tags=["generate"])


class ChapterItem(BaseModel):
    chapter_index: int
    title: str = ""
    content: str


class GenerateRequest(BaseModel):
    project_id: int
    chapters: list[ChapterItem]
    characters: list[dict]
    plot_events: list[dict]
    scenes: list[dict]
    metadata: dict


@router.post("/generate-script")
async def generate_script_api(req: GenerateRequest):
    """为每个场景生成剧本内容"""
    try:
        chapters = [ch.model_dump() for ch in req.chapters]
        yaml_content = await generate_script(
            chapters, req.characters, req.plot_events, req.scenes, req.metadata
        )
        return {"status": "completed", "yaml_content": yaml_content}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
