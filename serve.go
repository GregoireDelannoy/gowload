package main

import (
  "html/template"
  "os"
  "time"
  "math"
  "fmt"
  "strings"
  "io/ioutil"
  "log"
  "sort"
  "net/http"
  "path/filepath"
  "bytes"

  "github.com/pierrre/archivefile/zip"
)

// Global: holds HTML templates. Filled at startup
var templates map[string]*template.Template

type Link struct {
  Url string
  Name string
}

type DirTemplateContent struct {
  User string
  Path []Link
  Files []os.FileInfo
  IsPublic bool
}

// Allo to sort []fileinfo by the IsDir property
type ByDirectory []os.FileInfo
func (fi ByDirectory) Len() int {return len(fi)}
func (fi ByDirectory) Swap(i, j int) {fi[i], fi[j] = fi[j], fi[i]}
func (fi ByDirectory) Less(i, j int) bool {
  return fi[i].IsDir()
}

// Helper functions for templates
func humanDate (t time.Time) string {
  return strings.Split(t.Format(time.RFC3339), "T")[0]
}

var SIZE_PREFIXES = [5]string{"B", "kB", "MB", "GB", "TB"}
func humanSize(size int64) string {
  if size == 0 {
    return "0B"
  } else {
    floatSize := float64(size)
    i := math.Floor(math.Log(floatSize) / math.Log(1024))
    return fmt.Sprintf("%.2f", floatSize / math.Pow(1024, i)) + SIZE_PREFIXES[int(i)]
  }
}

func parseTemplates() (){
  templates = make(map[string]*template.Template)
  if files, err := ioutil.ReadDir(TEMPLATES_DIR) ; err != nil {
    msg := "Error reading templates directory: " + err.Error()
    log.Fatal(msg)
  } else {
    for _, f := range files {
      fmt.Println(f.Name())
      err = nil

      tpl, tplErr := template.New(f.Name()).Funcs(template.FuncMap{
        "humanDate": humanDate,
        "humanSize": humanSize,}).ParseFiles(TEMPLATES_DIR + "/" + f.Name())
      if tplErr != nil {
        log.Fatal("Error parsing template: " + tplErr.Error())
      } else {
        templates[f.Name()] = tpl
      }
    }
  }
  return
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
    return DirTemplateContent{Path: splitPath(userPath, isPublic), User: user, Files: files, IsPublic: isPublic}, nil
  }
}

func isSymlink(path string)(bool){
  newPath, err := filepath.EvalSymlinks(path)
  if err != nil {
    log.Print("Error checking symlink for path: " + path)
    return true
  }

  if newPath != strings.TrimSuffix(path, "/.") {
    log.Print("Mismatch after following symlinks. Should NOT HAPPEN. Path: " + path)
    return true
  } else {
    return false
  }
}

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
      if err := zip.Archive(internalPath + "/", writer, nil) ; err != nil {
        log.Print("Error sending ZIP archive: " + err.Error())
      }
      return
    }

    if string(request.URL.Path[len(request.URL.Path) - 1]) != "/"{
      http.Redirect(writer, request, request.URL.Path + "/", 302)
      return
    }

    if dir, err := buildTemplateObject(internalPath, userPath, user) ; err != nil {
      msg:= "Error building template object: " + err.Error()
      log.Print(msg)
      http.Error(writer, msg, 500)
      return
    } else {
      buf := &bytes.Buffer{}
      if err := templates["directory.html"].Execute(buf, dir) ; err != nil {
        msg:= "Error processing template: " + err.Error()
        log.Print(msg)
        http.Error(writer, msg, 500)
      } else {
        buf.WriteTo(writer)
      }
      return
    }
  } else {
    http.ServeFile(writer, request, internalPath)
    return
  }
}