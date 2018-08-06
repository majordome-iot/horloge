# build stage
FROM golang:alpine AS build-env
RUN apk update
RUN apk add git
ADD . /go/src/github.com/majordome/horloge
WORKDIR /go/src/github.com/majordome/horloge/
RUN cd ./cmd/horloge && go get
RUN go install
ENTRYPOINT /go/bin/horloge

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/bin/horloge /app/
ENTRYPOINT /app/horloge
