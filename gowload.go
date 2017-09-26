package main


import (
"net/http"
"log"
"strings"
"path"
)

const DATA_DIR = "datadir"
const DATA_URL = "files"
const LINKS_DIR = "linksdir"
const LINKS_URL = "links"
const TEMPLATES_DIR = "./templates"
const STATIC_DIR = "./static"
const PORT = "3000"

func linksRouter(writer http.ResponseWriter, request *http.Request){
  cleanPath := strings.TrimPrefix(path.Clean(strings.TrimPrefix(request.URL.Path, "/" + LINKS_URL + "/")), "/")
  if len(cleanPath) < 15{
    msg := "Links have to be links!"
    log.Print(msg)
    http.Error(writer, msg, 400)
    return
  }
  path := LINKS_DIR + "/" + cleanPath

  if(request.Method == "GET"){
    servePath(writer, request, path, cleanPath, "")
  } else {
    msg := "Only GET allowed on public links"
    log.Print(msg)
    http.Error(writer, msg, 500)
    return
  }
}

func filesRouter(writer http.ResponseWriter, request *http.Request){
  cleanPath := strings.TrimPrefix(path.Clean(strings.TrimPrefix(request.URL.Path, "/" + DATA_URL +"/")), "/")

  user, err := getUsername(request)
  if err != nil {
    msg := "Unable to identify user from request: " + err.Error()
    log.Print(msg)
    http.Error(writer, msg, 500)
    return
  }

  path := DATA_DIR + "/" + user + "/" + cleanPath

  if isSymlink(path){
    msg := "Symlinks or non-existing file. Path: " + cleanPath
    log.Print(msg)
    http.Error(writer, msg, 500)
    return
  }

  switch (request.Method){
  case "GET":
    servePath(writer, request, path, cleanPath, user)
  case "POST":
    if err := request.ParseForm() ; err != nil{
      msg := "Error parsing form: " + err.Error()
      log.Print(msg)
      http.Error(writer, msg, 500)
    } else {
      if request.PostFormValue("actionType") == "link" {
        getPublicLink(writer, request, path)
      } else {
        handleFileUpload(writer, request, path)
      }
    }
  default:
    http.Error(writer, "Unsupported method requested: " + request.Method, 400)
    return
  }
}

func main() {
  parseTemplates()
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(STATIC_DIR))))
  http.HandleFunc("/" + DATA_URL + "/", filesRouter)
  http.HandleFunc("/" + LINKS_URL + "/", linksRouter)
  http.ListenAndServe("127.0.0.1:" + PORT, nil)
}