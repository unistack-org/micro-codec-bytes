// Package bytes provides a bytes codec which does not encode or decode anything
package bytes

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/unistack-org/micro/v3/codec"
)

type Message struct {
	Header map[string]string
	Body   []byte
}

type Codec struct {
	Conn io.ReadWriteCloser
}

// Frame gives us the ability to define raw data to send over the pipes
type Frame struct {
	Data []byte
}

func (c *Codec) ReadHeader(conn io.ReadWriter, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *Codec) ReadBody(conn io.ReadWriter, b interface{}) error {
	// read bytes
	buf, err := ioutil.ReadAll(conn)
	if err != nil {
		return err
	}

	switch v := b.(type) {
	case *[]byte:
		*v = buf
	case *Frame:
		v.Data = buf
	default:
		return fmt.Errorf("failed to read body: %v is not type of *[]byte", b)
	}

	return nil
}

func (c *Codec) Write(conn io.ReadWriter, m *codec.Message, b interface{}) error {
	var v []byte
	switch vb := b.(type) {
	case nil:
		return nil
	case *Frame:
		v = vb.Data
	case *[]byte:
		v = *vb
	case []byte:
		v = vb
	default:
		return fmt.Errorf("failed to write: %v is not type of *[]byte or []byte", b)
	}
	_, err := conn.Write(v)
	return err
}

func (c *Codec) String() string {
	return "bytes"
}

func NewCodec() codec.Codec {
	return &Codec{}
}

func (n *Codec) Marshal(v interface{}) ([]byte, error) {
	switch ve := v.(type) {
	case *[]byte:
		return *ve, nil
	case []byte:
		return ve, nil
	case *Message:
		return ve.Body, nil
	}
	return nil, codec.ErrInvalidMessage
}

func (n *Codec) Unmarshal(d []byte, v interface{}) error {
	switch ve := v.(type) {
	case *[]byte:
		*ve = d
	case *Message:
		ve.Body = d
	}
	return codec.ErrInvalidMessage
}
