FROM golang:1.24.1 as builder
ENV CGO_ENABLED='0'
WORKDIR /home
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make executable ARTEFACT_PATH=/home/app

FROM alpine:3
COPY --from=builder /home/app /home/app
EXPOSE 8080
ENTRYPOINT [ "/home/app" ]