FROM golang:1.14-alpine AS build
LABEL maintainer="cccfs <408473944@qq.com>"

WORKDIR /srv
COPY . .
RUN apk add git
RUN go get ./...  \
    && go build -o kube-ldap-webhook

FROM alpine:3.13.4
WORKDIR /srv
COPY --from=build /srv/kube-ldap-webhook /srv/kube-ldap-webhook
EXPOSE 8080
ENTRYPOINT ./kube-ldap-webhook
