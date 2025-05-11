import asyncio

from grpc import aio

import data.v1.model_pb2 as pb
import data.v1.service_pb2_grpc as pb_grpc
from config import logger


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
        request = pb.VectorSearchRequest(
            query=query,
            sourceIds=source_ids,
            topK=top_k,
            threshold=threshold,
            useQuestions=use_questions,
        )
        return await self.stub.VectorSearch(request)

    async def get_documents(
        self, source_id: str, size: int = 10, page: int = 1
    ) -> pb.GetDocumentsOut:
        """Получает документы по source_id с пагинацией.

        Parameters
        ----------
        source_id : str
            Идентификатор источника.
        size : int, optional
            Количество документов на странице.
        page : int, optional
            Номер страницы.

        Returns
        -------
        pb.GetDocumentsOut
            Ответ с документами от сервера.
        """
        request = pb.GetDocumentsIn(sourceId=source_id, size=size, page=page)
        return await self.stub.GetDocuments(request)


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
