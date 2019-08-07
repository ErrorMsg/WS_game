//for two sides data stream grpc test
package main

import (
	"io"
	"log"
	"net"
	"strconv"
	"google.golang.org/grpc"
	proto "test"
)

//streamer server
type Streamer struct{}

func (s *Streamer) BidStream(stream proto.Test_BidStreamServer) error{
	ctx := stream.Context()
	for {
		select{
		case <- ctx.Done():   //terminat by client
			log.Println("client send terminate by context")
			return ctx.Err()
		default:   //default received from client
			rev, err := stream.Recv()
			if err == io.EOF{
				log.Println("receive from client end")
				return nil
			}
			if err != nil{
				log.Println("receive from client error")
				return err
			}
			
			switch rev.Input{
			case "ED\r\n":   //in window end with \r\n, in linux end with \n
				log.Println("received ED")
				if err := stream.Send(&proto.Response{Output: "get ED"});err!=nil{
					return err
				}
				return nil
			case "CD\r\n":   //in window end with \r\n, in linux end with \n
				log.Println("received CD")
				for i:=0;i<10;i++{
					if err := stream.Send(&proto.Response{Output: strconv.Itoa(i)});err!=nil{
						return err
					}
				}
			default:
				log.Println("received MSG ", rev.Input)
				if err := stream.Send(&proto.Response{Output: rev.Input});err!=nil{
					return err
				}
			}
		}
	}
}

func main(){
	log.Println("starting")
	server := grpc.NewServer()
	
	proto.RegisterTestServer(server, &Streamer{})
	
	addr, err := net.Listen("tcp", ":9002")
	if err != nil{
		panic(err)
	}
	if err := server.Serve(addr);err!=nil{
		panic(err)
	}
}