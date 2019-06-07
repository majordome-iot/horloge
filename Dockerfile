# build stage
FROM golang:alpine AS build-env
RUN apk update
RUN apk add git
ADD . /go/src/github.com/majordome-iot/horloge
WORKDIR /go/src/github.com/majordome-iot/horloge/
RUN cd ./cmd/horloge && go get
RUN go install

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/bin/horloge /app/
ENTRYPOINT ["/app/horloge", "-b", "0.0.0.0"]
