# Estágio de construção
FROM golang:1.23 AS builder

WORKDIR /app

# Copiar os arquivos de dependências
COPY go.mod go.sum ./

# Baixar as dependências
RUN go mod download

# Copiar o código-fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Estágio final
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar o executável do estágio de construção
COPY --from=builder /app/main .

COPY ./files .

# Comando para executar a aplicação
CMD ["./main"]