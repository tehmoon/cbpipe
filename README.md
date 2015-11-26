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
Usage of ./cbpipe:
         <-bucket name | -stdin> [-key regexp] [-filter JSON object] [-url url] [-or] [-pool name]

  -bucket="": Listen on Couchbase bucket name.
  -filter="{}": A JSON object as filter: {"username": "moon"}.
  -key="": A POSIX regexp to filter a key. Doesn't affect -stdin filtering.
  -or=false: Filter on Key OR filter.
  -pool="default": Couchbase pool name.
  -stdin=false: Listen on stdin.
  -url="http://localhost:8091": Couchbase URL.
```

Examples
--------
Listen from a Couchbase bucket twitter filter from json or the name of the key from a POSIX regexp.   
```shell:
$> ./cbpipe -bucket twitter -filter '{"username": "moon"}' -key '^twitter-' -or
```
Listen from a Couchbase bucket twitter filter from json and the name of the key from a POSIX regexp.  
Pipe it to another JSON filter and to logstash which waits for stdin json stream.  
```shell:
$> ./cbpipe -bucket twitter -filter '{"username": "moon"}' -key '^twitter-'  | ./cbpipe -stdin -filter '{"type": "tweet"}' | logstash -f stdin-elastic.conf
```
Listen from a Couchbase bucket, pipe it to another listener from another Couchbase bucket and stdout into a file.
```shell:
$> ./cbpipe -bucket twitter | ./cbpipe -stdin -bucket users > /tmp/dump.json
```

Docs
----
* [Couchbase TAP API](http://www.couchbase.com/wiki/display/couchbase/TAP+Protocol)
