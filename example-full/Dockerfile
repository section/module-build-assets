#BEGIN handler
FROM golang:1.12 as handler

# produce a statically compiled binary
ENV CGO_ENABLED=0

WORKDIR /go/src/app

# explicitly install dependencies to improve Docker re-build times
RUN go get -v \
  github.com/kisielk/errcheck

WORKDIR /app

# copy the required files, check and create the binary
COPY src/handler/app/go.mod src/handler/app/go.sum ./

RUN go mod download

COPY src/handler/app .

RUN gofmt -e -s -d . 2>&1 | tee /gofmt.out && test ! -s /gofmt.out
RUN go vet .
RUN errcheck ./...

RUN go install ./...

COPY src/environment_variables.lua.gotemplate ./testdata/
COPY src/module_schema.json ./testdata/
COPY src/tests/* ./testdata/
RUN go test -v ./...

#END handler

FROM ubuntu:16.04

# https://openresty.org/package/pubkey.gpg
COPY src/proxy/pubkey.gpg /

ARG openresty_version=1.13.6.1

ARG OPENRESTY_BASE_DIR=/etc/openresty
ARG SECTION_DIR=/opt/section/

USER root

# install the required packages
RUN printf 'deb http://openresty.org/package/ubuntu xenial main\n' >/etc/apt/sources.list.d/openresty.list && \
  apt-key add /pubkey.gpg && \
  apt-get update && \
  apt-get install --assume-yes --no-install-recommends \
    "openresty=${openresty_version}-*" \
    zlib1g \
    libssl1.0.0 libssl-dev

# clean the apt cache to reduce the size of the image
RUN apt-get clean && \
  rm -rf /var/lib/apt/lists/* /var/cache/* /var/log/*

# remove existing logs and link nginx logs
RUN rm -d /usr/local/openresty/nginx/logs && \
  ln -s /var/log/nginx /usr/local/openresty/nginx/logs

# create the required directories and assign them minimum required permissions
# /var/log/nginx : nginx logs
# /opt/section : module assets
# /var/openresty : directory for error_log.conf, resolver.conf, nginx pid file, other nginx temp files
RUN mkdir -p /var/log/nginx && \
  mkdir -p $SECTION_DIR && \
  mkdir -p /var/openresty && \
  chmod -R 0777 /var/log/nginx && \
  chmod -R 0755 $SECTION_DIR && \
  chmod -R 0777 /var/openresty

# give the nginx process capabilites to bind to port 80
# optionally verify it (-v flag)
# https://linux.die.net/man/8/setcap
RUN  setcap    cap_net_bind_service=+ep /usr/local/openresty/nginx/sbin/nginx \
  && setcap -v cap_net_bind_service=+ep /usr/local/openresty/nginx/sbin/nginx


# copy the module handler and assets
COPY --from=handler /go/bin/module-handler $SECTION_DIR/

COPY src/lib/ /usr/local/openresty/lualib/
COPY src/proxy/nginx.conf $OPENRESTY_BASE_DIR/nginx.conf
COPY src/proxy/content.lua $OPENRESTY_BASE_DIR/content.lua

COPY src/environment_variables.lua.gotemplate $SECTION_DIR/
COPY src/module_schema.json $SECTION_DIR/

# set the minimum required permissions
RUN chmod 0644 $OPENRESTY_BASE_DIR/nginx.conf && \
  chmod 0644 $OPENRESTY_BASE_DIR/content.lua && \
  chmod -R 0777 /usr/local/openresty/lualib/section/

# symlink module-handler to validate.sh
# used to run validation tests during git push
RUN ln -s $SECTION_DIR/module-handler $SECTION_DIR/validate.sh

# switch to the user created by openresty install
USER www-data

# use the default entrypoint provided by the go handler
ENTRYPOINT ["/opt/section/module-handler"]
