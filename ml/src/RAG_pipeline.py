from collections.abc import AsyncGenerator

import ml.v1.model_pb2 as ml_pb2_model
from config import (
    DATA_SERVICE_HOST,
    DATA_SERVICE_PORT,
    DEFAULT_RERANKER_NAME,
    DEVICE,
    OLLAMA_BASE_URL,
    RAG_PROMPT,
)
from data_client import AsyncDataServiceClient
from multi_query import get_multi_questions
from ollama_client import AsyncOllamaClient
from rerank import Reranker


class RAGPipeline:
    def __init__(self) -> None:
        self.ollama_client = AsyncOllamaClient(
            base_url=OLLAMA_BASE_URL,
        )
        self.data_client = AsyncDataServiceClient(
            host=DATA_SERVICE_HOST, port=DATA_SERVICE_PORT
        )
        self.reranker_model_name = DEFAULT_RERANKER_NAME
        self.reranker = Reranker(
            reranker_model_name=self.reranker_model_name,
            device=DEVICE,
        )

    async def generate(
        self,
        request: ml_pb2_model.ProcessQueryRequest,
    ) -> AsyncGenerator[tuple[ml_pb2_model.ProcessQueryResponse, list[str]], None]:
        questions = [request.query.content]
        if request.scenario.multiQuery.useMultiquery:
            questions += await get_multi_questions(
                client=self.ollama_client,
                user_prompt=request.query.content,
                n_questions=request.scenario.multiQuery.nQueries,
                model=request.scenario.multiQuery.queryModelName
                if request.scenario.multiQuery.queryModelName
                else request.scenario.model.modelName,
            )

        chunk_dict = {}
        for question in questions:
            search_result = await self.data_client.vector_search(
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
                max_length=request.scenario.reranker.rerankerMaxLength,
            )

        stream = await self.ollama_client.generate(
            prompt=RAG_PROMPT.format(query=request.query.content, docs=chunks),
            model=request.scenario.model.modelName,
            stream=True,
            temprature=request.scenario.model.temperature,
            top_k=request.scenario.model.topK,
            top_p=request.scenario.model.topP,
            system=request.scenario.model.systemPrompt,
        )
        # TODO: Refactor
        async for token in stream:
            yield token, chunks
