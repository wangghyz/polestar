FROM alpine:latest
COPY polestar /usr/local/bin/
COPY application.yaml /usr/local/bin/
WORKDIR /usr/local/bin
CMD polestar