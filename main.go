package main

import (
  "github.com/couchbase/go-couchbase"
  "github.com/couchbase/gomemcached/client"
  "log"
  "fmt"
  "flag"
  "regexp"
  "encoding/json"
  "os"
)

type Opts struct {
  RegexpKey *regexp.Regexp
  Key string
  Filter map[string]interface{}
  Or bool
  Url string
  Bucket string
  Pool string
  Stding bool
}

var opts = Opts{}

func usage(err string) {
  if err != "" {
    err = fmt.Sprintf("%s\n\n", err)
    fmt.Fprintf(os.Stderr, err)
  }
  flag.Usage()
  os.Exit(2)
}

func init() {
  flag.Usage = func () {
    fmt.Fprintf(os.Stderr, "Usage of %s: \n", os.Args[0])
    fmt.Fprintf(os.Stderr, "\t -stdin [-key regexp] [-filter JSON object] [-or]\n")
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

  opts.Stdin  = *stdin
  opts.Url    = *url
  opts.Or     = *or
  opts.Pool   = *pool
  opts.Bucket = *bucket
  opts.Key    = *key

  opts.RegexpKey = regexp.MustCompilePOSIX(opts.Key)

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

func getBucket(url, pool, bucket string) (*couchbase.Bucket, error) {
  c, err := couchbase.Connect(opts.Url)
  if err != nil {
    return nil, err
  }

  p, err := c.GetPool(pool)
  if err != nil {
    return nil, err
  }

  b, err := p.GetBucket(bucket)
  if err != nil {
    return nil, err
  }

  return b, nil
}

func main() {
  bucket, err := getBucket(opts.Url, opts.Pool, opts.Bucket)
  if err != nil {
    log.Fatal(err)
  }

  tapArgs := memcached.DefaultTapArguments()

  feed, err := bucket.StartTapFeed(&tapArgs)
  if err != nil {
    log.Fatal(err)
  }

  for {
    select {
      case e := <- feed.C:
        display := true

        if opts.Key != "" {
          if opts.RegexpKey.Match(e.Key) == false {
            display = false
          }
        }

        if opts.Or == false && display == false {
          break
        }

        if opts.Or == true && display == true {
          fmt.Println(string(e.Value[:]))
          break
        }

        display = true

        if len(opts.Filter) > 0 {
          var v interface{}

          err := json.Unmarshal(e.Value, &v)
          if err != nil {
            display = false
          }

          var doc map[string]interface{}

          if _, ok := v.(map[string]interface{}); ok {
            doc = v.(map[string]interface{})
          } else {
            display = false
          }

          for key, value := range opts.Filter {
            if _, exists := doc[key]; exists {
              if doc[key] == value {
                continue
              }

              display = false
              break
            } else {
              display = false
              break
            }
          }
        }

        if display == true {
          fmt.Println(string(e.Value[:]))
        }
    }
  }
}
