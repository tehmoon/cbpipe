package main

import (
  "github.com/couchbase/go-couchbase"
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
