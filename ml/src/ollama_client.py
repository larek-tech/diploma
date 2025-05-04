import json
from collections.abc import Iterator
from typing import Any

import requests

from config import OLLAMA_BASE_URL
from utils.logger import logger


class OllamaClient:
    def __init__(self, base_url: str = "http://localhost:11434") -> None:
        """Инициализация клиента Ollama API.

        Parameters
        ----------
        base_url : str, optional
            Базовый URL сервера Ollama, по умолчанию "http://localhost:11434"
        """
        self.base_url = base_url

    def generate(
        self,
        prompt: str,
        model: str,
        *,
        stream: bool = False,
        **kwargs: dict[str, Any],
    ) -> str | Iterator[str] | None:
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
        str | Iterator[str] | None
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
            **kwargs,
        }

        try:
            response = requests.post(
                url, json=payload, stream=stream, timeout=600
            )
            response.raise_for_status()

            if stream:
                return self._handle_stream_response(response)
            return self._handle_regular_response(response)

        except requests.exceptions.RequestException as e:
            raise RuntimeError(f"API request failed: {e}") from e

    def _handle_regular_response(self, response: requests.Response) -> str:
        """Обработка обычного (не потокового) ответа."""
        result = response.json()
        return result.get("response", "")

    def _handle_stream_response(
        self, response: requests.Response
    ) -> Iterator[str]:
        """Обработка потокового ответа."""
        for line in response.iter_lines():
            if line:
                chunk = json.loads(line.decode("utf-8"))
                yield chunk.get("response", "")


if __name__ == "__main__":
    client = OllamaClient(
        model="hf.co/yandex/YandexGPT-5-Lite-8B-instruct-GGUF:Q4_K_M",
        base_url=OLLAMA_BASE_URL,
    )

    for text in client.generate(
        prompt="Привет, как дела?",
        stream=True,
    ):
        logger.info(text, end=" ")
