FROM golang:1.20 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o admission-controller .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/admission-controller .
COPY tls.crt tls.key /root/
EXPOSE 8443
CMD ["./admission-controller"]
