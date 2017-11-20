package main

import (
  "net/http"
  "log"
  "strings"
  "path"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "os"
)

const CONFIGURATION_FILE="config.json"

type Config struct {
  ListenAddress string
  ListenPort string
  DataDir string
  DataUrl string
  LinksDir string
  LinksUrl string
  TemplatesDir string
  StaticDir string
}

var CONFIG Config

func parseConfiguration(){
  if raw, err := ioutil.ReadFile(CONFIGURATION_FILE); err != nil {
    fmt.Println("Cannot read configuration file. ABORT. " + err.Error())
    os.Exit(1)
  } else {
    if err := json.Unmarshal(raw, &CONFIG); err != nil{
      fmt.Println("Error unmarshalling JSON: " + err.Error())
      os.Exit(1)
    }
    fmt.Printf("Configuration for current run:\n%+v\n", CONFIG)
  }
}

func linksRouter(writer http.ResponseWriter, request *http.Request){
  cleanPath := strings.TrimPrefix(path.Clean(strings.TrimPrefix(request.URL.Path, "/" + CONFIG.LinksUrl + "/")), "/")
  if len(cleanPath) < 15{
    msg := "Links have to be links!"
    log.Print(msg)
    http.Error(writer, msg, 400)
    return
  }
  path := CONFIG.LinksDir + "/" + cleanPath

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
  cleanPath := strings.TrimPrefix(path.Clean(strings.TrimPrefix(request.URL.Path, "/" + CONFIG.DataUrl +"/")), "/")

  user, err := getUsername(request)
  if err != nil {
    msg := "Unable to identify user from request: " + err.Error()
    log.Print(msg)
    http.Error(writer, msg, 500)
    return
  }

  path := CONFIG.DataDir + "/" + user + "/" + cleanPath

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
  parseConfiguration()
  parseTemplates()
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(CONFIG.StaticDir))))
  http.HandleFunc("/" + CONFIG.DataUrl + "/", filesRouter)
  http.HandleFunc("/" + CONFIG.LinksUrl + "/", linksRouter)
  http.ListenAndServe(CONFIG.ListenAddress + ":" + CONFIG.ListenPort, nil)
}