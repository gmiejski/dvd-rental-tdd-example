FROM golang:1.12-stretch

WORKDIR /github.com/gmiejski/dvd-rental-tdd-example
COPY go.mod go.sum ./
RUN go mod download
COPY . /github.com/gmiejski/dvd-rental-tdd-example

# Build binary
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -a -installsuffix cgo -o main .

CMD ["/main"]
