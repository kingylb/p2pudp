package pack

import (
	"bytes"
	"encoding/gob"
)

type Hosts struct {
	Hosts map[interface{}]string
}
func NewHosts() *Hosts {
	return &Hosts{Hosts: make(map[interface{}]string)}
}

//添加主机到列表
func (this *Hosts) Add(host string, localhost string) {
	this.Hosts[host] = localhost
}

//
func (this *Hosts) Encode() []byte {
	buf := bytes.NewBuffer([]byte{})
	ge := gob.NewEncoder(buf)
	err := ge.Encode(this)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (this *Hosts) Decode(bs []byte) bool {
	b := bytes.NewBuffer(bs)
	gd := gob.NewDecoder(b)

	err := gd.Decode(this)

	return err == nil
}

//迭代
func (this *Hosts) GetElements() []interface{} {

	initlen := len(this.Hosts)
	snaphot := make([]interface{}, initlen)

	actuallen := 0

	for k, _ := range this.Hosts {
		if actuallen < initlen {
			snaphot[actuallen] = k
		} else {
			snaphot = append(snaphot, k)
		}
		actuallen++
	}

	if actuallen < initlen {
		snaphot = snaphot[:actuallen]
	}

	return snaphot
}


