FROM golang:1.22 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -o /out/code-executor .

FROM python:3.12-slim
COPY --from=build /out/code-executor /usr/local/bin/code-executor
EXPOSE 50051
ENTRYPOINT ["code-executor"]
