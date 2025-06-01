import asyncio

from grpc import aio
import random

import data.v1.model_pb2 as pb
import data.v1.service_pb2_grpc as pb_grpc
from config import logger

from opentelemetry import propagate, trace
from opentelemetry.trace.status import StatusCode, Status
from tracing import TracerProvider, tracer


class AsyncDataServiceClient:
    """Асинхронный клиент для взаимодействия с сервисом данных через gRPC.

    Parameters
    ----------
    host : str, optional
        Адрес сервера. По умолчанию "localhost".
    port : int, optional
        Порт сервера.
    """

    def __init__(self, host: str = "localhost", port: int = 9990) -> None:
        self.channel = aio.insecure_channel(f"{host}:{port}")
        self.stub = pb_grpc.DataServiceStub(self.channel)

    async def close(self) -> None:
        """Закрывает соединение с сервером по gRPC."""
        await self.channel.close()

    async def vector_search(
        self,
        query: str,
        source_ids: list[str],
        top_k: int,
        threshold: float,
        *,
        use_questions: bool,
    ) -> pb.VectorSearchResponse:
        """Выполняет векторный поиск на сервере.

        Parameters
        ----------
        query : str
            Строка поискового запроса
        source_ids : List[str]
            Список идентификаторов источников для ограничения поиска
        top_k : int
            Максимальное количество возвращаемых результатов
        threshold : float
            Минимальный порог схожести для результатов
        use_questions : bool
            Использовать поиск по вопросам

        Returns
        -------
        pb.VectorSearchResponse
            Ответ сервера с результатами поиска
        """
        with tracer.start_span(
            "AsyncDataServiceClient.vector_search",
            attributes={
                "query": query,
                "sourceIds": str(source_ids),
                "topK": top_k,
                "threshold": threshold,
                "useQuestions": use_questions,
            },
        ) as span:
            span.set_attribute("query", query)
            span.set_attribute("sourceIds", str(source_ids))
            span.set_attribute("topK", top_k)
            span.set_attribute("threshold", threshold)
            span.set_attribute("useQuestions", use_questions)
            request = pb.VectorSearchRequest(
                query=query,
                sourceIds=source_ids,
                topK=top_k,
                threshold=threshold,
                useQuestions=use_questions,
            )
            carrier = {}
            propagate.inject(carrier)


            grpc_metadata = []
            for key, value in carrier.items():
                grpc_metadata.append((key, value))

            current_span_context = span.get_span_context()
            if current_span_context.is_valid:
                trace_id_hex = format(current_span_context.trace_id, 'x')
                found_x_trace_id = False
                for k, _ in grpc_metadata:
                    if k == 'x-trace-id':
                        found_x_trace_id = True
                        break
                if not found_x_trace_id:
                     grpc_metadata.append(('x-trace-id', trace_id_hex))
            try:
                response = await self.stub.VectorSearch(request, metadata=grpc_metadata)
                span.set_status(Status(StatusCode.OK))
                # You can add attributes from the response if needed
                span.set_attribute("dataservice.response.num_chunks", len(response.chunks))
                return response
            except aio.AioRpcError as e:
                span.record_exception(e)
                span.set_status(Status(StatusCode.ERROR, description=f"gRPC error: {e.details()}"))
                span.set_attribute("rpc.grpc.status_code", e.code().value[0]) # Record gRPC status code
                raise
            except Exception as e:
                span.record_exception(e)
                span.set_status(Status(StatusCode.ERROR, description=str(e)))
                raise



    async def get_documents(
        self, source_id: str, n_samples: int
    ) -> pb.GetDocumentsOut:
        with tracer.start_span(
            "AsyncDataServiceClient.get_documents",
            attributes={
                "sourceId": source_id,
                "nSamples": n_samples,
            },
        ) as span:
            span.set_attribute("sourceId", source_id)
            span.set_attribute("nSamples", n_samples)
            """Получает случайные документы по source_id."""
            size = 10

            initial_request = pb.GetDocumentsIn(sourceId=source_id, size=1, page=1)
            initial_response = await self.stub.GetDocuments(initial_request)
            total_docs = initial_response.total

            if n_samples > total_docs:
                raise ValueError("Запрошено больше документов, чем существует")

            import random
            random_indices = random.sample(range(total_docs), n_samples)

            from collections import defaultdict
            page_offsets = defaultdict(list)
            for idx in random_indices:
                page = (idx // size) + 1
                offset = idx % size
                page_offsets[page].append(offset)

            tasks = [
                self.stub.GetDocuments(pb.GetDocumentsIn(sourceId=source_id, size=size, page=page))
                for page in page_offsets.keys()
            ]
            responses = await asyncio.gather(*tasks)

            selected_docs = []
            for page, resp in zip(page_offsets.keys(), responses):
                valid_offsets = [offset for offset in page_offsets[page] if offset < len(resp.documents)]
                selected_docs.extend(resp.documents[offset] for offset in valid_offsets)

            return pb.GetDocumentsOut(
                documents=selected_docs,
                total=len(selected_docs),
                page=0,
                size=len(selected_docs)
            )


async def main() -> None:
    client = AsyncDataServiceClient()
    try:
        response = await client.vector_search(
            query="Java разработчик",
            source_ids=["a6bfe96f-45bd-4e4b-8e6f-2c2ef53ca280"],
            top_k=5,
            threshold=0.1,
            use_questions=False,
        )

        logger.info(f"Received {len(response.chunks)} chunks:")
        for i, chunk in enumerate(response.chunks):
            logger.info(f"\nChunk {i + 1}:")
            logger.info(f"ID: {chunk.id}")
            logger.info(f"Index: {chunk.index}")
            logger.info(f"Similarity: {chunk.similarity:.2f}")
            logger.info(f"Content: {chunk.content}")
    finally:
        await client.close()


if __name__ == "__main__":
    asyncio.run(main())
