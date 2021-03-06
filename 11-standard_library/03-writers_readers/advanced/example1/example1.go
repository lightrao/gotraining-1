// All material is licensed under the GNU Free Documentation License
// https://github.com/gobridge/gotraining/blob/master/LICENSE

// http://play.golang.org/p/Tmt7v3fIQF

// https://github.com/extemporalgenome/watchpost/blob/master/main.go
// Sample code provided by Kevin Gillette

// Sample program to show how io.Writes can be embedded within
// other Writer calls to perform complex writes.
package main

import (
	"crypto"
	_ "crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

// main is the entry point for the application.
func main() {
	// Open the file for reading.
	file, err := os.Open("data.json")
	if err != nil {
		fmt.Println("Open File", err)
		return
	}

	// Schedule the file to be closed once
	// the function returns.
	defer file.Close()

	// Create an SHA1 hash value which implements io.Writer.
	// We want a hash value for the contents of the file.
	hash := crypto.SHA1.New()

	// Create a TeeReader that will write the content to the
	// hash Writer at the same time the file is being read and consumed
	// by io.Copy later on.
	hashReader := io.TeeReader(file, hash)

	// Create the multipart writer to copy the contents of the file
	// and the hash key into a single stream.
	mpWriter := multipart.NewWriter(os.Stdout)
	fileWriter, err := mpWriter.CreateFormFile("file", "data.json")

	// Read the TeeReader (file) by using hashReader and write it to
	// both the multipart form Writer (using fileWriter which is bound
	// to the pipeWriter) and the SHA1 hash Writer at the same time.
	//
	// When this happens the file bound to the hashReader is
	// read for this operation and written to both the hash Writer
	// and the multipart Writer at the same time.
	_, err = io.Copy(fileWriter, hashReader)
	if err != nil {
		fmt.Println("Write File", err)
		return
	}

	// Add the hash key generated by the io.Copy to the multipart document.
	mpWriter.WriteField("sha1", hex.EncodeToString(hash.Sum(nil)))

	// Close the pipeWriter which will cause the pipeReader to unblock.
	mpWriter.Close()
}
