package message

import (
	"encoding/json"

	"github.com/schollz/croc/v6/src/comm"
	"github.com/schollz/croc/v6/src/compress"
	"github.com/schollz/croc/v6/src/crypt"
	log "github.com/schollz/logger"
)

// Message is the possible payload for messaging
type Message struct {
	Type    string `json:"t,omitempty"`
	Message string `json:"m,omitempty"`
	Bytes   []byte `json:"b,omitempty"`
	Num     int    `json:"n,omitempty"`
}

func (m Message) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

// Send will send out
func Send(c *comm.Comm, key []byte, m Message) (err error) {
	mSend, err := Encode(key, m)
	if err != nil {
		return
	}
	log.Debugf("writing %s message (%d bytes)", m.Type, len(mSend))
	_, err = c.Write(mSend)
	return
}

// Encode will convert to bytes
func Encode(key []byte, m Message) (b []byte, err error) {
	b, err = json.Marshal(m)
	if err != nil {
		return
	}
	b = compress.Compress(b)
	if key != nil {
		b, err = crypt.Encrypt(b, key)
	}
	return
}

// Decode will convert from bytes
func Decode(key []byte, b []byte) (m Message, err error) {
	if key != nil {
		b, err = crypt.Decrypt(b, key)
		if err != nil {
			return
		}
	}
	b = compress.Decompress(b)
	err = json.Unmarshal(b, &m)
	return
}
