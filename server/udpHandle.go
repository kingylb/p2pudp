package server

import (
	"p2pudp/p2p/udp/utils"
	"net"
	"log"
	"errors"
	"p2pudp/p2p/udp/pack"
	funs "p2pudp/p2p/udp/utils"
	"strconv"
)

//保存所有的客户端连接
var hosts *utils.Set
var hosts_set *utils.HostSet
var udpConn *net.UDPConn

type UDPHandle struct {
	ListenPort string
}


func init(){
	hosts = utils.NewSet()
	hosts_set = utils.NewHostSet()
}

//返回hosts
func (this *UDPHandle) getHostAddrs() *utils.Set{
	return hosts
}

//制作监听
func (this *UDPHandle) MakUDPListener()  {
	UDPListenerAddress, UDPListenererr := net.ResolveUDPAddr("udp", "0.0.0.0:"+this.ListenPort)
	this.logErr(UDPListenererr,"Println")
	var UDPerr error
	udpConn, UDPerr = net.ListenUDP("udp", UDPListenerAddress)
	this.logErr(UDPerr,"Println")
	this.showMsg("UDP服务端启动成功，绑定端口:"+this.ListenPort)
}

func(this *UDPHandle) ReciveData(){

	go func(){

			for {
				//time.Sleep(time.Second)
				var buf [400]byte
				n, RemoteAddr, err := udpConn.ReadFromUDP(buf[:])
				if err != nil {
					this.logErr(err,"Println")
					continue
				}
				p := pack.Pack{}

				if !p.Decode(buf[:n]) {
					this.showMsg("收到来自"+RemoteAddr.String() + "的无效的数据包."+ err.Error())
					continue
				}
				RemoteOutIP :=RemoteAddr.String()
				//来源IP 是本机的情况转换ip(127.0.0.1)为内网IP
				if RemoteAddr.IP.String() == "127.0.0.1"{
					RemoteOutIP = funs.GetLocalIP()+":"+ strconv.Itoa(RemoteAddr.Port)
				}

				//判断为登陆状态
				if p.Type == pack.Join {

					//已有的ip重新登陆要先删除以前的登陆记录
					if hosts.Contains(RemoteOutIP) == true{
						hosts.Remove(RemoteOutIP)
					}
					//有新用户加入，向其他所有用户发送加入请求包，请求包包含了新用户ip：port信息，用于所有用户对他进行打洞
					if hosts.Add(RemoteOutIP,string(p.Data)) {
						this.showMsg("新用户"+ RemoteOutIP+" 报道.")

						//遍历主机列表 对自己及其它客户端回复信息
						for _, v := range hosts.Elements() {
							//向新用记回复收到消息
							if v.(string) == RemoteOutIP {
								tAddr, selferr := net.ResolveUDPAddr("udp", v.(string))
								this.logErr(selferr,"Pringln")
								//返回请求机的外网IP及端口
								selfpk := pack.Pack{pack.ResponesIP, []byte(RemoteOutIP)}
								selfbs := selfpk.Encode()
								//对请求者发送消息
								_, writeerr := udpConn.WriteToUDP(selfbs, tAddr)
								this.logErr(writeerr,"Println")
								continue
							}
							//组合请求者的外网IP及协带的内网IP 发送到其它客户端
							ap :=RemoteOutIP + "="+string(p.Data)
							tAddr, err := net.ResolveUDPAddr("udp", v.(string))
							pk := pack.Pack{pack.Join, []byte(ap)}
							bs := pk.Encode()
							//发送其它客户端
							_, err = udpConn.WriteToUDP(bs, tAddr)
							if err != nil {
								this.logErr(err,"Println")
							}
						}

						//新用户添加成功之后，把当前的所有用户信息都告诉【新用户】，用于一一打洞
						hosts_set.SetData(hosts)
						pk := pack.Pack{pack.JoinReturn, hosts_set.Encode()}
						bs := pk.Encode()
						udpConn.WriteToUDP(bs, RemoteAddr)
					}
				}
			}
	}()
}

//关闭UDP链接
func (this *UDPHandle) CloseUDPConn(){
	if udpConn != nil{
		err:=udpConn.Close()
		this.logErr(err,"Println")
	} else{
		this.logErr(errors.New("关闭udpConn失败，可能udpConn还是空对象"),"Println")
	}
}

//打印错误日志
func (this *UDPHandle) logErr(err error,flag string){
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
func (this *UDPHandle)showMsg(msg string){
	log.Println(msg)
}
