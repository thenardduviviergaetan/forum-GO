package forum

import (
	"log"
	"os"
	"os/exec"
)

// Auto generate keys and self-signed certificates for ssl/tls connections
func GenerateCert() {
	var missed bool
	files := []string{"ca-key.pem", "ca.pem", "cert-key.pem", "cert.pem", "cert.csr"}

	for _, file := range files {
		if !FileExist(file) {
			missed = true
		}
	}

	if missed {
		commands := []string{
			"openssl genrsa -aes256 -out ca-key.pem 4096",
			"openssl req -new -x509 -sha256 -days 365 -key ca-key.pem -out ca.pem -config openssl_forum.cnf",
			"openssl genrsa -out cert-key.pem 4096",
			"openssl req -new -sha256 -key cert-key.pem -out cert.csr -config openssl_forum.cnf",
			"openssl x509 -req -sha256 -days 365 -in cert.csr -CA ca.pem -CAkey ca-key.pem -out cert.pem -extfile openssl_forum.cnf",
		}

		for _, cmd := range commands {
			err := exec.Command("bash", "-c", cmd).Run()
			if err != nil {
				log.Fatalf("Command execution failed: %s", err)
			}
		}
	}
}

func FileExist(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
