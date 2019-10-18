package udp

import (
	"bytes"
	"net"
)

type UDPClient struct {
	addr string
	conn *net.UDPConn
}

func (c *UDPClient) connect() (err error) {
	c.conn, err = net.ListenUDP("udp", nil)
	if err != nil {
		return
	}
	return
}

func (c *UDPClient) Send(data []byte) (err error) {
	data = c.BuildData(data)
	to, err := net.ResolveUDPAddr("udp", c.addr)
	if err != nil {
		return err
	}

	_, err = c.conn.WriteToUDP(data, to)
	return err
}

func (c *UDPClient) BuildData(data []byte) []byte {
	first := header.GetHeader(len(data))
	buf := bytes.NewBuffer(first)
	buf.Write(data)
	return buf.Bytes()
}

func NewUDPClient(addr string) (*UDPClient, error) {
	u := &UDPClient{
		addr: addr,
	}
	return u, u.connect()
}

var (
	udpClient *UDPClient
)

func SendTo(addr string, data []byte) error {
	if udpClient == nil {
		var err error
		udpClient, err = NewUDPClient(addr)
		if err != nil {
			return err
		}
	}
	return udpClient.Send(data)
}
