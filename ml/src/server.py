import asyncio
from collections.abc import AsyncGenerator

import grpc
from grpc import aio

import ml.v1.model_pb2 as ml_pb2_model
import ml.v1.service_pb2_grpc as ml_pb2_grpc
from config import (
    DEFAULT_EMBEDER_MODEL,
    DEFAULT_REDIS_URL,
    ML_SERVICE_PORT,
)
from optuna_pipline import OptunaPipeline
from RAG_pipeline import RAGPipeline
from utils.logger import logger


class MLServiceServicer(ml_pb2_grpc.MLServiceServicer):
    def __init__(self) -> None:
        super().__init__()
        self.rag = RAGPipeline()
        self.optuna_optimizer = OptunaPipeline(
            redis_url=DEFAULT_REDIS_URL,
            embedings_model=DEFAULT_EMBEDER_MODEL,
        )

    async def ProcessQuery(  # noqa: N802
        self,
        request: ml_pb2_model.ProcessQueryRequest,
        context: aio.ServicerContext,
    ) -> AsyncGenerator[ml_pb2_model.ProcessQueryResponse]:
        client_ip = context.peer().split(":")[-1]
        request_id = f"{request.query.userId}-{hash(request.query.content)}"

        logger.info(
            f"New request [ID:{request_id}] from {client_ip}\n"
            f"User: {request.query.userId}\n"
            f"Content length: {len(request.query.content)}\n"
            f"Documents: {len(request.sourceIds)}"
        )
        try:
            async for token in self.rag.generate_stream(request=request):
                response = ml_pb2_model.ProcessQueryResponse(
                    chunk=ml_pb2_model.Chunk(content=f"{token}"),
                    sourceIds=request.sourceIds,
                )

                logger.debug(f"Sending chunk for request {request_id}")
                yield response
        except grpc.RpcError as e:
            logger.error(
                f"gRPC error processing request {request_id}:"
                f" {e.code()}: {e.details()}"
            )
            await context.abort(e.code(), e.details())
        except TimeoutError:
            logger.error(f"Timeout error processing request {request_id}")
            await context.abort(grpc.StatusCode.DEADLINE_EXCEEDED, "Timeout")

    async def ProcessScenario(  # noqa: N802
        self,
        request: ml_pb2_model.ProcessScenarioRequest,
        context: aio.ServicerContext,
    ) -> ml_pb2_model.ProcessScenarioResponse:
        client_ip = context.peer().split(":")[-1]

        logger.info(
            f"New request [For scenario {request.scenario.id} from {client_ip}"
            f"\nDocuments: {len(request.sourceIds)}"
        )
        background_tasks = set()
        task = asyncio.create_task(
            await self.optuna_optimizer.study(request.sourceIds)
        )
        background_tasks.add(task)
        task.add_done_callback(background_tasks.discard)
        return ml_pb2_model.ProcessScenarioResponse(
            status=True,
        )

async def serve() -> None:
    server = aio.server()
    ml_pb2_grpc.add_MLServiceServicer_to_server(MLServiceServicer(), server)
    server.add_insecure_port(f"0.0.0.0:{ML_SERVICE_PORT}")
    await server.start()
    logger.info(f"Server started on port {ML_SERVICE_PORT}")
    logger.info("Waiting for requests...")
    try:
        await server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("Shutting down server...")
        await server.stop(0)
        logger.info("Server stopped gracefully")


if __name__ == "__main__":
    asyncio.run(serve())
