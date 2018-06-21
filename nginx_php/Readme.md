 1. build image
 ```
 $docker build --force-rm . -t nginx_php:1.0
 ```

 2. run 
 ```
 $docker run -it --name php -p 55555:55555 nginx_php:1.0
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
