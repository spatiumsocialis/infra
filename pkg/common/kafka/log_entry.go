package kafka

import (
	"encoding/json"
	"log"
)

// ObjectLogEntry represents a log entry containing an object which gets produced
type ObjectLogEntry struct {
	Object json.RawMessage

	encoded []byte
	err     error
}

func (o *ObjectLogEntry) ensureEncoded() {
	if o.encoded == nil && o.err == nil {
		o.encoded, o.err = json.Marshal(o)
	}
}

// Length returns the number of bytes in the encoded ObjectLogEntry
func (o ObjectLogEntry) Length() int {
	o.ensureEncoded()
	return len(o.encoded)
}

// Encode encodes the ObjectLogEntry
func (o ObjectLogEntry) Encode() ([]byte, error) {
	o.ensureEncoded()
	return o.encoded, o.err
}

// NewObjectLogEntry returns a new ObjectLogEntry
func NewObjectLogEntry(object interface{}) ObjectLogEntry {
	o, err := json.Marshal(object)
	if err != nil {
		log.Fatalf("error marshalling object: %v", err)
	}
	return ObjectLogEntry{
		Object: o,
	}
}
