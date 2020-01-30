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

func Up() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {

			src, _, err := r.FormFile("file")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			filePath, err := saveFile(src)
			_ = src.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func() {
				_ = os.Remove(filePath)
			}()

			response, err := getInfo(filePath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, err = w.Write(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})
	log.Fatal(http.ListenAndServe(":8088", nil))
}

func genFilePath() string {
	return filepath.Join(os.TempDir(), time.Now().Format("20060102150405.999999999"))
}

func saveFile(src io.Reader) (string, error) {
	fileName := genFilePath()
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

func getInfo(fileName string) ([]byte, error) {
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
