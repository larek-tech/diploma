queue:
- [x] move each type of sources into its own queue
- [x] integrate with embedder

crawler:
- [x] sending result messages with progress and page_id
- [x] create embedder handler foe processing crawler messages
- [x] set up scheduled parsing tasks and define process for saving resources and do not calculate embeddings if page is not changed
- [x] add a way to display progress of each parsing job

s3:
- [ ] move to same implementation as crawler with multiple handlers and processes
- [ ] integrate with embedder

raw files:
- [ ] create implementation similar to crawler for processing raw files
- [ ] integrate with embedder

file storage:
- [x] migrate storing raw html and text files to s3 like storages
