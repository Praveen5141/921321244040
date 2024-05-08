package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const testServerURL = "http://20.244.56.144"

type Product struct {
	ProductID    string  `json:"productId"`
	ProductName  string  `json:"productName"`
	Price        float64 `json:"price"`
	Rating       float64 `json:"rating"`
	Discount     int     `json:"discount"`
	Availability string  `json:"availability"`
	Company      string  `json:"company"`
	Category     string  `json:"category"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
}

func fetchProducts(company, category string, top, minPrice, maxPrice int) ([]Product, error) {
	url := fmt.Sprintf("%s/test/companies/%s/categories/%s/products?top=%d&minPrice=%d&maxPrice=%d", testServerURL, company, category, top, minPrice, maxPrice)
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var products []Product
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func handleTopProducts(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	category := parts[2]

	topStr := r.URL.Query().Get("top")
	top, err := strconv.Atoi(topStr)
	if err != nil {
		http.Error(w, "Invalid 'top' parameter", http.StatusBadRequest)
		return
	}

	minPriceStr := r.URL.Query().Get("minPrice")
	minPrice, err := strconv.Atoi(minPriceStr)
	if err != nil {
		http.Error(w, "Invalid 'minPrice' parameter", http.StatusBadRequest)
		return
	}

	maxPriceStr := r.URL.Query().Get("maxPrice")
	maxPrice, err := strconv.Atoi(maxPriceStr)
	if err != nil {
		http.Error(w, "Invalid 'maxPrice' parameter", http.StatusBadRequest)
		return
	}

	if top < 1 || minPrice < 0 || maxPrice <= 0 || minPrice >= maxPrice {
		http.Error(w, "Invalid parameter values", http.StatusBadRequest)
		return
	}

	products, err := fetchProducts("AMZ", category, top, minPrice, maxPrice)
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func handleProductDetails(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	productID := parts[3]

	// Make request to get details of specific product
	// For simplicity, this is not implemented here
	// You should implement this by calling the appropriate API endpoint from the e-commerce companies

	// Dummy response for demonstration
	product := Product{
		ProductID:    productID,
		ProductName:  "Laptop 1",
		Price:        2236,
		Rating:       4.7,
		Discount:     63,
		Availability: "yes",
		Company:      "AMZ",
		Category:     "Laptop",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func main() {
	http.HandleFunc("/categories/", handleTopProducts)
	http.HandleFunc("/products/", handleProductDetails)

	port := 9876
	fmt.Printf("Server is running on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
