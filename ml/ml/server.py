import grpc
from concurrent import futures
import time
import logging
import pb.ml.v1.service_pb2 as ml_pb2
import pb.ml.v1.service_pb2_grpc as ml_pb2_grpc
from utils.logger import logger


class MLServiceServicer(ml_pb2_grpc.MLServiceServicer):
    def ProcessQuery(self, request, context):
        client_ip = context.peer().split(":")[-1]
        request_id = f"{request.query.userId}-{hash(request.query.content)}"

        logger.info(
            f"New request [ID:{request_id}] from {client_ip}\n"
            f"User: {request.query.userId}\n"
            f"Content length: {len(request.query.content)}\n"
            f"Scenario: {ml_pb2.ScenarioType.Name(request.scenario.customType)}\n"
            f"Documents: {len(request.documentIds)}"
        )

        try:
            content = request.query.content
            chunks = 123

            logger.debug(f"Processing {len(chunks)} chunks for request {request_id}")

            for i, chunk in enumerate(chunks):
                if not context.is_active():
                    logger.warning(f"Connection closed early for request {request_id}")
                    break

                response = ml_pb2.ProcessQueryResponse(
                    chunk=ml_pb2.Chunk(content=f"[{i + 1}] {chunk}"),
                    documentIds=[f"processed_{doc}" for doc in request.documentIds],
                )

                logger.debug(
                    f"Sending chunk {i + 1}/{len(chunks)} for request {request_id}"
                )
                yield response
                time.sleep(0.5)

            logger.info(f"Successfully completed request {request_id}")

        except Exception as e:
            logger.error(
                f"Error processing request {request_id}: {str(e)}\n"
                f"Traceback: {logging.traceback.format_exc()}"
            )
            context.abort(grpc.StatusCode.INTERNAL, "Internal server error")


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    ml_pb2_grpc.add_MLServiceServicer_to_server(MLServiceServicer(), server)
    server.add_insecure_port("[::]:50051")
    server.start()

    logger.info("Server started on port 50051")
    logger.info("Worker threads: 10")
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
