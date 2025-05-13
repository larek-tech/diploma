import os
from pathlib import Path

import dotenv

from utils.logger import logger

dotenv_path = Path(__file__).parents[1] / ".env"
logger.info(dotenv_path)
dotenv.load_dotenv(dotenv_path)

OLLAMA_BASE_URL = os.getenv("OLLAMA_BASE_URL")
GIGA_CHAT_API_KEY = os.getenv("GIGA_CHAT_API_KEY")
DEVICE = os.getenv("DEVICE")
DATA_SERVICE_PORT = os.getenv("DATA_SERVICE_PORT")
DATA_SERVICE_HOST = os.getenv("DATA_SERVICE_HOST")
DEFAULT_RERANKER_NAME = os.environ["DEFAULT_RERANKER_NAME"]
ML_SERVICE_PORT = os.getenv("ML_SERVICE_PORT")
HF_TOKEN = os.getenv("HF_TOKEN")
DEFAULT_REDIS_URL = os.getenv("DEFAULT_REDIS_URL")
DEFAULT_EMBEDER_MODEL = os.getenv("DEFAULT_EMBEDER_MODEL")

MULTI_QUESTION_PROMPT = """
Переформулируй вопрос: {query}. Предложите {n_questions} различных вариантов, которые помогут рассмотреть тему с разных сторон.
Требования:
Каждая версия вопроса должна сохранять исходный смысл, но использовать уникальные формулировки, стили (например, аналитический, гипотетический, практический) или акценты (например, на причину, последствия, примеры).
Избегайте простого перефразирования — стремитесь к разнообразию контекстов и аудиторий (например, для эксперта, новичка, студента).
Оформите ответ в виде нумерованного списка.
Пример:
Если исходный вопрос: «Почему небо голубое?», варианты могут быть:
«Какие физические явления приводят к голубому цвету неба?»
«Как объяснить ребенку, почему небо выглядит голубым?»
«Менялся ли цвет неба в истории Земли и почему?»
"""

# TODO: refactor
RAG_PROMPT = """
Вы — ИИ-ассистент, созданный для ответов на вопросы с использованием Retrieval-Augmented Generation (RAG). Ваша цель:  
1. Найти релевантную информацию в предоставленной базе знаний/документах.
2. Синтезировать эту информацию в ясный, точный и контекстно-осознанный ответ.

Вопрос пользователя:
{query}
Контекст:
{docs}
"""

# TODO: refactor
QA_PROMPT_TEMPLATE = """
Ты — интеллектуальный ассистент. На основе следующего текста сгенерируй один информативный вопрос и точный ответ.

Текст:
"{chunk}"

Формат ответа:
Вопрос: <вопрос>
Ответ: <ответ>
"""


FIRST_MESSAE_PROMPT = """Ты — интеллектуальный ассистент. На основе вопроса пользователя сделай навзвание чата и краткое описание.{message}"""
