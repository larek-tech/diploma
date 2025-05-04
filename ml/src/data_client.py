import grpc

import data.v1.model_pb2 as pb
import data.v1.service_pb2_grpc as pb_grpc
from utils.logger import logger


class DataServiceClient:
    """Клиент для взаимодействия с сервисом данных через gRPC.

    Parameters
    ----------
    host : str, optional
        Адрес сервера. По умолчанию "localhost".
    port : int, optional
        Порт сервера.
    """

    def __init__(self, host: str = "localhost", port: int = 9990) -> None:
        self.channel = grpc.insecure_channel(f"{host}:{port}")
        self.stub = pb_grpc.DataServiceStub(self.channel)

    def close(self) -> None:
        """Закрывает соединение с сервером по gRPC."""
        self.channel.close()

    def vector_search(
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
        return self.stub.VectorSearch(request)


if __name__ == "__main__":
    client = DataServiceClient()
    try:
        response = client.vector_search(
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
        client.close()
