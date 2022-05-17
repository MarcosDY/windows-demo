package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spiffe/go-spiffe/v2/logger"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

const (
	customerAPIURL = "https://customer-api:8443/"
	productAPIURL  = "https://host.docker.internal:9443/"
	port           = 8080
)

var (
	x509Source        = &workloadapi.X509Source{}
	bundleSource      = &workloadapi.BundleSource{}
	namedPipeNameFlag = flag.String("namedPipeName", "\\spire-agent\\public\\api", "Agent named pipe name")
)

type CustomersResponse struct {
	Customers []*Customer `json:"customers"`
}

type Customer struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type ProductsResponse struct {
	Products []*Product `json:"products"`
}

type Product struct {
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

func getCustomers() ([]*Customer, error) {
	customersID := spiffeid.RequireFromString("spiffe://example.org/customers-api")
	tlsConfig := tlsconfig.MTLSClientConfig(x509Source, bundleSource, tlsconfig.AuthorizeID(customersID))
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	resp, err := client.Get(customerAPIURL + "/customers")
	if err != nil {
		return nil, fmt.Errorf("error connecting to %q: %v", customerAPIURL, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	listResp := new(CustomersResponse)
	if err := json.NewDecoder(resp.Body).Decode(listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp.Customers, nil
}

func getProducts() ([]*Product, error) {
	serverID := spiffeid.RequireFromString("spiffe://example.org/products-api")
	tlsConfig := tlsconfig.MTLSClientConfig(x509Source, bundleSource, tlsconfig.AuthorizeID(serverID))
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	resp, err := client.Get(productAPIURL + "products")
	if err != nil {
		return nil, fmt.Errorf("error connecting to %q: %v", productAPIURL, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	listResp := new(ProductsResponse)
	if err := json.NewDecoder(resp.Body).Decode(listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp.Products, nil
}

func indexHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	customers, customersErr := getCustomers()
	products, productsErr := getProducts()
	if customersErr != nil {
		log.Printf("Failed to get customers: %v\n", customersErr)
	}
	if productsErr != nil {
		log.Printf("Failed to get products: %v\n", productsErr)
	}

	page.Execute(resp, map[string]interface{}{
		"Customers":    customers,
		"CustomersErr": customersErr,
		"Products":     products,
		"ProductsErr":  productsErr,
		"LastUpdated":  time.Now(),
	})
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logStdOut := logger.Writer(os.Stdout)
	var err error
	x509Source, err = workloadapi.NewX509Source(ctx, workloadapi.WithClientOptions(workloadapi.WithNamedPipeName(*namedPipeNameFlag),
		workloadapi.WithLogger(logStdOut)))
	if err != nil {
		log.Fatalf("unable to create X509Source %v", err)
	}
	defer x509Source.Close()
	svid, err := x509Source.GetX509SVID()
	if err == nil {
		log.Printf("Identity: %v", svid.ID)
	}

	bundleSource, err = workloadapi.NewBundleSource(ctx, workloadapi.WithClientOptions(workloadapi.WithNamedPipeName(*namedPipeNameFlag),
		workloadapi.WithLogger(logStdOut)))
	if err != nil {
		log.Fatalf("unable to create BundleSource %v", err)
	}
	defer bundleSource.Close()

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}
	http.HandleFunc("/", indexHandler)

	log.Printf("Webapp listening on port %d...", port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
