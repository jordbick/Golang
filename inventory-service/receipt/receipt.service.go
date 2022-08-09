package receipt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// create path variable
const receiptPath = "receipts"

// Handler to retrieve a list of files
func handleReceipts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// If get method called, use the GetReceipts function we made, Marshal the data and return to the client
	case http.MethodGet:
		receiptList, err := GetReceipts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(receiptList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}

	// If Post is called, then upload our file to the service
	case http.MethodPost:
		r.ParseMultipartForm(5 << 20) // 5Mb, limit size of in-memory data
		// Grab receipt out of our HTTP multipart form by using FormFile method
		file, handler, err := r.FormFile("receipt")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		// create a file on disk using the os.OpenFile func
		// Pass in the file path, built use the .Join function, using the ReceiptDirectory constant and the file name (which we can get from the filehandler)
		// Flags to specify write only and create with bitmask of 0666
		f, err := os.OpenFile(filepath.Join(ReceiptDirectory, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer f.Close()

		// io.Copy func to copy the byte data from the file that was passed in through the HTTP POST into the new file we created on our system
		io.Copy(f, file)
		w.WriteHeader(http.StatusCreated)

		// Implement CORS headers
	case http.MethodOptions:
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// SetupRoutes which will be called from the main
func SetupRoutes(apiBasePath string) {
	receiptHandler := http.HandlerFunc(handleReceipts)
	http.Handle(fmt.Sprintf("%s%s", apiBasePath, receiptPath), cors.Middleware(receiptHandler))
}
