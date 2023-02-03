# Copyright (C) 2023 Timu Eren
# 
# This file is part of getir-case.

FROM docker.io/golang:1.19-bullseye as getir-build
WORKDIR /go/src/getir-case
COPY . .
RUN go mod download 
RUN go build -o /app/getir-case

FROM debian:bullseye-slim
EXPOSE 8085/tcp 8085/tcp
RUN apt-get update && apt-get install -y gettext ca-certificates jq
WORKDIR /app
COPY lib/config/config.tmpl docker-entrypoint.sh ./
RUN chmod +x docker-entrypoint.sh
COPY --from=getir-build /app/getir-case ./
RUN chmod +x getir-case
ENV REDIS_SERVER=redis
ENV MONGO_SERVER=mongo
ENTRYPOINT [ "/app/docker-entrypoint.sh" ]

