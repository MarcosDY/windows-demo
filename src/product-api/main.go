package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/spiffe/go-spiffe/v2/logger"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

var (
	namedPipeNameFlag = flag.String("namedPipeName", "\\spire-agent\\public\\api", "Agent named pipe name")
)

type ProductsResponse struct {
	Products []*Product `json:"products"`
}

type Product struct {
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

func productsList(w http.ResponseWriter, r *http.Request) {
	log.Println("List products called...")
	if r.Method != http.MethodGet {
		log.Printf("Invalid http method: %q", r.Method)
		http.Error(w, "unexpected http method", http.StatusInternalServerError)
		return
	}

	// Attes app and write response on disk
	listResp := &ProductsResponse{
		Products: []*Product{
			{
				Name:  "Bags",
				Stock: 11,
			},
			{
				Name:  "Sweet potatos",
				Stock: 2,
			},
		},
	}

	if err := json.NewEncoder(w).Encode(listResp); err != nil {
		log.Printf("Error processing payload: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up a `/products` resource handler
	http.HandleFunc("/products", productsList)
	// Create logger for workload API
	logStdOut := logger.Writer(os.Stdout)

	x509Source, err := workloadapi.NewX509Source(ctx,
		workloadapi.WithClientOptions(workloadapi.WithNamedPipeName(*namedPipeNameFlag),
			workloadapi.WithLogger(logStdOut)))
	if err != nil {
		log.Fatalf("Unable to create X509Source: %v", err)
	}
	defer x509Source.Close()

	bundleSource, err := workloadapi.NewBundleSource(ctx, workloadapi.WithClientOptions(workloadapi.WithNamedPipeName(*namedPipeNameFlag),
		workloadapi.WithLogger(logStdOut)))
	if err != nil {
		log.Fatalf("Unable to create BundleSource %v", err)
	}
	defer bundleSource.Close()

	svid, err := x509Source.GetX509SVID()
	if err == nil {
		log.Printf("Identity: %v", svid.ID)
	}

	// Allowed SPIFFE ID
	clientID := spiffeid.RequireFromString("spiffe://example.org/webapp")

	// Create server listening on 9443
	tlsConfig := tlsconfig.MTLSServerConfig(x509Source, bundleSource, tlsconfig.AuthorizeID(clientID))
	server := &http.Server{
		Addr:      ":9443",
		TLSConfig: tlsConfig,
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Error on serve: %v", err)
	}
}
