import asyncio
import json
import time
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
from config import (
    DATA_SERVICE_HOST,
    DATA_SERVICE_PORT,
    DEFAULT_EMBEDER_MODEL,
    DEFAULT_REDIS_URL,
    OLLAMA_BASE_MODEL,
    OLLAMA_BASE_URL,
)
from data_client import AsyncDataServiceClient
from RAG_pipeline import RAGPipeline
from utils.logger import logger

optuna.logging.enable_propagation()


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
            base_url=OLLAMA_BASE_URL,
            model=embedings_model
            if embedings_model
            else DEFAULT_EMBEDER_MODEL,
        )
        self.context_metric = NonLLMContextPrecisionWithReference()
        self.generate_metric = SemanticSimilarity(
            embeddings=LangchainEmbeddingsWrapper(self.embeder)
        )
        self.rag_pipeline = RAGPipeline()

    async def calculate_ragas_metrics(
        self, model_outputs: list[dict[str, Any]]
    ) -> tuple[float, float]:
        """Расчёт метрик RAGAS (ContextPrecision и SemanticSimilarity).

        :param model_outputs: Список словарей с полями:
            - 'user_input' (вопрос пользователя)
            - 'reference' (референсный ответ)
            - 'retrieved_contexts' (список извлечённых контекстов)
            - 'response' (сгенерированный ответ)
        :return: Кортеж из двух float:
            - context_precision_score
            - semantic_similarity_score
        """
        context_result = 0.0
        total_samples = len(model_outputs)

        for sample in model_outputs:
            rag_sample = SingleTurnSample(**sample)
            context_result += await self.context_metric.single_turn_ascore(
                rag_sample
            )

        if total_samples == 0:
            return 0.0, 0.0

        context_score = context_result / total_samples

        logger.info("Context Precision: %s", context_score)

        return context_score

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
            query=ml_pb2_model.Query(
                content=entry.get("question", ""), userId=9
            ),
            scenario=ml_pb2_model.Scenario(
                model=ml_pb2_model.LlmModel(
                    modelName=params.get("model", {}).get("modelName"),
                    temperature=params.get("model", {}).get("temperature"),
                    topK=params.get("model", {}).get("topK"),
                    topP=params.get("model", {}).get("topP"),
                    systemPrompt="",
                ),
                multiQuery=ml_pb2_model.MultiQuery(
                    useMultiquery=params.get("multiQuery", {}).get(
                        "useMultiquery"
                    ),
                    nQueries=params.get("multiQuery", {}).get("nQueries"),
                    queryModelName=params.get("multiQuery", {}).get(
                        "queryModelName"
                    )
                    or "",
                ),
                vectorSearch=ml_pb2_model.VectorSearch(
                    topN=params.get("vectorSearch", {}).get("topN"),
                    threshold=params.get("vectorSearch", {}).get("threshold"),
                    searchByQuery=params.get("vectorSearch", {}).get(
                        "searchByQuery"
                    ),
                ),
                reranker=ml_pb2_model.Reranker(
                    useRerank=params.get("reranker", {}).get("useRerank"),
                    topK=params.get("reranker", {}).get("topK"),
                    rerankerMaxLength=params.get("reranker", {}).get(
                        "rerankerMaxLength"
                    ),
                    rerankerModel=params.get("reranker", {}).get(
                        "rerankerModel"
                    ),
                ),
            ),
            sourceIds=params.get("sourceIds"),
        )

    async def make_test_request(
        self, params: dict, test_dataset: list[dict[str, str]]
    ) -> Dataset:
        results = []
        for entry in test_dataset:
            request = self.build_request(entry, params)
            try:

                chunks = await self.rag_pipeline._prepare_chunks(request)

                logger.debug(
                    "Threashold: %s", params["vectorSearch"]["threshold"]
                )
                logger.debug("Chunks: %s", chunks)
                results.append(
                    {
                        "user_input": entry["question"],
                        "reference": entry["answer"],
                        "reference_contexts": [entry["context"]],
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
            model_answers= await self.make_test_request(
                params, test_dataset
            )
            (
                context_precision_score
                # semantic_score,
            ) = await self.calculate_ragas_metrics(model_answers)
            return context_precision_score
        except Exception as e:
            tb = traceback.format_exc()
            logger.error(f"Error during evaluation: {e}\n{tb}")
            return 0.0, 0.0, 0.0

    async def objective(
        self,
        trial: optuna.Trial,
        test_dataset: list[dict[str]],
        source_ids: list[str],
    ) -> float:
        params = {
            "vectorSearch": {
                "topN": trial.suggest_int("vectorSearch.topN", 8, 15),
                "threshold": trial.suggest_float(
                    "vectorSearch.threshold", 0.2, 0.9, step=0.01
                ),
                "searchByQuery": trial.suggest_categorical(
                    "vectorSearch.searchByQuery", [False]
                ),
            },
            # "reranker": {
            #     "useRerank": trial.suggest_categorical(
            #         "reranker.useRerank", [True, False]
            #     ),
            #     "topK": trial.suggest_int("reranker.topK", 4, 8)
            #     if trial.params["reranker.useRerank"]
            #     else None,
            #     "rerankerMaxLength": 4096,
            #     "rerankerModel": "BAAI/bge-reranker-v2-m3",
            # },
            "model": {
                "temperature": 0.2,
                "topK":20,
                "topP": 0.5,
                "modelName": OLLAMA_BASE_MODEL,
            },
            "multiQuery": {
                "useMultiquery": trial.suggest_categorical(
                    "multiQuery.useMultiquery", [True]
                ),
                "nQueries": trial.suggest_int("multiQuery.nQueries", 3, 5)
                if trial.params["multiQuery.useMultiquery"]
                else None,
                "queryModelName": None,
            },
            "sourceIds": source_ids,
        }
        return await self.evaluate_rag_pipeline(params, test_dataset)

    async def study(
        self,
        source_ids: list[str],
        params: dict | None = None,
    ) -> list[dict[str, Any]]:
        """Создает и запускает Optuna-исследование для оптимизации
        параметров RAGPipeline.
        Параметры оптимизации включают:
        - vectorSearch.topN
        - vectorSearch.threshold
        - vectorSearch.searchByQuery
        - reranker.useRerank
        - reranker.topK
        - model.temperature
        - model.topK
        - model.topP
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
        test_dataset = []
        logger.info("Loading test dataset from Redis")
        logger.info("Redis URL: %s", self.redis_url)
        logger.info("Source IDs: %s", source_ids)
        for source_id in source_ids:
            logger.info("Loading data for source ID: %s", source_id)
            test_dataset += await self.load_all_qa_samples_from_redis(
                key_prefix=source_id
            )
        if params:
            logger.info("Evaluating with provided params (без Optuna)")
            context_score, semantic_score = await self.evaluate_rag_pipeline(
                params, test_dataset
            )
            return [
                {
                    "context_precision": context_score,
                    "semantic_similarity": semantic_score,
                    "params": params,
                }
            ]
        storage_path = f"sqlite:///optuna_study_multiquery_rag_{'_'.join(source_ids)}.db"
        study = optuna.create_study(
            direction="maximize",
            study_name=f"Default RAG params in sources: {source_ids}",
            storage=storage_path,
            load_if_exists=True,
        )
        logger.info("Study name: %s", study.study_name)
        for trial_num in range(50):
            trial = study.ask()
            try:
                context_precision_score = await self.objective(trial, test_dataset, source_ids)
            except Exception as e:
                logger.error(f"Trial failed: {e}")
                context_precision_score = 0.0
            logger.info(
                "Trial #%d\n params: %s\n metric: %s\n",
                trial_num,
                trial.params,
                context_precision_score,
                # average_time,
            )
            study.tell(trial, context_precision_score)
        # best_trial = max(
        #     study.best_trials,
        #     key=lambda t: (
        #         t.values[0],
        #         t.values[1],
        #     ),  # сортировка по двум метрикам
        # )
        # logger.info("Best trial found:")
        # logger.info(
        #     " - Values: context=%.4f, semantic=%.4f",
        #     best_trial.values[0],
        #     best_trial.values[1],
        # )
        logger.info(" - Params: %s", study.best_params)
        return self.trial_to_model_params(study.best_params)

    def trial_to_model_params(self, params):
        model_params = ml_pb2_model.ModelParams()

        # MultiQuery
        if (
            "multiQuery.useMultiquery" in params
            and "multiQuery.nQueries" in params
        ):
            model_params.multiQuery.useMultiquery = params[
                "multiQuery.useMultiquery"
            ]
            model_params.multiQuery.nQueries = int(
                params["multiQuery.nQueries"]
            )

        # Reranker
        if all(
            k in params
            for k in [
                "reranker.useRerank",
                "reranker.topK",
                "reranker.rerankerMaxLength",
            ]
        ):
            model_params.reranker.useRerank = params["reranker.useRerank"]
            model_params.reranker.topK = int(params["reranker.topK"])
            model_params.reranker.rerankerMaxLength = int(
                params["reranker.rerankerMaxLength"]
            )
            model_params.reranker.rerankerModel = "BAAI/bge-reranker-v2-m3"

        # VectorSearch
        if all(
            k in params
            for k in [
                "vectorSearch.topN",
                "vectorSearch.threshold",
                "vectorSearch.searchByQuery",
            ]
        ):
            model_params.vectorSearch.topN = int(params["vectorSearch.topN"])
            model_params.vectorSearch.threshold = float(
                params["vectorSearch.threshold"]
            )
            model_params.vectorSearch.searchByQuery = params[
                "vectorSearch.searchByQuery"
            ]

        model_params.model.modelName = OLLAMA_BASE_MODEL
        # model_params.model.temperature = float(params["model.temperature"])
        # model_params.model.topK = int(params["model.topK"])
        # model_params.model.topP = float(params["model.topP"])
        model_params.model.systemPrompt = ""

        return model_params


async def main():
    pipline = OptunaPipeline(
        redis_url="redis://192.168.1.5:6379",
        embedings_model=DEFAULT_EMBEDER_MODEL,
    )
    data_client = AsyncDataServiceClient(
        host=DATA_SERVICE_HOST, port=DATA_SERVICE_PORT
    )
    source_ids = ["d028f055-c743-4b5f-9995-bd9fcf3b4330"]
    # await generate_dataset(source_ids, data_client)
    params = await pipline.study(
        source_ids=source_ids,
    )
    logger.info("Best params: %s", params)

def foo():
    default_rag_params = {
        "vectorSearch": {
            "topN": trial.suggest_int("vectorSearch.topN", 9, 20),
            "threshold": trial.suggest_float(
                "vectorSearch.threshold", 0.1, 0.9
            ),
            "searchByQuery": trial.suggest_categorical(
                "vectorSearch.searchByQuery", [False]
            ),
        },
        # "reranker": {
        #     "useRerank": trial.suggest_categorical(
        #         "reranker.useRerank", [True, False]
        #     ),
        #     "topK": trial.suggest_categorical("reranker.topK", 4, 8),
        #     "rerankerMaxLength": trial.suggest_int(
        #         "reranker.rerankerMaxLength",
        #         [4096, 8192]
        #     ),
        #     "rerankerModel": "BAAI/bge-reranker-v2-m3",
        # },
        "model": {
            "temperature": trial.suggest_float("model.temperature", 0.0, 1.0),
            "topK": trial.suggest_int("model.topK", 1, 50),
            "topP": trial.suggest_float("model.topP", 0.1, 1.0),
            "modelName": OLLAMA_BASE_MODEL,
        },
        # "multiQuery": {
        #     "useMultiquery": trial.suggest_categorical(
        #         "multiQuery.useMultiquery", [True, False]
        #     ),
        #     "nQueries": trial.suggest_int("multiQuery.nQueries", 3, 8),
        #     "queryModelName": None,
        # },
        "sourceIds": source_ids,
    }
    multquery_rag_params = {
        "vectorSearch": {
            "topN": trial.suggest_int("vectorSearch.topN", 9, 20),
            "threshold": trial.suggest_float(
                "vectorSearch.threshold", 0.1, 0.9
            ),
            "searchByQuery": trial.suggest_categorical(
                "vectorSearch.searchByQuery", [False]
            ),
        },
        # "reranker": {
        #     "useRerank": trial.suggest_categorical(
        #         "reranker.useRerank", [True, False]
        #     ),
        #     "topK": trial.suggest_categorical("reranker.topK", 4, 8),
        #     "rerankerMaxLength": trial.suggest_int(
        #         "reranker.rerankerMaxLength",
        #         [4096, 8192]
        #     ),
        #     "rerankerModel": "BAAI/bge-reranker-v2-m3",
        # },
        "model": {
            "temperature": trial.suggest_float("model.temperature", 0.0, 1.0),
            "topK": trial.suggest_int("model.topK", 1, 50),
            "topP": trial.suggest_float("model.topP", 0.1, 1.0),
            "modelName": OLLAMA_BASE_MODEL,
        },
        "multiQuery": {
            "useMultiquery": trial.suggest_categorical(
                "multiQuery.useMultiquery", [True]
            ),
            "nQueries": trial.suggest_int("multiQuery.nQueries", 3, 8),
            "queryModelName": None,
        },
        "sourceIds": source_ids,
    }
    hypotetical_question_rag_params = {
        "vectorSearch": {
            "topN": trial.suggest_int("vectorSearch.topN", 9, 20),
            "threshold": trial.suggest_float(
                "vectorSearch.threshold", 0.1, 0.9
            ),
            "searchByQuery": trial.suggest_categorical(
                "vectorSearch.searchByQuery", [True]
            ),
        },
        # "reranker": {
        #     "useRerank": trial.suggest_categorical(
        #         "reranker.useRerank", [True, False]
        #     ),
        #     "topK": trial.suggest_categorical("reranker.topK", 4, 8),
        #     "rerankerMaxLength": trial.suggest_int(
        #         "reranker.rerankerMaxLength",
        #         [4096, 8192]
        #     ),
        #     "rerankerModel": "BAAI/bge-reranker-v2-m3",
        # },
        "model": {
            "temperature": trial.suggest_float("model.temperature", 0.0, 1.0),
            "topK": trial.suggest_int("model.topK", 1, 50),
            "topP": trial.suggest_float("model.topP", 0.1, 1.0),
            "modelName": OLLAMA_BASE_MODEL,
        },
        # "multiQuery": {
        #     "useMultiquery": trial.suggest_categorical(
        #         "multiQuery.useMultiquery", [True]
        #     ),
        #     "nQueries": trial.suggest_int("multiQuery.nQueries", 3, 8),
        #     "queryModelName": None,
        # },
        "sourceIds": source_ids,
    }
    reranker_rag_params = {
        "vectorSearch": {
            "topN": trial.suggest_int("vectorSearch.topN", 9, 20),
            "threshold": trial.suggest_float(
                "vectorSearch.threshold", 0.1, 0.9
            ),
            "searchByQuery": trial.suggest_categorical(
                "vectorSearch.searchByQuery", [False]
            ),
        },
        "reranker": {
            "useRerank": trial.suggest_categorical(
                "reranker.useRerank", [True, False]
            ),
            "topK": trial.suggest_categorical("reranker.topK", 4, 8),
            "rerankerMaxLength": trial.suggest_int(
                "reranker.rerankerMaxLength", [4096, 8192]
            ),
            "rerankerModel": "BAAI/bge-reranker-v2-m3",
        },
        "model": {
            "temperature": trial.suggest_float("model.temperature", 0.0, 1.0),
            "topK": trial.suggest_int("model.topK", 1, 50),
            "topP": trial.suggest_float("model.topP", 0.1, 1.0),
            "modelName": OLLAMA_BASE_MODEL,
        },
        # "multiQuery": {
        #     "useMultiquery": trial.suggest_categorical(
        #         "multiQuery.useMultiquery", [True]
        #     ),
        #     "nQueries": trial.suggest_int("multiQuery.nQueries", 3, 8),
        #     "queryModelName": None,
        # },
        "sourceIds": source_ids,
    }
    hypotetical_multiquery_rag_params = {
        "vectorSearch": {
            "topN": trial.suggest_int("vectorSearch.topN", 9, 20),
            "threshold": trial.suggest_float(
                "vectorSearch.threshold", 0.1, 0.9
            ),
            "searchByQuery": trial.suggest_categorical(
                "vectorSearch.searchByQuery", [True]
            ),
        },
        # "reranker": {
        #     "useRerank": trial.suggest_categorical(
        #         "reranker.useRerank", [True, False]
        #     ),
        #     "topK": trial.suggest_categorical("reranker.topK", 4, 8),
        #     "rerankerMaxLength": trial.suggest_int(
        #         "reranker.rerankerMaxLength",
        #         [4096, 8192]
        #     ),
        #     "rerankerModel": "BAAI/bge-reranker-v2-m3",
        # },
        "model": {
            "temperature": trial.suggest_float("model.temperature", 0.0, 1.0),
            "topK": trial.suggest_int("model.topK", 1, 50),
            "topP": trial.suggest_float("model.topP", 0.1, 1.0),
            "modelName": OLLAMA_BASE_MODEL,
        },
        "multiQuery": {
            "useMultiquery": trial.suggest_categorical(
                "multiQuery.useMultiquery", [True]
            ),
            "nQueries": trial.suggest_int("multiQuery.nQueries", 3, 8),
            "queryModelName": None,
        },
        "sourceIds": source_ids,
    }
    hypotetical_reranker_rag_params = {
        "vectorSearch": {
            "topN": trial.suggest_int("vectorSearch.topN", 9, 20),
            "threshold": trial.suggest_float(
                "vectorSearch.threshold", 0.1, 0.9
            ),
            "searchByQuery": trial.suggest_categorical(
                "vectorSearch.searchByQuery", [True]
            ),
        },
        "reranker": {
            "useRerank": trial.suggest_categorical(
                "reranker.useRerank", [True, False]
            ),
            "topK": trial.suggest_categorical("reranker.topK", 4, 8),
            "rerankerMaxLength": trial.suggest_int(
                "reranker.rerankerMaxLength", [4096, 8192]
            ),
            "rerankerModel": "BAAI/bge-reranker-v2-m3",
        },
        "model": {
            "temperature": trial.suggest_float("model.temperature", 0.0, 1.0),
            "topK": trial.suggest_int("model.topK", 1, 50),
            "topP": trial.suggest_float("model.topP", 0.1, 1.0),
            "modelName": OLLAMA_BASE_MODEL,
        },
        # "multiQuery": {
        #     "useMultiquery": trial.suggest_categorical(
        #         "multiQuery.useMultiquery", [True]
        #     ),
        #     "nQueries": trial.suggest_int("multiQuery.nQueries", 3, 8),
        #     "queryModelName": None,
        # },
        "sourceIds": source_ids,
    }
    reranker_multiquery_rag_params = {
        "vectorSearch": {
            "topN": trial.suggest_int("vectorSearch.topN", 9, 20),
            "threshold": trial.suggest_float(
                "vectorSearch.threshold", 0.1, 0.9
            ),
            "searchByQuery": trial.suggest_categorical(
                "vectorSearch.searchByQuery", [False]
            ),
        },
        "reranker": {
            "useRerank": trial.suggest_categorical(
                "reranker.useRerank", [True, False]
            ),
            "topK": trial.suggest_categorical("reranker.topK", 4, 8),
            "rerankerMaxLength": trial.suggest_int(
                "reranker.rerankerMaxLength", [4096, 8192]
            ),
            "rerankerModel": "BAAI/bge-reranker-v2-m3",
        },
        "model": {
            "temperature": trial.suggest_float("model.temperature", 0.0, 1.0),
            "topK": trial.suggest_int("model.topK", 1, 50),
            "topP": trial.suggest_float("model.topP", 0.1, 1.0),
            "modelName": OLLAMA_BASE_MODEL,
        },
        "multiQuery": {
            "useMultiquery": trial.suggest_categorical(
                "multiQuery.useMultiquery", [True]
            ),
            "nQueries": trial.suggest_int("multiQuery.nQueries", 3, 8),
            "queryModelName": None,
        },
        "sourceIds": source_ids,
    }

    all_params = {
        "vectorSearch": {
            "topN": trial.suggest_int("vectorSearch.topN", 9, 20),
            "threshold": trial.suggest_float(
                "vectorSearch.threshold", 0.1, 0.7
            ),
            "searchByQuery": trial.suggest_categorical(
                "vectorSearch.searchByQuery", [True, False]
            ),
        },
        "reranker": {
            "useRerank": trial.suggest_categorical(
                "reranker.useRerank", [True, False]
            ),
            "topK": trial.suggest_categorical("reranker.topK", 4, 8),
            "rerankerMaxLength": trial.suggest_int(
                "reranker.rerankerMaxLength", [4096, 8192]
            ),
            "rerankerModel": "BAAI/bge-reranker-v2-m3",
        },
        "model": {
            "temperature": trial.suggest_float("model.temperature", 0.0, 1.0),
            "topK": trial.suggest_int("model.topK", 1, 100),
            "topP": trial.suggest_float("model.topP", 0.1, 1.0),
            "modelName": OLLAMA_BASE_MODEL,
        },
        "multiQuery": {
            "useMultiquery": trial.suggest_categorical(
                "multiQuery.useMultiquery", [True, False]
            ),
            "nQueries": trial.suggest_int("multiQuery.nQueries", 3, 8),
            "queryModelName": None,
        },
        "sourceIds": source_ids,
    }


if __name__ == "__main__":
    asyncio.run(main())
