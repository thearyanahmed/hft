FROM golang:1.21

ENV APP_DIR /app

# Install CompileDaemon for auto installing binary on changes
RUN go install github.com/githubnemo/CompileDaemon@latest

COPY cmd/pkg/ pkg/ .env Makefile $APP_DIR
WORKDIR $APP_DIR

RUN make deps
