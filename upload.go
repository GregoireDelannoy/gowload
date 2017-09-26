package main

import (
  "io"
  "os"
  "regexp"
  "net/http"
  "log"
  "fmt"
)

func handleFileUpload(writer http.ResponseWriter, request *http.Request, path string){
  if fileDescriptor, fileDescError := os.Stat(path) ; fileDescError != nil {
    msg := "Error reading path: " + path + " " + fileDescError.Error()
    log.Print(msg)
    http.Error(writer, msg, 500)
    return
  } else {
    if !fileDescriptor.IsDir(){
      msg := "Upload path is not a directory: " + path
      log.Print(msg)
      http.Error(writer, msg, 500)
      return
    }
  }

  if reader, err := request.MultipartReader() ; err != nil {
    msg := "Error reading multipart: " + err.Error()
    log.Print(msg)
    http.Error(writer, msg, 400)
    return
  } else {
    //copy each part to destination.
    for {
      part, err := reader.NextPart()
      if err == io.EOF {
        break
      }

      //if part.FileName() is empty, skip this iteration.
      if part.FileName() == "" {
        continue
      }
      
      // Check for correctness of Filename, to avoid /PATH/../../.. file creations
      if matched, err := regexp.MatchString(`^(\p{L}|\_|\-|\d|\ |\')+(\.\w+)?$`, part.FileName()); matched != true || err != nil {
        msg := "Error validating filename, keep it simple: " + part.FileName()
        log.Print(msg)
        http.Error(writer, msg, 500)
        return
      }

      if _, err := os.Stat(path + "/" + part.FileName()) ; err == nil {
        msg := "File already exists: " + part.FileName()
        log.Print(msg)
        http.Error(writer, msg, 500)
        return
      }

      dst, err := os.Create(path + "/" + part.FileName())
      defer dst.Close()

      if err != nil {
        msg := "Error creating file: " + err.Error()
        log.Print(msg)
        http.Error(writer, msg, 500)
        return
      }

      if _, err := io.Copy(dst, part); err != nil {
        msg := "Error copying content from request to file: " + err.Error()
        log.Print(msg)
        http.Error(writer, msg, 500)
        return
      }
    }
    fmt.Fprintf(writer, "Upload successfull")
  }
}