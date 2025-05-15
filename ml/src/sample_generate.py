import json

import anyio
import redis.asyncio as redis

from config import (
    OLLAMA_BASE_MODEL,
    OLLAMA_BASE_URL,
    QA_PROMPT_TEMPLATE,
)
from data_client import AsyncDataServiceClient
from ollama_client import AsyncOllamaClient
from utils.logger import logger


def generate_structured_output_schema(n: int = 1) -> dict:
    """Генерирует JSON Schema для списка объектов
    с полями 'question' и 'answer'.

    Parameters
    ----------
    n : int
        Ожидаемое количество объектов в списке (повторений).

    Returns
    -------
    dict
        JSON Schema как словарь Python.
    """
    return {
        "type": "array",
        "minItems": n,
        "maxItems": n,
        "items": {
            "type": "object",
            "properties": {
                "question": {"type": "string"},
                "answer": {"type": "string"},
            },
            "required": ["question", "answer"],
            "additionalProperties": False,
        },
    }


class SyntheticDatasetGenerator:
    def __init__(self, model: str, base_url: str = OLLAMA_BASE_URL) -> None:
        self.client = AsyncOllamaClient(base_url=base_url)
        self.model = OLLAMA_BASE_MODEL

    async def generate_qa_pair(
        self, chunks: list[str], n_questions: int
    ) -> list[dict[str, str]]:
        """Генерирует QA-пары для каждого куска текста.
        Parametrs.
        ----------
        chunks : list[str]
            Список текстовых кусков, для которых нужно сгенерировать QA-пары.
        n_questions : int
            Количество вопросов, которые нужно
            сгенерировать для каждого куска текста.

        Returns
        -------
        list[dict[str, str]]
            Список словарей, каждый из которых содержит QA-пару.
        """
        results = []

        for chunk in chunks:
            prompt = QA_PROMPT_TEMPLATE.format(chunk=chunk)

            response = await self.client.generate(
                prompt=prompt,
                model=self.model,
                format=generate_structured_output_schema(n=n_questions),
            )
            qa_pairs = json.loads(response)

            for pair in qa_pairs:
                results.append(  # noqa: PERF401
                    {
                        "question": pair["question"],
                        "answer": pair["answer"],
                        "context": chunk,
                    }
                )
        return results

    async def save_to_jsonl(
        self, data: list[dict[str, str]], output_path: str
    ) -> None:
        """Сохраняет список QA-пар в JSONL-файл."""
        async with await anyio.open_file(
            output_path, "w", encoding="utf-8"
        ) as f:
            content = json.dumps(data, indent=4, ensure_ascii=False)
            await f.write(content)

    async def save_to_redis(
        self,
        data: list[dict[str, str]],
        redis_url: str = "redis://localhost",
        key_prefix: str = "qa",
    ) -> None:
        """Сохраняет переданные данные в Redis с указанным префиксом ключей.

        Parametrs
        ----------
        data : list of dict[str, str]
            Список словарей для сохранения. Каждый словарь должен
            содержать строковые ключи и значения.
            Пример: [{"question": "текст", "answer": "текст"}, ...]
        redis_url : str, optional, default = "redis://localhost"
            URL-адрес Redis-сервера.
        key_prefix : str, optional, default = "qa"
            Префикс для формирования ключей в Redis.
            Ключи будут иметь вид: `{префикс}:{индекс}`.
        """
        redis_client = redis.from_url(
            redis_url, encoding="utf-8", decode_responses=True
        )
        for i, item in enumerate(data):
            key = f"{key_prefix}:{i}"
            await redis_client.set(key, json.dumps(item, ensure_ascii=False))
        await redis_client.close()


async def generate_dataset(
    source_ids: list[str], data_client: AsyncDataServiceClient
) -> None:
    generator = SyntheticDatasetGenerator(model=OLLAMA_BASE_MODEL)
    for source_id in source_ids:
        response = await data_client.get_documents(source_id, size=25, page=1)
        chunks = [doc.content for doc in response.documents]
        logger.info(chunks)
        dataset = await generator.generate_qa_pair(chunks, n_questions=5)

        output_path = "synthetic_dataset.jsonl"
        await generator.save_to_jsonl(dataset, output_path)

        await generator.save_to_redis(dataset, key_prefix=source_id)

    logger.info(f"Сохранено {len(dataset)} QA-пар в {output_path}")
