package websocket

import (
	"log"
	"net"
	"sync"
)

type WS struct{
	C net.Conn
	Closed chan struct{}
	Mu sync.Mutex
}

func (ws WS) Close(){
	ws.C.Close()
	close(ws.Closed)
}

func HandleWS(bufChan chan []byte, ws WS, cs *AuthConns){
	//get received frame from bufChan
	for buf := range bufChan{
		frame, err := ReadFrame(buf)
		if err != nil{
			log.Println("readframe failed: ", err)
		}
		switch frame.opcode{
		
		//if text frame
		case optext: 
			log.Println(frame.length)
			
			//get deocded data
			rData := HandleReceived(frame.application, frame.maskingkey) 
			
			//if this is a public message, add sender and broadcast it
			pData := append([]byte(ws.C.RemoteAddr().String()), rData...) 
			cs.PubChan <- PubMsg{Sender:ws.C, Message:pData}
			
			err := ws.SendFrame([][]byte{rData})
			if err != nil{
				log.Println("send message failed: ", err)
			}
			
		case opclose:
			ws.Close()
			
		//if it's PING frame, reply PONG, else if PONG, reply PING?.
		case opping:
			log.Println("receive ping")
			err := ws.SendPong()
			if err != nil{
				log.Println("send ping failed: ", err)
			}
		case oppong:
			log.Println("receive pong")
			err := ws.SendPong()
			if err != nil{
				log.Println("send pong failed: ", err)
			}
		default:
		}
	}
}