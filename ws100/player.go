//player info

package srv

import (
	"net"
	"log"
)

type Player struct{
	Position   //position in map
	Summary
	Equip
	Link
	HandCards []interface{}   //[]Card？
}

//player's position
type Position struct{
	P_level int
	P_block int   //blocks from start point to current point
	P_pos int   //real position
}

//player's summary info
type Summary struct{
	S_hp int
	S_maxhp int
	S_atk int
	S_def int
	S_spd int
	S_money int
	S_status int
}

//player's equipment info
type Equips struct{
	Equipments []int   //[head, weapon0, weapon1, weapon2, body, shoes, jewelry1, jewelry2]
	E_head int
	E_weapon0 int   //two hands weapon
	E_weapon1 int   //one hand weapon1
	E_weapon2 int   //one hand weapon2
	E_body int
	E_shoes int
	E_jewelry1 int
	E_jewelry2 int
}

//player's connect info
type Link struct{
	L_addr net.TCPAddr
	L_Name string
	L_ID int
}


//return player list for map update
var PlayerList map[int]Player