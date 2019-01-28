package main

import (
	"../compute"
	"../tcp"
	"fmt"
)

func main() {
	tcp.InitSVR()
	go tcp.Only().RunSVR("0.0.0.0", 5555)
	filepath := "./files"
	//读取文件夹内的所有文件 生成2级文件
	compute.Pack(filepath)
	//读取二级文件生成 hash string放入内存
	LevelBString := compute.LevelBtoBuff(filepath)
	LevelBMap := make(map[string]int)
	for k, v := range LevelBString { // 1000个
		LevelBMap[v] = k
	}
	for {
		reqstring := <-tcp.Only().DataChan
		_, has := LevelBMap[reqstring]
		if has {
			tcp.Only().Send([]byte("y"))
		} else {
			tcp.Only().Send([]byte("n"))
			key := <-tcp.Only().KeyChan
			levelAhashString := compute.LevelAtoBuff(filepath, int(key))
			LevelAMap := make(map[string]int)
			for x, y := range levelAhashString { //100000个
				LevelAMap[y] = x
			}
			reqstring = <-tcp.Only().DataChan
			_, has = LevelBMap[reqstring]
			if has {
				tcp.Only().Send([]byte("y"))
			} else {
				tcp.Only().Send([]byte("n"))
				key2 := <-tcp.Only().KeyChan
				differentFile := fmt.Sprintf("%d", key*100000+key2) //1级文件名
				fmt.Println(differentFile)
				return
			}

		}

	}

}
