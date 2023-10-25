openssl genrsa -aes256 -out ca-key.pem 4096
openssl req -new -x509 -sha256 -days 365 -key ca-key.pem -out ca.pem
openssl genrsa -out cert-key.pem 4096
openssl req -new -sha256 -subj "/CN=yourcn" -key cert-key.pem -out cert.csr
echo "subjectAltName=DNS:Forum,IP:127.0.0.1" >> extfile.cnf
openssl x509 -req -sha256 -days 365 -in cert.csr -CA ca.pem -CAkey ca-key.pem -out cert.pem -extfile extfile.cnf -CAcreateserial