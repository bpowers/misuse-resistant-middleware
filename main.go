package main

import (
	"crypto/tls"
	"flag"
	"goedge/middleware/a"
	"goedge/middleware/b"
	"log"
	"net/http"
	"time"

	"goji.io"
	"goji.io/pat"
)

func main() {
	addr := flag.String("addr", "localhost:443", "address to listen on")
	certFile := flag.String("cert", "localhost.pem", "TLS certificate file")
	keyFile := flag.String("key", "localhost-key.pem", "TLS private key file")
	flag.Parse()

	// misuse resistant middleware pattern.  You can't add NewEchoNameHandler to the mux
	// without having (a) constructed and (b) registered NewNAmeExtractorMiddleware with
	// the mux.
	nameMux := goji.SubMux()
	nameMiddleware := a.NewNameExtractorMiddleware("name")
	nameKey := nameMiddleware.Register(nameMux)
	nameMux.Handle(pat.New(""), b.NewEchoNameHandler(nameKey))

	mux := goji.NewMux()
	mux.Handle(pat.Get("/hello/:name"), nameMux)

	server := &http.Server{
		Addr: *addr,
		Handler: mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  15 * time.Second,
		TLSConfig: &tls.Config{
			// Causes servers to use Go's default ciphersuite preferences,
			// which are tuned to avoid attacks. Does nothing on clients.
			PreferServerCipherSuites: true,
			// Only use curves which have assembly implementations
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}

	log.Printf("serving on addr %s\n", *addr)
	if err := server.ListenAndServeTLS(*certFile, *keyFile); err != nil {
		log.Fatalf("ListenAndServeTLS: %s\n", err)
	}
}
