package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gocv.io/x/gocv"
	"net"
)

var (
	ADDR   string = "127.0.0.1:8080"
	CAM_ID int    = 0
)

func dataEncode(data []byte) []byte {
	lenData := int32(len(data)) // 占4个字节
	bytesBuffer := bytes.NewBuffer([]byte{})
	// 以大端模式写入，即高字节在高地址，低字节在低地址（也就是和正常理解的数字存放方式一样）
	_ = binary.Write(bytesBuffer, binary.BigEndian, lenData)
	_ = binary.Write(bytesBuffer, binary.BigEndian, data)
	fmt.Printf("lenData:%d\n", lenData)
	return bytesBuffer.Bytes()
}

func sendHandle(conn net.Conn) {
	defer conn.Close()
	cam, err := gocv.VideoCaptureDevice(CAM_ID)
	if err != nil {
		fmt.Printf("Open cammera failed, err:%v\n", err)
		return
	}
	defer cam.Close()
	frame := gocv.NewMat()
	/////////////////
	num := 0
	//////////////////
	for {
		ok := cam.Read(&frame)
		if !ok {
			break
		}
		byteFrame, _ := gocv.IMEncode(gocv.JPEGFileExt, frame)
		data := dataEncode(byteFrame)
		_, err := conn.Write(data)
		if err != nil {
			fmt.Printf("Connection closed, err:%v\n", err)
			break
		}
		//////////////////
		num++
		fmt.Printf("send frame nums:%d\n", num)
		//////////
	}
}

func main() {

	listen, err := net.Listen("tcp", ADDR)
	if err != nil {
		fmt.Printf("Listen failed, err:%v\n", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("Accept failed, err:%v\n", err)
			return
		}
		go sendHandle(conn)
	}
}
