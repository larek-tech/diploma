import asyncio
import grpc
import uuid

# Assuming your generated protobuf files are in a directory structure like:
# ml/
#   v1/
#     model_pb2.py
#     service_pb2_grpc.py
# You might need to adjust the import paths based on your project structure
# and how you've added them to your PYTHONPATH.
# For example, if 'ml' is a top-level directory in your project:
import ml.v1.model_pb2 as ml_pb2_model
import ml.v1.service_pb2_grpc as ml_pb2_grpc

# Configuration (matches server's default port, adjust if needed)
ML_SERVICE_HOST = "localhost"
ML_SERVICE_PORT = 8888  # Make sure this matches ML_SERVICE_PORT in your server config
MODEL = 'hf.co/t-tech/T-lite-it-1.0-Q8_0-GGUF:Q8_0'
TEST_QUERY = "docker"

async def run_process_query(stub: ml_pb2_grpc.MLServiceStub):
    print("--- Calling ProcessQuery ---")
    request = ml_pb2_model.ProcessQueryRequest(
        query=ml_pb2_model.Query(
            userId=12,
            content= TEST_QUERY,
        ),
        sourceIds=["52eccef1-8f92-4cc2-812c-dc441a236829"],
        scenario=ml_pb2_model.Scenario(
            model=ml_pb2_model.LlmModel(modelName=MODEL),
            vectorSearch=ml_pb2_model.VectorSearch(
                topN=5,
                threshold=0.1,
            ),
        )
        # Example source IDs
        # You can also populate modelParams here if needed
        # modelParams=ml_pb2_model.ModelParams(...)
    )
    try:
        async for response in stub.ProcessQuery(request):
            print(f"Received response for query: {response}")
            if response.HasField("chunk"):
                print(f"Received chunk: {response.chunk.content}")
            if response.sourceIds:
                print(f"Received source IDs at the end: {list(response.sourceIds)}")
        print("--- ProcessQuery Finished ---")
    except grpc.aio.AioRpcError as e:
        print(f"ProcessQuery RPC failed: {e.code()} - {e.details()}")


async def run_process_first_query(stub: ml_pb2_grpc.MLServiceStub):
    print("\n--- Calling ProcessFirstQuery ---")
    request = ml_pb2_model.ProcessFirstQueryRequest(
        query="What is the capital of France?"
    )
    try:
        response = await stub.ProcessFirstQuery(request)
        print(f"ProcessFirstQuery response: {response.query}")
        print("--- ProcessFirstQuery Finished ---")
    except grpc.aio.AioRpcError as e:
        print(f"ProcessFirstQuery RPC failed: {e.code()} - {e.details()}")


async def run_get_default_params(stub: ml_pb2_grpc.MLServiceStub):
    print("\n--- Calling GetDefaultParams ---")
    request = ml_pb2_model.GetDefaultParamsRequest() # Empty request
    try:
        response = await stub.GetDefaultParams(request)
        print(f"GetDefaultParams response:")
        print(f"  MultiQuery: use={response.multiQuery.useMultiquery}, nQueries={response.multiQuery.nQueries}")
        print(f"  Reranker: use={response.reranker.useRerank}, topK={response.reranker.topK}, model={response.reranker.rerankerModel}")
        print(f"  VectorSearch: topN={response.vectorSearch.topN}, threshold={response.vectorSearch.threshold}, byQuery={response.vectorSearch.searchByQuery}")
        print(f"  LLM Model: name={response.model.modelName}, temp={response.model.temperature}")
        print("--- GetDefaultParams Finished ---")
    except grpc.aio.AioRpcError as e:
        print(f"GetDefaultParams RPC failed: {e.code()} - {e.details()}")


async def run_get_optimal_params(stub: ml_pb2_grpc.MLServiceStub):
    print("\n--- Calling GetOptimalParams ---")
    request = ml_pb2_model.GetOptimalParamsRequest(
        sourceIds=["doc_A", "doc_B", "doc_C"] # Example source IDs
    )
    try:
        response = await stub.GetOptimalParams(request)
        print(f"GetOptimalParams response:")
        # Assuming it returns ModelParams, similar to GetDefaultParams
        print(f"  MultiQuery: use={response.multiQuery.useMultiquery}, nQueries={response.multiQuery.nQueries}")
        print(f"  Reranker: use={response.reranker.useRerank}, topK={response.reranker.topK}, model={response.reranker.rerankerModel}")
        print(f"  VectorSearch: topN={response.vectorSearch.topN}, threshold={response.vectorSearch.threshold}, byQuery={response.vectorSearch.searchByQuery}")
        print(f"  LLM Model: name={response.model.modelName}, temp={response.model.temperature}")
        print("--- GetOptimalParams Finished ---")
    except grpc.aio.AioRpcError as e:
        print(f"GetOptimalParams RPC failed: {e.code()} - {e.details()}")


async def main():
    server_address = f"{ML_SERVICE_HOST}:{ML_SERVICE_PORT}"
    async with grpc.aio.insecure_channel(server_address) as channel:
        stub = ml_pb2_grpc.MLServiceStub(channel)

        await run_process_query(stub)
        # await run_process_first_query(stub)
        # await run_get_default_params(stub)
        # await run_get_optimal_params(stub)


if __name__ == "__main__":
    asyncio.run(main())