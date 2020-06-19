package receipt

/*
	upload file
	receive file

*/
import (
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

var ReceiptDirectory string = filepath.Join("uploads")

type Receipt struct {
	ReceiptName string    `json:"name"`
	UploadDate  time.Time `json:"uploadDate"`
}

func GetReceipts() ([]Receipt, error) {
	receipts := make([]Receipt, 0)

	log.Println("[Get Receipts] ReceiptDirectory: ", ReceiptDirectory)

	files, err := ioutil.ReadDir(ReceiptDirectory)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		receipts = append(receipts, Receipt{ReceiptName: f.Name(), UploadDate: f.ModTime()})
	}
	return receipts, nil
}
