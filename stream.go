package main

type Stream struct {
  C chan *Event
  _ struct{}
}

func newStream() *Stream {
  stream := &Stream{
    C: make(chan *Event),
  }

  return stream
}
