FROM golang:1.21 as builder
ENV CGO_ENABLED='0'
WORKDIR /home
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ARTEFACT_NAME=contacts.app OUT_DIR=/home make local.build

FROM alpine:3
WORKDIR /home
COPY --from=builder /home/contacts.app .
EXPOSE 8080
ENTRYPOINT [ "/home/contacts.app" ]