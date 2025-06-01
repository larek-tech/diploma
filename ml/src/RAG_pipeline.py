from collections.abc import AsyncGenerator, AsyncIterator
import json
from document import Chunk

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
from ollama_client import AsyncOllamaClient, OllamaOptions
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

    async def _prepare_chunks(
        self, request: ml_pb2_model.ProcessQueryRequest
    ) -> list[Chunk]:
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

        document_chunk_dict = {}
        for question in questions:
            search_result = await self.data_client.vector_search(
                query=question,
                source_ids=request.sourceIds, # type: ignore
                top_k=request.scenario.vectorSearch.topN,
                threshold=request.scenario.vectorSearch.threshold,
                use_questions=request.scenario.vectorSearch.searchByQuery,
            )
            for document_chunk in search_result.chunks:
                document_chunk_dict[document_chunk.id] = {
                    "id": document_chunk.id,
                    "content": document_chunk.content,
                    "similarity": document_chunk.similarity,
                    "metadata": json.loads(document_chunk.metadata.decode("utf-8")),
                }
        chunks: list[Chunk] = [
            Chunk(
                content=chunk["content"],
                id=chunk["id"],
                similarity=chunk["similarity"],
                metadata=chunk["metadata"],
            )
            for chunk in sorted(
                document_chunk_dict.values(),
                key=lambda x: x["similarity"],
                reverse=True,
            )
            if "content" in chunk and "id" in chunk and "metadata" in chunk
        ]
        if not chunks:
            return []
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
        return chunks

    async def generate_stream(
        self,
        request: ml_pb2_model.ProcessQueryRequest,
    ) -> AsyncGenerator[tuple[str, list[Chunk]], None]:
        chunks = await self._prepare_chunks(request)

        stream = await self.ollama_client.generate(
            prompt=RAG_PROMPT.format(query=request.query.content, docs=chunks),
            model=request.scenario.model.modelName,
            stream_response=True,
            options=OllamaOptions(
                temperature=request.scenario.model.temperature,
                top_k=request.scenario.model.topK,
                top_p=request.scenario.model.topP,
                system=request.scenario.model.systemPrompt,
            )
        )
        if isinstance(stream, str):
            raise ValueError(
                "Expected an async iterator from the Ollama client."
            )
        async for token in stream: # type: ignore[assignment]
            yield token, chunks

    async def generate(
        self,
        request: ml_pb2_model.ProcessQueryRequest,
    ) -> tuple[str, list[Chunk]]:
        chunks = await self._prepare_chunks(request)

        responses = await self.ollama_client.generate(
            prompt=RAG_PROMPT.format(query=request.query.content, docs=chunks),
            model=request.scenario.model.modelName,
            stream_response=False,
            options=OllamaOptions(
                temperature=request.scenario.model.temperature,
                top_k=request.scenario.model.topK,
                top_p=request.scenario.model.topP,
                system=request.scenario.model.systemPrompt,
            )
        )
        if isinstance(responses, str):
            response = responses
        else:
            raise ValueError(
                "Expected a string response from the Ollama client."
            )
        return response, chunks
