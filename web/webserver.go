package web

import (
	"p2pudp/web/Base"
	"net/http"
	"log"
	"net"
	"p2pudp/p2p/udp/utils"
	"fmt"
	"strconv"
	"bytes"
	"p2pudp/client"
	"p2pudp/p2p/udp/pack"
)

var mux *Base.MyMux
var HostAddrs *utils.Set
var clientConn *net.UDPConn

func init(){
	mux = &Base.MyMux{}
	mux.MyMuxInit()
	mux.Routers["/"] = getList
	mux.Routers["/send"] = sendMsg
	mux.Routers["/file"] = sendFile
}

func SetClient(udpconn *net.UDPConn){
	clientConn = client.ClientUDPconn
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
//显示信息
func showMsg(msg string){
	log.Println(msg)
}
//建立web 服务端
func BuildServerWeb(port string){
	showMsg("开始打开"+port+"端口")
	err := http.ListenAndServe(":"+port,mux)
	logErr(err,"Println")

}
//查询列表
func getList(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"节点列表\n\r")
	fmt.Fprintf(w,"==========================================\n\r")

	for i, v := range client.GetHosts().Elements() {
		if v != nil{
			fmt.Fprintf(w, "id:" + strconv.Itoa(i) + " :    addr: "+  v.(string)+ "\n\r")
		}
	}
}
//发送消息
func sendMsg(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	ipaddrs, _ := r.Form [ "ip" ]
	ipaddr := ipaddrs[0]
	msg, _ := r.Form [ "msg" ]
	content := msg[0];
	action := r.Form [ "action" ][0]
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(content)

	pk := pack.Pack{pack.Msg, buf.Bytes()}

	if(action == "single"){
		tAddr, terr := net.ResolveUDPAddr("udp", ipaddr)
		if terr != nil{
			fmt.Fprintf(w,"single",terr)
		}
		_, err := clientConn.WriteToUDP(pk.Encode(),tAddr)
		if err == nil {
			fmt.Fprintf(w, "消息发送成功\n\r")
		}else{
			fmt.Fprintf(w,"single",err)
		}
	}else{

		for _, v := range client.GetHosts().Elements() {
			sAddr, sErr := net.ResolveUDPAddr("udp", v.(string))
			if sErr != nil{
				logErr(sErr,"Println")
			}
			_, err2 := clientConn.WriteToUDP(pk.Encode(), sAddr)
			if err2 == nil{
				fmt.Fprintf(w, "内容: "+content+" 发送成功\n\r")
			}else{
				logErr(err2,"Println")
			}
		}
	}
}

//发送文件
func sendFile(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	ipaddrs, _ := r.Form [ "ip" ]
	ip := ipaddrs[0]
	files, _ := r.Form [ "filename" ]
	filename := files[0]
	client.SendFileFuncs(filename,ip)
	fmt.Fprintln(w,"发送文件已完成,请查收")
}