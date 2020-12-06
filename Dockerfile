# build stage
FROM golang:alpine AS build-env
RUN apk update
RUN apk add git
ADD . /go/src/github.com/shinuza/horloge
WORKDIR /go/src/github.com/shinuza/horloge/
RUN cd ./cmd/horloge && go get
RUN go install

# final stage
# FROM scratch
# WORKDIR /app
# COPY --from=build-env /go/bin/horloge /app/
ENTRYPOINT ["/go/bin/horloge"]
