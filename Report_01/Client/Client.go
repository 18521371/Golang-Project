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
	"strconv"
	"time"
)

func main() {
	Option, Path, IntervalTime, TotalTime := Input()
	Client(Option, Path, IntervalTime, TotalTime)
}

//Input : enter information
func Input() (Option string, Path string, IntervalTime string, TotalTime string) {

	Option = InputOption()
	Path = InputPath()
	IntervalTime = InputTime()
	TotalTime = InputAllTime()
	fmt.Println("-------------------------")
	return
}

// Client do at Client
func Client(Option string, Path string, Time string, TotalTime string) {
	Time1, _ := time.ParseDuration(Time + "s")
	AllTime1, _ := time.ParseDuration(TotalTime + "s")

	Start := time.Now()
	//Path1 := "F:/Golang/Client/Image"
	FileCount := CountFile(Path)
	for {
		if time.Since(Start) > AllTime1 {
			break
		}
		if Option == "1" {
			fmt.Println("Option 1")
			Option1(FileCount, Path)
		} else if Option == "2" {
			fmt.Println("Option 2")
			Option2(FileCount, Path)
		}
		time.Sleep(Time1)
		continue
	}
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Print("Finished Time for all process: ")
	fmt.Println(time.Since(Start))
}

// Option1 : Send sequentially
func Option1(count int, Path string) {
	files, _ := ioutil.ReadDir(Path)
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".jpg" {
			continue
		}
		Start := time.Now()
		// Make a byte writer
		var bodyRequest bytes.Buffer
		// Make Multipart Writer of BodyRequest
		MultiWriter := multipart.NewWriter(&bodyRequest)

		FileImageUp, _ := os.Open(Path + "/" + file.Name())
		defer FileImageUp.Close()
		//Make an Form Upload Field
		WriterFile, _ := MultiWriter.CreateFormFile("FileUpload", file.Name())
		// Copy FileImageUp into WriterFile
		_, err := io.Copy(WriterFile, FileImageUp)
		if err != nil {
			log.Fatal(err)
		}
		MultiWriter.Close()
		MakeRequest(MultiWriter, bodyRequest)
		fmt.Print("Finished Duration: ")
		fmt.Println(time.Since(Start))
		/*err = os.Remove(Path + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		} */
	}
}

// Option2 : Send Parallel
func Option2(count int, Path string) {
	files, _ := ioutil.ReadDir(Path)
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".jpg" {
			continue
		}
		go func(file os.FileInfo) {
			Start := time.Now()
			// Make a byte writer
			var bodyRequest bytes.Buffer
			// Make Multipart Writer of BodyRequest
			MultiWriter := multipart.NewWriter(&bodyRequest)
			//----------------------------------------------------------------
			FileImageUp, _ := os.Open(Path + "/" + file.Name())
			defer FileImageUp.Close()
			//Make an Form Upload Field
			WriterFile, _ := MultiWriter.CreateFormFile("FileUpload", file.Name())
			// Copy FileImageUp into WriterFile
			_, err := io.Copy(WriterFile, FileImageUp)
			if err != nil {
				log.Fatal(err)
			}
			//----------------------------------------------------------------
			//Send Request
			MultiWriter.Close()
			MakeRequest(MultiWriter, bodyRequest)
			fmt.Print("Finished Duration: ")
			fmt.Println(time.Since(Start))
		}(file)
		time.Sleep(time.Millisecond)
	}
}

// ConvItoStr : convert int to string
func ConvItoStr(n int) string {
	return strconv.Itoa(n)
}

// CountFile : count files in the path
func CountFile(Path string) int {
	files, _ := ioutil.ReadDir(Path)
	return len(files)
}

// MakeRequest : Make and send request to server
func MakeRequest(MultiWriter *multipart.Writer, bodyRequest bytes.Buffer) {
	// Make request send to Server
	Req, err := http.NewRequest("POST", "http://localhost:12345/MultiFile", &bodyRequest)
	if err != nil {
		log.Fatal(err)
	}
	Req.Header.Set("Content-Type", MultiWriter.FormDataContentType())
	//Header ContentType có Vai trò gì??
	//Create Client Object
	ClientObject := &http.Client{} // Tại sao lại có {} ??
	Rep, err := ClientObject.Do(Req)
	if err != nil {
		log.Fatal(err)
	}
	//Rep, _ := http.Post("http://localhost:12345/MultiFile", multipartWriter.FormDataContentType(), &bodyRequest)
	defer Rep.Body.Close()
	ByteFile, _ := ioutil.ReadAll(Rep.Body)
	fmt.Println(string(ByteFile))
}

// InputOption :
func InputOption() (Option string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter option: ")
	fmt.Println("Enter '1' Send Sequence  ")
	fmt.Println("Enter '2' Send Parallel  ")
	scanner.Scan()
	Option = scanner.Text()
	return
}

// InputPath :
func InputPath() (Path string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter Path: ")
	fmt.Println("Ex: F:/Golang/Client/Image")
	scanner.Scan()
	Path = scanner.Text()
	return
}

// InputTime :
func InputTime() (Time string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter IntervalTime (Second): ")
	fmt.Println("Ex: '1' ")
	scanner.Scan()
	Time = scanner.Text()
	return
}

// InputAllTime :
func InputAllTime() (TotalTime string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter total time (Second): ")
	fmt.Println("Ex: greater than Interval time ")
	scanner.Scan()
	TotalTime = scanner.Text()
	return
}
