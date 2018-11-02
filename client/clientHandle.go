package client

import (
	"p2pudp/p2p/udp/utils"
	"p2pudp/p2p/udp/pack"
	"net"
	"log"
	"errors"
	"strings"
	"os"
	"io"
	"path"
	"strconv"
)

//保存所有的客户端连接
var Hosts *utils.Set
var ClientUDPconn *net.UDPConn

type ClientHandle struct {
	ListenPort string
	ServerIP string
	PublicIP string
}
func init(){
	Hosts = utils.NewSet()
}

func GetHosts() *utils.Set {
	return Hosts
}

//制作监听
func (this *ClientHandle) MakUDPListener()  {
	ClientUDPListenerAddress, ClientUDPListenererr := net.ResolveUDPAddr("udp", "0.0.0.0:"+this.ListenPort)
	logErr(ClientUDPListenererr,"Println")
	var UDPerr error
	ClientUDPconn, UDPerr = net.ListenUDP("udp", ClientUDPListenerAddress)
	logErr(UDPerr,"Println")
	showMsg("UDP客户端启动成功，绑定端口:"+this.ListenPort)
	this.ConnectServer()
	showMsg("向服务端发送请求.")
}

//关闭UDP
func (this *ClientHandle) CloseClientUDPconn(){
	if ClientUDPconn != nil{
		err:=ClientUDPconn.Close()
		logErr(err,"Println")
	} else{
		logErr(errors.New("关闭udpConn失败，可能udpConn还是空对象"),"Println")
	}
}

//打印错误日志
func logErr(err error,flag string){
	if err != nil {
		switch flag {
		case "Println":
			log.Println(err.Error())
			break;
		case "Fatalln":
			log.Fatalln(err.Error())
			break;
		case "Panicln":
			log.Panicln(err.Error())
			break;
		}
	}
}

//连接服务器
func (this *ClientHandle)ConnectServer(){

	ip:=utils.GetLocalIP()
	pk := pack.Pack{pack.Join, []byte(ip+":"+this.ListenPort)}
	bs := pk.Encode()

	RemoteUdpServerAddr, RemoteUdperr := net.ResolveUDPAddr("udp", this.ServerIP)
	logErr(RemoteUdperr,"Println")
	_, udperr := ClientUDPconn.WriteToUDP(bs, RemoteUdpServerAddr)  //udp 请求服务器
	logErr(udperr,"Println")
	showMsg("UDP新人向服务器申请加入")

}

//显示信息
func showMsg(msg string){
	log.Println(msg)
}

//客户端UPD接受数据
func (this *ClientHandle)ClientUdpReciveData(){

	for {
		//time.Sleep(time.Second)
		buf := make([]byte, 400)
		n, RemoteAddr, err := ClientUDPconn.ReadFromUDP(buf[:])
		logErr(err,"Println")
		p := pack.Pack{}
		//fmt.Println("收到来自", RemoteAddr, "的数据"+string(buf))
		if !p.Decode(buf[:n]) {
			showMsg("错误:"+ RemoteAddr.String()+ "的无效的数据包")
			continue
		}

		switch p.Type {
		case pack.Join:
			this.joinFuncs(p)
			break;
		case pack.JoinReturn:
			this.joinReturnFuncs(p)
			break;
		case pack.Msg:
			this.msgFuncs(RemoteAddr,p)
			break;
		case pack.Hole:
			this.holeFuncs(RemoteAddr,p)
			break;
		case pack.ResponesIP:
			this.responesIPFuncs(p)
			break;
		case pack.GetFile:
			this.GetFile(p)
			break;
		default:
			showMsg("default method is run.")
			break;
		}
	}
}

//有新人加入后处理方法
func (this *ClientHandle) joinFuncs(p pack.Pack){

	//确认与新人的IP
	var sendIP string
	a_r := strings.Split(this.PublicIP,":")
	other_addrs := string(p.Data[:])
	b_r := strings.Split(other_addrs,"=")
	c_r := strings.Split(b_r[0],":")
	if c_r[0] == a_r[0]{
		sendIP = b_r[1]
	} else {
		sendIP = b_r[0]
	}

	addr, err := net.ResolveUDPAddr("udp", sendIP)
	logErr(err,"Println")
	Hosts.Add(sendIP,sendIP)
	pk := pack.Pack{pack.Hole, []byte("welcome to join us")}
	bs := pk.Encode()
	_, err = ClientUDPconn.WriteToUDP(bs, addr)
	logErr(err,"Println")
	showMsg("向UDP新人"+sendIP+"打洞")
}

