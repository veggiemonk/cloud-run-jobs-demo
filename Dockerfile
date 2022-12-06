FROM golang:1.19 as build-env
ARG BUILD_REF
ARG BUILD_DATE

WORKDIR /go/src/app

## No need since we are vendoring
# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

COPY . /go/src/app/

RUN mkdir -p bin && CGO_ENABLED=0 GOOS=linux go build \
        -a \
        -v \
        -trimpath='true' \
        -buildmode='exe' \
        -buildvcs='true' \
        -compiler='gc' \
        -mod='vendor' \
        -ldflags "-X main.gitSha=${BUILD_REF} -X main.buildTime=${BUILD_DATE}" \
        -o bin/batch

#===============================================================================
FROM gcr.io/distroless/static
ARG BUILD_REF
ARG BUILD_DATE

WORKDIR /app

COPY --from=build-env /go/src/app/bin/batch /app/batch
CMD ["/app/batch"]

LABEL org.opencontainers.image.title="batch" \
        org.opencontainers.image.authors="Julien Bisconti" \
        org.opencontainers.image.source="https://github.com/veggiemonk/cloud-run-jobs-demo" \
        org.opencontainers.image.revision="${BUILD_REF}" \
        org.opencontainers.image.created="${BUILD_DATE}" \
        org.opencontainers.image.vendor="Vendor"
