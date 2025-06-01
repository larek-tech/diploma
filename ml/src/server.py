import asyncio
from collections.abc import AsyncGenerator
from concurrent import futures
import json
import grpc
from grpc import aio

import ml.v1.model_pb2 as ml_pb2_model
import ml.v1.service_pb2_grpc as ml_pb2_grpc
from config import (
    DEFAULT_EMBEDER_MODEL,
    DEFAULT_REDIS_URL,
    DEFAULT_RERANKER_NAME,
    FIRST_MESSAE_PROMPT,
    JSON_SCHEMA,
    ML_SERVICE_PORT,
    OLLAMA_BASE_MODEL,
)
from optuna_pipline import OptunaPipeline
from RAG_pipeline import RAGPipeline
from sample_generate import generate_dataset
from utils.logger import logger
from ollama_client import OllamaOptions
from opentelemetry import trace, baggage, propagate
from opentelemetry.trace.status import StatusCode, Status 

from typing import Optional
from tracing import *
import functools # Add this import
import inspect   # Add this impor

tracer = tracerProvider.get_tracer("ml_service")

def trace_grpc_method(func):
    if inspect.isasyncgenfunction(func):
        @functools.wraps(func)
        async def async_gen_wrapper(self, request, context: aio.ServicerContext, *args, **kwargs):
            class_name = self.__class__.__name__
            method_name = func.__name__
            span_name = f"rpc/{class_name}/{method_name}"
            invocation_md = context.invocation_metadata()
            carrier = {k: v for k, v in invocation_md} if invocation_md else {}
            parent_context = propagate.extract(carrier)
            print(f"Invocation metadata: {carrier}")

            with tracer.start_as_current_span(span_name, context=parent_context, kind=trace.SpanKind.SERVER) as span:
                # The line "metadata = context.invocation_metadata()" previously here is removed as it's redundant.
                span.set_attribute("rpc.system", "grpc")
                span.set_attribute("rpc.service", class_name)
                span.set_attribute("rpc.method", method_name)
                if context and hasattr(context, 'peer') and context.peer():
                    span.set_attribute("net.peer.name", context.peer())

                try:
                    async for item in func(self, request, context, *args, **kwargs):
                        yield item
                    span.set_status(Status(StatusCode.OK))
                except grpc.RpcError as e:
                    if e.code() is not None:
                        span.set_attribute("rpc.grpc.status_code", e.code().value[0])
                    span.record_exception(e)
                    span.set_status(Status(StatusCode.ERROR, description=f"gRPC Error: {e.details()}"))
                    raise
                except Exception as e:
                    span.record_exception(e)
                    span.set_status(Status(StatusCode.ERROR, description=str(e)))
                    raise
        return async_gen_wrapper
    else:  # func is a regular async function
        @functools.wraps(func)
        async def regular_async_wrapper(self, request, context: aio.ServicerContext, *args, **kwargs):
            class_name = self.__class__.__name__
            method_name = func.__name__
            span_name = f"rpc/{class_name}/{method_name}"

            invocation_md = context.invocation_metadata()
            carrier = {k: v for k, v in invocation_md} if invocation_md else {}
            print(f"Invocation metadata: {carrier}")
            parent_context = propagate.extract(carrier)


            with tracer.start_as_current_span(span_name, context=parent_context, kind=trace.SpanKind.SERVER) as span:
                span.set_attribute("rpc.system", "grpc")
                span.set_attribute("rpc.service", class_name)
                span.set_attribute("rpc.method", method_name)
                if context and hasattr(context, 'peer') and context.peer():
                    span.set_attribute("net.peer.name", context.peer())

                try:
                    result = await func(self, request, context, *args, **kwargs)
                    span.set_status(Status(StatusCode.OK))
                    return result
                except grpc.RpcError as e:
                    if e.code() is not None:
                        span.set_attribute("rpc.grpc.status_code", e.code().value[0])
                    span.record_exception(e)
                    span.set_status(Status(StatusCode.ERROR, description=f"gRPC Error: {e.details()}"))
                    raise
                except Exception as e:
                    span.record_exception(e)
                    span.set_status(Status(StatusCode.ERROR, description=str(e)))
                    raise
        return regular_async_wrapper

