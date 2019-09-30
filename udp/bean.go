package udp

import (
	"bytes"
	"encoding/binary"
	"time"
)

func NewHeader() *Header {
	h := &Header{}
	h.init()
	return h
}

type Header []byte

func (h *Header) Length() int {
	return len(*h) + 4 // 四个字节存包体长度
}
func (h *Header) Size(bs []byte) int {
	bodyByte := bs[len(*h):]
	return BytesToInt(bodyByte)
}
func (h *Header) GetHeader(bodySize int) []byte {
	header := *h
	buf := bytes.NewBuffer(header)
	binary.Write(buf, binary.LittleEndian, int32(bodySize))
	return buf.Bytes()
}
func (h *Header) init() {
	*h = []byte{0x7A, 0x78, 0x79, 0x31, 0x39, 0x39, 0x34} // zxy1994
}

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return int(x)
}

// message  body reference
type Message struct {
	Data interface{}
	Type string
	ID   string
	Time time.Time
}
