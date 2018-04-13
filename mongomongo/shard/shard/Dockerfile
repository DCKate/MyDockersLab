FROM mongo:3.4.14-jessie

RUN mkdir -p /opt/keyfile/
COPY mongo-keyfile /opt/keyfile/
COPY mongod.conf /etc/
RUN chown -R 999 /opt/keyfile/
EXPOSE 27019
