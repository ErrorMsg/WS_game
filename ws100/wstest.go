package main

import (
	"log"
	"../ws"
	"net"
	"time"
	"encoding/json"
)


type Player struct{
	Name string
	Position
	Hands []int
}

type Position struct{
	Level int
	Block int
}

func msgtest(ws *websocket.WS){
	time.Sleep(time.Second*10)
	msg := []byte("msgtest")
	err := ws.SendFrame([][]byte{msg})
	if err != nil{
		log.Println(err)
	}
}

func Dealmsg(data []byte) bool{
	switch string(data[0]){
	case "s":
		var player Player
		err := json.Unmarshal(data[1:],&player)
		if err != nil{
			log.Println(err)
			return false
		}else{
			log.Println(player)
		}
	default:
		return false
	}
	return false
}

func main(){
	ln, err := net.Listen("tcp","127.0.0.1:9001")
	if err != nil{
		log.Println("listen failed: ", err)
	}
	log.Println("listen")
	authConn := websocket.NewAuthConns()
	
	for {
		conn, err := ln.Accept()
		if err != nil{
			log.Println("conn failed once: ", err)
			continue
		}
		if c,ok:=authConn.Conns[conn];!ok || !c{
			ws := websocket.NewWS(conn)
			go ws.HandleConn(conn, authConn, Dealmsg)
		}
	}
}