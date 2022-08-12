package receipt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jordbick/Golang/inventory-service/cors"
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
			log.Print(err)
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

// Parses URL to get the filename
// Pass in the receipt's URL and the file name as the last part of the URL path
func handleDownload(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", receiptPath))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileName := urlPathSegments[1:][0]
	file, err := os.Open(filepath.Join(ReceiptDirectory, fileName))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set header info on response - Response might be processed differently dependent on client that is calling the code
	fHeader := make([]byte, 512)
	file.Read(fHeader)
	// http.DetectContentType function to determine what kind of file this is and use it to set the Content-Type header in our response
	fContentType := http.DetectContentType(fHeader)

	// check file size so that the client knows how much data it's going to be downloading in the response = file.Stat()
	stat, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fSize := strconv.FormatInt(stat.Size(), 10)
	// set response header to
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", fContentType)
	w.Header().Set("Content-Length", fSize)

	// Because we have already read the file we can use the file.Seek() function to set the seeker back to the start of the file
	file.Seek(0, 0)

	// Copy byte data from the file to our ResponseWriter, which will return the file to the client
	io.Copy(w, file)
}

// SetupRoutes which will be called from the main
func SetupRoutes(apiBasePath string) {
	receiptHandler := http.HandlerFunc(handleReceipts)
	downloadHandler := http.HandlerFunc(handleDownload)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, receiptPath), cors.Middleware(receiptHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, receiptPath), cors.Middleware(downloadHandler))
}
