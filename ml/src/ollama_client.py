import asyncio
import json
from collections.abc import AsyncIterator
from typing import Any, Union

import httpx

from config import OLLAMA_BASE_MODEL, OLLAMA_BASE_URL, NUM_CTX
from utils.logger import logger


class AsyncOllamaClient:
    def __init__(self, base_url: str = OLLAMA_BASE_URL) -> None:
        self.base_url = base_url

    async def generate(
        self,
        prompt: str,
        model: str,
        *,
        stream: bool = False,
        **kwargs: dict[str, Any],
    ) -> Union[str, AsyncIterator[str], None]:
        """Генерация текста с использованием модели Ollama."""
        url = f"{self.base_url}/api/generate"
        payload = {
            "model": model,
            "prompt": prompt,
            "stream": stream,
            "num_ctx": NUM_CTX,
            **kwargs,
        }

        try:
            logger.debug(f"Sending payload: {payload}")

            if stream:
                return self._handle_stream_response(payload)

            async with httpx.AsyncClient(timeout=600) as client:
                response = await client.post(url, json=payload)
                response.raise_for_status()
                return self._handle_regular_response(response)

        except httpx.HTTPStatusError as e:
            msg = f"Got bad status: {e}"
            raise RuntimeError(msg) from e
        except httpx.RequestError as e:
            msg = f"API request failed: {e}"
            raise RuntimeError(msg) from e

    def _handle_regular_response(self, response: httpx.Response) -> str:
        """Обработка обычного (не потокового) ответа."""
        result = response.json()
        return result.get("response", "")

    async def _handle_stream_response(
        self, payload: dict
    ) -> AsyncIterator[str]:
        """Обработка потокового ответа."""
        url = f"{self.base_url}/api/generate"
        logger.info(f"Sending payload for stream: {payload}")

        async with httpx.AsyncClient(timeout=600) as client:
            async with client.stream("POST", url, json=payload) as response:
                response.raise_for_status()

                async for line in response.aiter_lines():
                    if line.strip():
                        try:
                            data = json.loads(line)
                            text = data.get("response", "")
                            if text:
                                yield text
                        except json.JSONDecodeError as e:
                            logger.error(f"JSON decode error: {e}")


async def main() -> None:
    client = AsyncOllamaClient()

    stream = await client.generate(
        prompt="Привет, как дела?",
        model=OLLAMA_BASE_MODEL,
        stream=True,
    )

    async for text in stream:
        logger.info(text)


if __name__ == "__main__":
    asyncio.run(main())
