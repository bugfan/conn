package main

import (
	"conn/udp"
	"fmt"
	"testing"
	"time"
)

func TestUDPClient(*testing.T) {
	fmt.Println("-------begin------")
	host := "127.0.0.1:6666"
	m := []byte(`{"Request":{"Method":"GET","URL":{"Scheme":"","Opaque":"","User":null,"Host":"","Path":"/","RawPath":"","ForceQuery":false,"RawQuery":"","Fragment":""},"Path":"/","Host":"e.l.hao.io:9090","Proto":"HTTP/1.1","RawQuery":"","Query":{},"Header":{"Accept":["*/*"],"User-Agent":["curl/7.47.0"],"X-Forwarded-For":["127.0.0.1"]}},"Response":{"Status":200,"Header":{"Content-Type":["text/html"],"Date":["Mon, 30 Sep 2019 04:21:51 GMT"],"Etag":["W/\"5d4df6c8-2b70\""],"Last-Modified":["Fri, 09 Aug 2019 22:42:16 GMT"],"Server":["*****"],"Vary":["Accept-Encoding"]}},"Refer":"","Agent":"curl/7.47.0","Tags":{"location_operators":"其他","location_province":"其他","server_id":"70","server_name":"说的","ua_browser":"BrowserUnknown","ua_device_type":"DeviceUnknown","ua_os":"OSUnknown","ua_platform":"PlatformUnknown"},"Size":11120,"StartTime":"2019-09-30T12:21:51.4583009+08:00","EndTime":"2019-09-30T12:21:51.551810263+08:00","Date":"2019-09-30","IP":"127.0.0.1","City":"","Country":""}`)
	// u, err := NewUDPClient(host)
	// fmt.Println("udp error:", err)
	// err = u.Send(m)
	counter := 0
	for {
		go func() {
			err := udp.SendTo(host, m)
			fmt.Println(time.Now(), len(m), "udp1 send error:", err)
		}()
		go func() {
			err := udp.SendTo(host, m)
			fmt.Println(time.Now(), len(m), "udp2 send error:", err)
		}()
		go func() {
			err := udp.SendTo(host, m)
			fmt.Println(time.Now(), len(m), "udp3 send error:", err)
		}()
		time.Sleep(1e3)
		counter += 3
		fmt.Println("count:", counter)
	}
	fmt.Println("-------end------")
}
