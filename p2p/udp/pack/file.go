package pack

import (
	"bytes"
	"encoding/gob"
)

type File struct {
	Id int
	FileFullname string
	Content []byte
	FileShortName string
	FileExt string
	FilePath string
	FileHash string
	FileNameHash string
}

/**
序列化pack,返回序列化后的字节数组
**/
func (this *File) Encode() []byte {
	buf := bytes.NewBuffer([]byte{})
	ge := gob.NewEncoder(buf)
	err := ge.Encode(this)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()

}
/**
反序列化pack
**/
func (this *File) Decode(bs []byte) error {
	b := bytes.NewBuffer(bs)
	gd := gob.NewDecoder(b)
	err := gd.Decode(this)
	return err
}
