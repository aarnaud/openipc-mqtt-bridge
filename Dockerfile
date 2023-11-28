FROM golang:alpine AS builderimage
WORKDIR /go/src/openipc-speaker-bridge
COPY . .
RUN go build -o openipc-speaker-bridge main.go


###################################################################

FROM alpine
COPY --from=builderimage /go/src/openipc-speaker-bridge/openipc-speaker-bridge /app/
WORKDIR /app
CMD ["./openipc-speaker-bridge"]
