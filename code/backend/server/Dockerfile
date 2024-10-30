# 使用官方的 Go 1.21 映像作為構建階段
FROM golang:1.21 as builder

# 設定工作目錄
WORKDIR /app

# 將 go.mod 和 go.sum 複製到工作目錄
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製所有源代碼到工作目錄
COPY . .

# 編譯應用程式
RUN go build -o main .

# 使用輕量級的運行時映像
FROM debian:bullseye-slim

# 設定工作目錄
WORKDIR /app

# 從構建階段複製二進制文件到運行時映像
COPY --from=builder /app/main .

# 設定環境變數，根據需要修改
ENV PORT=8080

# 暴露應用程式埠
EXPOSE ${PORT}

# 執行應用程式
CMD ["./main"]
