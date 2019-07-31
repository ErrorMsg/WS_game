package main

import (
	"net"
	"log"
)

type Player struct{
	Position
	Summary
	Equipment
	Link
}

//player's position info in chess
type Position struct{
	P_level int
	P_block int
	P_pos int
}


//player's summary info
type Summary struct{
	S_hp int
	S_atk int
	S_def int
	S_spd int
	S_money int
	S_status int
}

//player's equipment info
type Equipment struct{
	E_head int
	E_weapon0 int
	E_weapon1 int
	E_weapon2 int
	E_body int
	E_shoes int
	E_jewelry1 int
	E_jewelry2 int
}

//player's connect info
type Link struct{
	addr net.TCPAddr
	name string
}