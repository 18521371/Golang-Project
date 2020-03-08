package main

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

func main() {
	http.HandleFunc("/MultiFile", GetMultiFile)
	http.ListenAndServe(":12345", nil)
}

// GetMultiFile do
func GetMultiFile(Rep http.ResponseWriter, Req *http.Request) {

	//Req.ParseMultipartForm(512) // ParseMultipartForm parse the submited MultipartForm
	if err := Req.ParseMultipartForm(12); err != nil {
		http.Error(Rep, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(Rep, "------------------------")
	for key, value := range Req.MultipartForm.File {
		if key == "FileUpload" {
			for _, MultiFile := range value {
				//
				SavetoFileLog(MultiFile, Rep)
				fmt.Fprintln(Rep, "------------------------")
			}
		}
	}
	fmt.Fprintln(Rep, "Sucessful!...")
}

// PrintInfo : Print info file
func PrintInfo(MultiFile *multipart.FileHeader, osFile *os.File, Rep http.ResponseWriter) {
	fmt.Fprintln(Rep, "File name: "+MultiFile.Filename)
	fmt.Fprintln(Rep, "Size: "+ConvItoStr(MultiFile.Size))
	fmt.Fprintln(Rep, "Path in Server: "+"F:\\Golang\\Server\\"+osFile.Name())
}

// SavetoFileLog : Save File upload into Folder log and print info
func SavetoFileLog(MultiFile *multipart.FileHeader, Rep http.ResponseWriter) {
	File, _ := MultiFile.Open()
	FileByte, _ := ioutil.ReadAll(File)
	FiletoSave, _ := ioutil.TempFile("log", "image-*.jpg")
	PrintInfo(MultiFile, FiletoSave, Rep)
	defer File.Close()

	FiletoSave.Write(FileByte)

	defer FiletoSave.Close()
}

// ConvItoStr : convert int to string
func ConvItoStr(n int64) string {
	return strconv.FormatInt(n, 10)
}
