import asyncio
from concurrent.futures import ThreadPoolExecutor
import json

import optuna
import redis.asyncio as redis

import ml.v1.model_pb2 as ml_pb2_model
from RAG_pipeline import RAGPipeline
from utils.logger import logger


async def load_all_qa_samples_from_redis(
    redis_url: str = "redis://localhost", key_prefix: str = "qa:"
) -> list[dict]:
    redis_client = redis.from_url(
        redis_url, encoding="utf-8", decode_responses=True
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
    entry: dict, params: dict
) -> ml_pb2_model.ProcessQueryRequest:
    return ml_pb2_model.ProcessQueryRequest(
        query=ml_pb2_model.Query(
            content=entry["question"], userId="test_user"
        ),
        scenario=ml_pb2_model.Scenario(
            model=ml_pb2_model.LLMMdel(
                modelName=params["model"]["modelName"],
                temprature=params["model"]["temprature"],
                topK=params["model"]["topK"],
                topP=params["model"]["topP"],
            ),
            content=params["content"],
            multiQuery=ml_pb2_model.MultiQuery(
                useMultiquery=params["multiQuery"]["useMultiquery"],
                nQueryes=params["multiQuery"]["nQueryes"],
                queryModelName=params["multiQuery"]["queryModelName"] or "",
            ),
            vectorSearch=ml_pb2_model.VectorSearch(
                topN=params["vectorSearch"]["topN"],
                threshold=params["vectorSearch"]["threshold"],
                searchByQuery=params["vectorSearch"]["searchByQuery"],
            ),
            reranker=ml_pb2_model.Reranker(
                useRerank=params["reranker"]["useRerank"],
                topK=params["reranker"]["topK"],
                rerankerMaxLenght=params["reranker"]["rerankerMaxLenght"],
                rerankerModel=params["reranker"]["rerankerModel"],
            ),
        ),
        sourceIds=params["sourceIds"],
    )
def calculate_score(model_answers: list[dict], test_dataset: list[dict]) -> float:
    # Placeholder for actual scoring logic
    # Compare generated answers with expected answers and compute a score
    return 0.5 


async def make_test_request(pipeline: RAGPipeline, params: dict, test_dataset: list[dict[str, str]]) -> list[dict]:
    results = []
    for entry in test_dataset:
        request = build_request(entry, params)
        output = []
        try:
            chunks = []
            async for token, request_chunks in pipeline.generate(request):
                chunks = request_chunks
                output.append(  # noqa: PERF401
                    token.chunk.content
                    if hasattr(token, "chunk")
                    else str(token)
                )
            generated_answer = "".join(output)

            results.append(
                {
                    "question": entry["question"],
                    "answer": entry["answer"],
                    "contexts": entry["contexts"],
                    "generated_answer": generated_answer,
                    "retrivment_content": chunks,
                }
            )
        except Exception as e:
            logger.info(f"Ошибка при обработке запроса: {e}")
            results.append(
                {
                    "question": entry["question"],
                    "answer": entry["answer"],
                    "contexts": entry["contexts"],
                    "generated_answer": "",
                }
            )

    return results


async def evaluate_rag_pipeline(params: dict, test_dataset: list[dict[str, str]]) -> float:
    pipeline = RAGPipeline()

    # model_asnwers = make_test_request(pipeline, params, test_dataset)
    try:
        model_answers = await make_test_request(pipeline, params, test_dataset)
        # Implement your scoring logic here based on model_answers and test_dataset
        # Example: Calculate accuracy or another metric
        score = calculate_score(model_answers, test_dataset)
        return score
    except Exception as e:
        logger.error(f"Error during evaluation: {e}")
        return 0.0


def objective(trial, test_dataset: list[dict[str]], source_ids: list[str]) -> float:
    params = {
        "vectorSearch": {
            "topN": trial.suggest_int("vectorSearch.topN", 3, 20),
            "threshold": trial.suggest_float(
                "vectorSearch.threshold", 0.1, 0.9
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
            "rerankerMaxLenght": trial.suggest_int(
                "reranker.rerankerMaxLenght", 128, 1024
            ),
            "rerankerModel": "your_reranker_model_name",
        },
        "model": {
            "temprature": trial.suggest_float("model.temprature", 0.0, 1.0),
            "topK": trial.suggest_int("model.topK", 1, 100),
            "topP": trial.suggest_float("model.topP", 0.1, 1.0),
            "modelName": "your_model_name",
        },
        "multiQuery": {
            "useMultiquery": trial.suggest_categorical(
                "multiQuery.useMultiquery", [True, False]
            ),
            "nQueryes": trial.suggest_int("multiQuery.nQueryes", 1, 5),
            "queryModelName": None,
        },
        "content": "You are a helpful assistant",  # system prompt
        "sourceIds": source_ids,
    }
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    try:
        return loop.run_until_complete(evaluate_rag_pipeline(params, test_dataset))
    finally:
        loop.close()
    # return asyncio.run(evaluate_rag_pipeline(params, test_dataset))


async def main() -> None:
    loop = asyncio.get_event_loop()
    study = optuna.create_study(direction="maximize")
    source_ids=["a6bfe96f-45bd-4e4b-8e6f-2c2ef53ca280"]
    test_dataset = await load_all_qa_samples_from_redis(
        redis_url="redis://localhost", key_prefix=source_ids[0]
    )
    with ThreadPoolExecutor(max_workers=1) as executor:
        study = optuna.create_study(direction="maximize")
        study.optimize(
            lambda trial: objective(trial, test_dataset, source_ids),
            n_trials=50,
            n_jobs=1,  # Important for async compatibility
            callbacks=[lambda study, trial: logger.info(f"Trial {trial.number} completed")],
            gc_after_trial=True
        )

    print("Best parameters:", study.best_params)

    # print("Лучшие параметры:", study.best_params)


# Запуск Optuna
if __name__ == "__main__":
    asyncio.run(main())

