import json

from config import MULTI_QUESTION_PROMPT
from ollama_client import AsyncOllamaClient
from utils.logger import logger

def generate_schema(n: int = 1) -> dict:
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
            },
            "required": ["question"],
            "additionalProperties": False,
        },
    }



async def get_multi_questions(
    client: AsyncOllamaClient,
    user_prompt: str,
    n_questions: int,
    model: str,
) -> list[str]:
    schema = generate_schema(
        n=n_questions,
        # schema_title="{i} перефразированный вопрос пользователя.",
        # title_template="Перефразированные вопрос пользователя.",
    )
    answers =  json.loads(
            await client.generate(
                prompt=MULTI_QUESTION_PROMPT.format(
                    query=user_prompt,
                    n_questions=n_questions,
                ),
                format=schema,
                model=model,
            )
        )
    logger.info(answers)
    return answers
