//used to deal with case from websocket

package srv

import (
	"net"
	"log"
)


var Done = make(chan struct{})

type MainHandler struct{
	interrupt chan os.Signal
	//complete chan error
	//timeout <-chan time.Time
	Handlers [string]func([]byte) error
	Srv SrvInfo
}


func creatMainHandler() *MainHandler{
	handlers := make([string]func([]byte) error)
	/* handlers["TEST"] = func(data []byte) error{
		log.Println(len(data))
		return nil
	} */
	srv, err := StartSrv(100)
	if err != nil{
		panic(err)
	}
	return &MainHandler{Handlers:handlers, Srv:srv}
}

func (mh *MainHandler)addMainHandler(hname, hfunc) error{
	if ok,_ := mh.Handlers[hname];!ok{
		mh.Handlers[hname] = hfunc
		return nil
	}else{
		return error.New(hname+" exists")
	}
}



func (mh *MainHandler)Run() error{
	for{
		select{
		case data := <- recvData:   //deal with data if received
			msgType := data[:4]
			go func(){
				mh.Handlers[msgType](data[4:])
			}()   //use worker factory ?
		case <- Done:   //if mainhandler closed, clean left received data
			for data := range recvData{
				msgType := data[:4]
				go func(){
					mh.Handlers[msgType](data[4:])
				}()
			}
			return nil
		default:
		}
	}
}


func (mh *MainHandler)Stop(){
	close(Done)
}