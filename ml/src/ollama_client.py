import asyncio
import json
from collections.abc import AsyncIterator
from typing import Any

import httpx

from config import OLLAMA_BASE_MODEL, OLLAMA_BASE_URL
from utils.logger import logger


class AsyncOllamaClient:
    def __init__(self, base_url: str = OLLAMA_BASE_URL) -> None:
        """Инициализация асинхронного клиента Ollama API.

        Parameters
        ----------
        base_url : str, optional
            Базовый URL сервера Ollama, по умолчанию "http://localhost:11434"
        """
        self.base_url = base_url

    async def generate(
        self,
        prompt: str,
        model: str,
        *,
        stream: bool = False,
        **kwargs: dict[str, Any],
    ) -> str | AsyncIterator[str] | None:
        """Генерация текста с помощью предустановленной модели.

        Parameters
        ----------
        prompt : str
            Текст промпта для генерации
        model: str
            Название модели для использования
        stream : bool, optional
            Режим потоковой передачи, по умолчанию False
        **kwargs
            Дополнительные параметры для API

        Returns
        -------
        str | AsyncIterator[str] | None
            Сгенерированный текст или итератор при потоковом режиме.
            Возвращает None в случае ошибки.

        Raises
        ------
        RuntimeError
            При возникновении ошибок сети или API
        """
        url = f"{self.base_url}/api/generate"
        payload = {
            "model": model,
            "prompt": prompt,
            "stream": stream,
            "num_ctx": NUM_CTX
            **kwargs,
        }

        async with httpx.AsyncClient(timeout=600) as client:
            try:
                response = await client.post(url, json=payload)
                response.raise_for_status()

                if stream:
                    return self._handle_stream_response(response)
                return self._handle_regular_response(response)

            except httpx.RequestError as e:
                msg = f"API request failed: {e}"
                raise RuntimeError(msg) from e

    def _handle_regular_response(self, response: httpx.Response) -> str:
        """Обработка обычного (не потокового) ответа."""
        result = response.json()
        return result.get("response", "")

    async def _handle_stream_response(
        self, response: httpx.Response
    ) -> AsyncIterator[str]:
        """Обработка потокового ответа."""
        async for line in response.aiter_lines():
            if line:
                chunk = json.loads(line)
                yield chunk.get("response", "")


async def main() -> None:
    client = AsyncOllamaClient()

    stream = await client.generate(
        prompt="Привет, как дела?",
        model="hf.co/t-tech/T-lite-it-1.0-Q8_0-GGUF:Q8_0",
        stream=True,
    )

    async for text in stream:
        logger.info(text)


if __name__ == "__main__":
    asyncio.run(main())
