package websocket

import (
	"strings"
	"strconv"
	"errors"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
)

//handshake get
/*GET /chat HTTP/1.1
Host: server.example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Origin: http://example.com
Sec-WebSocket-Protocol: chat, superchat
Sec-WebSocket-Version: 13*/

//handshake send
/*HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
Sec-WebSocket-Protocol: chat*/

//key guid
var GUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

//authorize header and split sec-websocket-key
func AuthorizeHeaders(header string) (string, error){
	subs := strings.Split(header, "\r\n")
	for _, attr := range subs{
		switch{
		case strings.Contains(attr, "GET"):
			value := strings.Split(attr, "HTTP/")
			if v,err := strconv.ParseFloat(value[len(value)-1],64);err==nil{
				if v < 1.1{
					return "", errors.New("GET incorrect")
				}
			}
		case strings.Contains(attr, "Upgrader"):
			value := strings.Split(attr, ": ")
			if value[len(value)-1] != "websocket"{
				return "", errors.New("Upgrader incorrect")
			}
		case strings.Contains(attr, "Connection"):
			value := strings.Split(attr, ": ")
			if value[len(value)-1] != "Upgrade"{
				return "", errors.New("Connection incorrect")
			}
		case strings.Contains(attr, "Sec-WebSocket-Version"):
			value := strings.Split(attr, ": ")
			if value[len(value)-1] != "13"{
				return "", errors.New("Sec-WebSocket-Version incorrect")
			}
		case strings.Contains(attr, "Sec-WebSocket-Protocol"):
			value := strings.Split(attr, ": ")
			if !strings.Contains(value[len(value)-1],"chat"){
				return "", errors.New("Sec-WebSocket-Protocol incorrect")
			}
		case strings.Contains(attr, "Origin"):
			value := strings.Split(attr, ": ")
			if len(value) < 2 {
				return "", errors.New("Origin incorrect")
			}
		case strings.Contains(attr, "Sec-WebSocket-Key"):
			value := strings.Split(attr, ": ")
			if len(value) != 2{
				return "", errors.New(value[0] + " no value")
			} 
			return value[1], nil
		}
	}
	return "", errors.New("header incorrect")
}

//add "258EAFA5-E914-47DA-95CA-C5AB0DC85B11", SHA1, base64
func PrepareKey(key string) string{
	p := make([]byte, len(key)+len(GUID))
	copy(p[:24], key)
	copy(p[24:],GUID)
	sum := sha1.Sum(p)
	return base64.StdEncoding.EncodeToString(sum[:])
}

//handle http headers
func PrepareHeaders(key string) []byte{
	header := "HTTP/1.1 101 Switching Protocols\r\n" + "Upgrade: websocket\r\n" + "Connection: Upgrade\r\n" + "Sec-WebSocket-Accept: " + key + "\r\nSec-WebSocket-Protocol: chat\r\n\r\n"
	return []byte(header)
}

func Int64ToBytes(i int64) []byte {
    var buf = make([]byte, 8)
    binary.BigEndian.PutUint64(buf, uint64(i))
    return buf
}

func BytesToInt64(buf []byte) int64 {
    return int64(binary.BigEndian.Uint64(buf))
}

