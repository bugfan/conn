package udp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	UDP_SERVER_ADDR   = ":6666"
	UDP_BUFFER_LENGTH = 1024
	UDP_INTERVAL      = 3
)

func NewUDPServer(addr string) (*UDPServer, error) {
	u := &UDPServer{
		addr:     addr,
		length:   UDP_BUFFER_LENGTH,
		receiver: make([]Spit, 0, 0),
	}
	u.buffer = make([]byte, u.length)
	return u, u.listen()
}

type UDPServer struct {
	addr     string
	length   int
	buffer   []byte
	conn     *net.UDPConn
	receiver []Spit
}

func (s *UDPServer) listen() error {
	sAddr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return err
	}
	s.conn, err = net.ListenUDP("udp", sAddr)
	if err != nil {
		return err
	}
	return nil
}

// 阻塞接收
func (s *UDPServer) Receive() {
	if s.conn == nil {
		logrus.Error("UDP conn is nil")
		return
	}
	for {
		n, clientAddr, err := s.conn.ReadFromUDP(s.buffer)
		if err != nil {
			// if err = s.listen(); err != nil {
			// 	logrus.Errorf("UDPServer lost,listen error: %v", err)
			// 	time.Sleep(time.Second * UDP_INTERVAL)
			// }
			continue
		}
		l := header.Length()
		if n > l {
			bodySize := header.Size(s.buffer[:l])
			if bodySize >= (n - l) { // finish this  receive
				s.write(clientAddr.String(), s.buffer[l:l+bodySize])
			}
		}
	}
}

func (s *UDPServer) Use(r Spit) {
	s.receiver = append(s.receiver, r)
}

func (s *UDPServer) SetBuffer(l int) {
	s.length = l
	s.buffer = make([]byte, l)
}

func (s *UDPServer) write(from string, data []byte) {
	for _, r := range s.receiver {
		go r.Write(from, data)
	}
}

type Spit interface {
	Write(from string, data []byte) error
}

var (
	header *Header
)

func init() {
	header = NewHeader()
}

func NewHeader() *Header {
	h := &Header{}
	h.init()
	return h
}

type Header []byte

// 头的长度
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

//TestReceive is a demo to implement Spit..
// you can write other struct in your project
type TestReceive struct{}

func (TestReceive) Write(from string, data []byte) error {
	fmt.Println(time.Now(), "receive test:", from, len(data), string(data))
	return nil
}
