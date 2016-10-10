//通讯协议处理，主要处理封包和解包的过程
package socket

import (
	"bytes"
	"encoding/binary"
)

const (
	constHeader       = "Header"
	constHeaderLength = 6
	constDataLength   = 4
)

//封包
func packet(message []byte) []byte {
	return append(append([]byte(constHeader), IntToBytes(len(message))...), message...)
}

//解包
func unpack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+constHeaderLength+constDataLength {
			break
		}
		if string(buffer[i:i+constHeaderLength]) == constHeader {
			messageLength := BytesToInt(buffer[i+constHeaderLength : i+constHeaderLength+constDataLength])
			if length < i+constHeaderLength+constDataLength+messageLength {
				break
			}
			data := buffer[i+constHeaderLength+constDataLength : i+constHeaderLength+constDataLength+messageLength]
			readerChannel <- data
			i += constHeaderLength + constDataLength + messageLength - 1
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
