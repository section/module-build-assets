FROM alpine:3.6

VOLUME /override

COPY . /src/
RUN chown -R root:root /src/* && chmod -R a+rX /src/*

CMD ["/bin/sh", "-c", "rm -rf /override/* && cp -vr /src/* /override/"]
