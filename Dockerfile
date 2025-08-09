FROM golang:1.24.4-alpine AS builder

ARG username="appuser"
ARG uid=3125

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && \
    adduser \
            --disabled-password \
            --gecos "" \
            --home "/nonexistent" \
            --shell "/sbin/nologin" \
            --no-create-home \
            --uid "${uid}" \
            "${username}"

COPY --chown=${username}:${username} . .

RUN apk --no-cache add ca-certificates make \
    && make run@build.prod \
    && chmod a+x /app/build/main \
    && chown -R ${username}:${username} /app

FROM gcr.io/distroless/static-debian12 AS final

ARG username

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

WORKDIR /app

USER root

# Copy application files with correct ownership
COPY --chown=${username}:${username} --from=builder /app/build/main .

VOLUME /app/.log

# Expose the necessary port
EXPOSE ${PORT}

# Use an unprivileged user.
USER ${username}:${username}

# Command to run the application
ENTRYPOINT ["/app/main"]