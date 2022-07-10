FROM golang:1.17-alpine as build

WORKDIR /app

COPY . ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o /spotify-search


FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=build /app/templates /templates
COPY --from=build /spotify-search /spotify-search


EXPOSE 8080

USER nonroot:nonroot

CMD ["/spotify-search", "web"]