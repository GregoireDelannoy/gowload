package main

import (
  "strings"
  "log"
  "path/filepath"
)

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