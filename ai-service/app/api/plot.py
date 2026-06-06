"""剧情事件链接口"""

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from app.pipeline.plot_event_builder import build_plot_events

router = APIRouter(prefix="/ai", tags=["plot"])


class ChapterItem(BaseModel):
    chapter_index: int
    title: str = ""
    content: str


class CharacterItem(BaseModel):
    id: str
    name: str
    role: str = ""
    description: str = ""


class PlotRequest(BaseModel):
    project_id: int
    chapters: list[ChapterItem]
    characters: list[CharacterItem]


@router.post("/build-plot-events")
async def build_plot_events_api(req: PlotRequest):
    """从章节内容中提取剧情事件链"""
    try:
        chapters = [ch.model_dump() for ch in req.chapters]
        characters = [c.model_dump() for c in req.characters]
        result = await build_plot_events(chapters, characters)
        return {"status": "completed", "data": result}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
