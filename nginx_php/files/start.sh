#!/bin/bash
service php7.0-fpm start
nginx -c /home/ubuntu/www/conf/nginx.conf