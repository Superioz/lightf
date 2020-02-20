FROM golang:1.13.7-alpine3.10 AS builder

# after Docker 1.10+ labels could be written in
# multiple lines without multiple layers, but
# to support prior versions ...
LABEL maintainer="Tobias Büser <tobias.bueser@yahoo.de>" \
      version="0.2.16" \
      description="Lightweight, secure and transient fileserver written in Go."

WORKDIR /build

# enables caching of modules as a Docker layer
COPY go.mod .
COPY go.sum .
RUN go mod download

# copy the rest of the directory to 
COPY . .

# build the program
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -o ./out/app ./cmd/lightf-serv

################
# run program
################

FROM alpine:3.10

COPY --from=builder /build/out/app .

CMD ["./app"]
