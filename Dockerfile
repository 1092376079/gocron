FROM alpine:3.12

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -g app app

ENV TZ Asia/Shanghai

WORKDIR /app

COPY ./bin/gocron .

RUN chown -R app:app ./

EXPOSE 5920

USER app

ENTRYPOINT ["/app/gocron", "web"]
