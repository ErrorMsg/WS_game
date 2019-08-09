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

func NewWS(conn net.Conn) *WS{
	return &WS{
		C:conn,
		Closed:make(chan struct{}),
	}
}

func (ws WS) Close(){
	ws.C.Close()
	close(ws.Closed)
}

func (ws *WS)HandleWS(bufChan chan []byte, cs *AuthConns, dealmsg func([]byte)bool){
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
			rData := HandleReceived(frame.application, frame.maskingkey, dealmsg) 
			
			//if this is a public message, add sender and broadcast it, cs.PubChan can store 10 public message only
			//pData := append([]byte(ws.C.RemoteAddr().String()), rData...) 
			//cs.PubChan <- PubMsg{Sender:ws.C, Message:pData}
			
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