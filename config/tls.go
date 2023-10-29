package forum

import (
	"log"
	"os/exec"
)

// Auto generate keys and self-signed certificates for ssl/tls connections
func GenerateCert() {

	commands := []string{
		"openssl genrsa -aes256 -out ca-key.pem 4096",
		"openssl req -new -x509 -sha256 -days 365 -key ca-key.pem -out ca.pem -config openssl_forum.cnf",
		"openssl genrsa -out cert-key.pem 4096",
		"openssl req -new -sha256 -subj \"/CN=forum_group\" -key cert-key.pem -out cert.csr",
		"openssl x509 -req -sha256 -days 365 -in cert.csr -CA ca.pem -CAkey ca-key.pem -out cert.pem -extfile openssl_forum.cnf -CAcreateserial",
	}

	for _, cmd := range commands {
		err := exec.Command("bash", "-c", cmd).Run()
		if err != nil {
			log.Fatalf("Command execution failed: %s", err)
		}
	}
}
