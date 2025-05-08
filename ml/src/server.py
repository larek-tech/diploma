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
    DEFAULT_RERANKER_NAME,
    DEVICE,
    ML_SERVICE_PORT,
    OLLAMA_BASE_URL,
    RAG_PROMPT,
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
        self.reranker_model_name = DEFAULT_RERANKER_NAME
        self.reranker = Reranker(
            reranker_model_name=self.reranker_model_name,
            device=DEVICE,
        )

    def ProcessQuery(  # noqa: N802
        self,
        request: ml_pb2_model.ProcessQueryRequest,
        context: grpc.ServicerContext,
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
            questions = [request.query.content]
            if request.scenario.multiQuery.useMultiquery:
                questions += get_multi_questions(
                    client=self.ollama_client,
                    user_prompt=request.query.content,
                    n_questions=request.scenario.multiQuery.nQueryes,
                    model=request.scenario.multiQuery.queryModelName
                    if request.scenario.multiQuery.queryModelName
                    else request.scenario.model.modelName,
                )

            chunk_dict = {}
            for question in questions:
                search_result = self.data_client.vector_search(
                    query=question,
                    source_ids=request.sourceIds,
                    top_k=request.scenario.vectorSearch.topN,
                    threshold=request.scenario.vectorSearch.threshold,
                    use_questions=request.scenario.vectorSearch.searchByQuery,
                )
                for chunk in search_result.chunks:
                    chunk_dict[chunk.id] = {
                        "content": chunk.content,
                        "similarity": chunk.similarity,
                    }
            chunks = [
                chunk["content"]
                for chunk in sorted(
                    chunk_dict.values(),
                    key=lambda x: x["similarity"],
                    reverse=True,
                )
            ]
            if request.scenario.reranker.useRerank:
                if (
                    request.scenario.reranker.rerankerModel
                    != self.reranker_model_name
                ):
                    self.reranker_model_name = (
                        request.scenario.reranker.rerankerModel
                    )
                    self.reranker = Reranker(
                        reranker_model_name=self.reranker_model_name,
                        device=DEVICE,
                    )
                chunks = self.reranker.rerank_documents(
                    query=request.query.content,
                    documents=chunks,
                    top_k=request.scenario.reranker.topK,
                    max_length=request.scenario.reranker.rerankerMaxLenght,
                )
            for i, token in enumerate(
                self.ollama_client.generate(
                    prompt=RAG_PROMPT.format(
                        query=request.query.content, docs=chunks
                    ),
                    model=request.scenario.model.modelName,
                    stream=True,
                    temprature=request.scenario.model.temprature,
                    top_k=request.scenario.model.topK,
                    top_p=request.scenario.model.topP,
                    system=request.scenario.content,
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
                f"Error processing request {request_id}: {e!s}\n"
                f"Traceback: {logging.traceback.format_exc()}"
            )
            context.abort(grpc.StatusCode.INTERNAL, "Internal server error")


def serve() -> None:
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    ml_pb2_grpc.add_MLServiceServicer_to_server(MLServiceServicer(), server)
    server.add_insecure_port(f"0.0.0.0:{ML_SERVICE_PORT}")
    server.start()

    logger.info(f"Server started on port {ML_SERVICE_PORT}")
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
