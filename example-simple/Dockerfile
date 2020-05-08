FROM ubuntu:18.04

# Required dependencies: libcurl4 libapr1 libjansson4 libaprutil1
RUN apt update && apt install -y --no-install-recommends libcurl4 nginx

#nginx config
COPY proxy/nginx.conf /etc/nginx/
COPY proxy_config/* /opt/proxy_config/

#platform scripts
COPY *.sh /opt/section/
RUN chmod +x /opt/section/*.sh

#log to stdout/stderr                            
RUN ln -sf /proc/$$/fd/1   /var/log/nginx/access.log
RUN ln -sf /proc/$$/fd/2   /var/log/nginx/error.log

#run nginx
CMD ["nginx"]
