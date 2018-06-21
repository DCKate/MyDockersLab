FROM ubuntu:16.04

RUN useradd -ms /bin/bash ubuntu \
    && usermod -a -G ubuntu www-data \
    &&  apt-get -qq update \
    && apt-get -qq install vim \
    && apt-get -qq install nginx \
    && apt-get -qq install php7.0-fpm \
    && apt-get -qq install php-redis \
    && apt-get -qq install php-xml \
    && apt-get -qq install curl \
    && apt-get -qq install git \
    && mkdir -p /home/ubuntu/www/conf/ \
    && mkdir -p /home/ubuntu/www/logs/ \
    && mkdir -p /home/ubuntu/www/php/sessions/ \
    && mkdir -p /home/ubuntu/www/key/

COPY files/limits.conf /etc/security/limits.conf
# COPY files/sysctl.conf /etc/sysctl.conf
COPY files/fastcgi_params /home/ubuntu/www/conf/fastcgi_params
COPY files/mime.types /home/ubuntu/www/conf/mime.types
COPY files/www.conf /etc/php/7.0/fpm/pool.d/www.conf
COPY files/nginx.conf /home/ubuntu/www/conf/nginx.conf
COPY files/info.php /home/ubuntu/www/php/info.php
COPY files/start.sh /home/ubuntu/www/start.sh

RUN chown -R ubuntu /home/ubuntu/www/ \
    && /sbin/sysctl -p 
    # && service php7.0-fpm start
EXPOSE 55555

CMD ["sh","/home/ubuntu/www/start.sh"]
