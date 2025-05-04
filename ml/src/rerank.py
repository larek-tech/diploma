import torch
from transformers import AutoModelForSequenceClassification, AutoTokenizer


class Reranker:
    """Класс для повторной оценки документов на основе заданного запроса
    с использованием предобученной модели трансформера.

    Parameters
    ----------
    reranker_model_name : str
        Название предобученной модели, используемой для повторной оценки.
    max_length : int
        Максимальная длина входных последовательностей.
    device : str
        Устройство, на котором будет выполняться модель ('cpu' или 'cuda').

    Methods
    -------
    rerank_documents:
        Повторно оценивает предоставленные документы на основе заданного
        запроса и возвращает top_k ранжированных документов.
    """

    def __init__(
        self, reranker_model_name: str, max_length: int, device: str
    ) -> None:
        self.reranker_tokenizer = AutoTokenizer.from_pretrained(
            reranker_model_name
        )
        self.reranker_model = (
            AutoModelForSequenceClassification.from_pretrained(
                reranker_model_name
            )
        )
        self.max_length = max_length
        self.device = device
        self.reranker_model.to(device)

    def rerank_documents(
        self, query: str, documents: list[str], top_k: int = 5
    ) -> list[tuple[str, float]]:
        """Повторно оценивает предоставленные документы на основе
          заданного запроса.

        Parameters
        ----------
        query : str
            Запрос, по которому будует оцениваться документы.
        documents : list[str]
            Список строк документов, которые необходимо повторно оценить.
        top_k : int, optional
            Количество документов с наивысшими оценками, которые нужно
            вернуть (по умолчанию 5).

        Returns
        -------
        list[tuple[str, float]]
            Список кортежей, где каждый кортеж содержит документ и
            соответствующую ему оценку, отсортированный в порядке убывания
            оценок, ограниченный top_k документами.
        """
        pairs = []
        for doc in documents:
            pairs.append((query, doc))

        features = self.reranker_tokenizer(
            pairs,
            padding=True,
            truncation=True,
            return_tensors="pt",
            max_length=self.max_length,
        ).to(self.device)

        with torch.no_grad():
            scores = (
                self.reranker_model(**features).logits.squeeze().cpu().numpy()
            )

        doc_score_pairs = list(zip(documents, scores, strict=False))
        ranked_docs = sorted(doc_score_pairs, key=lambda x: x[1], reverse=True)

        return ranked_docs[:top_k]
