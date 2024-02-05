FROM golang:1.21.6-alpine3.18
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /app .

FROM vimagick/youtube-dl
COPY --from=0 /app /app
ENTRYPOINT [ "/app" ]