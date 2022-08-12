package product

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

// define a new type to hold our filter fields
type ProductReportFilter struct {
	NameFilter         string `json:"productName"`
	ManufacturerFilter string `json:"manufacturer"`
	SKUFilter          string `json:"sku"`
}

// Handler to handle the incoming request
func handleProductReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// if type post we need to get the productFilter out of the request body
	case http.MethodPost:
		var productFilter ProductReportFilter
		// can use the JSON and Marshal function to get the data out of the request body
		// will show a different method
		// NewDecoder allows us to read data straight from a stream and stream the bytes directly from the request by calling Decode and passing in the pointer to our productFilter
		err := json.NewDecoder(r.Body).Decode(&productFilter)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Define function to get the products from the DB using these filters
		// The searchForProductData is declared in the product.data file
		products, err := searchForProductData(productFilter)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// define our template with name to match the file name (which is a formatted HTML document with values for a single product name and QuantityOnHand)
		// Extending functionality using Funcs method
		// FuncMap. key is "mod" and value is the atcual function
		// inline function, to return whether or not the input modulus something, is equal to 0
		// The mod function is now available to us within the template
		t := template.New("report.gotmpl").Funcs(template.FuncMap{"mod": func(i, x int) bool { return i%x == 0 }})
		// Call ParseFiles to Parse a file instead of just a string. Pass path with is directory with file name
		t, err = t.ParseFiles(path.Join("templates", "report.gotmpl"))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// return the first product from the list
		// define a new product variable and bytes.Buffer
		// if we get one product back then we'll call execute on our template
		var tmpl bytes.Buffer
		// var product Product
		if len(products) > 0 {
			// product = products[0]
			// Write the data to our bytes.Buffer variable
			err = t.Execute(&tmpl, products)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

		// In order to send this back as a file we need to define a NewReader to read the byte data into the response thats sent back to the client using http.ServeContent
		// Takes a responseWriter, a request, a file name, a modified time and our reader
		rdr := bytes.NewReader(tmpl.Bytes())
		w.Header().Set("Content-Disposition", "Attachment")
		http.ServeContent(w, r, "report.html", time.Now(), rdr)

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}

}

// Need to add to our SetupRoutes function in our product.service file
