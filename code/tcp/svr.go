package tcp

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
)

type Engine struct {
	Svr      net.Listener
	Conn     net.Conn
	route    map[int]handleFunc //[reqtype]handleFunc
	DataChan chan string
	KeyChan  chan int64
}

type handleFunc func(data []byte)

var (
	engine *Engine
	locksvr   sync.Once
)

func InitSVR() {
	locksvr.Do(func() {
		engine = new(Engine)
		engine.DataChan = make(chan string)
		engine.KeyChan = make(chan int64)
		engine.route = make(map[int]handleFunc)
		engine.route[ReqString] = engine.transData
		engine.route[ReqStringKey] = engine.transKey
	})

}
func Only() *Engine {
	return engine
}
func (this *Engine)RunSVR(ip string, port int) (err error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	fmt.Println("listen on ", addr)

	engine.Svr, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer engine.Svr.Close()

	fmt.Println("listen on ", addr)

	engine.Conn, err = engine.Svr.Accept()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go engine.parse()
	return nil
}

func (this *Engine) parse() {
	datas := make([]byte, 1024) //1k
	index := 0
	for {
	re:
		n, err := this.Conn.Read(datas[index:])
		if err != nil {
			fmt.Println(err.Error())
			this.Conn.Close()
			return
		}
		fmt.Println(datas)
		index += n
		if index < 7 {
			goto re
		}
		prefix := int(binary.BigEndian.Uint32(datas[:4]))
		if prefix != PREFIX {
			fmt.Println("err prefix !")
			this.Conn.Close()
			return
		}
		reqType := datas[4]
		datalen := int(binary.BigEndian.Uint16(datas[5:7]))
		if index < (7 + datalen) {
			goto re
		}
		data := make([]byte, 0)
		data = append(data, datas[7:7+datalen]...)
		this.route[int(reqType)](data)
		if index >= (7 + datalen) { //粘包
			copy(datas, datas[7+datalen:])
			index -= (7 + datalen)
		}

	}
}

func (this *Engine) transData(data []byte) {
	this.DataChan <- string(data)
	return
}
func (this *Engine) transKey(data []byte) {
	k := binary.BigEndian.Uint64(data)
	this.KeyChan <- int64(k)
	return
}

func (this *Engine) Send(data []byte) (err error) {
	nData := make([]byte,7)
	binary.BigEndian.PutUint32(nData[:4],PREFIX)
	nData[4] = RspBool
	binary.BigEndian.PutUint16(nData[5:7],uint16(len(data)))
	data = append(nData,data...)
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