class MLServiceServicer(ml_pb2_grpc.MLServiceServicer):
    def __init__(self) -> None:
        super().__init__()
        self.rag = RAGPipeline()
        self.optuna_optimizer = OptunaPipeline(
            redis_url=DEFAULT_REDIS_URL, embedings_model=DEFAULT_EMBEDER_MODEL
        )

    @trace_grpc_method
    async def ProcessQuery(  # noqa: N802
        self,
        request: ml_pb2_model.ProcessQueryRequest,
        context: aio.ServicerContext,
    ) -> AsyncGenerator[ml_pb2_model.ProcessQueryResponse]:
        metadata = context.invocation_metadata()
        
        client_ip = context.peer().split(":")[-1]
        request_id = f"{request.query.userId}-{hash(request.query.content)}"

        logger.info(
            f"New request [ID:{request_id}] from {client_ip}\n"
            f"User: {request.query.userId}\n"
            f"Content length: {len(request.query.content)}\n"
            f"Documents: {len(request.sourceIds)}"
        )
        try:
            # chunks = None
            chunk_dict = {}
            async for token, chunk in self.rag.generate_stream(
                request=request
            ):
                response = ml_pb2_model.ProcessQueryResponse(
                    chunk=ml_pb2_model.Chunk(content=f"{token}"),
                )
                if chunk:
                    for c in chunk:
                        if "id" in c.keys() and "metadata" in c.keys():
                            if  c["id"] not in chunk_dict:
                                chunk_dict[c["id"]] = c
                            else:
                                chunk_dict[c["id"]]["metadata"].update(c["metadata"])

                logger.debug(f"Sending chunk for request {request_id}")
                yield response
            meta = [
                json.dumps({
                    "metadata": chunk["metadata"],
                    "id": chunk["id"] if "id" in chunk else "",
                })
                for chunk in chunk_dict.values()
                if "metadata" in chunk and "id" in chunk
            ] if len(chunk_dict) > 0 else []
            print(len(meta))
            yield ml_pb2_model.ProcessQueryResponse(sourceIds=meta)
        except grpc.RpcError as e:
            logger.error(
                f"gRPC error processing request {request_id}:"
                f" {e.code()}: {e.details()}"
            )
            await context.abort(e.code(), e.details())
        except TimeoutError:
            logger.error(f"Timeout error processing request {request_id}")
            context.abort(grpc.StatusCode.DEADLINE_EXCEEDED, "Timeout")

    @trace_grpc_method
    async def ProcessFirstQuery(  # noqa: N802
        self,
        request: ml_pb2_model.ProcessFirstQueryRequest,
        context: aio.ServicerContext,
    ) -> ml_pb2_model.ProcessFirstQueryResponse:
        client_ip = context.peer().split(":")[-1]

        logger.info(f"New request [From {client_ip}\nQuery: {request.query}\n")
        try:
            json_response = await self.rag.ollama_client.generate(
                prompt=FIRST_MESSAE_PROMPT.format(message=request.query),
                model=OLLAMA_BASE_MODEL,
                stream_response=False,
                options=OllamaOptions(
                    format=JSON_SCHEMA,
                )
            )
            print(json_response)
            # response = json.loads(json_response) #FIXME: structured output
            return ml_pb2_model.ProcessFirstQueryResponse(query=json_response)
        except grpc.RpcError as e:
            logger.error(
                f"gRPC error processing request: {e.code()}: {e.details()}"
            )
            await context.abort(e.code(), e.details())
        except TimeoutError:
            logger.error("Timeout error processing request")
            await context.abort(grpc.StatusCode.DEADLINE_EXCEEDED, "Timeout")

    async def GetDefaultParams(  # noqa: N802
        self,
        request,
        context: aio.ServicerContext,
    ) -> ml_pb2_model.ModelParams:
        return ml_pb2_model.ModelParams(
            multiQuery=ml_pb2_model.MultiQuery(
                useMultiquery=True,
                nQueries=3,
            ),
            reranker=ml_pb2_model.Reranker(
                useRerank=True,
                topK=5,
                rerankerMaxLength=8192,
                rerankerModel=DEFAULT_RERANKER_NAME,
            ),
            vectorSearch=ml_pb2_model.VectorSearch(
                topN=10,
                threshold=0.1,
                searchByQuery=True,
            ),
            model=ml_pb2_model.LlmModel(
                modelName=OLLAMA_BASE_MODEL,
                temperature=0.7,
                topK=5,
                topP=0.9,
                systemPrompt="",
            ),
        )

    async def GetOptimalParams(  # noqa: N802
        self,
        request: ml_pb2_model.GetOptimalParamsRequest,
        context: aio.ServicerContext,
    ) -> ml_pb2_model.ModelParams:
        client_ip = context.peer().split(":")[-1]

        logger.info(
            f"New request [From {client_ip}"
            f"\nDocuments: {len(request.sourceIds)}"
        )
        # await generate_dataset(request.sourceIds, self.rag.data_client)
        return await self.optuna_optimizer.study(
            source_ids=request.sourceIds,
        )


async def serve() -> None:
    server = aio.server(futures.ThreadPoolExecutor(max_workers=10))
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
