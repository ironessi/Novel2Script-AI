"""统一 LLM 客户端，支持 DeepSeek / OpenAI / Ollama"""

import json
from typing import Optional

import httpx
from openai import AsyncOpenAI

from app.core.config import settings
from app.core.logger import logger


class LLMClient:
    """统一的 LLM 调用客户端"""

    def __init__(self):
        self._client: Optional[AsyncOpenAI] = None
        self._http_client: Optional[httpx.AsyncClient] = None

    def _get_openai_client(self) -> AsyncOpenAI:
        """获取 OpenAI 兼容客户端（DeepSeek / OpenAI）"""
        if self._client is None:
            self._client = AsyncOpenAI(
                api_key=settings.LLM_API_KEY,
                base_url=settings.LLM_BASE_URL,
                timeout=300.0,
            )
        return self._client

    def _get_http_client(self) -> httpx.AsyncClient:
        """获取 HTTP 客户端（Ollama）"""
        if self._http_client is None:
            self._http_client = httpx.AsyncClient(timeout=300.0)
        return self._http_client

    async def chat(
        self,
        system_prompt: str,
        user_prompt: str,
        temperature: float = 0.3,
        max_tokens: int = 4096,
    ) -> str:
        """调用 LLM 生成回复"""
        provider = settings.LLM_PROVIDER.lower()

        if provider == "ollama":
            return await self._chat_ollama(system_prompt, user_prompt, temperature)
        else:
            return await self._chat_openai(
                system_prompt, user_prompt, temperature, max_tokens
            )

    async def _chat_openai(
        self,
        system_prompt: str,
        user_prompt: str,
        temperature: float,
        max_tokens: int,
    ) -> str:
        """通过 OpenAI 兼容接口调用（DeepSeek / OpenAI）"""
        client = self._get_openai_client()
        model = settings.LLM_MODEL

        logger.info(f"Calling LLM: provider=openai, model={model}")

        try:
            response = await client.chat.completions.create(
                model=model,
                messages=[
                    {"role": "system", "content": system_prompt},
                    {"role": "user", "content": user_prompt},
                ],
                temperature=temperature,
                max_tokens=max_tokens,
            )
            content = response.choices[0].message.content
            logger.info(f"LLM response length: {len(content)} chars")
            return content
        except Exception as e:
            logger.error(f"LLM call failed: {e}")
            raise

    async def _chat_ollama(
        self, system_prompt: str, user_prompt: str, temperature: float
    ) -> str:
        """通过 Ollama API 调用本地模型"""
        client = self._get_http_client()
        model = settings.OLLAMA_MODEL

        logger.info(f"Calling LLM: provider=ollama, model={model}")

        try:
            response = await client.post(
                f"{settings.OLLAMA_BASE_URL}/api/chat",
                json={
                    "model": model,
                    "messages": [
                        {"role": "system", "content": system_prompt},
                        {"role": "user", "content": user_prompt},
                    ],
                    "stream": False,
                    "options": {"temperature": temperature},
                },
            )
            response.raise_for_status()
            data = response.json()
            content = data["message"]["content"]
            logger.info(f"LLM response length: {len(content)} chars")
            return content
        except Exception as e:
            logger.error(f"Ollama call failed: {e}")
            raise

    async def chat_json(
        self,
        system_prompt: str,
        user_prompt: str,
        temperature: float = 0.1,
        max_tokens: int = 8192,
    ) -> dict:
        """调用 LLM 并解析 JSON 响应"""
        # 在 system_prompt 末尾强制要求 JSON 输出
        json_system = (
            system_prompt
            + "\n\n你必须只返回合法的 JSON 对象，不要包含任何其他文字、解释或 markdown 标记。"
        )

        raw = await self.chat(json_system, user_prompt, temperature, max_tokens)

        # 清理响应，提取 JSON
        content = raw.strip()

        # 去除 markdown 代码块标记
        if content.startswith("```"):
            lines = content.split("\n")
            # 去掉第一行和最后一行
            start = 1
            end = len(lines)
            for i, line in enumerate(lines):
                if i > 0 and line.strip().startswith("```"):
                    end = i
                    break
            content = "\n".join(lines[start:end]).strip()

        try:
            return json.loads(content)
        except json.JSONDecodeError as e:
            logger.error(f"Failed to parse LLM JSON response: {e}")
            logger.error(f"Raw response (first 500 chars): {content[:500]}")
            raise ValueError(f"LLM 返回的内容不是合法 JSON: {e}")


# 全局单例
llm_client = LLMClient()
