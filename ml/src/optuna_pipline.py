import json
import traceback
from typing import Any

import optuna
import redis.asyncio as redis
from datasets import Dataset
from langchain_ollama import OllamaEmbeddings
from ragas.dataset_schema import SingleTurnSample
from ragas.embeddings import LangchainEmbeddingsWrapper
from ragas.metrics import (
    NonLLMContextPrecisionWithReference,
    SemanticSimilarity,
)

import ml.v1.model_pb2 as ml_pb2_model
from config import DEFAULT_EMBEDER_MODEL, DEFAULT_REDIS_URL
from optuna_rag_params import load_all_qa_samples_from_redis
from RAG_pipeline import RAGPipeline
from utils.logger import logger


class OptunaPipeline:
    def __init__(
        self, redis_url: str | None, embedings_model: str | None
    ) -> None:
        """Инициализация класса OptunaPipeline.

        Parameters
        ----------
        redis_url : str
            URL Redis-сервера.
        embedings_model : str
            Название модели эмбедера'.
        """
        self.redis_url = redis_url if redis_url else DEFAULT_REDIS_URL
        self.embedings_model = embedings_model
        self.embeder = OllamaEmbeddings(
            model=embedings_model
            if embedings_model
            else DEFAULT_EMBEDER_MODEL,
        )
        self.context_metric = NonLLMContextPrecisionWithReference()
        self.generate_metric = SemanticSimilarity(
            embeddings=LangchainEmbeddingsWrapper(self.embeder)
        )
        self.rag_pipeline = RAGPipeline()

    async def compute_metrics(self, dataset: list[dict[str]]) -> float:
        """Вычисляет метрики NonLLMContextPrecisionWithReference для
        набора данных и SemanticSimilarity для оценки генерации.

        :param dataset: Dataset, который состоит из колонок:
            - 'question'
            - 'answer' (референсный ответ)
            - 'contexts' (извлечённые контексты)
            - 'generated_answer' (не используется в этой метрике, но можно оставить)
            - 'retrivment_content' (можно использовать как синоним contexts)
        :return: значение метрики (float)
        """
        context_result = 0
        generate_result = 0
        for i in range(len(dataset)):
            context_result += await self.context_metric.single_turn_ascore(
                SingleTurnSample(**dataset[i])
            )
            generate_result += await self.generate_metric.single_turn_ascore(
                SingleTurnSample(**dataset[i])
            )
        logger.info(
            "Context precision metrric %s", context_result / len(dataset)
        )
        logger.info("Semantic score: %s", context_result / len(dataset))
        return context_result / len(dataset), generate_result / len(dataset)

    async def load_all_qa_samples_from_redis(
        self, key_prefix: str = "qa:"
    ) -> list[dict]:
        redis_client = redis.from_url(
            self.redis_url, encoding="utf-8", decode_responses=True
        )

        keys = await redis_client.keys(f"{key_prefix}:*")
        keys.sort()

        samples = []
        for key in keys:
            value = await redis_client.get(key)
            if value:
                samples.append(json.loads(value))

        await redis_client.close()
        return samples

    def build_request(
        self, entry: dict, params: dict
    ) -> ml_pb2_model.ProcessQueryRequest:
        return ml_pb2_model.ProcessQueryRequest(
            query=ml_pb2_model.Query(content=entry["question"], userId=9),
            scenario=ml_pb2_model.Scenario(
                model=ml_pb2_model.LlmModel(
                    modelName=params["model"]["modelName"],
                    temperature=params["model"]["temperature"],
                    topK=params["model"]["topK"],
                    topP=params["model"]["topP"],
                    systemPrompt="",
                ),
                multiQuery=ml_pb2_model.MultiQuery(
                    useMultiquery=params["multiQuery"]["useMultiquery"],
                    nQueries=params["multiQuery"]["nQueries"],
                    queryModelName=params["multiQuery"]["queryModelName"]
                    or "",
                ),
                vectorSearch=ml_pb2_model.VectorSearch(
                    topN=params["vectorSearch"]["topN"],
                    threshold=params["vectorSearch"]["threshold"],
                    searchByQuery=params["vectorSearch"]["searchByQuery"],
                ),
                reranker=ml_pb2_model.Reranker(
                    useRerank=params["reranker"]["useRerank"],
                    topK=params["reranker"]["topK"],
                    rerankerMaxLength=params["reranker"]["rerankerMaxLength"],
                    rerankerModel=params["reranker"]["rerankerModel"],
                ),
            ),
            sourceIds=params["sourceIds"],
        )

    async def make_test_request(
        self, params: dict, test_dataset: list[dict[str, str]]
    ) -> Dataset:
        results = []
        for entry in test_dataset:
            request = self.build_request(entry, params)
            try:
                generated_answer, chunks = await self.rag_pipeline.generate(
                    request
                )

                logger.debug("Content: %s", generated_answer)
                logger.debug(
                    "Threashold: %s", params["vectorSearch"]["threshold"]
                )
                logger.debug("Chunks: %s", chunks)
                results.append(
                    {
                        "user_input": entry["question"],
                        "reference": entry["answer"],
                        "reference_contexts": [entry["context"]],
                        "response": generated_answer,
                        "retrieved_contexts": chunks,
                    }
                )
            except Exception as e:
                tb = traceback.format_exc()
                logger.info(f"Ошибка при обработке запроса: {e}\t{tb}")
                results.append(
                    {
                        "user_input": entry["question"],
                        "reference": entry["answer"],
                        "reference_contexts": [entry["context"]],
                        "response": "",
                        "retrieved_contexts": chunks,
                    }
                )

        return results

    async def evaluate_rag_pipeline(
        self, params: dict, test_dataset: list[dict[str, str]]
    ) -> tuple[float, float]:
        try:
            model_answers = await self.make_test_request(params, test_dataset)
            (
                context_precision_score,
                semantic_score,
            ) = await self.compute_metrics(model_answers)
            return context_precision_score, semantic_score
        except Exception as e:
            tb = traceback.format_exc()
            logger.error(f"Error during evaluation: {e}\n{tb}")
            return 0.0

    async def objective(
        self, trial, test_dataset: list[dict[str]], source_ids: list[str]
    ) -> float:
        params = {
            "vectorSearch": {
                "topN": trial.suggest_int("vectorSearch.topN", 3, 20),
                "threshold": trial.suggest_float(
                    "vectorSearch.threshold", 0.1, 0.2
                ),
                "searchByQuery": trial.suggest_categorical(
                    "vectorSearch.searchByQuery", [True, False]
                ),
            },
            "reranker": {
                "useRerank": trial.suggest_categorical(
                    "reranker.useRerank", [True, False]
                ),
                "topK": trial.suggest_int("reranker.topK", 1, 10),
                "rerankerMaxLength": trial.suggest_int(
                    "reranker.rerankerMaxLength", 128, 1024
                ),
                "rerankerModel": "BAAI/bge-reranker-v2-m3",
            },
            "model": {
                "temperature": trial.suggest_float(
                    "model.temperature", 0.0, 1.0
                ),
                "topK": trial.suggest_int("model.topK", 1, 100),
                "topP": trial.suggest_float("model.topP", 0.1, 1.0),
                "modelName": "hf.co/yandex/YandexGPT-5-Lite-8B-instruct-GGUF:Q4_K_M",
            },
            "multiQuery": {
                "useMultiquery": trial.suggest_categorical(
                    "multiQuery.useMultiquery", [True, False]
                ),
                "nQueries": trial.suggest_int("multiQuery.nQueries", 1, 5),
                "queryModelName": None,
            },
            "sourceIds": source_ids,
        }
        return await self.evaluate_rag_pipeline(params, test_dataset)

    async def study(self, source_ids: list[str]) -> list[dict[str, Any]]:
        """Создает и запускает Optuna-исследование для оптимизации
        параметров RAGPipeline.
        Параметры оптимизации включают:
        - vectorSearch.topN
        - vectorSearch.threshold
        - vectorSearch.searchByQuery
        - reranker.useRerank
        - reranker.topK
        - reranker.rerankerMaxLength
        - reranker.rerankerModel
        - model.temperature
        - model.topK
        - model.topP
        - model.modelName
        - multiQuery.useMultiquery
        - multiQuery.nQueries
        - multiQuery.queryModelName
        - sourceIds
        Параметры оптимизации могут быть изменены в методе objective.
        Исследование будет выполнено на тестовом наборе данных,
        который загружается из Redis по ключу, переданному в
        параметре source_ids.
        Параметры.
        ----------
        source_ids : list[str]
            Список идентификаторов источников для ограничения поиска.
        Возвращает
        -------
        list[dict[str, Any]]
            Список лучших испытаний, найденных Optuna.
        """
        study = optuna.create_study(directions=["maximize", "maximize"])
        test_dataset = []
        logger.info("Loading test dataset from Redis")
        logger.info("Redis URL: %s", self.redis_url)
        logger.info("Source IDs: %s", source_ids)
        for source_id in source_ids:
            logger.info("Loading data for source ID: %s", source_id)
            test_dataset += await load_all_qa_samples_from_redis(
                redis_url=self.redis_url, key_prefix=source_id
            )

        for _ in range(5):
            trial = study.ask()
            try:
                (
                    context_precision_score,
                    similarity_score,
                ) = await self.objective(trial, test_dataset, source_ids)
            except Exception as e:
                logger.error(f"Trial failed: {e}")
                context_precision_score, similarity_score = 0.0, 0.0
            study.tell(trial, [context_precision_score, similarity_score])
        best_trial = max(
            study.best_trials,
            key=lambda t: (t.values[0], t.values[1])  # сортировка по двум метрикам
        )
        return study.best_trials
