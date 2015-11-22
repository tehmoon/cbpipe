package main

import (
  "regexp"
)

type Opts struct {
  RegexpKey *regexp.Regexp
  Key string
  Filter map[string]interface{}
  Or bool
  Url string
  Bucket string
  Pool string
  Stdin bool
}

var opts = Opts{}
