FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /tollgate ./cmd/tollgate

FROM scratch
COPY --from=build /tollgate /tollgate
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/tollgate"]
