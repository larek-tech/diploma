import logging
import time
from collections.abc import Iterator
from concurrent import futures

import grpc

import ml.v1.model_pb2 as ml_pb2_model
import ml.v1.service_pb2_grpc as ml_pb2_grpc
from config import (
    DATA_SERVICE_HOST,
    DATA_SERVICE_PORT,
    DEVICE,
    OLLAMA_BASE_URL,
    RAG_PROMPT,
    THRESHOLD,
    TOP_K,
)
from data_client import DataServiceClient
from multi_query import get_multi_questions
from ollama_client import OllamaClient
from rerank import Reranker
from utils.logger import logger


class MLServiceServicer(ml_pb2_grpc.MLServiceServicer):
    def __init__(self) -> None:
        super().__init__()
        self.ollama_client = OllamaClient(
            base_url=OLLAMA_BASE_URL,
        )
        self.data_client = DataServiceClient(
            host=DATA_SERVICE_HOST, port=DATA_SERVICE_PORT
        )
        self.reranker = Reranker(
            reranker_model_name="BAAI/bge-reranker-v2-m3",
            max_length=2048,
            device=DEVICE,
        )

    def ProcessQuery(  # noqa: N802
        self, request: ml_pb2_model.ProcessQueryRequest, context
    ) -> Iterator[ml_pb2_model.ProcessQueryResponse]:
        client_ip = context.peer().split(":")[-1]
        request_id = f"{request.query.userId}-{hash(request.query.content)}"

        logger.info(
            f"New request [ID:{request_id}] from {client_ip}\n"
            f"User: {request.query.userId}\n"
            f"Content length: {len(request.query.content)}\n"
            f"Scenario: {ml_pb2_model.ScenarioType.Name(request.scenario.customType)}\n"  # noqa: E501
            f"Documents: {len(request.sourceIds)}"
        )

        try:
            content = request.query.content
            perefrazing_questions = get_multi_questions(
                client=self.ollama_client,
                user_prompt=content,
                n_questions=5,
                model="hf.co/yandex/YandexGPT-5-Lite-8B-instruct-GGUF:Q4_K_M",
            )

            chunk_dict = {}
            for question in perefrazing_questions:
                search_result = self.data_client.vector_search(
                    query=question,
                    source_ids=request.sourceIds,
                    top_k=TOP_K,
                    threshold=THRESHOLD,
                    use_questions=False,
                )
                for chunk in search_result.chunks:
                    # TODO: Возможность добавить поля по типу similaryty
                    chunk_dict[chunk.id] = chunk.content
            chunks = list(chunk_dict.values())
            rerunked_chunks = self.reranker.rerank_documents(
                query=request.query.content,
                documents=chunks,
            )
            for i, token in enumerate(
                self.ollama_client.generate(
                    prompt=RAG_PROMPT.format(
                        query=request.query.content, docs=rerunked_chunks
                    ),
                    model="hf.co/yandex/YandexGPT-5-Lite-8B-instruct-GGUF:Q4_K_M",
                    stream=True,
                )
            ):
                response = ml_pb2_model.ProcessQueryResponse(
                    chunk=ml_pb2_model.Chunk(content=f"{token}"),
                    sourceIds=request.sourceIds,
                )

                logger.debug(f"Sending chunk {i + 1} for request {request_id}")
                yield response
        except Exception as e:
            logger.error(
                f"Error processing request {request_id}: {str(e)}\n"
                f"Traceback: {logging.traceback.format_exc()}"
            )
            context.abort(grpc.StatusCode.INTERNAL, "Internal server error")


def serve() -> None:
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    ml_pb2_grpc.add_MLServiceServicer_to_server(MLServiceServicer(), server)
    server.add_insecure_port("0.0.0.0:50051")
    server.start()

    logger.info("Server started on port 50051")
    logger.info("Waiting for requests...")

    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        logger.info("Shutting down server...")
        server.stop(0)
        logger.info("Server stopped gracefully")


if __name__ == "__main__":
    serve()
