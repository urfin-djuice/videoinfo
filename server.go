package main

import (
	"encoding/json"
	"github.com/zelenin/go-mediainfo"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func NewInfoServer(listenAddr string) *InfoServer {
	return &InfoServer{
		ListenAddr: listenAddr,
	}
}

type InfoServer struct {
	ListenAddr string
}

func (is InfoServer) Up() {
	log.Printf("Server start on %s\n", is.ListenAddr)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			log.Printf("Got request: %s %s%s\n", r.Method, r.Host, r.URL.String())
			src, _, err := r.FormFile("file")
			if err != nil {
				log.Println("Fail to get uploaded file")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			filePath, err := is.saveFile(src)
			_ = src.Close()
			if err != nil {
				log.Println("Fail to save uploaded file into temporary directory")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func() {
				_ = os.Remove(filePath)
			}()

			response, err := is.getInfo(filePath)
			if err != nil {
				log.Println("Fail to get media file info")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, err = w.Write(response)
			if err != nil {
				log.Println("Fail to send response")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Println("Request processed successfully")
		} else {
			log.Printf("Unknown request: %s %s%s\n", r.Method, r.Host, r.URL.String())
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})
	log.Fatal(http.ListenAndServe(is.ListenAddr, nil))
}

func (is InfoServer) genFilePath() string {
	return filepath.Join(os.TempDir(), time.Now().Format("20060102150405.999999999"))
}

func (is InfoServer) saveFile(src io.Reader) (string, error) {
	fileName := is.genFilePath()
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()
	_, err = io.Copy(file, src)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func (is InfoServer) getInfo(fileName string) ([]byte, error) {
	mi, err := mediainfo.Open(fileName)
	if err != nil {
		return []byte{}, err
	}
	defer mi.Close()

	informer := newInformer(mi)
	informerResult := informer.GetInfo()

	result, err := json.Marshal(informerResult)
	if err != nil {
		return []byte{}, err
	}

	return result, nil
}
