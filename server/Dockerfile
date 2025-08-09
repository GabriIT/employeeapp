FROM golang:1.22

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Accept git SHA & build time from Dokku (or use defaults)
ARG GIT_SHA=dev
ARG BUILD_TIME=unknown

# Build with ldflags injecting version data
RUN CGO_ENABLED=0 go build -ldflags "-X main.Version=${GIT_SHA} -X main.BuildTime=${BUILD_TIME}" -o main .

EXPOSE 8080
CMD ["./main"]
