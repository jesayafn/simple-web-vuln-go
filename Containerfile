FROM docker.io/library/golang:1.20.11 AS build

WORKDIR /app

COPY . .

RUN go mod download && CGO_ENABLED=0  GOOS=linux go build -o vuln-web .

FROM docker.io/redhat/ubi9-minimal:9.4

WORKDIR /app

COPY --from=build /app/vuln-web . 

ENTRYPOINT [ "/app/vuln-web" ]