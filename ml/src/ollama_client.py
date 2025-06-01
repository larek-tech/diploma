import asyncio

from collections.abc import AsyncIterator
from typing import Any, Union, Dict
import ollama


from config import OLLAMA_BASE_MODEL, OLLAMA_BASE_URL, NUM_CTX, FIRST_MESSAE_PROMPT, JSON_SCHEMA
from utils.logger import logger


class OllamaOptions (Dict[str, Any]):
    temperature: float = 0.7
    top_k: int = 40
    top_p: float = 0.95
    system= "You are a helpful assistant."
    format: dict | None = None

class AsyncOllamaClient:
    def __init__(self, base_url: str = OLLAMA_BASE_URL) -> None:
        self.base_url = base_url
        self.client = ollama.AsyncClient(host=base_url)

    async def generate(
        self,
        prompt: str,
        model: str,
        *,
        stream_response: bool = False,
        options: OllamaOptions | None = None,
    ) -> Union[str, AsyncIterator[str]]:

        ollama_options = {
            "num_ctx": int(NUM_CTX) if NUM_CTX else 64000,
            "temperature": options.temperature if options else 0.7,
            "top_k": options.top_k if options else 40,
            "top_p": options.top_p if options else 0.95,
            "system": options.system if options else "You are a helpful assistant.",
        }


        if stream_response:
            return self._streaming(
                model=model,
                prompt=prompt,
                options=ollama_options,
                format=options.format if options else None,
            )
        else:
            return await self._generate(
                model=model,
                prompt=prompt,
                options=ollama_options,
                format=options.format if options else None,
            )

    async def _generate(
        self,
        model: str,
        prompt: str,
        options: dict[str, Any],
        format: dict | None = None,
        **kwargs: Any,
    ) -> str:
        if format:
            response = await self.client.generate(
                    model=model,
                    prompt=prompt,
                    stream=False,
                    format=format,
                    options=options,
                )
        else:
            response = await self.client.generate(
                model=model,
                prompt=prompt,
                stream=False,
                options=options,
            )
        return response.response if response else ""

    def _streaming(
            self,
            model: str,
            prompt: str,
            options: dict[str, Any],
            format: dict | None = None,
        ) -> AsyncIterator[str]:

            async def _stream_helper() -> AsyncIterator[str]:
                if format:
                    stream = await self.client.generate(
                        model=model,
                        prompt=prompt,
                        stream=True,
                        format=format,
                        options=options,
                    )
                else:
                    stream = await self.client.generate(
                        model=model,
                        prompt=prompt,
                        stream=True,
                        options=options,
                    )
                async for chunk in stream:
                    text = chunk.response
                    if text:
                        yield text
            return _stream_helper()


async def main() -> None:
    client = AsyncOllamaClient()
    test_prompt: str = "Привет, как дела? Напиши мне что-нибудь интересное."
    test_query: str = "Какой сегодня день?"
    test_model: str = OLLAMA_BASE_MODEL if OLLAMA_BASE_MODEL else "hf.co/t-tech/T-lite-it-1.0-Q8_0-GGUF:Q8_0"

    # stream = client.generate(
    #     prompt=test_prompt,
    #     model=test_model,
    #     stream_response=True,
    # )

    # async for text in await stream:
    #     logger.info(text)

    # response  =   await client.generate(
    #     prompt=test_prompt,
    #     model=test_model,
    #     stream_response=False,
    # )
    # logger.info(f"Response: { response}")

    json_response = await client.generate(
        prompt=FIRST_MESSAE_PROMPT.format(message=test_query),
        model=OLLAMA_BASE_MODEL,
        stream_response=False,
        options=OllamaOptions(
            format=JSON_SCHEMA,
            )
        )
    logger.info(f"JSON Response: {json_response}")


if __name__ == "__main__":
    asyncio.run(main())
