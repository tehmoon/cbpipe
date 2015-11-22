package main

import (
)

func init() {
  err := parseOpts()
  if err != nil {
    usage(err.Error())
  }
}
