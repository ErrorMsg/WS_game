package main

import (
	"net"
	"log"
)

type Card struct{
	id int
	name string
	pic string
}

//experience card
type ExpCard struct{
	Card
	kind int
	effect string
}


//monster card
type MonCard struct{
	Card
	Summary
	effect string
	reword
}