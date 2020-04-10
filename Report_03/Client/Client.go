package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	Client()
	http.ListenAndServe(":8093", nil)
}

// Client :
func Client() {
	PathIn := "Image/"
	Option := "2"
	uploadImage(PathIn, Option)
}

func uploadImage(Path string, Option string) {

	files, _ := ioutil.ReadDir(Path)
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".jpg" {
			continue
		}
		var bodyRequest bytes.Buffer
		MultiWriter := multipart.NewWriter(&bodyRequest)

		FileImageUp, _ := os.Open(Path + "/" + file.Name())
		defer FileImageUp.Close()
		WriterFile, _ := MultiWriter.CreateFormFile("FileUpload", file.Name())
		_, err := io.Copy(WriterFile, FileImageUp)
		if err != nil {
			log.Fatal(err)
		}

		MultiWriter.WriteField("Option", Option)
		MultiWriter.Close()
		MakeRequest(MultiWriter, bodyRequest, Option)
	}
}

// MakeRequest : Make and send request to server
func MakeRequest(MultiWriter *multipart.Writer, bodyRequest bytes.Buffer, Option string) {
	Req, err := http.NewRequest("POST", "http://server:8090/", &bodyRequest)
	if err != nil {
		log.Fatal(err)
	}
	Req.Header.Set("Content-Type", MultiWriter.FormDataContentType())

	ClientObject := &http.Client{}
	Rep, err := ClientObject.Do(Req)
	if err != nil {
		log.Fatal(err)
	}
	defer Rep.Body.Close()

	if Option == "1" {
		Tempfile, _ := ioutil.TempFile("Result", "Image-*.jpg")
		defer Tempfile.Close()
		io.Copy(Tempfile, Rep.Body)
	} else {
		ByteFile, _ := ioutil.ReadAll(Rep.Body)
		fmt.Println(string(ByteFile))
	}
}

func inputOption(Option *string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter your Option")
	fmt.Println("Note")
	fmt.Println("Option 1: Get image which maked face")
	fmt.Println("Option 2: Get JSON file includes position of faces in Image")
	fmt.Print("Enter here('1' or '2'): ")
	scanner.Scan()
	*Option = scanner.Text()
	fmt.Println("----------------")
}
