package demo

import "encoding/binary"

//用多少个字节来表示长度
const LengthBytes = 8

func EncodeMsg(data []byte) []byte {
	resp := make([]byte, len(data)+LengthBytes)
	l := len(data)
	//大端是高字节存放到内存的低地址
	//
	//小端是高字节存放到内存的高地址
	binary.BigEndian.PutUint64(resp[:LengthBytes], uint64(l))
	copy(resp[LengthBytes:], data)
	return resp
}
