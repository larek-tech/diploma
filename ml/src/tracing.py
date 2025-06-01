from opentelemetry.sdk.resources import SERVICE_NAME, Resource

from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor

from config import TRACING_ENDPOINT
from utils.logger import logger

resource = Resource.create(attributes={
    SERVICE_NAME: "ml"
})

tracerProvider = TracerProvider(resource=resource)
processor = BatchSpanProcessor(OTLPSpanExporter(endpoint=f"{TRACING_ENDPOINT}"))
tracerProvider.add_span_processor(processor)
trace.set_tracer_provider(tracerProvider)
tracer = tracerProvider.get_tracer("ml_service")

logger.info("OpenTelemetry tracing initialized")