//本人加入后收到服务器反馈处理方法
func (this *ClientHandle) joinReturnFuncs(p pack.Pack){

	HostsList := utils.NewHostSet()
	if !HostsList.Decode(p.Data) {
		showMsg("向服务器申请加入后收到反馈消息为:"+string(p.Data))
		return
	}
	Hosts.SetInData(HostsList)

	for _, host := range Hosts.Elements() {
		var sendtoIP string
		sendlocalip :=Hosts.GetLocalAddr(host.(string))
		a_r := strings.Split(this.PublicIP,":")
		b_r := strings.Split(host.(string),":")

		//我自己的ip 跳过
		if this.PublicIP == host.(string){
			continue
		}

		if b_r[0] == a_r[0]{
			sendtoIP = sendlocalip
		} else {
			sendtoIP = host.(string)
		}
		Hosts.Add(sendtoIP,sendtoIP)
		addr, err := net.ResolveUDPAddr("udp", sendtoIP)
		if err != nil {
			logErr(err,"Println")
			continue
		}
		pk := pack.Pack{pack.Hole, []byte("hello I joined in")}
		bs := pk.Encode()
		_, err = ClientUDPconn.WriteToUDP(bs, addr)
		logErr(err,"Println")
		showMsg("本加入后向老用户"+sendtoIP+"打洞")
	}
}

//关闭UDP链接
func (this *ClientHandle) CloseUDPConn(){
	if ClientUDPconn != nil{
		err:=ClientUDPconn.Close()
		logErr(err,"Println")
	} else{
		logErr(errors.New("关闭udpConn失败，可能udpConn还是空对象"),"Println")
	}
}

//消息方法
func (this *ClientHandle) msgFuncs(RemoteAddr *net.UDPAddr,p pack.Pack){
	showMsg(RemoteAddr.String()+"UDP说:"+string(p.Data[:]))
}

//其他客户端打洞方法
func (this *ClientHandle) holeFuncs(RemoteAddr *net.UDPAddr,p pack.Pack){
	showMsg(RemoteAddr.String()+"向我打洞说:"+string(p.Data[:]))
}

//服务器返回本机请求的外网IP及端口
func (this *ClientHandle) responesIPFuncs(p pack.Pack){
	this.PublicIP =string(p.Data[:])
	showMsg("服务器ResponseIP:"+string(p.Data[:]))
}

//发送文件
func SendFileFuncs(filename string,ip string){
	filepath := "./source/"+filename
	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		logErr(err,"Println")
		defer file.Close()
		//os.Exit(0)
	}else {

		RemoteUdpServerAddr, RemoteUdperr := net.ResolveUDPAddr("udp", ip)
		logErr(RemoteUdperr,"Println")
		fileSuffix := path.Ext(filename) //获取文件后缀
		filenameOnly := strings.TrimSuffix(filename, fileSuffix)//获取文件名
		id:=1
		//fileNameHash := utils.HashEncryptFunc([]byte(filename))

		for {
			buf := make([]byte, 140)
			n, ferr := file.Read(buf)
			if ferr != nil && ferr != io.EOF{
				logErr(ferr,"Println")
				break
			}
			if n == 0{
				break
			}
			fileStruct := pack.File{Id:id,FileFullname:filename,Content:buf[:n],FileShortName:filenameOnly,FileExt:fileSuffix,FilePath:filepath,FileHash:"",FileNameHash:""}
			pk := pack.Pack{pack.GetFile, fileStruct.Encode()}
			bs := pk.Encode()
			_, udperr := ClientUDPconn.WriteToUDP(bs, RemoteUdpServerAddr)  //udp 请求服务器
			logErr(udperr,"Println")
			showMsg("SendFileFuncs 发送文件已经调用,数据包:"+strconv.Itoa(fileStruct.Id)+"。")
			id++
		}
	}
}

//获取文件
func (this *ClientHandle) GetFile(p pack.Pack) {

	fileStruct := pack.File{}

	if err:=fileStruct.Decode(p.Data);err != nil{
		logErr(err,"Println")
	}else{

		//utils.MakeDir()
		//分片 后补方法

		filepath := "./source/trans/"+fileStruct.FileFullname
		f, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
		if err != nil {
			logErr(err,"Println")
		}
		n, err := f.Write(fileStruct.Content)
		if err == nil && n < len(fileStruct.Content) {
			err = io.ErrShortWrite
			logErr(err,"Println")
		}
		if err1 := f.Close(); err == nil {
			logErr(err1,"Println")
		}
		showMsg("GetFile has recive file "+strconv.Itoa(fileStruct.Id)+" data.")
	}

}
