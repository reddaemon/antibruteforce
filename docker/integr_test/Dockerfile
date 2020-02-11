FROM golang:1.13 as build-env

ENV NAME "integr_tests"
WORKDIR /opt/${NAME}

COPY go.mod .
COPY go.sum .

COPY . .
