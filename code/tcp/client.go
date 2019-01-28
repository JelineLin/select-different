package tcp

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
)

type Client struct {
	Conn net.Conn
}

var (
	clent *Client
	lock  sync.Once
)

func InitClient() {
	lock.Do(func() {
		clent = new(Client)
	})
}
func OnlyClient() *Client {
	return clent
}
func RunClient(ip string, port int) {
	var err error
	addr := fmt.Sprintf("%s:%d", ip, port)
	clent.Conn, err = net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}
func (this *Client) SendReqString(data []byte) (err error) {
	buff := make([]byte, 7)
	binary.BigEndian.PutUint32(buff[:4], PREFIX)
	buff[4] = ReqString
	binary.BigEndian.PutUint16(buff[5:7], 32)
	data = append(buff, data...)
	index := 0
re:
	n, err := this.Conn.Write(data[index:])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if n < len(data) {
		index += n
		goto re
	}
	return nil

}
func (this *Client) SendReqStringKey(key int64) (err error) {
	buff := make([]byte, 15)
	binary.BigEndian.PutUint32(buff[:4], PREFIX)
	buff[4] = ReqStringKey
	binary.BigEndian.PutUint16(buff[5:7], 8)
	binary.BigEndian.PutUint64(buff[7:15], uint64(key))
	index := 0
re:
	n, err := this.Conn.Write(buff[index:])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if n < len(buff) {
		index += n
		goto re
	}
	return nil

}
func (this *Client) ReadRsp() string {
	datas := make([]byte, 1024) //1k
	index := 0

re:
	n, err := this.Conn.Read(datas[index:])
	if err != nil {
		fmt.Println(err.Error())
		this.Conn.Close()
		return ""
	}
	index += n
	if index < 7 {
		goto re
	}
	prefix := int(binary.BigEndian.Uint32(datas[:4]))
	if prefix != PREFIX {
		fmt.Println("err prefix !")
		this.Conn.Close()
		return ""
	}
	reqType := datas[4]
	if reqType != 2 {
		fmt.Println("parse err")
		return ""
	}
	datalen := int(binary.BigEndian.Uint16(datas[5:7]))
	if index < (7 + datalen) {
		goto re
	}
	return string(datas[7 : 7+datalen])
}
