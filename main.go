package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func serverHTTP(addr string, handler http.Handler, err chan<- error) {
	err <- http.ListenAndServe(addr, handler)
}

func serverHTTPS(addr, caCert, caKey string, handler http.Handler, err chan<- error) {
	err <- http.ListenAndServeTLS(addr, caCert, caKey, handler)
}

func main() {
	// Create the flags to get parameter around ca certs
	var port, sslPort int
	var caCertPath, caKeyPath, host string
	flag.IntVar(&port, "port", 8080, "Port to serve non tls proxy")
	flag.IntVar(&sslPort, "ssl-port", 8443, "SSL port to serve tls proxy")

	flag.StringVar(&caCertPath, "ca-cert-path", "", "CA cert path for TLS")
	flag.StringVar(&caKeyPath, "ca-key-path", "", "CA key path for TLS")
	flag.StringVar(&host, "host", "", "host (domain/IP) for alternate dns in certificate")

	flag.Parse()

	if (caKeyPath != "" && caCertPath == "") || (caKeyPath == "" && caCertPath != "") {
		log.Fatalf("ca-cert-path and ca-key-path should be used together")
	}

	// Create error channel
	errCh := make(chan error, 1)

	if caKeyPath == "" && caCertPath == "" {
		var err error
		hosts := []string{"localhost", "127.0.0.1", "192.168.122.1"}
		if host != "" {
			hosts = append(hosts, host)
		}
		caCertPath, caKeyPath, err = generateCA(hosts)
		if err != nil {
			log.Fatal("Failed to generate proxy CA")
		}
		log.Printf("CA cert file %s\n", caCertPath)
		log.Printf("CA key file %s\n", caKeyPath)
	}
	{
		// Get the system certificate pools
		caCertPool, err := x509.SystemCertPool()
		if err != nil {
			log.Printf("Not able to get system certificate pool %v", err)
			caCertPool = x509.NewCertPool()
		}
		caCert, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			log.Fatalf("Error reading %s file %v", caCertPath, err)
		}
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			log.Fatal("Failed to append proxy CA to system CAs")
		}

		proxy := goproxy.NewProxyHttpServer()
		proxy.Tr = &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS13,
				RootCAs:    caCertPool,
			},
		}
		proxy.Verbose = true
		p := fmt.Sprintf(":%d", sslPort)
		go serverHTTPS(p, caCertPath, caKeyPath, proxy, errCh)
	}

	p := fmt.Sprintf(":%d", port)
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	go serverHTTP(p, proxy, errCh)
	log.Fatal(<-errCh)
}
