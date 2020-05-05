FROM ubuntu:18.04

RUN apt update && apt install -y --no-install-recommends libcurl4 libapr1 libjansson4 libaprutil1 ca-certificates nginx

COPY proxy/nginx.conf /etc/nginx/

COPY proxy_config/* /opt/proxy_config/

# clean the apt cache to reduce the size of the image
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /var/cache/* /var/log/*

#platform scripts
COPY *.sh /opt/section/
RUN chmod +x /opt/section/*.sh

RUN mkdir /var/log/nginx/
RUN touch var/log/nginx/access.log
RUN touch var/log/nginx/error.log
RUN ln -sf /proc/$$/fd/1   /var/log/nginx/access.log
RUN ln -sf /proc/$$/fd/2   /var/log/nginx/error.log

CMD ["nginx"]
