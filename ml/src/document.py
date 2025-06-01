from typing import Dict, Any, Optional

class Chunk(Dict):
    id: Optional[str]
    content: str
    similarity: float
    metadata: Dict[str, Any]