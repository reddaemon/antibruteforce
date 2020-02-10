FROM golang:1.13 as build-env
ENV NAME "antibruteforce"
WORKDIR /opt/${NAME}
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM build-env AS build
ENV NAME "antibruteforce"
WORKDIR /opt/${NAME}
COPY . .
RUN CGO_ENABLED=0 go build -o bin/${NAME}

FROM alpine
ENV NAME "antibruteforce"
WORKDIR /opt/${NAME}
COPY --from=build /opt/${NAME}/bin/${NAME} ./${NAME}
RUN apk add --no-cache tzdata
EXPOSE 8080
CMD ./${NAME} api