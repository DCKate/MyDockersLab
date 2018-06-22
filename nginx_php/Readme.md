 1. build image
 ```
 $docker build --force-rm . -t nginx_php:1.0
 ```

 2. create network
 ```
 $ docker network create small
 ```
 2.1 connect network
 Connect a running container to a network
 ```
 $ docker network connect multi-host-network container1
 ```
 Connect a container to a network when it starts
 You can also use the docker run --network=<network-name> option to start a container and immediately connect it to a network.
 ```
 $ docker run -itd --network=multi-host-network busybox
 ```
 3. run 
 ```
 $docker run -itd --network=small --name php -p 55555:55555 nginx_php:1.0
 ```

 4. exec
 ```
 $docker exec -it -u ubuntu php bash
 ```
 
 5. run redis
 ```
 $docker run -v {{ docker_redis_folder }}:/data --network=small --name goredis -p 6379:6379 -d redis redis-server --appendonly yes
 ```


Note:
* How to sign ssl cert
Create an SSL CertificatePermalink
Let’s Encrypt automatically performs Domain Validation (DV) using a series of challenges. The Certificate Authority (CA) uses challenges to verify the authenticity of your computer’s domain. Once your Linode has been validated, the CA will issue SSL certificates to you.

Run Let’s Encrypt with the --standalone parameter. For each additional domain name requiring a certificate, add -d example.com to the end of the command.
Let’s Encrypt does not deploy wildcard certificates. Each subdomain requires its own certificate.
 ```
 $git clone https://github.com/letsencrypt/letsencrypt opt/letsencrypt
 $cd opt/letsencrypt
 $./letsencrypt-auto certonly --standalone -d example.com -d www.example.com
 ```
