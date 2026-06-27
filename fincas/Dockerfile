FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY bin/fincas .
EXPOSE 8082
CMD ["./fincas"]
