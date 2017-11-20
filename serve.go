package main

import (
  "os"
  "log"
  "net/http"
)

func servePath(writer http.ResponseWriter, request *http.Request, internalPath string, userPath string, user string){
  var isDir bool
  if fileDescriptor, fileDescError := os.Stat(internalPath) ; fileDescError != nil {
    msg := "Error reading path: " + internalPath + " " + fileDescError.Error()
    log.Print(msg)
    http.Error(writer, msg, 500)
    return
  } else {
    isDir = fileDescriptor.IsDir()
  }

  if isDir {
    if request.URL.Query().Get("o") == "zip" {
      serveZip(writer, internalPath)
    } else {
      serveDirectory(writer, request, internalPath, userPath, user)
    }
  } else {
    if request.URL.Query().Get("o") == "image" {
      serveGallery(writer, internalPath)
    } else {
      http.ServeFile(writer, request, internalPath)
    }
  }
}