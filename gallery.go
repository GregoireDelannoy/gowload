package main

import (
  "os"
  "fmt"
  "strings"
  "io/ioutil"
  "log"
  "path/filepath"
  "path"
  "mime"
  "net/http"
)

type GalleryTemplateContent struct {
  PreviousImagePath string
  CurrentImagePath string
  NextImagePath string
}

func isImage(file os.FileInfo)(bool) {
  return strings.HasPrefix(mime.TypeByExtension(path.Ext(file.Name())), "image")
}

func pickImages(files []os.FileInfo)(res []os.FileInfo){
  for _, v := range files {
    if isImage(v){
      res = append(res, v)
    }
  }
  return
}

func toFileName(info os.FileInfo)(string){
  _, file := filepath.Split(info.Name())
  return file
}

func findPreviousNext(images []os.FileInfo, currentImage string)(string, string, error){
  for k, v := range images {
      if toFileName(v) == currentImage {
        prev := images[(k+len(images)-1) % len(images)]
        next := images[(k+1) % len(images)]
          return toFileName(prev), toFileName(next), nil
      }
  }
  msg := "Cannot find previous and next images for current :" + currentImage
  log.Print(msg)
  return "", "", fmt.Errorf(msg)
}

func buildGalleryObject(internalPath string)(GalleryTemplateContent, error){
  dir, currentImage := filepath.Split(internalPath)
  if files, err := ioutil.ReadDir(dir) ; err != nil {
    msg := "Error reading directory: " + dir + " " + err.Error()
    log.Print(msg)
    return GalleryTemplateContent{}, fmt.Errorf(msg)
  } else {
    images := pickImages(files)
    if prev, next, err := findPreviousNext(images, currentImage) ; err != nil {
      msg := "Error finding previous/next images: " + err.Error()
      log.Print(msg)
      return GalleryTemplateContent{}, fmt.Errorf(msg)
    } else {
      return GalleryTemplateContent{PreviousImagePath: prev, CurrentImagePath: currentImage, NextImagePath: next}, nil
    }
  }
}

func serveGallery(writer http.ResponseWriter, internalPath string){
  if obj, err := buildGalleryObject(internalPath) ; err != nil {
    msg := "Error building gallery object: " + err.Error()
    log.Print(msg)
    http.Error(writer, msg, 500)
  } else {
    compileAndWriteTemplate("image.html", obj, writer)
  }
}