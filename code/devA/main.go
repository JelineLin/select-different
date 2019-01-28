package main

import (
	"../compute"
	"../tcp"
	"fmt"
)

func main() {
	//初始化
	tcp.InitClient()
	tcp.RunClient("127.0.0.1", 5555)
	filepath := "./files"
	//读取文件夹内的所有文件 生成2级文件
	compute.Pack(filepath)
	//读取二级文件生成 hash string放入内存
	LevelBString := compute.LevelBtoBuff(filepath)
	for k, v := range LevelBString { // 1000个
		tcp.OnlyClient().SendReqString([]byte(v))
		if tcp.OnlyClient().ReadRsp() == "n" { //hash string不相同
			//读取二级文件内容到 内存
			tcp.OnlyClient().SendReqStringKey(int64(k))
			levelAhashString := compute.LevelAtoBuff(filepath, k)
			for x, y := range levelAhashString { //100000个
				tcp.OnlyClient().SendReqString([]byte(y))
				if tcp.OnlyClient().ReadRsp() == "n" { //hash string不相同
					tcp.OnlyClient().SendReqStringKey(int64(x))
					differentFile := fmt.Sprintf("%d", k*100000+x) //1级文件名
					fmt.Println(differentFile)
					return
				}
			}
		}
	}
}
