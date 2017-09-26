package main

import (
  "log"
  "fmt"
  "strings"
  "net/http"
  "encoding/base64"
)

const HEADER_AUTH = "X-Basic-Auth"

func getUsername(request *http.Request)(string, error){
  if val, ok := request.Header[HEADER_AUTH]; ok {
    if decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(val[0], "Basic ")); err != nil {
      msg := "Unable to decode Basic Auth from base64: " + val[0] + err.Error()
      log.Print(msg)
      return "", fmt.Errorf(msg)
    } else {
      // Convert byte array to string and split part before ":"
      return strings.Split(string(decoded[:]), ":")[0], nil
    }
  } else {
    msg := "Unable to find username in headers"
    log.Print(msg)
    return "", fmt.Errorf(msg)
  }
}