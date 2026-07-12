FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY bin/server .
EXPOSE 8080
USER nobody:nobody
ENTRYPOINT ["./server"]
