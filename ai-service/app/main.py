"""Novel2Script-AI FastAPI 入口"""

from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse

from app.core.config import settings
from app.core.logger import logger

# 路由
from app.api.health import router as health_router
from app.api.analyze import router as analyze_router
from app.api.extract import router as extract_router
from app.api.plot import router as plot_router
from app.api.scene import router as scene_router
from app.api.generate import router as generate_router
from app.api.validate import router as validate_router
from app.api.hallucination import router as hallucination_router
from app.api.safety import router as safety_router
from app.api.repair import router as repair_router

app = FastAPI(
    title="Novel2Script-AI Service",
    description="小说转剧本 AI 服务",
    version="0.1.0",
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


# 内部认证中间件
@app.middleware("http")
async def verify_internal_token(request: Request, call_next):
    """验证内部服务调用 Token"""
    # 健康检查和文档无需认证
    if request.url.path in ("/health", "/docs", "/redoc", "/openapi.json"):
        return await call_next(request)

    # 验证 Token
    token = request.headers.get("Authorization", "").replace("Bearer ", "")
    if not settings.AI_SERVICE_TOKEN:
        return JSONResponse(
            status_code=503,
            content={"detail": "AI_SERVICE_TOKEN is not configured"},
        )
    if token != settings.AI_SERVICE_TOKEN:
        return JSONResponse(
            status_code=401,
            content={"detail": "Invalid internal token"},
        )

    return await call_next(request)


# 注册路由
app.include_router(health_router)
app.include_router(analyze_router)
app.include_router(extract_router)
app.include_router(plot_router)
app.include_router(scene_router)
app.include_router(generate_router)
app.include_router(validate_router)
app.include_router(hallucination_router)
app.include_router(safety_router)
app.include_router(repair_router)


@app.on_event("startup")
async def startup():
    logger.info(f"AI Service starting on port {settings.APP_PORT}")
    logger.info(f"LLM Provider: {settings.LLM_PROVIDER}")
    logger.info(f"LLM Model: {settings.LLM_MODEL}")


if __name__ == "__main__":
    import uvicorn

    uvicorn.run("app.main:app", host="0.0.0.0", port=settings.APP_PORT, reload=True)
