FROM golang:1.22.11

ENV TODO_PORT=7540
ENV TODO_DBFILE=data/scheduler.db
ENV TODO_PASSWORD=1234
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod tidy

# Вывод информации о зависимостях и модуле
RUN go list -m all
RUN go list -m

EXPOSE ${TODO_PORT}

# Указываем путь к main.go в папке cmd и проверяем создание файла myapp
RUN go build -o /myapp ./cmd && ls -la /myapp

CMD ["/myapp"]
