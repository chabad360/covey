package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/crypto/ssh"
)

var (
	clients []ssh.Client
)

func main() {
	http.HandleFunc("/add", addHost)
	http.HandleFunc("/run", hRun)

	http.ListenAndServe(":8081", nil)
}

func addHost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.URL.Query().Get("user"))
	fmt.Fprintf(w, r.URL.Query().Get("pass"))
	fmt.Fprintf(w, r.URL.Query().Get("server"))
	fmt.Fprintf(w, r.URL.Query().Get("port"))

	client, err := setupHost(r.URL.Query().Get("user"), r.URL.Query().Get("pass"), r.URL.Query().Get("server"), r.URL.Query().Get("port"))
	if err != nil {
		log.Fatal(err)
	}

	c, err := run(client, "/usr/bin/whoami")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, *c)
	clients = append(clients, *client)

}

func hRun(w http.ResponseWriter, r *http.Request) {
	for c := range clients {
		client := new(ssh.Client)
		*client = clients[c]
		go run(client, r.URL.Query().Get("cmd"))
	}
}

func run(client *ssh.Client, command string) (*string, error) {
	var b bytes.Buffer
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		return nil, err
	}
	r := new(string)
	*r = b.String()
	session.Close()
	return r, nil
}

func setupHost(username string, password string, server string, port string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	var b bytes.Buffer

	client, err := ssh.Dial("tcp", server+":"+port, config)
	if err != nil {
		return nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		return nil, err
	}
	fmt.Println(b.String())

	key, err := generatePrivateKey(4096)
	if err != nil {
		return nil, err
	}
	pubKey, err := generatePublicKey(&key.PublicKey)
	if err != nil {
		return nil, err
	}
	err = storePublicKeyOnHost(client, pubKey)
	if err != nil {
		return nil, err
	}
	session.Close()

	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		return nil, err
	}

	config = &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err = ssh.Dial("tcp", server+":"+port, config)
	if err != nil {
		return nil, err
	}
	session, err = client.NewSession()
	if err != nil {
		return nil, err
	}
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		return nil, err
	}
	session.Close()

	return client, nil
}

func storePublicKeyOnHost(client *ssh.Client, keyBytes []byte) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}

	if err := session.Run("echo '" + string(keyBytes) + "' | tee -a .ssh/authorized_keys"); err != nil {
		return err
	}
	log.Println("Stored key on host")
	session.Close()
	return nil
}

// generatePrivateKey creates a RSA Private Key of specified byte size
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	log.Println("Private Key generated")
	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	log.Println("Public key generated")
	return pubKeyBytes, nil
}

// writePemToFile writes keys to a file
func writeKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}

	log.Printf("Key saved to: %s", saveFileTo)
	return nil
}
