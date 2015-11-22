package main

type Event struct {
  Key []byte
  Value []byte
  _ struct{}
}

func newEvent(key, value []byte) *Event {
  event := &Event{
    Key: key,
    Value: value,
  }

  return event
}
