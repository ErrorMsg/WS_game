package websocket

import (
	"encoding/binary"
	"log"
)

const (
	opcontinuation byte = 0x0
	optext byte = 0x1
	opbinary byte = 0x2
	opclose byte = 0x8
	opping byte = 0x9
	oppong byte = 0xa
)

type Frame struct{
	fin bool//1 bit, 0-continue, 1-end
	rsv byte//RSV123, 1 bit each, must 000 or specific
	opcode byte//4 bits, x0-xF 0-continue 1-text 2-binary 8-close 9-ping A-pong
	masked bool//1 bit, 0-not, 1-masked
	length int64//7 bits, 7+16 bits, 7+64 bits, 0-125=normal, 126=length is next 16 bits, 127=length is next 64 bits (length = extension data length + applicaiton data length)
	//extendlength //16 bits or 64 bits
	maskingkey []byte//0 byte or 4 bytes
	//data //x+y bytes extension data + application data
	extension []byte//0 byte or x bytes, normally 0 if no specific by Sec-WebSocket-Extensions
	application []byte//y bytes
}

//read frame from data
func ReadFrame(data []byte) (Frame,error){
	frame := Frame{}
	firstF := data[0]
	//if ping or pong, return directly
	/*if firstF == 0x89{
		frame.opcode = opping
		//return frame
	}else if firstF == 0x8a{
		frame.opcode = oppong
		return frame, nil
	}*/
	//get fin, 1=end, 0=continue
	if firstF >> 7 == 0x1{
		frame.fin = true
	}else{
		frame.fin = false
	}
	//get rsv
	frame.rsv = (firstF >> 4) & 0x7
	//get opcode
	frame.opcode = firstF & 0x0f
	secondF := data[1]
	//get if masked, 1=true, 0=false
	if secondF >> 7 == 0x1{
		frame.masked = true
	}else{
		frame.masked = false
		frame.opcode = opclose
		return frame, nil
	}
	//get normal length
	length := secondF & 0x7f
	switch length{
	//if length=126, get length from next 2 bytes
	case 126:
		frame.length = int64(binary.BigEndian.Uint16(data[2:4])) //int64(binary.BigEndian.Uint16(data[2:4]))
		if frame.masked{
			frame.maskingkey = data[4:8]
			frame.application = data[8:frame.length+8]
		}else{
			frame.application = data[4:frame.length+4]
		}
	//if length=127, get length from next 8 bytes
	case 127:
		frame.length = int64(binary.BigEndian.Uint16(data[2:10]))
		if frame.masked{
			frame.maskingkey = data[10:14]
			frame.application = data[14:frame.length+14]
		}else{
			frame.application = data[10:frame.length+10]
		}
	default:
		frame.length = int64(length)
		if frame.masked{
			frame.maskingkey = data[2:6]
			frame.application = data[6:frame.length+6]
		}else{
			frame.application = data[2:frame.length+2]
		}
	}
	return frame, nil
}

//received = frame.application
//dealmsg is a function, func([]byte)bool, if deal with received data success, return true, else return false, use it to unmarshal json
func HandleReceived(received []byte, maskingkey []byte, dealmsg func([]byte)bool) []byte{
	if maskingkey != nil{
		for i:=0;i<len(received);i++{
			received[i] = received[i] ^ maskingkey[i%4]
		}
	}
	
	//rData := string(received)
	if ok:=dealmsg(received);!ok{
		log.Println("receive msg: ", string(received))
	}
	return received
	/*switch msgType := received[0];msgType{
	case 0x00: //normal message, broadcast it
		ws.Broadcast(received[1:])
	case 0x01: //friend message, transfer it
		ws.PrivateChat(received[1:])
	case 0x02: //in-game message, broadcast it in game
		ws.PublicChat(received[1:])
	case 0x10: //painting
		ws.Paint(received[1:])
	}*/
}

/*

*/