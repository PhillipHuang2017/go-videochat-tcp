package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"gocv.io/x/gocv"
	"net"
	"time"
)

var (
	ADDR string = "127.0.0.1:8080"
)

func dataDecode(reader *bufio.Reader) ([]byte, error) {
	lengthByte, err := reader.Peek(4)
	if err != nil {
		return nil, err
	}
	_, _ = reader.Discard(4)
	var lenData int32
	_ = binary.Read(bytes.NewBuffer(lengthByte), binary.BigEndian, &lenData)
	/////////////
	//fmt.Printf("lenData:%d\n", lenData)
	/////////////
	imgBuf := bytes.NewBuffer([]byte{})
	needRcv := lenData
	bufSize := int32(4096)
	byteBuf := make([]byte, bufSize)
	for needRcv > 0 {
		if needRcv < bufSize {
			byteBuf = byteBuf[:needRcv]
		}
		n, _ := reader.Read(byteBuf)
		imgBuf.Write(byteBuf[:n])
		needRcv = needRcv - int32(n)
	}
	return imgBuf.Bytes(), nil
}

func rcvHandle(conn net.Conn) {
	window := gocv.NewWindow("receiveVideo")
	defer window.Close()
	reader := bufio.NewReader(conn)
	preTime := time.Now().UnixNano()
	num := 0
	for {
		dataByte, err := dataDecode(reader)
		if err != nil {
			continue
		}
		frame, _ := gocv.IMDecode(dataByte, gocv.IMReadUnchanged)
		window.IMShow(frame)
		key := window.WaitKey(1)
		if key == int('q') {
			break
		}
		num++
		if num%10 == 0 {
			nowTime := time.Now().UnixNano()
			frameRate := float64(num) / (float64(nowTime-preTime) * 1e-9)
			preTime = nowTime
			num = 0
			fmt.Printf("frameRate:%.3f\n", frameRate)
		}
	}
}

func main() {
	conn, err := net.Dial("tcp", ADDR)
	if err != nil {
		fmt.Printf("Dial failed, err:%v\n", err)
		return
	}
	rcvHandle(conn)
}
