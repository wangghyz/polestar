FROM alpine:latest
COPY server /usr/local/bin/
COPY application.yaml /usr/local/bin/
WORKDIR /usr/local/bin
RUN chmod +x server
EXPOSE 9999
ENTRYPOINT ["server"]