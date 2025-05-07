FROM postgres:16-alpine

# Install build dependencies (grouped as 'build-deps' to make removal easy)
RUN apk add --no-cache --virtual build-deps \
        git build-base postgresql-dev wget cmake linux-headers gawk python3 clang llvm && \
    # Install pgq
    git clone --depth 1 https://github.com/pgq/pgq.git /tmp/pgq && \
    cd /tmp/pgq && \
    make && make install && \
    cd / && rm -rf /tmp/pgq && \
    # Install pgvector
    git clone --depth 1 --branch v0.5.1 https://github.com/pgvector/pgvector.git /tmp/pgvector && \
    cd /tmp/pgvector && \
    make && make install && \
    cd / && rm -rf /tmp/pgvector && \
    # Remove build dependencies
    apk del build-deps

# Add scripts to enable extensions on first start
RUN echo "CREATE EXTENSION IF NOT EXISTS pgq;" > /docker-entrypoint-initdb.d/10-pgq.sql && \
    echo "CREATE EXTENSION IF NOT EXISTS vector;" > /docker-entrypoint-initdb.d/11-pgvector.sql
