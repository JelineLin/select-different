package compute

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
)

//生成hash字符串
func ByteToHash(data []byte) (rsp string) {
	a := sha256.New() //256位 32B
	_, err := a.Write(data)
	if err != nil {
		fmt.Println(err.Error())
	}
	rsp = hex.EncodeToString(a.Sum(nil))

	fmt.Println("ai = ", rsp)
	return
}

//读取1级文件生成2级文件
func Pack(filePath string) {
	//创建2级文件文件夹
	newdir := filePath + "/" + "level2"
	if err := os.Mkdir(newdir, os.ModePerm); err != nil {
		fmt.Println("md dir err : ", err.Error())
		return
	}

	//读取文件夹
	dirFds, err := ioutil.ReadDir(filePath)
	if err != nil {
		fmt.Println("read dir err : ", err.Error())
		return
	}

	for k, v := range dirFds { //内存考察
		if v.IsDir() {
			continue
		}
		fmt.Println(v.Name())

		path := filePath + "/" + v.Name()
		fileData := make([]byte, 1024)
		fileFd, err := os.Open(path)
		if err != nil {
			fmt.Println("1 open file ", path, " err : ", err.Error())
			return
		}
		defer fileFd.Close()
		index := 0
	reRead:
		n, err := fileFd.Read(fileData[index:])
		index += n
		if index < 1024 {
			goto reRead
		}
		newString := ByteToHash(fileData)
		newfilename := fmt.Sprintf("%s/%d", newdir, k/100000) //1000个2级文件
		newfilefd, err := os.OpenFile(newfilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println("2 open file ", newfilename, " err : ", err.Error())
			return
		}
		defer newfilefd.Close()
		index = 0
	reWrite:
		n, err = newfilefd.Write([]byte(newString[index:])) //位移
		if err != nil {
			fmt.Println("write file ", newfilename, " err : ", err.Error())
			return
		}
		if n < 32 {
			index += n
			goto reWrite
		}

	}
	return

}

//读取2级文件hash写入内存
func LevelBtoBuff(filePath string) []string {
	dirpath := filePath + "/" + "level2"
	rsp := make([]string, 0)
	dirFds, err := ioutil.ReadDir(dirpath)
	if err != nil {
		fmt.Println("read dir err : ", err.Error())
		return nil
	}

	for _, v := range dirFds { //1000个
		path := dirpath + "/" + v.Name()
		fileData := make([]byte, 1024)
		fileFd, err := os.Open(path)
		if err != nil {
			fmt.Println("3 open file ", path, " err : ", err.Error())
			return nil
		}
		defer fileFd.Close()
		index := 0
	reRead:
		n, err := fileFd.Read(fileData[index:])
		if n < 3125 {
			index += n
			goto reRead
		}
		rsp = append(rsp, ByteToHash(fileData))

	}
	return rsp
}

//读取2级文件内容到内存
func LevelAtoBuff(filePath string, arrea int) []string {
	rsp := make([]string, 0) //100000个
	path := fmt.Sprintf("%s/%d", filePath+"/"+"level2", arrea)
	filefd, err := os.Open(path)
	if err != nil {
		fmt.Println("4 open file ", path, " err : ", err.Error())
		return nil
	}
	defer filefd.Close()

	index := 0
	buff := make([]byte, 32)
	for {

		_, err = filefd.ReadAt(buff, int64(32*index))
		if err != nil {
			fmt.Println("5 open file ", path, " err : ", err.Error())
			return nil
		}
		index++
		if index == 10000 {
			break
		}
		rsp = append(rsp, string(buff))
	}
	return rsp
}

