package handlers

import (
	"GoServer/api/middleware"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

// Define the Product struct
type Product struct {
	ID    string
	Name  string
	Price float64
}

// Create some valid products
var products = []Product{
	{ID: "1", Name: "Product 1", Price: 10.00},
	{ID: "2", Name: "Product 2", Price: 20.00},
	{ID: "3", Name: "Product 3", Price: 30.00},
}

// Create the cart
// Define a local list to store products in cart
var cart = struct {
	sync.RWMutex
	items map[string]int
}{items: make(map[string]int)}

// Render store
func RenderStore(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	if !middleware.IsAuthenticated(r) {
		http.Redirect(w, r, "/log-in", http.StatusSeeOther)
		return
	}

	tmplPath := filepath.Join("web", "templates", "store.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// Calculate total price
	totalPrice := 0.0
	cart.RLock()
	for id, quantity := range cart.items {
		for _, product := range products {
			if product.ID == id {
				totalPrice += product.Price * float64(quantity)
			}
		}
	}
	cart.RUnlock()

	data := struct {
		Products   []Product
		TotalPrice float64
	}{
		Products:   products,
		TotalPrice: totalPrice,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// Clear cart
func ClearCart() {
	cart.Lock()
	defer cart.Unlock()
	cart.items = make(map[string]int)
}

func ClearCartHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	if !middleware.IsAuthenticated(r) {
		http.Redirect(w, r, "/log-in", http.StatusSeeOther)
		return
	}

	ClearCart()
	// Recalculate total price and send it back to the client
	// Respond with updated total price for htmx
	totalPrice := calculateTotalPrice()
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("%.2f", totalPrice)))

}

// Add to cart handler
func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated (replace with your auth logic)
	if !middleware.IsAuthenticated(r) {
		http.Redirect(w, r, "/log-in", http.StatusSeeOther)
		return
	}

	// Get product ID from request
	productID := r.URL.Query().Get("id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	// Add product to cart
	addProductToCart(productID)

	// Respond with updated total price for htmx
	totalPrice := calculateTotalPrice()
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("%.2f", totalPrice)))
}

func addProductToCart(productID string) {
	cart.Lock()
	defer cart.Unlock()
	cart.items[productID]++
}

func calculateTotalPrice() float64 {
	totalPrice := 0.0
	cart.RLock()
	defer cart.RUnlock()
	for id, quantity := range cart.items {
		for _, product := range products {
			if product.ID == id {
				totalPrice += product.Price * float64(quantity)
			}
		}
	}
	return totalPrice
}
