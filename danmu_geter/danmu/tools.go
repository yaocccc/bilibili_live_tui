package danmu

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"reflect"
	"unsafe"

	"github.com/gorilla/websocket"
)

func zlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

func ByteArrToDecimal(src []byte) (sum int) {
	if src == nil {
		return 0
	}
	b := []byte(hex.EncodeToString(src))
	l := len(b)
	for i := l - 1; i >= 0; i-- {
		base := int(math.Pow(16, float64(l-i-1)))
		var mul int
		if int(b[i]) >= 97 {
			mul = int(b[i]) - 87
		} else {
			mul = int(b[i]) - 48
		}

		sum += base * mul
	}
	return
}

func (d *DanmuClient) sendPackage(packetlen uint32, magic uint16, ver uint16, typeID uint32, param uint32, data []byte) (err error) {
	packetHead := new(bytes.Buffer)

	if packetlen == 0 {
		packetlen = uint32(len(data) + 16)
	}
	var pdata = []interface{}{
		packetlen,
		magic,
		ver,
		typeID,
		param,
	}

	// 将包的头部信息以大端序方式写入字节数组
	for _, v := range pdata {
		if err = binary.Write(packetHead, binary.BigEndian, v); err != nil {
			fmt.Println("binary.Write err: ", err)
			return
		}
	}

	// 将包内数据部分追加到数据包内
	sendData := append(packetHead.Bytes(), data...)

	// fmt.Println("本次发包消息为：", sendData)

	if err = d.conn.WriteMessage(websocket.BinaryMessage, sendData); err != nil {
		fmt.Println("conn.Write err: ", err)
		return
	}

	return
}

func splitMsg(src []byte) (msgs [][]byte) {
	lens := ByteArrToDecimal(src[:4])
	totalLen := len(src)
	startLoc := 0
	for {
		if startLoc+lens <= totalLen {
			msgs = append(msgs, src[startLoc:startLoc+lens])
			startLoc += lens
			if startLoc < totalLen {
				lens = ByteArrToDecimal(src[startLoc : startLoc+4])
			} else {
				break
			}
		} else {
			break
		}
	}
	return msgs
}

func BytesToStringFast(b []byte) string {
    return *(*string)(unsafe.Pointer(&b))
}

func StringToBytes(s string) []byte {
    sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
    bh := reflect.SliceHeader{sh.Data, sh.Len, 0}
    return *(*[]byte)(unsafe.Pointer(&bh))
}
