package main

import (
  "fmt"
  "flag"
  "regexp"
  "log"
  "os"
  "encoding/json"
)

func init() {
  flag.Usage = func () {
    fmt.Fprintf(os.Stderr, "Usage of %s: \n", os.Args[0])
    fmt.Fprintf(os.Stderr, "\t -stdin [-filter JSON object] [-or]\n")
    fmt.Fprintf(os.Stderr, "\t -bucket name [-key regexp] [-filter JSON object] [-url url] [-or] [-pool name]\n\n")
    flag.CommandLine.PrintDefaults()
  }

  key    := flag.String("key", "", "A POSIX regexp to filter a key.")
  filter := flag.String("filter", "{}", "A JSON object as filter: {\"username\": \"moon\"}.")
  or     := flag.Bool("or", false, "Filter on Key OR filter.")
  url    := flag.String("url", "http://localhost:8091", "Couchbase URL.")
  bucket := flag.String("bucket", "", "Couchbase bucket name.")
  pool   := flag.String("pool", "default", "Couchbase pool name.")
  stdin  := flag.Bool("stdin", false, "Listen on stdin instead of Couchbase TAP.")

  flag.Parse()

  if *bucket == "" && *stdin == false {
    usage("When -stdin is not specified, use -bucket instead.")
  }

  if *stdin == true && *key != "" {
    usage("-key cannot be specified when -stdin is true.")
  } else {
    opts.RegexpKey = regexp.MustCompilePOSIX(*key)
  }

  opts.Stdin  = *stdin
  opts.Url    = *url
  opts.Or     = *or
  opts.Pool   = *pool
  opts.Bucket = *bucket
  opts.Key    = *key


  var v interface{}

  err := json.Unmarshal([]byte(*filter), &v)
  if err != nil {
    log.Fatal(err)
  }

  if filter, ok := v.(map[string]interface{}); ok {
    opts.Filter = filter
  } else {
    log.Fatal("Filter must be a JSON object.")
  }
}
