import json

from config import MULTI_QUESTION_PROMPT
from ollama_client import AsyncOllamaClient


def generate_schema(
    n: int,
    schema_title: str,
    title_template: str,
    base_name: str = "question_id",
) -> dict:
    """Сгенерировать словарь с настраиваемой структурой для JSON Schema."""
    properties = {}
    required = []

    for i in range(1, n + 1):
        field_name = f"{base_name}_{i}"
        properties[field_name] = {
            "title": title_template.format(n=i),
            "type": "string",
        }
        required.append(field_name)

    return {
        "properties": properties,
        "required": required,
        "title": schema_title,
        "type": "object",
    }


async def get_multi_questions(
    client: AsyncOllamaClient,
    user_prompt: str,
    n_questions: int,
    model: str,
) -> list[str]:
    schema = generate_schema(
        n=n_questions,
        schema_title="{i} перефразированный вопрос пользователя.",
        title_template="Перефразированные вопрос пользователя.",
    )
    return list(
        json.loads(
            await client.generate(
                prompt=MULTI_QUESTION_PROMPT.format(
                    query=user_prompt,
                    n_questions=n_questions,
                ),
                format=schema,
                model=model,
            )
        ).values()
    )
