package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"google.golang.org/grpc"
	proto "test"
)

func main(){
	conn, err := grpc.Dial("localhost:9002", grpc.WithInsecure())
	if err != nil{
		log.Println("connect fail")
		return
	}
	defer conn.Close()
	
	client := proto.NewTestClient(conn)
	
	ctx := context.Background()
	
	stream, err := client.BidStream(ctx)
	if err != nil{
		log.Println("create bidstream fail")
	}
	
	go func(){
		log.Println(">: ")
		data := bufio.NewReader(os.Stdin)
		for {
			msg, _ := data.ReadString('\n')   //in window end with \r\n, in linux end with \n
			if err := stream.Send(&proto.Request{Input: msg});err!=nil{
				return
			}
		}
	}()
	
	for{
		rev, err := stream.Recv()
		if err == io.EOF{
			log.Println("received from server end")
			break
		}
		if err != nil{
			log.Println("received from server fail")
		}
		
		log.Println(rev.Output)
	}
}