package websocket

import (
	"log"
	"net"
	"sync"
	"io"
	"bytes"
)

type PubMsg struct{
	Sender net.Conn
	Message []byte
}

type AuthConns struct{
	//value shows conn's status, true = connect/active; false = disconnect temporary/away; no value = offline/inactive
	Conns map[net.Conn]bool
	PubChan chan PubMsg
	mu sync.Mutex
}

//create new AuthConns, use map to manage conns
func NewAuthConns() *AuthConns{
	return &AuthConns{
		Conns: make(map[net.Conn]bool),
		PubChan: make(chan PubMsg, 10),
	}
}

//when conn login, set value true
func (cs *AuthConns)ConnIn(c net.Conn){
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.Conns[c] = true
}

//when conn disconnect, set value false before delete it
func (cs *AuthConns)ConnOut(c net.Conn){
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.Conns[c] = false
}

//if conn is confirmed offline(disconnect >= 3min?), delete it
func (cs *AuthConns)ConnOff(c net.Conn){
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if _,ok := cs.Conns[c]; ok{
		delete(cs.Conns, c)
	}
}

//used to handle public messages
func (cs *AuthConns)HandlePubConns(){
	for data := range cs.PubChan{
		for c,s := range cs.Conns{ 
		
			//if conn is active, broadcast public message
			if s && c != data.Sender{ 
				err := BroadcastFrame(c, data.Message)
				if err != nil{
					log.Printf("broadcast to %s failed: ", c.RemoteAddr().String())
					log.Println(err)
				}
			}
		}
	}
}

func BroadcastFrame(w io.Writer, data []byte) (err error){
	buf := bytes.NewBuffer([]byte{})
	
	err = buf.WriteByte(byte(0x81))
	if err != nil{
		return
	}
	err = SetupData(data, buf)
	if err != nil{
		return
	}
	err = SendWS(w, buf.Bytes())

	log.Println("broadcast message")
	return
}