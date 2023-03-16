FROM golang:alpine as builder

COPY . /usr/src/openstack-cni
WORKDIR /usr/src/openstack-cni

RUN apk add --no-cache make musl-dev && make clean && make build

FROM alpine:3
COPY --from=builder /usr/src/openstack-cni/bin/openstack-cni /usr/bin/
COPY --from=builder /usr/src/openstack-cni/bin/openstack-cni-daemon /usr/bin/
WORKDIR /

LABEL io.k8s.display-name="OPENSTACK CNI"

COPY ./entrypoint.sh /

ENTRYPOINT ["./entrypoint.sh"]