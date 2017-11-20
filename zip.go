package main

import (
  "log"
  "net/http"
  
  "github.com/pierrre/archivefile/zip"
)


func serveZip(writer http.ResponseWriter, internalPath string){
  if err := zip.Archive(internalPath + "/", writer, nil) ; err != nil {
    log.Print("Error sending ZIP archive: " + err.Error())
  }
}