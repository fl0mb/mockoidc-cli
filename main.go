package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"flag"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fl0mb/mockoidc"
	"golang.org/x/crypto/acme/autocert"
)

func genTLSmemory(cn string) *tls.Certificate {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		panic(err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: cn,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(5, 0, 0),
		BasicConstraintsValid: false,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	publicKey := &privateKey.PublicKey

	derEncodedCert, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey, privateKey)
	if err != nil {
		panic(err)
	}

	cert := &tls.Certificate{
		Certificate: [][]byte{derEncodedCert},
		PrivateKey:  privateKey,
	}

	return cert
}

func autocert_get(cn string) *tls.Config {

	autocertManager := &autocert.Manager{
		Cache:      autocert.DirCache("secret-dir"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(cn),
	}

	go func() {
		srv := &http.Server{
			Addr:         ":80",
			Handler:      autocertManager.HTTPHandler(nil),
			IdleTimeout:  time.Minute,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		err := srv.ListenAndServe()
		if err != nil {
			panic(err)
		}

	}()
	tlsConfig := autocertManager.TLSConfig()
	return tlsConfig
}

func main() {

	//brauch ich das noch ?
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	listener := flag.String("l", "127.0.0.1:8080", "listener to use <ip>:<port>")
	clientid := flag.String("client-id", "api://AzureADTokenExchange", "client-id to request token with")
	clientsecret := flag.String("client-secret", "secret", "client secret")
	https := flag.Bool("https", false, "Whether to serve using https instead of http. This will automatically create a self-signed certificate in memory")
	// bis hier geht alles jetzt noch tls lösen
	cn := flag.String("cn", "secret", "Define the common name for a automatically requested Let's Encrypt certificate. This option implies \"https\" and it will start an autocert listener on port 80 to complete a HTTP-01 challenge.")
	flag.Parse()

	ln, err := net.Listen("tcp", *listener)
	if err != nil {
		panic(err)
	}

	m, _ := mockoidc.NewServer(nil)
	//aud claim is set based on ClientID
	m.ClientID = *clientid
	m.ClientSecret = *clientsecret

	// request logging
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			fmt.Printf("\n%v %v \n %v", req.Method, req.RequestURI, req.Header)
			next.ServeHTTP(rw, req)
		})
	}
	m.AddMiddleware(middleware)

	if *https {
		cert := genTLSmemory(m.GetIssuer())
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{*cert},
		}
		m.Start(ln, tlsConfig)
		defer m.Shutdown()
	} else if *cn != "" {
		m.Issuer = *cn
		m.Start(ln, autocert_get(*cn))
	} else {
		m.Start(ln, nil)
	}

	fmt.Printf("Issuer: %s\nclient_id: %s\nclient_secret: %s\n", m.GetIssuer(), m.ClientID, m.ClientSecret)

	<-sigs

}
