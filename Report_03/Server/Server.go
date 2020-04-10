package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	http.HandleFunc("/", Server)
	http.ListenAndServe(":8090", nil)
}

// Server  :
func Server(Rep http.ResponseWriter, Req *http.Request) {
	CreateDatabase()
	HandleDetection(Rep, Req)
}

// HandleDetection do
func HandleDetection(Rep http.ResponseWriter, Req *http.Request) {
	if err := Req.ParseMultipartForm(1024 * 1024 * 10); err != nil {
		http.Error(Rep, err.Error(), http.StatusBadRequest)
		return
	}

	for key, value := range Req.MultipartForm.File {
		Option := Req.FormValue("Option")
		if key == "FileUpload" {
			for _, oneFileOfMultiFile := range value {
				check := IsAlreadyExistInDatabase(oneFileOfMultiFile)
				if check == true {
					if Option == "1" {
						ResponseImageToClient(oneFileOfMultiFile, Rep)
					} else {
						ResponseJSONtoClient(oneFileOfMultiFile, Rep)
					}
				} else {
					saveFileintoServer(oneFileOfMultiFile)
					saveOriginalFileToDatabase(oneFileOfMultiFile)
					Detection(Rep, Req)
					SaveResultImageIntoDatabase(oneFileOfMultiFile.Filename)
					saveResultJSONintoDatabse(oneFileOfMultiFile.Filename)
					if Option == "1" {
						ResponseImageToClient(oneFileOfMultiFile, Rep)
					} else {
						ResponseJSONtoClient(oneFileOfMultiFile, Rep)
					}
					os.Remove("ImageOut/" + oneFileOfMultiFile.Filename)
					os.Remove("output.json")
				}
			}
		}
	}
}

func saveOriginalFileToDatabase(oneFileOfMultiFile *multipart.FileHeader) {
	conn, _ := sql.Open("mysql", "root:123456@tcp(db)/facedetectionresult")
	defer conn.Close()
	FileName := oneFileOfMultiFile.Filename
	FileSize := strconv.FormatInt(oneFileOfMultiFile.Size, 10)
	insert, _ := conn.Query("INSERT INTO results(FileName,SIZE) VALUES ('" + FileName + "','" + FileSize + "')")

	defer insert.Close()
}

// saveFileintoServer : Save File upload into Folder log and print info
func saveFileintoServer(oneFileOfMultiFile *multipart.FileHeader) {
	File, _ := oneFileOfMultiFile.Open()
	FileByte, _ := ioutil.ReadAll(File)
	FiletoSave, _ := os.Create(oneFileOfMultiFile.Filename)
	FiletoSave.Write(FileByte)

	defer File.Close()
	defer FiletoSave.Close()
}

//SaveResultImageIntoDatabase :
func SaveResultImageIntoDatabase(PathFile string) {
	File, _ := os.Open("ImageOut/" + PathFile)
	defer File.Close()

	FileByte, _ := ioutil.ReadAll(File)
	FileString := base64.StdEncoding.EncodeToString(FileByte)

	conn, _ := sql.Open("mysql", "root:123456@tcp(db)/facedetectionresult")
	defer conn.Close()

	update, _ := conn.Query("UPDATE results SET FileByte='" + FileString + "' WHERE FileName='" + PathFile + "'")
	defer update.Close()
}

func saveResultJSONintoDatabse(PathFile string) {
	File, _ := os.Open("output.json")
	defer File.Close()
	FileByte, _ := ioutil.ReadAll(File)
	FileString := base64.StdEncoding.EncodeToString(FileByte)

	conn, _ := sql.Open("mysql", "root:123456@tcp(db)/facedetectionresult")
	defer conn.Close()

	update, _ := conn.Query("UPDATE results SET JSON='" + FileString + "' WHERE FileName='" + PathFile + "'")
	defer update.Close()

}

// IsAlreadyExistInDatabase :
func IsAlreadyExistInDatabase(FileUploadHeader *multipart.FileHeader) (check bool) {
	conn, _ := sql.Open("mysql", "root:123456@tcp(db)/facedetectionresult")
	defer conn.Close()

	check = false
	FileName := FileUploadHeader.Filename
	FileSize := strconv.FormatInt(FileUploadHeader.Size, 10)
	results, _ := conn.Query("SELECT FileName, SIZE FROM results")

	for results.Next() {
		var FileNametoCheck string
		var FileSizetoCheck string
		_ = results.Scan(&FileNametoCheck, &FileSizetoCheck)
		if (FileNametoCheck == FileName) && (FileSizetoCheck == FileSize) {
			check = true
			break
		}
	}
	return
}

// ResponseJSONtoClient :
func ResponseJSONtoClient(oneFileOfMultiFile *multipart.FileHeader, Rep http.ResponseWriter) {
	conn, _ := sql.Open("mysql", "root:123456@tcp(db)/facedetectionresult")
	defer conn.Close()

	FileByteImageInDatabase, _ := conn.Query("SELECT JSON FROM results WHERE FileName='" + oneFileOfMultiFile.Filename + "'")
	defer FileByteImageInDatabase.Close()

	for FileByteImageInDatabase.Next() {
		var FileString string
		_ = FileByteImageInDatabase.Scan(&FileString)

		var Response bytes.Buffer
		FileByte, _ := base64.StdEncoding.DecodeString(FileString)

		json.Indent(&Response, FileByte, "", "\t")
		Response.WriteTo(Rep)
	}
}

// ResponseImageToClient :
func ResponseImageToClient(oneFileOfMultiFile *multipart.FileHeader, Rep http.ResponseWriter) {
	conn, _ := sql.Open("mysql", "root:123456@tcp(db)/facedetectionresult")
	defer conn.Close()

	FileByteImageInDatabase, _ := conn.Query("SELECT FileByte FROM results WHERE FileName='" + oneFileOfMultiFile.Filename + "'")
	defer FileByteImageInDatabase.Close()

	for FileByteImageInDatabase.Next() {
		var FileString string
		var BytePage bytes.Buffer

		_ = FileByteImageInDatabase.Scan(&FileString)
		FileByte, _ := base64.StdEncoding.DecodeString(FileString)

		BytePage.Write(FileByte)
		BytePage.WriteTo(Rep)
	}

}

// Detection :
func Detection(Rep http.ResponseWriter, Req *http.Request) {
	Listfiles, _ := ioutil.ReadDir(".")
	for _, oneFileOfListfiles := range Listfiles {
		if filepath.Ext(oneFileOfListfiles.Name()) != ".jpg" {
			continue
		}
		filename := oneFileOfListfiles.Name()
		cmd := exec.Command("pigo", "-in", filename, "-out", "ImageOut/"+filename, "-cf", "cascade/facefinder", "-json")
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		os.Remove(filename)
	}
}

// CreateDatabase :
func CreateDatabase() {
	conn, _ := sql.Open("mysql", "root:123456@tcp(db)/")
	defer conn.Close()
	_, _ = conn.Query("CREATE DATABASE facedetectionresult; ")

	conn1, _ := sql.Open("mysql", "root:123456@tcp(db)/facedetectionresult")
	defer conn1.Close()
	_, _ = conn1.Query("CREATE TABLE results (FileName varchar(100) PRIMARY KEY, SIZE varchar(30), FileByte longblob, JSON longblob)")
}
