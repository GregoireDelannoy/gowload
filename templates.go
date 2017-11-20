package main

import (
  "html/template"
  "time"
  "math"
  "fmt"
  "strings"
  "io/ioutil"
  "log"
  "bytes"
  "net/http"
)


// Global: holds HTML templates. Filled at startup
var templates map[string]*template.Template

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

// Initialize all templates into global variable
func parseTemplates() (){
  templates = make(map[string]*template.Template)
  if files, err := ioutil.ReadDir(CONFIG.TemplatesDir) ; err != nil {
    msg := "Error reading templates directory: " + err.Error()
    log.Fatal(msg)
  } else {
    for _, f := range files {
      fmt.Println(f.Name())
      err = nil

      tpl, tplErr := template.New(f.Name()).Funcs(template.FuncMap{
        "humanDate": humanDate,
        "humanSize": humanSize,}).ParseFiles(CONFIG.TemplatesDir + "/" + f.Name())
      if tplErr != nil {
        log.Fatal("Error parsing template: " + tplErr.Error())
      } else {
        templates[f.Name()] = tpl
      }
    }
  }
  return
}

func compileAndWriteTemplate(templateName string, data interface{}, writer http.ResponseWriter){
  buf := &bytes.Buffer{}
  if err := templates[templateName].Execute(buf, data) ; err != nil {
    msg:= "Error processing image template: " + err.Error()
    log.Print(msg)
    http.Error(writer, msg, 500)
  } else {
    buf.WriteTo(writer)
  }
}