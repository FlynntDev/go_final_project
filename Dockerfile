FROM golang:1.22.11



WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o /myapp ./cmd

ENV TODO_PORT=7540
ENV TODO_DBFILE=data/scheduler.db
ENV TODO_PASSWORD=1234

EXPOSE ${TODO_PORT}

CMD ["/myapp"]
