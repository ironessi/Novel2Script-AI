import os
from dotenv import load_dotenv

# 加载 .env（从项目根目录）
load_dotenv(os.path.join(os.path.dirname(__file__), "../../../.env"))


class Settings:
    """AI 服务配置"""

    # 服务
    APP_ENV: str = os.getenv("APP_ENV", "local")
    APP_PORT: int = int(os.getenv("AI_SERVICE_PORT", "9000"))

    # AI Service 认证
    AI_SERVICE_TOKEN: str = os.getenv("AI_SERVICE_TOKEN", "")

    # LLM 配置
    LLM_PROVIDER: str = os.getenv("LLM_PROVIDER", "deepseek")
    LLM_API_KEY: str = os.getenv("LLM_API_KEY", "")
    LLM_MODEL: str = os.getenv("LLM_MODEL", "deepseek-chat")
    LLM_BASE_URL: str = os.getenv("LLM_BASE_URL", "https://api.deepseek.com/v1")

    # Ollama 配置（本地模型）
    OLLAMA_BASE_URL: str = os.getenv("OLLAMA_BASE_URL", "http://localhost:11434")
    OLLAMA_MODEL: str = os.getenv("OLLAMA_MODEL", "qwen2.5:14b")

    # Redis
    REDIS_HOST: str = os.getenv("REDIS_HOST", "localhost")
    REDIS_PORT: int = int(os.getenv("REDIS_PORT", "6379"))
    REDIS_PASSWORD: str = os.getenv("REDIS_PASSWORD", "")
    REDIS_DB: int = int(os.getenv("REDIS_DB", "2"))


settings = Settings()
