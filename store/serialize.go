package store

import (
	"bufio"
	"bytes"
	"encoding/gob"
)

// Encode an object to byte slice
func Encode(object interface{}) ([]byte, error) {
	var b bytes.Buffer
	buf := bufio.NewWriter(&b)
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(object); err != nil {
		return nil, err
	}
	buf.Flush()
	return b.Bytes(), nil
}

// Decode from byte slice to an object
func Decode(data []byte, object interface{}) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(object)
}
