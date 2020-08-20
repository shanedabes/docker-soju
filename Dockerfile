FROM golang:alpine AS build-env
Add . /src
WORKDIR /src
RUN apk add gcc libc-dev git
RUN mkdir /app
RUN go build -o /app/dbgen dbgen.go
RUN cp /src/entrypoint.sh /app
RUN git clone https://github.com/emersion/soju.git \
 && cd soju \
 && git checkout 81c7e80e0fa47619d80dc941104dfe7da73ca58c \
 && go build -o /app/soju cmd/soju/main.go

FROM alpine
RUN apk add tini ca-certificates openssl
WORKDIR /app
COPY --from=build-env /app /app
ENTRYPOINT ["tini", "--"]
CMD ./entrypoint.sh
