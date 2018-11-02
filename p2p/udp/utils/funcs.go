package utils

import (
	"fmt"
	"os"
	"net"
	"crypto/md5"
)

func GetLocalIP()string{

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var returnIP string
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//fmt.Println(ipnet.IP.String())
				returnIP = ipnet.IP.String()
				break
			}

		}
	}
	return returnIP
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//创建目录
func MakeDir(DirPath string)(bool,error){
	exist, err := PathExists(DirPath)
	if exist {
		return exist,err
	} else {
		// 创建文件夹
		err := os.Mkdir(DirPath, os.ModePerm)
		if err != nil {
			return false,err
		} else {
			return true,err
		}
	}
}

//2次md5 简单加密
func HashEncryptFunc(s []byte)string{

	has := md5.Sum(s)
	hastwo :=md5.Sum(has[:])
	md5str := fmt.Sprintf("%x", hastwo) //将[]byte转成16进制
	return md5str
}