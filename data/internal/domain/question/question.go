package question

type Questions struct {
	ID         string    `db:"id"`
	ChunkID    string    `db:"chunk_id"`
	Question   string    `db:"question"`
	Embeddings []float32 `db:"embeddings"`
}
