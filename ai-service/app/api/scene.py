"""场景拆分接口"""

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from app.pipeline.scene_planner import split_scenes

router = APIRouter(prefix="/ai", tags=["scene"])


class ChapterItem(BaseModel):
    chapter_index: int
    title: str = ""
    content: str


class SceneRequest(BaseModel):
    project_id: int
    chapters: list[ChapterItem]
    characters: list[dict]
    plot_events: list[dict]


@router.post("/split-scenes")
async def split_scenes_api(req: SceneRequest):
    """根据小说内容、人物和事件拆分场景"""
    try:
        chapters = [ch.model_dump() for ch in req.chapters]
        result = await split_scenes(chapters, req.characters, req.plot_events)
        return {"status": "completed", "data": result}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
