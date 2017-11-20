package main

import (
  "os"
  "fmt"
  "strings"
  "io/ioutil"
  "log"
  "sort"
  "net/http"
)

// Allo to sort []fileinfo by the IsDir property
type ByDirectory []os.FileInfo
func (fi ByDirectory) Len() int {return len(fi)}
func (fi ByDirectory) Swap(i, j int) {fi[i], fi[j] = fi[j], fi[i]}
func (fi ByDirectory) Less(i, j int) bool {
  return fi[i].IsDir()
}

type Link struct {
  Url string
  Name string
}

type DirTemplateContent struct {
  User string
  Path []Link
  Files []os.FileInfo
  IsPublic bool
  FirstImage string
}

func splitPath(path string, isPublic bool)(links []Link){
  split := strings.Split(path, "/")
  for i := 0; i < len(split); i++ {
    if(split[i] != "."){
      if isPublic {
        links = append(links, Link{Url: "/" + LINKS_URL + "/" + strings.Join(split[0:i+1], "/") + "/", Name: split[i]})
      } else {
        links = append(links, Link{Url: "/" + DATA_URL + "/" + strings.Join(split[0:i+1], "/") + "/", Name: split[i]})
      }
    }
  }
  return 
}

func buildTemplateObject(internalPath string, userPath string, user string)(DirTemplateContent, error){
  if files, err := ioutil.ReadDir(internalPath) ; err != nil {
    msg := "Error reading directory: " + internalPath + " " + err.Error()
    log.Print(msg)
    return DirTemplateContent{}, fmt.Errorf(msg)
  } else {
    sort.Sort(ByDirectory(files))
    isPublic := false
    if(user == ""){
      isPublic = true
    }

    firstImage, _ := firstImage(files)

    return DirTemplateContent{Path: splitPath(userPath, isPublic), User: user, Files: files, IsPublic: isPublic, FirstImage: firstImage}, nil
  }
}

func serveDirectory(writer http.ResponseWriter, request *http.Request, internalPath string, userPath string, user string){

  // Redirect requests that does not end with a trailing slash (but query a directory) to the correct path
  if string(request.URL.Path[len(request.URL.Path) - 1]) != "/"{
    http.Redirect(writer, request, request.URL.Path + "/", 302)
    return
  }

  if obj, err := buildTemplateObject(internalPath, userPath, user) ; err != nil {
    msg := "Error building template object: " + err.Error()
    log.Print(msg)
    http.Error(writer, msg, 500)
    return
  } else {
    compileAndWriteTemplate("directory.html", obj, writer)
  }
}