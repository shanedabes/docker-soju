FROM golang:alpine AS build-env
RUN apk add gcc libc-dev git
Add . /src
WORKDIR /src
RUN mkdir /app
RUN go build -o /app/dbgen dbgen.go
RUN cp /src/entrypoint.sh /app
RUN git clone https://github.com/emersion/soju.git \
 && cd soju \
 && git checkout v0.1.2 \
 && go build -o /app/soju cmd/soju/main.go

FROM alpine
RUN apk add ca-certificates openssl
WORKDIR /app
COPY --from=build-env /app /app
ENTRYPOINT ["./entrypoint.sh"]
