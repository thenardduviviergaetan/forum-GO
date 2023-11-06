# Commands to manually generate certificates and keys for SSL/TLS

```bash
openssl genrsa -aes256 -out ca-key.pem 4096
openssl req -new -x509 -sha256 -days 365 -key ca-key.pem -out ca.pem -config openssl_forum.cnf
openssl genrsa -out cert-key.pem 4096
openssl req -new -sha256 -key cert-key.pem -out cert.csr -config openssl_forum.cnf
openssl x509 -req -sha256 -days 365 -in cert.csr -CA ca.pem -CAkey ca-key.pem -out cert.pem -extfile forum.cnf -CAcreateserial
```
