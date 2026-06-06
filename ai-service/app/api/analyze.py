"""完整分析接口 - 多阶段流水线"""

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from app.core.logger import logger
from app.pipeline.character_extractor import extract_characters
from app.pipeline.plot_event_builder import build_plot_events
from app.pipeline.scene_planner import split_scenes
from app.pipeline.script_generator import generate_script
from app.pipeline.schema_validator import validate_yaml
from app.pipeline.hallucination_checker import check_hallucination
from app.pipeline.safety_checker import check_safety
from app.pipeline.yaml_repairer import repair_yaml

router = APIRouter(prefix="/ai", tags=["analyze"])


class ChapterItem(BaseModel):
    chapter_index: int
    title: str = ""
    content: str


class AnalyzeRequest(BaseModel):
    project_id: int
    chapters: list[ChapterItem]
    adaptation_mode: str = "screen_script"


@router.post("/analyze")
async def analyze_novel(req: AnalyzeRequest):
    """完整的小说分析流水线

    流程：人物抽取 → 剧情事件链 → 场景拆分 → 剧本生成 → 校验 → 修复
    """
    chapters = [ch.model_dump() for ch in req.chapters]
    logger.info(f"Starting full analysis for project {req.project_id}, {len(chapters)} chapters")

    try:
        # 1. 人物抽取
        logger.info("Step 1/6: Extracting characters")
        char_result = await extract_characters(chapters)
        characters = char_result["characters"]

        # 2. 剧情事件链
        logger.info("Step 2/6: Building plot events")
        plot_result = await build_plot_events(chapters, characters)
        plot_events = plot_result["plot_events"]

        # 3. 场景拆分
        logger.info("Step 3/6: Splitting scenes")
        scene_result = await split_scenes(chapters, characters, plot_events)
        scenes = scene_result["scenes"]

        # 4. 剧本生成
        logger.info("Step 4/6: Generating script")
        metadata = {
            "title": f"剧本 - 项目{req.project_id}",
            "source_title": f"小说 - 项目{req.project_id}",
            "adaptation_mode": req.adaptation_mode,
            "language": "zh-CN",
            "version": 1,
        }
        yaml_content = await generate_script(
            chapters, characters, plot_events, scenes, metadata
        )

        # 5. Schema 校验
        logger.info("Step 5/6: Validating schema")
        valid, issues = validate_yaml(yaml_content)

        # 6. 如果校验失败，尝试修复（最多 3 次）
        if not valid:
            for attempt in range(3):
                logger.info(f"Step 6/6: Attempting repair (attempt {attempt + 1})")
                repair_result = await repair_yaml(yaml_content, issues)
                if repair_result.get("success") and repair_result.get("fixed_yaml"):
                    yaml_content = repair_result["fixed_yaml"]
                    valid, issues = validate_yaml(yaml_content)
                    if valid:
                        break

        # 7. 幻觉检测
        logger.info("Checking hallucination")
        hallucination_result = await check_hallucination(yaml_content, chapters, characters)

        # 8. 安全审查
        logger.info("Checking safety")
        safety_result = await check_safety(yaml_content)

        return {
            "status": "completed",
            "yaml_content": yaml_content,
            "characters": characters,
            "plot_events": plot_events,
            "scenes": scenes,
            "validation": {"valid": valid, "issues": issues},
            "hallucination": hallucination_result,
            "safety": safety_result,
        }

    except Exception as e:
        logger.error(f"Analysis failed: {e}")
        raise HTTPException(status_code=500, detail=f"分析失败: {str(e)}")
