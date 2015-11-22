package main

import (
  "log"
)

func init() {
  err := parseOpts()
  if err != nil {
    log.Fatal(err)
  }
}
