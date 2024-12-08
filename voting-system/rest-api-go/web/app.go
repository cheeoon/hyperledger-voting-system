package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/rs/cors"
)

// OrgSetup contains organization's config to interact with the network.
type OrgSetup struct {
	OrgName      string
	MSPID        string
	CryptoPath   string
	CertPath     string
	KeyPath      string
	TLSCertPath  string
	PeerEndpoint string
	GatewayPeer  string
	Gateway      client.Gateway
}

// Serve starts http web server.
func Serve(setups OrgSetup) {
	// Setup CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:4200"}, // Your frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Allow credentials if needed
	})

	// Register routes
	http.HandleFunc("/query", setups.Query)
	http.HandleFunc("/invoke", setups.Invoke)

	// Apply the CORS middleware to all routes
	handler := c.Handler(http.DefaultServeMux)

	// Start the server with CORS enabled
	fmt.Println("Listening on http://localhost:3000...")
	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
