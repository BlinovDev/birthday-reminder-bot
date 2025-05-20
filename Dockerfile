# -- build stage --
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /birthday-reminder-bot main.go

# -- runtime stage --
FROM alpine
WORKDIR /app
COPY --from=build /birthday-reminder-bot .
# Copy or mount config.yml at runtime
CMD ["./birthday-reminder-bot"]
