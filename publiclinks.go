package main

import (
  "net/http"
  "io/ioutil"
  "log"
  "fmt"
  "crypto/rand"
  "os"
)

func getAllPublicLinks()(links [][]string){
  if files, err := ioutil.ReadDir(CONFIG.LinksDir) ; err != nil {
    msg := "Error reading public links directory: " + err.Error()
    log.Fatal(msg)
    return
  } else {
    for _, f := range(files){
      if dest, err := os.Readlink(CONFIG.LinksDir + "/" + f.Name()) ; err != nil {
        msg := "Error reading links in public directory: " + err.Error()
        log.Fatal(msg)
        return
      } else {
        links = append(links, []string{f.Name(), dest})
      }
    }
    return
  }
}

func getPublicLink(writer http.ResponseWriter, request *http.Request, path string){
  allLinks := getAllPublicLinks()

  for _, l := range(allLinks) {
    if l[1] == path {
      log.Print("Existing public link found: " + l[0] + " => " + l[1])
      writer.Write([]byte(fmt.Sprintf("/%s/%s", CONFIG.LinksUrl, l[0])))
      return
    }
  }

  // Link not found, create one
  randomBytes := make([]byte, 16)
  if _, err := rand.Read(randomBytes) ; err != nil {
    msg := "Error generating random for link: " + err.Error()
    log.Print(msg)
    http.Error(writer, msg, 500)
    return
  } else {
    if err := os.Symlink(path, fmt.Sprintf("%s/%X", CONFIG.LinksDir, randomBytes)) ; err != nil {
      msg := "Error creating symlink link: " + err.Error()
      log.Print(msg)
      http.Error(writer, msg, 500)
      return
    } else {
      log.Print("No existing link found, created one : " + fmt.Sprintf("/%s/%X", CONFIG.LinksUrl, randomBytes))
      getPublicLink(writer, request, path)
      return
    }
  }
}