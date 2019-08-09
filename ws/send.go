package websocket

import (
	"log"
	_ "net"
	"bytes"
	"math"
	_ "sync"
	_ "time"
	"io"
)

//var mu sync.Mutex

//real method to send any data
func SendWS(w io.Writer, data []byte) error{
	mu.Lock()
	//ws.C.SetWriteDeadline(time.Now().Add(100))
	_ ,err := w.Write(data)
	mu.Unlock()
	return err
}

//method desert, ignore
func (ws WS)SendMessage(op byte, data []byte, tail bool) error{
	var firstF byte
	switch op{
	case opping:
		err := ws.SendPong()
		return err
	case oppong:
		err := ws.SendPong()
		return err
	case optext:
		err := ws.SendFrame([][]byte{data})
		return err
	case opbinary:
		if tail{
			firstF = 0x82
		}else{
			firstF = 0x02
		}
	case opclose:
		firstF = 0x88
		err := ws.SendClose()
		return err
	default:
	}
	_ = firstF
	return nil
}

func (ws WS)SendPing() error{
	//frame.opcode = oppong
	//err := ws.SendFrame(frame)
	buf := bytes.NewBuffer([]byte{})
	buf.WriteByte(byte(0x89))
	buf.WriteByte(byte(0x04))
	buf.Write([]byte("PING"))
	err := SendWS(ws.C, buf.Bytes())
	log.Println("send ping")
	return err
}

func (ws WS)SendPong() error{
	//first = 0x8a, second = 0x04 (data = PONG)
	buf := bytes.NewBuffer([]byte{})
	buf.WriteByte(byte(0x8a))
	buf.WriteByte(byte(0x04))
	buf.Write([]byte("PONG"))
	err := SendWS(ws.C, buf.Bytes())
	log.Println("send pong")
	return err
}

func (ws WS)SendClose() error{
	//correct ???
	var firstF byte = 0x88
	err := SendWS(ws.C, []byte{firstF})
	log.Println("send close")
	return err
}

//use data = [][]byte ? if len(data)=1,tail=true,else tail=false
func (ws WS)SendFrame(data [][]byte) (err error){
	buf := bytes.NewBuffer([]byte{})
	
	//init header, for continue frames, first message's header = 0x01, middle message's header = 0x00, last message's header = 0x80
	//for single frame, header = 0x81
	switch len(data){
	
	//if no data, send PONG instead
	case 0: 
		err = ws.SendPong()
		return
		
	//single frame
	case 1:
		err = buf.WriteByte(byte(0x81))
		if err != nil{
			return
		}
		err = SetupData(data[0], buf)
		if err != nil{
			return
		}
		err = SendWS(ws.C, buf.Bytes())
		
	//cotninual frames
	default:
		for i := range data{
			switch i{
			
			//first one's header
			case 0: 
				err = buf.WriteByte(byte(0x01))
				
			//last one's header
			case len(data)-1:
				err = buf.WriteByte(byte(0x80))
				
			//middles' header
			default:
				err = buf.WriteByte(byte(0x00))
			}
			if err != nil{
				return
			}
			
			//setup length and data
			err = SetupData(data[i], buf)
			if err != nil{
				return
			}
			err = SendWS(ws.C, buf.Bytes())
			if err != nil{
				return
			}
			
			//reset for next one
			buf.Reset()
		}
	}
	

	//_, err := buf.Write([]byte{})
	//_, err = buf.WriteString("")
	//err = buf.WriteByte(byte)
	//_, err = buf.WriteTo(ws.C)??

	log.Println("send message")
	return
}

//setup length and data after header
func SetupData(sData []byte, buf *bytes.Buffer) (err error){
	length := int64(len(sData))
	switch{
	
	//math.MaxInt64 = 1<<63
	case length > math.MaxInt64:
		//length = 127 first, follow by real length
		//Int64ToBytes return []byte, length = 8
		err = buf.WriteByte(byte(0x7f))
		if err != nil{
			return
		}
		_, err = buf.Write(Int64ToBytes(length))
		if err != nil{
			return
		}
		_, err = buf.Write(sData)
		if err != nil{
			return
		}
		
	case length > 1<<15:
		//length = 126 first, follow by real length, 2 bytes
		err = buf.WriteByte(byte(0x7e))
		if err != nil{
			return
		}
		_, err = buf.Write(Int64ToBytes(length)[6:])
		if err != nil{
			return
		}
		_, err = buf.Write(sData)
		if err != nil{
			return
		}
		
	default:
		//message from server to client must not masked, so 0x00 | length
		err = buf.WriteByte(byte(0x00) | Int64ToBytes(length)[7])
		if err != nil{
			return err
		}
		_, err = buf.Write(sData)
		if err != nil{
			return
		}
	}
	return
}

//send normal data(test method)
func (ws WS)SendWelcome() error{
	data := []byte{byte(0x81),byte(0x07)}
	data = append(data, []byte("welcome")...)
	err := SendWS(ws.C, data)
	log.Println("send welcome")
	return err
}