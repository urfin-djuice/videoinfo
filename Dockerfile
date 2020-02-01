FROM golang:1.13

RUN apt-get update && apt-get install -y --no-install-recommends libmediainfo0v5 libmediainfo-dev

WORKDIR /src

COPY ./ ./

RUN go mod download
RUN go mod verify

RUN go build -o /usr/local/bin/videoinfo ./
