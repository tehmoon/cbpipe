package main

import (
  "flag"
  "os"
  "fmt"
)

func print(msg string) {
  fmt.Println(msg)
}

func usage(err string) {
  if err != "" {
    err = fmt.Sprintf("%s\n\n", err)
    fmt.Fprintf(os.Stderr, err)
  }
  flag.Usage()
  os.Exit(2)
}
