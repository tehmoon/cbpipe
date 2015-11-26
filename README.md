CBPIPE
======

A friendly tool to pipe [couchbase's](http://www.couchbase.com/) TAP event to output.
It can be chained from a previous stdin to an another cbpipe or everything accepting to treat json stream from stdin like [logstash](https://www.elastic.co/products/logstash).

How to install ?
----------------
Download golang and follow instructions at: [https://golang.org/dl/](https://golang.org/dl/)
and simply:

```shell:
$> go get github.com/tehmoon/cbpipe
$> go build -a github.com/tehmoon/cbpipe
```

How to run ?
------------
```shell:
$> ./cbpipe -help
```
