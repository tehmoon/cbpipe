package main

import (
  "github.com/couchbase/go-couchbase"
  "github.com/couchbase/gomemcached/client"
  "log"
)

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

func listenBucket (stream *Stream) {
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
}
