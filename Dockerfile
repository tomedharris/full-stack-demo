FROM golang:1.24.4 AS build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and
# only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN cp ./bin/wait-for-it.sh /usr/local/bin && chmod +x /usr/local/bin/wait-for-it.sh
RUN go build -v -o /usr/local/bin/ ./...

FROM debian:bookworm

EXPOSE 80 443

WORKDIR /
COPY --from=build /usr/local/bin/wait-for-it.sh /usr/local/bin
COPY --from=build /usr/local/bin/web /usr/local/bin

# ENTRYPOINT ["wait-for-it.sh", "-t", "3", "db:3306", "--"]
CMD ["web"]

