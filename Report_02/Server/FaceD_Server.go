package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type position struct {
	X int
	Y int
}

type facePosition struct {
	Min position
	Max position
}

func main() {

	http.HandleFunc("/", Server)
	http.ListenAndServe(":8080", nil)
}

// Server  :
func Server(Rep http.ResponseWriter, Req *http.Request) {
	GetMultiFile(Rep, Req)
	Detection(Rep, Req)
}

// GetMultiFile do
func GetMultiFile(Rep http.ResponseWriter, Req *http.Request) {
	if err := Req.ParseMultipartForm(1024 * 1024 * 10); err != nil {
		http.Error(Rep, err.Error(), http.StatusBadRequest)
		return
	}

	for key, value := range Req.MultipartForm.File {
		if key == "FileUpload" {
			for _, oneFileOfMultiFile := range value {
				saveFileintoServer(oneFileOfMultiFile)
			}
		}
	}
}

// saveFileintoServer : Save File upload into Folder log and print info
func saveFileintoServer(oneFileOfMultiFile *multipart.FileHeader) {

	File, _ := oneFileOfMultiFile.Open()
	FileByte, _ := ioutil.ReadAll(File)
	FiletoSave, _ := ioutil.TempFile("ImageIn", "image-*.jpg")

	FiletoSave.Write(FileByte)

	defer File.Close()
	defer FiletoSave.Close()
}

// Detection :
func Detection(Rep http.ResponseWriter, Req *http.Request) {
	Listfiles, _ := ioutil.ReadDir("ImageIn")
	for _, oneFileOfListfiles := range Listfiles {
		if filepath.Ext(oneFileOfListfiles.Name()) != ".jpg" {
			continue
		}
		filename := oneFileOfListfiles.Name()
		cmd := exec.Command("pigo", "-in", "ImageIn/"+filename, "-out", "ImageOut/"+filename, "-cf", "cascade/facefinder", "-json")
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		os.Remove("ImageIn/" + filename)
		if Req.FormValue("Option") == "1" {
			MakeResponseImage("ImageOut/"+filename, Rep)
		} else {
			MakeResponseJSONData("output.json", Rep)
		}
		os.Remove("ImageOut/" + filename)
	}
}

//MakeResponseImage :
func MakeResponseImage(Path string, Rep http.ResponseWriter) {
	FileResponse, _ := os.Open(Path)
	defer FileResponse.Close()
	_, err := io.Copy(Rep, FileResponse)
	if err != nil {
		log.Fatal(err)
	}
}

// MakeResponseJSONData :
func MakeResponseJSONData(Path string, Rep http.ResponseWriter) {
	var people []facePosition

	File, _ := os.Open("output.json")
	FileByte, _ := ioutil.ReadAll(File)

	_ = json.Unmarshal(FileByte, &people)

	var out bytes.Buffer
	json.Indent(&out, FileByte, "=", "\t")
	out.WriteTo(Rep)
}
