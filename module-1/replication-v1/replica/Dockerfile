FROM golang:1.21.3-bookworm
WORKDIR app
COPY . .
RUN go get ./...
RUN go build .
CMD ["./replica"]

