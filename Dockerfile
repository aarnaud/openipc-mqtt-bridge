FROM golang:alpine AS builderimage
WORKDIR /go/src/openipc-mqtt-bridge
COPY . .
RUN go build -o openipc-mqtt-bridge main.go


###################################################################

FROM alpine
COPY --from=builderimage /go/src/openipc-mqtt-bridge/openipc-mqtt-bridge /app/
WORKDIR /app
CMD ["./openipc-mqtt-bridge"]
