package handlers

import (
	"GoServer/api/middleware"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
    "strconv"
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
    {ID: "4", Name: "Product 4", Price: 40.00},
}

// Create the cart
// Define a local list to store products in cart
var cart = struct {
	sync.RWMutex
	items map[string]int
}{items: make(map[string]int)}

// Add products to the database (replace with your database logic)
func AddProducts(w http.ResponseWriter, r *http.Request) {
    // Add products to the database
    // for now we will add the to the products slice
    if !middleware.IsAuthenticated(r) {
        http.Redirect(w, r, "/log-in", http.StatusSeeOther)
        return
    }

    // Add products to the slice from the post Request
    // This is a simple example, in a real application you would

    // Parse the form
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }

    // Get the product details from the ParseForm
    productID := r.FormValue("product_id")
    productName := r.FormValue("product_name")
    productPrice := r.FormValue("product_price")

    // make productPrice a float64
    price, err := strconv.ParseFloat(productPrice, 64)
    if err != nil {
        http.Error(w, "Error parsing price", http.StatusBadRequest)
    }

    
    // Add the product to the products slice
    products = append(products, Product{
        ID:    productID,
        Name:  productName,
        Price: price,
    })

    // Redirect to the admin packag

    http.Redirect(w, r, "/admin", http.StatusSeeOther)

}

// Admin handler for rendering the page
func AdminHandler(w http.ResponseWriter, r *http.Request) {
    // Check if user is authenticated
    if !middleware.IsAuthenticated(r) {
        http.Redirect(w, r, "/log-in", http.StatusSeeOther)
        return
    }

    tmplPath := filepath.Join("web", "templates", "admin.html")
    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        http.Error(w, "Error parsing template", http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, nil)
    if err != nil {
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
        return
    }
}

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
	totalPrice := 0.0
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

// Checkout handler
// Checkout handler
func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
    // Check if user is authenticated
    if !middleware.IsAuthenticated(r) {
        http.Redirect(w, r, "/log-in", http.StatusSeeOther)
        return
    }

    // Go to checkout page where you can see all your Products
    // and the corresponding prices of each product
    // and the total price of all Products
    tmplPath := filepath.Join("web", "templates", "checkout.html")
    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        http.Error(w, "Error parsing template", http.StatusInternalServerError)
        return
    }

    // Calculate total prices and gather cart details
    totalPrice := 0.0
    var cartItems []struct {
        Product  Product
        Quantity int
    }

    cart.RLock()
    for id, quantity := range cart.items {
        for _, product := range products {
            if product.ID == id {
                totalPrice += product.Price * float64(quantity)
                cartItems = append(cartItems, struct {
                    Product  Product
                    Quantity int
                }{
                    Product:  product,
                    Quantity: quantity,
                })
                break
            }
        }
    }
    cart.RUnlock()

    data := struct {
        CartItems  []struct {
            Product  Product
            Quantity int
        }
        TotalPrice float64
    }{
        CartItems:  cartItems,
        TotalPrice: totalPrice,
    }

    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
        return
    }
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
