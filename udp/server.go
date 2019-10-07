package udp

import (
	"fmt"
	"net"
	"time"
)

const (
	UDP_SERVER_ADDR   = ":6666"
	UDP_BUFFER_LENGTH = 4096
	UDP_INTERVAL      = 3
)

// default buffer size is 4kb,if you use more than it ,you need  sign when you NewUDPServer
func NewUDPServer(addr string, args ...int) (*UDPServer, error) {
	bufferSize := UDP_BUFFER_LENGTH
	if len(args) > 0 && args[0] > bufferSize {
		bufferSize = args[0]
	}
	u := &UDPServer{
		addr:     addr,
		length:   bufferSize,
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
		fmt.Println("UDP conn is nil")
		return
	}
	for {
		n, clientAddr, err := s.conn.ReadFromUDP(s.buffer)
		if err != nil {
			// if err = s.listen(); err != nil {
			// 	log.Errorf("UDPServer lost,listen error: %v", err)
			// 	time.Sleep(time.Second * UDP_INTERVAL)
			// }
			continue
		}
		l := header.Length()
		if n <= s.length && l <= s.length && n > l {
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

//TestReceive is a demo to implement Spit..
// you can write other struct in your project
type TestReceive struct{}

func (TestReceive) Write(from string, data []byte) error {
	fmt.Println(time.Now(), "receive test:", from, len(data), string(data))
	return nil
}
