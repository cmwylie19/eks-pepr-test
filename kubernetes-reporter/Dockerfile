FROM golang:1.23.3-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o k8s-reporter .


FROM busybox

COPY --from=builder /app/k8s-reporter /k8s-reporter
ENTRYPOINT ["/k8s-reporter"]
