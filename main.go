package main

import (
  "github.com/couchbase/gomemcached/client"
  "log"
  "fmt"
  "encoding/json"
  "os"
  "bufio"
)

func main() {
  stream := newStream()

  if opts.Bucket != "" {
    go func() {
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
            stream.C <- newEvent(e.Key, e.Value)
        }
      }
    }()
  }
  if opts.Stdin != false {
    go func() {
      scanner := bufio.NewScanner(os.Stdin)
      for scanner.Scan() {
        // Switched from Bytes() to Text() because:
        // Bytes returns the most recent token generated by a call to Scan.
        // The underlying array may point to data that will be overwritten
        // by a subsequent call to Scan. It does no allocation.
        line := scanner.Text()
        stream.C <- newEvent(nil, []byte(line[:]))
      }
      if err := scanner.Err(); err != nil {
        log.Fatal(err)
      }
    }()
  }

  for {
    select {
      case e := <- stream.C:
        display := true

        if opts.Key != ""  && e.Key != nil {
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
