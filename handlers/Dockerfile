FROM golang:1.22.0
WORKDIR /app
COPY . .
RUN go get ./...
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /gorilla-products-maintenance
ENV PORT 9090
CMD ["/gorilla-products-maintenance"]