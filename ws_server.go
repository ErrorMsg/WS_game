package main

import (
	"log"
	"net"
	"sync"
	"io"
	"./ws"
)

var mu sync.Mutex

//handle first http connect, and convert connect to websocket
func HandleConn(conn net.Conn, authConn *websocket.AuthConns){
	defer conn.Close()
	ws := websocket.WS{
		C:conn,
		Closed:make(chan struct{}),
	}
	buf := make([]byte, 1024)
	bufChan := make(chan []byte, 10)
	defer close(bufChan)
	go websocket.HandleWS(bufChan, ws, authConn)
	for{
		_, err := conn.Read(buf)
		if err != nil{
			if err != io.EOF{ //if client ws closed abnormally, close ws
				log.Println("conn read failed: ", err)
				//ws.Close()
				mu.Lock()
				authConn.Conns[conn] = false
				mu.Unlock()
				return
			}
			continue
		}
		//if conn is in authConn, start websocket handle
		if authConn.Conns[conn]{
			//websocket.HandleWS(buf[:],ws)
			bufChan <- buf[:]
			select{ //if websocket closed, delete conn from authConn
			case <- ws.Closed:
				//delete(authConn, conn)
				mu.Lock()
				authConn.Conns[conn] = false
				mu.Unlock()
				return
			default:
			}
			continue
		}

		//get sec-websocket-key
		key,err := websocket.AuthorizeHeaders(string(buf))
		if err != nil{
			log.Println("auth failed: ", err)
			return
		}
		nkey := websocket.PrepareKey(key)
		nheader := websocket.PrepareHeaders(nkey)
		_, err = conn.Write(nheader)
		if err != nil{
			log.Println("new header send failed: ", err)
		}
	
		//if handshake success, add conn into authConn
		mu.Lock()
		authConn.Conns[conn] = true
		mu.Unlock()
		log.Println("new websocket from ",conn.RemoteAddr().String())
		err = ws.SendWelcome()
		if err != nil{
			log.Println("welcome failed: ", err)
		}
	}
}

func main(){
	ln, err := net.Listen("tcp","127.0.0.1:9001")
	if err != nil{
		log.Println("listen failed: ", err)
	}
	log.Println("listen")
	//wss := make(map[int]*WS)
	//authConn := make(map[net.Conn]bool)
	authConn := websocket.NewAuthConns()
	go authConn.HandlePubConns()
	for {
		conn, err := ln.Accept()
		if err != nil{
			log.Println("conn failed once: ", err)
			continue
		}
		if c,ok:=authConn.Conns[conn];!ok || !c{
			go HandleConn(conn, authConn)
		}
	}
}