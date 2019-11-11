FROM golang:1.13.4-alpine3.10 as build-env
LABEL maintainer="Elis Lulja <elis.lulja@gmail.com>"

RUN mkdir /polycube-sidecar-injector
WORKDIR /polycube-sidecar-injector
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/polycube-sidecar-injector
FROM scratch
COPY --from=build-env /go/bin/polycube-sidecar-injector /go/bin/polycube-sidecar-injector
ENTRYPOINT ["/go/bin/polycube-sidecar-injector"]