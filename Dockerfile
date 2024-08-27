# syntax=docker/dockerfile:1

# ▄▄▄▄    █    ██  ██▓ ██▓    ▓█████▄ ▓█████  ██▀███  
# ▓█████▄  ██  ▓██▒▓██▒▓██▒    ▒██▀ ██▌▓█   ▀ ▓██ ▒ ██▒
# ▒██▒ ▄██▓██  ▒██░▒██▒▒██░    ░██   █▌▒███   ▓██ ░▄█ ▒
# ▒██░█▀  ▓▓█  ░██░░██░▒██░    ░▓█▄   ▌▒▓█  ▄ ▒██▀▀█▄  
# ░▓█  ▀█▓▒▒█████▓ ░██░░██████▒░▒████▓ ░▒████▒░██▓ ▒██▒
# ░▒▓███▀▒░▒▓▒ ▒ ▒ ░▓  ░ ▒░▓  ░ ▒▒▓  ▒ ░░ ▒░ ░░ ▒▓ ░▒▓░
# ▒░▒   ░ ░░▒░ ░ ░  ▒ ░░ ░ ▒  ░ ░ ▒  ▒  ░ ░  ░  ░▒ ░ ▒░
#  ░    ░  ░░░ ░ ░  ▒ ░  ░ ░    ░ ░  ░    ░     ░░   ░ 
#  ░         ░      ░      ░  ░   ░       ░  ░   ░     
#       ░                       ░                      
#
FROM golang:1.23-alpine AS builder

ARG BUILD_VCS_REF
ARG BUILD_VERSION

WORKDIR /usr/src/esprofiler
COPY --chown=root:root . /usr/src/esprofiler
RUN go build -o /app/esprofiler ./cmd/esprofiler


# ██▀███   █    ██  ███▄    █  ███▄    █ ▓█████  ██▀███  
# ▓██ ▒ ██▒ ██  ▓██▒ ██ ▀█   █  ██ ▀█   █ ▓█   ▀ ▓██ ▒ ██▒
# ▓██ ░▄█ ▒▓██  ▒██░▓██  ▀█ ██▒▓██  ▀█ ██▒▒███   ▓██ ░▄█ ▒
# ▒██▀▀█▄  ▓▓█  ░██░▓██▒  ▐▌██▒▓██▒  ▐▌██▒▒▓█  ▄ ▒██▀▀█▄  
# ░██▓ ▒██▒▒▒█████▓ ▒██░   ▓██░▒██░   ▓██░░▒████▒░██▓ ▒██▒
# ░ ▒▓ ░▒▓░░▒▓▒ ▒ ▒ ░ ▒░   ▒ ▒ ░ ▒░   ▒ ▒ ░░ ▒░ ░░ ▒▓ ░▒▓░
#   ░▒ ░ ▒░░░▒░ ░ ░ ░ ░░   ░ ▒░░ ░░   ░ ▒░ ░ ░  ░  ░▒ ░ ▒░
#   ░░   ░  ░░░ ░ ░    ░   ░ ░    ░   ░ ░    ░     ░░   ░ 
#    ░        ░              ░          ░    ░  ░   ░     
#
FROM alpine:3 AS runner

COPY --from=builder /app/esprofiler /bin/esprofiler

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.description="Visualize Elasticsearch profiling results" \
      org.label-schema.name="esprofiler" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.url="https://github.com/rafaeljusto/esprofiler" \
      org.label-schema.vcs-url="https://github.com/rafaeljusto/esprofiler" \
      org.label-schema.vcs-ref=$BUILD_VCS_REF \
      org.label-schema.vendor="Rafael Dantas Justo" \
      org.label-schema.version=$BUILD_VERSION

EXPOSE 80
ENV ESPROFILER_PORT=80
ENTRYPOINT ["/bin/esprofiler"]