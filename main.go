package main

import (
  "github.com/couchbase/go-couchbase"
  "github.com/couchbase/gomemcached/client"
  "log"
  "fmt"
  "flag"
  "regexp"
  "encoding/json"
)

type Opts struct {
  RegexpKey *regexp.Regexp
  Key string
  Filter map[string]interface{}
  Or bool
}

var opts = Opts{}

func init() {
  key := flag.String("key", "", "A POSIX regexp to filter a key")
  filter := flag.String("filter", "{}", "A JSON object as filter: {\"username\": \"moon\"")
  or := flag.Bool("or", false, "Filter on Key OR filter")

  flag.Parse()

  opts.Or = *or

  opts.Key = *key
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

func main() {
  c, err := couchbase.Connect("http://10.0.0.6:8091")
  if err != nil {
    log.Fatal(err)
  }

  pool, err := c.GetPool("default")
  if err != nil {
    log.Fatal(err)
  }

  bucket, err := pool.GetBucket("twitter")
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
