FROM ubuntu:16.04

RUN useradd -ms /bin/bash ubuntu \
    && usermod -a -G ubuntu www-data \
    && apt-get -qq update \
    && apt-get -qq install vim \
    && apt-get -qq install nginx supervisor \
    && apt-get -qq install curl \
    && apt-get -qq install git \
    && mkdir -p /home/ubuntu/www/conf/sites-enabled/ \
    && mkdir -p /home/ubuntu/www/protected/key/ \
    && mkdir -p /home/ubuntu/www/logs/ \
    && mkdir -p /home/ubuntu/www/gbin/

COPY files/mime.types /home/ubuntu/www/conf/mime.types
COPY files/nginx.conf /home/ubuntu/www/conf/nginx.conf
COPY files/goweb.conf /home/ubuntu/www/conf/sites-enabled/goweb.conf
COPY files/start.conf /etc/supervisor/conf.d/start.conf
COPY files/enc.key /home/ubuntu/www/protected/key/
COPY files/goweb /home/ubuntu/www/gbin/
# ADD source/goweb/protected/demo.tar.gz /home/ubuntu/www/protected/
ADD source/goweb/protected/localdemo.tar.gz /home/ubuntu/www/protected/
ADD source/goweb/protected/encdemo.tar.gz /home/ubuntu/www/protected/

WORKDIR /home/ubuntu/www/

RUN chown -R ubuntu /home/ubuntu/www/ \
    && chown -R www-data:www-data /home/ubuntu/www \
    && chown -R www-data:www-data /home/ubuntu/www/protected

EXPOSE 55555

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/start.conf"]
