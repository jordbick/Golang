package receipt

import (
	"path/filepath"
	"time"
)

// create a variable to point to our recipt directory, which is where we will store the files on disk
// Directory called uploads
var ReceiptDirectory string = filepath.Join("uploads")

// Receipt type with 2 fields
type Receipt struct {
	ReceiptName string    `json:"name"`
	UploadDate  time.Time `json:"uploadDate"`
}

// function that will go through the list of files in our uploads folder and create a slice of receipt objects
func GetReceipts() ([]Receipt, error) {
	receipts := make([]Receipt, 0)
	files, err : = ioutil.ReadDir(RecReceiptDirectory)
	if err != nil {
		return nil, err
	}
	// loop over files returned and append those to the slice as receipt objects
	for _, f := range files {
		receipts = append(receipts, Receipt{ReceiptName: f.Name(), UploadDate: f.ModTime()})
	}
	return receipts, nil
}