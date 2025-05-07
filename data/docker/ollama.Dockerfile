FROM ollama/ollama:0.6.7
ENV MODELS="bge-m3 llama3.2"
ENV OLLAMA_KEEP_ALIVE=24h
ENTRYPOINT [ "/bin/bash", "-c", "(sleep 5 ; for m in $MODELS ; do ollama pull $m ; done) & exec /bin/ollama $0" ]
CMD [ "serve" ]