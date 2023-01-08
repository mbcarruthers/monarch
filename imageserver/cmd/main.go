package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

const (
	__port = 8025
)

var (
	port = fmt.Sprintf(":%d", __port)
)

type ImageFunc func(w http.ResponseWriter, r *http.Request)

type ImageHandler struct {
	*http.ServeMux
}

func NewImageHandler() *ImageHandler {
	return &ImageHandler{
		ServeMux: http.NewServeMux(),
	}
}

func ServeImage(imagefile string) ImageFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := filepath.Join("assets/", imagefile)
		res, err := os.ReadFile(filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

func CreateImageFileMap(directoryname string) map[string]ImageFunc {
	filemap := make(map[string]ImageFunc)
	items, err := os.ReadDir(directoryname)
	if err != nil {
		log.Fatalf(err.Error())
		return nil
	}
	for _, item := range items {
		urlPath := "/" + item.Name()
		filemap[urlPath] = ServeImage(item.Name())
	}
	return filemap
}

func main() {
	imageMap := CreateImageFileMap("assets/")
	imageHandle := NewImageHandler()

	for k, v := range imageMap {
		imageHandle.HandleFunc(k, v)
	}
	server := &http.Server{
		Addr:    port,
		Handler: imageHandle,
	}
	go func() {
		log.Printf("We are live at port:%d \n", __port)
		log.Println("Available Image Urls")
		for _url, _ := range imageMap {
			log.Printf("http://image-server:%d%s \n",
				__port, _url)
		}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("There was an error listening to the server:%s \n",
				err.Error())
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("image server closed.")
}
