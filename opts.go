package main

import (
  "errors"
  "fmt"
  "flag"
  "regexp"
  "os"
  "encoding/json"
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

func usageFunc() {
  fmt.Fprintf(os.Stderr, "Usage of %s: \n", os.Args[0])
  fmt.Fprintf(os.Stderr, "\t <-bucket name | -stdin> [-key regexp] [-filter JSON object] [-url url] [-or] [-pool name]\n\n")
  flag.CommandLine.PrintDefaults()
}

func parseOpts() error {
  flag.Usage = usageFunc

  key    := flag.String("key", "", "A POSIX regexp to filter a key. Doesn't affect -stdin filtering.")
  filter := flag.String("filter", "{}", "A JSON object as filter: {\"username\": \"moon\"}.")
  or     := flag.Bool("or", false, "Filter on Key OR filter.")
  url    := flag.String("url", "http://localhost:8091", "Couchbase URL.")
  bucket := flag.String("bucket", "", "Listen on Couchbase bucket name.")
  pool   := flag.String("pool", "default", "Couchbase pool name.")
  stdin  := flag.Bool("stdin", false, "Listen on stdin.")

  flag.Parse()

  if *bucket == "" && *stdin == false {
    return errors.New("Either -stdin or -bucket has to be specified.")
  }

  opts.RegexpKey = regexp.MustCompilePOSIX(*key)

  opts.Stdin  = *stdin
  opts.Url    = *url
  opts.Or     = *or
  opts.Pool   = *pool
  opts.Bucket = *bucket
  opts.Key    = *key

  var v interface{}

  err := json.Unmarshal([]byte(*filter), &v)
  if err != nil {
    return err
  }

  if filter, ok := v.(map[string]interface{}); ok {
    opts.Filter = filter
  } else {
    return errors.New("Filter must be a JSON object.")
  }

  return nil
}
