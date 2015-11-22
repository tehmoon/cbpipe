package main

import (
  "encoding/json"
)


//http://www.couchbase.com/wiki/display/couchbase/TAP+Protocol#

func main() {
  stream := newStream()

  if opts.Bucket != "" {
    go listenBucket(stream)
  }
  if opts.Stdin != false {
    go listenStdin(stream)
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
          print(string(e.Value[:]))
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
          print(string(e.Value[:]))
        }
    }
  }
}
