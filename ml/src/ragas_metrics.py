from ragas.metrics import (
    answer_relevancy,
    faithfulness,
    context_precision,
    context_recall,
    context_relevancy,
)
from ragas import evaluate
from datasets import Dataset


def compute_quality_metric(output_text: str) -> float:
    # Пример золотого ответа и используемого контекста — замени своими
    gold_data = {
        "question": ["What is the capital of France?"],
        "answer": ["Paris"],  # Эталонный ответ
        "contexts": [
            ["Paris is the capital of France. It is located in Europe."]
        ],  # Документы
        "generated_answer": [output_text],
    }

    dataset = Dataset.from_dict(gold_data)

    result = evaluate(
        dataset,
        metrics=[
            answer_relevancy,
            faithfulness,
            context_precision,
            context_recall,
            context_relevancy,
        ],
    )

    # Возврат одной метрики или усреднение:
    score = float(result["answer_relevancy"])  # или усреднение нескольких
    return score
