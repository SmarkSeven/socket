//通讯协议处理，主要处理封包和解包的过程
package protocol

import (
	"bytes"
	"encoding/binary"
)

const (
	ConstHeader       = "Header"
	ConstHeaderLength = 6
	ConstDataLength   = 4
)

//封包
func Packet(message []byte) []byte {
	return append(append([]byte(ConstHeader), IntToBytes(len(message))...), message...)
}

//解包
func Unpack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+ConstHeaderLength+ConstDataLength {
			break
		}
		if string(buffer[i:i+ConstHeaderLength]) == ConstHeader {
			messageLength := BytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstDataLength])
			if length < i+ConstHeaderLength+ConstDataLength+messageLength {
				break
			}
			data := buffer[i+ConstHeaderLength+ConstDataLength : i+ConstHeaderLength+ConstDataLength+messageLength]
			readerChannel <- data
			i += ConstHeaderLength + ConstDataLength + messageLength - 1
		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
