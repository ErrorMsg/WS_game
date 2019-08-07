//card info

package srv

import (
	"net"
	"log"
	"math/rand"
	"time"
)

//all cards store in db
type Card struct{
	C_id int
	C_name string
	C_pic string
	C_kind int   //0-expcard, 1-moncard, 2-equipcard, 3-freecard
}

type CardsInUse [2][]Card   //0=cover cards, 1=discard cards

//experience card, 0-1000
type ExpCard struct{
	Card
	//kind int
	C_effect string
}


//monster card, 1000-2000
type MonCard struct{
	Card
	Summary
	C_effect string
}


//equipment card, 2000-3000
type EquipCard struct{
	Card
	Summary
	Part int   //1-head,2-two hand weapon,3-one hand weapon,4-body,5-shoes,6-jewelry
	C_effect string
}

//free card, 3000-4000
type FreeCard struct{
	Card
	C_effect string
}


var moncards Cards   
var expcards Cards
var equipcards Cards
var freecards Cards
var pickcards []Card
var allcards map[int]Card

//use card to find specific card from db
func (c *Card)SpecificCard() (interface{},error){
	
}

//shuffle in one side
func (cs *Cards)shuffle() error{
	for i:=0;i<len(cs[1]);i++{
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(len(cs[1])-i)+i
		cs[1][i], cs[1][r] = cs[1][r], cs[1][i]
	}
	cs[0] = cs[1][:]
	cs[1] = nil
	return nil
}

//shuffle in two side
func (cs *Cards)shuffle2() error{
	for i:=0;i<len(cs[1]);i++{
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(len(cs[1]))
		cs[0][i] = cs[1][r]
	}
	return nil
}

func (cs *Cards)pickCard(n int) ([]Card, error){
	pickcards = nil
	for i:=0;i<n;i++{
		if len(cs) == 0{
			err := cs.shuffle()
			if err != nil{
				return pickcards,err
			}
		}
		c := cs[0]
		pickcards = append(pickcards, c)
		cs = cs[1:]
	}
	return pickcards, nil
}

func (cs *Card)exchangeCard(n map[int]int) error{
	for k,v := range n{
		if k < 0{
			k = len(cs) + k
		}else if v < 0{
			v = len(cs) + v
		}
		cs[k], cs[v] = cs[v], cs[k]
	}
	return nil
}