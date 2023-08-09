package common

import (
	"bytes"

	"github.com/vmihailenco/msgpack/v5"
)

func MsgpackMarshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	enc.SetCustomStructTag("json")

	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MsgpackUnmarshal(b []byte, v interface{}) error {
	dec := msgpack.NewDecoder(bytes.NewBuffer(b))
	dec.SetCustomStructTag("json")

	if err := dec.Decode(v); err != nil {
		return err
	}
	return nil
}
