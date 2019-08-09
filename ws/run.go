package websocket

import (
	"log"
	"net"
	"sync"
	"io"
)

var mu sync.Mutex

//handle first http connect, and convert connect to websocket
func (ws *WS)HandleConn(conn net.Conn, authConn *AuthConns, dealmsg func([]byte)bool){
	defer conn.Close()
	buf := make([]byte, 1024)
	bufChan := make(chan []byte, 10)
	defer close(bufChan)
	go ws.HandleWS(bufChan, authConn, dealmsg)
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
		key,err := AuthorizeHeaders(string(buf))
		if err != nil{
			log.Println("auth failed: ", err)
			return
		}
		nkey := PrepareKey(key)
		nheader := PrepareHeaders(nkey)
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
